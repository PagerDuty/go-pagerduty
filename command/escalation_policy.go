package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"strings"
)

type EpsLs struct {
	Meta
}

type ArrayFlags []string

func (a *ArrayFlags) String() string {
	return strings.Join(*a, ",")
}

func (a *ArrayFlags) Set(v string) error {
	if *a == nil {
		*a = make([]string, 0, 1)
	}
	*a = append(*a, v)
	return nil
}

func (c *EpsLs) Help() string {
	helpText := `
	ep ls  List escalation policies

	Options:

		 -query
		 -user
		 -team
		 -include
		 -sort

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EpsLs) Synopsis() string {
	return "List escalation policies"
}

func EpsLsCommand() (cli.Command, error) {
	return &EpsLs{}, nil
}

func (c *EpsLs) Run(args []string) int {
	var query string
	var sortBy string
	var userIDs []string
	var teamIDs []string
	var includes []string

	flags := c.Meta.FlagSet("eps ls")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&query, "query", "", "Show escalation policies whose names contain the query")
	flags.StringVar(&sortBy, "sort", "", "Sort results by name (ascending or descending)")
	flags.Var((*ArrayFlags)(&userIDs), "user", "Filter escalation policies by user ids (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&teamIDs), "team", "Filter escalation policies by team ids (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&includes), "include", "Additional details to include (can be specified multiple times)")

	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}
	if err := c.Meta.Validate(); err != nil {
		log.Error(err)
		return -1
	}
	c.Meta.SetupLogging()
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
		for _, p := range eps.EscalationPolicies {
			log.Info(p.Name, p.ID)
		}
	}
	return 0
}

type EpShow struct {
	Meta
}

func (c *EpShow) Help() string {
	helpText := `
	ep show  show escalation policies

	Options:

		 -id
		 -include

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EpShow) Synopsis() string {
	return "Show escalation policy"
}

func EpShowCommand() (cli.Command, error) {
	return &EpShow{}, nil
}

func (c *EpShow) Run(args []string) int {
	var includes []string
	var epID string

	flags := c.Meta.FlagSet("ep show")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&epID, "id", "", "Escalation policy id")
	flags.Var((*ArrayFlags)(&includes), "include", "Additional details to include (can be specified multiple times)")

	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}
	if err := c.Meta.Validate(); err != nil {
		log.Error(err)
		return -1
	}
	if epID == "" {
		log.Error("You must provide escalation policy id")
		return -1
	}
	c.Meta.SetupLogging()
	client := c.Meta.Client()
	o := &pagerduty.GetEscalationPolicyOptions{
		Includes: includes,
	}
	ep, err := client.GetEscalationPolicy(epID, o)
	if err != nil {
		log.Error(err)
		return -1
	}
	log.Println(ep)
	return 0
}
