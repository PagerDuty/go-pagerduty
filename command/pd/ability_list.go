package main

import (
	"fmt"
	"github.com/mitchellh/cli"
	"log"
	"strings"
)

type AbilityList struct {
	Meta
}

func (c *AbilityList) Help() string {
	helpText := `
	pd ability list

	List all abilities of the account
	 ` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *AbilityList) Synopsis() string {
	return "List all abilities of the account"
}

func AbilityListCommand() (cli.Command, error) {
	return &AbilityList{}, nil
}

func (c *AbilityList) Run(args []string) int {
	flags := c.Meta.FlagSet("ability list")
	flags.Usage = func() { fmt.Println(c.Help()) }

	if err := flags.Parse(args); err != nil {
		log.Println(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Println(err)
		return -1
	}
	client := c.Meta.Client()
	abilityList, err := client.ListAbilities()
	if err != nil {
		log.Println(err)
		return -1
	}
	for _, ability := range abilityList.Abilities {
		fmt.Println(ability)
	}
	return 0
}
