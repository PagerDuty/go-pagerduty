package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type TeamCreate struct {
	Meta
}

func TeamCreateCommand() (cli.Command, error) {
	return &TeamCreate{}, nil
}

func (c *TeamCreate) Help() string {
	helpText := `
	team create Create a team

	Options:
		-name           Name of the team
		-description    Description about the team
	`
	return strings.TrimSpace(helpText)
}

func (c *TeamCreate) Synopsis() string {
	return "Create a new team"
}

func (c *TeamCreate) Run(args []string) int {
	var name, desc string
	flags := c.Meta.FlagSet("team create")
	flags.Usage = func() { fmt.Println(c.Help()) }
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
	result, err := client.CreateTeam(team)
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
