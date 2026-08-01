package presentation

import (
	"time"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/rbac"
	accountInfrastructure "github.com/blocknextai/platform-api/internal/account/infrastructure"
	"github.com/blocknextai/platform-api/internal/account/presentation/auth"
	authEmail "github.com/blocknextai/platform-api/internal/account/presentation/auth/email"
	authMagicLink "github.com/blocknextai/platform-api/internal/account/presentation/auth/email/magiclink"
	authPassword "github.com/blocknextai/platform-api/internal/account/presentation/auth/password"
	"github.com/blocknextai/platform-api/internal/account/presentation/linkedaccounts"
	"github.com/blocknextai/platform-api/internal/account/presentation/sessions"
	"github.com/blocknextai/platform-api/internal/account/presentation/userpreferences"
	"github.com/blocknextai/platform-api/internal/account/presentation/users"
	usersEmail "github.com/blocknextai/platform-api/internal/account/presentation/users/email"
	usersPassword "github.com/blocknextai/platform-api/internal/account/presentation/users/password"
	"github.com/blocknextai/platform-api/internal/account/presentation/usersocials"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/gofiber/fiber/v3"
)

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	cacheMiddleware *cachemiddleware.Middleware,
	handlers *accountInfrastructure.Handlers,
) {
	registerAuthRoutes(router, cacheMiddleware, authMiddleware, handlers)
	registerUserRoutes(router, authMiddleware, handlers)
}

func registerAuthRoutes(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *accountInfrastructure.Handlers,
) {
	authRouterGroup := router.Group("/auth")

	authRouterGroup.Get(
		"/methods",
		cacheMiddleware.Cache(5*time.Minute),
		auth.NewGetAuthMethodsHandler(handlers.GetAuthMethods),
	)

	authRouterGroup.Post(
		"/token/refresh",
		auth.NewRefreshTokenHandler(handlers.RefreshToken),
	)

	authRouterGroup.Post(
		"/logout",
		authMiddleware.Authenticate(),
		sessions.NewLogoutHandler(handlers.Logout),
	)

	oauthRouterGroup := authRouterGroup.Group("/oauth")

	oauthRouterGroup.Post(
		"/nonce",
		auth.NewCreateUserNonceHandler(handlers.CreateUserNonce),
	)

	oauthRouterGroup.Post(
		"/token",
		auth.NewCreateUserTokenHandler(handlers.CreateUserToken),
	)

	passwordRouterGroup := authRouterGroup.Group("/password")

	passwordRouterGroup.Post(
		"/register",
		authPassword.NewRegisterHandler(handlers.PasswordRegister),
	)

	passwordRouterGroup.Post(
		"/login",
		authPassword.NewLoginHandler(handlers.PasswordLogin),
	)

	passwordRouterGroup.Post(
		"/forgot",
		authPassword.NewForgotHandler(handlers.PasswordForgot),
	)

	passwordRouterGroup.Post(
		"/reset",
		authPassword.NewResetHandler(handlers.PasswordReset),
	)

	emailRouterGroup := authRouterGroup.Group("/email")

	emailRouterGroup.Post(
		"/verify",
		authEmail.NewVerifyHandler(handlers.PasswordVerify),
	)

	emailRouterGroup.Post(
		"/resend-verification",
		authEmail.NewResendVerificationHandler(handlers.ResendVerification),
	)

	emailRouterGroup.Post(
		"/change/confirm",
		authEmail.NewConfirmEmailChangeHandler(handlers.ConfirmEmailChange),
	)

	emailRouterGroup.Post(
		"/magic-link/request",
		authMagicLink.NewMagicLinkRequestHandler(handlers.MagicLinkRequest),
	)

	emailRouterGroup.Post(
		"/magic-link/consume",
		authMagicLink.NewMagicLinkConsumeHandler(handlers.MagicLinkConsume),
	)
}

func registerUserRoutes(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *accountInfrastructure.Handlers,
) {
	meRouterGroup := router.Group("/users/me")

	meRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadUserPermission),
		users.NewGetProfileHandler(handlers.GetProfile),
	)

	meRouterGroup.Get(
		"/roles",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadUserPermission),
		users.NewGetRolesHandler(handlers.GetRoles),
	)

	emailMeRouterGroup := meRouterGroup.Group("/email")

	emailMeRouterGroup.Post(
		"/add",
		authMiddleware.Authenticate(),
		usersEmail.NewAddEmailHandler(handlers.AddEmail),
	)

	emailMeRouterGroup.Post(
		"/change",
		authMiddleware.Authenticate(),
		usersEmail.NewChangeEmailHandler(handlers.ChangeEmail),
	)

	passwordRouterGroup := meRouterGroup.Group("/password")

	passwordRouterGroup.Post(
		"/set",
		authMiddleware.Authenticate(),
		usersPassword.NewSetPasswordHandler(handlers.SetPassword),
	)

	passwordRouterGroup.Post(
		"/change",
		authMiddleware.Authenticate(),
		usersPassword.NewChangePasswordHandler(handlers.ChangePassword),
	)

	linkedAccountsRouterGroup := meRouterGroup.Group("/linked-accounts")

	linkedAccountsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadLinkedAccountsPermission),
		linkedaccounts.NewGetAllLinkedAccountsHandler(handlers.GetAllLinkedAccounts),
	)

	linkedAccountsRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.CreateLinkedAccountPermission),
		linkedaccounts.NewAddLinkedAccountHandler(handlers.AddLinkedAccount),
	)

	linkedAccountsRouterGroup.Delete(
		"/:linkedAccountId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.DeleteLinkedAccountPermission),
		linkedaccounts.NewDeleteLinkedAccountHandler(handlers.DeleteLinkedAccount),
	)

	socialsRouterGroup := meRouterGroup.Group("/socials")

	socialsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadSocialPermission),
		usersocials.NewGetAllUserSocialsHandler(handlers.GetAllUserSocials),
	)

	socialsRouterGroup.Put(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateSocialPermission),
		usersocials.NewUpdateUserSocialHandler(handlers.UpdateUserSocial),
	)

	preferencesRouterGroup := meRouterGroup.Group("/preferences")

	preferencesRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadUserPreferencesPermission),
		userpreferences.NewGetUserPreferencesHandler(handlers.GetUserPreferences),
	)

	preferencesRouterGroup.Patch(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateUserPreferencesPermission),
		userpreferences.NewUpdateUserPreferencesHandler(handlers.UpdateUserPreferences),
	)

	sessionsRouterGroup := meRouterGroup.Group("/sessions")

	sessionsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadSessionPermission),
		sessions.NewGetAllSessionsHandler(handlers.GetAllSessions),
	)

	sessionsRouterGroup.Post(
		"/revoke-all",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.DeleteSessionPermission),
		sessions.NewRevokeAllSessionsHandler(handlers.RevokeAllSessions),
	)

	sessionsRouterGroup.Post(
		"/:sessionId/revoke",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.DeleteSessionPermission),
		sessions.NewRevokeSessionHandler(handlers.RevokeSession),
	)
}
