import platformApi from '@/lib/platform-api'

const getAll = async (
  organizationId,
  { query = '', offset = 0, limit = 10 } = {},
) => {
  const params = new URLSearchParams()
  if (query) params.set('query', query)
  params.set('offset', String(offset))
  params.set('limit', String(limit))
  return await platformApi.get(
    `/organizations/${organizationId}/triggers?${params.toString()}`,
  )
}

const update = async (organizationId, triggerId, data) => {
  return await platformApi.patch(
    `/organizations/${organizationId}/triggers/${triggerId}`,
    data,
  )
}

const regenerateWebhookToken = async (organizationId, triggerId) => {
  return await platformApi.post(
    `/organizations/${organizationId}/triggers/${triggerId}/regenerate-webhook-token`,
  )
}

const del = async (organizationId, triggerId) => {
  return await platformApi.delete(
    `/organizations/${organizationId}/triggers/${triggerId}`,
  )
}

export default {
  getAll,
  update,
  regenerateWebhookToken,
  delete: del,
}
