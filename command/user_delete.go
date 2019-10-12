package main

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type UserDelete struct {
	Meta
}

func UserDeleteCommand() (cli.Command, error) {
	return &UserDelete{}, nil
}

func (c *UserDelete) Help() string {
	helpText := `
	pd user delete <USER_ID> Remove an existing user
	`
	return strings.TrimSpace(helpText)
}

func (c *UserDelete) Synopsis() string {
	return "Remove an existing user"
}

func (c *UserDelete) Run(args []string) int {
	var userID string
	flags := c.Meta.FlagSet("user delete")
	flags.StringVar(&userID, "id", "", "User ID to remove")
	flags.Usage = func() { fmt.Println(c.Help()) }
	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}

	client := c.Meta.Client()
	if err := client.DeleteUser(userID); err != nil {
		log.Error(err)
		return -1
	}
	log.Infof("Delete user %d", userID)
	return 0
}
