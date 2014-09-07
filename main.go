package main

import (
	_ "flag"
	"fmt"
	"github.com/mattes/fugu/file"
	"os"
	"path"
	_ "strings"
)

var fileNamePaths = []string{"fugu.yml", "fugu.yaml", ".fugu.yml", ".fugu.yaml"}

// fugu run fugu.yml
// fugu run fugu.yml label
// fugu run fugu.yml label --rm
// fugu run label
// fugu run --rm
// fugu run label --rm

// fugu run echo "hello world!"

// fugu run [fugu.yml-path] [label] [docker-run-options + --image] [command] [args]
// fugu build [fugu.yml-path] [label] [docker-build-options] [path=pwd|url|-]

func main() {
	if len(fileNamePaths) == 0 {
		panic("Specify at least one fileNamePath!")
	}

	args := os.Args
	argsLen := len(args)

	if argsLen <= 1 {
		fmt.Println("no cmd")
		os.Exit(1)
	}

	command := args[1]
	fugufilePath := ""
	fugufilePathGiven := false
	label := ""
	labelGiven := false

	args1 := make([]string, 0)
	if argsLen >= 4 {
		args1 = args[2:4]
	} else if argsLen >= 3 {
		args1 = args[2:3]
	}

	// check if fugufile
	exts := fileExtensions(fileNamePaths)
	for _, f := range args1 {
		if isFuguFile(f, exts) {
			fugufilePath = f
			fugufilePathGiven = true
			break
		}
	}

	if fugufilePath == "" {
		for _, f := range fileNamePaths {
			if fileExists(f) {
				fugufilePath = f
				break
			}
		}
		if fugufilePath == "" {
			fmt.Printf("%s not found\n", fileNamePaths[0])
			os.Exit(1)
		}

	} else {
		if !fileExists(fugufilePath) {
			fmt.Printf("%s not found\n", fugufilePath)
			os.Exit(1)
		}
	}

	// get all labels
	labels, err := file.GetLabels(fugufilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// check if label
	// dont worry here, if no label is found
	for _, l := range args1 {
		if isLabel(l, labels) {
			label = l
			labelGiven = true
			break
		}
	}

	// load fugu config
	fuguConfig, err := file.GetConfig(fugufilePath, label)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_ = fuguConfig

	// now do docker option parsing
	offsetCount := 2
	if fugufilePathGiven {
		offsetCount += 1
	}
	if labelGiven {
		offsetCount += 1
	}

	dockerArgs := []string{} // []string{strings.TrimSpace(fuguConfig["image"].(string))}
	if argsLen >= offsetCount {
		dockerArgs = append(dockerArgs, args[offsetCount:]...)
	}

	fmt.Println(dockerArgs)

	switch command {
	case "run":

	case "build":

	default:
		os.Exit(1)
	}
}

func isFuguFile(file string, extensions []string) bool {
	for _, ext := range extensions {
		if path.Ext(file) == ext {
			return true
		}
	}
	return false
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func isLabel(label string, search []string) bool {
	for _, l := range search {
		if label == l {
			return true
		}
	}
	return false
}

func fileExtensions(s []string) []string {
	su := make([]string, 0)
	exts := make(map[string]bool)
	for _, v := range s {
		ext := path.Ext(v)
		if _, ok := exts[ext]; !ok {
			exts[ext] = true
			su = append(su, ext)
		}
	}
	return su
}
