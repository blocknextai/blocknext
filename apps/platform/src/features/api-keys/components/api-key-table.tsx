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
import { Key, Info, RefreshCw, Trash2 } from 'lucide-react'

const ApiKeyTable = ({
  apiKeys,
  onDetails,
  onConfirmRegenerate,
  onConfirmDelete,
}) => {
  const { t } = useTranslation()

  return (
    <div className="px-6 pb-6">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{t('ui.text.name')}</TableHead>
            <TableHead>{t('ui.text.lastUsed')}</TableHead>
            <TableHead>{t('ui.text.updated')}</TableHead>
            <TableHead className="text-right">{t('ui.text.actions')}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {apiKeys?.map((k) => {
            const actionMenuItems = [
              {
                label: t('ui.text.details'),
                icon: <Info />,
                onClick: () => onDetails(k),
              },
              {
                label: t('ui.text.revokeAndRegenerate', 'Revoke & Regenerate'),
                icon: <RefreshCw />,
                onClick: () => onConfirmRegenerate(k),
              },
              {
                label: t('ui.text.delete'),
                icon: <Trash2 />,
                onClick: () => onConfirmDelete(k),
                variant: 'destructive',
              },
            ]

            return (
              <TableRow key={k.id}>
                <TableCell>
                  <Button
                    variant="link"
                    size="sm"
                    onClick={() => onDetails(k)}
                    className="px-0 text-foreground"
                  >
                    <Key className="size-5 text-icon" /> {k.name}
                  </Button>
                </TableCell>
                <TableCell>
                  {k.lastUsedAt ? (
                    <TimeAgoI18n date={k.lastUsedAt} />
                  ) : (
                    <span className="text-muted-foreground text-sm">
                      {t('ui.text.neverUsed')}
                    </span>
                  )}
                </TableCell>
                <TableCell>
                  {k.updatedAt ? (
                    <TimeAgoI18n date={k.updatedAt} />
                  ) : (
                    <span className="text-muted-foreground text-sm">—</span>
                  )}
                </TableCell>
                <TableCell className="p-4 w-full h-full flex items-center justify-end">
                  <ActionMenu items={actionMenuItems} />
                </TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>
    </div>
  )
}

ApiKeyTable.displayName = 'ApiKeyTable'

export { ApiKeyTable }
