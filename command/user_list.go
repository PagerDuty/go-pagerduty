package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type UserList struct {
}

func UserListCommand() (cli.Command, error) {
	return &UserList{}, nil
}

func (c *UserList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserList) Synopsis() string {
	return "List users of your PagerDuty account, optionally filtered by a search query"
}

func (c *UserList) Run(args []string) int {
	return 0
}
