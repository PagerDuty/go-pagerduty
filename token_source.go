package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type fileTokenSource struct {
	base           oauth2.TokenSource
	configFilePath string
}

func NewFileTokenSource(context context.Context, clientId string, clientSecret string, scopes []string, configFilePath string) oauth2.TokenSource {
	base := baseTokenSource(context, clientId, clientSecret, scopes)

	fts := &fileTokenSource{
		base:           base,
		configFilePath: configFilePath,
	}

	return oauth2.ReuseTokenSource(nil, fts)
}

func (c *fileTokenSource) loadToken() (*oauth2.Token, error) {
	log.Printf("[FTS] Loading token from file\n")

	_, err := os.Stat(c.configFilePath)
	if os.IsNotExist(err) {
		return &oauth2.Token{}, nil
	}

	data, err := os.ReadFile(c.configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", c.configFilePath, err)
	}
	t := &oauth2.Token{}
	if err := json.Unmarshal(data, t); err != nil {
		return nil, fmt.Errorf("failed to decode content of %s: %w", c.configFilePath, err)
	}

	return t, nil
}

func (c *fileTokenSource) saveToken(tok *oauth2.Token) error {
	// Note that if we continue to rely on oauth2.TokenSource, the `expires_in` field needs
	// to be taken into account.
	log.Printf("[FTS] Saving token to file\n")

	_, err := os.Stat(c.configFilePath)
	if os.IsNotExist(err) {
		if _, err := os.Create(c.configFilePath); err != nil {
			return fmt.Errorf("failed to create file %s: %w", c.configFilePath, err)
		}
	}

	data, err := json.Marshal(tok)
	if err != nil {
		return fmt.Errorf("failed to encode token into file %s: %w", c.configFilePath, err)
	}

	if err := os.WriteFile(c.configFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save token into file %s: %w", c.configFilePath, err)
	}

	return nil
}

func (c *fileTokenSource) Token() (t *oauth2.Token, err error) {
	t, _ = c.loadToken()
	if t != nil && t.Valid() {
		return t, nil
	}

	log.Printf("[FTS] Fetching new token\n")
	if t, err = c.base.Token(); err != nil {
		return nil, err
	}

	err = c.saveToken(t)

	return t, err
}

/*
Token Source Implementation for Terraform Provider
*/
type tfprovTokenSource struct {
	base           oauth2.TokenSource
	configFilePath string
	clientId       string
	scopes         []string
}

type tfprovPersistentToken struct {
	*oauth2.Token
	ClientId string `json:"clientId"`
	Scopes   string `json:"scopes"`
}

func NewTFProvFileTokenSource(context context.Context, clientId string, clientSecret string, scopes []string, configFilePath string) oauth2.TokenSource {
	base := baseTokenSource(context, clientId, clientSecret, scopes)

	fts := &tfprovTokenSource{
		base:           base,
		configFilePath: configFilePath,
		clientId:       clientId,
		scopes:         scopes,
	}

	return oauth2.ReuseTokenSource(nil, fts)
}

func (c *tfprovTokenSource) loadToken() (*oauth2.Token, error) {
	log.Printf("[FTS] Loading token from file\n")

	_, err := os.Stat(c.configFilePath)
	if os.IsNotExist(err) {
		return &oauth2.Token{}, nil
	}

	data, err := os.ReadFile(c.configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", c.configFilePath, err)
	}
	pt := tfprovPersistentToken{}
	if err := json.Unmarshal(data, &pt); err != nil {
		return nil, fmt.Errorf("failed to decode content of %s: %w", c.configFilePath, err)
	}

	t := pt.Token

	// [*] TODO: Persisted token is already expired in config file.
	needToRefreshTokenExpired := time.Now().After(pt.Expiry)
	// [*] TODO: Addresss case when scopes are missing and token is not refreshed
	// after scopes are corrected until config file is deleted.
	needToRefreshTokenNotSameScopes := !isSameScope(pt.Scopes, strings.Join(c.scopes, " "))
	// [*] TODO: Address case when credentials are swaped and previous generated
	// token is not refreshed untill config file is deleted.
	needToRefreshTokenNotSameCredentials := pt.ClientId != c.clientId
	if needToRefreshTokenExpired || needToRefreshTokenNotSameScopes || needToRefreshTokenNotSameCredentials {
		t.AccessToken = ""
	}

	return t, nil
}

func (c *tfprovTokenSource) saveToken(tok *oauth2.Token) error {
	log.Printf("[FTS] Saving token to file\n")

	_, err := os.Stat(c.configFilePath)
	if os.IsNotExist(err) {
		if _, err := os.Create(c.configFilePath); err != nil {
			return fmt.Errorf("failed to create file %s: %w", c.configFilePath, err)
		}
	}

	p := tfprovPersistentToken{
		Token:    tok,
		ClientId: c.clientId,
		Scopes:   strings.Join(c.scopes, " "),
	}
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to encode token into file %s: %w", c.configFilePath, err)
	}

	if err := os.WriteFile(c.configFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save token into file %s: %w", c.configFilePath, err)
	}

	return nil
}

func (c *tfprovTokenSource) Token() (t *oauth2.Token, err error) {
	// [*] TODO: Address case when token expires before its end of life, because
	// it's revoked.
	//    * This should be addressed in client when receiving the 401 response.
	t, _ = c.loadToken()
	if t != nil && t.Valid() {
		return t, nil
	}

	log.Printf("[FTS] Fetching new token\n")
	if t, err = c.base.Token(); err != nil {
		return nil, err
	}

	err = c.saveToken(t)

	return t, err
}

func isSameScope(a, b string) bool {
	ta := strings.TrimSpace(a)
	tb := strings.TrimSpace(b)

	if strings.Compare(ta, tb) == 0 {
		return true
	}

	return false
}
