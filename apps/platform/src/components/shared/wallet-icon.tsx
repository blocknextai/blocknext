import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { UserIcon2 } from '@/components/shared/custom-icons'

interface LinkedAccount {
  isPrimary?: boolean
  displayName?: string
}

interface WalletIconOwner {
  id?: string
  alias?: string
  linkedAccounts?: LinkedAccount[]
}

interface WalletIconProps {
  owner?: WalletIconOwner | null
}

const WalletIcon = ({ owner }: WalletIconProps) => {
  if (!owner) {
    return null
  }

  const alias = owner.alias || ''
  const displayName =
    owner.linkedAccounts?.find((account) => account?.isPrimary)?.displayName ||
    alias
  const seed = owner.id || ''

  return (
    <Tooltip>
      <TooltipTrigger className="flex items-center gap-2">
        <div className="size-5 overflow-hidden">
          <UserIcon2 seed={seed} className={'size-4 text-primary'} size={20} />
        </div>
        <span className="text-muted-foreground text-sm">{alias}</span>
      </TooltipTrigger>
      <TooltipContent
        side="top"
        className="max-w-md whitespace-pre-wrap text-center"
      >
        {displayName}
      </TooltipContent>
    </Tooltip>
  )
}
WalletIcon.displayName = 'WalletIcon'

export { WalletIcon }
