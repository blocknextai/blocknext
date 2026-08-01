import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useTranslation } from 'react-i18next'

const FlowSaveDialog = ({
  open,
  onOpenChange,
  flowData,
  updateFlowData,
  saveFlow,
}) => {
  const { t } = useTranslation()
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{t('ui.text.saveFlow')}</DialogTitle>
          <DialogDescription>
            {t('ui.text.saveFlowDescription')}
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4">
          <div className="grid gap-3">
            <Label htmlFor={`f-title`}>{t('ui.text.title')}</Label>
            <Input
              id={`f-title`}
              name={`f-title`}
              value={flowData.title}
              onChange={(e) => updateFlowData('title', e.target.value)}
            />
          </div>
          <div className="grid gap-3">
            <Label htmlFor={`f-desc`}>{t('ui.text.description')}</Label>
            <Input
              id={`f-desc`}
              name={`f-decs`}
              value={flowData.description}
              onChange={(e) => updateFlowData('description', e.target.value)}
            />
          </div>
        </div>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline" size="sm">
              {t('ui.text.cancel')}
            </Button>
          </DialogClose>
          <DialogClose asChild>
            <Button size="sm" onClick={saveFlow}>
              {t('ui.text.save')}
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
FlowSaveDialog.displayName = 'FlowSaveDialog'

export { FlowSaveDialog }
