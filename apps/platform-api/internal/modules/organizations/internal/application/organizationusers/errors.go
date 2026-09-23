package organizationusers

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToUpdateOrganizationUserRole = apperror.Internal("failed to update organization user role")
	ErrFailedToUpdateOrganizationUserInfo = apperror.Internal("failed to update organization user info")
	ErrFailedToDeleteOrganizationUser     = apperror.Internal("failed to delete organization user")

	ErrCouldBeOneOwner                    = apperror.Conflict("there must be at least one owner")
	ErrUserCannotChangeOwnRole            = apperror.Forbidden("user cannot change own role")
	ErrCantDeleteTheOwnerOrganizationUser = apperror.Forbidden("cannot delete the owner")

	ErrIdentifierIsRequired = apperror.Validation("identifier is required")
)
