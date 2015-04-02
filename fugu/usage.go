package main

import (
	"fmt"
	"github.com/mattes/go-collect"
	"github.com/mattes/go-collect/flags"
	"os"
	"strings"
)

func usage(c *collect.Collector, command string) {

	switch command {
	case "build":
		printMulti(`
    Usage: fugu build [LABEL] [OPTIONS] [PATH | URL]

    Build a new image from the source code at PATH`)

		c.PrintUsage()
		printSourceExampleUrls(c)

	case "run":
		printMulti(`
    Usage: fugu run [LABEL] [OPTIONS] [COMMAND] [ARG...]

    Run a command in a new container`)

		c.PrintUsage()
		printSourceExampleUrls(c)

	case "exec":
		printMulti(`
    Usage: fugu exec [LABEL] [OPTIONS] [COMMAND] [ARG...]

    Run a command in a running container`)

		c.PrintUsage()
		printSourceExampleUrls(c)

	case "shell":
		printMulti(`
    Usage: fugu shell [LABEL] [OPTIONS]

    Open a shell in a running container`)

		c.PrintUsage()
		printSourceExampleUrls(c)

	case "destroy":
		printMulti(`
    Usage: fugu destroy [LABEL]

    Kil a running container and remove it`)

		c.PrintUsage()
		printSourceExampleUrls(c)

	case "push":
		printMulti(`
    Usage: fugu push [LABEL] [OPTIONS] [TAG]

    Push an image or a repository to the registry`)

		c.PrintUsage()
		printSourceExampleUrls(c)

	case "pull":
		printMulti(`
    Usage: fugu pull [OPTIONS] [TAG]

    Pull an image or a repository from the registry`)

		c.PrintUsage()
		printSourceExampleUrls(c)

	case "images":
		printMulti(`
    Usage: fugu images [REGISTRY]

    List images (from remote registry)`)

		c.PrintUsage()
		printSourceExampleUrls(c)

	case "show-data":
		printMulti(`
    Usage: fugu show-data [LABEL] [OPTIONS]

    Show aggregated data for label`)

		c.PrintUsage()
		printSourceExampleUrls(c)

	case "show-labels":
		printMulti(`
    Usage: fugu show-labels [OPTIONS]

    Show all labels`)

		c.PrintUsage()
		printSourceExampleUrls(c)

	default:
		printMulti(`
    Usage: fugu COMMAND [LABEL] [arg...]

    Swiss Army knife for Docker.

    Commands:
        build        Build an image from a Dockerfile
        run          Run a command in a new container
        exec         Run a command in a running container
        shell        Open a shell in a running container
        destroy      Kil a running container and remove it
        push         Push an image or a repository to the registry
        pull         Pull an image or a repository from the registry
        images       List images (from remote registry)
        show-data    Show aggregated data for label
        show-labels  Show all labels
        help         Show help

    Run 'fugu help COMMAND' for more information on a command.`)
	}

}

func printMulti(out string) {
	out = strings.TrimSpace(out)
	for _, l := range strings.Split(out, "\n") {
		fmt.Fprintln(os.Stderr, strings.TrimPrefix(l, "    "))
	}
}

func printSourceExampleUrls(c *collect.Collector) {
	ex := collect.SourceExampleUrls()
	if len(ex) > 0 {
		fmt.Fprintln(os.Stderr, "\nExample source options:")

		if c.GetDefaultSource() != "" {
			fmt.Fprintln(os.Stderr, "  "+flags.Nice("source", c.GetDefaultSource())+" (default)")
		}

		for _, v := range ex {
			fmt.Fprintln(os.Stderr, "  "+v)
		}
	}
}
