package credential

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var ErrInvalidReference = errors.New("invalid credential reference")

type Scope string

const (
	CredentialScope  = "credential"
	CredentialPrefix = CredentialScope + ":"

	OrganizationCredentialScope Scope = "organization"

	CredentialMaskValue = "__BLOCKNEXT_CREDENTIAL_VALUE_1121a730-5165-4c9b-af0e-29759f08b65c__"
)

var (
	CredentialScopes = map[Scope]struct{}{
		OrganizationCredentialScope: {},
	}
)

func (s Scope) String() string {
	return string(s)
}

func (s Scope) IsValid() bool {
	_, ok := CredentialScopes[s]
	return ok
}

func BuildUIKey(ownerType string, id string) string {
	switch ownerType {
	case OrganizationCredentialScope.String():
		return CredentialPrefix + OrganizationCredentialScope.String() + ":" + id
	default:
		return id
	}
}

func ParseReference(ref string) (Scope, uuid.UUID, error) {
	parts := strings.Split(ref, ":")
	if len(parts) != 3 || parts[0] != CredentialScope {
		return "", uuid.Nil, ErrInvalidReference
	}

	scope := Scope(parts[1])
	if !scope.IsValid() {
		return "", uuid.Nil, ErrInvalidReference
	}

	credentialID, err := uuid.Parse(parts[2])
	if err != nil {
		return "", uuid.Nil, ErrInvalidReference
	}

	return scope, credentialID, nil
}
