package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type UserContactMethodUpdate struct {
}

func UserContactMethodUpdateCommand() (cli.Command, error) {
	return &UserContactMethodUpdate{}, nil
}

func (c *UserContactMethodUpdate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserContactMethodUpdate) Synopsis() string {
	return "Update a user's contact method"
}

func (c *UserContactMethodUpdate) Run(args []string) int {
	return 0
}
