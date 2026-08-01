import TimeAgoI18n from '@/components/shared/timeagoi18'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { KeyRound, Trash2, Pencil, Building2, User } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { useTranslation } from 'react-i18next'

const CredentialTable = ({ credentialsArray, onEdit, onDelete }) => {
  const { t } = useTranslation()

  const renderCredentials = () => {
    const renderArray: React.ReactNode[] = []
    for (let i = 0; i < credentialsArray?.length; i++) {
      const c = credentialsArray[i]
      const IconComponent = c.icon || KeyRound
      const isPlatform = c.sourceType === 'platform'
      renderArray.push(
        <TableRow key={i}>
          <TableCell>
            <Button
              variant="link"
              onClick={() => onEdit(c)}
              className="h-auto p-0 gap-2 text-foreground hover:text-primary"
            >
              <IconComponent className="size-5 text-icon" /> {c.title}
            </Button>
          </TableCell>
          <TableCell>
            <Badge
              variant={isPlatform ? 'secondary' : 'outline'}
              className={`gap-1 ${isPlatform ? 'bg-primary/10 text-primary border-primary/20' : ''}`}
            >
              {isPlatform ? (
                <Building2 className="size-3" />
              ) : (
                <User className="size-3" />
              )}
              {isPlatform
                ? t('ui.text.credentialSourceTypePlatform')
                : t('ui.text.credentialSourceTypeOwner')}
            </Badge>
          </TableCell>
          <TableCell>
            <TimeAgoI18n date={c.updatedAt} />
          </TableCell>
          <TableCell className="text-right">
            <div className="flex items-center justify-end gap-1">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="size-8"
                    onClick={() => onEdit(c)}
                  >
                    <Pencil className="size-4" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{t('ui.text.edit')}</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="size-8 text-destructive hover:text-destructive"
                    onClick={() => onDelete(c)}
                  >
                    <Trash2 className="size-4" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{t('ui.text.delete')}</TooltipContent>
              </Tooltip>
            </div>
          </TableCell>
        </TableRow>,
      )
    }
    return renderArray
  }

  return (
    <div className="px-6 pb-6">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{t('ui.text.name')}</TableHead>
            <TableHead>{t('ui.text.credentialSourceType')}</TableHead>
            <TableHead>{t('ui.text.lastUpdated')}</TableHead>
            <TableHead></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>{renderCredentials()}</TableBody>
      </Table>
    </div>
  )
}

CredentialTable.displayName = 'CredentialTable'

export { CredentialTable }
