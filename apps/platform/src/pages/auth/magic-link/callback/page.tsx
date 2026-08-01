import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Loading } from '@/components/shared/loading'
import { authService } from '@/features/auth'
import tokenManager from '@/lib/token-manager'

type Status = 'pending' | 'error'

function AuthMagicLinkCallbackPage() {
  const { t } = useTranslation()
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token') ?? ''
  const [status, setStatus] = useState<Status>('pending')

  useEffect(() => {
    if (!token) {
      setStatus('error')
      return
    }
    let cancelled = false
    authService
      .consumeMagicLink({ token })
      .then((response) => {
        if (cancelled) return
        const tokens = response.data
        if (!response.isSuccess || !tokens?.accessToken) {
          setStatus('error')
          return
        }
        tokenManager.setTokens(tokens.accessToken, tokens.refreshToken)
        const urlParams = new URLSearchParams(window.location.search)
        const returnUrl = urlParams.get('returnUrl')
        window.location.href = returnUrl || '/'
      })
      .catch(() => {
        if (!cancelled) setStatus('error')
      })
    return () => {
      cancelled = true
    }
  }, [token])

  return (
    <>
      <div className="mb-6 flex flex-col gap-y-3">
        <h1 className="text-2xl font-bold md:text-3xl">
          {t('ui.text.magic_link_consuming_title')}
        </h1>
      </div>

      {status === 'pending' && <Loading />}

      {status === 'error' && (
        <div className="flex flex-col gap-4">
          <div className="flex items-start gap-3 p-4 rounded-lg bg-red-500/10 border border-red-500/20 text-sm">
            <AlertCircle className="size-5 text-red-600 dark:text-red-400 shrink-0 mt-0.5" />
            <span className="text-red-700 dark:text-red-300">
              {t('ui.text.magic_link_consume_failed')}
            </span>
          </div>
          <Button asChild variant="outline" className="w-full">
            <Link to="/auth/login">{t('ui.text.sign_in')}</Link>
          </Button>
        </div>
      )}
    </>
  )
}

export default AuthMagicLinkCallbackPage
