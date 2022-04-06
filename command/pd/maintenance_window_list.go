package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type MaintenanceWindowList struct {
	Meta
}

func MaintenanceWindowListCommand() (cli.Command, error) {
	return &MaintenanceWindowList{}, nil
}

func (c *MaintenanceWindowList) Help() string {
	helpText := `
	pd maintenance-window list List maintenance windows

	Options:
		-filter     Filter result by state
		-include    Additional details to include
		-query      Filter results by query
		-service-id Filter result by service ids
		-team-id    Filter results by team ids
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *MaintenanceWindowList) Synopsis() string {
	return "List existing maintenance windows"
}

func (c *MaintenanceWindowList) Run(args []string) int {
	var includes, serviceIDs, teamIDs []string
	var query, filter string
	flags := c.Meta.FlagSet("maintenance-window list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.Var((*ArrayFlags)(&includes), "includes", "Additional details to include (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&serviceIDs), "service-id", "Show maintenance windows for the specified services only (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&teamIDs), "team-id", "Show maintenance windows for the specified teams only (can be specified multiple times)")
	flags.StringVar(&filter, "filter", "all", "Filter results by maintenance window state (past, future, ongoing, open, all)")
	flags.StringVar(&query, "query", "", "Filter results showing only tags whose labels match the query")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}

	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}

	client := c.Meta.Client()
	opts := pagerduty.ListMaintenanceWindowsOptions{
		Query:      query,
		Includes:   includes,
		TeamIDs:    teamIDs,
		ServiceIDs: serviceIDs,
		Filter:     filter,
	}

	mws, err := client.ListMaintenanceWindowsWithContext(context.Background(), opts)
	if err != nil {
		log.Error(err)
		return -1
	}

	for i, mw := range mws.MaintenanceWindows {
		fmt.Println("Entry: ", i)
		data, err := yaml.Marshal(mw)
		if err != nil {
			log.Error(err)
			return -1
		}
		fmt.Println(string(data))
	}

	return 0
}
