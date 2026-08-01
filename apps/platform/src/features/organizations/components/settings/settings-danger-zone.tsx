import { Button } from '@/components/ui/button'
import { useTranslation } from 'react-i18next'
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

const SettingsDangerZone = ({ delOpen, setDelOpen, onDelete }) => {
  const { t } = useTranslation()

  return (
    <AlertDialog open={delOpen} onOpenChange={setDelOpen}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{t('ui.text.areYouSure')}</AlertDialogTitle>
          <AlertDialogDescription>
            {t('ui.text.deleteOrgConfirmation')}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{t('ui.text.cancel')}</AlertDialogCancel>
          <AlertDialogAction asChild>
            <Button variant="destructive" onClick={onDelete}>
              {t('ui.text.delete')}
            </Button>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

SettingsDangerZone.displayName = 'SettingsDangerZone'

export { SettingsDangerZone }
