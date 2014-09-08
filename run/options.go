package run

import (
	_ "fmt"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/mflag"
	"io/ioutil"
)

func Parse(args []string) (map[string]interface{}, *mflag.FlagSet, error) {
	cmd := mflag.NewFlagSet("run", mflag.ContinueOnError)
	cmd.SetOutput(ioutil.Discard)
	cmd.Usage = nil

	// keep this in sync with:
	// https://github.com/docker/docker/blob/master/runconfig/parse.go

	var rm, detach, networking, privileged, publishAll,
		interactive, tty, sigProxy bool

	var cidfile, entrypoint, hostname, memory, user, workdir, cpuset, net,
		restart, name, image string

	var cpuShares int64

	var attachL = opts.NewListOpts(nil)
	var volumeL = opts.NewListOpts(nil)
	var linkL = opts.NewListOpts(nil)
	var deviceL = opts.NewListOpts(nil)
	var envL = opts.NewListOpts(nil)
	var envFileL = opts.NewListOpts(nil)
	var publishL = opts.NewListOpts(nil)
	var exposeL = opts.NewListOpts(nil)
	var dnsL = opts.NewListOpts(nil)
	var dnsSearchL = opts.NewListOpts(nil)
	var volumesFromL = opts.NewListOpts(nil)
	var lxcConfL = opts.NewListOpts(nil)
	var capAddL = opts.NewListOpts(nil)
	var capDropL = opts.NewListOpts(nil)

	// add new option for IMAGE
	cmd.StringVar(&image, []string{"#image", "-image"}, "", "")

	cmd.BoolVar(&rm, []string{"#rm", "-rm"}, false, "")
	cmd.BoolVar(&detach, []string{"d", "-detach"}, false, "")
	cmd.BoolVar(&networking, []string{"#n", "#-networking"}, true, "")
	cmd.BoolVar(&privileged, []string{"#privileged", "-privileged"}, false, "")
	cmd.BoolVar(&publishAll, []string{"P", "-publish-all"}, false, "")
	cmd.BoolVar(&interactive, []string{"i", "-interactive"}, false, "")
	cmd.BoolVar(&tty, []string{"t", "-tty"}, false, "")

	cmd.StringVar(&cidfile, []string{"#cidfile", "-cidfile"}, "", "")
	cmd.StringVar(&entrypoint, []string{"#entrypoint", "-entrypoint"}, "", "")
	cmd.StringVar(&hostname, []string{"h", "-hostname"}, "", "")
	cmd.StringVar(&memory, []string{"m", "-memory"}, "", "")
	cmd.StringVar(&user, []string{"u", "-user"}, "", "")
	cmd.StringVar(&workdir, []string{"w", "-workdir"}, "", "")
	cmd.Int64Var(&cpuShares, []string{"c", "-cpu-shares"}, 0, "")
	cmd.StringVar(&cpuset, []string{"-cpuset"}, "", "")
	cmd.StringVar(&net, []string{"-net"}, "bridge", "")
	cmd.StringVar(&restart, []string{"-restart"}, "", "")
	cmd.StringVar(&name, []string{"#name", "-name"}, "", "")
	cmd.BoolVar(&sigProxy, []string{"#sig-proxy", "-sig-proxy"}, true, "")

	cmd.Var(&attachL, []string{"a", "-attach"}, "")
	cmd.Var(&volumeL, []string{"v", "-volume"}, "")
	cmd.Var(&linkL, []string{"#link", "-link"}, "")
	cmd.Var(&deviceL, []string{"-device"}, "")
	cmd.Var(&envL, []string{"e", "-env"}, "")
	cmd.Var(&envFileL, []string{"-env-file"}, "")
	cmd.Var(&publishL, []string{"p", "-publish"}, "")
	cmd.Var(&exposeL, []string{"#expose", "-expose"}, "")
	cmd.Var(&dnsL, []string{"#dns", "-dns"}, "")
	cmd.Var(&dnsSearchL, []string{"-dns-search"}, "")
	cmd.Var(&volumesFromL, []string{"#volumes-from", "-volumes-from"}, "")
	cmd.Var(&lxcConfL, []string{"#lxc-conf", "-lxc-conf"}, "")
	cmd.Var(&capAddL, []string{"-cap-add"}, "")
	cmd.Var(&capDropL, []string{"-cap-drop"}, "")

	if err := cmd.Parse(args); err != nil {
		return nil, nil, err
	}

	config := make(map[string]interface{})

	config["rm"] = rm
	config["detach"] = detach
	config["networking"] = networking
	config["privileged"] = privileged
	config["publish-all"] = publishAll
	config["interactive"] = interactive
	config["tty"] = tty
	config["sig-proxy"] = sigProxy
	config["cidfile"] = cidfile
	config["entrypoint"] = entrypoint
	config["hostname"] = hostname
	config["memory"] = memory
	config["user"] = user
	config["workdir"] = workdir
	config["cpuset"] = cpuset
	config["net"] = net
	config["restart"] = restart
	config["name"] = name
	config["image"] = image

	config["attach"] = attachL.GetAll()
	config["volume"] = volumeL.GetAll()
	config["link"] = linkL.GetAll()
	config["device"] = deviceL.GetAll()
	config["env"] = envL.GetAll()
	config["env-file"] = envFileL.GetAll()
	config["publish"] = publishL.GetAll()
	config["expose"] = exposeL.GetAll()
	config["dns"] = dnsL.GetAll()
	config["dns-search"] = dnsSearchL.GetAll()
	config["volumes-from"] = volumesFromL.GetAll()
	config["lxc-conf"] = lxcConfL.GetAll()
	config["cap-add"] = capAddL.GetAll()
	config["cap-drop"] = capDropL.GetAll()

	return config, cmd, nil
}
