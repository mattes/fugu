package file

import (
	"github.com/mattes/go-collect/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		testDesc string
		body     string
		label    string

		data   *data.Data
		labels []string
		err    error
	}{
		{
			testDesc: "no labels",
			body: `
      image: test
      name: foobar
      foo:
        - bar
        - hu
      `,
			label: "",

			data: data.ToData(map[string][]string{
				"image": []string{"test"},
				"name":  []string{"foobar"},
				"foo":   []string{"bar", "hu"},
			}),
			labels: []string{"default"},
			err:    nil,
		},

		{
			testDesc: "with label",
			body: `
      label:
        image: test
        name: foobar
      `,
			label: "label",
			data: data.ToData(map[string][]string{
				"image": []string{"test"},
				"name":  []string{"foobar"},
			}),
			labels: []string{"label"},
			err:    nil,
		},

		{
			testDesc: "get first label",
			body: `
      label:
        image: test
        name: foobar

      another-label:
        image: test

      label5:
        image: test1
      `,
			label: "",
			data: data.ToData(map[string][]string{
				"image": []string{"test"},
				"name":  []string{"foobar"},
			}),
			labels: []string{"label", "another-label", "label5"},
			err:    nil,
		},

		{
			testDesc: "get default label if exists",
			body: `
      label:
        image: test
        name: foobar

      default:
        image: use-this

      another-label: ~
      `,
			label: "",
			data: data.ToData(map[string][]string{
				"image": []string{"use-this"},
			}),
			labels: []string{"label", "default", "another-label"},
			err:    nil,
		},

		{
			testDesc: "only support two levels of indentation",
			body: `
      label:
        image: test
        name: foobar
        fails:
          foo: # ok
            - bar # not ok
      `,
			label:  "",
			data:   data.ToData(map[string][]string{}),
			labels: nil,
			err:    ErrYamlLevels,
		},

		{
			testDesc: "simple inheritance",
			body: `
      label1:
        image: image
        name: name
      label2:
        <<: label1
        foo: foo
        image: take-this
      `,
			label: "label2",

			data: data.ToData(map[string][]string{
				"image": []string{"take-this"},
				"name":  []string{"name"},
				"foo":   []string{"foo"},
			}),
			labels: []string{"label1", "label2"},
			err:    nil,
		},

		{
			testDesc: "inheritance in inheritance",
			body: `
      label1:
        <<: label4  
        image: image
        name: name    
      label3: 
        <<: label2  
        bar: bar
        foo: foo3   
        name: name2 
      label2: 
        <<: label1
        foo: foo    
      label4: 
        oof: oof
        name: name4    
      `,
			label: "label3",

			data: data.ToData(map[string][]string{
				"image": []string{"image"},
				"name":  []string{"name2"},
				"foo":   []string{"foo3"},
				"bar":   []string{"bar"},
				"oof":   []string{"oof"},
			}),
			labels: []string{"label1", "label3", "label2", "label4"},
			err:    nil,
		},
	}

	for _, tt := range tests {
		if tt.testDesc != "inheritance" {
			continue
		}
		f := File{}
		f.body = []byte(tt.body)
		err := f.parse()
		assert.Equal(t, tt.err, err, tt.testDesc)
		if err == nil {
			assert.Equal(t, tt.data, f.getData(), tt.testDesc)
			assert.Equal(t, tt.labels, f.labels, tt.testDesc)
		}
	}
}

func TestSetPathFromUrl(t *testing.T) {
	// TODO
}

func TestLabels(t *testing.T) {
	f := File{}
	f.body = []byte(`
  label:
    image: test
  another-label:
    image: test`)
	assert.NoError(t, f.parse())
	assert.Contains(t, f.Labels(), "label")
	assert.Contains(t, f.Labels(), "another-label")
	assert.Equal(t, 2, len(f.Labels()))
}

func TestSelectLabel(t *testing.T) {
	var tests = []struct {
		labelsInFile []string
		labelFromArg string
		outLabel     string
	}{
		// all empty
		{[]string{}, "", ""},

		// If you don't ask for a specific label, it will ...
		// return label ``default`` if found, or
		{[]string{"foo", "default", "bar"}, "", "default"},

		// ... return first found label
		{[]string{"foo", "bar"}, "", "foo"},

		// If you ask for a specific label, it will ...
		// return this specific label if found, and
		{[]string{"foobar", "label"}, "label", "label"},

		// just don't return anything
		{[]string{}, "not-found", ""},
		{[]string{"foobar", "label"}, "not-found", ""},
	}

	for _, tt := range tests {
		f := File{}
		f.label = tt.labelFromArg
		f.labels = tt.labelsInFile
		assert.Equal(t, tt.outLabel, f.selectLabel(tt.labelFromArg))
	}
}
