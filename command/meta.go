package main

import (
	"flag"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"strings"
)

type Meta struct {
	authtoken string
	subdomain string
	loglevel  string
}

type FlagSetFlags uint

func (m *Meta) FlagSet(n string) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)
	f.StringVar(&m.authtoken, "authtoken", "", "PagerDuty API authentication token")
	f.StringVar(&m.subdomain, "subdomain", "", "PagerDuty account name (subdomain)")
	f.StringVar(&m.loglevel, "loglevel", "", "Logging level")
	return f
}

func (m *Meta) Client() *pagerduty.Client {
	return pagerduty.NewClient(m.subdomain, m.authtoken)
}

func (m *Meta) Help() string {
	helpText := `
	Generral options:

	-authtoken PagerDuty API authentication token
	-subdomain PagerDuty account ID
	-loglevel Logging level
`
	return strings.TrimSpace(helpText)
}

func (m *Meta) Validate() error {
	if m.authtoken == "" {
		return fmt.Errorf("Authtoken can not be blank")
	}
	if m.subdomain == "" {
		return fmt.Errorf("Subdomain can not be blank")
	}
	return nil
}

func (m *Meta) SetupLogging() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	switch m.loglevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Fatal("Unknown log level", m.loglevel)
	}
}
