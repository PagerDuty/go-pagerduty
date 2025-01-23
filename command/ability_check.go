package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

type AbilityTest struct {
	Meta
}

func (c *AbilityTest) Help() string {
	helpText := `
	pd ability test <ABILITY>

	Test if an account has given ability
	 ` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *AbilityTest) Synopsis() string {
	return "Test if an account has given ability"
}

func AbilityTestCommand() (cli.Command, error) {
	return &AbilityTest{}, nil
}

func (c *AbilityTest) Run(args []string) int {
	flags := c.Meta.FlagSet("ability test")
	flags.Usage = func() { fmt.Println(c.Help()) }

	if err := flags.Parse(args); err != nil {
		fmt.Println(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		fmt.Println(err)
		return -1
	}
	if len(flags.Args()) != 1 {
		fmt.Println("Please specify an ability")
		return -1
	}
	client := c.Meta.Client()
	if err := client.TestAbilityWithContext(context.Background(), flags.Arg(0)); err != nil {
		fmt.Println(err)
		return -1
	}
	return 0
}
