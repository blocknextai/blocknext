import platformApi from '@/lib/platform-api'

const getAll = async (
  organizationId,
  { query = '', offset = 0, limit = 10 } = {},
) => {
  const params = new URLSearchParams()
  if (query) {
    params.set('query', query)
  }
  params.set('offset', String(offset))
  params.set('limit', String(limit))
  return await platformApi.get(
    `/organizations/${organizationId}/tool-invocations?${params.toString()}`,
  )
}

const getById = async (organizationId, toolInvocationId) => {
  return await platformApi.get(
    `/organizations/${organizationId}/tool-invocations/${toolInvocationId}`,
  )
}

export default {
  getAll,
  getById,
}
