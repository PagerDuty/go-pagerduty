package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type EpList struct {
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

func (c *EpList) Help() string {
	helpText := `
	ep ls  List escalation policies

	Options:

		 -query     Filter escalation policies with certain name
		 -user      Filter escalation policies by user id(s)
		 -team      Filter escalation policies by team id(s)
		 -include   Additional details to include
		 -sort      Sort results by property (name:asc or name:dsc)

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EpList) Synopsis() string {
	return "List escalation policies"
}

func EpListCommand() (cli.Command, error) {
	return &EpList{}, nil
}

func (c *EpList) Run(args []string) int {
	var query string
	var sortBy string
	var userIDs []string
	var teamIDs []string
	var includes []string

	flags := c.Meta.FlagSet("ep list")
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
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	if epID == "" {
		log.Error("You must provide escalation policy id")
		return -1
	}
	client := c.Meta.Client()
	o := &pagerduty.GetEscalationPolicyOptions{
		Includes: includes,
	}
	ep, err := client.GetEscalationPolicy(epID, o)
	if err != nil {
		log.Error(err)
		return -1
	}
	data, err := yaml.Marshal(ep)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(string(data))
	return 0
}
