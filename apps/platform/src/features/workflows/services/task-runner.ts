import platformApi from '@/lib/platform-api'

const trigger = async (organizationId, data) => {
  return await platformApi.post(
    `/organizations/${organizationId}/task-runner/trigger`,
    data,
  )
}

const cancel = async (organizationId, data) => {
  return await platformApi.post(
    `/organizations/${organizationId}/task-runner/cancel`,
    data,
  )
}

const rerunAll = async (organizationId, data) => {
  return await platformApi.post(
    `/organizations/${organizationId}/task-runner/rerun-all`,
    data,
  )
}

const rerunFailed = async (organizationId, data) => {
  return await platformApi.post(
    `/organizations/${organizationId}/task-runner/rerun-failed`,
    data,
  )
}

export default {
  trigger,
  cancel,
  rerunAll,
  rerunFailed,
}
