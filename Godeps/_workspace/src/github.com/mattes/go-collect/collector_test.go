package collect

import (
	"errors"
	"github.com/mattes/go-collect/data"
	"github.com/mattes/go-collect/flags"
	"github.com/mattes/go-collect/source/urlquery"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	assert.NotNil(t, c.args)
	assert.NotNil(t, c.flags)
	assert.NotNil(t, c.sources)
}

func TestParseLabel(t *testing.T) {
	c := New()
	c.args = []string{"label", "--opt=value"}
	assert.Equal(t, "label", c.parseLabel())
	assert.Equal(t, []string{"--opt=value"}, c.args)

	c.args = []string{"label", "-opt"}
	assert.Equal(t, "label", c.parseLabel())
	assert.Equal(t, []string{"-opt"}, c.args)

	c.args = []string{"--opt=value"}
	assert.Equal(t, "", c.parseLabel())
	assert.Equal(t, []string{"--opt=value"}, c.args)

	c.args = []string{"", "--opt=value"}
	assert.Equal(t, "", c.parseLabel())
	assert.Equal(t, []string{"--opt=value"}, c.args)

	c.args = []string{"--foo=bar", "label", "--opt=value"}
	assert.Equal(t, "", c.parseLabel())
	assert.Equal(t, []string{"--foo=bar", "label", "--opt=value"}, c.args)
}

func TestParse(t *testing.T) {
	f := flags.New("")
	f.String([]string{"-foo"}, "", "")
	f.String([]string{"-bar"}, "", "")

	RegisterSource(&urlquery.UrlQuery{})

	var tests = []struct {
		testDesc      string
		args          []string
		flags         *flags.Flags
		d             *data.Data
		remainingArgs []string
		err           error
	}{
		{
			testDesc:      "simple",
			args:          []string{"label", "--foo=bar"},
			flags:         f,
			d:             data.ToData(map[string][]string{"foo": []string{"bar"}}),
			remainingArgs: []string{},
			err:           nil,
		},
		{
			testDesc:      "with remaining args",
			args:          []string{"label", "--foo=bar", "."},
			flags:         f,
			d:             data.ToData(map[string][]string{"foo": []string{"bar"}}),
			remainingArgs: []string{"."},
			err:           nil,
		},
		{
			// foo: bar -> from args
			// bar: foo -> from source
			testDesc:      "with source",
			args:          []string{"label", "--foo=bar", "--source=urlquery://foo=rab&bar=foo"},
			flags:         f,
			d:             data.ToData(map[string][]string{"foo": []string{"bar"}, "bar": []string{"foo"}}),
			remainingArgs: []string{},
			err:           nil,
		},
		{
			// foo: bar -> from args
			// bar: oof -> from last source
			testDesc:      "with multiple sources",
			args:          []string{"label", "--foo=bar", "--source=urlquery://foo=rab&bar=foo", "--source=urlquery://foo=arb&bar=oof"},
			flags:         f,
			d:             data.ToData(map[string][]string{"foo": []string{"bar"}, "bar": []string{"oof"}}),
			remainingArgs: []string{},
			err:           nil,
		},
		{
			testDesc:      "flag provided but not defined",
			args:          []string{"label", "--foo=bar", "--bogus='option is not defined'"},
			flags:         f,
			d:             data.New(),
			remainingArgs: []string{},
			err:           errors.New("flags: flag provided but not defined: --bogus"),
		},
		{
			testDesc:      "source gets defined automatically",
			args:          []string{"label", "--foo=bar", "--source=urlquery://"},
			flags:         f,
			d:             data.ToData(map[string][]string{"foo": []string{"bar"}}),
			remainingArgs: []string{},
			err:           nil,
		},
		{
			testDesc:      "fail on invalid source scheme",
			args:          []string{"label", "--foo=bar", "--source=bogus://"},
			flags:         f,
			d:             data.New(),
			remainingArgs: []string{},
			err:           ErrUnknownScheme,
		},
	}

	for _, tt := range tests {
		f := New()

		// capture Stderr as f.Parse fmt.Prints(os.stderr) ... sucks
		rescueStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		d, remainingArgs, err := f.Parse(tt.args, tt.flags)

		w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stderr = rescueStderr

		if !assert.Equal(t, tt.err, err, tt.testDesc) {
			t.Logf("%s", out)
		}

		assert.Equal(t, tt.err, err, tt.testDesc)
		if err == nil {
			assert.Equal(t, tt.d, d, tt.testDesc)
			assert.Equal(t, tt.remainingArgs, remainingArgs, tt.testDesc)
		}
	}
}

func TestAddFlags(t *testing.T) {
	c := New()
	f1 := flags.New("f1")
	f2 := flags.New("f2")
	c.AddFlags(f2, f1)
	assert.Equal(t, []*flags.Flags{f2, f1}, c.flags)
}

func TestAddSource(t *testing.T) {
	c := New()
	c.AddSource("dummy://", "another://")
	assert.Equal(t, []string{"dummy://", "another://"}, c.sources)
}

func TestSetDefaultSource(t *testing.T) {
	c := New()
	c.SetDefaultSource("dummy://foobar")
	assert.Equal(t, []string{"dummy://foobar"}, c.Sources())
}

func TestGetDefaultSource(t *testing.T) {
	c := New()
	c.SetDefaultSource("dummy://foobar")
	assert.Equal(t, "dummy://foobar", c.GetDefaultSource())
}

func TestSources(t *testing.T) {
	c := New()
	assert.Equal(t, []string{}, c.Sources())

	c.SetDefaultSource("dummy://")
	assert.Equal(t, []string{"dummy://"}, c.Sources())

	c.AddSource("another://")
	assert.Equal(t, []string{"another://"}, c.Sources())
}

func TestUpperFirst(t *testing.T) {
	assert.Equal(t, "Foobar", upperFirst("foobar"))
	assert.Equal(t, "Foobar", upperFirst("Foobar"))
	assert.Equal(t, "Foo bar", upperFirst("foo bar"))
	assert.Equal(t, "F", upperFirst("f"))
	assert.Equal(t, "", upperFirst(""))
}
