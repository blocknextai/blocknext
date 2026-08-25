import { useCallback, useMemo, useState } from 'react'
import {
  useAuthLinkedAccounts,
  useAuthMethods,
  useAuthActions,
  authService,
} from '@/features/auth'
import { ProfileSection } from '@/features/preferences/components/profile-section'
import tokenManager from '@/lib/token-manager'

const PreferencesProfilePage = () => {
  const {
    linkedAccounts,
    isLoading: loading,
    mutate: refetchLinked,
  } = useAuthLinkedAccounts()
  const { methods: authMethods } = useAuthMethods()
  const { unlinkAccount } = useAuthActions()

  const [deletingAccountId, setDeletingAccountId] = useState<string | null>(
    null,
  )
  const [resendingEmail, setResendingEmail] = useState<string | null>(null)

  const userId = tokenManager.getUserId() ?? ''

  const availableProviders = useMemo(
    () => [...(authMethods.crypto ?? []), ...(authMethods.oauth ?? [])],
    [authMethods],
  )

  const linkedProviderNames = useMemo(
    () => linkedAccounts.map((account: any) => account.authProvider),
    [linkedAccounts],
  )

  const getUnlinkedProviders = useCallback(
    () =>
      availableProviders.filter(
        (p: string) => !linkedProviderNames.includes(p),
      ),
    [availableProviders, linkedProviderNames],
  )

  const deleteLinkedAccount = useCallback(
    async (accountId: string) => {
      setDeletingAccountId(accountId)
      try {
        await unlinkAccount(accountId)
      } finally {
        setDeletingAccountId(null)
      }
    },
    [unlinkAccount],
  )

  const handleAccountLinked = useCallback(async () => {
    await refetchLinked()
  }, [refetchLinked])

  const handleResendVerification = useCallback(async (email: string) => {
    setResendingEmail(email)
    try {
      await authService.resendEmailVerification({ email })
    } finally {
      setResendingEmail(null)
    }
  }, [])

  return (
    <ProfileSection
      userId={userId}
      loading={loading}
      linkedAccounts={linkedAccounts}
      deletingAccountId={deletingAccountId}
      deleteLinkedAccount={deleteLinkedAccount}
      getUnlinkedProviders={getUnlinkedProviders}
      handleAccountLinked={handleAccountLinked}
      passwordEnabled={authMethods.password === true}
      resendingEmail={resendingEmail}
      onResendVerification={handleResendVerification}
    />
  )
}

export default PreferencesProfilePage
