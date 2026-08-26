import { lazy } from 'react'

export const AuthLoginPage = lazy(() => import('@/pages/auth/login/page'))
export const AuthLoginCallbackPage = lazy(
  () => import('@/pages/auth/login/callback/page'),
)
export const AuthRegisterPage = lazy(() => import('@/pages/auth/register/page'))
export const AuthForgotPasswordPage = lazy(
  () => import('@/pages/auth/forgot-password/page'),
)
export const AuthResetPasswordPage = lazy(
  () => import('@/pages/auth/reset-password/page'),
)
export const AuthVerifyEmailPage = lazy(
  () => import('@/pages/auth/verify-email/page'),
)
export const AuthMagicLinkCallbackPage = lazy(
  () => import('@/pages/auth/magic-link/callback/page'),
)
export const AuthEmailChangeConfirmPage = lazy(
  () => import('@/pages/auth/email-change/confirm/page'),
)

export const FlowsCreatePage = lazy(() => import('@/pages/flows/create/page'))
export const FlowsEditPage = lazy(() => import('@/pages/flows/edit/[id]/page'))

export const OrganizationsPage = lazy(
  () => import('@/pages/organizations/page'),
)
export const OrganizationCreatePage = lazy(
  () => import('@/pages/organizations/new/page'),
)
export const OrganizationDetailPage = lazy(
  () => import('@/pages/organizations/[organizationId]/(root)/page'),
)
export const OrganizationSettingsPage = lazy(
  () => import('@/pages/organizations/[organizationId]/settings/page'),
)
export const OrganizationCredentialsPage = lazy(
  () => import('@/pages/organizations/[organizationId]/credentials/page'),
)
export const OrganizationSettingsMembersPage = lazy(
  () => import('@/pages/organizations/[organizationId]/settings/members/page'),
)
export const OrganizationHistoryPage = lazy(
  () => import('@/pages/organizations/[organizationId]/history/page'),
)
export const OrganizationHistoryDetailPage = lazy(
  () => import('@/pages/organizations/[organizationId]/history/[runId]/page'),
)
export const OrganizationTriggersPage = lazy(
  () => import('@/pages/organizations/[organizationId]/triggers/page'),
)
export const OrganizationAPIKeysPage = lazy(
  () => import('@/pages/organizations/[organizationId]/api-keys/page'),
)
export const OrganizationNotificationsPage = lazy(
  () => import('@/pages/organizations/[organizationId]/notifications/page'),
)
export const OrganizationRunPage = lazy(
  () => import('@/pages/organizations/[organizationId]/run/[id]/page'),
)
export const PreferencesProfilePage = lazy(
  () => import('@/pages/preferences/profile/page'),
)
export const PreferencesAppearancePage = lazy(
  () => import('@/pages/preferences/appearance/page'),
)
export const PreferencesSessionsPage = lazy(
  () => import('@/pages/preferences/sessions/page'),
)
export const PreferencesNotificationsPage = lazy(
  () => import('@/pages/preferences/notifications/page'),
)

export const McpPage = lazy(() => import('@/pages/mcp/page'))

export const CredentialOAuth2CallbackPage = lazy(
  () => import('@/pages/credential-oauth/oauth2/callback/page'),
)

export const NotFoundPage = lazy(() => import('@/pages/notfound/page'))
