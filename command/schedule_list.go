package main

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
)

type ScheduleList struct {
	Meta
}

func ScheduleListCommand() (cli.Command, error) {
	return &ScheduleList{}, nil
}

func (c *ScheduleList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ScheduleList) Synopsis() string {
	return "List the on-call schedules"
}

func (c *ScheduleList) Run(args []string) int {
	var query string

	flags := c.Meta.FlagSet("Schedule list")
	flags.Usage = func() { fmt.Println(c.Help()) }
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
	opts := pagerduty.ListSchedulesOptions{
		Query: query,
	}
	if scheduleList, err := client.ListSchedules(opts); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, schedule := range scheduleList.Schedules {
			fmt.Println("Entry: ", i+1)
			data, err := yaml.Marshal(schedule)
			if err != nil {
				log.Error(err)
				return -1
			}
			fmt.Println(string(data))
		}
	}
	return 0
}
