package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type fileTokenSource struct {
	base           oauth2.TokenSource
	configFilePath string
	clientId       string
	scopes         []string
}

type persistedToken struct {
	*oauth2.Token
	ClientId string `json:"clientId"`
	Scopes   string `json:"scopes"`
}

// NewFileTokenSource creates an oauth2.TokenSource with a Token method which is
// able to load/save the token info in a file located at configFilePath (e.g.,
// "token.json" to use the file token.json at CWD).
func NewFileTokenSource(ctx context.Context, clientId string, clientSecret string, scopes []string, configFilePath string) oauth2.TokenSource {
	base := baseTokenSource(ctx, clientId, clientSecret, scopes)

	fts := &fileTokenSource{
		base:           base,
		configFilePath: configFilePath,
		clientId:       clientId,
		scopes:         scopes,
	}

	return oauth2.ReuseTokenSource(nil, fts)
}

func (c *fileTokenSource) loadToken() (*oauth2.Token, error) {
	_, err := os.Stat(c.configFilePath)
	if os.IsNotExist(err) {
		return &oauth2.Token{}, nil
	}

	data, err := os.ReadFile(c.configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", c.configFilePath, err)
	}
	pt := persistedToken{}
	if err := json.Unmarshal(data, &pt); err != nil {
		return nil, fmt.Errorf("failed to decode content of %s: %w", c.configFilePath, err)
	}

	t := pt.Token

	needToRefreshTokenExpired := time.Now().After(pt.Expiry)
	needToRefreshTokenNotSameScopes := !isSameScope(pt.Scopes, strings.Join(c.scopes, " "))
	needToRefreshTokenNotSameCredentials := pt.ClientId != c.clientId
	if needToRefreshTokenExpired || needToRefreshTokenNotSameScopes || needToRefreshTokenNotSameCredentials {
		t.AccessToken = ""
	}

	return t, nil
}

func (c *fileTokenSource) saveToken(tok *oauth2.Token) error {
	_, err := os.Stat(c.configFilePath)
	if os.IsNotExist(err) {
		if _, err := os.Create(c.configFilePath); err != nil {
			return fmt.Errorf("failed to create file %s: %w", c.configFilePath, err)
		}
	}

	p := persistedToken{
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

func (c *fileTokenSource) Token() (t *oauth2.Token, err error) {
	t, _ = c.loadToken()
	if t != nil && t.Valid() {
		return t, nil
	}

	if t, err = c.base.Token(); err != nil {
		return nil, err
	}

	err = c.saveToken(t)

	return t, err
}

func isSameScope(a, b string) bool {
	ta := strings.TrimSpace(a)
	tb := strings.TrimSpace(b)

	return strings.Compare(ta, tb) == 0
}
