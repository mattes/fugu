// Package docker parses docker command line arguments and passes those
// args into config.Values
package docker

import (
	"fmt"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/mflag"
	"github.com/mattes/fugu/config"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
)

var (
	ErrInvalidCmdType = fmt.Errorf("Unknown cmd.(type)")
)

// Load hijacks docker/pkg/mflag package for args parsing
func Load(args []string, conf *[]config.Value) error {
	cmd := mflag.NewFlagSet("", mflag.ContinueOnError)
	cmd.SetOutput(ioutil.Discard)
	cmd.Usage = nil

	// what do we want to parse?
	// create cmd.Vars from conf
	for _, c := range *conf {

		// prepare names by inversing slice and prepending hyphens
		// example: `attach, a` becomes `a, -attach`
		names := make([]string, 0)
		for _, n := range c.Names() {
			if len(n) > 1 {
				names = append(names, "-"+n)
			} else {
				names = append(names, n)
			}
		}
		sort.Sort(sort.Reverse(sort.StringSlice(names)))

		// create cmd.Vars for types
		switch c.(type) {
		case *config.StringValue:
			cmd.String(names, "", "")
		case *config.BoolValue:
			cmd.Bool(names, false, "")
		case *config.Int64Value:
			cmd.Int64(names, 0, "")
		case *config.StringSliceValue:
			var dummy = opts.NewListOpts(nil)
			cmd.Var(&dummy, names, "")
		}
	}

	// let docker parse the arguments
	if err := cmd.Parse(args); err != nil {
		return err
	}

	// check remaining args
	rrags := cmd.Args()
	var dockerCommand string
	var dockerArgs []string
	if len(rrags) == 1 {
		dockerCommand = rrags[0]
	} else if len(rrags) >= 2 {
		dockerCommand = rrags[0]
		dockerArgs = rrags[1:]
	}

	// populate into config.Value ...
	if dockerCommand != "" {
		var isSet int
		for _, c := range *conf {
			if c.Names()[0] == "command" {
				c.Set(dockerCommand)
				isSet += 1
			}
			if c.Names()[0] == "args" {
				c.Set(dockerArgs)
				isSet += 1
			}
			if isSet >= 2 {
				break
			}
		}
	}

	var err error
	cmd.VisitAll(func(f *mflag.Flag) {

		if len(f.Names) == 0 {
			err = fmt.Errorf("panic")
			return
		}
		name := f.Names[0]

		// TODO: *mflag.stringValue and all other *mflag.xxxValue
		// are not exported by mflag. Therefore we need to use reflect here.
		// see also http://golang.org/src/pkg/flag/flag.go

		// this is really really bad!
		// TODO: fix me!!!
		typeStr := reflect.TypeOf(f.Value).String()

		for _, c := range *conf {
			namesMatch := false
			for _, n := range c.Names() {
				if n == strings.Trim(name, "-") {
					namesMatch = true
					break
				}
			}

			if namesMatch {

				// TODO the `v != f.DefValue` comparison is hacky.
				// atm there is no way to check if an argument was defined/ given by the user

				switch typeStr {
				case "*mflag.stringValue":
					// TODO check if arg really was given by the user
					v := f.Value.String()
					if v != f.DefValue {
						err = c.Set(v)
					}

				case "*mflag.int64Value":
					// TODO check if arg really was given by the user
					v := f.Value.String()
					if v != f.DefValue {
						err = c.Set(v)
					}

				case "*mflag.boolValue":
					// TODO check if arg really was given by the user
					v := f.Value.String()
					if v != f.DefValue {
						err = c.Set(v)
					}

				case "*opts.ListOpts":
					// TODO check if arg really was given by the user
					v := f.Value.(*opts.ListOpts).GetAll()
					if len(v) > 0 {
						err = c.Set(v)
					}

				default:
					err = fmt.Errorf("panic")
				}
			}
		}
	})
	if err != nil {
		return ErrInvalidCmdType
	}

	return nil
}
