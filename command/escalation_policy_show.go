package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type EscalationPolicyShow struct {
	Meta
}

func (c *EscalationPolicyShow) Help() string {
	helpText := `
	escalation-policy show  Show escalation policies

	Options:

		 -id
		 -include

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EscalationPolicyShow) Synopsis() string {
	return "Show information about an existing escalation policy and its rules"
}

func EscalationPolicyShowCommand() (cli.Command, error) {
	return &EscalationPolicyShow{}, nil
}

func (c *EscalationPolicyShow) Run(args []string) int {
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
