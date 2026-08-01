import { useTranslation } from 'react-i18next'
import { RunStepDetail } from './run-step-detail'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Play, CheckCircle2, XCircle } from 'lucide-react'

const RunStepList = ({
  nodeExecutions,
  statusPrefs,
  expandedNodes,
  onToggleNode,
}) => {
  const { t } = useTranslation()

  const totalNodes = nodeExecutions?.length || 0
  const successfulNodes =
    nodeExecutions?.filter((n) => n.status === 'success').length || 0
  const failedNodes =
    nodeExecutions?.filter((n) => n.status === 'failed').length || 0

  return (
    <>
      {/* Workflow Summary */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="text-lg">
            {t('ui.text.workflowSummary')}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-full bg-blue-100 dark:bg-blue-900/20">
                <Play className="size-4 text-blue-600" />
              </div>
              <div>
                <div className="text-sm font-medium">
                  {t('ui.text.totalSteps')}
                </div>
                <div className="text-2xl font-bold">{totalNodes}</div>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-full bg-emerald-100 dark:bg-emerald-900/20">
                <CheckCircle2 className="size-4 text-emerald-600" />
              </div>
              <div>
                <div className="text-sm font-medium">
                  {t('ui.text.successful')}
                </div>
                <div className="text-2xl font-bold text-emerald-600">
                  {successfulNodes}
                </div>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-full bg-destructive/10">
                <XCircle className="size-4 text-destructive" />
              </div>
              <div>
                <div className="text-sm font-medium">{t('ui.text.failed')}</div>
                <div className="text-2xl font-bold text-destructive">
                  {failedNodes}
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Execution Steps */}
      {nodeExecutions?.length !== 0 && (
        <div className="overflow-y-auto">
          <h3 className="text-lg font-semibold mb-4">
            {t('ui.text.executionSteps')}
          </h3>
          {nodeExecutions?.map((execution, index) => (
            <RunStepDetail
              key={execution.id}
              execution={execution}
              index={index}
              statusPrefs={statusPrefs}
              isExpanded={expandedNodes.has(execution.id)}
              onToggle={() => onToggleNode(execution.id)}
            />
          ))}
        </div>
      )}
    </>
  )
}

RunStepList.displayName = 'RunStepList'

export { RunStepList }
