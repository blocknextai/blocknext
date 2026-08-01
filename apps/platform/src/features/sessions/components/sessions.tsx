import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { Monitor, Smartphone, Globe, Loader2, LogOut } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'
import { Badge } from '@/components/ui/badge'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Loading } from '@/components/shared/loading'
import TimeAgoI18n from '@/components/shared/timeagoi18'
import { AppPagination } from '@/components/shared/app-pagination'
import { ProviderIcon } from '@/features/auth/components/provider-icon'
import { getProviderName } from '@/features/auth/components/provider-utils'

function parseDevice(userAgent: string, unknownLabel: string) {
  if (!userAgent) return { label: unknownLabel, isMobile: false }

  const ua = userAgent.toLowerCase()

  let browser = 'Browser'
  if (ua.includes('firefox')) browser = 'Firefox'
  else if (ua.includes('edg')) browser = 'Edge'
  else if (ua.includes('chrome') && !ua.includes('edg')) browser = 'Chrome'
  else if (ua.includes('safari') && !ua.includes('chrome')) browser = 'Safari'

  let os = ''
  if (ua.includes('windows')) os = 'Windows'
  else if (ua.includes('mac os') || ua.includes('macintosh')) os = 'macOS'
  else if (ua.includes('linux') && !ua.includes('android')) os = 'Linux'
  else if (ua.includes('android')) os = 'Android'
  else if (ua.includes('iphone') || ua.includes('ipad')) os = 'iOS'

  const isMobile =
    ua.includes('mobile') || ua.includes('android') || ua.includes('iphone')

  const label = os ? `${browser} (${os})` : browser

  return { label, isMobile }
}

const Sessions = ({
  sessions,
  loading,
  revokingSessionId,
  revokingAll,
  onRevokeSession,
  onRevokeAllSessions,
  pagination,
  page,
  limit,
  offset,
  onPageChange,
}) => {
  const { t } = useTranslation()
  const otherSessions = useMemo(
    () => (sessions ?? []).filter((s) => !s.isCurrent),
    [sessions],
  )

  if (loading) {
    return <Loading />
  }

  if (!sessions || sessions.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-muted-foreground">
        <Globe className="w-12 h-12 mb-4 opacity-40" />
        <p className="text-lg font-medium">{t('ui.text.noActiveSessions')}</p>
        <p className="text-sm">{t('ui.text.sessionDataAppearAfterLogin')}</p>
      </div>
    )
  }

  return (
    <div className="space-y-6 p-4">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-lg font-semibold text-foreground">
            {t('ui.text.activeSessions')}
          </h2>
          <p className="text-sm text-muted-foreground">
            {t('ui.text.manageActiveSessions')}
          </p>
        </div>
        {otherSessions.length > 0 && (
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant="destructive" size="sm" disabled={revokingAll}>
                {revokingAll ? (
                  <Loader2 className="w-4 h-4 animate-spin mr-2" />
                ) : (
                  <LogOut className="w-4 h-4 mr-2" />
                )}
                {t('ui.text.revokeAllOtherSessions')}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>
                  {t('ui.text.revokeAllSessionsConfirm')}
                </AlertDialogTitle>
                <AlertDialogDescription>
                  {t('ui.text.revokeAllSessionsDescription')}
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>{t('ui.text.cancel')}</AlertDialogCancel>
                <AlertDialogAction onClick={onRevokeAllSessions}>
                  {t('ui.text.revokeAll')}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        )}
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{t('ui.text.device')}</TableHead>
            <TableHead>{t('ui.text.signedInVia')}</TableHead>
            <TableHead>{t('ui.text.created')}</TableHead>
            <TableHead>{t('ui.text.lastActive')}</TableHead>
            <TableHead className="text-right">{t('ui.text.actions')}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {sessions.map((session) => {
            const device = parseDevice(session.userAgent, t('ui.text.unknown'))
            const DeviceIcon = device.isMobile ? Smartphone : Monitor

            return (
              <TableRow key={session.sessionId}>
                <TableCell>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <div className="flex items-center gap-2">
                        <DeviceIcon className="w-4 h-4 text-muted-foreground" />
                        <span className="truncate max-w-[200px]">
                          {device.label}
                        </span>
                        {session.isCurrent && (
                          <Badge
                            variant="default"
                            className="bg-green-600/15 text-green-600 border-green-600/20 ml-1"
                          >
                            {t('ui.text.currentSession')}
                          </Badge>
                        )}
                      </div>
                    </TooltipTrigger>
                    <TooltipContent side="top">
                      {session.userAgent || 'Unknown'}
                    </TooltipContent>
                  </Tooltip>
                </TableCell>
                <TableCell>
                  {session.authProvider ? (
                    <Badge variant="outline" className="gap-1.5">
                      <ProviderIcon
                        provider={session.authProvider}
                        className="w-3 h-3"
                      />
                      {getProviderName(session.authProvider)}
                    </Badge>
                  ) : (
                    <span className="text-sm text-muted-foreground">—</span>
                  )}
                </TableCell>
                <TableCell>
                  <TimeAgoI18n date={session.createdAt} />
                </TableCell>
                <TableCell>
                  <TimeAgoI18n date={session.updatedAt} />
                </TableCell>
                <TableCell className="text-right">
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <Button
                        variant="ghost"
                        size="sm"
                        disabled={revokingSessionId === session.sessionId}
                      >
                        {revokingSessionId === session.sessionId ? (
                          <Loader2 className="w-4 h-4 animate-spin" />
                        ) : (
                          <LogOut className="w-4 h-4" />
                        )}
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>
                          {session.isCurrent
                            ? t('ui.text.revokeCurrentSession')
                            : t('ui.text.revokeThisSession')}
                        </AlertDialogTitle>
                        <AlertDialogDescription>
                          {session.isCurrent
                            ? t('ui.text.revokeCurrentSessionDescription')
                            : t('ui.text.revokeThisSessionDescription', {
                                device: device.label,
                              })}
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>
                          {t('ui.text.cancel')}
                        </AlertDialogCancel>
                        <AlertDialogAction
                          onClick={() => onRevokeSession(session.sessionId)}
                        >
                          {t('ui.text.revoke')}
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>

      {pagination && (
        <AppPagination
          pagination={{
            total: pagination.total,
            page,
            limit,
            offset,
            hasPrev: pagination.hasPrev,
            hasNext: pagination.hasNext,
          }}
          onPageChange={onPageChange}
        />
      )}
    </div>
  )
}

export default Sessions
