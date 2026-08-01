import { intlFormatter } from '@/lib/date'
import TimeAgo from 'react-timeago'
import { format } from 'date-fns'
import type { ComponentProps } from 'react'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'

type TimeAgoProps = ComponentProps<typeof TimeAgo>

export default function TimeAgoI18n(props: TimeAgoProps) {
  if (!props.date) {
    return null
  }
  const fullDate = format(props.date, 'dd/MM/yyyy HH:mm:ss')
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span>
          <TimeAgo
            {...props}
            formatter={intlFormatter() as TimeAgoProps['formatter']}
            title=""
          />
        </span>
      </TooltipTrigger>
      <TooltipContent>{fullDate}</TooltipContent>
    </Tooltip>
  )
}
