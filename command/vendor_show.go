package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type VendorShow struct {
	Meta
}

func VendorShowCommand() (cli.Command, error) {
	return &VendorShow{}, nil
}

func (c *VendorShow) Help() string {
	helpText := `
	vendor show Show details of an individual vendor

	Options:
		 -id   Vendor ID
	`
	return strings.TrimSpace(helpText)
}

func (c *VendorShow) Synopsis() string { return "Get details about a vendor" }

func (c *VendorShow) Run(args []string) int {
	log.Println(args)
	var id string
	flags := c.Meta.FlagSet("vendor show")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&id, "id", "", "Vendor ID")
	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}

	if id == "" {
		log.Error("You must specify the vendor id using -id flag")
		return -1
	}
	client := c.Meta.Client()
	result, err := client.GetVendor(id)
	if err != nil {
		log.Error(err)
		return -1
	}
	data, err := yaml.Marshal(result)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(string(data))
	return 0
}
