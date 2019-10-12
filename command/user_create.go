package main

import (
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type UserCreate struct {
	Meta
}

func UserCreateCommand() (cli.Command, error) {
	return &UserCreate{}, nil
}

func (c *UserCreate) Help() string {
	helpText := `
	pd user create <EMAIL> Create a user
	`
	return strings.TrimSpace(helpText)
}

func (c *UserCreate) Synopsis() string {
	return "Create a new user"
}

func (c *UserCreate) Run(args []string) int {
	var email string
	flags := c.Meta.FlagSet("user create")
	flags.StringVar(&email, "email", "", "Email of the user")
	flags.Usage = func() { fmt.Println(c.Help()) }
	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	u := pagerduty.User{
		Email: email,
	}
	if len(flags.Args()) != 1 {
	}
	if _, err := client.CreateUser(u); err != nil {
		log.Error(err)
		return -1
	}

	return 0
}
