import useSWR from 'swr'
import mcpService from '@/features/mcp/services/mcp'
import type { McpServer } from '@/features/mcp/services/mcp'

const fetchServers = async (): Promise<McpServer[]> => {
  const response = await mcpService.getServers()
  return response.data ?? []
}

export function useMcpServers() {
  const { data, error, isLoading, mutate } = useSWR(
    ['mcp-servers'],
    fetchServers,
    {
      keepPreviousData: true,
    },
  )
  return {
    servers: data ?? [],
    isLoading,
    error,
    mutate,
  }
}
