package fugu

import (
	"errors"
	"fmt"
	"github.com/docker/docker/registry"
	"github.com/howeyc/gopass"
	"github.com/mattes/go-collect"
	"github.com/mattes/go-collect/data"
	"gopkg.in/mattes/go-expand-tilde.v1"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
)

// DockerCommands return string that is always run
var DockerCommands = make(map[string]func(c *collect.Collector, p *data.Data, args []string) (str string, err error))

// Commands have more freedom and usually print to stdout/stderr directly
var Commands = make(map[string]func(c *collect.Collector, p *data.Data, args []string) (err error))

var (
	ErrTooManyArgs   = errors.New("too many arguments given")
	ErrMissingImage  = errors.New("image option is missing")
	ErrMissingName   = errors.New("name option is missing")
	ErrUnknownLabel  = errors.New("unknown label")
	ErrTagGitBranch  = errors.New("tag-git-branch failed")
	ErrMissingFlag   = errors.New("missing required flag")
	ErrNoCredentials = errors.New("missing required credentials")
)

func init() {

	DockerCommands["build"] = func(c *collect.Collector, p *data.Data, args []string) (str string, err error) {
		if p.Get("image") == "" {
			return "", ErrMissingImage
		}

		if p.IsTrue("tag-git-branch") {
			branchName, err := currentGitBranch()
			if err != nil {
				return "", ErrTagGitBranch
			}
			p.Set("tag", branchName)
		}

		if p.Get("tag") != "" {
			p.Set("tag", p.Get("image")+":"+p.Get("tag"))
		} else {
			p.Set("tag", p.Get("image"))
		}

		path := "."
		if pathh := p.Get("path"); pathh != "" {
			path = pathh
		}
		if url := p.Get("url"); url != "" {
			path = url
		}

		if len(args) == 1 {
			path = args[0]
		}

		if len(args) > 1 {
			return "", ErrTooManyArgs
		}

		pf, err := filterDockerFlags(p, "build")
		if err != nil {
			return "", err
		}

		return buildDockerStr("build", pf, path), nil
	}

	DockerCommands["run"] = func(c *collect.Collector, p *data.Data, args []string) (str string, err error) {
		if p.Get("image") == "" {
			return "", ErrMissingImage
		}

		dockerArgCommand := p.Get("command")
		if len(args) > 0 {
			dockerArgCommand = args[0]
			args = args[1:]
		}

		dockerArgArgs := p.GetAll("arg")

		if len(args) > 0 {
			dockerArgArgs = args
			args = args[:]
		}

		nargs := []string{p.Get("image")}
		if dockerArgCommand != "" {
			nargs = append(nargs, dockerArgCommand)
		}
		if len(dockerArgArgs) > 0 {
			nargs = append(nargs, dockerArgArgs...)
		}

		pf, err := filterDockerFlags(p, "run")
		if err != nil {
			return "", err
		}

		return buildDockerStr("run", pf, nargs...), nil
	}

	DockerCommands["exec"] = func(c *collect.Collector, p *data.Data, args []string) (str string, err error) {
		if p.Get("name") == "" {
			return "", ErrMissingName
		}

		dockerArgCommand := p.Get("command")
		if len(args) > 0 {
			dockerArgCommand = args[0]
			args = args[1:]
		}

		dockerArgArgs := p.GetAll("arg")
		if len(args) > 0 {
			dockerArgArgs = args
			args = args[:]
		}

		nargs := []string{p.Get("name")}
		if dockerArgCommand != "" {
			nargs = append(nargs, dockerArgCommand)
		}
		if len(dockerArgArgs) > 0 {
			nargs = append(nargs, dockerArgArgs...)
		}

		pf, err := filterDockerFlags(p, "exec")
		if err != nil {
			return "", err
		}

		return buildDockerStr("exec", pf, nargs...), nil
	}

	DockerCommands["shell"] = func(c *collect.Collector, p *data.Data, args []string) (str string, err error) {

		p.SetTrue("interactive")
		p.SetTrue("tty")
		p.SetFalse("detach")

		name := p.Get("name")
		if name == "" {
			return "", ErrMissingName
		}

		if !p.Exists("shell") {
			p.Set("shell", "/bin/bash")
		}

		pf, err := filterDockerFlags(p, "exec")
		if err != nil {
			return "", err
		}

		return buildDockerStr("exec", pf, name, p.Get("shell")), nil
	}

	DockerCommands["destroy"] = func(c *collect.Collector, p *data.Data, args []string) (str string, err error) {
		if p.Get("name") == "" {
			return "", ErrMissingName
		}

		if len(args) > 0 {
			return "", ErrTooManyArgs
		}

		return "docker rm -f " + p.Get("name"), nil
	}

	DockerCommands["push"] = func(c *collect.Collector, p *data.Data, args []string) (str string, err error) {
		if p.Get("image") == "" {
			return "", ErrMissingImage
		}

		tag := ""
		if tagg := p.Get("tag"); tagg != "" {
			tag = tagg
		}
		if len(args) > 0 {
			tag = args[0]
			args = args[1:]
		}

		if len(args) > 0 {
			return "", ErrTooManyArgs
		}

		image := p.Get("image")
		if tag != "" {
			image += ":" + tag
		}

		pf, err := filterDockerFlags(p, "push")
		if err != nil {
			return "", err
		}

		return buildDockerStr("push", pf, image), nil
	}

	DockerCommands["pull"] = func(c *collect.Collector, p *data.Data, args []string) (str string, err error) {
		if p.Get("image") == "" {
			return "", ErrMissingImage
		}

		tag := ""
		if tagg := p.Get("tag"); tagg != "" {
			tag = tagg
		}
		if len(args) > 0 {
			tag = args[0]
			args = args[1:]
		}

		if len(args) > 0 {
			return "", ErrTooManyArgs
		}

		image := p.Get("image")
		if tag != "" {
			image += ":" + tag
		}

		pf, err := filterDockerFlags(p, "pull")
		if err != nil {
			return "", err
		}

		return buildDockerStr("pull", pf, image), nil
	}

	Commands["show-data"] = func(c *collect.Collector, p *data.Data, args []string) error {
		if len(args) > 1 {
			return ErrTooManyArgs
		}

		out, err := yaml.Marshal(p.RawEnhanced())
		if err != nil {
			return err
		}
		outStr := strings.TrimSpace(fmt.Sprintf("%s", out))
		if outStr != "{}" {
			fmt.Println(outStr)
		}
		return nil
	}

	Commands["show-labels"] = func(c *collect.Collector, p *data.Data, args []string) error {
		if c.Label() != "" || len(args) > 0 {
			return ErrTooManyArgs
		}
		labels := c.Labels()
		if len(labels) > 0 {
			sort.Sort(sort.StringSlice(labels))
			for _, l := range labels {
				fmt.Println(l)
			}
		}
		return nil
	}

	Commands["images"] = func(c *collect.Collector, p *data.Data, args []string) error {
		registryStr := ""
		if len(args) > 0 {
			registryStr = args[0]
			args = args[1:]
		}

		if len(args) > 0 {
			return ErrTooManyArgs
		}

		// show local docker images
		if registryStr == "" {
			if os.Getenv("GOTEST") != "" {
				fmt.Println("<local docker images shown>")
			} else {
				DockerExec("docker images", false)
			}
			return nil
		}

		// show docker images in other registryStr ...
		if (p.Exists("user") || p.Exists("password") || p.Exists("password-stdin")) && p.Exists("file") {
			return ErrTooManyArgs
		}
		if (p.Exists("password") || p.Exists("password-stdin")) && !p.Exists("user") {
			return ErrMissingFlag
		}

		user := ""
		if p.Exists("user") {
			user = p.Get("user")
		}

		password := ""
		if p.Exists("password") {
			password = p.Get("password")
		} else if p.Exists("password-stdin") {
			fmt.Print("Password for '" + registryStr + "': ")
			password = string(gopass.GetPasswd())
		}

		if !p.Exists("user") {
			// read credentials from file
			file := p.Get("file")
			if file == "" {
				file = "~/.dockercfg"
			}

			rootPath, err := tilde.Expand(file)
			if err != nil {
				return fmt.Errorf("images: %v", err.Error())
			}

			// TODO docker.registry sets CONFIGFILE as const
			// so we can't change it, see
			// https://github.com/docker/docker/blob/v1.5.0/registry/auth.go#L22
			// using at least path information from file flag as a work-around
			rootPath = filepath.Dir(rootPath)

			allConfig, err := registry.LoadConfig(rootPath)
			if err != nil {
				return fmt.Errorf("images: %v", err.Error())
			}

			if config, ok := allConfig.Configs[registryStr]; ok {
				user = config.Username
				password = config.Password
			}
		}

		if user == "" || password == "" {
			return ErrNoCredentials
		}

		images, err := ListImages(registryStr, user, password)
		if err != nil {
			return fmt.Errorf("images: %v", err.Error())
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)
		for _, v := range images {
			fmt.Fprintf(w, "%v\t%v\n", v.Name, strings.Join(v.Tags, ", "))
		}
		w.Flush()

		return nil
	}
}
