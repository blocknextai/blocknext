import CredentialTab from '@/features/credentials/components/credential-tab'
import { useCredentials } from '@/features/credentials/hooks/use-credentials'
import {
  credentialsService,
  credentialOAuthService,
} from '@/features/credentials'

const PreferencesCredentialsPage = () => {
  const credentialData = useCredentials(null, true)

  const handleCredentialSave = async (payload) => {
    try {
      if (!payload.ref) {
        await credentialsService.createUserCredential(payload)
      } else {
        await credentialsService.updateUserCredential(payload.ref, payload)
      }
      credentialData.getSecrets()
    } catch (error) {
      console.error('Error saving credentials:', error)
    }
  }

  const handleCredentialDelete = async (ref) => {
    try {
      await credentialsService.deleteUserCredential(ref)
      credentialData.getSecrets()
    } catch (error) {
      console.error('Error deleting credential:', error)
    }
  }

  const handleAuthorizeOAuth2 = async (payload) => {
    try {
      const resResponse = await credentialsService.createUserCredential(payload)
      const res = resResponse.data
      const authUrlResponse = await credentialOAuthService.authorizeUser({
        credentialId: res.id,
      })
      const authUrl = authUrlResponse.data
      const popupWindow = window.open('', '_blank')
      if (popupWindow) {
        popupWindow.location.href = authUrl.url
      }
    } catch (error) {
      console.error('Error during OAuth2 authorization:', error)
    }
  }

  const handleReauthorizeOAuth2 = async (credentialId) => {
    try {
      const authUrlResponse = await credentialOAuthService.authorizeUser({
        credentialId,
      })
      const authUrl = authUrlResponse.data
      const popupWindow = window.open('', '_blank')
      if (popupWindow) {
        popupWindow.location.href = authUrl.url
      }
    } catch (error) {
      console.error('Error during OAuth2 re-authorization:', error)
    }
  }

  return (
    <CredentialTab
      isUserMode={true}
      loading={credentialData.loading}
      secrets={credentialData.secrets}
      credentialsArray={credentialData.credentialsArray}
      pagination={credentialData.pagination}
      hasAccess={credentialData.hasAccess}
      onFetch={credentialData.fetchCredentials}
      getSecrets={credentialData.getSecrets}
      getSecretDetails={credentialData.getSecretDetails}
      onSave={handleCredentialSave}
      onDelete={handleCredentialDelete}
      onAuthorizeOAuth2={handleAuthorizeOAuth2}
      onReauthorizeOAuth2={handleReauthorizeOAuth2}
    />
  )
}

export default PreferencesCredentialsPage
