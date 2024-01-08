package pagerduty

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

type fileTokenSource struct {
	base     oauth2.TokenSource
	filename string
	//TODO: delete me
	tmpFileToken *oauth2.Token
}

func NewFileTokenSource(context context.Context, clientId string, clientSecret string, scopes []string, filename string) oauth2.TokenSource {
	base := baseTokenSource(context, clientId, clientSecret, scopes)

	fts := &fileTokenSource{
		base:     base,
		filename: filename,
	}

	return oauth2.ReuseTokenSource(nil, fts)
}

func (c *fileTokenSource) saveToken(tok *oauth2.Token) error {
	// Fake code for saving token
	// TODO: write real code to save the token to the config file.
	// Note that if we continue to rely on oauth2.TokenSource, the `expires_in` field needs
	// to be taken into account.
	fmt.Printf("[FTS] Saving token to file\n")
	c.tmpFileToken = tok
	return nil
}
func (c *fileTokenSource) loadToken() (*oauth2.Token, error) {
	// Fake code for loading the token
	// TODO: write real code to load the token from a config file.
	fmt.Printf("[FTS] Loading token from file\n")
	return c.tmpFileToken, nil
}

func (c *fileTokenSource) Token() (t *oauth2.Token, err error) {
	// Load token and try to use it
	t, _ = c.loadToken()
	if t != nil && t.Valid() {
		return t, nil
	}

	// Fetch a new token
	fmt.Printf("[FTS] Fetching new token\n")
	if t, err = c.base.Token(); err != nil {
		return nil, err
	}

	// Attempt to save it
	err = c.saveToken(t)

	return t, err
}
