import platformApi from '@/lib/platform-api'

const buildQuery = ({ offset = 0, limit = 10, query = '' } = {}) => {
  const params = new URLSearchParams()
  params.set('offset', String(offset))
  params.set('limit', String(limit))
  if (query) {
    params.set('query', query)
  }
  return params.toString()
}

// --- Organization scope (/organizations/:organizationId/notifications) ---
const getOrganizationNotifications = async (
  organizationId,
  pagination = {},
) => {
  return await platformApi.get(
    `/organizations/${organizationId}/notifications?${buildQuery(pagination)}`,
  )
}

const getOrganizationNotificationCounts = async (organizationId) => {
  return await platformApi.get(
    `/organizations/${organizationId}/notifications/count`,
  )
}

const markAllOrganizationNotificationsSeen = async (organizationId) => {
  return await platformApi.post(
    `/organizations/${organizationId}/notifications/seen`,
  )
}

const markAllOrganizationNotificationsRead = async (organizationId) => {
  return await platformApi.post(
    `/organizations/${organizationId}/notifications/read-all`,
  )
}

const markOrganizationNotificationRead = async (
  organizationId,
  recipientId,
) => {
  return await platformApi.patch(
    `/organizations/${organizationId}/notifications/${recipientId}/read`,
  )
}

const deleteOrganizationNotification = async (organizationId, recipientId) => {
  return await platformApi.delete(
    `/organizations/${organizationId}/notifications/${recipientId}`,
  )
}

// --- User scope (/users/me/notifications) ---
const getUserNotifications = async (pagination = {}) => {
  return await platformApi.get(
    `/users/me/notifications?${buildQuery(pagination)}`,
  )
}

const getUserNotificationCounts = async () => {
  return await platformApi.get(`/users/me/notifications/count`)
}

const markAllUserNotificationsSeen = async () => {
  return await platformApi.post(`/users/me/notifications/seen`)
}

const markAllUserNotificationsRead = async () => {
  return await platformApi.post(`/users/me/notifications/read-all`)
}

const markUserNotificationRead = async (recipientId) => {
  return await platformApi.patch(`/users/me/notifications/${recipientId}/read`)
}

const deleteUserNotification = async (recipientId) => {
  return await platformApi.delete(`/users/me/notifications/${recipientId}`)
}

export default {
  getOrganizationNotifications,
  getOrganizationNotificationCounts,
  markAllOrganizationNotificationsSeen,
  markAllOrganizationNotificationsRead,
  markOrganizationNotificationRead,
  deleteOrganizationNotification,
  getUserNotifications,
  getUserNotificationCounts,
  markAllUserNotificationsSeen,
  markAllUserNotificationsRead,
  markUserNotificationRead,
  deleteUserNotification,
}
