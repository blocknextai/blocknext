import { ethers } from 'ethers'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import tokenManager from '@/lib/token-manager'
import { getReturnUrl } from '@/lib/auth-redirect'
import { ProviderIcon } from '@/features/auth/components/provider-icon'
import { useAuthRedirect } from '@/features/auth/hooks/use-auth-redirect'

const MetaMaskLoginButton = ({ mode = 'login', authActions }) => {
  const { t } = useTranslation()

  const isInstalled = typeof window.ethereum !== 'undefined'
  const redirectAfterAuth = useAuthRedirect()

  const handleMetaMaskLogin = async () => {
    try {
      const provider = new ethers.BrowserProvider(window.ethereum)
      const signer = await provider.getSigner()
      const address = await signer.getAddress()

      const nonceResult = await authActions.getNonce({
        authProvider: 'metamask',
        providerId: address,
      })
      const { nonce, loginMessage } = nonceResult.data
      const signature = await signer.signMessage(loginMessage)

      if (mode === 'linked_accounts') {
        await authActions.linkAccount({
          authProvider: 'metamask',
          flow: 'linked_accounts',
          payload: {
            nonce,
            walletAddress: address,
            signature,
          },
        })
      } else {
        const tokenResponse = await authActions.getToken({
          authProvider: 'metamask',
          flow: 'login',
          payload: {
            nonce,
            walletAddress: address,
            signature,
          },
        })
        const tokens = tokenResponse.data
        tokenManager.setTokens(tokens.accessToken, tokens.refreshToken)
        await redirectAfterAuth(getReturnUrl())
      }
    } catch (error) {
      console.error('MetaMask login error:', error)
    }
  }

  return (
    <Button
      onClick={handleMetaMaskLogin}
      variant="default"
      disabled={!isInstalled}
      title={!isInstalled ? t('ui.text.metamaskNotInstalled') : undefined}
      className="w-full cursor-pointer bg-primary text-primary-foreground"
    >
      <ProviderIcon provider="metamask" className="w-5 h-5 mr-2" />
      {t('ui.text.continue_with', { provider: 'MetaMask' })}
    </Button>
  )
}

MetaMaskLoginButton.displayName = 'MetaMaskLoginButton'

export { MetaMaskLoginButton }
