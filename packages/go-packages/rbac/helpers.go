package rbac

// OrganizationPermissions returns the permission codes granted to the given
// organization role, or nil if the role is unknown.
func OrganizationPermissions(role string) []string {
	perms, ok := OrganizationRolePermissions[role]
	if !ok {
		return nil
	}

	result := make([]string, 0, len(perms))
	for p := range perms {
		result = append(result, p)
	}
	return result
}

// UserPermissions returns the permission codes granted to the given user role,
// or nil if the role is unknown.
func UserPermissions(role string) []string {
	perms, ok := UserRolePermissions[role]
	if !ok {
		return nil
	}

	result := make([]string, 0, len(perms))
	for p := range perms {
		result = append(result, p)
	}
	return result
}

// HasOrganizationPermission reports whether the given organization role grants
// the specified permission.
func HasOrganizationPermission(role string, permission string) bool {
	perms, ok := OrganizationRolePermissions[role]
	if !ok {
		return false
	}
	_, exists := perms[permission]
	return exists
}

// HasUserPermission reports whether the given user role grants the specified
// permission.
func HasUserPermission(role string, permission string) bool {
	perms, ok := UserRolePermissions[role]
	if !ok {
		return false
	}
	_, exists := perms[permission]
	return exists
}
