import { Button } from '@/components/ui/button'
import {
  DialogClose,
  DialogFooter,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  Dialog,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useTranslation } from 'react-i18next'

const EditDialog = ({
  open,
  onOpenChange,
  selectedWorkflow,
  setSelectedWorkflow,
  onConfirmEdit,
}) => {
  const { t } = useTranslation()

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{t('ui.text.editFlow')}</DialogTitle>
          <DialogDescription>
            {t('ui.text.editFlowDescription')}
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4">
          <div className="grid gap-3">
            <Label htmlFor={`edited-title`}>{t('ui.text.title')}</Label>
            <Input
              id={`edited-title`}
              name={`edited-title`}
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
            <Label htmlFor={`edited-desc`}>{t('ui.text.description')}</Label>
            <Input
              id={`edited-desc`}
              name={`edited-desc`}
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

          <Button onClick={onConfirmEdit}>{t('ui.text.save')}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
EditDialog.displayName = 'EditDialog'

export { EditDialog }
