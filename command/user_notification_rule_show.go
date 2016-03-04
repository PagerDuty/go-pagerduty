package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type UserNotificationRuleShow struct {
}

func UserNotificationRuleShowCommand() (cli.Command, error) {
	return &UserNotificationRuleShow{}, nil
}

func (c *UserNotificationRuleShow) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserNotificationRuleShow) Synopsis() string {
	return "Get details about a user's notification rule"
}

func (c *UserNotificationRuleShow) Run(args []string) int {
	return 0
}
