package file

import (
	"errors"
	"fmt"
	"github.com/mattes/yaml"
	"io/ioutil"
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
		return &FuguFile{}, err
	}

	if len(f.Data) == 0 {
		return &FuguFile{}, errors.New("Empy configuration given.")
	}

	// are labels used?
	imageFoundInLevel1 := false
	newFlatData := make(map[interface{}]interface{}) // as per yaml pkg default
	for k, v := range f.Data {
		newFlatData[k.Value] = v
		if k.Value == "image" {
			// found image variable in level 1, thus no labels are used
			imageFoundInLevel1 = true
		}
	}
	if imageFoundInLevel1 {
		// set default label, because no labels are used
		f.Data = map[yaml.StringIndex]interface{}{yaml.StringIndex{Index: 0, Value: "default"}: newFlatData}
	}

	if !imageFoundInLevel1 {
		// make sure all labels have a image variable
		for k, v := range f.Data {
			switch v.(type) {

			case string:
				return &FuguFile{}, errors.New(fmt.Sprintf("Missing 'image' variable for label %s.", k.Value))

			case map[interface{}]interface{}:
				imageFoundInLevel2 := false
				for k2, _ := range v.(map[interface{}]interface{}) {
					if k2.(string) == "image" {
						imageFoundInLevel2 = true
						break
					}
				}
				if !imageFoundInLevel2 {
					return &FuguFile{}, errors.New(fmt.Sprintf("Missing 'image' variable for label %s.", k.Value))
				}

			default:
				return &FuguFile{}, errors.New(fmt.Sprintf("Unable to parse configuration for label %s.", k.Value))
			}
		}
	}

	return f, nil
}

func Find(dir string) (filepath string, err error) {
	return "", nil
}

func GetConfig(filepath, label string) (config map[string]interface{}, err error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	fugufile, err := Parse(content)
	if err != nil {
		return nil, err
	}

	var foundConfig map[interface{}]interface{}
	if label == "" {
		// get first label from fugufile
		for k, v := range fugufile.Data {
			if k.Index == 0 {
				foundConfig = v.(map[interface{}]interface{})
				break
			}
		}
	} else {
		// get the actual label
		for k, v := range fugufile.Data {
			if k.Value == label {
				foundConfig = v.(map[interface{}]interface{})
				break
			}
		}
	}

	if foundConfig == nil {
		if label == "" {
			return nil, errors.New("No configuration found.")
		} else {
			return nil, errors.New(fmt.Sprintf("No configuration found for label %s.", label))
		}
	}

	// map[interface{}]interface{} -> map[string]interface{}
	config = make(map[string]interface{})
	for k, v := range foundConfig {
		config[k.(string)] = v
	}

	return config, nil
}
