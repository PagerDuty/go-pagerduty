package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type UserUpdate struct {
}

func UserUpdateCommand() (cli.Command, error) {
	return &UserUpdate{}, nil
}

func (c *UserUpdate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserUpdate) Synopsis() string {
	return "Update an existing user"
}

func (c *UserUpdate) Run(args []string) int {
	return 0
}
