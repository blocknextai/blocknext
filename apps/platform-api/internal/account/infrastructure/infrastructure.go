package infrastructure

import (
	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/go-packages/database"
	pkgEmail "github.com/blocknextai/go-packages/email"
	"github.com/blocknextai/go-packages/hashing"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusernonce"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusertoken"
	emailadd "github.com/blocknextai/platform-api/internal/account/application/auth/email/add"
	"github.com/blocknextai/platform-api/internal/account/application/auth/email/change"
	"github.com/blocknextai/platform-api/internal/account/application/auth/email/confirm"
	magiclinkconsume "github.com/blocknextai/platform-api/internal/account/application/auth/email/magiclink/consume"
	magiclinkrequest "github.com/blocknextai/platform-api/internal/account/application/auth/email/magiclink/request"
	"github.com/blocknextai/platform-api/internal/account/application/auth/email/resendverification"
	"github.com/blocknextai/platform-api/internal/account/application/auth/getauthmethods"
	"github.com/blocknextai/platform-api/internal/account/application/auth/mailer"
	passwordchange "github.com/blocknextai/platform-api/internal/account/application/auth/password/change"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/forgot"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/login"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/register"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/reset"
	passwordset "github.com/blocknextai/platform-api/internal/account/application/auth/password/set"
	"github.com/blocknextai/platform-api/internal/account/application/auth/password/verify"
	"github.com/blocknextai/platform-api/internal/account/application/auth/passwordpolicy"
	"github.com/blocknextai/platform-api/internal/account/application/auth/refreshtoken"
	"github.com/blocknextai/platform-api/internal/account/application/auth/tokenissuer"
	"github.com/blocknextai/platform-api/internal/account/application/auth/verificationtokenissuer"
	"github.com/blocknextai/platform-api/internal/account/application/events/emailadded"
	"github.com/blocknextai/platform-api/internal/account/application/events/emailchangerequested"
	"github.com/blocknextai/platform-api/internal/account/application/events/emailverificationrequested"
	"github.com/blocknextai/platform-api/internal/account/application/events/magiclinkcreated"
	"github.com/blocknextai/platform-api/internal/account/application/events/passwordresetrequested"
	"github.com/blocknextai/platform-api/internal/account/application/events/registrationexistingemailnotified"
	"github.com/blocknextai/platform-api/internal/account/application/events/userwelcome"
	"github.com/blocknextai/platform-api/internal/account/application/linkedaccounts/addlinkedaccount"
	"github.com/blocknextai/platform-api/internal/account/application/linkedaccounts/deletelinkedaccount"
	"github.com/blocknextai/platform-api/internal/account/application/linkedaccounts/getalllinkedaccounts"
	accountApplicationSessions "github.com/blocknextai/platform-api/internal/account/application/sessions"
	"github.com/blocknextai/platform-api/internal/account/application/sessions/getallsessions"
	"github.com/blocknextai/platform-api/internal/account/application/sessions/logout"
	"github.com/blocknextai/platform-api/internal/account/application/sessions/revokeallsessions"
	"github.com/blocknextai/platform-api/internal/account/application/sessions/revokesession"
	"github.com/blocknextai/platform-api/internal/account/application/userpreferences/getuserpreferences"
	"github.com/blocknextai/platform-api/internal/account/application/userpreferences/updateuserpreferences"
	"github.com/blocknextai/platform-api/internal/account/application/users/getprofile"
	"github.com/blocknextai/platform-api/internal/account/application/users/getroles"
	"github.com/blocknextai/platform-api/internal/account/application/usersocials/getallusersocials"
	"github.com/blocknextai/platform-api/internal/account/application/usersocials/updateusersocial"
	"github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	"github.com/blocknextai/platform-api/internal/account/domain/passwordcredentials"
	"github.com/blocknextai/platform-api/internal/account/domain/sessions"
	"github.com/blocknextai/platform-api/internal/account/domain/usernonces"
	"github.com/blocknextai/platform-api/internal/account/domain/userpreferences"
	"github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/blocknextai/platform-api/internal/account/domain/usersocials"
	"github.com/blocknextai/platform-api/internal/account/domain/verificationtokens"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/application/idempotency"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
)

type Handlers struct {
	CreateUserToken       cqrs.Handler[*createusertoken.CreateUserTokenCommand, *accountApplicationAuth.AccessTokenResponse]
	RefreshToken          cqrs.Handler[*refreshtoken.RefreshTokenCommand, *accountApplicationAuth.AccessTokenResponse]
	GetProfile            cqrs.Handler[*getprofile.GetProfileQuery, *getprofile.GetProfileResponse]
	GetRoles              cqrs.Handler[*getroles.GetRolesQuery, *getroles.GetRolesResponse]
	GetAuthMethods        cqrs.Handler[*getauthmethods.GetAuthMethodsQuery, *getauthmethods.GetAuthMethodsResponse]
	GetAllLinkedAccounts  cqrs.Handler[*getalllinkedaccounts.GetAllLinkedAccountsQuery, *getalllinkedaccounts.GetAllLinkedAccountsResponse]
	AddLinkedAccount      cqrs.Handler[*addlinkedaccount.AddLinkedAccountCommand, *addlinkedaccount.AddLinkedAccountResponse]
	DeleteLinkedAccount   cqrs.Handler[*deletelinkedaccount.DeleteLinkedAccountCommand, *deletelinkedaccount.DeleteLinkedAccountResponse]
	CreateUserNonce       cqrs.Handler[*createusernonce.CreateUserNonceCommand, *createusernonce.CreateUserNonceResponse]
	UpdateUserSocial      cqrs.Handler[*updateusersocial.UpdateUserSocialCommand, *updateusersocial.UpdateUserSocialResponse]
	GetAllUserSocials     cqrs.Handler[*getallusersocials.GetAllUserSocialsQuery, *getallusersocials.GetAllUserSocialsResponse]
	GetUserPreferences    cqrs.Handler[*getuserpreferences.GetUserPreferencesQuery, *getuserpreferences.GetUserPreferencesResponse]
	UpdateUserPreferences cqrs.Handler[*updateuserpreferences.UpdateUserPreferencesCommand, *updateuserpreferences.UpdateUserPreferencesResponse]
	GetAllSessions        cqrs.Handler[*getallsessions.GetAllSessionsQuery, *getallsessions.GetAllSessionsResponse]
	RevokeSession         cqrs.Handler[*revokesession.RevokeSessionCommand, *revokesession.RevokeSessionResponse]
	Logout                cqrs.Handler[*logout.LogoutCommand, *logout.LogoutResponse]
	RevokeAllSessions     cqrs.Handler[*revokeallsessions.RevokeAllSessionsCommand, *revokeallsessions.RevokeAllSessionsResponse]
	PasswordRegister      cqrs.Handler[*register.RegisterCommand, *accountApplicationAuth.AccessTokenResponse]
	PasswordLogin         cqrs.Handler[*login.LoginCommand, *accountApplicationAuth.AccessTokenResponse]
	PasswordVerify        cqrs.Handler[*verify.VerifyCommand, *accountApplicationAuth.AccessTokenResponse]
	PasswordForgot        cqrs.Handler[*forgot.ForgotCommand, *forgot.ForgotResponse]
	PasswordReset         cqrs.Handler[*reset.ResetCommand, *reset.ResetResponse]
	AddEmail              cqrs.Handler[*emailadd.AddEmailCommand, *emailadd.AddEmailResponse]
	ChangeEmail           cqrs.Handler[*change.ChangeEmailCommand, *change.ChangeEmailResponse]
	ConfirmEmailChange    cqrs.Handler[*confirm.ConfirmEmailChangeCommand, *confirm.ConfirmEmailChangeResponse]
	ResendVerification    cqrs.Handler[*resendverification.ResendVerificationCommand, *resendverification.ResendVerificationResponse]
	SetPassword           cqrs.Handler[*passwordset.SetPasswordCommand, *passwordset.SetPasswordResponse]
	ChangePassword        cqrs.Handler[*passwordchange.ChangePasswordCommand, *passwordchange.ChangePasswordResponse]
	MagicLinkRequest      cqrs.Handler[*magiclinkrequest.MagicLinkRequestCommand, *magiclinkrequest.MagicLinkRequestResponse]
	MagicLinkConsume      cqrs.Handler[*magiclinkconsume.MagicLinkConsumeCommand, *accountApplicationAuth.AccessTokenResponse]
}

type RegisterInfrastructureDeps struct {
	TransactionManager       database.TransactionManager
	EventBus                 *eventbus.Bus
	EventBusPublisherService publishing.PublisherService
	EventBusInboxService     *idempotency.InboxService

	AuthOptions       config.AuthOptions
	PlatformUIBaseURL string

	UserRepository               users.UserRepository
	UserNonceRepository          usernonces.UserNonceRepository
	LinkedAccountRepository      linkedaccounts.LinkedAccountRepository
	UserSocialRepository         usersocials.UserSocialRepository
	UserPreferenceRepository     userpreferences.UserPreferenceRepository
	SessionRepository            sessions.SessionRepository
	PasswordCredentialRepository passwordcredentials.PasswordCredentialRepository
	VerificationTokenRepository  verificationtokens.VerificationTokenRepository
	SessionService               accountApplicationSessions.SessionService
	AuthProviderRegistry         createusertoken.AuthProviderRegistry
	AuthJWTService               jwt.AuthJWTService
	PasswordHasher               hashing.Hasher
	EmailSender                  pkgEmail.EmailSender
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	tokenIssuer := tokenissuer.NewService(deps.AuthJWTService, deps.SessionService, deps.UserRepository)
	verificationIssuer := verificationtokenissuer.NewService(deps.VerificationTokenRepository)
	policyChecker := passwordpolicy.NewChecker()
	mailerInstance := mailer.NewMailer(deps.EmailSender, deps.PlatformUIBaseURL)

	emailverificationrequested.New(mailerInstance, deps.EventBus)
	emailadded.New(mailerInstance, deps.EventBus)
	passwordresetrequested.New(mailerInstance, deps.EventBus)
	emailchangerequested.New(mailerInstance, deps.EventBus)
	magiclinkcreated.New(mailerInstance, deps.EventBus)
	registrationexistingemailnotified.New(mailerInstance, deps.EventBus)
	userwelcome.New(mailerInstance, deps.EventBus)

	return &Handlers{
		CreateUserToken:       cqrs.ValidationBehavior(createusertoken.New(deps.UserRepository, deps.UserNonceRepository, deps.LinkedAccountRepository, deps.UserPreferenceRepository, tokenIssuer, deps.AuthProviderRegistry, deps.EventBusPublisherService, deps.TransactionManager)),
		RefreshToken:          cqrs.ValidationBehavior(refreshtoken.New(deps.UserRepository, deps.SessionService, deps.AuthJWTService)),
		GetProfile:            cqrs.ValidationBehavior(getprofile.New(deps.UserRepository)),
		GetRoles:              cqrs.ValidationBehavior(getroles.New()),
		GetAuthMethods:        cqrs.ValidationBehavior(getauthmethods.New(deps.AuthProviderRegistry, deps.AuthOptions.Password.Enabled, deps.AuthOptions.MagicLink.Enabled)),
		GetAllLinkedAccounts:  cqrs.ValidationBehavior(getalllinkedaccounts.New(deps.UserRepository, deps.LinkedAccountRepository)),
		AddLinkedAccount:      cqrs.ValidationBehavior(addlinkedaccount.New(deps.UserRepository, deps.UserNonceRepository, deps.LinkedAccountRepository, deps.AuthProviderRegistry, deps.TransactionManager)),
		DeleteLinkedAccount:   cqrs.ValidationBehavior(deletelinkedaccount.New(deps.LinkedAccountRepository, deps.PasswordCredentialRepository, deps.TransactionManager)),
		CreateUserNonce:       cqrs.ValidationBehavior(createusernonce.New(deps.UserNonceRepository, deps.AuthProviderRegistry, deps.TransactionManager)),
		UpdateUserSocial:      cqrs.ValidationBehavior(updateusersocial.New(deps.UserSocialRepository, deps.TransactionManager)),
		GetAllUserSocials:     cqrs.ValidationBehavior(getallusersocials.New(deps.UserSocialRepository)),
		GetUserPreferences:    cqrs.ValidationBehavior(getuserpreferences.New(deps.UserPreferenceRepository)),
		UpdateUserPreferences: cqrs.ValidationBehavior(updateuserpreferences.New(deps.UserPreferenceRepository)),
		GetAllSessions:        cqrs.ValidationBehavior(getallsessions.New(deps.SessionRepository)),
		RevokeSession:         cqrs.ValidationBehavior(revokesession.New(deps.SessionService, deps.SessionRepository)),
		Logout:                cqrs.ValidationBehavior(logout.New(deps.SessionService)),
		RevokeAllSessions:     cqrs.ValidationBehavior(revokeallsessions.New(deps.SessionService)),
		PasswordRegister: cqrs.ValidationBehavior(register.New(
			deps.UserRepository,
			deps.LinkedAccountRepository,
			deps.PasswordCredentialRepository,
			deps.UserPreferenceRepository,
			verificationIssuer,
			tokenIssuer,
			policyChecker,
			deps.PasswordHasher,
			deps.AuthOptions.Password.VerificationTokenTTL,
			deps.EventBusPublisherService,
			deps.TransactionManager,
		)),
		PasswordLogin: cqrs.ValidationBehavior(login.New(
			deps.LinkedAccountRepository,
			deps.PasswordCredentialRepository,
			tokenIssuer,
			deps.PasswordHasher,
			deps.TransactionManager,
		)),
		PasswordVerify: cqrs.ValidationBehavior(verify.New(
			deps.LinkedAccountRepository,
			deps.VerificationTokenRepository,
			tokenIssuer,
			deps.TransactionManager,
		)),
		PasswordForgot: cqrs.ValidationBehavior(forgot.New(
			deps.LinkedAccountRepository,
			verificationIssuer,
			deps.AuthOptions.Password.PasswordResetTokenTTL,
			deps.EventBusPublisherService,
			deps.TransactionManager,
		)),
		PasswordReset: cqrs.ValidationBehavior(reset.New(
			deps.LinkedAccountRepository,
			deps.PasswordCredentialRepository,
			deps.VerificationTokenRepository,
			policyChecker,
			deps.PasswordHasher,
			deps.TransactionManager,
		)),
		AddEmail: cqrs.ValidationBehavior(emailadd.New(
			deps.LinkedAccountRepository,
			verificationIssuer,
			deps.AuthOptions.Password.VerificationTokenTTL,
			deps.EventBusPublisherService,
			deps.TransactionManager,
		)),
		ChangeEmail: cqrs.ValidationBehavior(change.New(
			deps.LinkedAccountRepository,
			deps.PasswordCredentialRepository,
			verificationIssuer,
			deps.PasswordHasher,
			deps.AuthOptions.Password.VerificationTokenTTL,
			deps.EventBusPublisherService,
			deps.TransactionManager,
		)),
		ConfirmEmailChange: cqrs.ValidationBehavior(confirm.New(
			deps.LinkedAccountRepository,
			deps.VerificationTokenRepository,
			deps.TransactionManager,
		)),
		ResendVerification: cqrs.ValidationBehavior(resendverification.New(
			deps.LinkedAccountRepository,
			verificationIssuer,
			deps.AuthOptions.Password.VerificationTokenTTL,
			deps.EventBusPublisherService,
			deps.TransactionManager,
		)),
		SetPassword: cqrs.ValidationBehavior(passwordset.New(
			deps.LinkedAccountRepository,
			deps.PasswordCredentialRepository,
			policyChecker,
			deps.PasswordHasher,
			deps.TransactionManager,
		)),
		ChangePassword: cqrs.ValidationBehavior(passwordchange.New(
			deps.LinkedAccountRepository,
			deps.PasswordCredentialRepository,
			policyChecker,
			deps.PasswordHasher,
			deps.EventBusPublisherService,
			deps.TransactionManager,
		)),
		MagicLinkRequest: cqrs.ValidationBehavior(magiclinkrequest.New(
			deps.LinkedAccountRepository,
			verificationIssuer,
			deps.AuthOptions.MagicLink.TokenTTL,
			deps.EventBusPublisherService,
			deps.TransactionManager,
		)),
		MagicLinkConsume: cqrs.ValidationBehavior(magiclinkconsume.New(
			deps.LinkedAccountRepository,
			deps.VerificationTokenRepository,
			tokenIssuer,
			deps.TransactionManager,
		)),
	}
}
