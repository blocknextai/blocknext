import { useUserNotifications } from '@/features/notifications'
import NotificationsTab from '@/features/notifications/components/notifications-tab'

function PreferencesNotificationsPage() {
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
  } = useUserNotifications()

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

export default PreferencesNotificationsPage
