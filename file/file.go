// Package file parses YAML files and passes configuration into config.Values
package file

import (
	"errors"
	"github.com/mattes/fugu/config"
	"github.com/mattes/yaml"
	"io/ioutil"
)

var (
	ErrInvalidYaml = errors.New("Invalid YAML format. Did you set an `image` variable?")
)

type Label struct {
	Name   string
	Config map[string]interface{}
}

func LoadFile(filepath, label string, conf *[]config.Value) (allLabelNames []string, err error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return Load(data, label, conf)
}

func Load(data []byte, label string, conf *[]config.Value) (allLabelNames []string, err error) {
	if label == "" {
		label = "default"
	}

	labels, err := parse(data)
	if err != nil {
		return nil, err
	}

	if len(labels) > 0 {

		// TODO(mattes): this is buggy, because we cannot garantuee the sort order
		// see https://github.com/go-yaml/yaml/issues/30

		useLabel := labels[0]
		for _, l := range labels {
			allLabelNames = append(allLabelNames, l.Name)
			if label == l.Name {
				useLabel = l
			}
		}

		// populate Label.Config into config.Value
		for _, c := range *conf {
			for _, name := range c.Names() {
				for name2, val := range useLabel.Config {
					if name == name2 {
						c.Set(val)
					}
				}
			}
		}
		return allLabelNames, nil

	} else {
		return nil, nil
	}
}

func parse(data []byte) ([]Label, error) {
	config := make([]Label, 0)

	// test if labels are used
	var labels bool
	preConfig := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &preConfig); err != nil {
		return config, err
	}

	if len(preConfig) == 0 {
		return config, ErrInvalidYaml
	}

	if _, ok := preConfig["image"]; ok {
		labels = false
	} else {
		labels = true
	}

	// align config depending on labels
	if !labels {
		// when no labels are used, prepend `default` label
		config = append(config, Label{
			Name:   "default",
			Config: preConfig, // is map[string]interface{}
		})
	} else {
		// labels are used, convert some types

		// TODO(mattes): see comment from above
		// preConfig should reflect the order from the yaml file itself

		for label, v := range preConfig {
			vNew := make(map[string]interface{})
			switch v.(type) {
			case map[interface{}]interface{}:
				imageFound := false
				for k2, v2 := range v.(map[interface{}]interface{}) {
					vNew[k2.(string)] = v2
					if k2.(string) == "image" {
						imageFound = true
					}
				}
				if !imageFound {
					return config, ErrInvalidYaml
				}
			default:
				return config, ErrInvalidYaml
			}
			config = append(config, Label{
				Name:   label,
				Config: vNew, // is map[string]interface{}
			})
		}
	}

	return config, nil
}
