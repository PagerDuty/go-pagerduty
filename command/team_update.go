package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type TeamUpdate struct {
	Meta
}

func TeamUpdateCommand() (cli.Command, error) {
	return &TeamUpdate{}, nil
}

func (c *TeamUpdate) Help() string {
	helpText := `
	team update Update a team

	Options:
		-id             ID of the team to update
		-name           Name of the team
		-description    Description about the team
	`
	return strings.TrimSpace(helpText)
}

func (c *TeamUpdate) Synopsis() string {
	return "Update an existing team"
}

func (c *TeamUpdate) Run(args []string) int {
	var id, name, desc string
	flags := c.Meta.FlagSet("team update")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&id, "id", "", "ID of the team")
	flags.StringVar(&name, "name", "", "Name of the team")
	flags.StringVar(&desc, "description", "", "Description about the team")
	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	team := &pagerduty.Team{
		Name:        name,
		Description: desc,
	}
	client := c.Meta.Client()
	result, err := client.UpdateTeam(id, team)
	if err != nil {
		log.Error(err)
		return -1
	}
	data, err := json.Marshal(result)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(string(data))
	return 0
}
