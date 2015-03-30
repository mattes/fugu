package flags

import (
	"github.com/mattes/go-collect/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	f := New("foo")
	assert.NotNil(t, f.flagset)
	assert.Equal(t, "foo", f.Name)
}

func TestKeys(t *testing.T) {
	f := New("")
	f.String([]string{"-foo"}, "", "")
	f.String([]string{"-bar"}, "", "")

	keys, err := f.Keys()
	assert.NoError(t, err)
	assert.Contains(t, keys, "foo")
	assert.Contains(t, keys, "bar")
	assert.Len(t, keys, 2)
}

func TestParse(t *testing.T) {
	// assume that github.com/docker/docker/pkg/mflag is tested

	f := New("")
	f.String([]string{"-foo"}, "", "")

	dr := data.New()
	dr.Set("foo", "bar")

	d, err := f.Parse(&[]string{"--foo=bar"})
	assert.NoError(t, err)
	assert.Equal(t, dr, d)

	_, err = f.Parse(nil)
	assert.NoError(t, err)
}

func TestParseOrder(t *testing.T) {
	f := New("")
	f.Var([]string{"-foo"}, "")

	dr := data.New()
	dr.Set("foo", "b", "a", "c", "f", "e", "d", "in", "this", "order")

	d, err := f.Parse(&[]string{
		"--foo=b", "--foo=a", "--foo=c",
		"--foo=f", "--foo=e", "--foo=d",
		"--foo=in", "--foo=this", "--foo=order"})
	assert.NoError(t, err)
	assert.Equal(t, dr, d)
}

func TestFlagCount(t *testing.T) {
	f := New("")
	f.String([]string{"-foo"}, "", "")
	assert.Equal(t, 1, f.FlagCount())
}

func TestExists(t *testing.T) {
	f := New("")
	f.String([]string{"-foo"}, "", "")
	assert.True(t, f.Exists("foo"))
	assert.False(t, f.Exists("bogus"))
}

func TestVar(t *testing.T) {
	f := New("")
	f.Var([]string{"-foo"}, "")

	dr := data.New()
	dr.Set("foo", "bar", "rab", "in-this-order")

	d, err := f.Parse(&[]string{"--foo=bar", "--foo=rab", "--foo=in-this-order"})
	assert.NoError(t, err)
	assert.Equal(t, dr, d)
}

func TestString(t *testing.T) {
	f := New("")
	f.String([]string{"-foo"}, "", "")

	dr := data.New()
	dr.Set("foo", "bar")

	d, err := f.Parse(&[]string{"--foo=bar"})
	assert.NoError(t, err)
	assert.Equal(t, dr, d)
}

func TestBool(t *testing.T) {
	f := New("")
	f.Bool([]string{"-foo"}, true, "")

	dr := data.New()
	dr.Set("foo", "true")

	d, err := f.Parse(&[]string{"--foo"})
	assert.NoError(t, err)
	assert.Equal(t, dr, d)

	dr = data.New()
	dr.Set("foo", "false")

	d, err = f.Parse(&[]string{"--foo=false"})
	assert.NoError(t, err)
	assert.Equal(t, dr, d)
}

func TestInt64(t *testing.T) {
	f := New("")
	f.Int64([]string{"-foo"}, 5, "")

	dr := data.New()
	dr.Set("foo", "5")

	d, err := f.Parse(&[]string{"--foo=5"})
	assert.NoError(t, err)
	assert.Equal(t, dr, d)
}

func TestMerge(t *testing.T) {
	f1 := New("")
	f1.String([]string{"-foo"}, "", "")
	f1.String([]string{"-foo1"}, "", "")

	f2 := New("")
	f2.String([]string{"-foo"}, "", "")
	f1.String([]string{"-foo2"}, "", "")

	fr := New("")
	fr.String([]string{"-foo"}, "", "")
	fr.String([]string{"-foo1"}, "", "")
	fr.String([]string{"-foo2"}, "", "")

	fMerged := Merge(f1, f2)
	assert.Equal(t, fr, fMerged)

	fMerged2 := Merge(f2, f1)
	assert.Equal(t, fr, fMerged2)
}

func TestNice(t *testing.T) {
	assert.Equal(t, "--foo=bar", Nice("foo", "bar"))
	assert.Equal(t, "--foo='b a r'", Nice("foo", "b a r"))
	assert.Equal(t, "--foo", Nice("foo", "true"))
	assert.Equal(t, "--foo=false", Nice("foo", "false"))
	assert.Equal(t, "--foo=''", Nice("foo", ""))
	assert.Equal(t, "--foo=5", Nice("foo", "5"))
	assert.Equal(t, "--foo=5:5", Nice("foo", "5:5"))
}

func TestGetLongName(t *testing.T) {
	assert.Equal(t, "foo", getLongName([]string{"f", "-foo"}))
	assert.Equal(t, "foo", getLongName([]string{"-foo"}))
	assert.Equal(t, "", getLongName([]string{}))
	assert.Equal(t, "", getLongName([]string{"f"}))
}
