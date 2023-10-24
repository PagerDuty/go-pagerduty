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

type StandardList struct {
	Meta
}

func (c *StandardList) Help() string {
	helpText := `
	pd standard list List all of the existing standards

	Options:

    -active
    -resource-type    Filter by resource type. Allowed values: technical_service

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *StandardList) Synopsis() string {
	return "List all of the existing standards"
}

func StandardListCommand() (cli.Command, error) {
	return &StandardList{}, nil
}

func (c *StandardList) Run(args []string) int {
	var active bool
	var resourceType string

	flags := c.Meta.FlagSet("standard list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.BoolVar(&active, "active", false, "")
	flags.StringVar(&resourceType, "resource-type", "", "Filter by resource type. Allowed values: technical_service")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	opts := pagerduty.ListStandardsOptions{
		Active:       active,
		ResourceType: resourceType,
	}
	if res, err := client.ListStandards(context.Background(), opts); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, r := range res.Standards {
			if i > 0 {
				fmt.Println("---")
			}
			data, err := yaml.Marshal(r)
			if err != nil {
				log.Error(err)
				return -1
			}
			fmt.Println(string(data))
		}
	}
	return 0
}
