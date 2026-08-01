package rbac

// Organization-scoped roles and their lookup collections.
var (
	// OrganizationOwnerRole is the highest-authority organization role.
	OrganizationOwnerRole = &Role{
		Name:  "organization:owner",
		Score: 100,
	}
	// OrganizationAdminRole is the administrative organization role.
	OrganizationAdminRole = &Role{
		Name:  "organization:admin",
		Score: 75,
	}
	// OrganizationEditorRole is the organization role with edit access.
	OrganizationEditorRole = &Role{
		Name:  "organization:editor",
		Score: 50,
	}
	// OrganizationViewerRole is the read-only organization role.
	OrganizationViewerRole = &Role{
		Name:  "organization:viewer",
		Score: 25,
	}

	// OrganizationRoles lists all organization roles ordered by descending authority.
	OrganizationRoles = []*Role{
		OrganizationOwnerRole,
		OrganizationAdminRole,
		OrganizationEditorRole,
		OrganizationViewerRole,
	}

	// OrganizationRolesMap indexes the organization roles by their name.
	OrganizationRolesMap = map[string]*Role{
		OrganizationOwnerRole.Name:  OrganizationOwnerRole,
		OrganizationAdminRole.Name:  OrganizationAdminRole,
		OrganizationEditorRole.Name: OrganizationEditorRole,
		OrganizationViewerRole.Name: OrganizationViewerRole,
	}
)

// IsValidOrganizationRole reports whether role is a known organization role.
func IsValidOrganizationRole(role string) bool {
	_, ok := OrganizationRolesMap[role]
	return ok
}
