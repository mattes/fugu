package cli

import (
	"github.com/mattes/fugu/config"
	"github.com/mattes/fugu/docker"
	"github.com/mattes/fugu/file"
	"io/ioutil"
	"os"
	"sort"
)

func MergeConfig(fugufileData []byte, args []string, label string, conf *[]config.Value) error {
	// first: parse fugufile
	_, err := file.Load(fugufileData, label, conf)
	if err != nil {
		return err
	}

	// second: parse docker args
	err = docker.Load(args, conf)
	if err != nil {
		return err
	}

	return nil
}

func BuildRunArgs(conf *[]config.Value) []string {
	dockerImage := ""
	dockerCommand := ""
	dockerArgs := make([]string, 0)

	args := make([]string, 0)
	for _, c := range *conf {
		if c.Names()[0] == "image" {
			dockerImage = c.Get().(string)
		} else if c.Names()[0] == "command" {
			dockerCommand = c.Get().(string)
		} else if c.Names()[0] == "args" {
			dockerArgs = c.Get().([]string)
		} else {
			v := c.Arg()
			if len(v) > 0 {
				args = append(args, c.Arg()...)
			}
		}
	}

	sort.Sort(sort.StringSlice(args))

	if dockerImage != "" {
		// check if dockerImage != "" although this should never happen!
		args = append(args, dockerImage)
	}
	if dockerCommand != "" {
		args = append(args, dockerCommand)
	}
	if len(dockerArgs) > 0 {
		args = append(args, dockerArgs...)
	}

	return args
}

func FindFugufile(filepath string, searchFilePaths []string) (fugufilePath string, viaFugufilePath bool) {
	// does the file exist?
	if _, err := os.Stat(filepath); err == nil {
		return filepath, false
	}

	// filepath does not exist.
	// let's try our searchFilePaths
	for _, f := range searchFilePaths {
		if _, err := os.Stat(f); err == nil {
			return f, false
		}
	}

	// no file found in searchFilePaths
	return "", false
}

// GetAllLabels reads all labels from a given yaml file
func GetAllLabels(filepath string) ([]string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	c := []config.Value{}
	return file.Load(data, "", &c)
}
