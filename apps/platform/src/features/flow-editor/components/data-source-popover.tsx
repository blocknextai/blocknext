import { Braces, Webhook } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { useTranslation } from 'react-i18next'

interface ContextMenuItem {
  label: string
  value: string
  isEditable?: boolean
}

interface DataSourcePopoverProps {
  contextMenu: ContextMenuItem[]
  triggerVariables?: string[]
  onSelect: (value: string) => void
}

export function DataSourcePopover({
  contextMenu,
  triggerVariables,
  onSelect,
}: DataSourcePopoverProps) {
  const { t } = useTranslation()

  const hasContextMenu = contextMenu && contextMenu.length > 0
  const hasTriggerVariables = triggerVariables && triggerVariables.length > 0

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="icon-sm"
          className="shrink-0"
          type="button"
        >
          <Braces className="size-4" />
        </Button>
      </PopoverTrigger>
      <PopoverContent align="end" className="w-64 max-h-72 overflow-y-auto p-2">
        {hasTriggerVariables && (
          <>
            <div className="text-xs font-medium text-muted-foreground px-2 py-1 mb-1 flex items-center gap-1.5">
              <Webhook className="size-3" />
              {t('ui.text.trigger')}
            </div>
            <div className="grid gap-1 mb-2">
              {triggerVariables.map((variable) => (
                <button
                  key={variable}
                  type="button"
                  className="flex items-center gap-2 w-full rounded-md px-2 py-1.5 text-left text-xs font-mono text-muted-foreground hover:bg-accent hover:text-foreground transition-colors cursor-pointer"
                  onClick={() => onSelect(variable)}
                >
                  {variable}
                </button>
              ))}
            </div>
          </>
        )}

        {hasContextMenu && (
          <>
            <div className="text-xs font-medium text-muted-foreground px-2 py-1 mb-1 flex items-center gap-1.5">
              <Braces className="size-3" />
              {t('ui.text.dataSource')}
            </div>
            <div className="grid gap-1">
              {contextMenu.map((item, index) => (
                <button
                  key={`${item.value}-${index}`}
                  type="button"
                  className="hover:bg-accent flex w-full cursor-pointer flex-col items-start gap-0.5 rounded-md px-2 py-1.5 text-left transition-colors"
                  onClick={() => onSelect(item.value)}
                >
                  <span className="text-foreground w-full truncate text-xs">
                    {item.label || item.value}
                  </span>
                  <span className="text-muted-foreground w-full truncate font-mono text-[10px]">
                    {item.value}
                  </span>
                </button>
              ))}
            </div>
          </>
        )}
      </PopoverContent>
    </Popover>
  )
}
