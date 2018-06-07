package main

import (
	"fmt"
	"strings"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
)

// IncidentCreate holds the meta data
type IncidentCreate struct {
	Meta
}

// IncidentCreateCommand creates the struct
func IncidentCreateCommand() (cli.Command, error) {
	return &IncidentCreate{}, nil
}

// Help returns how to run the command
func (c *IncidentCreate) Help() string {
	helpText := `
	pd incident create create incident

	`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a summary of the command
func (c *IncidentCreate) Synopsis() string {
	return "Create incidents"
}

// Run runs the command
func (c *IncidentCreate) Run(args []string) int {
	var service string
	var title string
	var from string
	var body string
	flags := c.Meta.FlagSet("incident create")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&service, "service", "", "service ID")
	flags.StringVar(&title, "title", "", "incident title")
	flags.StringVar(&from, "from", "", "user creating the ticket")
	flags.StringVar(&body, "body", "", "detailed ticket description")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}

	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}

	client := c.Meta.Client()
	opts := pagerduty.CreateIncidentOptions{
		Service: pagerduty.APIReference{
			Type: "service",
			ID:   service,
		},
		Title: title,
	}
	if body != "" {
		opts.Body = pagerduty.Body{
			Type:    "incident_body",
			Details: body,
		}
	}
	if incident, err := client.CreateIncident(from, opts); err != nil {
		log.Error(err)
		return -1
	} else {
		data, err := yaml.Marshal(incident)
		if err != nil {
			log.Error(err)
			return -1
		}
		fmt.Println(string(data))

	}
	return 0
}
