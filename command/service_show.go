package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type ServiceShow struct {
	Meta
}

func ServiceShowCommand() (cli.Command, error) {
	return &ServiceShow{}, nil
}

func (c *ServiceShow) Help() string {
	helpText := `
	service show Show specific service

	Options:

		 -id
		 -include

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *ServiceShow) Synopsis() string {
	return "Get details about an existing service"
}

func (c *ServiceShow) Run(args []string) int {
	var includes []string
	flags := c.Meta.FlagSet("service show")
	flags.Usage = func() { fmt.Println(c.Help()) }
	servID := flags.String("id", "", "Service ID")
	flags.Var((*ArrayFlags)(&includes), "include", "Additional details to include (can be specified multiple times)")
	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	if servID == nil {
		log.Error("You must provide a service id")
		return -1
	}
	client := c.Meta.Client()
	o := &pagerduty.GetServiceOptions{
		Includes: includes,
	}
	servicerecord, err := client.GetService(*servID, o)
	if err != nil {
		log.Error(err)
		return -1
	}
	data, err := yaml.Marshal(servicerecord)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(string(data))
	return 0
}
