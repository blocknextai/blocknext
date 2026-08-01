import { useCallback } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import organizationsService from '@/features/organizations/services/organizations'
import { unwrap } from '@/lib/swr'

export function useOrganizations() {
  const { data, error, isLoading, mutate } = useSWR('organizations', () =>
    unwrap(organizationsService.getAll()),
  )
  return {
    organizations: data ?? [],
    isLoading,
    error,
    mutate,
  }
}

export function useOrganization(organizationId: string | null | undefined) {
  const key = organizationId ? ['organization', organizationId] : null
  const { data, error, isLoading, mutate } = useSWR(key, () =>
    unwrap(organizationsService.getById(organizationId)),
  )
  return {
    organization: data,
    isLoading,
    error,
    mutate,
  }
}

export function useOrganizationMe(organizationId: string | null | undefined) {
  const key = organizationId ? ['organization-me', organizationId] : null
  const { data, error, isLoading, mutate } = useSWR(key, () =>
    unwrap(organizationsService.getMyMembership(organizationId)),
  )
  return {
    me: data,
    isLoading,
    error,
    mutate,
  }
}

export function useOrganizationActions() {
  const { mutate } = useSWRConfig()

  const create = useCallback(
    async (payload: unknown) => {
      const result = await unwrap(organizationsService.create(payload as never))
      await mutate('organizations')
      return result
    },
    [mutate],
  )

  const update = useCallback(
    async (organizationId: string, payload: unknown) => {
      const result = await unwrap(
        organizationsService.update(organizationId, payload as never),
      )
      await Promise.all([
        mutate('organizations'),
        mutate(['organization', organizationId]),
      ])
      return result
    },
    [mutate],
  )

  const remove = useCallback(
    async (organizationId: string) => {
      const result = await unwrap(organizationsService.delete(organizationId))
      await mutate('organizations')
      return result
    },
    [mutate],
  )

  return { create, update, remove }
}
