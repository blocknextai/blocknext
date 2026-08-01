import { useCallback } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import organizationsService from '@/features/organizations/services/organizations'
import { unwrap } from '@/lib/swr'

type MembersParams = { offset?: number; limit?: number; query?: string }

export function useOrganizationMembers(
  organizationId: string | null | undefined,
  params: MembersParams = {},
) {
  const { offset = 0, limit = 10, query = '' } = params
  const key = organizationId
    ? ['organization-members', organizationId, offset, limit, query]
    : null
  const { data, error, isLoading, mutate } = useSWR(
    key,
    () =>
      organizationsService.getMembers(organizationId, { offset, limit, query }),
    { keepPreviousData: true },
  )
  return {
    members: data?.data ?? [],
    pagination: data?.meta?.pagination,
    isLoading,
    error,
    mutate,
  }
}

export function useOrganizationRoles() {
  const { data, error, isLoading } = useSWR('organization-roles', () =>
    unwrap(organizationsService.getRoles()),
  )
  return {
    roles: data ?? [],
    isLoading,
    error,
  }
}

export function useOrganizationMemberActions(
  organizationId: string | null | undefined,
) {
  const { mutate } = useSWRConfig()

  const invalidate = useCallback(
    () =>
      mutate(
        (k) =>
          Array.isArray(k) &&
          k[0] === 'organization-members' &&
          k[1] === organizationId,
        undefined,
        { revalidate: true },
      ),
    [organizationId, mutate],
  )

  const invite = useCallback(
    async (payload: unknown) => {
      if (!organizationId) return null
      const result = await unwrap(
        organizationsService.inviteMember(organizationId, payload as never),
      )
      await invalidate()
      return result
    },
    [organizationId, invalidate],
  )

  const updateRole = useCallback(
    async (memberId: string, payload: unknown) => {
      if (!organizationId) return null
      const result = await unwrap(
        organizationsService.updateMemberRole(
          organizationId,
          memberId,
          payload as never,
        ),
      )
      await invalidate()
      return result
    },
    [organizationId, invalidate],
  )

  const updateInfo = useCallback(
    async (memberId: string, payload: unknown) => {
      if (!organizationId) return null
      const result = await unwrap(
        organizationsService.updateMemberInfo(
          organizationId,
          memberId,
          payload as never,
        ),
      )
      await invalidate()
      return result
    },
    [organizationId, invalidate],
  )

  const remove = useCallback(
    async (memberId: string) => {
      if (!organizationId) return null
      const result = await unwrap(
        organizationsService.deleteMember(organizationId, memberId),
      )
      await invalidate()
      return result
    },
    [organizationId, invalidate],
  )

  return { invite, updateRole, updateInfo, remove }
}
