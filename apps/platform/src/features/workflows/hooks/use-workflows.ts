import { useCallback } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import workflowsService from '@/features/workflows/services/workflows'
import { unwrap } from '@/lib/swr'

type WorkflowsParams = { offset?: number; limit?: number; query?: string }

export function useWorkflows(
  organizationId: string | null | undefined,
  params: WorkflowsParams = {},
) {
  const { offset = 0, limit = 12, query = '' } = params
  const key = organizationId
    ? ['workflows', organizationId, offset, limit, query]
    : null
  const { data, error, isLoading, mutate } = useSWR(
    key,
    () => workflowsService.getAll(organizationId, { offset, limit, query }),
    { keepPreviousData: true },
  )
  return {
    workflows: (data?.data ?? []) as any[],
    pagination: data?.meta?.pagination,
    isLoading,
    error,
    mutate,
  }
}

export function useWorkflow(
  organizationId: string | null | undefined,
  workflowId: string | null | undefined,
) {
  const key =
    organizationId && workflowId
      ? ['workflow', organizationId, workflowId]
      : null
  const { data, error, isLoading, mutate } = useSWR(key, () =>
    unwrap(workflowsService.getById(organizationId, workflowId)),
  )
  return {
    workflow: data,
    isLoading,
    error,
    mutate,
  }
}

export function useWorkflowRun(
  organizationId: string | null | undefined,
  workflowId: string | null | undefined,
) {
  const key =
    organizationId && workflowId
      ? ['workflow-run', organizationId, workflowId]
      : null
  const { data, error, isLoading, mutate } = useSWR(key, () =>
    unwrap(workflowsService.getRun(organizationId, workflowId)),
  )
  return {
    run: data,
    isLoading,
    error,
    mutate,
  }
}

export function useWorkflowActions(organizationId: string | null | undefined) {
  const { mutate } = useSWRConfig()

  const invalidate = useCallback(async () => {
    await Promise.all([
      mutate(
        (k) =>
          Array.isArray(k) && k[0] === 'workflows' && k[1] === organizationId,
        undefined,
        { revalidate: true },
      ),
      mutate(
        (k) =>
          Array.isArray(k) && k[0] === 'workflow' && k[1] === organizationId,
        undefined,
        { revalidate: true },
      ),
    ])
  }, [organizationId, mutate])

  const create = useCallback(
    async (payload: unknown) => {
      if (!organizationId) {
        return null
      }
      const result = await unwrap(
        workflowsService.create(organizationId, payload as never),
      )
      await invalidate()
      return result
    },
    [organizationId, invalidate],
  )

  const update = useCallback(
    async (workflowId: string, payload: unknown) => {
      if (!organizationId) {
        return null
      }
      const result = await unwrap(
        workflowsService.update(organizationId, workflowId, payload as never),
      )
      await invalidate()
      return result
    },
    [organizationId, invalidate],
  )

  const remove = useCallback(
    async (workflowId: string) => {
      if (!organizationId) {
        return null
      }
      const result = await unwrap(
        workflowsService.delete(organizationId, workflowId),
      )
      await invalidate()
      return result
    },
    [organizationId, invalidate],
  )

  const duplicate = useCallback(
    async (workflowId: string, payload: unknown) => {
      if (!organizationId) {
        return null
      }
      const result = await unwrap(
        workflowsService.duplicate(
          organizationId,
          workflowId,
          payload as never,
        ),
      )
      await invalidate()
      return result
    },
    [organizationId, invalidate],
  )

  return { create, update, remove, duplicate }
}
