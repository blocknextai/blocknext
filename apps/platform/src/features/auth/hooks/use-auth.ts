import { useCallback } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import authService from '@/features/auth/services/auth'
import { unwrap } from '@/lib/swr'
import tokenManager from '@/lib/token-manager'

export interface AuthMethods {
  crypto?: string[]
  oauth?: string[]
  password?: boolean
  magicLink?: boolean
}

export function useAuthMe() {
  const hasToken = !!tokenManager.getAccessToken()
  const { data, error, isLoading, mutate } = useSWR(
    hasToken ? 'auth-me' : null,
    () => unwrap(authService.getMe()),
  )
  return { me: data, isLoading, error, mutate }
}

export function useAuthRoles() {
  const hasToken = !!tokenManager.getAccessToken()
  const { data, error, isLoading, mutate } = useSWR(
    hasToken ? 'auth-roles' : null,
    () => unwrap(authService.getRoles()),
  )
  return { roles: data, isLoading, error, mutate }
}

export function useAuthPreferences() {
  const hasToken = !!tokenManager.getAccessToken()
  const { data, error, isLoading, mutate } = useSWR(
    hasToken ? 'auth-preferences' : null,
    () => unwrap(authService.getPreferences()),
  )
  return { preferences: data, isLoading, error, mutate }
}

export function useAuthLinkedAccounts() {
  const hasToken = !!tokenManager.getAccessToken()
  const { data, error, isLoading, mutate } = useSWR(
    hasToken ? 'auth-linked-accounts' : null,
    () => unwrap(authService.getLinkedAccounts()),
  )
  return { linkedAccounts: data ?? [], isLoading, error, mutate }
}

export function useAuthSessions(offset = 0, limit = 10) {
  const hasToken = !!tokenManager.getAccessToken()
  const key = hasToken ? ['auth-sessions', offset, limit] : null
  const { data, error, isLoading, mutate } = useSWR(key, () =>
    authService.getSessions(offset, limit),
  )
  return {
    sessions: data?.data ?? [],
    pagination: data?.meta?.pagination,
    isLoading,
    error,
    mutate,
  }
}

export function useAuthMethods() {
  const { data, error, isLoading } = useSWR<AuthMethods>('auth-methods', () =>
    unwrap(authService.getMethods()),
  )
  return { methods: (data ?? {}) as AuthMethods, isLoading, error }
}

export function useAuthActions() {
  const { mutate } = useSWRConfig()

  const updatePreferences = useCallback(
    async (payload: unknown) => {
      const result = await unwrap(
        authService.updatePreferences(payload as never),
      )
      await mutate('auth-preferences')
      return result
    },
    [mutate],
  )

  const linkAccount = useCallback(
    async (payload: unknown) => {
      const result = await unwrap(authService.linkAccount(payload as never))
      await mutate('auth-linked-accounts')
      return result
    },
    [mutate],
  )

  const unlinkAccount = useCallback(
    async (accountId: string) => {
      const result = await unwrap(authService.unlinkAccount(accountId as never))
      await mutate('auth-linked-accounts')
      return result
    },
    [mutate],
  )

  const revokeSession = useCallback(
    async (sessionId: string) => {
      const result = await unwrap(authService.revokeSession(sessionId as never))
      await mutate(
        (k) => Array.isArray(k) && k[0] === 'auth-sessions',
        undefined,
        { revalidate: true },
      )
      return result
    },
    [mutate],
  )

  const revokeAllSessions = useCallback(async () => {
    const result = await unwrap(authService.revokeAllSessions())
    await mutate(
      (k) => Array.isArray(k) && k[0] === 'auth-sessions',
      undefined,
      { revalidate: true },
    )
    return result
  }, [mutate])

  return {
    updatePreferences,
    linkAccount,
    unlinkAccount,
    revokeSession,
    revokeAllSessions,
  }
}
