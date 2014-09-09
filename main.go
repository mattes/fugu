package main

import (
	"fmt"
	"github.com/mattes/fugu/cli"
	"os"
)

var fugufileSearchFiles = []string{"fugu.yml", "fugu.yaml", ".fugu.yml", ".fugu.yaml"}

func main() {
	args := os.Args
	argsLen := len(args)
	if argsLen <= 1 {
		fmt.Println("no cmd")
		os.Exit(1)
	}
	command := args[1]

	// check command upfront
	if command != "run" && command != "build" {
		fmt.Println("invalid command")
		os.Exit(1)
	}

	// verfiy first two args if fugufile and/or label
	fugufilePath := ""
	fugufilePathGiven := false
	label := ""
	labelGiven := false

	// extract possible fugufile and possible label if possibly given
	args1 := make([]string, 0)
	if argsLen >= 4 {
		args1 = args[2:4]
	} else if argsLen >= 3 {
		args1 = args[2:3]
	}

	// find fugufile with given information
	var possibleFugufilePath string
	if len(args1) > 0 {
		possibleFugufilePath = args1[0]
	}
	fugufilePath, fugufilePathGiven = cli.FindFugufile(possibleFugufilePath, fugufileSearchFiles)
	if fugufilePath == "" {
		fmt.Println("Please create fugufile or specify path to fugufile.")
		os.Exit(1)
	}

	// get all labels
	labels, err := cli.GetAllLabels(fugufilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// check if label is given
	for _, l := range args1 {
		if isLabel(l, labels) {
			label = l
			labelGiven = true
			break
		}
	}

	// get remaining options
	offsetCount := 2
	if fugufilePathGiven {
		offsetCount += 1
	}
	if labelGiven {
		offsetCount += 1
	}
	dockerArgs := args[offsetCount:]

	// fmt.Println(dockerArgs)

	switch command {
	case "run":
		cli.CmdRun(fugufilePath, dockerArgs, label)
	}
}

func isLabel(label string, search []string) bool {
	for _, l := range search {
		if label == l {
			return true
		}
	}
	return false
}
