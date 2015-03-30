package collect

import (
	"github.com/mattes/go-collect/source/urlquery"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSource(t *testing.T) {
	assert.Equal(t, nil, GetSource("bogus"))

	RegisterSource(&urlquery.UrlQuery{})
	assert.Equal(t, &urlquery.UrlQuery{}, GetSource("urlquery"))
}
