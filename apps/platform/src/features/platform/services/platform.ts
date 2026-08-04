import platformApi from '@/lib/platform-api'

const getFeatures = async () => {
  return await platformApi.get('/platform/features')
}

export default {
  getFeatures,
}
