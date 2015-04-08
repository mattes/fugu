package main

import (
	"fmt"
	"github.com/mattes/fugu"
	"github.com/mattes/go-collect"
	"os"

	// Import sources here and register in main()
	fileSource "github.com/mattes/go-collect/source/file"
)

var FugufileSearchpaths = []string{
	"fugu.yml", "fugu.yaml", ".fugu.yml", ".fugu.yaml"}

func main() {
	// Register sources ...
	collect.RegisterSource(&fileSource.File{})

	var args = os.Args[1:]

	// get command
	var command string
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	}

	// create new collect object
	c := collect.New()

	// set default source if fugu file is found in search path
	for _, p := range FugufileSearchpaths {
		if _, err := os.Stat(p); err == nil {
			// TODO expand ~ in p
			c.SetDefaultSource("file://" + p)
			break
		}
	}

	switch command {
	case "--version":
		fmt.Println(Version)
		os.Exit(0)

	case "":
		fallthrough
	case "--help":
		usage(c, "")
		os.Exit(0)

	case "help":
		if len(args) > 0 {
			command = args[0]
			c.AddFlags(
				// order determines order in printUsage()
				fugu.FuguFlags[command],
				fugu.DockerFlags[command],
			)
			usage(c, command)
		} else {
			usage(c, "")
		}
		os.Exit(0)

	case "build":
		fallthrough
	case "run":
		fallthrough
	case "exec":
		fallthrough
	case "shell":
		fallthrough
	case "destroy":
		fallthrough
	case "push":
		fallthrough
	case "pull":
		dockerCommand(c, command, args)

	case "images":
		fallthrough
	case "show-data":
		fallthrough
	case "show-labels":
		fuguCommand(c, command, args)

	default:
		usage(c, "")
		fmt.Println()
		fuguErrExit("unkown command")
	}
}

func dockerCommand(c *collect.Collector, command string, args []string) {
	data, remainingArgs, err := c.Parse(args,
		fugu.FuguFlags[command],
		fugu.DockerFlags[command],
	)

	if err != nil {
		fuguErrExit(err)
	}

	if data.IsTrue("help") {
		usage(c, command)
		os.Exit(0)
	}

	cmdStr, err := fugu.DockerCommands[command](c, data, remainingArgs)
	if err != nil {
		fuguErrExit(err)
	}

	if data.IsTrue("dry-run") {
		fmt.Println(cmdStr)
	} else {
		fugu.DockerExec(cmdStr, true)
	}
}

func fuguCommand(c *collect.Collector, command string, args []string) {
	data, remainingArgs, err := c.Parse(args,
		fugu.FuguFlags[command],
		fugu.DockerFlags[command],
	)
	if err != nil {
		fuguErrExit(err)
	}

	if data.IsTrue("help") {
		usage(c, command)
		os.Exit(0)
	}

	if err := fugu.Commands[command](c, data, remainingArgs); err != nil {
		fuguErrExit(err)
	}
}

func fuguErrExit(msg interface{}) {
	if msg != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", msg)
	}
	os.Exit(1)
}
