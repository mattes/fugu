package main

import (
	"gopkg.in/yaml.v1"
)

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

	// parse labels

	return f, nil
}
