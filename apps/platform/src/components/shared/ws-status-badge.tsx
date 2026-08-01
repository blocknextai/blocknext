import { Loader2, Wifi, WifiOff } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { useWsStatus } from '@/hooks/use-ws-status'

const WsStatusBadge = ({ className = '' }: { className?: string }) => {
  const { t } = useTranslation()
  const status = useWsStatus()

  if (status === 'online') {
    return (
      <Badge
        variant="outline"
        className={`gap-1 bg-green-500/10 text-green-600 border-green-500/20 ${className}`}
      >
        <Wifi className="size-3" />
        {t('ui.text.liveOnline', 'Live')}
      </Badge>
    )
  }

  if (status === 'connecting' || status === 'reconnecting') {
    return (
      <Badge
        variant="outline"
        className={`gap-1 bg-amber-500/10 text-amber-600 border-amber-500/20 ${className}`}
      >
        <Loader2 className="size-3 animate-spin" />
        {t('ui.text.liveReconnecting', 'Reconnecting')}
      </Badge>
    )
  }

  return (
    <Badge
      variant="outline"
      className={`gap-1 bg-zinc-500/10 text-zinc-600 border-zinc-500/20 ${className}`}
    >
      <WifiOff className="size-3" />
      {t('ui.text.liveOffline', 'Offline')}
    </Badge>
  )
}

WsStatusBadge.displayName = 'WsStatusBadge'

export { WsStatusBadge }
