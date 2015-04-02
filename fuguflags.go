package fugu

import (
	"github.com/mattes/go-collect/flags"
)

var FuguFlags = make(map[string]*flags.Flags)

func init() {

	FuguCommon := flags.New("")
	FuguCommon.Var([]string{"-source"}, "Get data from this source")
	FuguCommon.Bool([]string{"-dry-run"}, false, "Just print commands")

	// Define FuguFlags["build"]
	FuguFlags["build"] = flags.New("fugu")
	FuguFlags["build"].String([]string{"-image"}, "", "Name of the image")
	FuguFlags["build"].String([]string{"-path"}, "", "PATH")
	FuguFlags["build"].String([]string{"-url"}, "", "URL")
	FuguFlags["build"].Bool([]string{"-tag-git-branch"}, false, "Tag with current git branch")
	FuguFlags["build"] = flags.Merge(FuguCommon, FuguFlags["build"])
	FuguFlags["build"].Name = "fugu"

	// Define FuguFlags["run"]
	FuguFlags["run"] = flags.New("fugu")
	FuguFlags["run"].String([]string{"-image"}, "", "Name of the image")
	FuguFlags["run"].String([]string{"-command"}, "", "COMMAND")
	FuguFlags["run"].Var([]string{"-arg"}, "ARG")
	FuguFlags["run"] = flags.Merge(FuguCommon, FuguFlags["run"])
	FuguFlags["run"].Name = "fugu"

	// Define FuguFlags["exec"]
	FuguFlags["exec"] = flags.New("fugu")
	FuguFlags["exec"].String([]string{"-name"}, "", "Name of the container")
	FuguFlags["exec"].String([]string{"-command"}, "", "COMMAND")
	FuguFlags["exec"].Var([]string{"-arg"}, "ARG")
	FuguFlags["exec"] = flags.Merge(FuguCommon, FuguFlags["exec"])
	FuguFlags["exec"].Name = "fugu"

	// Define FuguFlags["shell"]
	FuguFlags["shell"] = flags.New("fugu")
	FuguFlags["shell"].String([]string{"-name"}, "", "Name of the container")
	FuguFlags["shell"].String([]string{"-shell"}, "/bin/bash", "Path to shell")
	FuguFlags["shell"] = flags.Merge(FuguCommon, FuguFlags["shell"])
	FuguFlags["shell"].Name = "fugu"

	// Define FuguFlags["destroy"]
	FuguFlags["destroy"] = flags.New("fugu")
	FuguFlags["destroy"].String([]string{"-name"}, "", "Name of the container to be destroyed")
	FuguFlags["destroy"] = flags.Merge(FuguCommon, FuguFlags["destroy"])
	FuguFlags["destroy"].Name = "fugu"

	// Define FuguFlags["push"]
	FuguFlags["push"] = flags.New("fugu")
	FuguFlags["push"].String([]string{"-image"}, "", "Name of the image")
	FuguFlags["push"].String([]string{"-tag"}, "", "Push this tag of the image")
	FuguFlags["push"] = flags.Merge(FuguCommon, FuguFlags["push"])
	FuguFlags["push"].Name = "fugu"

	// Define FuguFlags["pull"]
	FuguFlags["pull"] = flags.New("fugu")
	FuguFlags["pull"].String([]string{"-image"}, "", "Name of the image")
	FuguFlags["pull"].String([]string{"-tag"}, "", "Pull this tag of the image")
	FuguFlags["pull"] = flags.Merge(FuguCommon, FuguFlags["pull"])
	FuguFlags["pull"].Name = "fugu"

	// Define FuguFlags["images"]
	FuguFlags["images"] = flags.New("fugu")
	FuguFlags["images"].String([]string{"-registry"}, "", "URL of the registry")
	FuguFlags["images"].String([]string{"-user"}, "", "Use this username")
	FuguFlags["images"].String([]string{"-password"}, "", "Use this password")
	FuguFlags["images"].String([]string{"-file"}, "~/.dockercfg", "Read credentials from this file")
	FuguFlags["images"].Bool([]string{"-password-stdin"}, false, "Ask for password")

	// Define FuguFlags["show-data"]
	FuguFlags["show-data"] = flags.New("fugu")
	FuguFlags["show-data"].Var([]string{"-source"}, "Get data from this source")

	// Define FuguFlags["show-labels"]
	FuguFlags["show-labels"] = flags.New("fugu")
	FuguFlags["show-labels"].Var([]string{"-source"}, "Get data from this source")
}
