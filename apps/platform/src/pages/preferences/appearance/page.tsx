import { useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { useThemeStore } from '@/stores/theme-store'
import { useAuthActions } from '@/features/auth'
import { AppearanceSection } from '@/features/preferences/components/appearance-section'

const PreferencesAppearancePage = () => {
  const { i18n } = useTranslation()
  const setTheme = useThemeStore((s) => s.setTheme)
  const setColor = useThemeStore((s) => s.setColor)
  const mode = useThemeStore((s) => s.mode)
  const getMode = useThemeStore((s) => s.getMode)
  const color = useThemeStore((s) => s.color)

  const { updatePreferences } = useAuthActions()

  const changeLanguage = useCallback(
    (lng: string) => {
      i18n.changeLanguage(lng)
      updatePreferences({ language: lng }).catch((err: unknown) =>
        console.error('Error syncing preferences:', err),
      )
    },
    [i18n, updatePreferences],
  )

  return (
    <AppearanceSection
      mode={mode}
      setTheme={setTheme}
      getMode={getMode}
      color={color}
      setColor={setColor}
      changeLanguage={changeLanguage}
    />
  )
}

export default PreferencesAppearancePage
