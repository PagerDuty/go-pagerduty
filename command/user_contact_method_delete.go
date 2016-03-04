package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type UserContactMethodDelete struct {
}

func UserContactMethodDeleteCommand() (cli.Command, error) {
	return &UserContactMethodDelete{}, nil
}

func (c *UserContactMethodDelete) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserContactMethodDelete) Synopsis() string {
	return "Remove a user's contact method"
}

func (c *UserContactMethodDelete) Run(args []string) int {
	return 0
}
