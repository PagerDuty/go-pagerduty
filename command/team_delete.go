package main

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type TeamDelete struct {
	Meta
}

func TeamDeleteCommand() (cli.Command, error) {
	return &TeamDelete{}, nil
}

func (c *TeamDelete) Help() string {
	helpText := `
	team delete Delete a team

	Option:
		-id     ID of the team
	`
	return strings.TrimSpace(helpText)
}

func (c *TeamDelete) Synopsis() string {
	return "Remove an existing team"
}

func (c *TeamDelete) Run(args []string) int {
	var id string
	flags := c.Meta.FlagSet("team delete")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&id, "id", "", "ID of the team")
	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}

	client := c.Meta.Client()
	err := client.DeleteTeam(id)
	if err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
