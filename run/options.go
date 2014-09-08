package run

import (
	_ "fmt"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/mflag"
	"io/ioutil"
	"strings"
)

func Parse(args []string) (map[string]interface{}, *mflag.FlagSet, error) {
	cmd := mflag.NewFlagSet("run", mflag.ContinueOnError)
	cmd.SetOutput(ioutil.Discard)
	cmd.Usage = nil
	dummy := opts.NewListOpts(nil)

	// add new option for IMAGE
	cmd.String([]string{"#image", "-image"}, "", "")

	// keep this in sync with:
	// https://github.com/docker/docker/blob/master/runconfig/parse.go

	cmd.Bool([]string{"#rm", "-rm"}, false, "")
	cmd.Bool([]string{"d", "-detach"}, false, "")
	cmd.Bool([]string{"#n", "#-networking"}, true, "")
	cmd.Bool([]string{"#privileged", "-privileged"}, false, "")
	cmd.Bool([]string{"P", "-publish-all"}, false, "")
	cmd.Bool([]string{"i", "-interactive"}, false, "")
	cmd.Bool([]string{"t", "-tty"}, false, "")
	cmd.String([]string{"#cidfile", "-cidfile"}, "", "")
	cmd.String([]string{"#entrypoint", "-entrypoint"}, "", "")
	cmd.String([]string{"h", "-hostname"}, "", "")
	cmd.String([]string{"m", "-memory"}, "", "")
	cmd.String([]string{"u", "-user"}, "", "")
	cmd.String([]string{"w", "-workdir"}, "", "")
	cmd.Int64([]string{"c", "-cpu-shares"}, 0, "")
	cmd.String([]string{"-cpuset"}, "", "")
	cmd.String([]string{"-net"}, "bridge", "")
	cmd.String([]string{"-restart"}, "", "")
	cmd.Bool([]string{"#sig-proxy", "-sig-proxy"}, true, "")
	cmd.String([]string{"#name", "-name"}, "", "")

	cmd.Var(&dummy, []string{"a", "-attach"}, "")
	cmd.Var(&dummy, []string{"v", "-volume"}, "")
	cmd.Var(&dummy, []string{"#link", "-link"}, "")
	cmd.Var(&dummy, []string{"-device"}, "")
	cmd.Var(&dummy, []string{"e", "-env"}, "")
	cmd.Var(&dummy, []string{"-env-file"}, "")
	cmd.Var(&dummy, []string{"p", "-publish"}, "")
	cmd.Var(&dummy, []string{"#expose", "-expose"}, "")
	cmd.Var(&dummy, []string{"#dns", "-dns"}, "")
	cmd.Var(&dummy, []string{"-dns-search"}, "")
	cmd.Var(&dummy, []string{"#volumes-from", "-volumes-from"}, "")
	cmd.Var(&dummy, []string{"#lxc-conf", "-lxc-conf"}, "")
	cmd.Var(&dummy, []string{"-cap-add"}, "")
	cmd.Var(&dummy, []string{"-cap-drop"}, "")

	if err := cmd.Parse(args); err != nil {
		return nil, nil, err
	}

	config := make(map[string]interface{})
	cmd.Visit(func(f *mflag.Flag) {
		var name string
		if len(f.Names) > 1 {
			name = f.Names[1]
		} else {
			name = f.Names[0]
		}
		name = strings.Trim(name, "#-")
		config[name] = f.Value
	})

	return config, cmd, nil
}
