package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type IncidentNoteCreate struct {
}

func IncidentNoteCreateCommand() (cli.Command, error) {
	return &IncidentNoteCreate{}, nil
}

func (c *IncidentNoteCreate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *IncidentNoteCreate) Synopsis() string {
	return "Create a new note for the specified incident"
}

func (c *IncidentNoteCreate) Run(args []string) int {
	return 0
}
