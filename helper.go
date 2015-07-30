package fugu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/github/hub/github"
	"github.com/mattes/go-collect/data"
	"github.com/mattes/go-collect/flags"
	"gopkg.in/mattes/go-expand-tilde.v1"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
)

// DockerExec runs a docker command
func DockerExec(cmd string, printCmd bool) {
	if printCmd {
		fmt.Println(cmd)
	}
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	if err := c.Run(); err != nil {
		os.Exit(2)
	}
	os.Exit(0)
}

func filterDockerFlags(p *data.Data, command string) (*data.Data, error) {
	df, err := DockerFlags[command].Keys()
	if err != nil {
		return nil, err
	}

	return data.Filter(p, df), nil
}

// buildDockerStr builds the docker command with all options and args
func buildDockerStr(command string, p *data.Data, args ...string) string {
	str := []string{}
	for _, n := range p.Keys() {
		for _, o := range p.GetAll(n) {

			o, err := tilde.Expand(o)
			if err != nil {
				fmt.Println(err)
			}

			nice := strings.TrimSpace(flags.Nice(n, o))
			if nice != "" {
				str = append(str, nice)
			}
		}
	}
	sort.Sort(sort.StringSlice(str))

	str = append([]string{"docker", command}, str...)
	str = append(str, args...)
	return strings.Join(str, " ")
}

func currentGitBranch() (shortname string, err error) {
	if os.Getenv("GOTEST") != "" {
		return "current-branch", nil
	}

	localRepo, err := github.LocalRepo()
	if err != nil {
		return "", err
	}
	branch, err := localRepo.CurrentBranch()
	if err != nil {
		return "", err
	}
	if branch.ShortName() == "" {
		return "", errors.New("empty branch name")
	}
	return branch.ShortName(), nil
}

type RegistryDockerImage struct {
	Name string
	Tags []string
}

type RegistryDockerImages []RegistryDockerImage

func (a RegistryDockerImages) Len() int           { return len(a) }
func (a RegistryDockerImages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RegistryDockerImages) Less(i, j int) bool { return a[i].Name < a[j].Name }

func ListImages(registry, user, password string) ([]RegistryDockerImage, error) {
	// # see https://docs.docker.com/reference/api/registry_api/
	// TODO maybe use github.com/docker/docker/registry instead

	var data struct {
		Results []struct {
			Name        string
			Description string
		}
	}
	// TODO allow http:// too?
	if err := registryJsonRequest("GET", "https://"+registry+"/v1/search", user, password, &data); err != nil {
		return nil, errors.New("reading from registry failed")
	}

	result := make(RegistryDockerImages, 0)
	var wg sync.WaitGroup
	for _, i := range data.Results {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			r := RegistryDockerImage{}
			r.Name = name
			var tagData map[string]string
			if err := registryJsonRequest("GET", "https://"+registry+"/v1/repositories/"+name+"/tags", user, password, &tagData); err != nil {
				return // TODO handle this error?
			}
			for tag, _ := range tagData {
				r.Tags = append(r.Tags, tag)
			}
			sort.Sort(sort.StringSlice(r.Tags))
			result = append(result, r)
		}(i.Name)
	}
	wg.Wait()
	sort.Sort(result)
	return result, nil
}

func registryJsonRequest(method, url, username, password string, data interface{}) error {
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return json.Unmarshal(body, &data)
}
