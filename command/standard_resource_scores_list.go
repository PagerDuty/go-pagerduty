package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type StandardListResourceScores struct {
	Meta
}

func (c *StandardListResourceScores) Help() string {
	helpText := `
	pd standard resource-scores list List all of the existing standards scores for a specific resource

	Options:

    -id
    -resource-type    List scores applied to resource type (default: technical_services). Allowed values: technical_services

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *StandardListResourceScores) Synopsis() string {
	return "List all of the existing standards scores applied to a specific resource"
}

func StandardListResourceScoresCommand() (cli.Command, error) {
	return &StandardListResourceScores{}, nil
}

func (c *StandardListResourceScores) Run(args []string) int {
	var resID string
	var resourceType string

	flags := c.Meta.FlagSet("standard resource-scores list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&resID, "id", "", "Resource id")
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
	if res, err := client.ListResourceStandardScores(context.Background(), resID, resourceType); err != nil {
		log.Error(err)
		return -1
	} else {
		data, err := yaml.Marshal(res)
		if err != nil {
			log.Error(err)
			return -1
		}
		fmt.Println(string(data))
	}
	return 0
}
