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
import { useTranslation } from 'react-i18next'

const DeleteDialog = ({ open, onOpenChange, onConfirmDelete }) => {
  const { t } = useTranslation()

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{t('ui.text.areYouSure')}</AlertDialogTitle>
          <AlertDialogDescription>
            {t('ui.text.deleteFlowConfirmation')}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{t('ui.text.cancel')}</AlertDialogCancel>
          <AlertDialogAction onClick={onConfirmDelete}>
            {t('ui.text.continue')}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
DeleteDialog.displayName = 'DeleteDialog'

export { DeleteDialog }
