import platformApi from '@/lib/platform-api'

const getCredentials = async () => {
  return await platformApi.get('/node-engine/credentials')
}

const getCredentialById = async (id: string) => {
  return await platformApi.get(`/node-engine/credentials/${id}`)
}

const getNodes = async () => {
  return await platformApi.get('/node-engine/nodes')
}

const getWebhookSources = async () => {
  return await platformApi.get('/node-engine/webhook-sources')
}

const getTriggerVariables = async () => {
  return await platformApi.get('/node-engine/trigger-variables')
}

export default {
  getCredentials,
  getCredentialById,
  getNodes,
  getWebhookSources,
  getTriggerVariables,
}
