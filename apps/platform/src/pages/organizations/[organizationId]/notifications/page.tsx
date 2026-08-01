import { useParams } from 'react-router'
import { useOrganizationNotifications } from '@/features/notifications'
import NotificationsTab from '@/features/notifications/components/notifications-tab'

function OrganizationNotificationsPage() {
  const { organizationId } = useParams()

  const {
    loading,
    notifications,
    pagination,
    unreadCount,
    fetchNotifications,
    handleMarkAllAsSeen,
    handleMarkAsRead,
    handleMarkAllAsRead,
    handleDelete,
  } = useOrganizationNotifications(organizationId)

  return (
    <NotificationsTab
      loading={loading}
      notifications={notifications}
      pagination={pagination}
      unreadCount={unreadCount}
      onFetch={fetchNotifications}
      onMarkAllAsSeen={handleMarkAllAsSeen}
      onMarkAsRead={handleMarkAsRead}
      onMarkAllAsRead={handleMarkAllAsRead}
      onDelete={handleDelete}
    />
  )
}

export default OrganizationNotificationsPage
