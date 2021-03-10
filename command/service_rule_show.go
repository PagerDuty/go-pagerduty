package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ServiceRuleShow struct {
	Meta
}

func ServiceRuleShowCommand() (cli.Command, error) {
	return &ServiceRuleShow{}, nil
}

func (c *ServiceRuleShow) Help() string {
	helpText := `
  pd service rules show <SERVICE_ID> <RULE_ID> Show specific service
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *ServiceRuleShow) Synopsis() string {
	return "Get details about an existing service rule"
}

func (c *ServiceRuleShow) Run(args []string) int {
	flags := c.Meta.FlagSet("service rules show")
	flags.Usage = func() { fmt.Println(c.Help()) }
	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	if len(flags.Args()) != 2 {
		log.Error("Please specify service id and rule id")
		return -1
	}
	log.Info("Service id is:", flags.Arg(0))
	log.Info("Rule id is:", flags.Arg(1))

	client := c.Meta.Client()
	rule, err := client.GetServiceRule(context.Background(), flags.Arg(0), flags.Arg(1))
	if err != nil {
		log.Error(err)
		return -1
	}
	data, err := yaml.Marshal(rule)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(string(data))
	return 0
}
