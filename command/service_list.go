package main

import (
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ServiceList struct {
	Meta
}

func ServiceListCommand() (cli.Command, error) {
	return &ServiceList{}, nil
}

func (c *ServiceList) Help() string {
	helpText := `
	pd service list List services

	Options:
		 -include    Include additional details
		 -team-id    Filter result by team (can be specified multiple times)
		 -sort-by    Sort result (name:asc, name:dsc)
		 -query      Filter result by pattern (name or service key(
	`
	return strings.TrimSpace(helpText)
}

func (c *ServiceList) Synopsis() string {
	return "List existing services"
}

func (c *ServiceList) Run(args []string) int {
	var teamIDs []string
	var timeZone string
	var sortBy string
	var query string
	var includes []string
	flags := c.Meta.FlagSet("service list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.Var((*ArrayFlags)(&includes), "include", "Additional details to include (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&teamIDs), "team-id", "Only show for team ID (can be specified multiple times)")
	flags.StringVar(&timeZone, "time-zone", "", "Time Zone")
	flags.StringVar(&sortBy, "sort-by", "", "sort by")
	flags.StringVar(&query, "query", "", "Query")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	opts := pagerduty.ListServiceOptions{
		TeamIDs:  teamIDs,
		TimeZone: timeZone,
		SortBy:   sortBy,
		Query:    query,
		Includes: includes,
	}
	if serviceList, err := client.ListServices(opts); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, service := range serviceList.Services {
			fmt.Println("Entry: ", i+1)
			data, err := yaml.Marshal(service)
			if err != nil {
				log.Error(err)
				return -1
			}
			fmt.Println(string(data))
		}
	}
	return 0
}
