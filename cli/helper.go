package cli

import (
	"github.com/mattes/fugu/config"
	"github.com/mattes/fugu/docker"
	"github.com/mattes/fugu/file"
)

func MergeConfig(fugufileData []byte, args []string, label string, conf *[]config.Value) error {
	// first: get fugufile
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

func BuildArgs(conf *[]config.Value) []string {
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
				args = append(args, c.Arg())
			}
		}
	}

	if dockerImage != "" {
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

func FindFugufile(searchFiles []string) (string, error) {

	return "", nil
}
