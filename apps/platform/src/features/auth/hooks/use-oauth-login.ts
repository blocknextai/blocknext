import { useCallback } from 'react'
import authService from '@/features/auth/services/auth'

export function useOAuthLogin(provider: string, mode: string) {
  return useCallback(async () => {
    try {
      const result = await authService.getNonce({ authProvider: provider })
      const response = result.data

      if (response.url) {
        const urlParams = new URLSearchParams(window.location.search)
        const returnUrl = urlParams.get('returnUrl')

        const stateObj = {
          state: response.nonce,
          mode: mode || 'login',
          returnUrl: returnUrl || '',
        }

        const encodedState = btoa(JSON.stringify(stateObj))

        const url = new URL(response.url)
        url.searchParams.set('state', encodedState)
        window.location.href = url.toString()
      }
    } catch (error) {}
  }, [provider, mode])
}
