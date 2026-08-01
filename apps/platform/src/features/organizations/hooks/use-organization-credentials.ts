import { useCallback } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import organizationsService from '@/features/organizations/services/organizations'
import { unwrap } from '@/lib/swr'

export function useOrganizationCredentials(
  organizationId: string | null | undefined,
) {
  const key = organizationId
    ? ['organization-credentials', organizationId]
    : null
  const { data, error, isLoading, mutate } = useSWR(key, () =>
    unwrap(organizationsService.getOrganizationCredentials(organizationId)),
  )
  return {
    credentials: data ?? [],
    isLoading,
    error,
    mutate,
  }
}

export function useOrganizationCredential(
  organizationId: string | null | undefined,
  credentialId: string | null | undefined,
) {
  const key =
    organizationId && credentialId
      ? ['organization-credential', organizationId, credentialId]
      : null
  const { data, error, isLoading, mutate } = useSWR(key, () =>
    unwrap(
      organizationsService.getOrganizationCredentialById(
        organizationId,
        credentialId,
      ),
    ),
  )
  return {
    credential: data,
    isLoading,
    error,
    mutate,
  }
}

export function useOrganizationCredentialActions(
  organizationId: string | null | undefined,
) {
  const { mutate } = useSWRConfig()
  const listKey = organizationId
    ? ['organization-credentials', organizationId]
    : null

  const create = useCallback(
    async (payload: unknown) => {
      if (!organizationId) return null
      const result = await unwrap(
        organizationsService.createOrganizationCredential(
          organizationId,
          payload as never,
        ),
      )
      await mutate(listKey)
      return result
    },
    [organizationId, listKey, mutate],
  )

  const update = useCallback(
    async (credentialId: string, payload: unknown) => {
      if (!organizationId) return null
      const result = await unwrap(
        organizationsService.updateOrganizationCredential(
          organizationId,
          credentialId,
          payload as never,
        ),
      )
      await Promise.all([
        mutate(listKey),
        mutate(['organization-credential', organizationId, credentialId]),
      ])
      return result
    },
    [organizationId, listKey, mutate],
  )

  const remove = useCallback(
    async (credentialId: string) => {
      if (!organizationId) return null
      const result = await unwrap(
        organizationsService.deleteOrganizationCredential(
          organizationId,
          credentialId,
        ),
      )
      await mutate(listKey)
      return result
    },
    [organizationId, listKey, mutate],
  )

  return { create, update, remove }
}
