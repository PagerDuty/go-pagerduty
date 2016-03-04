package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type TeamList struct {
	Meta
}

func TeamListCommand() (cli.Command, error) {
	return &TeamList{}, nil
}

func (c *TeamList) Help() string {
	helpText := `
	team list List teams

	Options:
		 -query     Filter teams with certain name or email address
	`
	return strings.TrimSpace(helpText)
}

func (c *TeamList) Synopsis() string {
	return "List teams of your PagerDuty account, optionally filtered by a search query"
}

func (c *TeamList) Run(args []string) int {
	var query string
	flags := c.Meta.FlagSet("team list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&query, "query", "", "Show teams whose names or email address contain the query")
	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	o := pagerduty.ListTeamOptions{
		Query: query,
	}
	result, err := client.ListTeams(o)
	if err != nil {
		log.Error(err)
		return -1
	}
	for i, p := range result.Teams {
		fmt.Println("Entry: ", i)
		data, err := yaml.Marshal(p)
		if err != nil {
			log.Error(err)
			return -1
		}
		fmt.Println(string(data))
	}
	return 0
}
