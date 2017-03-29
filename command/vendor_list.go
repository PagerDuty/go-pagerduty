package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type VendorList struct {
	Meta
}

func VendorListCommand() (cli.Command, error) {
	return &VendorList{}, nil
}

func (c *VendorList) Help() string {
	helpText := `
	pd vendor list List vendors

	Options:

		 -query   Filter vendors with certain name

	`
	return strings.TrimSpace(helpText)
}

func (c *VendorList) Synopsis() string {
	return "List vendors within PagerDuty, optionally filtered by a search query"
}

func (c *VendorList) Run(args []string) int {
	var query string

	flags := c.Meta.FlagSet("user list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&query, "query", "", "Show vendors whose names contain the query")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	opts := pagerduty.ListVendorOptions{
		Query: query,
	}
	if resp, err := client.ListVendors(opts); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, p := range resp.Vendors {
			fmt.Println("Entry: ", i)
			data, err := yaml.Marshal(p)
			if err != nil {
				log.Error(err)
				return -1
			}
			fmt.Println(string(data))
		}
	}

	return 0
}
