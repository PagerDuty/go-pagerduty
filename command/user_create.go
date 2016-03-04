package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type UserCreate struct {
}

func UserCreateCommand() (cli.Command, error) {
	return &UserCreate{}, nil
}

func (c *UserCreate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserCreate) Synopsis() string {
	return "Create a new user"
}

func (c *UserCreate) Run(args []string) int {
	return 0
}
