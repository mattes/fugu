package collect

import (
	"github.com/mattes/go-collect/data"
	"github.com/mattes/go-collect/flags"
	"net/url"
	"strings"
)

type Source interface {
	// The funcs below get called exactly in this order ...

	// Return the source scheme, like file or etcd
	Scheme() string

	// Return example url with scheme, ie file://config.yml
	ExampleUrl() string

	// Gets called after the flags are parsed
	// Read or load any additional sources here
	// Return data for label
	Load(label string, u *url.URL) (*data.Data, error)

	// Returns all labels
	Labels() []string
}

var sources = make(map[string]Source)

func RegisterSource(s Source) {
	if s == nil {
		panic("source: Register source is nil")
	}
	if s.Scheme() == "" {
		panic("source: Scheme() should not be empty")
	}
	if strings.HasSuffix(s.Scheme(), "://") {
		panic("source: Scheme() should not end with ://")
	}
	if !strings.HasPrefix(s.ExampleUrl(), s.Scheme()+"://") {
		panic("source: ExampleUrl() must start with Scheme()")
	}
	// TODO needed?!
	// if _, dup := sources[s.Scheme()]; dup {
	// 	panic("source: RegisterSource called twice for source " + s.Scheme())
	// }
	sources[s.Scheme()] = s
}

func Sources() map[string]Source {
	return sources
}

func GetSource(scheme string) Source {
	s, ok := sources[scheme]
	if !ok {
		return nil
	}
	return s
}

func SourceExampleUrls() []string {
	out := []string{}
	if len(sources) > 0 {
		for _, v := range sources {
			out = append(out, flags.Nice("source", v.ExampleUrl()))
		}
	}
	return out
}
