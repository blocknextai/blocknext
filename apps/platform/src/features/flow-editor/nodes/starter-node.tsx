import { memo, useState, useMemo, useEffect } from 'react'
import { Handle, Position, useUpdateNodeInternals } from '@xyflow/react'
import { Card, CardContent } from '@/components/ui/card'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'
import { useTranslation } from 'react-i18next'
import { flowCategoryPreferences } from '@/lib/flow-categories'

const StarterNode = memo(({ selected, data }) => {
  const updateNodeInternals = useUpdateNodeInternals()
  const [position] = useState('l-r')
  const { t } = useTranslation()

  const { prefs, Icon } = useMemo(() => {
    if (!data?.category || !flowCategoryPreferences[data.category]) {
      return { prefs: { color: '#ffffff' }, Icon: () => null }
    }
    const p = flowCategoryPreferences[data.category]
    return { prefs: p, Icon: p.icon }
  }, [data?.category])

  useEffect(() => {
    updateNodeInternals(data?.id)
  }, [position, updateNodeInternals, data?.id])

  return (
    <div
      className={`rounded-xl relative bg-background transition-all duration-200 hover:shadow-md border-2 border-current spread-transition
        ${selected ? 'spread-container' : ''}`}
      style={{ color: prefs?.color }}
    >
      <div className="bg-current/20 rounded-lg">
        <Card className="items-start bg-transparent py-2! shadow-none border-none">
          <CardContent className="pl-3 pr-4">
            <div className="flex items-center gap-4 cursor-pointer">
              <div
                style={{ color: prefs?.color }}
                className="size-10 rounded-lg flex items-center justify-center bg-current shrink-0"
              >
                {Icon && <Icon className="size-6 text-zinc-50" />}
              </div>
              <div className="flex flex-1 flex-col items-start justify-start gap-0">
                <div
                  className="capitalize text-sm"
                  style={{ color: prefs?.color }}
                >
                  {prefs?.labelKey ? t(prefs.labelKey) : data?.category}
                </div>
                <div className="text-sm">
                  {data?.title ? t(data.title) : null}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Tooltip>
          <TooltipTrigger asChild>
            <Handle
              type="source"
              position={Position.Right}
              className="border-ring hover:shadow-md"
            />
          </TooltipTrigger>
          <TooltipContent side="right">
            {t('ui.text.dataOutput')}
          </TooltipContent>
        </Tooltip>
      </div>
    </div>
  )
})

StarterNode.displayName = 'StarterNode'

export default StarterNode
