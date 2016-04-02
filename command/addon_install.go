package main

import (
	"encoding/json"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

type AddonInstall struct {
	Meta
}

func AddonInstallCommand() (cli.Command, error) {
	return &AddonInstall{}, nil
}

func (c *AddonInstall) Help() string {
	helpText := `
	pd addon install <FILE> Install a new addon
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *AddonInstall) Synopsis() string {
	return "Install a new addon"
}

func (c *AddonInstall) Run(args []string) int {
	flags := c.Meta.FlagSet("addon install")
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
	var a pagerduty.Addon
	if len(flags.Args()) != 1 {
		log.Error("Please specify input json file")
		return -1
	}
	log.Info("Input file is:", flags.Arg(0))
	f, err := os.Open(flags.Arg(0))
	if err != nil {
		log.Error(err)
		return -1
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&a); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}
	log.Debugf("%#v", a)
	if err := client.InstallAddon(a); err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
