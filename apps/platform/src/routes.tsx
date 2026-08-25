import { Suspense } from 'react'
import { Route, Routes as RRoutes, BrowserRouter, Navigate } from 'react-router'
import * as Pages from '@/pages'
import {
  OrganizationLayout,
  FlowLayout,
  RootLayout,
  AuthLayout,
  ProtectedLayout,
  RequireOrganization,
} from '@/layouts'
import PreferencesLayout from '@/features/preferences/components/preferences-layout'
import { ErrorBoundary } from '@/components/shared/error-boundary'
import { PageLoading } from '@/components/shared/loading'

function Routes() {
  return (
    <BrowserRouter>
      <ErrorBoundary>
        <Suspense fallback={<PageLoading />}>
          <RRoutes>
            <Route element={<RootLayout />}>
              <Route element={<AuthLayout />}>
                <Route path="auth">
                  <Route path="login">
                    <Route index element={<Pages.AuthLoginPage />} />
                    <Route
                      path=":provider"
                      element={<Pages.AuthLoginCallbackPage />}
                    />
                  </Route>
                  <Route path="register" element={<Pages.AuthRegisterPage />} />
                  <Route
                    path="forgot-password"
                    element={<Pages.AuthForgotPasswordPage />}
                  />
                  <Route
                    path="reset-password"
                    element={<Pages.AuthResetPasswordPage />}
                  />
                  <Route
                    path="verify-email"
                    element={<Pages.AuthVerifyEmailPage />}
                  />
                  <Route
                    path="magic-link/callback"
                    element={<Pages.AuthMagicLinkCallbackPage />}
                  />
                  <Route
                    path="email-change/confirm"
                    element={<Pages.AuthEmailChangeConfirmPage />}
                  />
                </Route>
              </Route>

              <Route path="credential-oauth/oauth2">
                <Route
                  path="callback"
                  element={<Pages.CredentialOAuth2CallbackPage />}
                />
              </Route>

              <Route element={<ProtectedLayout />}>
                <Route
                  path="organizations/new"
                  element={<Pages.OrganizationCreatePage />}
                />

                <Route element={<RequireOrganization />}>
                  <Route element={<OrganizationLayout />}>
                    <Route
                      path="/"
                      element={<Navigate to="/organizations" replace />}
                    />

                    <Route
                      path="organizations"
                      element={<Pages.OrganizationsPage />}
                    />

                    <Route path="mcp" element={<Pages.McpPage />} />

                    <Route path="preferences" element={<PreferencesLayout />}>
                      <Route
                        index
                        element={<Navigate to="/preferences/profile" replace />}
                      />
                      <Route
                        path="profile"
                        element={<Pages.PreferencesProfilePage />}
                      />
                      <Route
                        path="appearance"
                        element={<Pages.PreferencesAppearancePage />}
                      />
                      <Route
                        path="sessions"
                        element={<Pages.PreferencesSessionsPage />}
                      />
                      <Route
                        path="notifications"
                        element={<Pages.PreferencesNotificationsPage />}
                      />
                    </Route>
                  </Route>
                  <Route path="/organizations/:organizationId">
                    <Route element={<OrganizationLayout />}>
                      <Route index element={<Pages.OrganizationDetailPage />} />

                      <Route path="history">
                        <Route
                          index
                          element={<Pages.OrganizationHistoryPage />}
                        />
                        <Route
                          path=":runId"
                          element={<Pages.OrganizationHistoryDetailPage />}
                        />
                      </Route>

                      <Route path="run">
                        <Route
                          path=":id"
                          element={<Pages.OrganizationRunPage />}
                        />
                      </Route>

                      <Route path="triggers">
                        <Route
                          index
                          element={<Pages.OrganizationTriggersPage />}
                        />
                      </Route>

                      <Route path="api-keys">
                        <Route
                          index
                          element={<Pages.OrganizationAPIKeysPage />}
                        />
                      </Route>

                      <Route path="notifications">
                        <Route
                          index
                          element={<Pages.OrganizationNotificationsPage />}
                        />
                      </Route>

                      <Route path="credentials">
                        <Route
                          index
                          element={<Pages.OrganizationCredentialsPage />}
                        />
                      </Route>

                      <Route path="settings">
                        <Route
                          index
                          element={<Pages.OrganizationSettingsPage />}
                        />
                        <Route
                          path="members"
                          element={<Pages.OrganizationSettingsMembersPage />}
                        />
                      </Route>
                    </Route>

                    <Route element={<FlowLayout />}>
                      <Route
                        path="create"
                        element={<Pages.FlowsCreatePage />}
                      />
                      <Route
                        path="edit/:id"
                        element={<Pages.FlowsEditPage />}
                      />
                      <Route path="welcome" element={<Pages.WelcomePage />} />
                    </Route>
                  </Route>
                </Route>
              </Route>

              <Route path="*" element={<Pages.NotFoundPage />} />
            </Route>
          </RRoutes>
        </Suspense>
      </ErrorBoundary>
    </BrowserRouter>
  )
}

export default Routes
