import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useThemeStore } from '@/stores/theme-store'
import { Image } from '@/components/shared/image'

const PlatformLogo = ({ width = 200, height = 40, className = '' }) => {
  const { t } = useTranslation()
  const mode = useThemeStore((s) => s.mode)
  const getMode = useThemeStore((s) => s.getMode)

  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const [logoName, setLogoName] = useState(
    prefersDark ? 'logo-w.png' : 'logo-b.png',
  )

  useEffect(() => {
    const theme = getMode()
    if (theme === 'dark') {
      setLogoName('logo-b.png')
    } else if (theme === 'light') {
      setLogoName('logo-w.png')
    } else {
      setLogoName(prefersDark ? 'logo-b.png' : 'logo-w.png')
    }
  }, [mode, getMode])

  return (
    <Image
      priority
      src={`/assets/images/${logoName}`}
      alt={t('ui.text.platformLogo')}
      width={width}
      height={height}
      className={className}
    />
  )
}
PlatformLogo.displayName = 'PlatformLogo'

export { PlatformLogo }
