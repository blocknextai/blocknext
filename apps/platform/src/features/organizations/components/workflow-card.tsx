import { memo } from 'react'
import { useTranslation } from 'react-i18next'
import { FlowIcon } from '@/components/shared/custom-icons'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
} from '@/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Button } from '@/components/ui/button'
import {
  Copy,
  EllipsisVertical,
  GitPullRequest,
  Pencil,
  Pin,
  Play,
  Trash2,
} from 'lucide-react'
import { Link } from 'react-router'
import { VerifiedBadge } from '@/components/shared/verified-tooltip'
import { CopyIdMenuItem } from '@/components/shared/copy-id-menu-item'

const WorkflowCard = memo(
  ({
    workflow,
    currentUserId,
    handleFavorite,
    handleEdit,
    handleDuplicate,
    handleDelete,
  }) => {
    const { t } = useTranslation()

    return (
      <Card className="group dark:border-transparent flex flex-col justify-between w-full h-full p-0! overflow-hidden relative transition-all duration-300 ease-in-out hover:scale-[1.03] hover:shadow-lg hover:z-10">
        <CardHeader className="p-0">
          <div className="flex flex-col relative z-2 overflow-hidden">
            <div
              className="w-full absolute -z-1 top-0 left-0"
              style={{
                backgroundImage: `${FlowIcon(workflow.id)}`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
                height: '100%',
              }}
            >
              <div className="absolute inset-0 -bottom-1 group-hover:-bottom-2 transition-all -left-2 w-200 bg-linear-to-b from-transparent to-card to-95% will-change-transform transform-gpu"></div>
            </div>

            <div className="flex justify-end items-center w-full pt-2 pr-2">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant={'ghost'}
                    size={'icon-sm'}
                    className="bg-card/80 text-foreground backdrop-blur-sm hover:bg-card"
                  >
                    <EllipsisVertical className="size-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuItem onClick={() => handleFavorite(workflow)}>
                    <Pin fill={workflow.isPinned ? 'currentColor' : 'none'} />{' '}
                    {workflow.isPinned ? t('ui.text.unpin') : t('ui.text.pin')}
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => handleEdit(workflow)}>
                    <Pencil /> {t('ui.text.edit')}
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
            <div className="flex items-center px-6 gap-2">
              <div className="heading-md whitespace-normal break-all line-clamp-2">
                {workflow.title}
              </div>
            </div>
            {workflow.owner?.id !== currentUserId && (
              <div className="flex items-center gap-1.5 px-6 py-1 text-sm text-muted-foreground">
                <span className="truncate">
                  {t('ui.text.createdBy')}{' '}
                  <span className="text-primary underline">
                    {workflow.owner?.alias}
                  </span>
                </span>
                {workflow.owner?.isVerified && <VerifiedBadge type="user" />}
              </div>
            )}
            <CardDescription className="px-6 pt-2 text-muted-foreground/80">
              {workflow.description.length > 80
                ? `${workflow.description.slice(0, 80)}...`
                : workflow.description}
            </CardDescription>
          </div>
        </CardHeader>

        <CardContent className="flex justify-start w-full p-6 gap-2">
          <Link
            className="flex-1"
            to={`/organizations/${workflow.organizationId}/edit/${workflow.id}`}
          >
            <Button className="w-full" variant={'secondary'}>
              <GitPullRequest /> {t('ui.text.edit')}
            </Button>
          </Link>
          <Link
            className="flex-1"
            to={`/organizations/${workflow.organizationId}/run/${workflow.id}`}
          >
            <Button className="w-full">
              <Play /> {t('ui.text.run')}
            </Button>
          </Link>
        </CardContent>
      </Card>
    )
  },
)
WorkflowCard.displayName = 'WorkflowCard'

export { WorkflowCard }
