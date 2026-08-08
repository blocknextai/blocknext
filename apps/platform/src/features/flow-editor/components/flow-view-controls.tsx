import { memo } from 'react'
import { useReactFlow } from '@xyflow/react'
import { useTranslation } from 'react-i18next'
import { Lock, LockOpen, Maximize, Minus, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'

type FlowViewControlsProps = {
  locked: boolean
  setLocked: (locked: boolean) => void
  previewMode?: boolean
}

const FlowViewControls = memo(
  ({ locked, setLocked, previewMode }: FlowViewControlsProps) => {
    const { zoomIn, zoomOut, fitView } = useReactFlow()
    const { t } = useTranslation()

    if (previewMode) {
      return null
    }

    const actions = [
      {
        key: 'zoomIn',
        icon: Plus,
        label: t('ui.text.zoomIn'),
        run: () => zoomIn(),
      },
      {
        key: 'zoomOut',
        icon: Minus,
        label: t('ui.text.zoomOut'),
        run: () => zoomOut(),
      },
      {
        key: 'fitView',
        icon: Maximize,
        label: t('ui.text.fitView'),
        run: () => fitView({ padding: 0.2, duration: 200 }),
      },
      {
        key: 'lock',
        icon: locked ? Lock : LockOpen,
        label: locked ? t('ui.text.unlockCanvas') : t('ui.text.lockCanvas'),
        run: () => setLocked(!locked),
      },
    ]

    return (
      <div
        data-tour="flow-view-controls"
        className="bg-background dark:bg-accent border-zinc-200 text-foreground/80 relative z-40 flex gap-1.5 rounded-xl border p-1.5 shadow-xs dark:border-zinc-700"
      >
        {actions.map(({ key, icon: Icon, label, run }) => (
          <Tooltip key={key}>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                className="size-8"
                onClick={run}
                aria-label={label}
              >
                <Icon className="size-5" aria-hidden="true" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>{label}</TooltipContent>
          </Tooltip>
        ))}
      </div>
    )
  },
)
FlowViewControls.displayName = 'FlowViewControls'

export { FlowViewControls }
