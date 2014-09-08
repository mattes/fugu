package file

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/mattes/fugu/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var parseTests = []struct {
	in  []byte
	out []Label
	err bool
}{

	// test if default label is set if none is present
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
			&config.StringValue{Name: []string{"name"}, Value: "test", Present: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar", Present: true},
			&config.StringSliceValue{Name: []string{"publish"}, Value: []string{"8080:80"}, Present: true},
			&config.StringValue{Name: []string{"non-exist"}, Present: false},
		},
		false,
	},

	{
		labelData1,
		"label",
		[]string{"label", "another-label"},
		[]config.Value{
			&config.StringValue{Name: []string{"name"}, Value: "test", Present: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar", Present: true},
			&config.StringSliceValue{Name: []string{"publish"}, Value: []string{"8080:80"}, Present: true},
			&config.StringValue{Name: []string{"non-exist"}, Present: false},
		},
		false,
	},

	{
		labelData1,
		"another-label",
		[]string{"label", "another-label"},
		[]config.Value{
			&config.StringValue{Name: []string{"name"}, Value: "halligalli", Present: true},
			&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar2", Present: true},
			&config.StringSliceValue{Name: []string{"publish"}, Value: []string{"55:66"}, Present: true},
			&config.StringValue{Name: []string{"non-exist"}, Present: false},
		},
		false,
	},

	// test if first label is fetched
	// TODO(mattes): this is buggy, because we cannot garantuee the sort order
	// see comments in file.go
	// {
	// 	labelData1,
	// 	"",
	// 	[]string{"label", "another-label"},
	// 	[]config.Value{
	// 		&config.StringValue{Name: []string{"name"}, Value: "test", Present: true},
	// 		&config.StringValue{Name: []string{"image"}, Value: "mattes/foobar", Present: true},
	// 		&config.StringSliceValue{Name: []string{"publish"}, Value: []string{"8080:80"}, Present: true},
	// 		&config.StringValue{Name: []string{"non-exist"}, Present: false},
	// 	},
	// 	false,
	// },
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
