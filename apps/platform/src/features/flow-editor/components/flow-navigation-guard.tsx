import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { useTranslation } from 'react-i18next'

const FlowNavigationGuard = ({ open, onOpenChange, onConfirm, onCancel }) => {
  const { t } = useTranslation()
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{t('ui.text.unsavedChanges')}</DialogTitle>
          <DialogDescription>
            {t('ui.text.unsavedChangesDescription')}
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" size="sm" onClick={onCancel}>
            {t('ui.text.stayOnPage')}
          </Button>
          <Button size="sm" variant="destructive" onClick={onConfirm}>
            {t('ui.text.leaveWithoutSaving')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
FlowNavigationGuard.displayName = 'FlowNavigationGuard'

export { FlowNavigationGuard }
