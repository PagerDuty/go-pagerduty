package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ServiceShow struct {
}

func ServiceShowCommand() (cli.Command, error) {
	return &ServiceShow{}, nil
}

func (c *ServiceShow) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ServiceShow) Synopsis() string {
	return "Get details about an existing service"
}

func (c *ServiceShow) Run(args []string) int {
	return 0
}
