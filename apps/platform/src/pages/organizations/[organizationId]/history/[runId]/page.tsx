import { useParams, useNavigate } from 'react-router'
import { formatDistanceStrict } from 'date-fns'
import { Loading } from '@/components/shared/loading'
import TimeAgoI18n from '@/components/shared/timeagoi18'
import { useTranslation } from 'react-i18next'
import { useState, useEffect, useCallback } from 'react'
import { executionsService, taskRunnerService } from '@/features/workflows'
import { CheckCircle2, XCircle, Ban, Clock, SkipForward } from 'lucide-react'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { WalletIcon } from '@/components/shared/wallet-icon'
import { RunDetailHeader } from '@/features/organizations/components/history/run-detail-header'
import { RunStepList } from '@/features/organizations/components/history/run-step-list'
import { useOrganizationEvents } from '@/hooks/use-organization-events'
import type { OrganizationEvent } from '@/lib/ws-manager'

const statusPrefs = {
  pending: {
    icon: <Clock className="size-4" />,
    color: 'text-zinc-500',
    bgColor: 'bg-zinc-100 dark:bg-zinc-800',
  },
  running: {
    icon: <span className="loading loading-ring loading-xs"></span>,
    color: 'text-amber-500',
    bgColor: 'bg-amber-100 dark:bg-amber-900/20',
  },
  success: {
    icon: <CheckCircle2 className="size-4" />,
    color: 'text-emerald-500',
    bgColor: 'bg-emerald-100 dark:bg-emerald-900/20',
  },
  failed: {
    icon: <XCircle className="size-4" />,
    color: 'text-destructive',
    bgColor: 'bg-destructive/10',
  },
  cancelled: {
    icon: <Ban className="size-4" />,
    color: 'text-zinc-500',
    bgColor: 'bg-zinc-100 dark:bg-zinc-800',
  },
  skipped: {
    icon: <SkipForward className="size-4" />,
    color: 'text-zinc-400',
    bgColor: 'bg-zinc-100 dark:bg-zinc-800',
  },
}

function OrganizationHistoryDetailPage() {
  const { organizationId, runId } = useParams()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [history, setHistory] = useState()
  const [open, setOpen] = useState(false)
  const [openT, setOpenT] = useState(false)
  const [expandedNodes, setExpandedNodes] = useState(new Set())
  const { t } = useTranslation()

  const getHistory = async () => {
    const data = await executionsService.getById(organizationId, runId)
    setHistory(data.data)
    setLoading(false)
  }

  const deleteExecution = async () => {
    await executionsService.delete(organizationId, runId)
    getHistory()
  }

  const cancelExecution = async () => {
    await taskRunnerService.cancel(organizationId, {
      id: runId,
    })
    getHistory()
  }

  const rerunAllExecution = async () => {
    await taskRunnerService.rerunAll(organizationId, {
      id: runId,
    })
    navigate(`/organizations/${organizationId}/history`)
  }

  const rerunFailedExecution = async () => {
    await taskRunnerService.rerunFailed(organizationId, {
      id: runId,
    })
    navigate(`/organizations/${organizationId}/history`)
  }

  const toggleNodeExpansion = (nodeId) => {
    const newExpanded = new Set(expandedNodes)
    if (newExpanded.has(nodeId)) {
      newExpanded.delete(nodeId)
    } else {
      newExpanded.add(nodeId)
    }
    setExpandedNodes(newExpanded)
  }

  const renderDetails = () => {
    if (!history) return null

    const renderDate = (dateString) => {
      if (!dateString) return t('ui.text.workInProgress')
      const date = new Date(dateString)
      if (isNaN(date.getTime())) return t('ui.text.invalidDate')
      return <TimeAgoI18n date={dateString} />
    }

    const calculateDuration = () => {
      try {
        if (!history.startedAt || !history.completedAt)
          return t('ui.text.workInProgress')
        const startDate = new Date(history.startedAt)
        const endDate = new Date(history.completedAt)
        if (isNaN(startDate.getTime()) || isNaN(endDate.getTime()))
          return t('ui.text.invalid')
        return formatDistanceStrict(startDate, endDate)
      } catch (error) {
        return t('ui.text.error')
      }
    }

    return (
      <div className="flex flex-col gap-2 mb-6">
        {history.errorMessage && (
          <div className="flex gap-2 items-center p-3 bg-destructive/10 rounded-lg">
            <XCircle className="size-4 text-destructive" />
            <span className="text-destructive">{t(history.errorMessage)}</span>
          </div>
        )}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="flex gap-2 items-center">
            <Clock className="size-4 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              {t('ui.text.started')}:
            </span>
            <span className="text-sm">{renderDate(history.startedAt)}</span>
          </div>
          <div className="flex gap-2 items-center">
            <Clock className="size-4 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              {t('ui.text.finished')}:
            </span>
            <span className="text-sm">{renderDate(history.completedAt)}</span>
          </div>
          <div className="flex gap-2 items-center">
            <span className="text-sm text-muted-foreground">
              {t('ui.text.runBy')}:
            </span>
            <WalletIcon owner={history.triggeredByUser} />
          </div>
          <div className="flex gap-2 items-center">
            <span className="text-sm text-muted-foreground">
              {t('ui.text.duration')}:
            </span>
            <span className="text-sm">{calculateDuration()}</span>
          </div>
        </div>
      </div>
    )
  }

  const handleEvent = useCallback(
    (event: OrganizationEvent) => {
      if (event.id === runId) {
        getHistory()
      }
    },
    [runId],
  )

  useOrganizationEvents(handleEvent)

  useEffect(() => {
    getHistory()
  }, [])

  if (loading) {
    return <Loading />
  }

  return (
    <div className="w-full h-full flex p-6 flex-col overflow-hidden overflow-x-hidden">
      {history && (
        <RunDetailHeader
          history={history}
          statusPrefs={statusPrefs}
          onRerunAll={rerunAllExecution}
          onRerunFailed={rerunFailedExecution}
          onCancelConfirm={() => setOpenT(true)}
          onDeleteConfirm={() => setOpen(true)}
        />
      )}

      {renderDetails()}

      {history && (
        <RunStepList
          nodeExecutions={history.nodeExecutions}
          statusPrefs={statusPrefs}
          expandedNodes={expandedNodes}
          onToggleNode={toggleNodeExpansion}
        />
      )}

      <div className="mt-6">
        <AlertDialog open={open} onOpenChange={setOpen}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>{t('ui.text.areYouSure')}</AlertDialogTitle>
              <AlertDialogDescription>
                {t('ui.text.deleteHistoryConfirmation')}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>{t('ui.text.cancel')}</AlertDialogCancel>
              <AlertDialogAction
                onClick={deleteExecution}
                variant="destructive"
              >
                {t('ui.text.delete')}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
        <AlertDialog open={openT} onOpenChange={setOpenT}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>{t('ui.text.areYouSure')}</AlertDialogTitle>
              <AlertDialogDescription>
                {t('ui.text.cancelExecutionConfirmation')}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>{t('ui.text.cancel')}</AlertDialogCancel>
              <AlertDialogAction
                onClick={cancelExecution}
                variant="destructive"
              >
                {t('ui.text.terminate')}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </div>
  )
}
export default OrganizationHistoryDetailPage
