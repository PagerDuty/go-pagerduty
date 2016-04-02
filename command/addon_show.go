package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"strings"
)

type AddonShow struct {
	Meta
}

func (c *AddonShow) Help() string {
	helpText := `
	pd addon show <ID> Show details of an addon
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *AddonShow) Synopsis() string {
	return "Show details of an addon"
}

func AddonShowCommand() (cli.Command, error) {
	return &AddonShow{}, nil
}

func (c *AddonShow) Run(args []string) int {
	flags := c.Meta.FlagSet("addon show")
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
	data, err := yaml.Marshal(a)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(string(data))
	return 0
}
