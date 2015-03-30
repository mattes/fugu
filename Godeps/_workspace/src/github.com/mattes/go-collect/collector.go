package collect

import (
	"errors"
	"fmt"
	"github.com/mattes/go-collect/data"
	"github.com/mattes/go-collect/flags"
	"net/url"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrUnknownScheme = errors.New("source: scheme is unknown")
)

type Collector struct {
	// args are options given by the user via the cli
	args []string

	// flags define all formal possible args
	flags []*flags.Flags

	// label is a arg given by the user via the cli
	label string

	// sources is a list of url strings to source locations
	sources []string

	defaultSource string
}

func New() *Collector {
	return &Collector{
		args:          make([]string, 0),
		flags:         make([]*flags.Flags, 0),
		label:         "",
		sources:       make([]string, 0),
		defaultSource: "",
	}
}

func (c *Collector) Parse(args []string, f ...*flags.Flags) (p *data.Data, remainingArgs []string, err error) {
	c.args = args
	c.label = c.parseLabel()
	c.AddFlags(f...)
	combinedFlags := flags.New("")

	if !c.flagDefined("source") {
		f := flags.New("")
		f.Var([]string{"-source"}, "Get data from this source")
		combinedFlags = flags.Merge(combinedFlags, f)
	}

	for _, f := range c.flags {
		combinedFlags = flags.Merge(combinedFlags, f)
	}
	argsData, err := combinedFlags.Parse(&c.args)
	if err != nil {
		return nil, nil, err
	}

	c.AddSource(argsData.PickAll("source")...)

	sourceData := data.New()
	for _, sarg := range c.Sources() {
		s, err := c.getSourceFromScheme(sarg)
		if err != nil {
			return nil, nil, err
		}

		// TODO do this async
		u, _ := url.Parse(sarg)
		p, err := s.Load(c.label, u)
		if err != nil {
			return nil, nil, err
		}

		// merge data from Load
		sourceData.Merge(p)
	}

	// overwrite with args data
	sourceData.Merge(argsData)

	return sourceData, c.args, nil
}

func (c *Collector) getSourceFromScheme(source string) (Source, error) {
	u, err := url.Parse(source)
	if err != nil {
		return nil, ErrUnknownScheme
	}

	s := GetSource(u.Scheme)
	if s == nil {
		return nil, ErrUnknownScheme
	}

	return s, nil
}

func (c *Collector) Labels() []string {
	rl := make([]string, 0)
	for _, sarg := range c.Sources() {
		s, err := c.getSourceFromScheme(sarg)
		if err == nil {
			if labels := s.Labels(); labels != nil {
				rl = append(rl, labels...)
			}
		}
	}
	return rl
}

func (c *Collector) AddFlags(f ...*flags.Flags) {
	for _, ff := range f {
		if ff != nil {
			c.flags = append(c.flags, ff)
		}
	}
}

func (c *Collector) AddSource(s ...string) {
	c.sources = append(c.sources, s...)
}

func (c *Collector) flagDefined(name string) bool {
	for _, f := range c.flags {
		if f.Exists(name) {
			return true
		}
	}
	return false
}

func (c *Collector) parseLabel() (label string) {
	return ParseLabel(&c.args)
}

func (c *Collector) Label() string {
	return c.label
}

func (c *Collector) Sources() []string {
	if len(c.sources) == 0 && c.defaultSource != "" {
		return []string{c.defaultSource}
	} else {
		return c.sources
	}
}

func (c *Collector) SetDefaultSource(s string) {
	// TODO set default source per scheme/source
	c.defaultSource = s
}

func (c *Collector) GetDefaultSource() string {
	// TODO get default source per scheme/source
	return c.defaultSource
}

func (c *Collector) PrintUsage() {
	for _, f := range c.flags {
		fmt.Fprintln(os.Stderr, "")
		if f.Name != "" {
			fmt.Fprintln(os.Stderr, upperFirst(f.Name)+" options:")
		}
		f.PrintUsage()
	}
}

func ParseLabel(args *[]string) (label string) {
	if len(*args) > 0 && !strings.HasPrefix((*args)[0], "-") {
		label = (*args)[0]
		*args = (*args)[1:]
		return label
	}
	return ""
}

// upperFirst makes first letter uppercase
func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}
