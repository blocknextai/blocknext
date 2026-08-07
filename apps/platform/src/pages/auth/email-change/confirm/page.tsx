import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { AlertCircle, CheckCircle2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Loading } from '@/components/shared/loading'
import { authService } from '@/features/auth'

type Status = 'pending' | 'success' | 'error'

function AuthEmailChangeConfirmPage() {
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
      .confirmEmailChange({ token })
      .then((response) => {
        if (cancelled) {
          return
        }
        setStatus(response.isSuccess ? 'success' : 'error')
      })
      .catch(() => {
        if (!cancelled) {
          setStatus('error')
        }
      })
    return () => {
      cancelled = true
    }
  }, [token])

  return (
    <>
      <div className="mb-6 flex flex-col gap-y-3">
        <h1 className="text-2xl font-bold md:text-3xl">
          {t('ui.text.confirm_email_change_title')}
        </h1>
      </div>

      {status === 'pending' && <Loading />}

      {status === 'success' && (
        <div className="flex flex-col gap-4">
          <div className="flex items-start gap-3 p-4 rounded-lg bg-green-500/10 border border-green-500/20 text-sm">
            <CheckCircle2 className="size-5 text-green-600 dark:text-green-400 shrink-0 mt-0.5" />
            <span className="text-green-700 dark:text-green-300">
              {t('ui.text.email_change_confirmed')}
            </span>
          </div>
          <Button asChild className="w-full">
            <Link to="/">{t('ui.text.continue')}</Link>
          </Button>
        </div>
      )}

      {status === 'error' && (
        <div className="flex flex-col gap-4">
          <div className="flex items-start gap-3 p-4 rounded-lg bg-red-500/10 border border-red-500/20 text-sm">
            <AlertCircle className="size-5 text-red-600 dark:text-red-400 shrink-0 mt-0.5" />
            <span className="text-red-700 dark:text-red-300">
              {t('ui.text.email_change_failed')}
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

export default AuthEmailChangeConfirmPage
