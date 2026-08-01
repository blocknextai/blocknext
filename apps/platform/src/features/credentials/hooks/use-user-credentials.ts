import { useCallback } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import credentialsService from '@/features/credentials/services/credentials'
import credentialOAuthService from '@/features/credentials/services/credential-oauth'
import { unwrap } from '@/lib/swr'

export function useUserCredentials() {
  const { data, error, isLoading, mutate } = useSWR('user-credentials', () =>
    unwrap(credentialsService.getUserCredentials()),
  )
  return {
    credentials: (data ?? []) as any[],
    isLoading,
    error,
    mutate,
  }
}

export function useUserCredential(credentialId: string | null | undefined) {
  const key = credentialId ? ['user-credential', credentialId] : null
  const { data, error, isLoading, mutate } = useSWR(key, () =>
    unwrap(credentialsService.getUserCredentialById(credentialId)),
  )
  return {
    credential: data,
    isLoading,
    error,
    mutate,
  }
}

export function useUserCredentialActions() {
  const { mutate } = useSWRConfig()

  const create = useCallback(
    async (payload: unknown) => {
      const result = await unwrap(
        credentialsService.createUserCredential(payload as never),
      )
      await mutate('user-credentials')
      return result
    },
    [mutate],
  )

  const update = useCallback(
    async (credentialId: string, payload: unknown) => {
      const result = await unwrap(
        credentialsService.updateUserCredential(credentialId, payload as never),
      )
      await Promise.all([
        mutate('user-credentials'),
        mutate(['user-credential', credentialId]),
      ])
      return result
    },
    [mutate],
  )

  const remove = useCallback(
    async (credentialId: string) => {
      const result = await unwrap(
        credentialsService.deleteUserCredential(credentialId),
      )
      await mutate('user-credentials')
      return result
    },
    [mutate],
  )

  return { create, update, remove }
}

export function useCredentialOAuthActions() {
  const authorizeUser = useCallback(async (credentialId: string) => {
    const result = await unwrap(
      credentialOAuthService.authorizeUser({ credentialId }),
    )
    return result
  }, [])

  const authorizeOrganization = useCallback(
    async (organizationId: string, credentialId: string) => {
      const result = await unwrap(
        credentialOAuthService.authorizeOrganization(organizationId, {
          credentialId,
        }),
      )
      return result
    },
    [],
  )

  const callback = useCallback(async (queryString: string) => {
    const result = await unwrap(credentialOAuthService.callback(queryString))
    return result
  }, [])

  return { authorizeUser, authorizeOrganization, callback }
}
