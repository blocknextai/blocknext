import { Handle } from '@xyflow/react'
import type { Position } from '@xyflow/react'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'
import type { NodeHandle } from '@/features/flow-editor/types'
import { cn } from '@/lib/utils'

type NodeHandlesProps = {
  handles: NodeHandle[]
  type: 'source' | 'target'
  position: Position
  tooltip: string
}

// Labels sit outside the node so they never cover its content.
const LABEL_SIDE: Record<string, string> = {
  left: 'right-full mr-2 -translate-y-1/2',
  right: 'left-full ml-2 -translate-y-1/2',
  top: 'bottom-full mb-2 -translate-x-1/2',
  bottom: 'top-full mt-2 -translate-x-1/2',
}

const NodeHandles = ({
  handles,
  type,
  position,
  tooltip,
}: NodeHandlesProps) => {
  const isVertical = position === 'left' || position === 'right'

  return handles.map((handle, index) => {
    const offset = ((index + 1) / (handles.length + 1)) * 100
    const placement = isVertical
      ? { top: `${offset}%` }
      : { left: `${offset}%` }

    return (
      <div key={handle.key}>
        <Tooltip>
          <TooltipTrigger asChild>
            <Handle
              id={handle.key}
              type={type}
              position={position}
              style={placement}
              className="border-ring hover:shadow-md"
            />
          </TooltipTrigger>
          <TooltipContent side={isVertical ? 'left' : 'top'}>
            {handle.label || tooltip}
          </TooltipContent>
        </Tooltip>
        {handle.label && (
          <span
            aria-hidden
            className={cn(
              'text-muted-foreground pointer-events-none absolute text-[10px] whitespace-nowrap',
              LABEL_SIDE[position],
            )}
            style={placement}
          >
            {handle.label}
          </span>
        )}
      </div>
    )
  })
}

export { NodeHandles }
