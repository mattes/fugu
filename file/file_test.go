package file

import (
	"fmt"
	"github.com/mattes/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var parseTests = []struct {
	data     []byte
	fuguFile *FuguFile
	err      bool
}{

	// test if default label is set if none is present
	{
		[]byte(`
name: test
image: mattes/foobar
`),
		&FuguFile{
			Data: map[yaml.StringIndex]interface{}{yaml.StringIndex{"default", 0}: map[interface{}]interface{}{"name": "test", "image": "mattes/foobar"}},
		},
		false,
	},

	// test labels
	{
		[]byte(`
default:
  name: test
  image: mattes/foobar
`),
		&FuguFile{
			Data: map[yaml.StringIndex]interface{}{yaml.StringIndex{"default", 0}: map[interface{}]interface{}{"name": "test", "image": "mattes/foobar"}},
		},
		false,
	},

	// test maps
	{
		[]byte(`
default:
  name: test
  image: mattes/foobar
  publish:
    - 8080:80
`),
		&FuguFile{
			Data: map[yaml.StringIndex]interface{}{yaml.StringIndex{"default", 0}: map[interface{}]interface{}{"name": "test", "image": "mattes/foobar", "publish": []interface{}{"8080:80"}}},
		},
		false,
	},

	// test inheritance
	{
		[]byte(`
foo: &foo
  name: test
  image: mattes/foobar

bar:
  <<: *foo
  image: mattes/foobar2
`),
		&FuguFile{
			Data: map[yaml.StringIndex]interface{}{yaml.StringIndex{"foo", 0}: map[interface{}]interface{}{"name": "test", "image": "mattes/foobar"}, yaml.StringIndex{"bar", 1}: map[interface{}]interface{}{"name": "test", "image": "mattes/foobar2"}},
		},
		false,
	},

	// test missing image
	{
		[]byte(`
name: test
`),
		&FuguFile{},
		true,
	},

	// test image only
	{
		[]byte(`
image: test
`),
		&FuguFile{
			Data: map[yaml.StringIndex]interface{}{yaml.StringIndex{"default", 0}: map[interface{}]interface{}{"image": "test"}},
		},
		false,
	},

	// 	test missing image when using labels
	{
		[]byte(`
default:
  name: test
`),
		&FuguFile{},
		true,
	},

	// test image only when using labels
	{
		[]byte(`
label1:
  image: test
`),
		&FuguFile{
			Data: map[yaml.StringIndex]interface{}{yaml.StringIndex{"label1", 0}: map[interface{}]interface{}{"image": "test"}},
		},
		false,
	},

	// test empy
	{
		[]byte(""),
		&FuguFile{},
		true,
	},

	// test bullshit
	{
		[]byte(`label1:`),
		&FuguFile{},
		true,
	},
	{
		[]byte(`default`),
		&FuguFile{},
		true,
	},
	{
		[]byte(`default:`),
		&FuguFile{},
		true,
	},
	{
		[]byte(`
default:
  - hmm
`),
		&FuguFile{},
		true,
	},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		fuguFile, err := Parse(tt.data)
		if !tt.err {
			require.NoError(t, err, fmt.Sprintf("%s", tt))
		} else if tt.err {
			require.Error(t, err, fmt.Sprintf("%s", tt))
		}
		assert.Equal(t, tt.fuguFile.Data, fuguFile.Data)
	}
}

var getConfigTests = []struct {
	filepath string
	label    string
	config   map[string]interface{}
	err      bool
}{
	{
		"../fugu.example.yml",
		"",
		map[string]interface{}{"image": "a-team/action", "detach": true},
		false,
	},
	{
		"../fugu.example.yml",
		"a-team",
		map[string]interface{}{"image": "a-team/action", "detach": true},
		false,
	},
	{
		"../fugu.example.yml",
		"hannibal",
		map[string]interface{}{"image": "a-team/action", "detach": false, "command": "echo", "args": []interface{}{"I love it when a plan comes together."}, "rm": true},
		false,
	},
}

func TestGetConfig(t *testing.T) {
	for _, tt := range getConfigTests {
		config, err := GetConfig(tt.filepath, tt.label)
		if !tt.err {
			require.NoError(t, err, fmt.Sprintf("%s", tt))
		} else if tt.err {
			require.Error(t, err, fmt.Sprintf("%s", tt))
		}
		assert.Equal(t, tt.config, config)
	}
}
