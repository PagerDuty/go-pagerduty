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

type ServiceIntegrationCreate struct {
	Meta
}

func ServiceIntegrationCreateCommand() (cli.Command, error) {
	return &ServiceIntegrationCreate{}, nil
}

func (c *ServiceIntegrationCreate) Help() string {
	helpText := `
	pd service integration create <SERVICE_ID> <FILE> Create a new integration from json file
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *ServiceIntegrationCreate) Synopsis() string {
	return "Create a new integration within service"
}

func (c *ServiceIntegrationCreate) Run(args []string) int {
	flags := c.Meta.FlagSet("service integration create")
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
	var i pagerduty.Integration
	if len(flags.Args()) != 2 {
		log.Error("Please specify service id and input json file")
		return -1
	}
	log.Info("Service id is:", flags.Arg(0))
	log.Info("Input file is:", flags.Arg(1))
	f, err := os.Open(flags.Arg(1))
	if err != nil {
		log.Error(err)
		return -1
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&i); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}
	log.Debugf("%#v", i)
	if _, err := client.CreateIntegration(flags.Arg(0), i); err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
