import { Badge } from '@/components/ui/badge'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'
import { GhostInput } from '@/components/shared/ghost-input'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Pencil, UserMinus, Check } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { UserIcon2 } from '@/components/shared/custom-icons'

const MemberTable = ({
  members,
  roles,
  roleColorMap,
  roleIcons,
  openPopoverId,
  onOpenPopoverChange,
  onAliasChange,
  onAliasFocus,
  onAliasBlur,
  onUpdateAlias,
  onUpdateRole,
  onRemoveMember,
}) => {
  const { t } = useTranslation()

  const renderMembers = () => {
    const mArray: React.ReactNode[] = []
    for (let i = 0; i < members.length; i++) {
      const m = members[i]
      mArray.push(
        <TableRow key={i}>
          <TableCell>
            <UserIcon2 seed={m.id} className="size-4 text-primary" size={40} />
          </TableCell>
          <TableCell className="font-medium">
            {m.linkedAccounts?.find((a) => a.isPrimary)?.displayName}
          </TableCell>
          <TableCell className="font-medium">
            <div className="flex bg-input/35 w-fit items-center p-2 rounded-md">
              <GhostInput
                defaultValue={m.alias}
                className="peer"
                onFocus={(e) =>
                  onAliasFocus((e.target as HTMLInputElement).value)
                }
                onBlur={() => onAliasBlur(m.alias)}
                placeholder={t('ui.text.optional')}
                onChange={(e) =>
                  onAliasChange(m.id, (e.target as HTMLInputElement).value)
                }
              />
              <Button
                size={'icon'}
                className="opacity-0 peer-focus:opacity-100 cursor-default size-4! peer-focus:cursor-pointer"
                onClick={() => onUpdateAlias(m)}
              >
                <Check className="size-3" strokeWidth={3} />
              </Button>
            </div>
          </TableCell>
          <TableCell>
            {m.role === 'organization:owner' ? (
              <Badge
                className={`${roleColorMap[m.role].badgeBg} ${roleColorMap[m.role].badgeText} flex gap-1 w-fit`}
              >
                {roleIcons[m.role]} {m.role}
              </Badge>
            ) : (
              <DropdownMenu
                key={m.id}
                open={openPopoverId === m.id}
                onOpenChange={(isOpen) => {
                  onOpenPopoverChange(isOpen ? m.id : null)
                }}
              >
                <DropdownMenuTrigger>
                  <Badge
                    className={`cursor-pointer hover:bg-muted hover:text-foreground ${roleColorMap[m.role].badgeBg} ${roleColorMap[m.role].badgeText} flex gap-1 w-fit`}
                  >
                    {roleIcons[m.role]} {m.role} <Pencil size={14} />
                  </Badge>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-auto">
                  <DropdownMenuLabel>{t('ui.text.roles')}</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  {roles.map((r, s) => {
                    return (
                      <DropdownMenuItem
                        key={s}
                        onClick={() => onUpdateRole(m.id, r.name)}
                      >
                        {roleIcons[r.name]} {r.name}
                      </DropdownMenuItem>
                    )
                  })}
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </TableCell>
          <TableCell className="p-4 w-full h-full flex items-center justify-end">
            <AlertDialog key={m.id}>
              <AlertDialogTrigger asChild>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button variant="ghost" className="cursor-pointer">
                      <UserMinus />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent side="top">
                    {t('ui.text.remove')}
                  </TooltipContent>
                </Tooltip>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>{t('ui.text.areYouSure')}</AlertDialogTitle>
                  <AlertDialogDescription>
                    {t('ui.text.removeMemberDescription')}
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>{t('ui.text.cancel')}</AlertDialogCancel>
                  <AlertDialogAction asChild>
                    <Button
                      variant="destructive"
                      onClick={() => onRemoveMember(m.id)}
                    >
                      {t('ui.text.removeMember')}
                    </Button>
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </TableCell>
        </TableRow>,
      )
    }
    return mArray
  }

  return (
    <div className="flex flex-cols gap-4 py-6">
      <Table>
        <TableCaption className="text-muted-foreground/50">
          {t('ui.text.memberListCaption')}
        </TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead></TableHead>
            <TableHead>{t('ui.text.displayName')}</TableHead>
            <TableHead>{t('ui.text.alias')}</TableHead>
            <TableHead>{t('ui.text.role')}</TableHead>
            <TableHead className="text-right">{t('ui.text.actions')}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>{renderMembers()}</TableBody>
      </Table>
    </div>
  )
}

MemberTable.displayName = 'MemberTable'

export { MemberTable }
