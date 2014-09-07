package main

import (
	"gopkg.in/yaml.v1"
)

var fileNamePaths = []string{"fugu.yml", "fugu.yaml", ".fugu.yml", ".fugu.yaml"}

type FuguFile struct {
	Path     string
	FileName string
	Data     map[string]interface{}
}

func Parse(data []byte) (*FuguFile, error) {
	f := &FuguFile{}
	if err := yaml.Unmarshal(data, &f.Data); err != nil {
		return nil, err
	}

	// are labels used?
	if _, ok := f.Data["image"]; ok {
		// no labels are used ...
		f.Data = map[string]interface{}{"default": f.Data}
	}

	return f, nil
}
