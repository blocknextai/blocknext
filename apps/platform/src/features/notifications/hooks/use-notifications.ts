import { useCallback, useState } from 'react'
import useSWR, { useSWRConfig } from 'swr'
import notificationsService from '@/features/notifications/services/notifications'

const defaultPagination = {
  total: 0,
  page: 1,
  limit: 10,
  offset: 0,
  hasNext: false,
  hasPrev: false,
}

type Pagination = { offset?: number; limit?: number; query?: string }

export function useOrganizationNotifications(
  organizationId: string | null | undefined,
) {
  const [pagination, setPagination] = useState<Pagination>({
    offset: 0,
    limit: 10,
    query: '',
  })
  const { mutate } = useSWRConfig()

  const listKey = organizationId
    ? [
        'org-notifications',
        organizationId,
        pagination.offset,
        pagination.limit,
        pagination.query ?? '',
      ]
    : null

  const { data, error, isLoading } = useSWR(
    listKey,
    () =>
      notificationsService.getOrganizationNotifications(
        organizationId,
        pagination,
      ),
    { keepPreviousData: true },
  )

  const { data: countsData } = useSWR(
    organizationId ? ['org-notifications-counts', organizationId] : null,
    () =>
      notificationsService.getOrganizationNotificationCounts(organizationId),
  )

  // read / delete change both the list and the counts
  const revalidate = useCallback(
    () =>
      mutate(
        (k) =>
          Array.isArray(k) &&
          String(k[0]).startsWith('org-notifications') &&
          k[1] === organizationId,
        undefined,
        { revalidate: true },
      ),
    [organizationId, mutate],
  )

  // "seen" only affects the bell counter, not the list
  const revalidateCounts = useCallback(
    () =>
      mutate(
        (k) =>
          Array.isArray(k) &&
          k[0] === 'org-notifications-counts' &&
          k[1] === organizationId,
        undefined,
        { revalidate: true },
      ),
    [organizationId, mutate],
  )

  const fetchNotifications = useCallback(
    (offset = 0, limit = 10, query = '') => {
      setPagination({ offset, limit, query })
    },
    [],
  )

  const handleMarkAllAsSeen = useCallback(async () => {
    if (!organizationId) {
      return
    }
    await notificationsService.markAllOrganizationNotificationsSeen(
      organizationId,
    )
    await revalidateCounts()
  }, [organizationId, revalidateCounts])

  const handleMarkAsRead = useCallback(
    async (recipientId: string) => {
      if (!organizationId) {
        return
      }
      await notificationsService.markOrganizationNotificationRead(
        organizationId,
        recipientId,
      )
      await revalidate()
    },
    [organizationId, revalidate],
  )

  const handleMarkAllAsRead = useCallback(async () => {
    if (!organizationId) {
      return
    }
    await notificationsService.markAllOrganizationNotificationsRead(
      organizationId,
    )
    await revalidate()
  }, [organizationId, revalidate])

  const handleDelete = useCallback(
    async (recipientId: string) => {
      if (!organizationId) {
        return
      }
      await notificationsService.deleteOrganizationNotification(
        organizationId,
        recipientId,
      )
      await revalidate()
    },
    [organizationId, revalidate],
  )

  return {
    loading: isLoading,
    notifications: data?.data ?? [],
    pagination: data?.meta?.pagination ?? defaultPagination,
    unreadCount: countsData?.data?.unread ?? 0,
    unseenCount: countsData?.data?.unseen ?? 0,
    error,
    fetchNotifications,
    handleMarkAllAsSeen,
    handleMarkAsRead,
    handleMarkAllAsRead,
    handleDelete,
  }
}

const RECENT_LIMIT = 6

// Powers the header bell popover: the badge count is always live, while the
// recent list is only fetched once the popover is opened (`enabled`).
export function useOrganizationNotificationsBell(
  organizationId: string | null | undefined,
  enabled: boolean,
) {
  const { mutate } = useSWRConfig()

  const { data: countsData } = useSWR(
    organizationId ? ['org-notifications-counts', organizationId] : null,
    () =>
      notificationsService.getOrganizationNotificationCounts(organizationId),
  )

  const { data: listData, isLoading } = useSWR(
    enabled && organizationId
      ? ['org-notifications-recent', organizationId]
      : null,
    () =>
      notificationsService.getOrganizationNotifications(organizationId, {
        offset: 0,
        limit: RECENT_LIMIT,
      }),
    { keepPreviousData: true },
  )

  const revalidate = useCallback(
    () =>
      mutate(
        (k) =>
          Array.isArray(k) &&
          String(k[0]).startsWith('org-notifications') &&
          k[1] === organizationId,
        undefined,
        { revalidate: true },
      ),
    [organizationId, mutate],
  )

  const revalidateCounts = useCallback(
    () =>
      mutate(
        (k) =>
          Array.isArray(k) &&
          k[0] === 'org-notifications-counts' &&
          k[1] === organizationId,
        undefined,
        { revalidate: true },
      ),
    [organizationId, mutate],
  )

  const handleMarkAllAsSeen = useCallback(async () => {
    if (!organizationId) {
      return
    }
    await notificationsService.markAllOrganizationNotificationsSeen(
      organizationId,
    )
    await revalidateCounts()
  }, [organizationId, revalidateCounts])

  const handleMarkAsRead = useCallback(
    async (recipientId: string) => {
      if (!organizationId) {
        return
      }
      await notificationsService.markOrganizationNotificationRead(
        organizationId,
        recipientId,
      )
      await revalidate()
    },
    [organizationId, revalidate],
  )

  const handleMarkAllAsRead = useCallback(async () => {
    if (!organizationId) {
      return
    }
    await notificationsService.markAllOrganizationNotificationsRead(
      organizationId,
    )
    await revalidate()
  }, [organizationId, revalidate])

  return {
    loading: isLoading,
    notifications: listData?.data ?? [],
    unreadCount: countsData?.data?.unread ?? 0,
    unseenCount: countsData?.data?.unseen ?? 0,
    handleMarkAllAsSeen,
    handleMarkAsRead,
    handleMarkAllAsRead,
  }
}

export function useUserNotificationsCount() {
  const { data: countsData } = useSWR(['user-notifications-counts'], () =>
    notificationsService.getUserNotificationCounts(),
  )

  return {
    unreadCount: countsData?.data?.unread ?? 0,
    unseenCount: countsData?.data?.unseen ?? 0,
  }
}

export function useUserNotifications() {
  const [pagination, setPagination] = useState<Pagination>({
    offset: 0,
    limit: 10,
    query: '',
  })
  const { mutate } = useSWRConfig()

  const listKey = [
    'user-notifications',
    pagination.offset,
    pagination.limit,
    pagination.query ?? '',
  ]

  const { data, error, isLoading } = useSWR(
    listKey,
    () => notificationsService.getUserNotifications(pagination),
    { keepPreviousData: true },
  )

  const { data: countsData } = useSWR(['user-notifications-counts'], () =>
    notificationsService.getUserNotificationCounts(),
  )

  const revalidate = useCallback(
    () =>
      mutate(
        (k) =>
          Array.isArray(k) && String(k[0]).startsWith('user-notifications'),
        undefined,
        { revalidate: true },
      ),
    [mutate],
  )

  const revalidateCounts = useCallback(
    () =>
      mutate(
        (k) => Array.isArray(k) && k[0] === 'user-notifications-counts',
        undefined,
        { revalidate: true },
      ),
    [mutate],
  )

  const fetchNotifications = useCallback(
    (offset = 0, limit = 10, query = '') => {
      setPagination({ offset, limit, query })
    },
    [],
  )

  const handleMarkAllAsSeen = useCallback(async () => {
    await notificationsService.markAllUserNotificationsSeen()
    await revalidateCounts()
  }, [revalidateCounts])

  const handleMarkAsRead = useCallback(
    async (recipientId: string) => {
      await notificationsService.markUserNotificationRead(recipientId)
      await revalidate()
    },
    [revalidate],
  )

  const handleMarkAllAsRead = useCallback(async () => {
    await notificationsService.markAllUserNotificationsRead()
    await revalidate()
  }, [revalidate])

  const handleDelete = useCallback(
    async (recipientId: string) => {
      await notificationsService.deleteUserNotification(recipientId)
      await revalidate()
    },
    [revalidate],
  )

  return {
    loading: isLoading,
    notifications: data?.data ?? [],
    pagination: data?.meta?.pagination ?? defaultPagination,
    unreadCount: countsData?.data?.unread ?? 0,
    unseenCount: countsData?.data?.unseen ?? 0,
    error,
    fetchNotifications,
    handleMarkAllAsSeen,
    handleMarkAsRead,
    handleMarkAllAsRead,
    handleDelete,
  }
}
