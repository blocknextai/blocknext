import platformApi from '@/lib/platform-api'

const getAll = async (
  organizationId,
  params?: { offset?: number; limit?: number; query?: string },
) => {
  const search = new URLSearchParams()
  if (params?.offset !== undefined) search.set('offset', String(params.offset))
  if (params?.limit !== undefined) search.set('limit', String(params.limit))
  if (params?.query) search.set('query', params.query)
  const qs = search.toString()
  return await platformApi.get(
    `/organizations/${organizationId}/workflows${qs ? `?${qs}` : ''}`,
  )
}

const getById = async (organizationId, workflowId) => {
  return await platformApi.get(
    `/organizations/${organizationId}/workflows/${workflowId}`,
  )
}

const create = async (organizationId, data) => {
  return await platformApi.post(
    `/organizations/${organizationId}/workflows`,
    data,
  )
}

const update = async (organizationId, workflowId, data) => {
  return await platformApi.patch(
    `/organizations/${organizationId}/workflows/${workflowId}`,
    data,
  )
}

const del = async (organizationId, workflowId) => {
  return await platformApi.delete(
    `/organizations/${organizationId}/workflows/${workflowId}`,
  )
}

const duplicate = async (organizationId, workflowId, data) => {
  return await platformApi.post(
    `/organizations/${organizationId}/workflows/${workflowId}/duplicate`,
    data,
  )
}

const getRun = async (organizationId, workflowId) => {
  return await platformApi.get(
    `/organizations/${organizationId}/workflows/${workflowId}/run`,
  )
}

export default {
  getAll,
  getById,
  create,
  update,
  delete: del,
  duplicate,
  getRun,
}
