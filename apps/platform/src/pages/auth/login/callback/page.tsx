import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams, useParams, Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import toast from '@/lib/toast'
import { authService } from '@/features/auth'
import { getProviderName } from '@/features/auth/components/provider-utils'
import { ProviderIcon } from '@/features/auth/components/provider-icon'
import { Button } from '@/components/ui/button'
import { Loading } from '@/components/shared/loading'
import tokenManager from '@/lib/token-manager'

function AuthLoginCallbackPage() {
  const { t } = useTranslation()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const { provider } = useParams()
  const [isLoading, setIsLoading] = useState(true)
  const [hasError, setHasError] = useState(false)
  const [mode, setMode] = useState('login')

  useEffect(() => {
    const handleCallback = async () => {
      try {
        // Validate provider
        if (!provider) {
          toast.error(t('ui.text.invalidCallbackUrl'))
          setHasError(true)
          setIsLoading(false)
          return
        }

        // Decode state parameter to get mode, original state, and returnUrl
        const encodedState = searchParams.get('state')
        let originalState = ''
        let currentMode = 'login'
        let returnUrl = ''

        if (encodedState) {
          try {
            const stateObj = JSON.parse(atob(encodedState))
            currentMode = stateObj.mode || 'login'
            originalState = stateObj.state || ''
            returnUrl = stateObj.returnUrl || ''
          } catch (error) {
            console.error('Failed to decode state:', error)
            toast.error(t('ui.text.invalidStateParameter'))
            setHasError(true)
            setIsLoading(false)
            return
          }
        }

        // Set mode in state
        setMode(currentMode)
        const isLinkedAccounts = currentMode === 'linked_accounts'

        if (provider === 'metamask') {
          // MetaMask doesn't use OAuth, so we don't expect code/state
          // This would be handled differently if needed
          toast.error(t('ui.text.metamaskCallbackNotSupported'))
          setHasError(true)
          setIsLoading(false)
          return
        }

        // For OAuth providers (Google, GitHub)
        const code = searchParams.get('code')

        if (!code || !originalState) {
          toast.error(t('ui.text.invalidOAuthResponse'))
          setHasError(true)
          setIsLoading(false)
          return
        }

        if (isLinkedAccounts) {
          // For linked accounts, send to different endpoint
          await authService.linkAccount({
            authProvider: provider,
            payload: {
              code,
              state: originalState,
            },
          })

          // Redirect back to preferences
          navigate('/preferences')
        } else {
          // For login, use the original flow
          const tokens = await authService.getToken({
            authProvider: provider,
            payload: {
              code,
              state: originalState,
            },
          })

          if (!tokens) {
            throw new Error(
              `${provider} authentication failed. Please try again.`,
            )
          }

          // Set tokens and redirect
          tokenManager.setTokens(
            tokens.data?.accessToken,
            tokens.data?.refreshToken,
          )

          // Redirect to returnUrl from state or dashboard
          navigate(returnUrl || '/')
        }
      } catch (error) {
        console.error(
          `${getProviderName(provider)} authentication error:`,
          error,
        )
        setHasError(true)
        setIsLoading(false)
      }
    }

    handleCallback()
  }, [searchParams, navigate, provider])

  const getActionText = () => {
    const isLinkedAccounts = mode === 'linked_accounts'
    return isLinkedAccounts
      ? t('ui.text.linkingAccount')
      : t('ui.text.signingIn')
  }

  return (
    <div className="flex flex-col gap-y-3 items-start justify-center">
      {/* Provider Logo */}
      <div className="flex">
        <ProviderIcon provider={provider} className="w-16 h-16" />
      </div>

      {isLoading && (
        <>
          {/* Title and Description */}
          <h1 className="text-2xl font-semibold">
            {t('ui.text.completingAuth', {
              provider: getProviderName(provider),
              action: getActionText(),
            })}
          </h1>
          <p className="text-sm leading-relaxed">
            {t('ui.text.authWaitDescription', {
              provider: getProviderName(provider),
            })}
          </p>

          {/* Loading Animation */}
          <Loading />
        </>
      )}

      {hasError && (
        <>
          {/* Error Title and Description */}
          <h1 className="text-2xl font-semibold">
            {t('ui.text.authenticationFailed')}
          </h1>
          <p className="text-sm leading-relaxed">
            {t('ui.text.authFailedDescription', {
              provider: getProviderName(provider),
            })}
          </p>

          {/* Back to Login Button */}
          <Link to="/auth/login">
            <Button>{t('ui.text.backToLogin')}</Button>
          </Link>
        </>
      )}
    </div>
  )
}

export default AuthLoginCallbackPage
