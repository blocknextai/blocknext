import { User, Palette, ShieldCheck, KeyRound, Bell } from 'lucide-react'
import { Outlet, useLocation, Link } from 'react-router'
import { Card, CardContent } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { useTranslation } from 'react-i18next'

type Tab = {
  id: string
  label: string
  icon: React.ComponentType<{ className?: string }>
  path: string
  group: 'account' | 'workspace'
}

const tabs: Tab[] = [
  {
    id: 'profile',
    label: 'ui.text.profile',
    icon: User,
    path: '/preferences/profile',
    group: 'account',
  },
  {
    id: 'notifications',
    label: 'ui.text.notifications',
    icon: Bell,
    path: '/preferences/notifications',
    group: 'account',
  },
  {
    id: 'appearance',
    label: 'ui.text.appearance',
    icon: Palette,
    path: '/preferences/appearance',
    group: 'account',
  },
  {
    id: 'sessions',
    label: 'ui.text.sessions',
    icon: ShieldCheck,
    path: '/preferences/sessions',
    group: 'account',
  },
  {
    id: 'credentials',
    label: 'ui.text.credentials',
    icon: KeyRound,
    path: '/preferences/credentials',
    group: 'workspace',
  },
]

const PreferencesLayout = () => {
  const location = useLocation()
  const { t } = useTranslation()

  const getActiveTab = () => {
    const path = location.pathname
    const match = tabs.find((tab) => path === tab.path)
    return match?.id || 'profile'
  }

  const activeTab = getActiveTab()

  const renderTab = (tab: Tab) => {
    const Icon = tab.icon
    return (
      <Link
        key={tab.id}
        to={tab.path}
        className={`w-full flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-md transition-colors ${
          activeTab === tab.id
            ? 'bg-primary text-primary-foreground'
            : 'text-muted-foreground hover:text-foreground hover:bg-muted'
        }`}
      >
        <Icon className="w-4 h-4" />
        {t(tab.label)}
      </Link>
    )
  }

  const accountTabs = tabs.filter((tab) => tab.group === 'account')
  const workspaceTabs = tabs.filter((tab) => tab.group === 'workspace')

  return (
    <div className="p-6">
      <div className="flex gap-8">
        <div className="w-64 shrink-0">
          <Card className="border-0 bg-card p-0">
            <nav className="p-3 flex flex-col gap-3">
              {accountTabs.map(renderTab)}
              <Separator className="my-1" />
              {workspaceTabs.map(renderTab)}
            </nav>
          </Card>
        </div>

        <div className="flex-1">
          <Card className="border-0 bg-card p-0!">
            <CardContent className="p-3">
              <Outlet />
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}

export default PreferencesLayout
