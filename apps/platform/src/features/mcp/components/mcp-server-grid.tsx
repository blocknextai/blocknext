import { McpServerCard } from '@/features/mcp/components/mcp-server-card'
import type { McpServer } from '@/features/mcp/services/mcp'

type McpServerGridProps = {
  servers: McpServer[]
  onCardClick?: (server: McpServer) => void
}

const McpServerGrid = ({ servers, onCardClick }: McpServerGridProps) => {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3 sm:gap-4">
      {servers.map((server) => (
        <McpServerCard key={server.id} server={server} onClick={onCardClick} />
      ))}
    </div>
  )
}

McpServerGrid.displayName = 'McpServerGrid'

export { McpServerGrid }
