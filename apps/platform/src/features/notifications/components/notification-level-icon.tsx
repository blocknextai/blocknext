import { Bell, Info, CheckCircle2, AlertTriangle, XCircle } from 'lucide-react'
import { cn } from '@/lib/utils'

type NotificationLevel = 'info' | 'success' | 'warning' | 'error'

const levelConfig: Record<
  NotificationLevel,
  { icon: typeof Info; className: string }
> = {
  info: { icon: Info, className: 'text-sky-500' },
  success: { icon: CheckCircle2, className: 'text-emerald-500' },
  warning: { icon: AlertTriangle, className: 'text-amber-500' },
  error: { icon: XCircle, className: 'text-destructive' },
}

export function NotificationLevelIcon({
  level,
  className,
}: {
  level?: string
  className?: string
}) {
  const config = levelConfig[level as NotificationLevel] ?? {
    icon: Bell,
    className: 'text-muted-foreground',
  }
  const Icon = config.icon
  return <Icon className={cn('size-4 shrink-0', config.className, className)} />
}

export const isExternalUrl = (url: string) => /^https?:\/\//i.test(url)
