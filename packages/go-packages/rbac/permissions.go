package rbac

// Permissions enumerates every access-control permission recognized by the platform.
var (
	CreateUserPermission = NewPermission("users:create")
	ReadUserPermission   = NewPermission("users:read")
	UpdateUserPermission = NewPermission("users:update")
	DeleteUserPermission = NewPermission("users:delete")

	CreateUserNoncePermission = NewPermission("user_nonces:create")
	ReadUserNoncePermission   = NewPermission("user_nonces:read")
	UpdateUserNoncePermission = NewPermission("user_nonces:update")
	DeleteUserNoncePermission = NewPermission("user_nonces:delete")

	ReadUserCredentialsPermission   = NewPermission("user_credentials:read")
	CreateUserCredentialsPermission = NewPermission("user_credentials:create")
	UpdateUserCredentialsPermission = NewPermission("user_credentials:update")
	DeleteUserCredentialsPermission = NewPermission("user_credentials:delete")

	ReadNodeEnginePermission = NewPermission("node_engine:read")

	CreateOrganizationPermission = NewPermission("organizations:create")
	ReadOrganizationPermission   = NewPermission("organizations:read")
	UpdateOrganizationPermission = NewPermission("organizations:update")
	DeleteOrganizationPermission = NewPermission("organizations:delete")

	CreateOrganizationUserPermission     = NewPermission("organization_users:create")
	ReadOrganizationUserPermission       = NewPermission("organization_users:read")
	UpdateOrganizationUserPermission     = NewPermission("organization_users:update")
	DeleteOrganizationUserPermission     = NewPermission("organization_users:delete")
	UpdateOrganizationUserInfoPermission = NewPermission("organization_users:update:info")
	UpdateOrganizationUserRolePermission = NewPermission("organization_users:update:role")

	CreateOrganizationCredentialsPermission = NewPermission("organization_credentials:create")
	ReadOrganizationCredentialsPermission   = NewPermission("organization_credentials:read")
	UpdateOrganizationCredentialsPermission = NewPermission("organization_credentials:update")
	DeleteOrganizationCredentialsPermission = NewPermission("organization_credentials:delete")

	CreateOAuth2Permission = NewPermission("oauth2:create")

	TriggerTaskPermission = NewPermission("task_runner:trigger")
	CancelTaskPermission  = NewPermission("task_runner:cancel")
	RetryTaskPermission   = NewPermission("task_runner:retry")

	CreateWorkflowPermission = NewPermission("workflows:create")
	ReadWorkflowPermission   = NewPermission("workflows:read")
	UpdateWorkflowPermission = NewPermission("workflows:update")
	DeleteWorkflowPermission = NewPermission("workflows:delete")

	CreateWorkflowGenerationSessionPermission = NewPermission("workflows:generation:sessions:create")
	ReadWorkflowGenerationSessionPermission   = NewPermission("workflows:generation:sessions:read")
	UpdateWorkflowGenerationSessionPermission = NewPermission("workflows:generation:sessions:update")
	DeleteWorkflowGenerationSessionPermission = NewPermission("workflows:generation:sessions:delete")

	CreateWorkflowGenerationMessagePermission = NewPermission("workflows:generation:messages:create")
	ReadWorkflowGenerationMessagePermission   = NewPermission("workflows:generation:messages:read")

	ReadTaskExecutionPermission   = NewPermission("executions:read")
	DeleteTaskExecutionPermission = NewPermission("executions:delete")

	CreateTriggersPermission = NewPermission("triggers:create")
	ReadTriggersPermission   = NewPermission("triggers:read")
	UpdateTriggersPermission = NewPermission("triggers:update")
	DeleteTriggersPermission = NewPermission("triggers:delete")

	CreateAPIKeyPermission = NewPermission("api_keys:create")
	ReadAPIKeyPermission   = NewPermission("api_keys:read")
	UpdateAPIKeyPermission = NewPermission("api_keys:update")
	DeleteAPIKeyPermission = NewPermission("api_keys:delete")

	CreateLinkedAccountPermission = NewPermission("linked_accounts:create")
	ReadLinkedAccountsPermission  = NewPermission("linked_accounts:read")
	UpdateLinkedAccountPermission = NewPermission("linked_accounts:update")
	DeleteLinkedAccountPermission = NewPermission("linked_accounts:delete")

	ReadSocialPermission   = NewPermission("socials:read")
	UpdateSocialPermission = NewPermission("socials:update")

	ReadPlatformCredentialPermission = NewPermission("platform_credentials:read")

	ReadSessionPermission   = NewPermission("sessions:read")
	DeleteSessionPermission = NewPermission("sessions:delete")

	ReadUserPreferencesPermission   = NewPermission("user_preferences:read")
	UpdateUserPreferencesPermission = NewPermission("user_preferences:update")

	ReadNotificationPermission   = NewPermission("notifications:read")
	UpdateNotificationPermission = NewPermission("notifications:update")
	DeleteNotificationPermission = NewPermission("notifications:delete")
)
