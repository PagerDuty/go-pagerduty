package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type UserContactMethodList struct {
}

func UserContactMethodListCommand() (cli.Command, error) {
	return &UserContactMethodList{}, nil
}

func (c *UserContactMethodList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserContactMethodList) Synopsis() string {
	return "List contact methods of your PagerDuty user"
}

func (c *UserContactMethodList) Run(args []string) int {
	return 0
}
