import { useCallback, useState } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import apiKeysService from '@/features/api-keys/services/api-keys'

const defaultPagination = {
  total: 0,
  page: 1,
  limit: 10,
  offset: 0,
  hasNext: false,
  hasPrev: false,
}

type Pagination = { offset?: number; limit?: number; query?: string }

export function useOrganizationApiKeys(
  organizationId: string | null | undefined,
) {
  const [pagination, setPagination] = useState<Pagination>({
    offset: 0,
    limit: 10,
    query: '',
  })
  const { mutate } = useSWRConfig()

  const key = organizationId
    ? [
        'org-api-keys',
        organizationId,
        pagination.offset,
        pagination.limit,
        pagination.query ?? '',
      ]
    : null

  const { data, error, isLoading } = useSWR(
    key,
    () => apiKeysService.getOrganizationApiKeys(organizationId, pagination),
    { keepPreviousData: true },
  )

  const revalidate = useCallback(
    () =>
      mutate(
        (k) =>
          Array.isArray(k) &&
          k[0] === 'org-api-keys' &&
          k[1] === organizationId,
        undefined,
        { revalidate: true },
      ),
    [organizationId, mutate],
  )

  const fetchKeys = useCallback((offset = 0, limit = 10, query = '') => {
    setPagination({ offset, limit, query })
  }, [])

  const handleCreate = useCallback(
    async (payload: unknown) => {
      if (!organizationId) return null
      const res = await apiKeysService.createOrganizationApiKey(
        organizationId,
        payload as never,
      )
      if (res.isSuccess) {
        await revalidate()
        return res.data
      }
      return null
    },
    [organizationId, revalidate],
  )

  const handleRegenerate = useCallback(
    async (apiKeyId: string, payload: unknown) => {
      if (!organizationId) return null
      const res = await apiKeysService.regenerateOrganizationApiKey(
        organizationId,
        apiKeyId,
        payload as never,
      )
      if (res.isSuccess) {
        await revalidate()
        return res.data
      }
      return null
    },
    [organizationId, revalidate],
  )

  const handleDelete = useCallback(
    async (apiKeyId: string) => {
      if (!organizationId) return
      await apiKeysService.deleteOrganizationApiKey(organizationId, apiKeyId)
      await revalidate()
    },
    [organizationId, revalidate],
  )

  return {
    loading: isLoading,
    apiKeys: data?.data ?? [],
    pagination: data?.meta?.pagination ?? defaultPagination,
    error,
    fetchKeys,
    handleCreate,
    handleRegenerate,
    handleDelete,
  }
}

export function useApiKeyScopes() {
  const { data, error, isLoading } = useSWR('api-key-scopes', () =>
    apiKeysService.getApiKeyScopes(),
  )

  return {
    scopes: data?.data ?? [],
    loadingScopes: isLoading,
    error,
  }
}
