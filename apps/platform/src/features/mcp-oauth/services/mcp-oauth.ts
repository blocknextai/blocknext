import mcpApi from '@/lib/mcp-api'

export type McpOAuthClient = {
  clientId: string
  name: string
  logoUri?: string
  clientUri?: string
  isMetadataDocument: boolean
}

export type McpOAuthScope = {
  scope: string
  description: string
}

export type McpAuthorizationRequest = {
  id: string
  client: McpOAuthClient
  scopes: McpOAuthScope[]
  resource: string
  redirectUri: string
  status: 'pending' | 'approved' | 'denied'
  isPending: boolean
  expiresAt: string
}

export type McpAuthorizationRedirect = {
  redirectUri: string
}

const getAuthorizationRequest = async (authorizationRequestId: string) => {
  return await mcpApi.get<McpAuthorizationRequest>(
    `/mcp-oauth/authorization-requests/${authorizationRequestId}`,
  )
}

const approveAuthorizationRequest = async (
  authorizationRequestId: string,
  data: { organizationId: string; scopes?: string[] },
) => {
  return await mcpApi.post<McpAuthorizationRedirect>(
    `/mcp-oauth/authorization-requests/${authorizationRequestId}/approve`,
    data,
  )
}

const denyAuthorizationRequest = async (authorizationRequestId: string) => {
  return await mcpApi.post<McpAuthorizationRedirect>(
    `/mcp-oauth/authorization-requests/${authorizationRequestId}/deny`,
  )
}

export default {
  getAuthorizationRequest,
  approveAuthorizationRequest,
  denyAuthorizationRequest,
}
