package urlquery

import (
	"github.com/mattes/go-collect"
	"github.com/mattes/go-collect/flags"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	collect.RegisterSource(&UrlQuery{})
	c := collect.New()

	f := flags.New("")
	args := []string{"--source", "urlquery://key=value&key2=value2&k=1&k=2"}

	data, _, _ := c.Parse(args, f)

	assert.Equal(t, "value", data.Get("key"))
	assert.Equal(t, "value2", data.Get("key2"))
	assert.Equal(t, []string{"1", "2"}, data.GetAll("k"))
}
