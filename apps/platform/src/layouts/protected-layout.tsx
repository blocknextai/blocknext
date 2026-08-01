import { useEffect } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router'
import { useTranslation } from 'react-i18next'
import tokenManager from '@/lib/token-manager'
import { useAuthPreferences } from '@/features/auth'
import { useThemeStore } from '@/stores/theme-store'

const ProtectedLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { i18n } = useTranslation()
  const syncFromBackend = useThemeStore((s) => s.syncFromBackend)
  const { preferences } = useAuthPreferences()

  useEffect(() => {
    const token = tokenManager.getAccessToken()
    if (!token) {
      const currentPath = location.pathname + location.search
      if (
        currentPath &&
        currentPath !== '/' &&
        !currentPath.startsWith('/auth/login')
      ) {
        const returnUrl = encodeURIComponent(currentPath)
        navigate(`/auth/login?returnUrl=${returnUrl}`, { replace: true })
      } else {
        navigate('/auth/login', { replace: true })
      }
    }
  }, [navigate, location])

  useEffect(() => {
    if (!preferences) return
    syncFromBackend(preferences.theme ?? {})
    if (preferences.language && preferences.language !== i18n.language) {
      i18n.changeLanguage(preferences.language)
    }
  }, [preferences, syncFromBackend, i18n])

  return (
    <div className="min-h-app-screen bg-background">
      <Outlet />
    </div>
  )
}

export default ProtectedLayout
