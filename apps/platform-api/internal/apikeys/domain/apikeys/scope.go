package apikeys

import (
	"slices"
)

type Scope string

const (
	ScopeWorkflowsTrigger Scope = "workflows:trigger"
	ScopeMCPInvoke        Scope = "mcp:invoke"
)

var (
	AllScopes = map[Scope]struct{}{
		ScopeWorkflowsTrigger: {},
		ScopeMCPInvoke:        {},
	}
)

func (s Scope) IsValid() bool {
	_, ok := AllScopes[s]
	return ok
}

type Scopes []Scope

func (s Scopes) Has(scope Scope) bool {
	return slices.Contains(s, scope)
}
