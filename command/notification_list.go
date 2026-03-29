package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type NotificationList struct {
	Meta
}

func NotificationListCommand() (cli.Command, error) {
	return &NotificationList{}, nil
}

func (c *NotificationList) Help() string {
	helpText := `
	notification list List notifications

	Options:
		 -since     Beginning timestamp
		 -until     Ending timestamp
		 -filter    Return notification of specified type (sms_notification, email_notification, phone_notification, push_notification)
		 -time-zone Time zone in which results will be rendered (default is account time zone)
		 -limit     Maximum number of results to return
	`
	return strings.TrimSpace(helpText)
}

func (c *NotificationList) Synopsis() string {
	return "List notifications for a given time range"
}

func (c *NotificationList) Run(args []string) int {
	var (
		since, until string
		filter       string
		timeZone     string
		limit        int
	)
	flags := c.Meta.FlagSet("notification list")
	flags.Usage = func() { fmt.Println(c.Help()) }

	flags.StringVar(&since, "since", "", "Beginning timestamp")
	flags.StringVar(&until, "until", "", "Ending timestamp")
	flags.StringVar(&filter, "filter", "", "Return notification of specified type (sms_notification, email_notification, phone_notification, push_notification)")
	flags.StringVar(&timeZone, "time-zone", "", "Time zone in which results will be rendered (default is account time zone)")
	flags.IntVar(&limit, "limit", 0, "Maximum number of results to return")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}

	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()

	opts := pagerduty.ListNotificationOptions{
		Limit:    100,
		TimeZone: timeZone,
		Since:    since,
		Until:    until,
		Filter:   filter,
	}

	notifs, err := depaginateNotifications(context.Background(), client, opts, limit)
	if err != nil {
		log.Error(err)
		return -1
	}

	for i, notif := range notifs {
		fmt.Println("Entry: ", i+1)
		data, err := yaml.Marshal(notif)
		if err != nil {
			log.Error(err)
			return -1
		}
		fmt.Println(string(data))
	}

	return 0
}

// depaginateNotifications collects all available pages of notifications
func depaginateNotifications(ctx context.Context, cli *pagerduty.Client, opts pagerduty.ListNotificationOptions, limit int) ([]pagerduty.Notification, error) {
	var res []pagerduty.Notification

	for {
		if limit > 0 {
			remaining := limit - len(res)
			opts.Limit = uint(clamp(remaining, 1, 100))
		}

		notifs, err := cli.ListNotificationsWithContext(ctx, opts)
		if err != nil {
			return nil, err
		}

		res = append(res, notifs.Notifications...)

		if !notifs.More || len(res) >= limit {
			break
		}

		opts.Offset = notifs.Offset
	}

	return res, nil
}

func clamp(n, min, max int) int {
	if n > max {
		return max
	}
	if n < min {
		return min
	}
	return n
}
