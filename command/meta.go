package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ArrayFlags []string

func (a *ArrayFlags) String() string {
	return strings.Join(*a, ",")
}

func (a *ArrayFlags) Set(v string) error {
	if *a == nil {
		*a = make([]string, 0, 1)
	}
	*a = append(*a, v)
	return nil
}

type Meta struct {
	Authtoken    string `yaml:"authtoken"`
	Loglevel     string `yaml:"loglevel"`
	OutputFormat string `yaml:"outputformat"`
	Marshaler    func(any) ([]byte, error)
}

type FlagSetFlags uint

func (m *Meta) FlagSet(n string) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)
	f.StringVar(&m.Authtoken, "authtoken", "", "PagerDuty API authentication token (default: value of $PAGERDUTY_API_KEY)")
	f.StringVar(&m.Loglevel, "loglevel", "", "Logging level")
	f.StringVar(&m.OutputFormat, "outputformat", "", "Output format (valid values: json yaml)")
	return f
}

func (m *Meta) Client() *pagerduty.Client {
	return pagerduty.NewClient(m.Authtoken)
}

func (m *Meta) Help() string {
	helpText := `
	Common options:

	-authtoken PagerDuty API authentication token
	-loglevel Logging level
`
	return strings.TrimSpace(helpText)
}

func (m *Meta) validate() error {
	if m.Authtoken == "" {
		return fmt.Errorf("No authentication token provided")
	}
	return nil
}

func (m *Meta) Setup() error {
	m.setupLogging()
	if err := m.loadConfig(); err != nil {
		log.Warn(err)
	}
	m.setupMarshaler()
	return m.validate()
}

func (m *Meta) setupMarshaler() {
	switch {
	case m.OutputFormat == "json":
		m.Marshaler = json.Marshal
	case m.OutputFormat == "yaml":
		m.Marshaler = yaml.Marshal
	default:
		log.Fatalf("invalid output format %q; must be one of: json yaml", m.OutputFormat)
	}
}

func (m *Meta) setupLogging() {
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	switch m.Loglevel {
	case "info", "":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Fatal("Unknown log level", m.Loglevel)
	}
}

func (m *Meta) loadConfig() error {
	path, err := homedir.Dir()
	if err != nil {
		return err
	}
	configFile := filepath.Join(path, ".pd.yml")
	if _, err := os.Stat(configFile); err != nil {
		return err
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	other := &Meta{}
	if err := yaml.Unmarshal(data, other); err != nil {
		return err
	}
	if m.Authtoken == "" {
		m.Authtoken = other.Authtoken
	}
	if m.Loglevel == "" {
		m.Loglevel = other.Loglevel
	}
	if m.OutputFormat == "" {
		if other.OutputFormat != "" {
			m.OutputFormat = other.OutputFormat
		} else {
			m.OutputFormat = "yaml"
		}
	}
	return nil
}
