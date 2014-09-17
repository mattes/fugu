package cli

import (
	"fmt"
	"github.com/mattes/fugu/config"
	"github.com/mattes/fugu/docker"
	"github.com/mattes/fugu/file"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
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

	// add --interactive and --tty when command or args are given
	// interactive := false
	// dockerCommand := config.Get(conf, "command")
	// if dockerCommand != nil {
	// 	if dockerCommand.Get().(string) != "" {
	// 		interactive = true
	// 	}
	// }
	// dockerArgs := config.Get(conf, "args")
	// if dockerArgs != nil {
	// 	if len(dockerArgs.Get().([]string)) > 0 {
	// 		interactive = true
	// 	}
	// }
	// if interactive {
	// 	config.Set(&conf, "interactive", true)
	// 	config.Set(&conf, "tty", true)
	// }

	// finally build args
	a := BuildRunArgs(&conf)

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

func CmdBuild(fugufilePath string, args []string, label string) {

	// use image as --tag option
	var fugufileConf = []config.Value{
		&config.StringValue{Name: []string{"image"}},
	}

	// docker options
	// see http://docs.docker.com/reference/commandline/cli/#build

	var dockerBuildConf = []config.Value{
		// add new option: path
		&config.StringValue{Name: []string{"path"}},

		&config.BoolValue{Name: []string{"force-rm"}},
		&config.BoolValue{Name: []string{"no-cache"}},
		&config.BoolValue{Name: []string{"quit", "q"}},
		&config.BoolValue{Name: []string{"rm"}},
		&config.StringValue{Name: []string{"tag", "t"}},
	}

	// read fugufile
	data, err := ioutil.ReadFile(fugufilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// parse fugufile
	_, err = file.Load(data, label, &fugufileConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// parse docker args
	err = docker.Load(args, &dockerBuildConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// build docker build args
	dockerTag := ""
	dockerPath := ""
	a := make([]string, 0)
	for _, c := range dockerBuildConf {
		if c.Names()[0] == "tag" {
			dockerTag = c.Get().(string)
		} else if c.Names()[0] == "path" {
			dockerPath = c.Get().(string)
		} else {
			v := c.Arg()
			if len(v) > 0 {
				a = append(a, c.Arg()...)
			}
		}
	}

	// append dockerImage
	if dockerTag == "" {
		dockerImage := fugufileConf[0].Get().(string)
		a = append(a, fmt.Sprintf(`--tag="%v"`, dockerImage))
	}

	sort.Sort(sort.StringSlice(a))

	// get path | url | -
	if dockerPath != "" {
		a = append(a, dockerPath)
	} else {
		a = append(a, ".")
	}

	a = append(a, "")
	copy(a[1:], a[0:])
	a[0] = "build"

	fmt.Println("docker", strings.Join(a, " "))

	cmd := exec.Command("docker", a...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func CmdDestroy(fugufilePath string, args []string, label string) {

	var fugufileConf = []config.Value{
		&config.StringValue{Name: []string{"name"}},
	}

	// read fugufile
	data, err := ioutil.ReadFile(fugufilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// parse fugufile
	_, err = file.Load(data, label, &fugufileConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// parse docker args
	err = docker.Load(args, &fugufileConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	a := make([]string, 0)
	a = append(a, []string{"rm", "-f"}...)

	dockerName := fugufileConf[0].Get()
	if dockerName != nil && dockerName != "" {
		a = append(a, dockerName.(string))
	} else {
		fmt.Println("Could not find container name.")
		os.Exit(1)
	}

	fmt.Println("docker", strings.Join(a, " "))

	cmd := exec.Command("docker", a...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}
