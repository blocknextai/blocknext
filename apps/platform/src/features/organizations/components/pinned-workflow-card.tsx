import { FlowIcon } from '@/components/shared/custom-icons'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '@/components/ui/button'
import { Copy, EllipsisVertical, Pencil, Pin, Play, Trash2 } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { VerifiedBadge } from '@/components/shared/verified-tooltip'
import { CopyIdMenuItem } from '@/components/shared/copy-id-menu-item'

const PinnedWorkflowCard = ({
  workflow,
  handleFavorite,
  handleDuplicate,
  handleDelete,
}) => {
  const { t } = useTranslation()

  return (
    <div className="group bg-card p-3 rounded-lg dark:border-transparent justify-between w-full h-full overflow-hidden relative transition-all duration-300 ease-in-out hover:scale-[1.03] shadow-sm hover:shadow-lg hover:z-10">
      <div className="flex justify-between items-center gap-3">
        <div className="flex gap-4 items-center">
          <div
            className="p-7 rounded-md"
            style={{
              backgroundImage: `${FlowIcon(workflow.id)}`,
              backgroundSize: 'cover',
              backgroundPosition: 'center',
            }}
          ></div>
          <div className="flex flex-col gap-1">
            <div
              className="heading-xs whitespace-normal"
              title={workflow.title}
            >
              {workflow.title.length > 15
                ? `${workflow.title.slice(0, 15)}...`
                : workflow.title}
            </div>
            <div className="flex items-center gap-1.5 py-1 text-sm text-muted-foreground">
              <span className="truncate">
                {t('ui.text.createdBy')}{' '}
                <span className="text-primary underline">
                  {workflow.owner?.alias}
                </span>
              </span>
              {workflow.owner?.isVerified && <VerifiedBadge type="user" />}
            </div>
          </div>
        </div>
        <div className="flex gap-2">
          <Link
            className="flex-1 cursor-pointer!"
            to={`/organizations/${workflow.organizationId}/run/${workflow.id}`}
          >
            <Button size={'icon'} className="cursor-pointer!">
              <Play />
            </Button>
          </Link>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant={'ghost'}
                size={'icon-sm'}
                className="cursor-pointer!"
              >
                <EllipsisVertical className="size-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={() => handleFavorite(workflow)}>
                <Pin fill={workflow.isPinned ? 'currentColor' : 'none'} />{' '}
                {workflow.isPinned ? t('ui.text.unpin') : t('ui.text.pin')}
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link
                  to={`/organizations/${workflow.organizationId}/edit/${workflow.id}`}
                >
                  <Pencil /> {t('ui.text.edit')}
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleDuplicate(workflow)}>
                <Copy /> {t('ui.text.duplicate')}
              </DropdownMenuItem>
              <CopyIdMenuItem id={workflow.id} />
              <DropdownMenuItem
                onClick={() => handleDelete(workflow)}
                variant="destructive"
              >
                <Trash2 /> {t('ui.text.delete')}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </div>
  )
}
PinnedWorkflowCard.displayName = 'PinnedWorkflowCard'

export { PinnedWorkflowCard }
