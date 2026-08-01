import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { MailCheck } from 'lucide-react'
import { authService, useAuthMethods } from '@/features/auth'
import { ProviderLoginButton } from '@/features/auth/components/provider-login-button'
import { PasswordRegisterForm } from '@/features/auth/components/providers/password-register-form'
import { Button } from '@/components/ui/button'
import { Loading } from '@/components/shared/loading'
import { useOrganizationStore } from '@/stores/organization'

const Divider = ({ label }: { label: string }) => (
  <div className="flex items-center gap-2 my-2">
    <div className="flex-1 h-px bg-border" />
    <span className="text-xs uppercase text-muted-foreground tracking-wider">
      {label}
    </span>
    <div className="flex-1 h-px bg-border" />
  </div>
)

function AuthRegisterPage() {
  const { t } = useTranslation()
  const { methods, isLoading } = useAuthMethods()

  const authActions = useMemo(
    () => ({
      getNonce: (data) => authService.getNonce(data),
      getToken: (data) => authService.getToken(data),
      linkAccount: (data) => authService.linkAccount(data),
    }),
    [],
  )

  const providers = useMemo(
    () => [...(methods.oauth ?? []), ...(methods.crypto ?? [])],
    [methods],
  )

  const [submittedEmail, setSubmittedEmail] = useState<string | null>(null)

  useEffect(() => {
    useOrganizationStore.getState().reset()
  }, [])

  if (submittedEmail) {
    return (
      <>
        <div className="mb-6 flex flex-col gap-y-3">
          <h1 className="text-2xl font-bold md:text-3xl">
            {t('ui.text.register_verification_title')}
          </h1>
        </div>
        <div className="flex flex-col gap-4">
          <div className="flex items-start gap-3 p-4 rounded-lg bg-green-500/10 border border-green-500/20 text-sm">
            <MailCheck className="size-5 text-green-600 dark:text-green-400 shrink-0 mt-0.5" />
            <div className="text-green-700 dark:text-green-300">
              <p className="font-medium">
                {t('ui.text.register_verification_sent')}
              </p>
              <p className="mt-1 text-xs opacity-90">
                {t('ui.text.register_verification_sent_description', {
                  email: submittedEmail,
                })}
              </p>
            </div>
          </div>
          <Button asChild variant="outline" className="w-full">
            <Link to="/auth/login">{t('ui.text.back_to_sign_in')}</Link>
          </Button>
        </div>
      </>
    )
  }

  return (
    <>
      <div className="mb-6 flex flex-col gap-y-3">
        <h1 className="text-2xl font-bold md:text-3xl">
          {t('ui.text.sign_up.title')}
        </h1>
        <p className="text-sm text-muted-foreground">
          {t('ui.text.sign_up.description')}
        </p>
      </div>

      {isLoading ? (
        <Loading />
      ) : (
        <div className="flex flex-col space-y-4">
          {methods.password && (
            <PasswordRegisterForm onSubmitted={setSubmittedEmail} />
          )}

          {providers.length > 0 && (
            <>
              {methods.password && <Divider label={t('ui.text.or', 'or')} />}
              <div className="flex flex-col space-y-3">
                {providers.map((provider) => (
                  <ProviderLoginButton
                    key={provider}
                    provider={provider}
                    mode="register"
                    authActions={authActions}
                  />
                ))}
              </div>
            </>
          )}

          <div className="text-center text-sm text-muted-foreground mt-4">
            {t('ui.text.already_have_account')}{' '}
            <Link to="/auth/login" className="text-primary hover:underline">
              {t('ui.text.sign_in')}
            </Link>
          </div>
        </div>
      )}
    </>
  )
}

export default AuthRegisterPage
