import { memo } from 'react'

import { GroupNode } from '@/features/flow-editor/nodes/core/labeled-group-node'

const CoreGroup = memo(({ selected, data }) => {
  const nodeCount = data?.nodeCount ?? 0
  const width = nodeCount > 1 ? nodeCount * 600 : 600

  return (
    <GroupNode
      className="h-100 rounded-lg z-0! "
      style={{ width: `${width}px` }}
      selected={selected}
      label={data.label}
    />
  )
})
CoreGroup.displayName = 'CoreGroup'

export default CoreGroup
