import { useTranslation } from 'react-i18next'
import TimeAgoI18n from '@/components/shared/timeagoi18'
import ActionMenu from '@/components/shared/action-menu'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Trash2, Info, Zap, Play, Pause, ClockFading } from 'lucide-react'
import { FlowIcon } from '@/components/shared/custom-icons'
import cronstrue from 'cronstrue/i18n'

const getTriggerTypeIcon = (type: string) => {
  switch (type) {
    case 'schedule':
      return <ClockFading className="size-4" />
    case 'webhook':
      return <Zap className="size-4" />
    case 'event':
      return <Play className="size-4" />
    default:
      return <Zap className="size-4" />
  }
}

const TRIGGER_TYPE_LABELS: Record<string, string> = {
  schedule: 'ui.text.scheduled',
  webhook: 'ui.text.webhook',
  event: 'ui.text.event',
}

const TriggerTable = ({
  filteredTriggers,
  lang,
  onToggleTrigger,
  onConfirmDelete,
  onDetails,
}) => {
  const { t: translate } = useTranslation()

  const renderTriggers = () => {
    const renderArray: React.ReactNode[] = []
    for (let i = 0; i < filteredTriggers?.length; i++) {
      const t = filteredTriggers[i]

      const actionMenuItems: Array<{
        label: string
        icon: React.ReactNode
        onClick?: () => void
        variant?: string
      }> = [
        {
          label: translate('ui.text.details'),
          icon: <Info />,
          onClick: () => onDetails(t),
        },
      ]

      if (t.isActive) {
        actionMenuItems.push({
          label: translate('ui.text.pause'),
          icon: <Pause />,
          onClick: () => onToggleTrigger(t),
        })
      } else {
        actionMenuItems.push({
          label: translate('ui.text.activate'),
          icon: <Play />,
          onClick: () => onToggleTrigger(t),
        })
      }

      actionMenuItems.push({
        label: translate('ui.text.delete'),
        icon: <Trash2 />,
        onClick: () => onConfirmDelete(t.id),
        variant: 'destructive',
      })

      let cronText = ''
      if (t.cronPattern) {
        try {
          cronText = cronstrue.toString(t.cronPattern, { locale: lang })
        } catch {
          cronText = t.cronPattern
        }
      }

      const renderConfiguration = () => {
        if (t.type === 'schedule') {
          return (
            <div className="flex flex-col">
              <span className="text-sm">{cronText || '—'}</span>
              {t.timezone && (
                <span className="text-xs text-muted-foreground">
                  {t.timezone}
                </span>
              )}
            </div>
          )
        }
        return <span className="text-muted-foreground text-sm">—</span>
      }

      renderArray.push(
        <TableRow key={i}>
          <TableCell>
            <Button
              variant="link"
              size="sm"
              onClick={() => onDetails(t)}
              className="px-0 text-foreground"
            >
              <div
                className="size-6 rounded-sm background-cover background-center"
                style={{
                  backgroundImage: `${FlowIcon(t.workflow.id)}`,
                }}
              ></div>
              <div>{t.workflow.title}</div>
            </Button>
          </TableCell>
          <TableCell>
            <div className="flex gap-2 items-center">
              <div className="p-1 rounded-sm bg-muted">
                {getTriggerTypeIcon(t.type)}
              </div>
              <div className="font-medium">
                {TRIGGER_TYPE_LABELS[t.type]
                  ? translate(TRIGGER_TYPE_LABELS[t.type])
                  : t.type}
              </div>
            </div>
          </TableCell>
          <TableCell>{renderConfiguration()}</TableCell>
          <TableCell>
            <span
              className={`text-sm ${t.isActive ? 'text-green-600' : 'text-muted-foreground'}`}
            >
              {t.isActive
                ? translate('ui.text.active')
                : translate('ui.text.inactive')}
            </span>
          </TableCell>
          <TableCell>
            {t.updatedAt ? (
              <TimeAgoI18n date={t.updatedAt} />
            ) : (
              <span className="text-muted-foreground text-sm">—</span>
            )}
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
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>{translate('ui.text.name')}</TableHead>
          <TableHead>{translate('ui.text.type')}</TableHead>
          <TableHead>
            {translate('ui.text.configuration', 'Configuration')}
          </TableHead>
          <TableHead>{translate('ui.text.status')}</TableHead>
          <TableHead>{translate('ui.text.updated')}</TableHead>
          <TableHead className="text-right">
            {translate('ui.text.actions')}
          </TableHead>
        </TableRow>
      </TableHeader>

      <TableBody>{renderTriggers()}</TableBody>
    </Table>
  )
}

TriggerTable.displayName = 'TriggerTable'

export { TriggerTable }
