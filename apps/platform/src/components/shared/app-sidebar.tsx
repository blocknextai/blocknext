import { useTranslation } from 'react-i18next'
import { useOrganizationStore } from '@/stores/organization'
import { NavMain } from '@/features/navigation/components/main-nav'
import { NavUser } from '@/features/navigation/components/user-nav'
import { OrganizationSwitcher } from '@/features/navigation/components/organization-switcher'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
  useSidebar,
} from '@/components/ui/sidebar'
import { Button } from '@/components/ui/button'
import { Link } from 'react-router'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'
import {
  Building,
  DatabaseZap,
  Workflow,
  Zap,
  Settings,
  SquareTerminal,
  Plus,
  KeyRound,
} from 'lucide-react'

interface AppSidebarProps {
  onCreateOrganization: (data: unknown) => Promise<unknown> | unknown
  linkedAccounts: unknown[]
  [key: string]: unknown
}

export function AppSidebar({
  onCreateOrganization,
  linkedAccounts,
  ...props
}: AppSidebarProps) {
  const { t } = useTranslation()
  const organizations = useOrganizationStore((s) => s.organizations)
  const organizationId = useOrganizationStore((s) => s.organizationId)
  const { state } = useSidebar()

  const sidebarData = [
    {
      isOrganization: false,
      items: [
        {
          title: t('ui.text.myOrganizations'),
          url: '/organizations',
          icon: Building,
        },
        {
          title: t('ui.text.mcpServers'),
          url: '/mcp',
          icon: DatabaseZap,
        },
      ],
    },
    {
      title: t('ui.text.organization'),
      isOrganization: true,
      items: [
        {
          title: t('ui.text.flows'),
          url: '/organizations/:organizationId',
          icon: Workflow,
        },
        {
          title: t('ui.text.triggers'),
          url: '/organizations/:organizationId/triggers',
          icon: Zap,
        },
        {
          title: t('ui.text.history'),
          url: '/organizations/:organizationId/history',
          icon: SquareTerminal,
        },
        {
          title: t('ui.text.credentials'),
          url: '/organizations/:organizationId/credentials',
          icon: KeyRound,
        },
        {
          title: t('ui.text.settings'),
          icon: Settings,
          items: [
            {
              title: t('ui.text.general'),
              url: '/organizations/:organizationId/settings',
            },
            {
              title: t('ui.text.members'),
              url: '/organizations/:organizationId/settings/members',
            },
            {
              title: t('ui.text.apiKeys'),
              url: '/organizations/:organizationId/api-keys',
            },
          ],
        },
      ],
    },
  ]

  return (
    <Sidebar className="border-none" collapsible="icon" {...props}>
      <SidebarHeader>
        <OrganizationSwitcher
          teams={organizations || []}
          onCreateOrganization={onCreateOrganization}
        />
      </SidebarHeader>

      <SidebarContent>
        <Tooltip>
          <TooltipTrigger asChild>
            <Link
              to={`/organizations/${organizationId}/create`}
              className={
                state === 'collapsed'
                  ? 'p-0 flex items-center justify-center'
                  : 'p-4'
              }
            >
              <Button className={state === 'collapsed' ? 'size-8' : 'w-full'}>
                <div className={`flex items-center gap-2`}>
                  <Plus />
                  {state !== 'collapsed' && (
                    <span>{t('ui.text.createFlow')}</span>
                  )}
                </div>
              </Button>
            </Link>
          </TooltipTrigger>
          <TooltipContent
            side="right"
            align="center"
            hidden={state !== 'collapsed'}
          >
            {t('ui.text.createFlow')}
          </TooltipContent>
        </Tooltip>

        {sidebarData.map((category, index) => (
          <NavMain
            key={category.title || index}
            title={category.title}
            items={category.items}
          />
        ))}
      </SidebarContent>
      <SidebarFooter>
        <NavUser linkedAccounts={linkedAccounts} alone={false} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
