package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type TeamAddUser struct {
}

func TeamAddUserCommand() (cli.Command, error) {
	return &TeamAddUser{}, nil
}

func (c *TeamAddUser) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *TeamAddUser) Synopsis() string {
	return "Add a user to a team"
}

func (c *TeamAddUser) Run(args []string) int {
	return 0
}
