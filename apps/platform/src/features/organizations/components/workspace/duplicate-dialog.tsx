import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { DialogClose } from '@radix-ui/react-dialog'
import { useTranslation } from 'react-i18next'

const DuplicateDialog = ({
  open,
  onOpenChange,
  selectedWorkflow,
  setSelectedWorkflow,
  onConfirmDuplicate,
}) => {
  const { t } = useTranslation()

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{t('ui.text.duplicateFlow')}</DialogTitle>
          <DialogDescription>
            {t('ui.text.duplicateFlowDescription')}
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4">
          <div className="grid gap-3">
            <Label htmlFor={`duplicated-title`}>{t('ui.text.title')}</Label>
            <Input
              id={`duplicated-title`}
              name={`duplicated-title`}
              defaultValue={selectedWorkflow?.title}
              onChange={(e) =>
                setSelectedWorkflow({
                  ...selectedWorkflow,
                  title: e.target.value,
                })
              }
            />
          </div>
          <div className="grid gap-3">
            <Label htmlFor={`duplicated-desc`}>
              {t('ui.text.description')}
            </Label>
            <Input
              id={`duplicated-desc`}
              name={`duplicated-desc`}
              defaultValue={selectedWorkflow?.description}
              onChange={(e) =>
                setSelectedWorkflow({
                  ...selectedWorkflow,
                  description: e.target.value,
                })
              }
            />
          </div>
        </div>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline">{t('ui.text.cancel')}</Button>
          </DialogClose>

          <Button onClick={onConfirmDuplicate}>{t('ui.text.duplicate')}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
DuplicateDialog.displayName = 'DuplicateDialog'

export { DuplicateDialog }
