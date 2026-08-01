import platformApi from '@/lib/platform-api'

const getOrganizationApiKeys = async (
  organizationId,
  { offset = 0, limit = 10, query = '' } = {},
) => {
  const params = new URLSearchParams()
  params.set('offset', String(offset))
  params.set('limit', String(limit))
  if (query) params.set('query', query)
  return await platformApi.get(
    `/organizations/${organizationId}/api-keys?${params.toString()}`,
  )
}

const createOrganizationApiKey = async (organizationId, data) => {
  return await platformApi.post(
    `/organizations/${organizationId}/api-keys`,
    data,
  )
}

const regenerateOrganizationApiKey = async (organizationId, apiKeyId, data) => {
  return await platformApi.post(
    `/organizations/${organizationId}/api-keys/${apiKeyId}/regenerate`,
    data,
  )
}

const deleteOrganizationApiKey = async (organizationId, apiKeyId) => {
  return await platformApi.delete(
    `/organizations/${organizationId}/api-keys/${apiKeyId}`,
  )
}

const getApiKeyScopes = async () => {
  return await platformApi.get('/api-keys/scopes')
}

export default {
  getOrganizationApiKeys,
  createOrganizationApiKey,
  regenerateOrganizationApiKey,
  deleteOrganizationApiKey,
  getApiKeyScopes,
}
