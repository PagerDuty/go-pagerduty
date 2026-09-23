package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type IncidentShow struct {
	Meta
}

func IncidentShowCommand() (cli.Command, error) {
	return &IncidentShow{}, nil
}

func (c *IncidentShow) Help() string {
	helpText := `
	Usage: pd incident show -id INCIDENT_ID

	Show detailed information about an incident as YAML.
	Requests time out after 30 seconds.

	Options:
		-id Incident ID (required)
	`
	return strings.TrimSpace(helpText) + "\n\n" + c.Meta.Help()
}

func (c *IncidentShow) Synopsis() string {
	return "Show detailed information about an incident"
}

func (c *IncidentShow) Run(args []string) int {
	var id string
	flags := c.Meta.FlagSet("incident show")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&id, "id", "", "Incident ID")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		log.Error(err)
		return -1
	}
	if strings.TrimSpace(id) == "" {
		log.Error("You must specify the incident id using -id flag")
		return -1
	}
	if flags.NArg() != 0 {
		log.Error("Unexpected arguments; specify the incident id using -id flag")
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}

	client := c.Meta.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	incident, err := client.GetIncidentWithContext(ctx, url.PathEscape(id))
	if err != nil {
		log.Error(err)
		return -1
	}
	data, err := yaml.Marshal(incident)
	if err != nil {
		log.Error(err)
		return -1
	}
	if _, err := fmt.Println(string(data)); err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
