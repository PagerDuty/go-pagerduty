package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type UserShow struct {
}

func UserShowCommand() (cli.Command, error) {
	return &UserShow{}, nil
}

func (c *UserShow) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserShow) Synopsis() string {
	return "Get details about an existing user"
}

func (c *UserShow) Run(args []string) int {
	return 0
}
