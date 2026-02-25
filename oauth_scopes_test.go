package pagerduty

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClassicScope_String(t *testing.T) {
	tests := []struct {
		name  string
		scope ClassicScope
		want  string
	}{
		{
			name:  "read",
			scope: ClassicScopeRead,
			want:  "read",
		},
		{
			name:  "write",
			scope: ClassicScopeWrite,
			want:  "write",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.String(); got != tt.want {
				t.Errorf("ClassicScope.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClassicScopesToStrings(t *testing.T) {
	tests := []struct {
		name   string
		scopes []ClassicScope
		want   []string
	}{
		{
			name:   "empty",
			scopes: []ClassicScope{},
			want:   []string{},
		},
		{
			name:   "read_only",
			scopes: []ClassicScope{ClassicScopeRead},
			want:   []string{"read"},
		},
		{
			name:   "both",
			scopes: []ClassicScope{ClassicScopeRead, ClassicScopeWrite},
			want:   []string{"read", "write"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassicScopesToStrings(tt.scopes)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ClassicScopesToStrings() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidateClassicScopes(t *testing.T) {
	tests := []struct {
		name    string
		scopes  []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid_read",
			scopes:  []string{"read"},
			wantErr: false,
		},
		{
			name:    "valid_write",
			scopes:  []string{"write"},
			wantErr: false,
		},
		{
			name:    "valid_both",
			scopes:  []string{"read", "write"},
			wantErr: false,
		},
		{
			name:    "invalid_scoped_oauth_scope",
			scopes:  []string{"incidents.read"},
			wantErr: true,
			errMsg:  "invalid scope 'incidents.read' for Classic OAuth",
		},
		{
			name:    "invalid_mixed",
			scopes:  []string{"read", "incidents.read"},
			wantErr: true,
			errMsg:  "invalid scope 'incidents.read' for Classic OAuth",
		},
		{
			name:    "empty",
			scopes:  []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateClassicScopes(tt.scopes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClassicScopes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if got := err.Error(); got != tt.errMsg {
					// Check if errMsg is a prefix (the full error includes the hint)
					if len(got) < len(tt.errMsg) || got[:len(tt.errMsg)] != tt.errMsg {
						t.Errorf("ValidateClassicScopes() error = %q, should start with %q", got, tt.errMsg)
					}
				}
			}
		})
	}
}

func TestValidateScopedOAuthScopes(t *testing.T) {
	tests := []struct {
		name    string
		scopes  []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid_granular_scope",
			scopes:  []string{"incidents.read"},
			wantErr: false,
		},
		{
			name:    "valid_multiple_scopes",
			scopes:  []string{"incidents.read", "services.write", "users.read"},
			wantErr: false,
		},
		{
			name:    "valid_with_account_scope",
			scopes:  []string{"as_account-us.mycompany", "incidents.read"},
			wantErr: false,
		},
		{
			name:    "valid_subresource_scope",
			scopes:  []string{"users:contact_methods.read"},
			wantErr: false,
		},
		{
			name:    "valid_openid",
			scopes:  []string{"openid"},
			wantErr: false,
		},
		{
			name:    "invalid_classic_read",
			scopes:  []string{"read"},
			wantErr: true,
			errMsg:  "invalid scope 'read' for Scoped OAuth",
		},
		{
			name:    "invalid_classic_write",
			scopes:  []string{"write"},
			wantErr: true,
			errMsg:  "invalid scope 'write' for Scoped OAuth",
		},
		{
			name:    "invalid_mixed_with_classic",
			scopes:  []string{"incidents.read", "read"},
			wantErr: true,
			errMsg:  "invalid scope 'read' for Scoped OAuth",
		},
		{
			name:    "invalid_format_no_dot",
			scopes:  []string{"incidents"},
			wantErr: true,
			errMsg:  "invalid scope 'incidents' for Scoped OAuth",
		},
		{
			name:    "empty",
			scopes:  []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateScopedOAuthScopes(tt.scopes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateScopedOAuthScopes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if got := err.Error(); got != tt.errMsg {
					// Check if errMsg is a prefix (the full error includes the hint)
					if len(got) < len(tt.errMsg) || got[:len(tt.errMsg)] != tt.errMsg {
						t.Errorf("ValidateScopedOAuthScopes() error = %q, should start with %q", got, tt.errMsg)
					}
				}
			}
		})
	}
}

func TestInvalidScopeError_Error(t *testing.T) {
	err := &InvalidScopeError{
		Scope:     "incidents.read",
		OAuthType: "Classic",
		Hint:      "Classic OAuth only supports 'read' and 'write' scopes",
	}

	want := "invalid scope 'incidents.read' for Classic OAuth: Classic OAuth only supports 'read' and 'write' scopes"
	if got := err.Error(); got != want {
		t.Errorf("InvalidScopeError.Error() = %q, want %q", got, want)
	}
}
