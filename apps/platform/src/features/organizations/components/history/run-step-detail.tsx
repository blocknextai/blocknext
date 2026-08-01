import { useTranslation } from 'react-i18next'
import { formatDistanceStrict, format } from 'date-fns'
import { Button } from '@/components/ui/button'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  XCircle,
  Clock,
  ChevronDown,
  ChevronRight,
  ArrowRightToLine,
} from 'lucide-react'
import { OutputRenderer } from '@/features/organizations/components/history/output-renderer'

const RunStepDetail = ({
  execution,
  index,
  statusPrefs,
  isExpanded,
  onToggle,
}) => {
  const { t } = useTranslation()

  const formatDuration = (startedAt?: string, completedAt?: string) => {
    try {
      if (!startedAt || !completedAt) return t('ui.text.workInProgress')
      const start = new Date(startedAt)
      const end = new Date(completedAt)

      if (isNaN(start.getTime()) || isNaN(end.getTime())) {
        return t('ui.text.invalid')
      }

      const diff = end.getTime() - start.getTime()
      return formatDistanceStrict(0, diff)
    } catch {
      return t('ui.text.error')
    }
  }

  const formatTime = (dateString?: string) => {
    try {
      if (!dateString) return t('ui.text.workInProgress')
      const date = new Date(dateString)
      if (isNaN(date.getTime())) return t('ui.text.invalid')
      return format(date, 'HH:mm:ss.SSS')
    } catch {
      return t('ui.text.error')
    }
  }

  const duration = formatDuration(execution.startedAt, execution.completedAt)

  return (
    <Collapsible open={isExpanded} onOpenChange={() => onToggle()}>
      <Card className="mb-4">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div
                className={`p-2 rounded-full ${statusPrefs[execution.status]?.bgColor}`}
              >
                {statusPrefs[execution.status]?.icon}
              </div>
              <div className="flex flex-col">
                <CollapsibleTrigger asChild className="cursor-pointer">
                  <CardTitle className="text-sm font-medium">
                    {t('ui.text.step')} {index + 1}:{' '}
                    {t(execution.inputs?.[0]?.title) || t(execution.nodeType)}
                  </CardTitle>
                </CollapsibleTrigger>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Badge
                variant={
                  execution.status === 'success'
                    ? 'default'
                    : execution.status === 'failed'
                      ? 'destructive'
                      : 'secondary'
                }
              >
                {t(`ui.text.status.${execution.status}`)}
              </Badge>
              <div className="text-xs text-muted-foreground flex items-center gap-1">
                <Clock className="size-3" />
                {duration}
              </div>
              <CollapsibleTrigger asChild className="cursor-pointer">
                <Button variant="ghost" size="sm">
                  {isExpanded ? (
                    <ChevronDown className="size-4" />
                  ) : (
                    <ChevronRight className="size-4" />
                  )}
                </Button>
              </CollapsibleTrigger>
            </div>
          </div>
        </CardHeader>

        <CollapsibleContent>
          <CardContent className="pt-0">
            <div className="space-y-4">
              {/* Inputs */}
              {Array.isArray(execution.inputs) &&
                execution.inputs.some(
                  (input: any) =>
                    input?.runtimeInstruction || input?.runtimePrompt,
                ) && (
                  <div>
                    <h4 className="text-sm font-medium mb-3 flex items-center gap-2">
                      <ArrowRightToLine className="size-4" />
                      {t('ui.text.inputs')}
                    </h4>
                    <div className="space-y-3">
                      {execution.inputs.map(
                        (input: any, inputIndex: number) => {
                          if (
                            !input?.runtimeInstruction &&
                            !input?.runtimePrompt
                          ) {
                            return null
                          }
                          return (
                            <div key={inputIndex} className="space-y-2">
                              {input.runtimeInstruction && (
                                <div>
                                  <div className="text-xs text-muted-foreground mb-1 font-mono">
                                    runtimeInstruction
                                  </div>
                                  <div className="bg-muted/50 rounded-lg p-3 text-sm whitespace-pre-wrap break-words">
                                    {input.runtimeInstruction}
                                  </div>
                                </div>
                              )}
                              {input.runtimePrompt && (
                                <div>
                                  <div className="text-xs text-muted-foreground mb-1 font-mono">
                                    runtimePrompt
                                  </div>
                                  <div className="bg-muted/50 rounded-lg p-3 text-sm whitespace-pre-wrap break-words">
                                    {input.runtimePrompt}
                                  </div>
                                </div>
                              )}
                            </div>
                          )
                        },
                      )}
                    </div>
                  </div>
                )}

              {/* Outputs */}
              {execution.outputs && (
                <OutputRenderer outputs={execution.outputs} />
              )}

              {/* Error Message */}
              {execution.errorMessage && (
                <div>
                  <h4 className="text-sm font-medium mb-2 flex items-center gap-2 text-destructive">
                    <XCircle className="size-4" />
                    {t('ui.text.error')}
                  </h4>
                  <div className="bg-destructive/10 rounded-lg p-3 text-xs text-destructive">
                    {execution.errorMessage}
                  </div>
                </div>
              )}

              {/* Timing */}
              <div className="grid grid-cols-2 gap-4 text-xs">
                <div>
                  <span className="text-muted-foreground">
                    {t('ui.text.started')}:
                  </span>
                  <div className="font-mono">
                    {formatTime(execution.startedAt)}
                  </div>
                </div>
                <div>
                  <span className="text-muted-foreground">
                    {t('ui.text.completed')}:
                  </span>
                  <div className="font-mono">
                    {formatTime(execution.completedAt)}
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </CollapsibleContent>
      </Card>
    </Collapsible>
  )
}

RunStepDetail.displayName = 'RunStepDetail'

export { RunStepDetail }
