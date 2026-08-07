import { useCallback } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import triggersService from '@/features/workflows/services/triggers'

type TriggersParams = { query?: string; offset?: number; limit?: number }

export function useTriggers(
  organizationId: string | null | undefined,
  params: TriggersParams = {},
) {
  const { query = '', offset = 0, limit = 10 } = params
  const key = organizationId
    ? ['triggers', organizationId, query, offset, limit]
    : null
  const { data, error, isLoading, mutate } = useSWR(
    key,
    () => triggersService.getAll(organizationId, { query, offset, limit }),
    { keepPreviousData: true },
  )
  return {
    triggers: data?.data ?? [],
    pagination: data?.meta?.pagination,
    isLoading,
    error,
    mutate,
  }
}

export function useTriggerActions(organizationId: string | null | undefined) {
  const { mutate } = useSWRConfig()

  const invalidate = useCallback(
    () =>
      mutate(
        (k) =>
          Array.isArray(k) && k[0] === 'triggers' && k[1] === organizationId,
        undefined,
        { revalidate: true },
      ),
    [organizationId, mutate],
  )

  const update = useCallback(
    async (triggerId: string, payload: unknown) => {
      if (!organizationId) {
        return null
      }
      const result = await triggersService.update(
        organizationId,
        triggerId,
        payload as never,
      )
      await invalidate()
      return result.data
    },
    [organizationId, invalidate],
  )

  const regenerateWebhookToken = useCallback(
    async (triggerId: string) => {
      if (!organizationId) {
        return null
      }
      const result = await triggersService.regenerateWebhookToken(
        organizationId,
        triggerId,
      )
      await invalidate()
      return result.data
    },
    [organizationId, invalidate],
  )

  const remove = useCallback(
    async (triggerId: string) => {
      if (!organizationId) {
        return null
      }
      const result = await triggersService.delete(organizationId, triggerId)
      await invalidate()
      return result.data
    },
    [organizationId, invalidate],
  )

  return { update, regenerateWebhookToken, remove }
}
