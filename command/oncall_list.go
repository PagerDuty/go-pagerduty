package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type OncallList struct {
}

func OncallListCommand() (cli.Command, error) {
	return &OncallList{}, nil
}

func (c *OncallList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *OncallList) Synopsis() string {
	return "List the on-call entries during a given time range"
}

func (c *OncallList) Run(args []string) int {
	return 0
}
