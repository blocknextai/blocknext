import { flowCategoryPreferences } from '@/lib/flow-categories'

const FALLBACK_COLOR = '#888888'

const EdgeGradients = ({ edges, nodes }) => {
  return (
    <svg style={{ width: 0, height: 0, position: 'absolute', zIndex: -1 }}>
      <defs>
        {edges.map((edge) => {
          const sourceNode = nodes.find((node) => node.id === edge.source)
          const targetNode = nodes.find((node) => node.id === edge.target)

          const sourceColor =
            flowCategoryPreferences[sourceNode?.data?.category]?.color ||
            FALLBACK_COLOR
          const targetColor =
            flowCategoryPreferences[targetNode?.data?.category]?.color ||
            FALLBACK_COLOR

          const gradientId = `gradient-${edge.id}`

          return (
            <linearGradient
              key={gradientId}
              id={gradientId}
              x1="0%"
              y1="0%"
              x2="100%"
              y2="0%"
            >
              <stop
                offset="0%"
                style={{
                  stopColor: sourceColor,
                  stopOpacity: 1,
                }}
              />
              <stop
                offset="100%"
                style={{
                  stopColor: targetColor,
                  stopOpacity: 1,
                }}
              />
            </linearGradient>
          )
        })}
      </defs>
    </svg>
  )
}
EdgeGradients.displayName = 'EdgeGradients'

export { EdgeGradients }
