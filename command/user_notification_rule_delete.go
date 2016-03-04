package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type UserNotificationRuleDelete struct {
}

func UserNotificationRuleDeleteCommand() (cli.Command, error) {
	return &UserNotificationRuleDelete{}, nil
}

func (c *UserNotificationRuleDelete) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserNotificationRuleDelete) Synopsis() string {
	return "Remove a user's notification rule"
}

func (c *UserNotificationRuleDelete) Run(args []string) int {
	return 0
}
