import { OAuthLoginButton } from '@/features/auth/components/providers/oauth-login-button'

const ProviderLoginButton = ({ provider, mode = 'login' }) => {
  return (
    <div className="space-y-2">
      <OAuthLoginButton provider={provider} mode={mode} />
    </div>
  )
}

ProviderLoginButton.displayName = 'ProviderLoginButton'

export { ProviderLoginButton }
