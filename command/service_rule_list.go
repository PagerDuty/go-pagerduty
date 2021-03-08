package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ServiceRuleList struct {
	Meta
}

func ServiceRuleListCommand() (cli.Command, error) {
	return &ServiceRuleList{}, nil
}

func (c *ServiceRuleList) Help() string {
	helpText := `
  pd service rules list <SERVICE_ID> List rules for service
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *ServiceRuleList) Synopsis() string {
	return "List existing Rules for a service"
}

func (c *ServiceRuleList) Run(args []string) int {
	flags := c.Meta.FlagSet("service rules list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	if len(flags.Args()) != 1 {
		log.Error("Please specify service id")
		return -1
	}
	log.Info("Service id is:", flags.Arg(0))

	client := c.Meta.Client()
	if rulesList, err := client.ListServiceRulesPaginated(context.Background(), flags.Arg(0)); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, rule := range rulesList {
			fmt.Println("Entry: ", i+1)
			data, err := yaml.Marshal(rule)
			if err != nil {
				log.Error(err)
				return -1
			}
			fmt.Println(string(data))
		}
	}
	return 0
}
