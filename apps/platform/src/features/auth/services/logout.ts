import authService from '@/features/auth/services/auth'
import tokenManager from '@/lib/token-manager'
import { LOGIN_PATH } from '@/lib/auth-redirect'
import { useOrganizationStore } from '@/stores/organization'

export const logoutAndRedirect = async () => {
  try {
    await authService.logout()
  } catch {
    // proceed with local logout even if backend call fails
  }
  tokenManager.clearTokens()
  useOrganizationStore.getState().reset()
  window.location.href = LOGIN_PATH
}
