package urlquery

import (
	"fmt"
	"github.com/mattes/go-collect/data"
	"net/url"
)

type UrlQuery struct{}

func (s *UrlQuery) Scheme() string {
	return "urlquery"
}

func (s *UrlQuery) ExampleUrl() string {
	return "urlquery://key=value&..."
}

func (s *UrlQuery) Load(label string, u *url.URL) (*data.Data, error) {
	// rebuild query to avoid leading ?
	var err error
	u, err = url.Parse("?" + u.Host)
	if err != nil {
		return nil, fmt.Errorf("source: urlquery: %v", err.Error())
	}

	d := data.New()
	q := u.Query()
	for k, v := range q {
		d.Set(k, v...)
	}

	return d, nil
}

func (s *UrlQuery) Labels() []string {
	return nil
}
