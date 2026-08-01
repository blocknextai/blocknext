import { memo } from 'react'
import { Separator } from '@/components/ui/separator'
import { PinnedWorkflowCard } from '@/features/organizations/components/pinned-workflow-card'
import { Pin } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { EmptyContent } from './empty-content'

const PinnedWorkflowsContent = memo(
  ({ workflows, handleFavorite, handleDuplicate, handleDelete }) => {
    const { t } = useTranslation()
    const hasItems = workflows.length > 0

    return (
      <>
        <h2 className="text-xl font-semibold flex items-center gap-2">
          <Pin color="#FB6D6C" fill="#FB6D6C" />
          {t('ui.text.pinned')}
        </h2>

        <Separator className="my-4" />

        {!hasItems ? (
          <EmptyContent />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {workflows.map((workflow) => (
              <PinnedWorkflowCard
                key={workflow.id}
                workflow={workflow}
                handleFavorite={handleFavorite}
                handleDuplicate={handleDuplicate}
                handleDelete={handleDelete}
              />
            ))}
          </div>
        )}
      </>
    )
  },
)
PinnedWorkflowsContent.displayName = 'PinnedWorkflowsContent'

export { PinnedWorkflowsContent }
