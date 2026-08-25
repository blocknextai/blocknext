import platformApi from '@/lib/platform-api'

const authorizeOrganization = async (
  organizationId: string,
  data: { credentialId: string },
) => {
  return await platformApi.post(
    `/organizations/${organizationId}/credential-oauth/oauth2/auth`,
    data,
  )
}

const callback = async (queryString) => {
  return await platformApi.get(
    `/credential-oauth/oauth2/callback${queryString}`,
  )
}

export default {
  authorizeOrganization,
  callback,
}
