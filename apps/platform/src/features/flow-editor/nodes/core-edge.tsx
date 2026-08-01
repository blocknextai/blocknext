import { getBezierPath } from '@xyflow/react'

const CoreEdge = ({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  style = {},
  markerEnd,
  selected,
}) => {
  const xEqual = sourceX === targetX
  const yEqual = sourceY === targetY
  const [edgePath] = getBezierPath({
    sourceX: xEqual ? sourceX + 0.0001 : sourceX,
    sourceY: yEqual ? sourceY + 0.0001 : sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  })
  const gradientId = `gradient-${id}`
  return (
    <>
      <path
        id={id}
        className="react-flow__edge-path"
        d={edgePath}
        style={{
          stroke: selected ? 'var(--foreground)' : `url(#${gradientId})`,
          strokeWidth: 4,
          ...style,
        }}
        markerEnd={markerEnd}
      />
    </>
  )
}
CoreEdge.displayName = 'CoreEdge'

export default CoreEdge
