package rbac

// Global user-scoped roles and their lookup collections.
var (
	// GlobalOwnerRole is the highest-authority global user role.
	GlobalOwnerRole = &Role{
		Name:  "global:owner",
		Score: 100,
	}
	// GlobalAdminRole is the administrative global user role.
	GlobalAdminRole = &Role{
		Name:  "global:admin",
		Score: 75,
	}
	// GlobalViewerRole is the read-only global user role.
	GlobalViewerRole = &Role{
		Name:  "global:viewer",
		Score: 25,
	}

	// UserRoles lists all global user roles ordered by descending authority.
	UserRoles = []*Role{
		GlobalOwnerRole,
		GlobalAdminRole,
		GlobalViewerRole,
	}

	// UserRolesMap indexes the global user roles by their name.
	UserRolesMap = map[string]*Role{
		GlobalOwnerRole.Name:  GlobalOwnerRole,
		GlobalAdminRole.Name:  GlobalAdminRole,
		GlobalViewerRole.Name: GlobalViewerRole,
	}
)

// IsValidUserRole reports whether role is a known global user role.
func IsValidUserRole(role string) bool {
	_, ok := UserRolesMap[role]
	return ok
}
