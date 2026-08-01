import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { ProviderIcon } from '@/features/auth/components/provider-icon'
import { getProviderName } from '@/features/auth/components/provider-utils'
import { useOAuthLogin } from '@/features/auth/hooks/use-oauth-login'

const OAuthLoginButton = ({ provider, mode = 'login' }) => {
  const { t } = useTranslation()
  const handleLogin = useOAuthLogin(provider, mode)

  const text = t('ui.text.continue_with', {
    provider: getProviderName(provider),
  })

  return (
    <Button
      onClick={handleLogin}
      variant="outline"
      className="w-full cursor-pointer"
    >
      <ProviderIcon provider={provider} className="w-5 h-5 mr-2" />
      {text}
    </Button>
  )
}

OAuthLoginButton.displayName = 'OAuthLoginButton'

export { OAuthLoginButton }
