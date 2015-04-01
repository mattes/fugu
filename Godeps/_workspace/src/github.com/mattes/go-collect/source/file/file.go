package file

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mattes/go-collect/data"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrEmptyPath   = errors.New("source: file: no path given")
	ErrYamlParsing = errors.New("source: file: yaml parsing failed")
	ErrYamlLevels  = errors.New("source: file: only 2 levels of indentation allowed")
)

// File implements Source interface
type File struct {
	label string
	url   *url.URL
	path  string

	body []byte

	// TODO implement json, toml, ...
	yaml map[string]map[string][]string

	labels []string
}

func (s *File) Scheme() string {
	return "file"
}

func (s *File) ExampleUrl() string {
	return "file://config.yml"
}

func (s *File) Load(label string, u *url.URL) (*data.Data, error) {
	s.url = u
	s.setPathFromUrl()

	if err := s.readFile(); err != nil {
		return nil, err
	}
	if err := s.parse(); err != nil {
		return nil, err
	}

	s.label = s.selectLabel(label)

	return s.getData(), nil
}

func (s *File) Labels() []string {
	return s.labels
}

func (s *File) setPathFromUrl() {
	// TODO what about windows and file://paths?

	if s.url.Host == "" {
		// assume absolute path
		// file:///home/config.yml
		s.path = s.url.Path
	} else {
		// assume relative path, this is not standard conform though
		// file://config.yml
		s.path = s.url.Host + "/" + s.url.Path
		s.path, _ = filepath.Abs(s.path)
	}
}

func (s *File) readFile() error {
	if s.path == "" {
		return ErrEmptyPath
	}
	body, err := ioutil.ReadFile(s.path)
	if err != nil {
		return fmt.Errorf("source: file: %v", err.Error())
	}
	s.body = body
	return nil
}

// parse parses the file content into a yaml struct
func (s *File) parse() error {

	// inject env vars
	s.body = injectEnvVars(s.body)

	// Parse '<<:' to '<:' so we don't trigger the internal
	// yaml pkg inheritance parsing. it will fail because of
	// `map[string]map[string]interface{}`.
	// yaml: map merge requires map or sequence of maps as the value
	// We also don't want to rely on *label markers for inheritance.

	// this major wtf replaces <<: with <: in every line
	reRepl := regexp.MustCompile("^\\s*(<<:)")
	spl := bytes.Split(s.body, []byte("\n"))
	for i := 0; i < len(spl); i++ {
		spl[i] = reRepl.ReplaceAllFunc(spl[i], func(in []byte) []byte {
			// whitespace to replacement in order to allow <<:label, instead of <<: label
			return bytes.Replace(in, []byte("<<:"), []byte("<: "), 1)
		})
	}
	s.body = bytes.Join(spl, []byte("\n"))

	// try to unmarshal with labels
	hasLabels := true
	var yamlWithLabels map[string]map[string]interface{}
	if err := yaml.Unmarshal(s.body, &yamlWithLabels); err != nil {
		// try without labels
		var yamlNoLabels map[string]interface{}
		if err := yaml.Unmarshal(s.body, &yamlNoLabels); err != nil {
			return ErrYamlParsing
		}

		yamlWithLabels = map[string]map[string]interface{}{
			"default": yamlNoLabels,
		}
		hasLabels = false
	}

	// map[string]map[string]interface{} -> map[string]map[string][]string
	s.yaml = make(map[string]map[string][]string)
	for k, v := range yamlWithLabels {
		s.yaml[k] = make(map[string][]string)
		for k2, v2 := range v {
			s.yaml[k][k2] = make([]string, 0)
			switch v2.(type) {
			case []interface{}:
				s.yaml[k][k2] = interfaceSliceToStringSlice(v2.([]interface{}))

			case map[interface{}]interface{}:
				return ErrYamlLevels

			default:
				s.yaml[k][k2] = []string{fmt.Sprintf("%v", v2)}
			}
		}
	}

	// parse inheritance, do this until all referenced labels are referenced
	allParsed := false
	for !allParsed {
		allParsed = true
		// loop over all labels
		for label, vs := range s.yaml {
			// loop over all key:values in label
			for k, v := range vs {
				// look for '<' keys
				if k == "<" {
					if len(v) != 1 {
						// <<: [label, label] is not possible
						return ErrYamlParsing
					}
					// delete '<' key
					delete(s.yaml[label], k)
					// check if referenced label exists
					if refLabel, ok := s.yaml[v[0]]; ok {
						// merge key:values from referenced label into this label
						// loop over all keys:values in referenced label
						for k, v := range refLabel {
							if _, exists := vs[k]; !exists {
								s.yaml[label][k] = v
							}
							// this label references another label
							if k == "<" {
								allParsed = false
							}
						}
					}
				}
			}
		}
	}

	// get labels and their order
	if hasLabels {
		var labelOrder yaml.MapSlice
		if err := yaml.Unmarshal(s.body, &labelOrder); err != nil {
			return ErrYamlParsing
		}
		s.labels = make([]string, 0)
		for _, v := range labelOrder {
			s.labels = append(s.labels, v.Key.(string))
		}
	} else {
		s.labels = []string{"default"}
	}

	return nil
}

// labelExists returns bool if label exists in yaml
func (s *File) labelExists(label string) bool {
	for _, l := range s.labels {
		if l == label {
			return true
		}
	}
	return false
}

// selectLabel returns the label that should be used
func (s *File) selectLabel(label string) string {
	useLabel := ""

	if s.labelExists(label) {
		// use label
		useLabel = label

	} else if label != "" && !s.labelExists(label) {
		// don't return anything if specifc label was not found
		useLabel = ""

	} else if s.labelExists("default") {
		// use "default" label
		useLabel = "default"

	} else if len(s.labels) > 0 {
		// get first label found in file
		useLabel = s.labels[0]
	}

	return useLabel
}

// getYamlForLabel internally selects the right label to return
// given label (if exists), default label (if exist) or first label
func (s *File) getYamlForLabel() (p map[string][]string, ok bool) {
	for k, v := range s.yaml {
		if k == s.label {
			return v, true
		}
	}
	return nil, false
}

// Data returns param.Data for given label
func (s *File) getData() *data.Data {
	ps, ok := s.getYamlForLabel()
	if !ok {
		return data.New()
	}

	d := data.New()
	for k, v := range ps {
		d.Set(k, v...)
	}
	return d
}

// interfaceSliceToStringSlice converts:
// []interface{} -> []string
func interfaceSliceToStringSlice(in []interface{}) []string {
	out := make([]string, 0)
	for _, v := range in {
		out = append(out, fmt.Sprintf("%v", v))
	}
	return out
}

// injectEnvVars replaces $ENV vars with their actual value
func injectEnvVars(value []byte) []byte {
	envVarRegex := regexp.MustCompile(`\$[a-zA-Z_]+[a-zA-Z0-9_]*`)
	return envVarRegex.ReplaceAllFunc(value, func(match []byte) []byte {
		return []byte(os.Getenv(string(bytes.TrimPrefix(match, []byte("$")))))
	})
}
