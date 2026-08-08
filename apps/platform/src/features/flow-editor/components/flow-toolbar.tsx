import { RotateCcw, Play, Save, Code2 } from 'lucide-react'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'
import { Button } from '@/components/ui/button'
import { useTranslation } from 'react-i18next'

const FlowToolbar = ({
  previewMode,
  hasUnsavedChanges,
  deleteLastSaved,
  runFlow,
  setOpen,
  openApiSheet,
  viewControls,
}) => {
  const { t } = useTranslation()

  if (previewMode) {
    return null
  }

  return (
    <div className="absolute bottom-2 flex w-full items-center justify-center gap-2 p-4">
      {viewControls}
      <div
        data-tour="flow-toolbar"
        className="flex
                gap-1.5 bg-background border
              dark:bg-accent
              dark:border-zinc-700
              border-zinc-200 p-1.5
              text-foreground/80 rounded-xl relative z-40 shadow-xs"
      >
        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              variant={'ghost'}
              className="size-8"
              onClick={deleteLastSaved}
              aria-label={t('ui.text.revertAllChanges')}
            >
              <RotateCcw className="size-5" aria-hidden="true" />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="top" size="size-0">
            {t('ui.text.revertAllChanges')}
          </TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              variant={'ghost'}
              className="size-8"
              onClick={runFlow}
              aria-label={t('ui.text.runFlow')}
            >
              <Play className="size-5" aria-hidden="true" />
            </Button>
          </TooltipTrigger>
          <TooltipContent side="top" size="size-0">
            {t('ui.text.runFlow')}
          </TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <div className="relative">
              <Button
                variant={'ghost'}
                className="size-8"
                onClick={() => setOpen(true)}
                aria-label={
                  hasUnsavedChanges
                    ? `${t('ui.text.saveFlow')} (${t('ui.text.hasUnsavedChanges')})`
                    : t('ui.text.saveFlow')
                }
              >
                <Save className="size-5" aria-hidden="true" />
              </Button>
              {hasUnsavedChanges && (
                <div
                  className="absolute -top-1 -right-1 w-3 h-3 bg-yellow-500 rounded-full animate-pulse"
                  aria-hidden="true"
                ></div>
              )}
            </div>
          </TooltipTrigger>
          <TooltipContent side="top" size="size-0">
            {t('ui.text.saveFlow')}{' '}
            {hasUnsavedChanges && `(${t('ui.text.hasUnsavedChanges')})`}
          </TooltipContent>
        </Tooltip>
        {openApiSheet && (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant={'ghost'}
                className="size-8"
                onClick={openApiSheet}
                aria-label={t('ui.text.apiTrigger')}
              >
                <Code2 className="size-5" aria-hidden="true" />
              </Button>
            </TooltipTrigger>
            <TooltipContent side="top" size="size-0">
              {t('ui.text.apiTrigger')}
            </TooltipContent>
          </Tooltip>
        )}
      </div>
    </div>
  )
}
FlowToolbar.displayName = 'FlowToolbar'

export { FlowToolbar }
