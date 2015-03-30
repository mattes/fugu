package flags

import (
	"fmt"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/mflag"
	"github.com/docker/docker/utils"
	"github.com/mattes/go-collect/data"
	"strings"
)

type Flags struct {
	Name    string
	flagset *mflag.FlagSet
}

func New(name string) *Flags {
	d := &Flags{}
	d.Name = name
	d.flagset = mflag.NewFlagSet("", mflag.ContinueOnError)
	return d
}

// Keys returns all flag names
// Don't rely on it's order!
func (d *Flags) Keys() ([]string, error) {
	da, err := d.parse(nil, false)
	if err != nil {
		return nil, err
	}
	return da.Keys(), nil
}

// Parse parses the flags and returns the data
func (d *Flags) Parse(args *[]string) (*data.Data, error) {
	return d.parse(args, true)
}

func (d *Flags) parse(args *[]string, setOnly bool) (*data.Data, error) {
	p := data.New()

	if args == nil {
		args = &[]string{}
	}

	// filter out --help as this would trigger fmt.Println() stuff in mflag
	for i, v := range *args {
		if v == "--help" || v == "-help" || v == "-h" {
			*args = append((*args)[:i], (*args)[i+1:]...)
			p.Set("help", "true")
			break
		}
	}

	if err := utils.ParseFlags(d.flagset, *args, false); err != nil {
		return nil, fmt.Errorf("flags: %v", err.Error())
	}

	// the visitor func
	fn := func(m *mflag.Flag) {
		longName := getLongName(m.Names)
		if longName == "" {
			return
		}

		// make sure we don't have any further dependencies to mflag or opts
		switch m.Value.(type) {
		case *opts.ListOpts:
			l := m.Value.(*opts.ListOpts).GetAll()
			p.Set(longName, l...)

		default:
			p.Set(longName, m.Value.String())
		}
	}

	if setOnly {
		d.flagset.Visit(fn)
	} else {
		d.flagset.VisitAll(fn)
	}

	*args = d.flagset.Args()

	return p, nil
}

// FlagCount returns the number of flags that have been defined.
func (d *Flags) FlagCount() int {
	return d.flagset.FlagCount()
}

func (d *Flags) PrintUsage() {
	d.flagset.PrintDefaults()
}

func (d *Flags) Exists(name string) bool {
	if d.flagset.Lookup("-"+name) != nil {
		return true
	}
	return false
}

func (d *Flags) Var(names []string, usage string) *Flags {
	v := opts.NewListOpts(nil)
	d.flagset.Var(&v, names, usage)
	return d
}

func (d *Flags) String(names []string, value string, usage string) *Flags {
	d.flagset.String(names, value, usage)
	return d
}

func (d *Flags) Bool(names []string, value bool, usage string) *Flags {
	d.flagset.Bool(names, value, usage)
	return d
}

func (d *Flags) Int64(names []string, value int64, usage string) *Flags {
	d.flagset.Int64(names, value, usage)
	return d
}

// Merge merges one or more Flags
// Later flags don't overwrite previous ones.
func Merge(f ...*Flags) *Flags {
	n := New("")
	for _, v := range f {
		if v == nil || v.flagset == nil {
			continue
		}
		v.flagset.VisitAll(func(m *mflag.Flag) {

			// dont override if flag already exists
			if n.flagset.Lookup("-"+getLongName(m.Names)) != nil {
				return
			}

			switch m.Value.(type) {
			case *opts.ListOpts:
				n.Var(m.Names, m.Usage)
			default:
				value := m.Value.(mflag.Getter).Get()
				switch value.(type) {
				case string:
					n.String(m.Names, value.(string), m.Usage)
				case bool:
					n.Bool(m.Names, value.(bool), m.Usage)
				case int64:
					n.Int64(m.Names, value.(int64), m.Usage)
				default:
					panic("unsupported flag type")
				}
			}
		})
	}
	return n
}

// Nice formats flag nicely
func Nice(key, value string) string {
	if value == "true" {
		return fmt.Sprintf("--%v", key)
	}

	if strings.Contains(value, " ") || value == "" {
		return fmt.Sprintf("--%v='%v'", key, value)
	} else {
		return fmt.Sprintf("--%v=%v", key, value)
	}
}

// getLongName returns long docker option name
func getLongName(names []string) string {
	for _, n := range names {
		if strings.HasPrefix(n, "-") {
			return strings.TrimLeft(n, "-")
		}
	}
	return ""
}
