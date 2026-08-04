import { useCallback, useEffect } from 'react'
import { Outlet } from 'react-router'
import { AppSidebar } from '@/components/shared/app-sidebar'
import { ErrorBoundary } from '@/components/shared/error-boundary'
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { Footer } from '@/features/navigation/components/footer'
import { WsStatusBadge } from '@/components/shared/ws-status-badge'
import { NotificationsBell } from '@/features/notifications/components/notifications-bell'
import {
  useOrganizations,
  useOrganizationActions,
} from '@/features/organizations'
import { usePlatformFeatures } from '@/features/platform'
import { useOrganizationStore } from '@/stores/organization'
import { useAuthLinkedAccounts } from '@/features/auth'
import wsManager from '@/lib/ws-manager'

const OrganizationLayout = () => {
  const setOrganizations = useOrganizationStore((s) => s.setOrganizations)
  const organizationId = useOrganizationStore((s) => s.organizationId)
  const setOrganizationId = useOrganizationStore((s) => s.setOrganizationId)
  const { organizations } = useOrganizations()
  const { linkedAccounts } = useAuthLinkedAccounts()
  const { create: createOrganization } = useOrganizationActions()
  usePlatformFeatures()

  useEffect(() => {
    if (!organizations) return
    setOrganizations(organizations)
    if (organizationId === null && organizations.length > 0) {
      setOrganizationId(organizations[0].id)
    }
  }, [organizations, organizationId, setOrganizations, setOrganizationId])

  useEffect(() => {
    if (organizationId) {
      wsManager.connect(organizationId)
    }
    return () => {
      wsManager.disconnect()
    }
  }, [organizationId])

  const handleCreateOrganization = useCallback(
    async (data: unknown) => {
      return createOrganization(data)
    },
    [createOrganization],
  )

  return (
    <SidebarProvider>
      <AppSidebar
        onCreateOrganization={handleCreateOrganization}
        linkedAccounts={linkedAccounts}
      />
      <SidebarInset>
        <header className="flex h-16 shrink-0 pr-4 justify-between items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
          <div className="flex  items-center gap-2 px-4">
            <SidebarTrigger className="-ml-1" />
          </div>
          <div className="flex items-center gap-4">
            <NotificationsBell />
            <WsStatusBadge />
          </div>
        </header>
        <div className="flex flex-1 flex-col gap-4 p-0 relative">
          <ErrorBoundary scope="section">
            <Outlet />
          </ErrorBoundary>
          <Footer />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}

export default OrganizationLayout
