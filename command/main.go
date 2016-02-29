package main

import (
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

/*
	Entrypoint for pd cli client based on the go bindings of API.
	good cli examples:
		https://github.com/Netflix-Skunkworks/go-jira
		https://github.com/docker/docker
		https://github.com/hashicorp/serf
*/

const (
	version = "0.1"
)

func loadCommands() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"ep list": EpLsCommand,
		"ep show": EpShowCommand,
	}
}

func main() {
	os.Exit(invokeCLI())
}

func invokeCLI() int {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: loadCommands(),
		HelpFunc: cli.BasicHelpFunc("pd"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
