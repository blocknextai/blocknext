import { useTranslation } from 'react-i18next'
import { Info, Wrench } from 'lucide-react'
import TimeAgoI18n from '@/components/shared/timeagoi18'
import StatusBadge from '@/components/shared/status-badge'
import ActionMenu from '@/components/shared/action-menu'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

const ToolInvocationsTable = ({ toolInvocations, onSelect }) => {
  const { t } = useTranslation()

  return (
    <Table className="table-fixed">
      <TableHeader>
        <TableRow>
          <TableHead className="w-auto min-w-[12rem]">
            {t('ui.text.tool')}
          </TableHead>
          <TableHead className="w-32">{t('ui.text.status')}</TableHead>
          <TableHead className="w-28">{t('ui.text.source')}</TableHead>
          <TableHead className="w-44">{t('ui.text.completedAt')}</TableHead>
          <TableHead className="w-44">{t('ui.text.startedAt')}</TableHead>
          <TableHead className="w-20 text-right">
            {t('ui.text.actions')}
          </TableHead>
        </TableRow>
      </TableHeader>

      <TableBody>
        {toolInvocations.map((invocation) => {
          const actionMenuItems = [
            {
              label: t('ui.text.details'),
              icon: <Info />,
              onClick: () => onSelect(invocation.id),
            },
          ]

          return (
            <TableRow key={invocation.id}>
              <TableCell>
                <Button
                  variant="link"
                  size="sm"
                  onClick={() => onSelect(invocation.id)}
                  className="px-0 text-foreground max-w-full"
                >
                  <Wrench className="size-5 text-icon shrink-0" />
                  <span className="truncate">{invocation.toolId}</span>
                </Button>
              </TableCell>
              <TableCell>
                <StatusBadge
                  status={invocation.status}
                  title={invocation.errorMessage}
                />
              </TableCell>
              <TableCell>
                <Badge variant="outline" className="font-normal">
                  {invocation.source}
                </Badge>
              </TableCell>
              <TableCell className="text-muted-foreground">
                <TimeAgoI18n date={invocation.completedAt} />
              </TableCell>
              <TableCell className="text-muted-foreground">
                <TimeAgoI18n date={invocation.startedAt} />
              </TableCell>
              <TableCell className="p-4 w-full h-full flex items-center justify-end">
                <ActionMenu items={actionMenuItems} />
              </TableCell>
            </TableRow>
          )
        })}
      </TableBody>
    </Table>
  )
}

ToolInvocationsTable.displayName = 'ToolInvocationsTable'

export { ToolInvocationsTable }
