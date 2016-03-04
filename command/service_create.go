package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ServiceCreate struct {
}

func ServiceCreateCommand() (cli.Command, error) {
	return &ServiceCreate{}, nil
}

func (c *ServiceCreate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ServiceCreate) Synopsis() string {
	return "Create a new service"
}

func (c *ServiceCreate) Run(args []string) int {
	return 0
}
