package main

import "testing"

func TestGetOAuthEndpoints(t *testing.T) {
	tests := []struct {
		name      string
		oauthType OAuthType
		wantAuth  string
		wantToken string
	}{
		{
			name:      "classic",
			oauthType: OAuthTypeClassic,
			wantAuth:  classicAuthURL,
			wantToken: classicTokenURL,
		},
		{
			name:      "scoped",
			oauthType: OAuthTypeScoped,
			wantAuth:  scopedAuthURL,
			wantToken: scopedTokenURL,
		},
		{
			name:      "default",
			oauthType: "",
			wantAuth:  classicAuthURL,
			wantToken: classicTokenURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAuth, gotToken := getOAuthEndpoints(tt.oauthType)
			if gotAuth != tt.wantAuth {
				t.Errorf("auth URL = %q, want %q", gotAuth, tt.wantAuth)
			}
			if gotToken != tt.wantToken {
				t.Errorf("token URL = %q, want %q", gotToken, tt.wantToken)
			}
		})
	}
}
