package main

import (
	"encoding/json"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"strings"
)

type AddonUpdate struct {
	Meta
}

func (c *AddonUpdate) Help() string {
	helpText := `
	pd addon update <ID> Update details of an addon
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *AddonUpdate) Synopsis() string {
	return "Update details of an addon"
}

func AddonUpdateCommand() (cli.Command, error) {
	return &AddonUpdate{}, nil
}

func (c *AddonUpdate) Run(args []string) int {
	flags := c.Meta.FlagSet("addon update")
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
	if len(flags.Args()) != 1 {
		log.Error("Please specify addon id")
		return -1
	}
	a, err := client.GetAddon(flags.Arg(0))
	if err != nil {
		log.Error(err)
		return -1
	}

	oldData, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Error(err)
		return -1
	}

	newData, err := TextEditor(oldData)
	if err != nil {
		log.Error(err)
		return -1
	}
	var newAddon pagerduty.Addon
	if err := json.Unmarshal(newData, &newAddon); err != nil {
		log.Error(err)
		return -1
	}
	if err := client.UpdateAddon(newAddon.ID, newAddon); err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
