import { Separator } from '@/components/ui/separator'
import { formatDistanceStrict } from 'date-fns'
import ActionMenu from '@/components/shared/action-menu'
import { useTranslation } from 'react-i18next'
import { Trash2, RotateCw, XCircle, Clock, PlayCircle } from 'lucide-react'

const RunDetailHeader = ({
  history,
  statusPrefs,
  onRerunAll,
  onRerunFailed,
  onCancelConfirm,
  onDeleteConfirm,
}) => {
  const { t } = useTranslation()

  let duration = t('ui.text.workInProgress')
  try {
    if (history.startedAt && history.completedAt) {
      const startDate = new Date(history.startedAt)
      const endDate = new Date(history.completedAt)

      if (!isNaN(startDate.getTime()) && !isNaN(endDate.getTime())) {
        const diff = endDate.getTime() - startDate.getTime()
        duration = formatDistanceStrict(0, diff)
      }
    }
  } catch (error) {
    console.warn('Error calculating duration:', error)
  }

  const actionMenuItems: Array<{
    label: string
    icon: React.ReactNode
    onClick?: () => void
    variant?: string
  }> = []

  if (history.status !== 'running' && history.status !== 'pending') {
    actionMenuItems.push({
      label: t('ui.text.reRunAll'),
      icon: <RotateCw />,
      onClick: onRerunAll,
    })
  }
  if (history.status === 'failed') {
    actionMenuItems.push({
      label: t('ui.text.reRunFailed'),
      icon: <PlayCircle />,
      onClick: onRerunFailed,
    })
  }
  if (history.status === 'running' || history.status === 'pending') {
    actionMenuItems.push({
      label: t('ui.text.cancel'),
      icon: <XCircle />,
      onClick: onCancelConfirm,
    })
  }
  if (
    history.status === 'success' ||
    history.status === 'cancelled' ||
    history.status === 'failed'
  ) {
    actionMenuItems.push({
      label: t('ui.text.delete'),
      icon: <Trash2 />,
      onClick: onDeleteConfirm,
      variant: 'destructive',
    })
  }

  return (
    <div className="flex flex-col gap-2 mb-6">
      <div className="flex gap-2 justify-between w-full items-center">
        <div className="flex gap-2 flex-col items-start">
          <span className="heading-md">
            {history.workflow?.title || t('ui.text.untitledWorkflow')}
          </span>
          <span className="text-muted-foreground text-sm">
            {history.workflow?.description || ''}
          </span>
        </div>
        <div className="flex gap-4 items-center">
          <div
            title={history.errorMessage}
            className={`${statusPrefs[history.status]?.color} font-medium flex gap-1 items-center capitalize`}
          >
            {statusPrefs[history.status]?.icon}{' '}
            {t(`ui.text.status.${history.status}`)}
          </div>
          <div className="text-xs text-muted-foreground flex gap-1 items-center">
            <Clock className="size-3" /> {duration}
          </div>
          <ActionMenu items={actionMenuItems} />
        </div>
      </div>
      <Separator className="my-2 p-0" decorative={true} />
    </div>
  )
}

RunDetailHeader.displayName = 'RunDetailHeader'

export { RunDetailHeader }
