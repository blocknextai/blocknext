package rbac

// UserRolePermissions maps each global user role name to the set of permission
// codes it grants.
var (
	UserRolePermissions = map[string]map[string]struct{}{
		GlobalOwnerRole.Name: {
			ReadUserPermission.Code:   {},
			UpdateUserPermission.Code: {},

			ReadUserCredentialsPermission.Code:   {},
			CreateUserCredentialsPermission.Code: {},
			UpdateUserCredentialsPermission.Code: {},
			DeleteUserCredentialsPermission.Code: {},

			ReadNodeEnginePermission.Code: {},

			CreateOrganizationPermission.Code: {},
			ReadOrganizationPermission.Code:   {},

			CreateOAuth2Permission.Code: {},

			TriggerTaskPermission.Code: {},
			CancelTaskPermission.Code:  {},
			RetryTaskPermission.Code:   {},

			CreateLinkedAccountPermission.Code: {},
			ReadLinkedAccountsPermission.Code:  {},
			UpdateLinkedAccountPermission.Code: {},
			DeleteLinkedAccountPermission.Code: {},

			ReadSocialPermission.Code:   {},
			UpdateSocialPermission.Code: {},

			ReadPlatformCredentialPermission.Code: {},

			ReadSessionPermission.Code:   {},
			DeleteSessionPermission.Code: {},

			CreateAPIKeyPermission.Code: {},
			ReadAPIKeyPermission.Code:   {},
			UpdateAPIKeyPermission.Code: {},
			DeleteAPIKeyPermission.Code: {},

			ReadUserPreferencesPermission.Code:   {},
			UpdateUserPreferencesPermission.Code: {},

			ReadNotificationPermission.Code:   {},
			UpdateNotificationPermission.Code: {},
			DeleteNotificationPermission.Code: {},
		},
		GlobalAdminRole.Name: {
			ReadUserPermission.Code:   {},
			UpdateUserPermission.Code: {},

			ReadUserCredentialsPermission.Code:   {},
			CreateUserCredentialsPermission.Code: {},
			UpdateUserCredentialsPermission.Code: {},
			DeleteUserCredentialsPermission.Code: {},

			ReadNodeEnginePermission.Code: {},

			CreateOrganizationPermission.Code: {},
			ReadOrganizationPermission.Code:   {},

			CreateOAuth2Permission.Code: {},

			TriggerTaskPermission.Code: {},
			CancelTaskPermission.Code:  {},
			RetryTaskPermission.Code:   {},

			CreateLinkedAccountPermission.Code: {},
			ReadLinkedAccountsPermission.Code:  {},
			UpdateLinkedAccountPermission.Code: {},
			DeleteLinkedAccountPermission.Code: {},

			ReadSocialPermission.Code:   {},
			UpdateSocialPermission.Code: {},

			ReadPlatformCredentialPermission.Code: {},

			ReadSessionPermission.Code:   {},
			DeleteSessionPermission.Code: {},

			CreateAPIKeyPermission.Code: {},
			ReadAPIKeyPermission.Code:   {},
			UpdateAPIKeyPermission.Code: {},
			DeleteAPIKeyPermission.Code: {},

			ReadUserPreferencesPermission.Code:   {},
			UpdateUserPreferencesPermission.Code: {},

			ReadNotificationPermission.Code:   {},
			UpdateNotificationPermission.Code: {},
			DeleteNotificationPermission.Code: {},
		},
		GlobalViewerRole.Name: {
			ReadUserPermission.Code: {},

			ReadUserCredentialsPermission.Code: {},

			ReadNodeEnginePermission.Code: {},

			ReadOrganizationPermission.Code: {},

			CreateLinkedAccountPermission.Code: {},
			ReadLinkedAccountsPermission.Code:  {},
			UpdateLinkedAccountPermission.Code: {},
			DeleteLinkedAccountPermission.Code: {},

			ReadSocialPermission.Code:   {},
			UpdateSocialPermission.Code: {},

			ReadPlatformCredentialPermission.Code: {},

			ReadSessionPermission.Code:   {},
			DeleteSessionPermission.Code: {},

			ReadAPIKeyPermission.Code: {},

			ReadUserPreferencesPermission.Code: {},

			ReadNotificationPermission.Code:   {},
			UpdateNotificationPermission.Code: {},
			DeleteNotificationPermission.Code: {},
		},
	}
)
