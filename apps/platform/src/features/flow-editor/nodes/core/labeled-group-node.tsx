import { forwardRef } from 'react'
import { Panel } from '@xyflow/react'
import { BaseNode } from '@/features/flow-editor/nodes/core/base-node'
import { cn } from '@/lib/utils'

/* GROUP NODE Label ------------------------------------------------------- */

const GroupNodeLabel = forwardRef(({ children, className, ...props }, ref) => {
  return (
    <div ref={ref} className="h-full w-full" {...props}>
      <div
        className={cn(
          'w-fit bg-secondary p-2 text-xs text-card-foreground',
          className,
        )}
      >
        {children}
      </div>
    </div>
  )
})

GroupNodeLabel.displayName = 'GroupNodeLabel'

/* GROUP NODE -------------------------------------------------------------- */

const GroupNode = forwardRef(({ selected, label, position, ...props }, ref) => {
  const getLabelClassName = (position) => {
    switch (position) {
      case 'top-left':
        return 'rounded-br-sm'
      case 'top-center':
        return 'rounded-b-sm'
      case 'top-right':
        return 'rounded-bl-sm'
      case 'bottom-left':
        return 'rounded-tr-sm'
      case 'bottom-right':
        return 'rounded-tl-sm'
      case 'bottom-center':
        return 'rounded-t-sm'
      default:
        return 'rounded-br-sm'
    }
  }

  return (
    <BaseNode
      ref={ref}
      selected={selected}
      className="h-full overflow-hidden rounded-sm bg-white bg-opacity-50 p-0"
      {...props}
    >
      <Panel className={cn('m-0 p-0')} position={position}>
        {label && (
          <GroupNodeLabel className={getLabelClassName(position)}>
            {label}
          </GroupNodeLabel>
        )}
      </Panel>
    </BaseNode>
  )
})

GroupNode.displayName = 'GroupNode'

export { GroupNodeLabel, GroupNode }
