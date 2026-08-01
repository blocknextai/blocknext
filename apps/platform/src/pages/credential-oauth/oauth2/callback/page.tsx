import { useEffect, useState } from 'react'
import { hideAuthDialog } from '@/lib/broadcast'
import { credentialOAuthService } from '@/features/credentials'
import { PlatformLogo } from '@/features/navigation/components/logo'
import { useTranslation } from 'react-i18next'

const CredentialOAuth2CallbackPage = () => {
  const { t } = useTranslation()
  const [seconds, setSeconds] = useState(5)
  const goAuth2 = async () => {
    const q = window.location.search
    try {
      await credentialOAuthService.callback(q)
    } catch (error) {
      console.error('OAuth callback error:', error)
    } finally {
      hideAuthDialog()
      setTimeout(() => {
        window.close()
      }, 5000)
    }
  }

  useEffect(() => {
    goAuth2()

    const int = setInterval(() => {
      setSeconds((seconds) => {
        if (seconds === 0) {
          hideAuthDialog()
          window.close()
          clearInterval(int)
          return 0
        } else {
          return seconds - 1
        }
      })
    }, 1000)
  }, [])

  return (
    <div className="flex w-full h-app-screen items-center justify-center">
      <div className="bg-accent max-w-sm flex flex-col gap-6 p-4 rounded-lg shadow-xs">
        <PlatformLogo />
        <div className="heading-md">{t('ui.text.completed')}</div>
        <div className="text-muted-foreground text-sm">
          {t('ui.text.authCompletedDescription')}
        </div>
        <div className="text-xs">{t('ui.text.redirectionIn', { seconds })}</div>
      </div>
    </div>
  )
}

export default CredentialOAuth2CallbackPage
