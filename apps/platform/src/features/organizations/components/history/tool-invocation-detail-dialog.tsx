import { useTranslation } from 'react-i18next'
import { formatDistanceStrict } from 'date-fns'
import { XCircle } from 'lucide-react'
import TimeAgoI18n from '@/components/shared/timeagoi18'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { CodeBlock } from '@/components/shared/code-block'
import { Loading } from '@/components/shared/loading'
import StatusBadge from '@/components/shared/status-badge'
import { useToolInvocation } from '@/features/workflows'

type ToolInvocationDetailDialogProps = {
  organizationId: string
  toolInvocationId: string | null
  onOpenChange: (open: boolean) => void
}

const calculateDuration = (startedAt: string, completedAt: string) => {
  if (!startedAt || !completedAt) {
    return '-'
  }

  const start = new Date(startedAt).getTime()
  const end = new Date(completedAt).getTime()

  if (isNaN(start) || isNaN(end)) {
    return '-'
  }

  return formatDistanceStrict(0, end - start)
}

const ToolInvocationDetailDialog = ({
  organizationId,
  toolInvocationId,
  onOpenChange,
}: ToolInvocationDetailDialogProps) => {
  const { t } = useTranslation()
  const { toolInvocation, isLoading } = useToolInvocation(
    organizationId,
    toolInvocationId,
  )

  const duration = calculateDuration(
    toolInvocation?.startedAt,
    toolInvocation?.completedAt,
  )

  return (
    <Dialog open={!!toolInvocationId} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            {toolInvocation?.toolId || t('ui.text.toolCall')}
            {toolInvocation?.source && (
              <Badge variant="outline">{toolInvocation.source}</Badge>
            )}
          </DialogTitle>
          <DialogDescription className="sr-only">
            {t('ui.text.toolCall')}
          </DialogDescription>
        </DialogHeader>

        {isLoading || !toolInvocation ? (
          <Loading />
        ) : (
          <div className="flex flex-col gap-4 max-h-[60vh] overflow-y-auto">
            <div className="flex items-center gap-2">
              <StatusBadge
                status={toolInvocation.status}
                className="w-fit"
                title={toolInvocation.errorMessage}
              />
              {toolInvocation.apiKeyName && (
                <span className="text-muted-foreground text-sm">
                  {t('ui.text.apiKey')}: {toolInvocation.apiKeyName}
                </span>
              )}
            </div>

            <div className="grid grid-cols-1 gap-2 sm:grid-cols-3">
              <div className="flex flex-col gap-1">
                <span className="text-muted-foreground text-xs">
                  {t('ui.text.started')}
                </span>
                <span className="text-sm">
                  <TimeAgoI18n date={toolInvocation.startedAt} />
                </span>
              </div>
              <div className="flex flex-col gap-1">
                <span className="text-muted-foreground text-xs">
                  {t('ui.text.finished')}
                </span>
                <span className="text-sm">
                  <TimeAgoI18n date={toolInvocation.completedAt} />
                </span>
              </div>
              <div className="flex flex-col gap-1">
                <span className="text-muted-foreground text-xs">
                  {t('ui.text.duration')}
                </span>
                <span className="text-sm">{duration}</span>
              </div>
            </div>

            {toolInvocation.errorMessage && (
              <div className="flex gap-2 items-center p-3 bg-destructive/10 rounded-lg">
                <XCircle className="size-4 text-destructive shrink-0" />
                <span className="text-destructive text-sm">
                  {toolInvocation.errorMessage}
                </span>
              </div>
            )}

            <div className="flex flex-col gap-2">
              <span className="text-sm font-medium">
                {t('ui.text.parameters')}
              </span>
              <CodeBlock
                code={JSON.stringify(toolInvocation.parameters ?? {}, null, 2)}
              />
            </div>

            {toolInvocation.outputs && (
              <div className="flex flex-col gap-2">
                <span className="text-sm font-medium">
                  {t('ui.text.outputs')}
                </span>
                <CodeBlock
                  code={JSON.stringify(toolInvocation.outputs, null, 2)}
                />
              </div>
            )}
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}

ToolInvocationDetailDialog.displayName = 'ToolInvocationDetailDialog'

export { ToolInvocationDetailDialog }
