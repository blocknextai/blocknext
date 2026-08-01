import { Link } from 'react-router'
import TimeAgoI18n from '@/components/shared/timeagoi18'
import ActionMenu from '@/components/shared/action-menu'
import StatusBadge from '@/components/shared/status-badge'
import { useTranslation } from 'react-i18next'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Trash2,
  RotateCw,
  Info,
  XCircle,
  PlayCircle,
  MousePointerClick,
  Webhook,
  ClockFading,
  Code,
} from 'lucide-react'
import { FlowIcon } from '@/components/shared/custom-icons'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button.tsx'

const executionTypeConfig = {
  manual: { icon: MousePointerClick, labelKey: 'ui.text.manual' },
  schedule: { icon: ClockFading, labelKey: 'ui.text.schedule' },
  webhook: { icon: Webhook, labelKey: 'ui.text.webhook' },
  api: { icon: Code, labelKey: 'ui.text.api' },
}

const HistoryTable = ({
  organizationId,
  filteredHistory,
  onRerunAll,
  onRerunFailed,
  onConfirmCancel,
  onConfirmDelete,
}) => {
  const { t } = useTranslation()

  const renderHistory = () => {
    const renderArray: React.ReactNode[] = []
    for (let i = 0; i < filteredHistory?.length; i++) {
      const h = filteredHistory[i]
      const href = `/organizations/${organizationId}/history/${h.id}`
      const actionMenuItems: Array<{
        label: string
        icon: React.ReactNode
        href?: string
        onClick?: () => void
        variant?: string
      }> = [
        {
          label: t('ui.text.details'),
          icon: <Info />,
          href: href,
        },
      ]
      if (h.status !== 'running' && h.status !== 'pending') {
        actionMenuItems.push({
          label: t('ui.text.reRunAll'),
          icon: <RotateCw />,
          onClick: () => onRerunAll(h),
        })
      }
      if (h.status === 'failed') {
        actionMenuItems.push({
          label: t('ui.text.reRunFailed'),
          icon: <PlayCircle />,
          onClick: () => onRerunFailed(h),
        })
      }
      if (h.status === 'running' || h.status === 'pending') {
        actionMenuItems.push({
          label: t('ui.text.cancel'),
          icon: <XCircle />,
          onClick: () => onConfirmCancel(h.id),
        })
      }
      if (
        h.status === 'success' ||
        h.status === 'cancelled' ||
        h.status === 'failed'
      ) {
        actionMenuItems.push({
          label: t('ui.text.delete'),
          icon: <Trash2 />,
          onClick: () => onConfirmDelete(h.id),
          variant: 'destructive',
        })
      }

      renderArray.push(
        <TableRow key={i}>
          <TableCell>
            <Button
              asChild
              variant="link"
              className="px-0 text-foreground max-w-full"
            >
              <Link to={href}>
                <div className="flex items-center gap-1 min-w-0">
                  <div
                    className="size-6 shrink-0 rounded-sm background-cover background-center"
                    style={{
                      backgroundImage: `${FlowIcon(h.workflow.id)}`,
                    }}
                  ></div>
                  <div className="truncate">{h.workflow.title}</div>
                </div>
              </Link>
            </Button>
          </TableCell>
          <TableCell>
            <StatusBadge status={h.status} title={h.errorMessage} />
          </TableCell>
          <TableCell>
            {(() => {
              const config =
                executionTypeConfig[h.executionType] ||
                executionTypeConfig.manual
              const Icon = config.icon
              return (
                <Badge variant="outline" className="gap-1.5 font-normal">
                  <Icon className="size-3" />
                  {t(config.labelKey)}
                </Badge>
              )
            })()}
          </TableCell>
          <TableCell>
            {h.status === 'success' ||
            h.status === 'cancelled' ||
            h.status === 'failed' ? (
              <TimeAgoI18n date={h.completedAt} />
            ) : (
              '-'
            )}
          </TableCell>
          <TableCell>
            <TimeAgoI18n date={h.startedAt} />
          </TableCell>
          <TableCell className="p-4 w-full h-full flex items-center justify-end">
            <ActionMenu items={actionMenuItems} />
          </TableCell>
        </TableRow>,
      )
    }
    return renderArray
  }

  return (
    <Table className="table-fixed">
      <TableHeader>
        <TableRow>
          <TableHead className="w-auto min-w-[12rem]">
            {t('ui.text.flow')}
          </TableHead>
          <TableHead className="w-32">{t('ui.text.status')}</TableHead>
          <TableHead className="w-32">{t('ui.text.trigger')}</TableHead>
          <TableHead className="w-44">{t('ui.text.completedAt')}</TableHead>
          <TableHead className="w-44">{t('ui.text.startedAt')}</TableHead>
          <TableHead className="w-20 text-right">
            {t('ui.text.actions')}
          </TableHead>
        </TableRow>
      </TableHeader>

      <TableBody>{renderHistory()}</TableBody>
    </Table>
  )
}

HistoryTable.displayName = 'HistoryTable'

export { HistoryTable }
