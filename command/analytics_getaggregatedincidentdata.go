package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type AnalyticsAggregatedIncidentData struct {
	Meta
}

func AnalyticsGetAggregatedIncidentDataCommand() (cli.Command, error) {
	return &AnalyticsAggregatedIncidentData{}, nil
}

func (c *AnalyticsAggregatedIncidentData) Help() string {
	helpText := `
	AnalyticsAggregatedIncidentData show

	Options:

		 -(service_id|team_id) #provide service or team id to stats on.
		 -start #RFC3339 format default : 7 days ago
		 -end #RFC3339 format default : now
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *AnalyticsAggregatedIncidentData) Synopsis() string {
	return "Get analytics aggregated incident data"
}

func (c *AnalyticsAggregatedIncidentData) Run(args []string) int {
	flags := c.Meta.FlagSet("analyticsAggregatedIncidentData show")
	flags.Usage = func() { fmt.Println(c.Help()) }
	servID := flags.String("service_id", "", "Service ID")
	now := time.Now()
	sevenDaysAgo :=now.Add(time.Duration(-24*7)*time.Hour)
	start := flags.String("start", sevenDaysAgo.Format(time.RFC3339), "start date")
	end := flags.String("end", now.Format(time.RFC3339), "end date")
	urgency := flags.String("urgency","", "high|low")
	teamID := flags.String("team_id","", "team ID")
	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	serviceIds := make([]string,1)
	if *servID == "" {
		serviceIds = nil
	} else {
		serviceIds[0] = *servID
	}
	teamIds := make([]string,1)
	if *teamID=="" {
		teamIds=nil
	} else {
		teamIds[0] = *teamID
	}
	analyticsFilter := pagerduty.AnalyticsFilter{
		CreatedAtStart: *start,
		CreatedAtEnd:   *end,
		Urgency:        *urgency,
		ServiceIds:     serviceIds,
		TeamIds: teamIds,
	}
	analytics := pagerduty.Analytics{
		AnalyticsFilter: &analyticsFilter,
		AggregateUnit:   "day",
		TimeZone:        "Etc/UTC",
	}
	aggregatedIncidentData, err := client.GetAggregatedIncidentData(analytics)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(aggregatedIncidentData)
	return 0
}
