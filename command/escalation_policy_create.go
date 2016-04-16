package main

import (
	"encoding/json"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

type EscalationPolicyCreate struct {
	Meta
}

func EscalationPolicyCreateCommand() (cli.Command, error) {
	return &EscalationPolicyCreate{}, nil
}

func (c *EscalationPolicyCreate) Help() string {
	helpText := `
	pd escalation-policy create <FILE> Create a new escalation policy
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EscalationPolicyCreate) Synopsis() string {
	return "Creates a new escalation policy"
}

func (c *EscalationPolicyCreate) Run(args []string) int {
	flags := c.Meta.FlagSet("escalation-policy create")
	flags.Usage = func() { fmt.Println(c.Help()) }
	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	var ep pagerduty.EscalationPolicy
	if len(flags.Args()) != 1 {
		log.Error("Please specify input json file")
		return -1
	}
	log.Info("Input file is:", flags.Arg(0))
	f, err := os.Open(flags.Arg(0))
	if err != nil {
		log.Error(err)
		return -1
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&ep); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}
	log.Debugf("%#v", ep)
	if err := client.CreateEscalationPolicy(ep); err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
