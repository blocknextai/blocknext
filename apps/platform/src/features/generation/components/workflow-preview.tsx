import { useTranslation } from 'react-i18next'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Workflow, Check } from 'lucide-react'

// Fix over-escaped backslashes from AI response
// AI sometimes outputs \\" instead of \" for escaped quotes in strings
function sanitizeJson(jsonStr) {
  // Replace \\" with \" (the AI double-escapes quotes)
  return jsonStr.replace(/\\\\"/g, '\\"')
}

const WorkflowPreview = ({ json, onApplyWorkflow }) => {
  const { t } = useTranslation()
  let parsed = null
  let parseError = null

  // Try parsing with sanitization if direct parse fails
  try {
    parsed = JSON.parse(json)
  } catch (e) {
    try {
      parsed = JSON.parse(sanitizeJson(json))
    } catch (e2) {
      parseError = e2.message
    }
  }

  if (!parsed) {
    return (
      <Card className="my-3 border-destructive/20 bg-destructive/5">
        <CardContent className="p-4 flex flex-col gap-2">
          <div className="flex items-center gap-2">
            <Workflow className="w-4 h-4 text-destructive" />
            <span className="text-sm font-medium text-destructive">
              {t('ui.text.generation.invalidWorkflow')}
            </span>
          </div>
          {parseError && (
            <p className="text-xs text-muted-foreground">{parseError}</p>
          )}
        </CardContent>
      </Card>
    )
  }

  const nodes = parsed.nodes || []
  const edges = parsed.edges || []

  return (
    <Card className="my-3 border-primary/20">
      <CardContent className="p-4 flex flex-col gap-3">
        <div className="flex items-center gap-2">
          <Workflow className="w-4 h-4 text-primary" />
          <span className="text-sm font-medium">
            {t('ui.text.generation.workflow')}
          </span>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="secondary" className="text-xs">
            {t('ui.text.generation.nodesCount', { count: nodes.length })}
          </Badge>
          <Badge variant="outline" className="text-xs">
            {t('ui.text.generation.edgesCount', { count: edges.length })}
          </Badge>
        </div>
        <Button
          size="sm"
          className="w-full"
          onClick={() => onApplyWorkflow({ nodes, edges })}
        >
          <Check className="w-3.5 h-3.5 mr-1.5" />
          {t('ui.text.generation.applyToCanvas')}
        </Button>
      </CardContent>
    </Card>
  )
}

WorkflowPreview.displayName = 'WorkflowPreview'

export default WorkflowPreview
