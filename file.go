package main

import (
	_ "fmt"
	"github.com/mattes/yaml"
)

var fileNamePaths = []string{"fugu.yml", "fugu.yaml", ".fugu.yml", ".fugu.yaml"}

type FuguFile struct {
	Path     string
	FileName string
	Data     map[yaml.StringIndex]interface{}
}

func Parse(data []byte) (*FuguFile, error) {
	f := &FuguFile{}
	if err := yaml.Unmarshal(data, &f.Data); err != nil {
		return nil, err
	}

	// are labels used?
	// if _, ok := f.Data["image"]; ok {
	// 	// no labels are used ...
	// 	f.Data = map[string]interface{}{"default": f.Data}
	// }

	return f, nil
}
