package pagerduty

import "strings"

// ClassicScope represents scopes for Classic OAuth apps.
// Classic OAuth only supports two broad permission levels: read and write.
// These scopes grant access to all resources the authenticated user can access.
type ClassicScope string

const (
	// ClassicScopeRead grants read access to all resources the user can access.
	ClassicScopeRead ClassicScope = "read"

	// ClassicScopeWrite grants write access to all resources the user can modify.
	ClassicScopeWrite ClassicScope = "write"
)

// String returns the string representation of the ClassicScope.
func (s ClassicScope) String() string {
	return string(s)
}

// ClassicScopesToStrings converts a slice of ClassicScope to []string for oauth2 compatibility.
func ClassicScopesToStrings(scopes []ClassicScope) []string {
	result := make([]string, len(scopes))
	for i, s := range scopes {
		result[i] = s.String()
	}
	return result
}

// ValidateClassicScopes checks if the provided string scopes are valid Classic OAuth scopes.
// Returns an error if any scope is not "read" or "write".
func ValidateClassicScopes(scopes []string) error {
	for _, s := range scopes {
		switch s {
		case "read", "write":
		default:
			return &InvalidScopeError{
				Scope:     s,
				OAuthType: "Classic",
				Hint:      "Classic OAuth only supports 'read' and 'write' scopes",
			}
		}
	}
	return nil
}

// ValidateScopedOAuthScopes checks if the provided scopes are valid for Scoped OAuth.
// Returns an error if a classic scope ("read" or "write") is detected.
func ValidateScopedOAuthScopes(scopes []string) error {
	for _, s := range scopes {
		if s == "read" || s == "write" {
			return &InvalidScopeError{
				Scope:     s,
				OAuthType: "Scoped",
				Hint:      "Use granular scopes like 'incidents.read' or 'services.write' for Scoped OAuth apps",
			}
		}

		if strings.HasPrefix(s, "as_account-") {
			continue
		}

		if !strings.Contains(s, ".") && s != "openid" {
			return &InvalidScopeError{
				Scope:     s,
				OAuthType: "Scoped",
				Hint:      "Scoped OAuth scopes should follow the pattern 'resource.permission' (e.g., 'incidents.read')",
			}
		}
	}
	return nil
}

// InvalidScopeError is returned when an invalid scope is used for an OAuth type.
type InvalidScopeError struct {
	Scope     string
	OAuthType string
	Hint      string
}

func (e *InvalidScopeError) Error() string {
	return "invalid scope '" + e.Scope + "' for " + e.OAuthType + " OAuth: " + e.Hint
}
