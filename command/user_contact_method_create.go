package cmd

import (
	"strings"

	"github.com/mitchellh/cli"
)

type UserContactMethodCreate struct {
}

func UserContactMethodCreateCommand() (cli.Command, error) {
	return &UserContactMethodCreate{}, nil
}

func (c *UserContactMethodCreate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *UserContactMethodCreate) Synopsis() string {
	return "Create a new contact method"
}

func (c *UserContactMethodCreate) Run(args []string) int {
	return 0
}
