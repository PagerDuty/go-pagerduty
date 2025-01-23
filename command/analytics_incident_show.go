package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type AnalyticsShow struct {
	Meta
}

// AnalyticsGetAggregatedIncidentDataCommand gets the aggregated incident analytics for the requested data.
func AnalyticsGetAggregatedIncidentDataCommand() (cli.Command, error) {
	return &AnalyticsShow{}, nil
}

// Help displays information on how to use the analytics incident cli command.
func (c *AnalyticsShow) Help() string {
	helpText := `
	analytics incident show

	Options:

		 -(service_id|team_id) #MANDATORY provide service or team id to stats on.
		 -start #Optional RFC3339 format default : 7 days ago
		 -end #Optional RFC3339 format default : now
		 -urgency #Optional (high|low)
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

// Synopsis returns a synopsis of the analytics incident cli command.
func (c *AnalyticsShow) Synopsis() string {
	return "Get aggregated incident data analytics"
}

//Run executes analytics cli command and displays incident analytics for the requested data.
func (c *AnalyticsShow) Run(args []string) int {
	flags := c.Meta.FlagSet("analytics incident show")
	flags.Usage = func() { fmt.Println(c.Help()) }
	servID := flags.String("service_id", "", "Service ID")
	now := time.Now()
	sevenDaysAgo := now.Add(time.Duration(-24*7) * time.Hour)
	start := flags.String("start", sevenDaysAgo.Format(time.RFC3339), "start date")
	end := flags.String("end", now.Format(time.RFC3339), "end date")
	urgency := flags.String("urgency", "", "high|low")
	teamID := flags.String("team_id", "", "team ID")

	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}

	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}

	client := c.Meta.Client()

	serviceIds := make([]string, 1)
	if *servID == "" {
		serviceIds = nil
	} else {
		serviceIds[0] = *servID
	}

	teamIds := make([]string, 1)
	if *teamID == "" {
		teamIds = nil
	} else {
		teamIds[0] = *teamID
	}

	analyticsFilter := pagerduty.AnalyticsFilter{
		CreatedAtStart: *start,
		CreatedAtEnd:   *end,
		Urgency:        *urgency,
		ServiceIDs:     serviceIds,
		TeamIDs:        teamIds,
	}

	analytics := pagerduty.AnalyticsRequest{
		Filters:       &analyticsFilter,
		AggregateUnit: "day",
		TimeZone:      "Etc/UTC",
	}

	aggregatedIncidentData, err := client.GetAggregatedIncidentData(context.Background(), analytics)
	if err != nil {
		log.Error(err)
		return -1
	}

	aggregatedIncidentDataBytes, err := json.Marshal(aggregatedIncidentData)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(string(aggregatedIncidentDataBytes))
	return 0
}
