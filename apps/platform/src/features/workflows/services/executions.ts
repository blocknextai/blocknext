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
    `/organizations/${organizationId}/executions?${params.toString()}`,
  )
}

const getById = async (organizationId, executionId) => {
  return await platformApi.get(
    `/organizations/${organizationId}/executions/${executionId}`,
  )
}

const del = async (organizationId, executionId) => {
  return await platformApi.delete(
    `/organizations/${organizationId}/executions/${executionId}`,
  )
}

const bulkDelete = async (organizationId, data) => {
  return await platformApi.delete(
    `/organizations/${organizationId}/executions/bulk`,
    data,
  )
}

export default {
  getAll,
  getById,
  delete: del,
  bulkDelete,
}
