// Package rbac defines the role-based access control model for the platform,
// including roles, permissions, and the mappings that associate each role with
// its granted permissions.
package rbac

// Permission represents a single access-control permission identified by its code.
type Permission struct {
	Code string
}

// NewPermission returns a new Permission with the given code.
func NewPermission(code string) *Permission {
	permission := &Permission{
		Code: code,
	}
	return permission
}
