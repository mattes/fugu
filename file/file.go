// Package file parses YAML files and passes configuration into config.Values
package file

import (
	"fmt"
	"github.com/mattes/fugu/config"
	"github.com/mattes/yaml"
)

var (
	ErrInvalidYaml = fmt.Errorf("Invalid YAML format. Did you set an `image` variable?")
)

type Label struct {
	Name   string
	Config map[string]interface{}
}

func Load(data []byte, label string, conf *[]config.Value) (allLabelNames []string, err error) {
	if conf == nil {
		panic("Provide conf *[]config.Value")
	}

	if label == "" {
		label = "default"
	}

	labels, err := parse(data)
	if err != nil {
		return nil, err
	}

	if len(labels) > 0 {
		useLabel := labels[0] // use first label per default
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

	// TODO take care of yaml map[] order
	// we use a fork atm, see https://github.com/mattes/yaml
	// this fork includes a quick work-around to keep track of
	// the order of the labels in the yaml file. (see yaml.StringIndex type)
	//
	// let's use this until https://github.com/go-yaml/yaml/issues/30 is done

	// test if labels are used
	var labels bool
	preConfig := make(map[yaml.StringIndex]interface{})
	if err := yaml.Unmarshal(data, &preConfig); err != nil {
		return config, err
	}

	if len(preConfig) == 0 {
		return config, ErrInvalidYaml
	}

	// range over yaml.StringIndex to check if image tag is in top-level
	labels = true
	for k, _ := range preConfig {
		if k.Value == "image" {
			labels = false
			break
		}
	}

	// align config depending on labels
	if !labels {
		// when no labels are used, prepend `default` label

		// but first convert from map[yaml.StringIndex]interface{} to map[string]interface{}
		newPreConf := make(map[string]interface{})
		for k, v := range preConfig {
			newPreConf[k.Value] = v
		}

		config = append(config, Label{
			Name:   "default",
			Config: newPreConf, // is map[string]interface{}
		})
	} else {

		// labels are used, convert some types
		preConfigIndex := 0
		for _, _ = range preConfig {
			for label, v := range preConfig {
				if label.Index == preConfigIndex {
					preConfigIndex += 1

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
						Name:   label.Value,
						Config: vNew, // is map[string]interface{}
					})
				}
			}
		}
	}

	return config, nil
}
