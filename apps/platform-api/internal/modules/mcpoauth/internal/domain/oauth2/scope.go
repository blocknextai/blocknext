package oauth2

import (
	"slices"
	"strings"
)

type Scope string

const (
	MCPInvokeScope     Scope = "mcp:invoke"
	PlatformReadScope  Scope = "platform:read"
	PlatformWriteScope Scope = "platform:write"
	OfflineAccessScope Scope = "offline_access"
)

var (
	AllScopes = map[Scope]struct{}{
		MCPInvokeScope:     {},
		PlatformReadScope:  {},
		PlatformWriteScope: {},
		OfflineAccessScope: {},
	}

	DefaultScopes = Scopes{MCPInvokeScope, PlatformReadScope}

	ScopeDescriptions = map[Scope]string{
		MCPInvokeScope:     "Run platform tools on behalf of your organization",
		PlatformReadScope:  "Read your profile, organizations, workflows and executions",
		PlatformWriteScope: "Manage your workflows and platform settings",
		OfflineAccessScope: "Stay connected without asking you to sign in again",
	}
)

func (s Scope) String() string {
	return string(s)
}

func (s Scope) IsValid() bool {
	_, ok := AllScopes[s]
	return ok
}

func (s Scope) Description() string {
	return ScopeDescriptions[s]
}

type Scopes []Scope

func ParseScopes(raw string) Scopes {
	fields := strings.Fields(raw)
	scopes := make(Scopes, 0, len(fields))
	for _, field := range fields {
		scope := Scope(field)
		if !scopes.Has(scope) {
			scopes = append(scopes, scope)
		}
	}
	return scopes
}

func ScopesFromStrings(values []string) Scopes {
	scopes := make(Scopes, 0, len(values))
	for _, value := range values {
		scopes = append(scopes, Scope(value))
	}
	return scopes
}

func (s Scopes) Has(scope Scope) bool {
	return slices.Contains(s, scope)
}

func (s Scopes) HasAll(other Scopes) bool {
	for _, scope := range other {
		if !s.Has(scope) {
			return false
		}
	}
	return true
}

func (s Scopes) IsValid() bool {
	for _, scope := range s {
		if !scope.IsValid() {
			return false
		}
	}
	return true
}

func (s Scopes) Strings() []string {
	values := make([]string, 0, len(s))
	for _, scope := range s {
		values = append(values, scope.String())
	}
	return values
}

func (s Scopes) String() string {
	return strings.Join(s.Strings(), " ")
}
