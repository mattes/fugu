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
	usesLabels := true
	newFlatData := make(map[interface{}]interface{}) // as per yaml pkg default
	for k, v := range f.Data {
		newFlatData[k.Value] = v
		if k.Value == "image" {
			// found image variable in level 1, thus no labels are used
			usesLabels = false
		}
	}
	if !usesLabels {
		// set default label
		f.Data = map[yaml.StringIndex]interface{}{yaml.StringIndex{Index: 0, Value: "default"}: newFlatData}
	}

	return f, nil
}
