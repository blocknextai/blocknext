import type { ReactNode } from 'react'
import { Badge } from '@/components/ui/badge'
import { CheckCircle2, XCircle, Ban, Clock } from 'lucide-react'
import { useTranslation } from 'react-i18next'

type StatusKey = 'pending' | 'running' | 'success' | 'failed' | 'cancelled'

type StatusEntry = {
  icon: ReactNode
  color: string
  bgColor: string
  textColor: string
}

const statusConfig: Record<StatusKey, StatusEntry> = {
  pending: {
    icon: <Clock className="size-3" />,
    color: 'text-zinc-500',
    bgColor: 'bg-zinc-100 dark:bg-zinc-800',
    textColor: 'text-zinc-800 dark:text-zinc-200',
  },
  running: {
    icon: <span className="loading loading-ring loading-xs"></span>,
    color: 'text-amber-500',
    bgColor: 'bg-amber-100 dark:bg-amber-900',
    textColor: 'text-amber-800 dark:text-amber-200',
  },
  success: {
    icon: <CheckCircle2 className="size-3" />,
    color: 'text-emerald-500',
    bgColor: 'bg-emerald-100 dark:bg-emerald-900',
    textColor: 'text-emerald-800 dark:text-emerald-200',
  },
  failed: {
    icon: <XCircle className="size-3" />,
    color: 'text-destructive',
    bgColor: 'bg-destructive/10',
    textColor: 'text-destructive',
  },
  cancelled: {
    icon: <Ban className="size-3" />,
    color: 'text-zinc-500',
    bgColor: 'bg-zinc-100 dark:bg-zinc-800',
    textColor: 'text-zinc-800 dark:text-zinc-200',
  },
}

type BadgeVariant = 'default' | 'outline' | 'secondary' | 'destructive'

interface StatusBadgeProps {
  status: string
  variant?: BadgeVariant
  className?: string
  title?: string
}

const StatusBadge = ({
  status,
  variant = 'outline',
  className = '',
  title,
}: StatusBadgeProps) => {
  const { t } = useTranslation()
  const config =
    (statusConfig[status as StatusKey] as StatusEntry) ?? statusConfig.pending

  return (
    <Badge
      variant={variant}
      className={`${config.bgColor} ${config.textColor} ${className}`}
      title={title}
    >
      <span className="mr-1">{config.icon}</span>
      <span className="capitalize">
        {String(t(`ui.text.status.${status}`, status))}
      </span>
    </Badge>
  )
}

export default StatusBadge
