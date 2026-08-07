import {
  Bell,
  LogOut,
  User,
  Wallet,
  Copy,
  Check,
  ChevronRight,
  ChevronDown,
  KeyRound,
} from 'lucide-react'
import { useState } from 'react'
import { Avatar } from '@/components/ui/avatar'
import { UserIcon2 } from '@/components/shared/custom-icons'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuShortcut,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { ModeToggle } from '@/features/theme/components/mode-toggle'
import tokenManager from '@/lib/token-manager'
import { authService } from '@/features/auth'
import { useOrganizationStore } from '@/stores/organization'
import { useUserNotificationsCount } from '@/features/notifications'

const NavUser = ({ linkedAccounts, alone }) => {
  const { t } = useTranslation()
  const { unseenCount } = useUserNotificationsCount()
  let ib = false
  let st = 'expanded'

  // Only use useSidebar when not alone and within SidebarProvider
  if (!alone) {
    try {
      const sidebar = useSidebar()
      ib = sidebar.isMobile
      st = sidebar.state
    } catch {
      // If useSidebar fails (not within SidebarProvider), use defaults
      ib = false
      st = 'expanded'
    }
  }

  const primaryAccount = linkedAccounts?.find((a) => a.isPrimary)
  const address = primaryAccount?.identifier || ''
  const displayName = address ? address.slice(0, 10) + '...' : ''

  const doLogOut = async () => {
    try {
      await authService.logout()
    } catch {
      // proceed with local logout even if backend call fails
    }
    tokenManager.clearTokens()
    useOrganizationStore.getState().reset()
    window.location.href = '/auth/login'
  }
  const [copied, setCopied] = useState(false)
  const copyClipboard = () => {
    navigator.clipboard.writeText(address)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const userMenu = (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <div
          className={`flex gap-4 ${st && st === 'collapsed' ? 'mt-4 p-0' : 'gap-4 py-3 px-4 '} items-center cursor-pointer rounded-md 
      data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground
      hover:bg-sidebar-accent hover:text-sidebar-accent-foreground`}
        >
          <div className="relative shrink-0">
            <Avatar
              className={
                st && st === 'collapsed' ? 'h-8 w-8 p-0' : 'h-8 w-8 rounded-lg'
              }
            >
              <UserIcon2
                seed={address}
                className="size-4 text-primary"
                size={30}
              />
            </Avatar>
            {unseenCount > 0 && (
              <span className="absolute -top-1 -right-1 flex items-center justify-center min-w-4 h-4 px-1 rounded-full bg-primary text-primary-foreground text-[10px] font-medium leading-none">
                {unseenCount > 99 ? '99+' : unseenCount}
              </span>
            )}
          </div>
          <div className="grid flex-1 text-left text-sm leading-tight">
            <span className="truncate font-medium">{displayName}</span>
          </div>
          {alone ? (
            <ChevronDown className="ml-auto size-4" />
          ) : (
            <ChevronRight className="ml-auto size-4" />
          )}
        </div>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        side={ib ? 'bottom' : 'right'}
        align="end"
        sideOffset={4}
      >
        <DropdownMenuLabel className="p-0 font-normal">
          <div className="px-3 py-2 text-md font-semibold">
            {' '}
            {t('ui.text.settings')}{' '}
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem
            onSelect={(e) => e.preventDefault()}
            onClick={copyClipboard}
          >
            <Wallet size={16} />
            {displayName}
            <DropdownMenuShortcut>
              {copied ? (
                <Check size={16} className="text-green-500" />
              ) : (
                <Copy size={16} />
              )}
            </DropdownMenuShortcut>
          </DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <Link to="/preferences/profile">
            <DropdownMenuItem>
              <User size={16} />
              {t('ui.text.profile')}
            </DropdownMenuItem>
          </Link>
          <Link to="/preferences/notifications">
            <DropdownMenuItem>
              <Bell size={16} />
              {t('ui.text.notifications')}
              {unseenCount > 0 && (
                <span className="ml-auto flex items-center justify-center min-w-4 h-4 px-1 rounded-full bg-primary text-primary-foreground text-[10px] font-medium leading-none">
                  {unseenCount > 99 ? '99+' : unseenCount}
                </span>
              )}
            </DropdownMenuItem>
          </Link>
          <Link to="/preferences/credentials">
            <DropdownMenuItem>
              <KeyRound size={16} />
              {t('ui.text.credentials')}
            </DropdownMenuItem>
          </Link>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem variant="destructive" onClick={doLogOut}>
          <LogOut size={16} />
          {t('ui.text.logOut')}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )

  if (alone) {
    return userMenu
  }

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <ModeToggle
          className="w-full"
          direction={st && st === 'collapsed' ? 'vertical' : 'horizontal'}
        />
      </SidebarMenuItem>
      <SidebarMenuItem>{userMenu}</SidebarMenuItem>
    </SidebarMenu>
  )
}

NavUser.displayName = 'NavUser'

export { NavUser }
