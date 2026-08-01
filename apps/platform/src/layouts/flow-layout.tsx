import { useCallback, useEffect, useState } from 'react'
import { Link, Outlet, useParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { Home, Workflow } from 'lucide-react'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import { ErrorBoundary } from '@/components/shared/error-boundary'
import { ModeToggle } from '@/features/theme/components/mode-toggle'
import { NavUser } from '@/features/navigation/components/user-nav'
import { authService } from '@/features/auth'
import { useBreadcrumbStore } from '@/features/flow-editor/stores/breadcrumb-store'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { workflowsService } from '@/features/workflows'

const FlowLayout = () => {
  const { t } = useTranslation()
  const { organizationId, id: flowId } = useParams()

  const { flowName, setFlowName } = useBreadcrumbStore()

  const [open, setOpen] = useState(false)
  const [tempFlowName, setTempFlowName] = useState(flowName)

  const [linkedAccounts, setLinkedAccounts] = useState([])

  useEffect(() => {
    const fetchAccounts = async () => {
      const res = await authService.getLinkedAccounts()
      setLinkedAccounts(res.data)
    }
    fetchAccounts()
  }, [])

  useEffect(() => {
    setTempFlowName(flowName || '')
  }, [flowName])

  const handleSave = useCallback(async () => {
    setFlowName(tempFlowName)
    if (organizationId && flowId) {
      await workflowsService.update(organizationId, flowId, {
        title: tempFlowName,
      })
    }
    setOpen(false)
  }, [setFlowName, tempFlowName, organizationId, flowId])

  return (
    <div className="flex flex-col h-app-screen">
      <header className="flex h-16 shrink-0 items-center justify-between gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
        <div className="flex items-center gap-2 px-4">
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem className="hidden md:block">
                <BreadcrumbLink asChild>
                  <Link
                    to={
                      organizationId ? `/organizations/${organizationId}` : '/'
                    }
                  >
                    <Button variant="ghost" className="cursor-pointer">
                      <Home /> {t('ui.text.home')}
                    </Button>
                  </Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator className="hidden md:block" />
              <BreadcrumbItem>
                <BreadcrumbPage>
                  <Button variant="ghost" onClick={() => setOpen(true)}>
                    <Workflow />
                    {flowName}
                  </Button>
                </BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </div>
        <div className="flex items-center gap-2 px-4">
          <ModeToggle />
          <NavUser linkedAccounts={linkedAccounts} alone={true} />
        </div>
      </header>
      <div className="flex-1 min-h-0 gap-4 p-4 pt-0">
        <ErrorBoundary scope="section">
          <Outlet />
        </ErrorBoundary>
      </div>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{t('ui.text.nameFlow')}</DialogTitle>
            <DialogDescription>
              {t('ui.text.changeFlowTitle')}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4">
            <div className="grid gap-3">
              <Label htmlFor={`f-title`}>{t('ui.text.title')}</Label>
              <Input
                id={`f-title`}
                name={`f-title`}
                defaultValue={tempFlowName}
                onChange={(e) => setTempFlowName(e.target.value)}
              />
            </div>
          </div>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline" size="sm">
                {t('ui.text.cancel')}
              </Button>
            </DialogClose>
            <DialogClose asChild>
              <Button size="sm" onClick={handleSave}>
                {t('ui.text.save')}
              </Button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default FlowLayout
