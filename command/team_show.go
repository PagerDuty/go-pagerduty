package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type TeamShow struct {
	Meta
}

func TeamShowCommand() (cli.Command, error) {
	return &TeamShow{}, nil
}

func (c *TeamShow) Help() string {
	helpText := `
	team show Show details of an individual team

	Options:
		 -id     Team ID
	`
	return strings.TrimSpace(helpText)
}

func (c *TeamShow) Synopsis() string { return "Get details about an existing team" }

func (c *TeamShow) Run(args []string) int {
	var id string
	flags := c.Meta.FlagSet("team show")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&id, "id", "", "Team ID")
	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}

	if id == "" {
		log.Error("You must specify the team id using -id flag")
		return -1
	}
	client := c.Meta.Client()
	result, err := client.GetTeam(id)
	if err != nil {
		log.Error(err)
		return -1
	}
	data, err := yaml.Marshal(result)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(string(data))
	return 0
}
