import { useCallback, useMemo, useState } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import credentialsService from '@/features/credentials/services/credentials'
import { organizationsService } from '@/features/organizations'
import { nodeEngineService } from '@/features/flow-editor'
import { useIconResolver } from '@/features/flow-editor/icons'
import { unwrap } from '@/lib/swr'

const schemaToFields = (schema: any) => {
  if (!schema || !schema.properties) {
    return []
  }
  const required = new Set<string>(schema.required ?? [])
  return Object.entries(schema.properties).map(
    ([key, value]: [string, any]) => ({
      key,
      type: value.type,
      title: value.title,
      description: value.description,
      default: value.default,
      hidden: !!value.hidden,
      readOnly: !!value.readOnly,
      writeOnly: !!value.writeOnly,
      copyable: !!value.copyable,
      required: required.has(key),
    }),
  )
}

const NODE_ENGINE_KEY = 'node-engine-credentials'

const defaultPagination = {
  total: 0,
  page: 1,
  limit: 10,
  offset: 0,
  hasNext: false,
  hasPrev: false,
}

type Pagination = { offset?: number; limit?: number; query?: string }

export function useCredentials(
  organizationId: string | null,
  isUserMode: boolean,
) {
  const hasAccess = isUserMode || !!organizationId
  const resolveIcon = useIconResolver()

  const [pagination, setPagination] = useState<Pagination>({
    offset: 0,
    limit: 10,
    query: '',
  })

  const { data: nodeEngineData, isLoading: nodeEngineLoading } = useSWR(
    NODE_ENGINE_KEY,
    () => unwrap(nodeEngineService.getCredentials()),
  )

  const savedKey = !hasAccess
    ? null
    : isUserMode
      ? [
          'user-credentials',
          pagination.offset,
          pagination.limit,
          pagination.query ?? '',
        ]
      : [
          'organization-credentials',
          organizationId,
          pagination.offset,
          pagination.limit,
          pagination.query ?? '',
        ]

  const {
    data: savedResponse,
    isLoading: savedLoading,
    mutate: mutateSaved,
  } = useSWR(
    savedKey,
    () =>
      isUserMode
        ? credentialsService.getUserCredentials(pagination)
        : organizationsService.getOrganizationCredentials(
            organizationId,
            pagination,
          ),
    { keepPreviousData: true },
  )

  const ca: any[] = (nodeEngineData as any[]) ?? []
  const sa: any[] = (savedResponse?.data as any[]) ?? []

  const credentialsArray = useMemo(() => {
    const credentialMap: Record<string, any> = {}
    ca.forEach((c) => {
      credentialMap[c.id] = c
    })
    return sa.map((secret: any) => {
      const credential = credentialMap[secret.key]
      return {
        ...secret,
        icon: (credential && resolveIcon(credential.icon)) ?? null,
        isSupportPlatform: credential?.isSupportPlatform || false,
        sourceType: secret.sourceType || 'owner',
      }
    })
  }, [ca, sa, resolveIcon])

  const secrets = useMemo(() => {
    return ca.map((c) => {
      return {
        id: c.id,
        name: c.name,
        description: c.description,
        isOAuth2: c.isOAuth2,
        isSupportPlatform: c.isSupportPlatform,
        provider: schemaToFields(c.schema),
        icon: resolveIcon(c.icon),
      }
    })
  }, [ca, resolveIcon])

  const fetchCredentials = useCallback((offset = 0, limit = 10, query = '') => {
    setPagination({ offset, limit, query })
  }, [])

  const getSecrets = useCallback(async () => {
    await mutateSaved()
  }, [mutateSaved])

  const getSecretDetails = useCallback(
    async (secretId: string) => {
      if (!hasAccess || !secretId) {
        return null
      }
      try {
        const response = isUserMode
          ? await credentialsService.getUserCredentialById(secretId)
          : await organizationsService.getOrganizationCredentialById(
              organizationId,
              secretId,
            )
        return response.data
      } catch (error) {
        console.error('Error fetching secret details:', error)
        return null
      }
    },
    [hasAccess, isUserMode, organizationId],
  )

  return {
    loading: hasAccess ? nodeEngineLoading || savedLoading : nodeEngineLoading,
    secrets,
    credentialsArray,
    pagination: savedResponse?.meta?.pagination ?? defaultPagination,
    hasAccess,
    fetchCredentials,
    getSecrets,
    getSecretDetails,
  }
}

export function useCredentialsRevalidate() {
  const { mutate } = useSWRConfig()
  return useCallback(() => {
    return Promise.all([
      mutate(
        (k) => Array.isArray(k) && k[0] === 'user-credentials',
        undefined,
        { revalidate: true },
      ),
      mutate(
        (k) => Array.isArray(k) && k[0] === 'organization-credentials',
        undefined,
        { revalidate: true },
      ),
    ])
  }, [mutate])
}
