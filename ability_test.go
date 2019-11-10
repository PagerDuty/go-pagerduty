package pagerduty

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAbility_ListAbilities(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/abilities", func(w http.ResponseWriter, r *http.Request) {
		require.Equal("GET", r.Method)
		w.Write([]byte(`{"abilities": ["sso"]}`))
	})

	client := &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}

	res, err := client.ListAbilities()
	require.NoError(err)
	require.Equal(&ListAbilityResponse{Abilities: []string{"sso"}}, res)
}
