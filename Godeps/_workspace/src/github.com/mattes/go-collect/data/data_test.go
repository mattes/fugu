package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	p := New()
	assert.NotNil(t, p.data)
}

func TestKeys(t *testing.T) {
	p := New()
	p.Set("foo", "")
	p.Set("bar", "")
	n := p.Keys()
	assert.Contains(t, n, "foo")
	assert.Contains(t, n, "bar")
	assert.Len(t, n, 2)
}

func TestExists(t *testing.T) {
	p := New()
	p.Set("foo", "")
	assert.True(t, p.Exists("foo"))
	assert.False(t, p.Exists("bar"))
}

func TestGetAll(t *testing.T) {
	p := New()
	p.Set("foo", "b", "a", "r")
	assert.Equal(t, []string{"b", "a", "r"}, p.GetAll("foo"))
	assert.Equal(t, []string{}, p.GetAll("non-existing"))
}

func TestGet(t *testing.T) {
	p := New()
	p.Set("foo", "bar")
	assert.Equal(t, "bar", p.Get("foo"))

	p.Set("foo2", "b", "a", "r")
	assert.Equal(t, "", p.Get("foo2"))
}

func TestSet(t *testing.T) {
	p := New()
	p.Set("foo", "bar")
	assert.Equal(t, "bar", p.Get("foo"))

	// test overwrite
	p.Set("foo", "rab")
	assert.Equal(t, "rab", p.Get("foo"))
}

func TestAdd(t *testing.T) {
	p := New()
	p.Add("foo", "bar")
	p.Add("foo", "rab")
	p.Add("foo", "foo", "bar")
	assert.Equal(t, []string{"bar", "rab", "foo", "bar"}, p.GetAll("foo"))
}

func TestDelete(t *testing.T) {
	p := New()
	p.Set("foo", "bar")
	p.Delete("foo")
	assert.Equal(t, []string{}, p.GetAll("foo"))
	p.Delete("non-existing")
}

func TestPickAll(t *testing.T) {
	p := New()
	p.Set("foo", "b", "a", "r")
	assert.Equal(t, []string{"b", "a", "r"}, p.PickAll("foo"))
	assert.False(t, p.Exists("foo"))
}

func TestPick(t *testing.T) {
	p := New()
	p.Set("foo", "bar")
	assert.Equal(t, "bar", p.Pick("foo"))
	assert.False(t, p.Exists("foo"))
}

func TestMerge(t *testing.T) {
	p1 := New()
	p1.Set("1", "1")
	p1.Set("2", "2")
	p1.Set("3", "3")

	p2 := New()
	p2.Set("4", "4")
	p2.Set("2", "b")

	r := New()
	r.Set("1", "1")
	r.Set("2", "b")
	r.Set("3", "3")
	r.Set("4", "4")

	n := Merge(p1, p2)
	assert.Equal(t, r, n)

	p1.Merge(p2)
	assert.Equal(t, r, p1)
}

func TestFilter(t *testing.T) {
	p1 := New()
	p1.Set("1", "1")
	p1.Set("2", "2")
	p1.Set("3", "3")
	p1.Set("4", "4")

	p2 := New()
	p2.Set("2", "2")
	p2.Set("4", "4")

	p1.Filter(p2)
	assert.Equal(t, p2, p1)
}

func TestFilterWithStringsSlice(t *testing.T) {
	p1 := New()
	p1.Set("1", "1")
	p1.Set("2", "2")
	p1.Set("3", "3")
	p1.Set("4", "4")

	p2 := New()
	p2.Set("2", "2")
	p2.Set("4", "4")

	p1.Filter([]string{"4", "2"})
	assert.Equal(t, p2, p1)
}

func TestFilterWithString(t *testing.T) {
	p1 := New()
	p1.Set("1", "1")
	p1.Set("2", "2")
	p1.Set("3", "3")
	p1.Set("4", "4")

	p2 := New()
	p2.Set("2", "2")

	p1.Filter("2")
	assert.Equal(t, p2, p1)
}

func TestToData(t *testing.T) {
	p := New()
	p.Set("foo", "bar")
	tp := ToData(map[string][]string{"foo": []string{"bar"}})
	assert.Equal(t, p, tp)
}

func TestRaw(t *testing.T) {
	// TODO
}

func TestRawEnhanced(t *testing.T) {
	// TODO
}

func TestIsTrue(t *testing.T) {
	// TODO
}

func TestIsFalse(t *testing.T) {
	// TODO
}
