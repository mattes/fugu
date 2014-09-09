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
	if conf == nil {
		panic("Provide conf *[]config.Value")
	}

	cmd := mflag.NewFlagSet("", mflag.ContinueOnError)
	cmd.SetOutput(ioutil.Discard)
	cmd.Usage = nil

	// what do we want to parse?
	// create cmd.Vars from conf
	for _, c := range *conf {

		// prepare names by inverting slice and prepending hyphens
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

	// let docker/pkg/mflag parse the arguments
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

		// TODO: fix me!!!
		// this is really really bad! we are no comparing strings instead of real types
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
				if flagWasDefined(c.Names(), args) {
					switch typeStr {
					case "*mflag.stringValue":
						err = c.Set(f.Value.String())

					case "*mflag.int64Value":
						err = c.Set(f.Value.String())

					case "*mflag.boolValue":
						err = c.Set(f.Value.String())

					case "*opts.ListOpts":
						err = c.Set(f.Value.(*opts.ListOpts).GetAll())

					default:
						err = fmt.Errorf("panic")
					}
				}
			}
		}
	})
	if err != nil {
		return ErrInvalidCmdType
	}

	return nil
}

// flagWasDefined is a hack to check if a flag has been set
// mflag package won't tell us
func flagWasDefined(findName []string, args []string) bool {
	for _, a := range args {
		for _, n := range findName {
			n2 := strings.Trim(n, "-")
			if len(n) == 1 {
				n2 = "-" + n2
			} else if len(n) > 1 {
				n2 = "--" + n2
			}

			if strings.HasPrefix(a, n2) {
				if len(n) == 1 {
					// -p != -publish
					if strings.HasPrefix(a, n2+"=") || strings.HasPrefix(a, n2+" ") || a == n2 {
						return true
					}
				} else {
					return true
				}
			}
		}
	}
	return false
}
