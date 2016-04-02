package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"strings"
)

type AddonDelete struct {
	Meta
}

func (c *AddonDelete) Help() string {
	helpText := `
	pd addon delete <ID> Delete details of an addon
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *AddonDelete) Synopsis() string {
	return "Delete details of an addon"
}

func AddonDeleteCommand() (cli.Command, error) {
	return &AddonDelete{}, nil
}

func (c *AddonDelete) Run(args []string) int {
	flags := c.Meta.FlagSet("addon delete")
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
	err := client.DeleteAddon(flags.Arg(0))
	if err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
