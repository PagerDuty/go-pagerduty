package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type AddonList struct {
	Meta
}

func (c *AddonList) Help() string {
	helpText := `
	pd addon list List all of the existing escalation policies

	Options:

		 -include    Include additional details
		 -service-id Filter result by service ids
		 -filter     Filter result by type
	 ` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *AddonList) Synopsis() string {
	return "List all addons"
}

func AddonListCommand() (cli.Command, error) {
	return &AddonList{}, nil
}

func (c *AddonList) Run(args []string) int {
	var includes []string
	var serviceIDs []string
	var filter string
	flags := c.Meta.FlagSet("addon list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.Var((*ArrayFlags)(&includes), "include", "Additional details to include (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&serviceIDs), "service-id", "Show addons only for specified services (can be specified multiple times)")
	flags.StringVar(&filter, "filter", "", "Filter results, showing only add-ons of the given type")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	opts := pagerduty.ListAddonOptions{
		Includes:   includes,
		ServiceIDs: serviceIDs,
		Filter:     filter,
	}
	if addonList, err := client.ListAddons(opts); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, addon := range addonList.Addons {
			fmt.Println("Entry: ", i)
			data, err := yaml.Marshal(addon)
			if err != nil {
				log.Error(err)
				return -1
			}
			fmt.Println(string(data))
		}
	}
	return 0
}
