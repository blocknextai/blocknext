import platformApi from '@/lib/platform-api'

const getUserCredentials = async ({
  offset = 0,
  limit = 10,
  query = '',
} = {}) => {
  const params = new URLSearchParams()
  params.set('offset', String(offset))
  params.set('limit', String(limit))
  if (query) {
    params.set('query', query)
  }
  return await platformApi.get(`/users/me/credentials?${params.toString()}`)
}

const getUserCredentialsByNodes = async (nodeIds: string[]) => {
  const params = new URLSearchParams()
  for (const nodeId of nodeIds) {
    params.append('nodeIds', nodeId)
  }
  return await platformApi.get(
    `/users/me/credentials/by-nodes?${params.toString()}`,
  )
}

const getUserCredentialById = async (credentialId) => {
  return await platformApi.get(`/users/me/credentials/${credentialId}`)
}

const createUserCredential = async (data) => {
  return await platformApi.post('/users/me/credentials', data)
}

const updateUserCredential = async (credentialId, data) => {
  return await platformApi.put(`/users/me/credentials/${credentialId}`, data)
}

const deleteUserCredential = async (credentialId) => {
  return await platformApi.delete(`/users/me/credentials/${credentialId}`)
}

export default {
  getUserCredentials,
  getUserCredentialsByNodes,
  getUserCredentialById,
  createUserCredential,
  updateUserCredential,
  deleteUserCredential,
}
