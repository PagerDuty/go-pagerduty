package main

import (
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
)

type OncallList struct {
	Meta
}

func OncallListCommand() (cli.Command, error) {
	return &OncallList{}, nil
}

func (c *OncallList) Help() string {
	helpText := `
	pd oncall list List on-calls

	`
	return strings.TrimSpace(helpText)
}

func (c *OncallList) Synopsis() string {
	return "List the on-call entries during a given time range"
}

func (c *OncallList) Run(args []string) int {
	var escalationPolicyIDs []string
	var includes []string
	var scheduleIDs []string
	var timeZone string
	var userIDs []string
	/* not yet implemented, defaults to now
	var earliest string
	var until string
	var since string
	*/
	flags := c.Meta.FlagSet("on-call list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.Var((*ArrayFlags)(&includes), "include", "Additional details to include (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&scheduleIDs), "schedule-id", "Only show for schedule ID (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&userIDs), "user-id", "Only show for user ID (can be specified multiple times)")
	flags.Var((*ArrayFlags)(&escalationPolicyIDs), "escalationPolicy-id", "Only show for escalationPolicy ID (can be specified multiple times)")
	flags.StringVar(&timeZone, "time-zone", "", "Time Zone")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	opts := pagerduty.ListOnCallOptions{
		UserIDs:             userIDs,
		Includes:            includes,
		TimeZone:            timeZone,
		EscalationPolicyIDs: escalationPolicyIDs,
		ScheduleIDs:         scheduleIDs,
	}
	if oncs, err := client.ListOnCalls(opts); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, p := range oncs.OnCalls {
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
