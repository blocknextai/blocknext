import platformApi from '@/lib/platform-api'

const getAll = async () => {
  return await platformApi.get('/organizations')
}

const getById = async (organizationId) => {
  return await platformApi.get(`/organizations/${organizationId}`)
}

const create = async (data) => {
  return await platformApi.post('/organizations', data)
}

const update = async (organizationId, data) => {
  return await platformApi.put(`/organizations/${organizationId}`, data)
}

const del = async (organizationId) => {
  return await platformApi.delete(`/organizations/${organizationId}`)
}

const getRoles = async () => {
  return await platformApi.get('/organizations/roles')
}

const getMembers = async (
  organizationId,
  params?: { offset?: number; limit?: number; query?: string },
) => {
  const search = new URLSearchParams()
  if (params?.offset !== undefined) search.set('offset', String(params.offset))
  if (params?.limit !== undefined) search.set('limit', String(params.limit))
  if (params?.query) search.set('query', params.query)
  const qs = search.toString()
  return await platformApi.get(
    `/organizations/${organizationId}/users${qs ? `?${qs}` : ''}`,
  )
}

const inviteMember = async (organizationId, data) => {
  return await platformApi.post(`/organizations/${organizationId}/users`, data)
}

const updateMemberRole = async (organizationId, memberId, data) => {
  return await platformApi.put(
    `/organizations/${organizationId}/users/${memberId}/role`,
    data,
  )
}

const updateMemberInfo = async (organizationId, memberId, data) => {
  return await platformApi.put(
    `/organizations/${organizationId}/users/${memberId}/info`,
    data,
  )
}

const deleteMember = async (organizationId, memberId) => {
  return await platformApi.delete(
    `/organizations/${organizationId}/users/${memberId}`,
  )
}

const getOrganizationCredentials = async (
  organizationId,
  { offset = 0, limit = 10, query = '' } = {},
) => {
  const params = new URLSearchParams()
  params.set('offset', String(offset))
  params.set('limit', String(limit))
  if (query) params.set('query', query)
  return await platformApi.get(
    `/organizations/${organizationId}/credentials?${params.toString()}`,
  )
}

const getOrganizationCredentialsByNodes = async (
  organizationId,
  nodeIds: string[],
) => {
  const params = new URLSearchParams()
  for (const nodeId of nodeIds) {
    params.append('nodeIds', nodeId)
  }
  return await platformApi.get(
    `/organizations/${organizationId}/credentials/by-nodes?${params.toString()}`,
  )
}

const getOrganizationCredentialById = async (organizationId, credentialId) => {
  return await platformApi.get(
    `/organizations/${organizationId}/credentials/${credentialId}`,
  )
}

const createOrganizationCredential = async (organizationId, data) => {
  return await platformApi.post(
    `/organizations/${organizationId}/credentials`,
    data,
  )
}

const updateOrganizationCredential = async (
  organizationId,
  credentialId,
  data,
) => {
  return await platformApi.put(
    `/organizations/${organizationId}/credentials/${credentialId}`,
    data,
  )
}

const deleteOrganizationCredential = async (organizationId, credentialId) => {
  return await platformApi.delete(
    `/organizations/${organizationId}/credentials/${credentialId}`,
  )
}

const getMyMembership = async (organizationId) => {
  return await platformApi.get(`/organizations/${organizationId}/me`)
}

export default {
  getAll,
  getById,
  create,
  update,
  delete: del,
  getRoles,
  getMembers,
  inviteMember,
  updateMemberRole,
  updateMemberInfo,
  deleteMember,
  getOrganizationCredentials,
  getOrganizationCredentialsByNodes,
  getOrganizationCredentialById,
  createOrganizationCredential,
  updateOrganizationCredential,
  deleteOrganizationCredential,
  getMyMembership,
}
