package rbac

// OrganizationRolePermissions maps each organization role name to the set of
// permission codes it grants.
var (
	OrganizationRolePermissions = map[string]map[string]struct{}{
		OrganizationOwnerRole.Name: {
			CreateUserPermission.Code: {},
			ReadUserPermission.Code:   {},
			UpdateUserPermission.Code: {},
			DeleteUserPermission.Code: {},

			CreateUserNoncePermission.Code: {},
			ReadUserNoncePermission.Code:   {},
			UpdateUserNoncePermission.Code: {},
			DeleteUserNoncePermission.Code: {},

			ReadNodeEnginePermission.Code: {},

			CreateOrganizationPermission.Code: {},
			ReadOrganizationPermission.Code:   {},
			UpdateOrganizationPermission.Code: {},
			DeleteOrganizationPermission.Code: {},

			CreateOrganizationUserPermission.Code:     {},
			ReadOrganizationUserPermission.Code:       {},
			UpdateOrganizationUserInfoPermission.Code: {},
			UpdateOrganizationUserRolePermission.Code: {},
			DeleteOrganizationUserPermission.Code:     {},

			CreateOrganizationCredentialsPermission.Code: {},
			ReadOrganizationCredentialsPermission.Code:   {},
			UpdateOrganizationCredentialsPermission.Code: {},
			DeleteOrganizationCredentialsPermission.Code: {},

			TriggerTaskPermission.Code: {},
			CancelTaskPermission.Code:  {},
			RetryTaskPermission.Code:   {},

			CreateWorkflowPermission.Code: {},
			ReadWorkflowPermission.Code:   {},
			UpdateWorkflowPermission.Code: {},
			DeleteWorkflowPermission.Code: {},

			CreateWorkflowGenerationSessionPermission.Code: {},
			ReadWorkflowGenerationSessionPermission.Code:   {},
			UpdateWorkflowGenerationSessionPermission.Code: {},
			DeleteWorkflowGenerationSessionPermission.Code: {},

			CreateWorkflowGenerationMessagePermission.Code: {},
			ReadWorkflowGenerationMessagePermission.Code:   {},

			ReadTaskExecutionPermission.Code:   {},
			DeleteTaskExecutionPermission.Code: {},

			CreateTriggersPermission.Code: {},
			ReadTriggersPermission.Code:   {},
			UpdateTriggersPermission.Code: {},
			DeleteTriggersPermission.Code: {},

			CreateAPIKeyPermission.Code: {},
			ReadAPIKeyPermission.Code:   {},
			UpdateAPIKeyPermission.Code: {},
			DeleteAPIKeyPermission.Code: {},

			CreateOAuth2Permission.Code: {},

			CreateLinkedAccountPermission.Code: {},
			ReadLinkedAccountsPermission.Code:  {},
			UpdateLinkedAccountPermission.Code: {},
			DeleteLinkedAccountPermission.Code: {},

			ReadSocialPermission.Code:   {},
			UpdateSocialPermission.Code: {},

			ReadPlatformCredentialPermission.Code: {},

			ReadNotificationPermission.Code:   {},
			UpdateNotificationPermission.Code: {},
			DeleteNotificationPermission.Code: {},
		},
		OrganizationAdminRole.Name: {
			CreateUserPermission.Code: {},
			ReadUserPermission.Code:   {},
			UpdateUserPermission.Code: {},

			CreateUserNoncePermission.Code: {},
			ReadUserNoncePermission.Code:   {},
			UpdateUserNoncePermission.Code: {},

			ReadNodeEnginePermission.Code: {},

			CreateOrganizationPermission.Code: {},
			ReadOrganizationPermission.Code:   {},
			UpdateOrganizationPermission.Code: {},
			DeleteOrganizationPermission.Code: {},

			CreateOrganizationUserPermission.Code:     {},
			ReadOrganizationUserPermission.Code:       {},
			UpdateOrganizationUserInfoPermission.Code: {},
			UpdateOrganizationUserRolePermission.Code: {},
			DeleteOrganizationUserPermission.Code:     {},

			CreateOrganizationCredentialsPermission.Code: {},
			ReadOrganizationCredentialsPermission.Code:   {},
			UpdateOrganizationCredentialsPermission.Code: {},
			DeleteOrganizationCredentialsPermission.Code: {},

			TriggerTaskPermission.Code: {},
			CancelTaskPermission.Code:  {},
			RetryTaskPermission.Code:   {},

			CreateWorkflowPermission.Code: {},
			ReadWorkflowPermission.Code:   {},
			UpdateWorkflowPermission.Code: {},
			DeleteWorkflowPermission.Code: {},

			CreateWorkflowGenerationSessionPermission.Code: {},
			ReadWorkflowGenerationSessionPermission.Code:   {},
			UpdateWorkflowGenerationSessionPermission.Code: {},
			DeleteWorkflowGenerationSessionPermission.Code: {},

			CreateWorkflowGenerationMessagePermission.Code: {},
			ReadWorkflowGenerationMessagePermission.Code:   {},

			ReadTaskExecutionPermission.Code:   {},
			DeleteTaskExecutionPermission.Code: {},

			CreateTriggersPermission.Code: {},
			ReadTriggersPermission.Code:   {},
			UpdateTriggersPermission.Code: {},
			DeleteTriggersPermission.Code: {},

			CreateAPIKeyPermission.Code: {},
			ReadAPIKeyPermission.Code:   {},
			UpdateAPIKeyPermission.Code: {},
			DeleteAPIKeyPermission.Code: {},

			CreateOAuth2Permission.Code: {},

			CreateLinkedAccountPermission.Code: {},
			ReadLinkedAccountsPermission.Code:  {},
			UpdateLinkedAccountPermission.Code: {},
			DeleteLinkedAccountPermission.Code: {},

			ReadSocialPermission.Code:   {},
			UpdateSocialPermission.Code: {},

			ReadPlatformCredentialPermission.Code: {},

			ReadNotificationPermission.Code:   {},
			UpdateNotificationPermission.Code: {},
			DeleteNotificationPermission.Code: {},
		},
		OrganizationEditorRole.Name: {
			ReadUserPermission.Code:   {},
			UpdateUserPermission.Code: {},

			CreateUserNoncePermission.Code: {},
			ReadUserNoncePermission.Code:   {},
			UpdateUserNoncePermission.Code: {},

			ReadNodeEnginePermission.Code: {},

			ReadOrganizationPermission.Code:   {},
			UpdateOrganizationPermission.Code: {},

			ReadOrganizationUserPermission.Code:       {},
			UpdateOrganizationUserInfoPermission.Code: {},

			CreateOrganizationCredentialsPermission.Code: {},
			ReadOrganizationCredentialsPermission.Code:   {},
			UpdateOrganizationCredentialsPermission.Code: {},

			TriggerTaskPermission.Code: {},
			CancelTaskPermission.Code:  {},
			RetryTaskPermission.Code:   {},

			CreateWorkflowPermission.Code: {},
			ReadWorkflowPermission.Code:   {},
			UpdateWorkflowPermission.Code: {},

			CreateWorkflowGenerationSessionPermission.Code: {},
			ReadWorkflowGenerationSessionPermission.Code:   {},
			UpdateWorkflowGenerationSessionPermission.Code: {},

			CreateWorkflowGenerationMessagePermission.Code: {},
			ReadWorkflowGenerationMessagePermission.Code:   {},

			ReadTaskExecutionPermission.Code: {},

			CreateTriggersPermission.Code: {},
			ReadTriggersPermission.Code:   {},
			UpdateTriggersPermission.Code: {},
			DeleteTriggersPermission.Code: {},

			CreateAPIKeyPermission.Code: {},
			ReadAPIKeyPermission.Code:   {},
			UpdateAPIKeyPermission.Code: {},
			DeleteAPIKeyPermission.Code: {},

			CreateOAuth2Permission.Code: {},

			CreateLinkedAccountPermission.Code: {},
			ReadLinkedAccountsPermission.Code:  {},
			UpdateLinkedAccountPermission.Code: {},
			DeleteLinkedAccountPermission.Code: {},

			ReadSocialPermission.Code:   {},
			UpdateSocialPermission.Code: {},

			ReadPlatformCredentialPermission.Code: {},

			ReadNotificationPermission.Code:   {},
			UpdateNotificationPermission.Code: {},
			DeleteNotificationPermission.Code: {},
		},
		OrganizationViewerRole.Name: {
			ReadUserPermission.Code: {},

			ReadUserNoncePermission.Code: {},

			ReadNodeEnginePermission.Code: {},

			ReadOrganizationPermission.Code: {},

			ReadOrganizationUserPermission.Code:       {},
			UpdateOrganizationUserInfoPermission.Code: {},

			TriggerTaskPermission.Code: {},
			CancelTaskPermission.Code:  {},
			RetryTaskPermission.Code:   {},

			ReadWorkflowPermission.Code: {},

			ReadWorkflowGenerationSessionPermission.Code: {},

			ReadWorkflowGenerationMessagePermission.Code: {},

			ReadTaskExecutionPermission.Code: {},

			CreateTriggersPermission.Code: {},
			ReadTriggersPermission.Code:   {},
			UpdateTriggersPermission.Code: {},
			DeleteTriggersPermission.Code: {},

			CreateAPIKeyPermission.Code: {},
			ReadAPIKeyPermission.Code:   {},
			UpdateAPIKeyPermission.Code: {},
			DeleteAPIKeyPermission.Code: {},

			CreateOAuth2Permission.Code: {},

			CreateLinkedAccountPermission.Code: {},
			ReadLinkedAccountsPermission.Code:  {},
			UpdateLinkedAccountPermission.Code: {},
			DeleteLinkedAccountPermission.Code: {},

			ReadSocialPermission.Code: {},

			ReadPlatformCredentialPermission.Code: {},

			ReadNotificationPermission.Code:   {},
			UpdateNotificationPermission.Code: {},
			DeleteNotificationPermission.Code: {},
		},
	}
)
