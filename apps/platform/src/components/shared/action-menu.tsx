import type { ComponentType, MouseEvent, ReactNode } from 'react'
import { Link } from 'react-router'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { EllipsisVertical } from 'lucide-react'

type ButtonVariant =
  'default' | 'outline' | 'secondary' | 'ghost' | 'destructive' | 'link'
type ButtonSize =
  'default' | 'xs' | 'sm' | 'lg' | 'icon' | 'icon-xs' | 'icon-sm' | 'icon-lg'

export interface ActionMenuItem {
  label: ReactNode
  icon?: ReactNode
  href?: string
  onClick?: (event: MouseEvent) => void
  variant?: 'default' | 'destructive'
}

interface ActionMenuProps {
  items?: ActionMenuItem[]
  triggerIcon?: ComponentType<{ className?: string }>
  triggerVariant?: ButtonVariant
  triggerSize?: ButtonSize
  triggerClassName?: string
}

const ActionMenu = ({
  items = [],
  triggerIcon: TriggerIcon = EllipsisVertical,
  triggerVariant = 'ghost',
  triggerSize = 'sm',
  triggerClassName = '',
}: ActionMenuProps) => {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant={triggerVariant}
          size={triggerSize}
          className={`cursor-pointer ${triggerClassName}`}
        >
          <TriggerIcon className="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        {items.map((item, index) => {
          if (item.href) {
            return (
              <DropdownMenuItem key={index} asChild>
                <Link to={item.href}>
                  {item.icon} {item.label}
                </Link>
              </DropdownMenuItem>
            )
          }

          return (
            <DropdownMenuItem
              key={index}
              onClick={item.onClick}
              variant={item.variant}
            >
              {item.icon} {item.label}
            </DropdownMenuItem>
          )
        })}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

export default ActionMenu
