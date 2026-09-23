package server

import (
	"slices"
	"strings"

	"github.com/google/uuid"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/blocknextai/go-packages/apperror"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
)

var (
	ErrMissingOwner              = apperror.Unauthorized("missing authenticated owner")
	ErrOrganizationOwnerRequired = apperror.Forbidden("platform tools require an organization scoped token")
	ErrScopeRequired             = apperror.Forbidden("the access token does not have the scope required to perform this operation")
	ErrInvalidWorkflowID         = apperror.Validation("invalid workflow id")
	ErrInvalidTriggerID          = apperror.Validation("invalid trigger id")
)

type caller struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Scopes         []string
}

func callerFromRequest(req *mcpsdk.CallToolRequest) (*caller, error) {
	if req.Extra == nil {
		return nil, ErrMissingOwner
	}

	if commonDomain.OwnerType(req.Extra.Header.Get(commonAuth.OwnerTypeHeader)) != commonDomain.OwnerTypeOrganization {
		return nil, ErrOrganizationOwnerRequired
	}

	organizationID, err := uuid.Parse(req.Extra.Header.Get(commonAuth.OwnerIDHeader))
	if err != nil {
		return nil, ErrMissingOwner.WithCause(err)
	}

	userID, err := uuid.Parse(req.Extra.Header.Get(commonAuth.UserIDHeader))
	if err != nil {
		return nil, ErrMissingOwner.WithCause(err)
	}

	return &caller{
		UserID:         userID,
		OrganizationID: organizationID,
		Scopes:         strings.Fields(req.Extra.Header.Get(commonAuth.ScopesHeader)),
	}, nil
}

func requireCaller(req *mcpsdk.CallToolRequest, requiredScope string) (*caller, error) {
	authenticated, err := callerFromRequest(req)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(authenticated.Scopes, requiredScope) {
		return nil, ErrScopeRequired
	}

	return authenticated, nil
}

func normalizePagination(offset int, limit int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = defaultPageSize
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}
	return offset, limit
}
