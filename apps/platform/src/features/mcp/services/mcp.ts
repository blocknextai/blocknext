import type { IconSource } from '@/features/flow-editor/types'
import mcpApi from '@/lib/mcp-api'

export type McpToolSchemaProperty = {
  type?: string
  description?: string
  enum?: unknown[]
  items?: unknown
  properties?: Record<string, unknown>
  required?: string[]
  [key: string]: unknown
}

export type McpToolSchema = {
  type?: string
  description?: string
  properties?: Record<string, McpToolSchemaProperty>
  required?: string[]
  [key: string]: unknown
} | null

export type McpTool = {
  id: string
  version: string
  name: string
  description?: string
  icon?: IconSource
  inputSchema?: McpToolSchema
  outputSchema?: McpToolSchema
  supportedCredentials?: string[]
}

export type McpServer = {
  id: string
  name: string
  description?: string
  icon?: IconSource
  version: string
  url?: string
  tools: McpTool[]
}

const getServers = async () => {
  return await mcpApi.get<McpServer[]>(`/servers`)
}

export default {
  getServers,
}
