package file

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/mattes/fugu/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var parseTests = []struct {
	in  []byte
	out []Label
	err bool
}{

	// test if default label is set if none is defined
	{
		[]byte(`
name: test
image: mattes/foobar
`),
		[]Label{
			Label{
				Name: "default",
				Config: map[string]interface{}{
					"name":  "test",
					"image": "mattes/foobar",
				},
			},
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
		[]Label{
			Label{
				Name: "default",
				Config: map[string]interface{}{
					"name":  "test",
					"image": "mattes/foobar",
				},
			},
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
		[]Label{
			Label{
				Name: "default",
				Config: map[string]interface{}{
					"name":    "test",
					"image":   "mattes/foobar",
					"publish": []interface{}{"8080:80"},
				},
			},
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
		[]Label{
			Label{
				Name: "foo",
				Config: map[string]interface{}{
					"name":  "test",
					"image": "mattes/foobar",
				},
			},
			Label{
				Name: "bar",
				Config: map[string]interface{}{
					"name":  "test",
					"image": "mattes/foobar2",
				},
			},
		},
		false,
	},

	// test missing image
	{
		[]byte(`
name: test
`),
		[]Label{},
		true,
	},

	// test image only
	{
		[]byte(`
image: test
`),
		[]Label{
			Label{
				Name: "default",
				Config: map[string]interface{}{
					"image": "test",
				},
			},
		},
		false,
	},

	// 	test missing image when using labels
	{
		[]byte(`
default:
  name: test
`),
		[]Label{},
		true,
	},

	// test when only image variable is given in label
	{
		[]byte(`
label1:
  image: test
`),
		[]Label{
			Label{
				Name: "label1",
				Config: map[string]interface{}{
					"image": "test",
				},
			},
		},
		false,
	},

	// test empy
	{
		[]byte(""),
		[]Label{},
		true,
	},

	{
		[]byte(`
image: foo:bar
`),
		[]Label{
			Label{
				Name: "default",
				Config: map[string]interface{}{
					"image": "foo:bar",
				},
			},
		},
		false,
	},

	// test non-sense
	{
		[]byte(`label1:`),
		[]Label{},
		true,
	},
	{
		[]byte(`default`),
		[]Label{},
		true,
	},
	{
		[]byte(`default:`),
		[]Label{},
		true,
	},
	{
		[]byte(`
	default:
	  - hmm
	`),
		[]Label{},
		true,
	},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		labels, err := parse(tt.in)
		if !tt.err {
			require.NoError(t, err, spew.Sdump(tt))
		} else if tt.err {
			require.Error(t, err, spew.Sdump(tt))
		}
		assert.Equal(t, tt.out, labels)
	}
}

var labelData1 = []byte(`
label:
  name: test
  image: mattes/foobar
  publish:
    - 8080:80

another-label:
  name: halligalli
  image: mattes/foobar2
  publish:
    - 55:66
`)

var loadTests = []struct {
	in            []byte
	label         string
	allLabelNames []string
	out           []config.Value
	err           bool
}{
	{
		[]byte(`
default:
  name: test
  image: mattes/foobar
  publish:
    - 8080:80
`),
		"default",
		[]string{"default"},
		[]config.Value{
			&config.StringValue{Name: []string{"name"}, Value: "test", Defined: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar", Defined: true},
			&config.StringSliceValue{Name: []string{"publish"}, Value: []string{"8080:80"}, Defined: true},
			&config.StringValue{Name: []string{"non-exist"}, Defined: false},
		},
		false,
	},

	{
		labelData1,
		"label",
		[]string{"label", "another-label"},
		[]config.Value{
			&config.StringValue{Name: []string{"name"}, Value: "test", Defined: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar", Defined: true},
			&config.StringSliceValue{Name: []string{"publish"}, Value: []string{"8080:80"}, Defined: true},
			&config.StringValue{Name: []string{"non-exist"}, Defined: false},
		},
		false,
	},

	{
		labelData1,
		"another-label",
		[]string{"label", "another-label"},
		[]config.Value{
			&config.StringValue{Name: []string{"name"}, Value: "halligalli", Defined: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar2", Defined: true},
			&config.StringSliceValue{Name: []string{"publish"}, Value: []string{"55:66"}, Defined: true},
			&config.StringValue{Name: []string{"non-exist"}, Defined: false},
		},
		false,
	},

	// test if first label is fetched
	{
		labelData1,
		"",
		[]string{"label", "another-label"},
		[]config.Value{
			&config.StringValue{Name: []string{"name"}, Value: "test", Defined: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar", Defined: true},
			&config.StringSliceValue{Name: []string{"publish"}, Value: []string{"8080:80"}, Defined: true},
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
			&config.StringSliceValue{Name: []string{"publish"}},
			&config.StringValue{Name: []string{"non-exist"}},
		}

		allLabelNames, err := Load(tt.in, tt.label, &c)
		if !tt.err {
			require.NoError(t, err, spew.Sdump(tt))
		} else if tt.err {
			require.Error(t, err, spew.Sdump(tt))
		}
		require.Equal(t, tt.out, c, spew.Sdump(tt), spew.Sdump(c))
		require.Equal(t, tt.allLabelNames, allLabelNames, spew.Sdump(tt))
	}
}

var injectEnvVarsTests = []struct {
	in  []byte
	out []byte
}{
	{
		[]byte("$FUGU_TEST123"),
		[]byte("ok"),
	},
	{
		[]byte("foo $FUGU_TEST123 bar"),
		[]byte("foo ok bar"),
	},
	{
		[]byte("FUGU_TEST123"),
		[]byte("FUGU_TEST123"),
	},
	{
		[]byte("$fugu_test123"),
		nil,
	},
	{
		[]byte("$NON_EXISTING_ENV_VAR_123456789"),
		nil,
	},
}

func TestInjectEnvVars(t *testing.T) {
	os.Setenv("FUGU_TEST123", "ok")

	for _, tt := range injectEnvVarsTests {
		out := injectEnvVars(tt.in)
		assert.Equal(t, tt.out, out, spew.Sdump(tt))
	}
}
