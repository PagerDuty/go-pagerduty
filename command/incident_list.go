package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type IncidentList struct {
	Meta
}

func IncidentListCommand() (cli.Command, error) {
	return &IncidentList{}, nil
}

func (c *IncidentList) Help() string {
	helpText := `
	pd incident list List incidents

	`
	return strings.TrimSpace(helpText)
}

func (c *IncidentList) Synopsis() string {
	return "List existing incidents"
}

func (c *IncidentList) Run(args []string) int {
	var teamIDs []string
	var timeZone string
	var sortBy string
	var includes []string
	flags := c.Meta.FlagSet("incident list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.Var((*ArrayFlags)(&includes), "include", "Additional details to include (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&teamIDs), "team-id", "Only show for team ID (can be specified multiple times)")
	flags.StringVar(&timeZone, "time-zone", "", "Time Zone")
	flags.StringVar(&sortBy, "sort-by", "", "sort by")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	opts := pagerduty.ListIncidentsOptions{
		TeamIDs:  teamIDs,
		TimeZone: timeZone,
		SortBy:   sortBy,
		Includes: includes,
	}
	if incidentList, err := client.ListIncidents(opts); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, incident := range incidentList.Incidents {
			fmt.Println("Entry: ", i+1)
			data, err := yaml.Marshal(incident)
			if err != nil {
				log.Error(err)
				return -1
			}
			fmt.Println(string(data))
		}
	}
	return 0
}
