import { useCallback, useEffect, useState } from 'react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { Loading } from '@/components/shared/loading'
import { Separator } from '@/components/ui/separator'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import ConfirmationDialog from '@/components/shared/confirmation-dialog'
import ActionMenu from '@/components/shared/action-menu'
import TimeAgoI18n from '@/components/shared/timeagoi18'
import { AppPagination } from '@/components/shared/app-pagination'
import { SearchInput } from '@/components/shared/search-input'
import { SearchEmptyState } from '@/components/shared/search-empty-state'
import {
  NotificationLevelIcon,
  isExternalUrl,
} from '@/features/notifications/components/notification-level-icon'
import { cn } from '@/lib/utils'
import { Bell, CheckCheck, Check, Trash2, ExternalLink } from 'lucide-react'

const NotificationsTab = ({
  loading,
  notifications,
  pagination,
  unreadCount,
  onFetch,
  onMarkAllAsSeen,
  onMarkAsRead,
  onMarkAllAsRead,
  onDelete,
}) => {
  const { t } = useTranslation()
  const [query, setQuery] = useState('')
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [confirmData, setConfirmData] = useState<any>(null)

  useEffect(() => {
    onFetch()
    // Opening the page clears the bell badge (marks everything seen).
    onMarkAllAsSeen?.()
  }, [])

  const handleSearch = useCallback(
    (value: string) => {
      setQuery(value)
      onFetch(0, pagination?.limit ?? 10, value)
    },
    [onFetch, pagination?.limit],
  )

  const handlePageChange = (newPage: number) => {
    const newOffset = (newPage - 1) * pagination.limit
    onFetch(newOffset, pagination.limit, query)
  }

  const confirmDelete = (notification: any) => {
    setConfirmData({
      description: t('ui.text.deleteNotificationConfirmation'),
      action: () => onDelete(notification.id),
      label: t('ui.text.delete'),
    })
    setConfirmOpen(true)
  }

  return (
    <div className="w-full h-full flex flex-col">
      <div className="flex flex-col gap-2 mb-4 p-6 shrink-0">
        <div className="flex justify-between items-center">
          <h2 className="text-2xl font-bold flex items-center gap-2">
            {t('ui.text.notifications')}
            {unreadCount > 0 && (
              <Badge variant="secondary">{unreadCount}</Badge>
            )}
          </h2>
          {unreadCount > 0 && (
            <Button
              size="sm"
              variant="outline"
              onClick={() => onMarkAllAsRead()}
            >
              <CheckCheck />
              <span>{t('ui.text.markAllAsRead')}</span>
            </Button>
          )}
        </div>
        <Separator className="my-2 p-0" decorative={true} />
        <div className="max-w-sm">
          <SearchInput
            placeholder={t('ui.text.searchNotifications')}
            onSearch={handleSearch}
          />
        </div>
      </div>

      <div className="flex-1 min-h-0">
        {loading && !pagination?.total ? (
          <Loading />
        ) : !notifications?.length && query ? (
          <SearchEmptyState />
        ) : !notifications?.length ? (
          <div className="flex w-full h-full flex-col items-center justify-center text-center">
            <div className="flex flex-col items-center gap-3">
              <div className="p-3 rounded-full w-12 h-12 flex items-center justify-center">
                <Bell className="size-6 text-muted-foreground/80" />
              </div>
              <div className="flex flex-col gap-2">
                <span className="text-md font-semibold">
                  {t('ui.text.noNotifications')}
                </span>
                <span className="text-muted-foreground text-xs max-w-xs">
                  {t('ui.text.noNotificationsDescription')}
                </span>
              </div>
            </div>
          </div>
        ) : (
          <div className="px-6 pb-6">
            <ul className="flex flex-col gap-1">
              {notifications.map((notification: any) => {
                const isRead = !!notification.readAt
                const actionUrl: string | undefined = notification.actionUrl
                const external = actionUrl && isExternalUrl(actionUrl)
                const onActionClick = () => {
                  if (!isRead) {
                    onMarkAsRead(notification.id)
                  }
                }

                const titleEl = (
                  <span
                    className={cn(
                      'text-sm truncate',
                      isRead ? 'font-normal' : 'font-semibold',
                      actionUrl &&
                        'hover:underline underline-offset-2 cursor-pointer',
                    )}
                  >
                    {notification.title}
                  </span>
                )

                const titleNode = !actionUrl ? (
                  titleEl
                ) : external ? (
                  <a
                    href={actionUrl}
                    target="_blank"
                    rel="noreferrer"
                    onClick={onActionClick}
                    className="inline-flex items-center gap-1 min-w-0"
                  >
                    {titleEl}
                    <ExternalLink className="size-3 shrink-0 text-muted-foreground" />
                  </a>
                ) : (
                  <Link
                    to={actionUrl}
                    onClick={onActionClick}
                    className="min-w-0"
                  >
                    {titleEl}
                  </Link>
                )

                return (
                  <li
                    key={notification.id}
                    className={cn(
                      'flex items-start gap-3 rounded-md border p-3 transition-colors',
                      isRead ? 'bg-card' : 'bg-muted/50 border-primary/20',
                    )}
                  >
                    <div className="mt-0.5">
                      <NotificationLevelIcon level={notification.level} />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        {!isRead && (
                          <span className="size-2 rounded-full bg-primary shrink-0" />
                        )}
                        {titleNode}
                      </div>
                      {notification.body && (
                        <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                          {notification.body}
                        </p>
                      )}
                      <span className="text-xs text-muted-foreground/70 mt-1 block">
                        <TimeAgoI18n date={notification.createdAt} />
                      </span>
                    </div>
                    <ActionMenu
                      items={[
                        ...(isRead
                          ? []
                          : [
                              {
                                label: t('ui.text.markAsRead'),
                                icon: <Check className="size-4" />,
                                onClick: () => onMarkAsRead(notification.id),
                              },
                            ]),
                        {
                          label: t('ui.text.delete'),
                          icon: <Trash2 className="size-4" />,
                          variant: 'destructive',
                          onClick: () => confirmDelete(notification),
                        },
                      ]}
                    />
                  </li>
                )
              })}
            </ul>
            <AppPagination
              pagination={pagination}
              onPageChange={handlePageChange}
            />
          </div>
        )}
      </div>

      <ConfirmationDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title={t('ui.text.areYouSure')}
        description={confirmData?.description}
        confirmText={confirmData?.label}
        onConfirm={confirmData?.action}
        variant="destructive"
      />
    </div>
  )
}

NotificationsTab.displayName = 'NotificationsTab'

export default NotificationsTab
