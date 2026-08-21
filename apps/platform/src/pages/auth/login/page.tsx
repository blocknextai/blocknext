import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { Eye, EyeOff, Loader2, MailCheck } from 'lucide-react'
import { authService, useAuthMethods, useAuthRedirect } from '@/features/auth'
import { BackToSignInLink } from '@/features/auth/components/back-to-sign-in-link'
import { ProviderLoginButton } from '@/features/auth/components/provider-login-button'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Loading } from '@/components/shared/loading'
import tokenManager from '@/lib/token-manager'
import { getReturnUrl } from '@/lib/auth-redirect'
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

type Step = 'email' | 'methods' | 'magic-sent'

function AuthLoginPage() {
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

  const emailMethodsEnabled = !!methods.password || !!methods.magicLink

  const [step, setStep] = useState<Step>('email')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [signingIn, setSigningIn] = useState(false)
  const [sendingMagicLink, setSendingMagicLink] = useState(false)
  const redirectAfterAuth = useAuthRedirect()

  useEffect(() => {
    useOrganizationStore.getState().reset()
  }, [])

  const handleContinue = (e: React.FormEvent) => {
    e.preventDefault()
    if (!email.trim()) {
      return
    }
    setStep('methods')
  }

  const handlePasswordLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email || !password) {
      return
    }
    setSigningIn(true)
    try {
      const response = await authService.loginWithPassword({ email, password })
      const tokens = response.data
      if (!response.isSuccess || !tokens?.accessToken) {
        return
      }
      tokenManager.setTokens(tokens.accessToken, tokens.refreshToken)
      await redirectAfterAuth(getReturnUrl())
    } catch (error) {
      console.error('Password login error:', error)
    } finally {
      setSigningIn(false)
    }
  }

  const handleSendMagicLink = async () => {
    if (!email.trim()) {
      return
    }
    setSendingMagicLink(true)
    try {
      const response = await authService.requestMagicLink({ email })
      if (!response.isSuccess) {
        return
      }
      setStep('magic-sent')
    } catch (error) {
      console.error('Magic link error:', error)
    } finally {
      setSendingMagicLink(false)
    }
  }

  const resetToEmail = () => {
    setStep('email')
    setPassword('')
  }

  return (
    <>
      <div className="mb-6 flex flex-col gap-y-3">
        <h1 className="text-2xl font-bold md:text-3xl">
          {t('ui.text.sign_in.title')}
        </h1>
        <p className="text-sm text-muted-foreground">
          {t('ui.text.sign_in.description')}
        </p>
      </div>

      {step !== 'email' && <BackToSignInLink onClick={resetToEmail} />}

      {isLoading ? (
        <Loading />
      ) : (
        <div className="flex flex-col space-y-4">
          {providers.length > 0 && step === 'email' && (
            <div className="flex flex-col space-y-3">
              {providers.map((provider) => (
                <ProviderLoginButton
                  key={provider}
                  provider={provider}
                  authActions={authActions}
                />
              ))}
            </div>
          )}

          {emailMethodsEnabled && step === 'email' && (
            <>
              {providers.length > 0 && (
                <Divider label={t('ui.text.or', 'or')} />
              )}
              <form onSubmit={handleContinue} className="space-y-3">
                <div className="grid gap-2">
                  <Label htmlFor="login-email">
                    {t('ui.text.emailAddress')}
                  </Label>
                  <Input
                    id="login-email"
                    type="email"
                    autoComplete="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                  />
                </div>
                <Button type="submit" className="w-full">
                  {t('ui.text.continue')}
                </Button>
              </form>
              <div className="text-center text-sm text-muted-foreground">
                {t('ui.text.no_account')}{' '}
                <Link
                  to="/auth/register"
                  className="text-primary hover:underline"
                >
                  {t('ui.text.create_account')}
                </Link>
              </div>
            </>
          )}

          {step === 'methods' && (
            <div className="flex flex-col space-y-4">
              <div className="grid gap-2">
                <Label htmlFor="login-email-locked">
                  {t('ui.text.emailAddress')}
                </Label>
                <Input
                  id="login-email-locked"
                  type="email"
                  value={email}
                  disabled
                  className="bg-muted/50"
                />
              </div>

              {methods.password && (
                <form onSubmit={handlePasswordLogin} className="space-y-3">
                  <div className="grid gap-2">
                    <div className="flex items-center justify-between">
                      <Label htmlFor="login-password">
                        {t('ui.text.password')}
                      </Label>
                      <Link
                        to="/auth/forgot-password"
                        className="text-xs text-primary hover:underline"
                      >
                        {t('ui.text.forgot_password')}
                      </Link>
                    </div>
                    <div className="relative">
                      <Input
                        id="login-password"
                        type={showPassword ? 'text' : 'password'}
                        autoComplete="current-password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className="pr-10"
                        required
                        autoFocus
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        tabIndex={-1}
                        aria-label={
                          showPassword ? t('ui.text.hide') : t('ui.text.show')
                        }
                        className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 p-0"
                        onClick={() => setShowPassword((s) => !s)}
                      >
                        {showPassword ? (
                          <EyeOff className="h-3.5 w-3.5" />
                        ) : (
                          <Eye className="h-3.5 w-3.5" />
                        )}
                      </Button>
                    </div>
                  </div>
                  <Button type="submit" className="w-full" disabled={signingIn}>
                    {signingIn && (
                      <Loader2 className="w-4 h-4 animate-spin mr-2" />
                    )}
                    {t('ui.text.sign_in')}
                  </Button>
                </form>
              )}

              {methods.password && methods.magicLink && (
                <Divider label={t('ui.text.or', 'or')} />
              )}

              {methods.magicLink && (
                <Button
                  type="button"
                  variant="outline"
                  className="w-full"
                  disabled={sendingMagicLink}
                  onClick={handleSendMagicLink}
                >
                  {sendingMagicLink && (
                    <Loader2 className="w-4 h-4 animate-spin mr-2" />
                  )}
                  {t('ui.text.send_magic_link')}
                </Button>
              )}
            </div>
          )}

          {step === 'magic-sent' && (
            <div className="flex items-start gap-3 p-4 rounded-lg bg-green-500/10 border border-green-500/20 text-sm">
              <MailCheck className="size-5 text-green-600 dark:text-green-400 shrink-0 mt-0.5" />
              <div className="text-green-700 dark:text-green-300">
                <p className="font-medium">{t('ui.text.magic_link_sent')}</p>
                <p className="mt-1 text-xs opacity-90">
                  {t('ui.text.magic_link_sent_description', { email })}
                </p>
              </div>
            </div>
          )}
        </div>
      )}
    </>
  )
}

export default AuthLoginPage
