import { useRef, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { MoreVertical, Pencil, Trash2, Check, X, Loader2 } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
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

const SessionItemActions = ({
  session,
  isEditing,
  editingTitle,
  actionLoading,
  onEditingTitleChange,
  onStartRename,
  onCancelRename,
  onConfirmRename,
  onRenameKeyDown,
  onSelectSession,
  onDeleteRequest,
}) => {
  const { t } = useTranslation()
  const editInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (isEditing) {
      setTimeout(() => editInputRef.current?.focus(), 50)
    }
  }, [isEditing])

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString(undefined, {
      month: 'short',
      day: 'numeric',
    })
  }

  if (isEditing) {
    return (
      <div className="flex-1 flex items-center gap-2">
        <Input
          ref={editInputRef}
          value={editingTitle}
          onChange={(e) => onEditingTitleChange(e.target.value)}
          onKeyDown={onRenameKeyDown}
          className="h-8 text-sm"
          disabled={actionLoading}
        />
        <Button
          size="icon"
          variant="ghost"
          className="h-7 w-7 shrink-0"
          onClick={onConfirmRename}
          disabled={actionLoading || !editingTitle.trim()}
        >
          {actionLoading ? (
            <Loader2 className="w-3.5 h-3.5 animate-spin" />
          ) : (
            <Check className="w-3.5 h-3.5" />
          )}
        </Button>
        <Button
          size="icon"
          variant="ghost"
          className="h-7 w-7 shrink-0"
          onClick={onCancelRename}
          disabled={actionLoading}
        >
          <X className="w-3.5 h-3.5" />
        </Button>
      </div>
    )
  }

  return (
    <>
      <button
        onClick={() => onSelectSession(session)}
        className="flex-1 text-left min-w-0"
      >
        <div className="text-sm font-medium truncate">{session.title}</div>
        <div className="text-xs text-muted-foreground mt-0.5">
          {formatDate(session.updatedAt)}
        </div>
      </button>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            size="icon"
            variant="ghost"
            className="h-7 w-7 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
            onClick={(e) => e.stopPropagation()}
          >
            <MoreVertical className="w-4 h-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem onClick={(e) => onStartRename(session, e)}>
            <Pencil className="w-4 h-4 mr-2" />
            {t('ui.text.rename')}
          </DropdownMenuItem>
          <DropdownMenuItem
            onClick={(e) => {
              e.stopPropagation()
              onDeleteRequest(session)
            }}
            variant="destructive"
          >
            <Trash2 className="w-4 h-4 mr-2" />
            {t('ui.text.delete')}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  )
}

const DeleteSessionDialog = ({
  session,
  actionLoading,
  onOpenChange,
  onConfirmDelete,
}) => {
  const { t } = useTranslation()

  return (
    <AlertDialog
      open={!!session}
      onOpenChange={(open) => !open && onOpenChange(false)}
    >
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            {t('ui.text.generation.deleteChat')}
          </AlertDialogTitle>
          <AlertDialogDescription>
            {t('ui.text.generation.deleteChatConfirm', {
              title: session?.title,
            })}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={actionLoading}>
            {t('ui.text.cancel')}
          </AlertDialogCancel>
          <AlertDialogAction
            onClick={onConfirmDelete}
            disabled={actionLoading}
            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
          >
            {actionLoading ? (
              <Loader2 className="w-4 h-4 animate-spin mr-2" />
            ) : null}
            {t('ui.text.delete')}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

SessionItemActions.displayName = 'SessionItemActions'
DeleteSessionDialog.displayName = 'DeleteSessionDialog'

export { SessionItemActions, DeleteSessionDialog }
