import useSWR from 'swr'
import mcpService from '@/features/mcp/services/mcp'

export function useMcpServers() {
  const { data, error, isLoading, mutate } = useSWR(
    ['mcp-servers'],
    () => mcpService.getServers(),
    { keepPreviousData: true },
  )
  return {
    servers: data?.data ?? [],
    isLoading,
    error,
    mutate,
  }
}
