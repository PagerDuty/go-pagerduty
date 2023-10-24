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

type StandardListMultiResourcesScores struct {
	Meta
}

func (c *StandardListMultiResourcesScores) Help() string {
	helpText := `
	pd standard multi-resources-scores list List all of the existing standards scores for multiple resources

	Options:

    -id    List standard scores applied to resource (can be specified multiple times)
    -resource-type    List scores applied to resource type (default: technical_services). Allowed values: technical_services

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *StandardListMultiResourcesScores) Synopsis() string {
	return "List all of the existing standards scores applied to multiple resources"
}

func StandardListMultiResourcesScoresCommand() (cli.Command, error) {
	return &StandardListMultiResourcesScores{}, nil
}

func (c *StandardListMultiResourcesScores) Run(args []string) int {
	var resIDs []string
	var resourceType string

	flags := c.Meta.FlagSet("standard multi-resources-scores list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.Var((*ArrayFlags)(&resIDs), "id", "Display scores for standard applied to resource (can be specified multiple times)")
	flags.StringVar(&resourceType, "resource-type", "technical_services", "Scores applied to resource type (default: technical_services). Allowed values: technical_services")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	opts := pagerduty.ListMultiResourcesStandardScoresOptions{
		IDs: resIDs,
	}
	if res, err := client.ListMultiResourcesStandardScores(context.Background(), resourceType, opts); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, r := range res.Resources {
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
