import { useEffect } from 'react'
import { Navigate, Outlet, useLocation } from 'react-router'
import { useTranslation } from 'react-i18next'
import tokenManager from '@/lib/token-manager'
import { buildLoginUrl } from '@/lib/auth-redirect'
import { useAuthPreferences } from '@/features/auth'
import { useThemeStore } from '@/stores/theme-store'

const ProtectedLayout = () => {
  const location = useLocation()
  const { i18n } = useTranslation()
  const syncFromBackend = useThemeStore((s) => s.syncFromBackend)
  const hasToken = !!tokenManager.getAccessToken()
  const { preferences } = useAuthPreferences()

  useEffect(() => {
    if (!preferences) {
      return
    }
    syncFromBackend(preferences.theme ?? {})
    if (preferences.language && preferences.language !== i18n.language) {
      i18n.changeLanguage(preferences.language)
    }
  }, [preferences, syncFromBackend, i18n])

  if (!hasToken) {
    return (
      <Navigate
        to={buildLoginUrl(location.pathname + location.search)}
        replace
      />
    )
  }

  return (
    <div className="min-h-app-screen bg-background">
      <Outlet />
    </div>
  )
}

export default ProtectedLayout
