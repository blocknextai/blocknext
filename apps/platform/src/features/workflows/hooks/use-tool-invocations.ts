import useSWR from 'swr'
import toolInvocationsService from '@/features/workflows/services/tool-invocations'

type ToolInvocationsParams = { query?: string; offset?: number; limit?: number }

export function useToolInvocations(
  organizationId: string | null | undefined,
  params: ToolInvocationsParams = {},
) {
  const { query = '', offset = 0, limit = 10 } = params
  const key = organizationId
    ? ['tool-invocations', organizationId, query, offset, limit]
    : null
  const { data, error, isLoading, mutate } = useSWR(
    key,
    () =>
      toolInvocationsService.getAll(organizationId, { query, offset, limit }),
    { keepPreviousData: true },
  )
  return {
    toolInvocations: data?.data ?? [],
    pagination: data?.meta?.pagination,
    isLoading,
    error,
    mutate,
  }
}

export function useToolInvocation(
  organizationId: string | null | undefined,
  toolInvocationId: string | null | undefined,
) {
  const key =
    organizationId && toolInvocationId
      ? ['tool-invocation', organizationId, toolInvocationId]
      : null
  const { data, error, isLoading } = useSWR(key, () =>
    toolInvocationsService.getById(organizationId, toolInvocationId),
  )
  return {
    toolInvocation: data?.data,
    isLoading,
    error,
  }
}
