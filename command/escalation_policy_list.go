package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type EscalationPolicyList struct {
	Meta
}

func (c *EscalationPolicyList) Help() string {
	helpText := `
	pd escalation-policy list List all of the existing escalation policies

	Options:

		 -query     Filter escalation policies with certain name
		 -user      Filter escalation policies by user id(s)
		 -team      Filter escalation policies by team id(s)
		 -include   Additional details to include
		 -sort      Sort results by property (name:asc or name:dsc)

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EscalationPolicyList) Synopsis() string {
	return "List all of the existing escalation policies"
}

func EscalationPolicyListCommand() (cli.Command, error) {
	return &EscalationPolicyList{}, nil
}

func (c *EscalationPolicyList) Run(args []string) int {
	var query string
	var sortBy string
	var userIDs []string
	var teamIDs []string
	var includes []string

	flags := c.Meta.FlagSet("escalation-policy list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&query, "query", "", "Show escalation policies whose names contain the query")
	flags.StringVar(&sortBy, "sort", "", "Sort results by name (ascending or descending)")
	flags.Var((*ArrayFlags)(&userIDs), "user", "Filter escalation policies by user ids (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&teamIDs), "team", "Filter escalation policies by team ids (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&includes), "include", "Additional details to include (can be specified multiple times)")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	opts := pagerduty.ListEscalationPoliciesOptions{
		Query:    query,
		UserIDs:  userIDs,
		TeamIDs:  teamIDs,
		Includes: includes,
		SortBy:   sortBy,
	}
	if eps, err := client.ListEscalationPolicies(opts); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, p := range eps.EscalationPolicies {
			fmt.Println("Entry: ", i)
			data, err := yaml.Marshal(p)
			if err != nil {
				log.Error(err)
				return -1
			}
			fmt.Println(string(data))
		}
	}
	return 0
}
