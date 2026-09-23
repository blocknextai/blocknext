package http

import (
	"time"

	"github.com/gofiber/fiber/v3"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/http/auth"
	authEmail "github.com/blocknextai/platform-api/internal/modules/account/internal/http/auth/email"
	authMagicLink "github.com/blocknextai/platform-api/internal/modules/account/internal/http/auth/email/magiclink"
	authPassword "github.com/blocknextai/platform-api/internal/modules/account/internal/http/auth/password"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/http/linkedaccounts"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/http/sessions"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/http/userpreferences"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/http/users"
	usersEmail "github.com/blocknextai/platform-api/internal/modules/account/internal/http/users/email"
	usersPassword "github.com/blocknextai/platform-api/internal/modules/account/internal/http/users/password"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/http/usersocials"
	accountUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases"
)

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	cacheMiddleware *cachemiddleware.Middleware,
	useCases *accountUseCases.Services,
) {
	registerAuthRoutes(router, cacheMiddleware, authMiddleware, useCases)
	registerUserRoutes(router, authMiddleware, useCases)
}

func registerAuthRoutes(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *accountUseCases.Services,
) {
	authRouterGroup := router.Group("/auth")

	authRouterGroup.Get(
		"/methods",
		cacheMiddleware.Cache(5*time.Minute),
		auth.NewGetAuthMethodsHandler(useCases.Auth),
	)

	authRouterGroup.Post(
		"/token/refresh",
		auth.NewRefreshTokenHandler(useCases.Auth),
	)

	authRouterGroup.Post(
		"/logout",
		authMiddleware.Authenticate(),
		sessions.NewLogoutHandler(useCases.Sessions),
	)

	oauthRouterGroup := authRouterGroup.Group("/oauth")

	oauthRouterGroup.Post(
		"/nonce",
		auth.NewCreateUserNonceHandler(useCases.Auth),
	)

	oauthRouterGroup.Post(
		"/token",
		auth.NewCreateUserTokenHandler(useCases.Auth),
	)

	passwordRouterGroup := authRouterGroup.Group("/password")

	passwordRouterGroup.Post(
		"/register",
		authPassword.NewRegisterHandler(useCases.Auth),
	)

	passwordRouterGroup.Post(
		"/login",
		authPassword.NewLoginHandler(useCases.Auth),
	)

	passwordRouterGroup.Post(
		"/forgot",
		authPassword.NewForgotHandler(useCases.Auth),
	)

	passwordRouterGroup.Post(
		"/reset",
		authPassword.NewResetHandler(useCases.Auth),
	)

	emailRouterGroup := authRouterGroup.Group("/email")

	emailRouterGroup.Post(
		"/verify",
		authEmail.NewVerifyHandler(useCases.Auth),
	)

	emailRouterGroup.Post(
		"/resend-verification",
		authEmail.NewResendVerificationHandler(useCases.Auth),
	)

	emailRouterGroup.Post(
		"/change/confirm",
		authEmail.NewConfirmEmailChangeHandler(useCases.Auth),
	)

	emailRouterGroup.Post(
		"/magic-link/request",
		authMagicLink.NewMagicLinkRequestHandler(useCases.Auth),
	)

	emailRouterGroup.Post(
		"/magic-link/consume",
		authMagicLink.NewMagicLinkConsumeHandler(useCases.Auth),
	)
}

func registerUserRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *accountUseCases.Services,
) {
	meRouterGroup := router.Group("/users/me")

	meRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadUserPermission),
		users.NewGetProfileHandler(useCases.Users),
	)

	meRouterGroup.Get(
		"/roles",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadUserPermission),
		users.NewGetRolesHandler(useCases.Users),
	)

	emailMeRouterGroup := meRouterGroup.Group("/email")

	emailMeRouterGroup.Post(
		"/add",
		authMiddleware.Authenticate(),
		usersEmail.NewAddEmailHandler(useCases.Auth),
	)

	emailMeRouterGroup.Post(
		"/change",
		authMiddleware.Authenticate(),
		usersEmail.NewChangeEmailHandler(useCases.Auth),
	)

	passwordRouterGroup := meRouterGroup.Group("/password")

	passwordRouterGroup.Post(
		"/set",
		authMiddleware.Authenticate(),
		usersPassword.NewSetPasswordHandler(useCases.Auth),
	)

	passwordRouterGroup.Post(
		"/change",
		authMiddleware.Authenticate(),
		usersPassword.NewChangePasswordHandler(useCases.Auth),
	)

	linkedAccountsRouterGroup := meRouterGroup.Group("/linked-accounts")

	linkedAccountsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadLinkedAccountsPermission),
		linkedaccounts.NewGetAllLinkedAccountsHandler(useCases.LinkedAccounts),
	)

	linkedAccountsRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.CreateLinkedAccountPermission),
		linkedaccounts.NewAddLinkedAccountHandler(useCases.LinkedAccounts),
	)

	linkedAccountsRouterGroup.Delete(
		"/:linkedAccountId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.DeleteLinkedAccountPermission),
		linkedaccounts.NewDeleteLinkedAccountHandler(useCases.LinkedAccounts),
	)

	socialsRouterGroup := meRouterGroup.Group("/socials")

	socialsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadSocialPermission),
		usersocials.NewGetAllUserSocialsHandler(useCases.UserSocials),
	)

	socialsRouterGroup.Put(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateSocialPermission),
		usersocials.NewUpdateUserSocialHandler(useCases.UserSocials),
	)

	preferencesRouterGroup := meRouterGroup.Group("/preferences")

	preferencesRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadUserPreferencesPermission),
		userpreferences.NewGetUserPreferencesHandler(useCases.Userpreferences),
	)

	preferencesRouterGroup.Patch(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateUserPreferencesPermission),
		userpreferences.NewUpdateUserPreferencesHandler(useCases.Userpreferences),
	)

	sessionsRouterGroup := meRouterGroup.Group("/sessions")

	sessionsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadSessionPermission),
		sessions.NewGetAllSessionsHandler(useCases.Sessions),
	)

	sessionsRouterGroup.Post(
		"/revoke-all",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.DeleteSessionPermission),
		sessions.NewRevokeAllSessionsHandler(useCases.Sessions),
	)

	sessionsRouterGroup.Post(
		"/:sessionId/revoke",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.DeleteSessionPermission),
		sessions.NewRevokeSessionHandler(useCases.Sessions),
	)
}
