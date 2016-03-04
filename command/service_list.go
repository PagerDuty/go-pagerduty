package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ServiceList struct {
}

func ServiceListCommand() (cli.Command, error) {
	return &ServiceList{}, nil
}

func (c *ServiceList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ServiceList) Synopsis() string {
	return "List existing services"
}

func (c *ServiceList) Run(args []string) int {
	return 0
}
