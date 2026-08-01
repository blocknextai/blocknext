import platformApi from '@/lib/platform-api'

const getMe = async () => {
  return await platformApi.get('/users/me')
}

const getNonce = async (data) => {
  return await platformApi.post('/auth/oauth/nonce', data)
}

const getToken = async (data) => {
  return await platformApi.post('/auth/oauth/token', data)
}

const refreshToken = async (token: string) => {
  return await platformApi.post('/auth/token/refresh', { refreshToken: token })
}

const register = async (data: { email: string; password: string }) => {
  return await platformApi.post('/auth/password/register', data)
}

const loginWithPassword = async (data: { email: string; password: string }) => {
  return await platformApi.post('/auth/password/login', data)
}

const verifyEmail = async (data: { token: string }) => {
  return await platformApi.post('/auth/email/verify', data)
}

const requestPasswordReset = async (data: { email: string }) => {
  return await platformApi.post('/auth/password/forgot', data)
}

const resetPassword = async (data: { token: string; newPassword: string }) => {
  return await platformApi.post('/auth/password/reset', data)
}

const setPassword = async (data: { password: string }) => {
  return await platformApi.post('/users/me/password/set', data)
}

const addEmail = async (data: { email: string }) => {
  return await platformApi.post('/users/me/email/add', data)
}

const changePassword = async (data: {
  currentPassword: string
  newPassword: string
}) => {
  return await platformApi.post('/users/me/password/change', data)
}

const changeEmail = async (data: {
  newEmail: string
  currentPassword: string
}) => {
  return await platformApi.post('/users/me/email/change', data)
}

const confirmEmailChange = async (data: { token: string }) => {
  return await platformApi.post('/auth/email/change/confirm', data)
}

const requestMagicLink = async (data: { email: string }) => {
  return await platformApi.post('/auth/email/magic-link/request', data)
}

const consumeMagicLink = async (data: { token: string }) => {
  return await platformApi.post('/auth/email/magic-link/consume', data)
}

const getMethods = async () => {
  return await platformApi.get('/auth/methods')
}

const getRoles = async () => {
  return await platformApi.get('/users/me/roles')
}

const getLinkedAccounts = async () => {
  return await platformApi.get('/users/me/linked-accounts')
}

const linkAccount = async (data) => {
  return await platformApi.post('/users/me/linked-accounts', data)
}

const unlinkAccount = async (accountId) => {
  return await platformApi.delete(`/users/me/linked-accounts/${accountId}`)
}

const resendEmailVerification = async (data: { email: string }) => {
  return await platformApi.post('/auth/email/resend-verification', data)
}

const getPreferences = async () => {
  return await platformApi.get('/users/me/preferences')
}

const updatePreferences = async (data) => {
  return await platformApi.patch('/users/me/preferences', data)
}

const getSessions = async (offset = 0, limit = 10) => {
  return await platformApi.get(
    `/users/me/sessions?offset=${offset}&limit=${limit}`,
  )
}

const revokeSession = async (sessionId) => {
  return await platformApi.post(`/users/me/sessions/${sessionId}/revoke`)
}

const revokeAllSessions = async () => {
  return await platformApi.post('/users/me/sessions/revoke-all')
}

const logout = async () => {
  return await platformApi.post('/auth/logout')
}

export default {
  getMe,
  getNonce,
  getToken,
  refreshToken,
  register,
  loginWithPassword,
  verifyEmail,
  requestPasswordReset,
  resetPassword,
  setPassword,
  addEmail,
  changePassword,
  changeEmail,
  confirmEmailChange,
  requestMagicLink,
  consumeMagicLink,
  getMethods,
  getRoles,
  getLinkedAccounts,
  linkAccount,
  unlinkAccount,
  resendEmailVerification,
  getPreferences,
  updatePreferences,
  getSessions,
  revokeSession,
  revokeAllSessions,
  logout,
}
