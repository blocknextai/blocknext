import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Footer } from '@/features/navigation/components/footer'
import { PlatformLogo } from '@/features/navigation/components/logo'
import { Outlet } from 'react-router'
import { landingUrl } from '@/constants'
import { ModeToggle } from '@/features/theme/components/mode-toggle'
import { Image } from '@/components/shared/image'
import { useThemeStore } from '@/stores/theme-store'

const AuthLayout = () => {
  const { t } = useTranslation()
  const [landingImg, setLandingImg] = useState('landing-w.png')
  const mode = useThemeStore((s) => s.mode)
  const getMode = useThemeStore((s) => s.getMode)

  useEffect(() => {
    const resolved = getMode()
    setLandingImg(resolved === 'dark' ? 'landing-b.png' : 'landing-w.png')
  }, [mode, getMode])

  return (
    <div className="flex flex-col h-app-screen">
      <div className="flex min-h-app-screen flex-col bg-background">
        <div className="flex flex-1 md:flex">
          <Image
            src={`/assets/images/${landingImg}`}
            alt={t('ui.text.platformPreview')}
            className="hidden w-1/2 object-cover md:block"
          />
          <div className="flex flex-1 flex-col md:w-1/2">
            <ModeToggle
              className="shrink-0 absolute top-4 right-4"
              aria-label={t('ui.text.chooseTheme')}
            />
            <div className="flex flex-1 items-center justify-center">
              <div className="w-full max-w-sm">
                <div className="mb-6 flex items-start justify-between gap-4">
                  <a
                    href={landingUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <PlatformLogo />
                  </a>
                </div>

                <Outlet />
              </div>
            </div>
            <Footer />
          </div>
        </div>
      </div>
    </div>
  )
}

export default AuthLayout
