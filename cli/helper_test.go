package cli

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/mattes/fugu/config"
	"github.com/stretchr/testify/require"
	"testing"
)

var mergeConfigTests = []struct {
	fugufileData []byte
	args         []string
	out          []config.Value
	err          bool
}{
	{
		[]byte(`
name: test
image: mattes/foobar
`),
		[]string{"--name", "foobar", "--rm"},
		[]config.Value{
			&config.StringValue{Name: []string{"command"}, Value: "", Defined: false},
			&config.StringSliceValue{Name: []string{"args"}, Value: nil, Defined: false},
			&config.StringValue{Name: []string{"name"}, Value: "foobar", Defined: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar", Defined: true},
			&config.BoolValue{Name: []string{"rm"}, Value: true, Defined: true},
			&config.StringSliceValue{Name: []string{"publish", "p"}, Value: nil, Defined: false},
			&config.Int64Value{Name: []string{"non-exist"}, Defined: false},
		},
		false,
	},
	{
		[]byte(`
name: test
image: mattes/foobar
  `),
		[]string{"--name", "foobar", "--rm", "echo", `"hello world"`},
		[]config.Value{
			&config.StringValue{Name: []string{"command"}, Value: "echo", Defined: true},
			&config.StringSliceValue{Name: []string{"args"}, Value: []string{`"hello world"`}, Defined: true},
			&config.StringValue{Name: []string{"name"}, Value: "foobar", Defined: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar", Defined: true},
			&config.BoolValue{Name: []string{"rm"}, Value: true, Defined: true},
			&config.StringSliceValue{Name: []string{"publish", "p"}, Value: nil, Defined: false},
			&config.Int64Value{Name: []string{"non-exist"}, Defined: false},
		},
		false,
	},
}

func TestMergeConfig(t *testing.T) {
	for _, tt := range mergeConfigTests {
		c := []config.Value{
			&config.StringValue{Name: []string{"command"}},
			&config.StringSliceValue{Name: []string{"args"}},
			&config.StringValue{Name: []string{"name"}},
			&config.StringValue{Name: []string{"image"}},
			&config.BoolValue{Name: []string{"rm"}},
			&config.StringSliceValue{Name: []string{"publish", "p"}},
			&config.Int64Value{Name: []string{"non-exist"}},
		}

		err := MergeConfig(tt.fugufileData, tt.args, "", &c)
		if !tt.err {
			require.NoError(t, err, spew.Sdump(tt))
		} else if tt.err {
			require.Error(t, err, spew.Sdump(tt))
		}

		require.Equal(t, tt.out, c, spew.Sdump(tt), spew.Sdump(c))
	}
}

var buildArgsTest = []struct {
	in  []config.Value
	out []string
}{
	{
		[]config.Value{
			&config.StringValue{Name: []string{"command"}, Value: "echo", Defined: true},
			&config.StringSliceValue{Name: []string{"args"}, Value: []string{"hello", "world"}, Defined: true},
			&config.StringValue{Name: []string{"name"}, Value: "foobar", Defined: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/image", Defined: true},
			&config.BoolValue{Name: []string{"rm"}, Value: false, Defined: true},
			&config.StringSliceValue{Name: []string{"publish", "p"}, Value: []string{"50:60", "70:80", "90:100"}, Defined: true},
			&config.Int64Value{Name: []string{"non-exist"}, Value: 5, Defined: false},
		},
		[]string{`--name="foobar"`, `--rm=false`,
			`--publish="50:60" --publish="70:80" --publish="90:100"`,
			"mattes/image", "echo", "hello", "world"},
	},
}

func TestBuildArgs(t *testing.T) {
	for _, tt := range buildArgsTest {
		args := BuildArgs(&tt.in)
		require.Equal(t, tt.out, args, spew.Sdump(tt))
	}
}
