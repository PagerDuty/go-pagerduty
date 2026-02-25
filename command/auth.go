package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// OAuth endpoints for Classic and Scoped OAuth apps
const (
	// Classic OAuth endpoints (app.pagerduty.com)
	classicAuthURL  = "https://app.pagerduty.com/oauth/authorize"
	classicTokenURL = "https://app.pagerduty.com/oauth/token"

	// Scoped OAuth endpoints (identity.pagerduty.com)
	scopedAuthURL  = "https://identity.pagerduty.com/oauth/authorize"
	scopedTokenURL = "https://identity.pagerduty.com/oauth/token"
)

// getOAuthEndpoints returns the appropriate auth and token URLs for the given OAuth type
func getOAuthEndpoints(oauthType OAuthType) (authURL, tokenURL string) {
	switch oauthType {
	case OAuthTypeScoped:
		return scopedAuthURL, scopedTokenURL
	default:
		return classicAuthURL, classicTokenURL
	}
}

type AuthLoginCommand struct {
	Meta
}

func NewAuthLoginCommand() (cli.Command, error) {
	return &AuthLoginCommand{}, nil
}

func (c *AuthLoginCommand) Help() string {
	helpText := `
Usage: pd auth login [options]

	Authenticate with PagerDuty using OAuth. Opens a browser for you to log in
	and authorize the application.

Options:

	-client-id      OAuth client ID (required)
	-client-secret  OAuth client secret (required)
	-scope          OAuth scope (can be specified multiple times)
	-token-file     Path to store OAuth token (defaults to ~/.pd-token.json)
	-port           Local port for OAuth callback (default: 8080)
`
	return strings.TrimSpace(helpText)
}

func (c *AuthLoginCommand) Synopsis() string {
	return "Authenticate with PagerDuty using OAuth"
}

func (c *AuthLoginCommand) Run(args []string) int {
	var port int
	flags := c.Meta.FlagSet("auth login")
	flags.IntVar(&port, "port", 8080, "Local port for OAuth callback")
	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return 1
	}

	if err := c.Meta.loadConfig(); err != nil {
		log.Debug(err)
	}

	if c.Meta.ClientID == "" || c.Meta.ClientSecret == "" {
		log.Error("OAuth client-id and client-secret are required")
		return 1
	}

	if len(c.Meta.Scopes) == 0 {
		log.Error("At least one OAuth scope is required")
		return 1
	}

	tokenFile := c.Meta.TokenFile
	if tokenFile == "" {
		path, err := homedir.Dir()
		if err != nil {
			log.Warnf("Failed to get home directory: %v, using current directory", err)
			tokenFile = ".pd-token.json"
		} else {
			tokenFile = filepath.Join(path, ".pd-token.json")
		}
	}

	redirectURL := fmt.Sprintf("http://localhost:%d/callback", port)

	authURL, tokenURL := getOAuthEndpoints(c.Meta.OAuthType)
	config := &oauth2.Config{
		ClientID:     c.Meta.ClientID,
		ClientSecret: c.Meta.ClientSecret,
		Scopes:       c.Meta.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:   authURL,
			TokenURL:  tokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: redirectURL,
	}

	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Errorf("Failed to start local server: %v", err)
		return 1
	}

	server := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errMsg := r.URL.Query().Get("error")
			errDesc := r.URL.Query().Get("error_description")
			if errMsg != "" {
				errChan <- fmt.Errorf("OAuth error: %s - %s", errMsg, errDesc)
			} else {
				errChan <- fmt.Errorf("no authorization code received")
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body><h1>Authentication Failed</h1><p>You can close this window.</p></body></html>")
			return
		}

		codeChan <- code
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body><h1>Authentication Successful!</h1><p>You can close this window and return to the terminal.</p></body></html>")
	})

	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error: %v", err)
		}
	}()

	state := fmt.Sprintf("%d", time.Now().UnixNano())
	authorizeURL := config.AuthCodeURL(state)

	fmt.Printf("Opening browser for authentication...\n")
	fmt.Printf("If the browser doesn't open, visit this URL:\n%s\n\n", authorizeURL)

	if err := openBrowser(authorizeURL); err != nil {
		log.Warnf("Failed to open browser: %v", err)
	}

	fmt.Println("Waiting for authentication...")

	var code string
	select {
	case code = <-codeChan:
	case err := <-errChan:
		log.Error(err)
		server.Close()
		return 1
	case <-time.After(5 * time.Minute):
		log.Error("Authentication timed out")
		server.Close()
		return 1
	}

	server.Close()

	ctx := context.Background()
	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Errorf("Failed to exchange code for token: %v", err)
		return 1
	}

	if err := saveAuthToken(tokenFile, token, c.Meta.ClientID, c.Meta.Scopes); err != nil {
		log.Errorf("Failed to save token: %v", err)
		return 1
	}

	fmt.Printf("\nAuthentication successful! Token saved to %s\n", tokenFile)
	return 0
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}

// AuthToken extends oauth2.Token with metadata for validation
type AuthToken struct {
	*oauth2.Token
	ClientID string   `json:"client_id"`
	Scopes   []string `json:"scopes"`
}

func saveAuthToken(path string, token *oauth2.Token, clientID string, scopes []string) error {
	authToken := AuthToken{
		Token:    token,
		ClientID: clientID,
		Scopes:   scopes,
	}

	data, err := json.MarshalIndent(authToken, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

func loadAuthToken(path string) (*AuthToken, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var token AuthToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	return &token, nil
}

// authCodeTokenSource implements oauth2.TokenSource for Authorization Code flow
// with automatic token refresh and file persistence.
type authCodeTokenSource struct {
	config       *oauth2.Config
	tokenFile    string
	clientID     string
	scopes       []string
	currentToken *oauth2.Token
}

// NewAuthCodeTokenSource creates a token source that loads tokens from file
// and refreshes them automatically using the refresh token.
func NewAuthCodeTokenSource(oauthType OAuthType, clientID, clientSecret string, scopes []string, tokenFile string) (oauth2.TokenSource, error) {
	authURL, tokenURL := getOAuthEndpoints(oauthType)
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:   authURL,
			TokenURL:  tokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	ts := &authCodeTokenSource{
		config:    config,
		tokenFile: tokenFile,
		clientID:  clientID,
		scopes:    scopes,
	}

	authToken, err := loadAuthToken(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("no valid token found, run 'pd auth login' first: %w", err)
	}

	if authToken.ClientID != clientID {
		return nil, fmt.Errorf("token was created with different client ID, run 'pd auth login' again")
	}

	ts.currentToken = authToken.Token

	return oauth2.ReuseTokenSource(ts.currentToken, ts), nil
}

func (ts *authCodeTokenSource) Token() (*oauth2.Token, error) {
	if ts.currentToken.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token available, run 'pd auth login' again")
	}

	ctx := context.Background()
	newToken, err := ts.config.TokenSource(ctx, ts.currentToken).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w (run 'pd auth login' to re-authenticate)", err)
	}

	if err := saveAuthToken(ts.tokenFile, newToken, ts.clientID, ts.scopes); err != nil {
		log.Warnf("Failed to save refreshed token: %v", err)
	}

	ts.currentToken = newToken
	return newToken, nil
}
