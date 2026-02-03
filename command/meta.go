package main

import (
	"flag"
	"fmt"
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

// OAuthType represents the type of OAuth app being used
type OAuthType string

const (
	// OAuthTypeClassic uses broad read/write scopes via app.pagerduty.com
	OAuthTypeClassic OAuthType = "classic"

	// OAuthTypeScoped uses granular scopes via identity.pagerduty.com
	OAuthTypeScoped OAuthType = "scoped"
)

type Meta struct {
	Authtoken    string
	Loglevel     string
	OAuthType    OAuthType  `yaml:"oauth_type"`
	ClientID     string     `yaml:"client_id"`
	ClientSecret string     `yaml:"client_secret"`
	Scopes       ArrayFlags `yaml:"scopes"`
	TokenFile    string     `yaml:"token_file"`
}

type FlagSetFlags uint

func (m *Meta) FlagSet(n string) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)
	f.StringVar(&m.Authtoken, "authtoken", "", "PagerDuty API authentication token")
	f.StringVar(&m.Loglevel, "loglevel", "", "Logging level")
	f.StringVar((*string)(&m.OAuthType), "oauth-type", "", "OAuth app type: 'classic' or 'scoped'")
	f.StringVar(&m.ClientID, "client-id", "", "OAuth client ID")
	f.StringVar(&m.ClientSecret, "client-secret", "", "OAuth client secret")
	f.Var(&m.Scopes, "scope", "OAuth scope (can be specified multiple times)")
	f.StringVar(&m.TokenFile, "token-file", "", "Path to store OAuth token (defaults to ~/.pd-token.json)")
	return f
}

func (m *Meta) Client() *pagerduty.Client {
	if m.useOAuth() {
		tokenFile := m.getTokenFilePath()

		tokenSource, err := NewAuthCodeTokenSource(m.OAuthType, m.ClientID, m.ClientSecret, m.Scopes, tokenFile)
		if err != nil {
			log.Fatalf("OAuth error: %v", err)
		}

		// Use appropriate client option based on OAuth type
		switch m.OAuthType {
		case OAuthTypeScoped:
			return pagerduty.NewClient("", pagerduty.WithScopedOAuthAppTokenSource(tokenSource))
		default:
			// Classic OAuth - token is already obtained, use it directly
			return pagerduty.NewClient("", pagerduty.WithScopedOAuthAppTokenSource(tokenSource))
		}
	}
	return pagerduty.NewClient(m.Authtoken)
}

func (m *Meta) getTokenFilePath() string {
	if m.TokenFile != "" {
		return m.TokenFile
	}
	path, err := homedir.Dir()
	if err != nil {
		log.Warnf("Failed to get home directory: %v, using current directory for token file", err)
		return ".pd-token.json"
	}
	return filepath.Join(path, ".pd-token.json")
}

func (m *Meta) useOAuth() bool {
	return m.OAuthType != "" && m.ClientID != "" && m.ClientSecret != "" && len(m.Scopes) > 0
}

func (m *Meta) Help() string {
	helpText := `
	Common options:

	-authtoken      PagerDuty API authentication token
	-loglevel       Logging level

	OAuth options (alternative to -authtoken):

	-oauth-type     OAuth app type: 'classic' or 'scoped'
	                  classic: uses broad read/write scopes
	                  scoped:  uses granular resource-based scopes
	-client-id      OAuth client ID
	-client-secret  OAuth client secret
	-scope          OAuth scope (can be specified multiple times)
	                  Classic scopes: read, write
	                  Scoped scopes:  incidents.read, services.write, etc.
	-token-file     Path to store OAuth token (defaults to ~/.pd-token.json)
`
	return strings.TrimSpace(helpText)
}

func (m *Meta) validate() error {
	hasAuthtoken := m.Authtoken != ""
	hasOAuth := m.useOAuth()
	hasPartialOAuth := (m.ClientID != "" || m.ClientSecret != "" || len(m.Scopes) > 0 || m.OAuthType != "") && !hasOAuth

	if hasPartialOAuth {
		var missing []string
		if m.OAuthType == "" {
			missing = append(missing, "-oauth-type")
		}
		if m.ClientID == "" {
			missing = append(missing, "-client-id")
		}
		if m.ClientSecret == "" {
			missing = append(missing, "-client-secret")
		}
		if len(m.Scopes) == 0 {
			missing = append(missing, "-scope")
		}
		return fmt.Errorf("incomplete OAuth configuration, missing: %s", strings.Join(missing, ", "))
	}

	if !hasAuthtoken && !hasOAuth {
		return fmt.Errorf("authentication required: provide either -authtoken or OAuth credentials (-oauth-type, -client-id, -client-secret, -scope)")
	}

	// Validate OAuth type
	if hasOAuth {
		if m.OAuthType != OAuthTypeClassic && m.OAuthType != OAuthTypeScoped {
			return fmt.Errorf("invalid oauth-type %q: must be 'classic' or 'scoped'", m.OAuthType)
		}

		// Validate scopes match OAuth type
		if err := m.validateScopes(); err != nil {
			return err
		}
	}

	return nil
}

func (m *Meta) validateScopes() error {
	scopes := []string(m.Scopes)

	switch m.OAuthType {
	case OAuthTypeClassic:
		if err := pagerduty.ValidateClassicScopes(scopes); err != nil {
			return fmt.Errorf("scope validation failed: %w", err)
		}
	case OAuthTypeScoped:
		if err := pagerduty.ValidateScopedOAuthScopes(scopes); err != nil {
			return fmt.Errorf("scope validation failed: %w", err)
		}
	}

	return nil
}

func (m *Meta) Setup() error {
	m.setupLogging()
	if err := m.loadConfig(); err != nil {
		log.Warn(err)
	}
	return m.validate()
}

func (m *Meta) setupLogging() {
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
	data, err := os.ReadFile(configFile)
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
	if m.OAuthType == "" {
		m.OAuthType = other.OAuthType
	}
	if m.ClientID == "" {
		m.ClientID = other.ClientID
	}
	if m.ClientSecret == "" {
		m.ClientSecret = other.ClientSecret
	}
	if len(m.Scopes) == 0 {
		m.Scopes = other.Scopes
	}
	if m.TokenFile == "" {
		m.TokenFile = other.TokenFile
	}
	return nil
}
