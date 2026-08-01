import { memo } from 'react'
import { Separator } from '@/components/ui/separator'
import { WorkflowCard } from '@/features/organizations/components/workflow-card'
import { useTranslation } from 'react-i18next'
import { EmptyContent } from './empty-content'

const ByYouWorkflowsContent = memo(
  ({
    workflows,
    currentUserId,
    handleFavorite,
    handleEdit,
    handleDuplicate,
    handleDelete,
  }) => {
    const { t } = useTranslation()

    return (
      <>
        <h2 className="text-xl font-semibold flex items-center gap-2">
          {t('ui.text.createdByYou')}
        </h2>

        <Separator className="my-4" />

        {workflows.length === 0 ? (
          <EmptyContent />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {workflows.map((workflow) => (
              <WorkflowCard
                key={workflow.id}
                workflow={workflow}
                currentUserId={currentUserId}
                handleFavorite={handleFavorite}
                handleEdit={handleEdit}
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
ByYouWorkflowsContent.displayName = 'ByYouWorkflowsContent'

export { ByYouWorkflowsContent }
