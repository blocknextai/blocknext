import { useTranslation } from 'react-i18next'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
  TooltipProvider,
} from '@/components/ui/tooltip'
import { BadgeCheck } from 'lucide-react'

type VerifiedType = 'default' | 'user' | 'organization' | 'flow'

const typeMessages: Record<VerifiedType, string> = {
  default: 'ui.text.verified',
  user: 'ui.text.verifiedUser',
  organization: 'ui.text.verifiedOrganization',
  flow: 'ui.text.verifiedFlow',
}

interface VerifiedBadgeProps {
  type?: VerifiedType
  size?: number
  className?: string
}

export function VerifiedBadge({
  type = 'default',
  size = 20,
  className = '',
}: VerifiedBadgeProps) {
  const { t } = useTranslation()
  const message = t(typeMessages[type] || 'ui.text.verified')

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <BadgeCheck
            className={` cursor-pointer ${className}`}
            style={{ width: size, height: size }}
          />
        </TooltipTrigger>
        <TooltipContent>{message}</TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}
