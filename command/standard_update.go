package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"

	"github.com/mitchellh/cli"
)

type StandardUpdate struct {
	Meta
}

func StandardUpdateCommand() (cli.Command, error) {
	return &StandardUpdate{}, nil
}

func (c *StandardUpdate) Help() string {
	helpText := `
	pd standard update <FILE> Update an standard from json file
	Options:

		 -id

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *StandardUpdate) Synopsis() string {
	return "Update an existing standard"
}

func (c *StandardUpdate) Run(args []string) int {
	var standardID string

	flags := c.Meta.FlagSet("standard update")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&standardID, "id", "", "standard id")
	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	if standardID == "" {
		log.Error("You must provide standard id")
		return -1
	}

	client := c.Meta.Client()
	var standard pagerduty.Standard
	if len(flags.Args()) != 1 {
		log.Error("Please specify input json file")
		return -1
	}

	log.Info("Input file is:", flags.Arg(0))
	f, err := os.ReadFile(flags.Arg(0))
	if err != nil {
		log.Error(err)
		return -1
	}

	if err := json.Unmarshal(f, &standard); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}

	log.Debugf("%#v", standard)
	if _, err := client.UpdateStandard(context.Background(), standardID, standard); err != nil {
		log.Error(err)
		return -1
	}

	return 0
}
