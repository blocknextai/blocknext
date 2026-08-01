import { useEffect } from 'react'
import { Navigate, Outlet } from 'react-router'
import { PageLoading } from '@/components/shared/loading'
import { useOrganizations } from '@/features/organizations'
import { useOrganizationStore } from '@/stores/organization'

const RequireOrganization = () => {
  const { organizations, isLoading, error } = useOrganizations()
  const organizationId = useOrganizationStore((s) => s.organizationId)
  const reset = useOrganizationStore((s) => s.reset)

  const isEmpty = !isLoading && !error && organizations.length === 0

  useEffect(() => {
    // drop the persisted selection when the user no longer belongs to any organization
    if (isEmpty && organizationId) {
      reset()
    }
  }, [isEmpty, organizationId, reset])

  if (isLoading) {
    return <PageLoading />
  }

  if (isEmpty) {
    return <Navigate to="/organizations/new" replace />
  }

  return <Outlet />
}

export default RequireOrganization
