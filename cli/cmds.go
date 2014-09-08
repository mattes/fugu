package cli

import (
	"fmt"
	"github.com/mattes/fugu/config"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func CmdRun(fugufilePath string, args []string, label string) {

	// docker options
	// see http://docs.docker.com/reference/commandline/cli/#run

	var conf = []config.Value{

		// add new options: image, command and args
		&config.StringValue{Name: []string{"image"}},
		&config.StringValue{Name: []string{"command"}},
		&config.StringSliceValue{Name: []string{"args"}},

		// official docker options ...
		&config.BoolValue{Name: []string{"rm"}},
		&config.BoolValue{Name: []string{"detach", "d"}},
		&config.BoolValue{Name: []string{"networking", "n"}},
		&config.BoolValue{Name: []string{"privileged"}},
		&config.BoolValue{Name: []string{"publish-all", "P"}},
		&config.BoolValue{Name: []string{"interactive", "i"}},
		&config.BoolValue{Name: []string{"tty", "t"}},

		&config.StringValue{Name: []string{"cidfile"}},
		&config.StringValue{Name: []string{"entrypoint"}},
		&config.StringValue{Name: []string{"hostname", "h"}},
		&config.StringValue{Name: []string{"memory", "m"}},
		&config.StringValue{Name: []string{"user", "u"}},
		&config.StringValue{Name: []string{"workdir", "w"}},
		&config.Int64Value{Name: []string{"cpu-shares", "c"}},
		&config.StringValue{Name: []string{"cpuset"}},
		&config.StringValue{Name: []string{"net"}},
		&config.StringValue{Name: []string{"restart"}},
		&config.StringValue{Name: []string{"name"}},
		&config.BoolValue{Name: []string{"sig-proxy"}},

		&config.StringSliceValue{Name: []string{"attach", "a"}},
		&config.StringSliceValue{Name: []string{"volume", "v"}},
		&config.StringSliceValue{Name: []string{"link"}},
		&config.StringSliceValue{Name: []string{"device"}},
		&config.StringSliceValue{Name: []string{"env", "e"}},
		&config.StringSliceValue{Name: []string{"env-file"}},
		&config.StringSliceValue{Name: []string{"publish", "p"}},
		&config.StringSliceValue{Name: []string{"expose"}},
		&config.StringSliceValue{Name: []string{"dns"}},
		&config.StringSliceValue{Name: []string{"dns-search"}},
		&config.StringSliceValue{Name: []string{"volumes-from"}},
		&config.StringSliceValue{Name: []string{"lxc-conf"}},
		&config.StringSliceValue{Name: []string{"cap-add"}},
		&config.StringSliceValue{Name: []string{"cap-drop"}},
	}

	// read fugufile
	data, err := ioutil.ReadFile(fugufilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = MergeConfig(data, args, label, &conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	a := BuildArgs(&conf)
	a = append(a, "")
	copy(a[1:], a[0:])
	a[0] = "run"

	fmt.Println("docker", strings.Join(a, " "))

	cmd := exec.Command("docker", a...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

// func CmdBuild(fugufilePath string, args []string, label string) {
// 	conf := config.RunConfig
// 	err := MergeConfig(nil, args, label, &conf)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }
