import { useState } from 'react'
import { Link } from 'react-router'
import { Bell, CheckCheck } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import TimeAgoI18n from '@/components/shared/timeagoi18'
import { cn } from '@/lib/utils'
import { useOrganizationStore } from '@/stores/organization'
import { useOrganizationNotificationsBell } from '@/features/notifications'
import {
  NotificationLevelIcon,
  isExternalUrl,
} from '@/features/notifications/components/notification-level-icon'

function BellRow({
  notification,
  onActivate,
}: {
  notification: any
  onActivate: (notification: any) => void
}) {
  const isRead = !!notification.readAt
  const actionUrl: string | undefined = notification.actionUrl
  const external = actionUrl && isExternalUrl(actionUrl)

  const inner = (
    <>
      <NotificationLevelIcon level={notification.level} className="mt-0.5" />
      <div className="flex-1 min-w-0 text-left">
        <div className="flex items-center gap-2">
          {!isRead && (
            <span className="size-1.5 rounded-full bg-primary shrink-0" />
          )}
          <span
            className={cn(
              'text-sm truncate',
              isRead ? 'font-normal' : 'font-semibold',
            )}
          >
            {notification.title}
          </span>
        </div>
        {notification.body && (
          <p className="text-xs text-muted-foreground line-clamp-1 mt-0.5">
            {notification.body}
          </p>
        )}
        <span className="text-[11px] text-muted-foreground/70 mt-0.5 block">
          <TimeAgoI18n date={notification.createdAt} />
        </span>
      </div>
    </>
  )

  const className = cn(
    'flex items-start gap-2.5 rounded-md p-2 w-full transition-colors hover:bg-muted text-left',
    !isRead && 'bg-muted/40',
  )

  const handleClick = () => onActivate(notification)

  if (actionUrl && external) {
    return (
      <a
        href={actionUrl}
        target="_blank"
        rel="noreferrer"
        onClick={handleClick}
        className={className}
      >
        {inner}
      </a>
    )
  }

  if (actionUrl) {
    return (
      <Link to={actionUrl} onClick={handleClick} className={className}>
        {inner}
      </Link>
    )
  }

  return (
    <button type="button" onClick={handleClick} className={className}>
      {inner}
    </button>
  )
}

export function NotificationsBell() {
  const { t } = useTranslation()
  const organizationId = useOrganizationStore((s) => s.organizationId)
  const [open, setOpen] = useState(false)

  const {
    loading,
    notifications,
    unreadCount,
    unseenCount,
    handleMarkAllAsSeen,
    handleMarkAsRead,
    handleMarkAllAsRead,
  } = useOrganizationNotificationsBell(organizationId, open)

  if (!organizationId) {
    return null
  }

  const handleOpenChange = (next: boolean) => {
    setOpen(next)
    // Opening the panel clears the bell badge (marks everything seen).
    if (next && unseenCount > 0) {
      handleMarkAllAsSeen()
    }
  }

  const handleActivate = (notification: any) => {
    if (!notification.readAt) {
      handleMarkAsRead(notification.id)
    }
    setOpen(false)
  }

  return (
    <Popover open={open} onOpenChange={handleOpenChange}>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="relative cursor-pointer"
          aria-label={t('ui.text.notifications')}
        >
          <Bell className="size-4" />
          {unseenCount > 0 && (
            <span className="absolute -top-0.5 -right-0.5 flex items-center justify-center min-w-4 h-4 px-1 rounded-full bg-primary text-primary-foreground text-[10px] font-medium leading-none">
              {unseenCount > 99 ? '99+' : unseenCount}
            </span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent align="end" className="w-96 p-0 gap-0">
        <div className="flex items-center justify-between gap-2 p-3 border-b">
          <span className="font-semibold text-sm flex items-center gap-2">
            {t('ui.text.notifications')}
            {unreadCount > 0 && (
              <Badge variant="secondary">{unreadCount}</Badge>
            )}
          </span>
          {unreadCount > 0 && (
            <Button
              variant="ghost"
              size="xs"
              className="text-xs"
              onClick={() => handleMarkAllAsRead()}
            >
              <CheckCheck className="size-3.5" />
              {t('ui.text.markAllAsRead')}
            </Button>
          )}
        </div>

        <div className="max-h-96 overflow-y-auto p-1.5">
          {loading && !notifications.length ? (
            <div className="flex items-center justify-center py-10">
              <span className="loading loading-infinity loading-md" />
            </div>
          ) : !notifications.length ? (
            <div className="flex flex-col items-center justify-center gap-2 py-10 text-center">
              <Bell className="size-6 text-muted-foreground/70" />
              <span className="text-sm text-muted-foreground">
                {t('ui.text.noNotifications')}
              </span>
            </div>
          ) : (
            <div className="flex flex-col gap-0.5">
              {notifications.map((notification: any) => (
                <BellRow
                  key={notification.id}
                  notification={notification}
                  onActivate={handleActivate}
                />
              ))}
            </div>
          )}
        </div>

        <div className="border-t p-1.5">
          <Button
            asChild
            variant="ghost"
            className="w-full cursor-pointer"
            onClick={() => setOpen(false)}
          >
            <Link to={`/organizations/${organizationId}/notifications`}>
              {t('ui.text.viewAll')}
            </Link>
          </Button>
        </div>
      </PopoverContent>
    </Popover>
  )
}

NotificationsBell.displayName = 'NotificationsBell'
