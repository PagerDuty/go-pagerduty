package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type UserContactMethodShow struct {
}

func UserContactMethodShowCommand() (cli.Command, error) {
	return &UserContactMethodShow{}, nil
}

func (c *UserContactMethodShow) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserContactMethodShow) Synopsis() string {
	return "Get details about a user's contact method"
}

func (c *UserContactMethodShow) Run(args []string) int {
	return 0
}
