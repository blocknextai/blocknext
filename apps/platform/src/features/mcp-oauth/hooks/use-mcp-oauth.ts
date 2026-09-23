import { useCallback, useState } from 'react'
import useSWR from 'swr'
import mcpOAuthService from '@/features/mcp-oauth/services/mcp-oauth'
import type { McpAuthorizationRedirect } from '@/features/mcp-oauth/services/mcp-oauth'
import { useOrganizations } from '@/features/organizations'
import { useOrganizationStore } from '@/stores/organization'
import { unwrap } from '@/lib/swr'
import toast from '@/lib/toast'

type Organization = {
  id: string
  title: string
}

export function useMcpAuthorizationRequest(
  authorizationRequestId: string | null,
) {
  const key = authorizationRequestId
    ? ['mcp-authorization-request', authorizationRequestId]
    : null

  const { data, error, isLoading, mutate } = useSWR(key, () =>
    unwrap(
      mcpOAuthService.getAuthorizationRequest(authorizationRequestId as string),
    ),
  )

  return {
    authorizationRequest: data,
    isLoading,
    error,
    mutate,
  }
}

export function useMcpAuthorizationConsent(
  authorizationRequestId: string | null,
) {
  const { authorizationRequest, isLoading, error } = useMcpAuthorizationRequest(
    authorizationRequestId,
  )
  const { organizations, isLoading: isLoadingOrganizations } =
    useOrganizations()
  const activeOrganizationId = useOrganizationStore((s) => s.organizationId)

  const [selectedOrganizationId, setSelectedOrganizationId] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const availableOrganizations = organizations as Organization[]
  const organizationId =
    selectedOrganizationId ||
    availableOrganizations.find(
      (organization) => organization.id === activeOrganizationId,
    )?.id ||
    availableOrganizations[0]?.id ||
    ''

  const submit = useCallback(
    async (request: () => Promise<McpAuthorizationRedirect>) => {
      setIsSubmitting(true)
      try {
        const result = await request()
        window.location.assign(result.redirectUri)
      } catch (submitError) {
        toast.error((submitError as Error).message)
        setIsSubmitting(false)
      }
    },
    [],
  )

  const approve = useCallback(
    () =>
      submit(() =>
        unwrap(
          mcpOAuthService.approveAuthorizationRequest(
            authorizationRequestId as string,
            { organizationId },
          ),
        ),
      ),
    [submit, authorizationRequestId, organizationId],
  )

  const deny = useCallback(
    () =>
      submit(() =>
        unwrap(
          mcpOAuthService.denyAuthorizationRequest(
            authorizationRequestId as string,
          ),
        ),
      ),
    [submit, authorizationRequestId],
  )

  return {
    authorizationRequest,
    organizations: availableOrganizations,
    organizationId,
    setOrganizationId: setSelectedOrganizationId,
    isLoading: isLoading || isLoadingOrganizations,
    isUnavailable:
      !authorizationRequestId || !!error || !authorizationRequest?.isPending,
    isSubmitting,
    approve,
    deny,
  }
}
