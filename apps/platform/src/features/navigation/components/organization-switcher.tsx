import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { ChevronsUpDown, Plus, Settings } from 'lucide-react'
import { useThemeStore } from '@/stores/theme-store'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { CreateOrganizationForm } from '@/features/organizations/components/create-organization-form'

import { Link, useParams, useNavigate } from 'react-router'
import { CompanyIcon } from '@/components/shared/custom-icons'
import { useOrganizationStore } from '@/stores/organization'

const OrganizationSwitcher = ({ teams, onCreateOrganization }) => {
  const { t } = useTranslation()
  const { isMobile } = useSidebar()
  const [activeTeam, setActiveTeam] = useState(null)
  const [open, setOpen] = useState(false)
  const [avatarColor, setAvatarColor] = useState('fafafa')
  const mode = useThemeStore((s) => s.mode)
  const { organizationId } = useParams()
  const setOrganizationId = useOrganizationStore((s) => s.setOrganizationId)
  const storedOrgId = useOrganizationStore((s) => s.organizationId)
  const navigate = useNavigate()

  const createOrganization = (result) => {
    setOrganizationId(result.id)
    setOpen(false)
    navigate(`/organizations/${result.id}`)
  }

  useEffect(() => {
    if (teams && teams.length > 0) {
      const currentOrgId = organizationId || storedOrgId
      const selectedTeam = teams.find((t) => t.id === currentOrgId) || teams[0]
      setActiveTeam(selectedTeam)
      setOrganizationId(selectedTeam.id)
    }
  }, [teams, organizationId])

  useEffect(() => {
    const prefersDark = window.matchMedia(
      '(prefers-color-scheme: dark)',
    ).matches
    if (mode === 'dark') {
      setAvatarColor('2a2627')
    } else if (mode === 'light') {
      setAvatarColor('fafafa')
    } else {
      setAvatarColor(prefersDark ? '2a2627' : 'fafafa')
    }
  }, [mode])

  if (!activeTeam || !teams || teams.length === 0) {
    return null
  }

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <div className="relative bg-primary text-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg overflow-hidden">
                <CompanyIcon seed={activeTeam.id} avatarColor={avatarColor} />
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{activeTeam.title}</span>
              </div>
              <ChevronsUpDown className="ml-auto" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="start"
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className="text-muted-foreground text-xs">
              {t('ui.text.organizations')}
            </DropdownMenuLabel>
            {teams?.map((team, index) => (
              <DropdownMenuItem
                onClick={() => {
                  setActiveTeam(team)
                  setOrganizationId(team.id)
                  setTimeout(() => {
                    if (organizationId) {
                      const nPath = window.location.pathname.replace(
                        organizationId,
                        team.id,
                      )
                      navigate(nPath)
                    }
                  }, 250)
                }}
                className={`gap-2 p-2 flex items-center justify-between ${activeTeam?.id === team.id ? 'bg-accent' : ''}`}
                key={index}
              >
                <div className="flex gap-2 items-center">
                  <div className="relative bg-primary text-primary-foreground flex aspect-square size-6 items-center justify-center rounded-lg overflow-hidden">
                    <CompanyIcon seed={team.id} avatarColor={avatarColor} />
                  </div>
                  {team.title}
                </div>

                <Link
                  key={team.title}
                  to={`/organizations/${team.id}/settings`}
                  onClick={(e) => e.stopPropagation()}
                >
                  <DropdownMenuShortcut>
                    <Settings size={16} />
                  </DropdownMenuShortcut>
                </Link>
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="gap-2 p-2"
              onClick={() => setOpen(true)}
            >
              <div className="flex size-6 items-center justify-center rounded-md bg-transparent">
                <Plus className="size-4" />
              </div>
              <div className="text-muted-foreground font-medium">
                {t('ui.text.addOrganization')}
              </div>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>{t('ui.text.organizationDetails')}</DialogTitle>
              <DialogDescription>
                {t('ui.text.editOrganizationDetails')}
              </DialogDescription>
            </DialogHeader>
            <div className="mt-4">
              <CreateOrganizationForm
                onSubmit={onCreateOrganization}
                onCreated={createOrganization}
                onCancel={() => setOpen(false)}
              />
            </div>
          </DialogContent>
        </Dialog>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}

OrganizationSwitcher.displayName = 'OrganizationSwitcher'

export { OrganizationSwitcher }
