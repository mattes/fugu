package docker

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/mattes/fugu/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var loadTests = []struct {
	in  []string
	out []config.Value
	err bool
}{
	{
		[]string{"--name", "test", "--image", "mattes/foobar", "--publish", "8080:80", "-p", "55:66"},
		[]config.Value{
			&config.StringValue{Name: []string{"name"}, Value: "test", Defined: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar", Defined: true},
			&config.StringSliceValue{Name: []string{"publish", "p"}, Value: []string{"8080:80", "55:66"}, Defined: true},
			&config.StringValue{Name: []string{"non-exist"}, Defined: false},
		},
		false,
	},
	{
		[]string{"--name", "test", "--image", "mattes/foobar", "-p", "55:66"},
		[]config.Value{
			&config.StringValue{Name: []string{"name"}, Value: "test", Defined: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar", Defined: true},
			&config.StringSliceValue{Name: []string{"publish", "p"}, Value: []string{"55:66"}, Defined: true},
			&config.StringValue{Name: []string{"non-exist"}, Defined: false},
		},
		false,
	},
}

func TestLoad(t *testing.T) {
	for _, tt := range loadTests {
		c := []config.Value{
			&config.StringValue{Name: []string{"name"}},
			&config.StringValue{Name: []string{"image"}},
			&config.StringSliceValue{Name: []string{"publish", "p"}},
			&config.StringValue{Name: []string{"non-exist"}},
		}

		err := Load(tt.in, &c)
		if !tt.err {
			require.NoError(t, err, spew.Sdump(tt))
		} else if tt.err {
			require.Error(t, err, spew.Sdump(tt))
		}
		require.Equal(t, tt.out, c, spew.Sdump(tt), spew.Sdump(c))
	}
}

var flagWasDefinedTests = []struct {
	findName []string
	args     []string
	defined  bool
}{
	{
		[]string{"publish", "p"},
		[]string{"--bar", "--publish", "--foobar", "-foo"},
		true,
	},
	{
		[]string{"publish", "p"},
		[]string{"--bar", "-p", "--foobar", "-foo"},
		true,
	},
	{
		[]string{"publish", "p"},
		[]string{"--bar", "--foobar", "-foo"},
		false,
	},
	{
		[]string{"publish", "p"},
		[]string{"--bar", "--p", "--foobar", "-foo"},
		false,
	},
	{
		[]string{"publish", "p"},
		[]string{"--bar", "-publish", "--foobar", "-foo"},
		false,
	},
	{
		[]string{"rm"},
		[]string{"--bar", "-rm", "--foobar", "-foo"},
		false,
	},
	{
		[]string{"rm"},
		[]string{"--bar", "--rm", "--foobar", "-foo"},
		true,
	},
	{
		[]string{"rm"},
		[]string{"--bar", "--rm=true", "--foobar", "-foo"},
		true,
	},
	{
		[]string{"rm"},
		[]string{"--bar", "--rm=false", "--foobar", "-foo"},
		true,
	},
}

func TestFlagWasDefined(t *testing.T) {
	for _, tt := range flagWasDefinedTests {
		assert.Equal(t, tt.defined, flagWasDefined(tt.findName, tt.args), spew.Sdump(tt))
	}
}
