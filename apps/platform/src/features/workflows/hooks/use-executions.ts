import { useCallback } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import executionsService from '@/features/workflows/services/executions'
import taskRunnerService from '@/features/workflows/services/task-runner'

type ExecutionsParams = { query?: string; offset?: number; limit?: number }

export function useExecutions(
  organizationId: string | null | undefined,
  params: ExecutionsParams = {},
) {
  const { query = '', offset = 0, limit = 10 } = params
  const key = organizationId
    ? ['executions', organizationId, query, offset, limit]
    : null
  const { data, error, isLoading, mutate } = useSWR(
    key,
    () => executionsService.getAll(organizationId, { query, offset, limit }),
    { keepPreviousData: true },
  )
  return {
    executions: data?.data ?? [],
    pagination: data?.meta?.pagination,
    isLoading,
    error,
    mutate,
  }
}

export function useExecution(
  organizationId: string | null | undefined,
  executionId: string | null | undefined,
) {
  const key =
    organizationId && executionId
      ? ['execution', organizationId, executionId]
      : null
  const { data, error, isLoading, mutate } = useSWR(key, () =>
    executionsService.getById(organizationId, executionId),
  )
  return {
    execution: data?.data,
    isLoading,
    error,
    mutate,
  }
}

export function useExecutionActions(organizationId: string | null | undefined) {
  const { mutate } = useSWRConfig()

  const invalidateList = useCallback(
    () =>
      mutate(
        (k) =>
          Array.isArray(k) && k[0] === 'executions' && k[1] === organizationId,
        undefined,
        { revalidate: true },
      ),
    [organizationId, mutate],
  )

  const remove = useCallback(
    async (executionId: string) => {
      if (!organizationId) return null
      await executionsService.delete(organizationId, executionId)
      await invalidateList()
    },
    [organizationId, invalidateList],
  )

  const bulkRemove = useCallback(
    async (payload: unknown) => {
      if (!organizationId) return null
      await executionsService.bulkDelete(organizationId, payload as never)
      await invalidateList()
    },
    [organizationId, invalidateList],
  )

  const trigger = useCallback(
    async (payload: unknown) =>
      taskRunnerService.trigger(organizationId, payload as never),
    [organizationId],
  )

  const cancel = useCallback(
    async (payload: unknown) =>
      taskRunnerService.cancel(organizationId, payload as never),
    [organizationId],
  )

  const rerunAll = useCallback(
    async (payload: unknown) =>
      taskRunnerService.rerunAll(organizationId, payload as never),
    [organizationId],
  )

  const rerunFailed = useCallback(
    async (payload: unknown) =>
      taskRunnerService.rerunFailed(organizationId, payload as never),
    [organizationId],
  )

  return { remove, bulkRemove, trigger, cancel, rerunAll, rerunFailed }
}
