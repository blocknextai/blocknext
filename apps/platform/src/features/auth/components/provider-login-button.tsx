import { lazy, Suspense } from 'react'
import { OAuthLoginButton } from '@/features/auth/components/providers/oauth-login-button'

const MetaMaskLoginButton = lazy(() =>
  import('@/features/auth/components/providers/metamask').then((m) => ({
    default: m.MetaMaskLoginButton,
  })),
)

const ProviderLoginButton = ({ provider, mode = 'login', authActions }) => {
  if (provider === 'metamask') {
    return (
      <div className="space-y-2">
        <Suspense fallback={null}>
          <MetaMaskLoginButton
            provider={provider}
            mode={mode}
            authActions={authActions}
          />
        </Suspense>
      </div>
    )
  }

  return (
    <div className="space-y-2">
      <OAuthLoginButton provider={provider} mode={mode} />
    </div>
  )
}

ProviderLoginButton.displayName = 'ProviderLoginButton'

export { ProviderLoginButton }
