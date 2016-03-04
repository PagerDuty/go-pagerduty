package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type IncidentNoteList struct {
}

func IncidentNoteListCommand() (cli.Command, error) {
	return &IncidentNoteList{}, nil
}

func (c *IncidentNoteList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *IncidentNoteList) Synopsis() string {
	return "List existing notes for the specified incident"
}

func (c *IncidentNoteList) Run(args []string) int {
	return 0
}
