package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var tests = []struct {
	data     []byte
	fuguFile *FuguFile
	err      bool
}{
	// set default label
	{
		[]byte(`
name: test
image: mattes/foobar
`),
		&FuguFile{
			Data: map[string]interface{}{"default": map[string]interface{}{"name": "test", "image": "mattes/foobar"}},
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
			Data: map[string]interface{}{"default": map[interface{}]interface{}{"name": "test", "image": "mattes/foobar"}},
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
			Data: map[string]interface{}{"default": map[interface{}]interface{}{"name": "test", "image": "mattes/foobar", "publish": []interface{}{"8080:80"}}},
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
			Data: map[string]interface{}{"foo": map[interface{}]interface{}{"name": "test", "image": "mattes/foobar"}, "bar": map[interface{}]interface{}{"name": "test", "image": "mattes/foobar2"}},
		},
		false,
	},
}

func TestParse(t *testing.T) {
	for _, tt := range tests {
		fuguFile, err := Parse(tt.data)
		if err != nil && tt.err == false {
			t.Fatalf("Expected no err, but got one: %v.\n%s", err, tt)
		} else if err == nil && tt.err == true {
			t.Fatalf("Expected err, but gone none. %s", tt)
		}
		assert.Equal(t, tt.fuguFile.Data, fuguFile.Data)
	}
}
