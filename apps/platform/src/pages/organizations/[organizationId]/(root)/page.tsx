import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useCallback, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate, useParams, useSearchParams } from 'react-router'
import { useWorkflows, useWorkflowActions } from '@/features/workflows'
import { useOrganizationMe } from '@/features/organizations'
import { Loading } from '@/components/shared/loading'
import { PinnedWorkflowsContent } from '@/features/organizations/components/workspace/pinned-workflows-content'
import { ByYouWorkflowsContent } from '@/features/organizations/components/workspace/by-you-workflows-content'
import { ByOthersWorkflowsContent } from '@/features/organizations/components/workspace/by-others-workflows-content'
import { EditDialog } from '@/features/organizations/components/workspace/edit-dialog'
import { DuplicateDialog } from '@/features/organizations/components/workspace/duplicate-dialog'
import { DeleteDialog } from '@/features/organizations/components/workspace/delete-dialog'
import { EmptyWorkspace } from '@/features/organizations/components/workspace/empty-workspace'

function OrganizationDetailPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const { organizationId } = useParams()

  const { me, isLoading: meLoading } = useOrganizationMe(organizationId)
  const { workflows, isLoading: wfLoading } = useWorkflows(organizationId, {
    limit: 100,
  })

  const {
    update: updateWorkflow,
    remove: removeWorkflow,
    duplicate: duplicateWorkflow,
  } = useWorkflowActions(organizationId)

  const [selectedWorkflow, setSelectedWorkflow] = useState<any>(null)
  const [isEditedDialogOpen, setIsEditedDialogOpen] = useState(false)
  const [isDuplicateDialogOpen, setIsDuplicateDialogOpen] = useState(false)
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)

  const activeTab = searchParams.get('tab') || 'all'
  const currentUserId = me?.id ?? null

  const { pinnedWorkflows, byYouWorkflows, byOthersWorkflows } = useMemo(() => {
    return {
      pinnedWorkflows: workflows.filter((w) => w.isPinned),
      byYouWorkflows: workflows.filter((w) => w.owner?.id === currentUserId),
      byOthersWorkflows: workflows.filter((w) => w.owner?.id !== currentUserId),
    }
  }, [workflows, currentUserId])

  const isLoading = meLoading || wfLoading

  const isEmptyWorkspace =
    !isLoading &&
    pinnedWorkflows.length === 0 &&
    byYouWorkflows.length === 0 &&
    byOthersWorkflows.length === 0

  const handleTabChange = useCallback(
    (value: string) => {
      navigate(`/organizations/${organizationId}?tab=${value}`)
    },
    [navigate, organizationId],
  )

  const updateFavorite = useCallback(
    async (workflow) => {
      if (!workflow) {
        return
      }
      await updateWorkflow(workflow.id, { isPinned: !workflow.isPinned })
    },
    [updateWorkflow],
  )

  const handleEdit = useCallback((workflow) => {
    setSelectedWorkflow(workflow)
    setIsEditedDialogOpen(true)
  }, [])

  const confirmEdit = useCallback(async () => {
    if (!selectedWorkflow) {
      return
    }
    await updateWorkflow(selectedWorkflow.id, {
      title: selectedWorkflow.title || '',
      description: selectedWorkflow.description || '',
    })
    setIsEditedDialogOpen(false)
    setSelectedWorkflow(null)
  }, [selectedWorkflow, updateWorkflow])

  const handleDuplicate = useCallback((workflow) => {
    setSelectedWorkflow(workflow)
    setIsDuplicateDialogOpen(true)
  }, [])

  const confirmDuplicate = useCallback(async () => {
    if (!selectedWorkflow) {
      return
    }
    await duplicateWorkflow(selectedWorkflow.id, {
      title: selectedWorkflow.title || '',
      description: selectedWorkflow.description || '',
    })
    setIsDuplicateDialogOpen(false)
    setSelectedWorkflow(null)
  }, [selectedWorkflow, duplicateWorkflow])

  const handleDelete = useCallback((workflow) => {
    setSelectedWorkflow(workflow)
    setIsDeleteDialogOpen(true)
  }, [])

  const confirmDelete = useCallback(async () => {
    if (!selectedWorkflow) {
      return
    }
    await removeWorkflow(selectedWorkflow.id)
    setIsDeleteDialogOpen(false)
    setSelectedWorkflow(null)
  }, [selectedWorkflow, removeWorkflow])

  const closeEditDialog = useCallback(() => {
    setIsEditedDialogOpen(false)
    setSelectedWorkflow(null)
  }, [])

  const closeDuplicateDialog = useCallback(() => {
    setIsDuplicateDialogOpen(false)
    setSelectedWorkflow(null)
  }, [])

  const closeDeleteDialog = useCallback(() => {
    setIsDeleteDialogOpen(false)
    setSelectedWorkflow(null)
  }, [])

  if (isLoading) {
    return <Loading />
  }

  if (isEmptyWorkspace) {
    return <EmptyWorkspace organizationId={organizationId} />
  }

  return (
    <div className="p-6 pt-3">
      <Tabs defaultValue={activeTab} onValueChange={handleTabChange}>
        <TabsList className="mb-4 w-full justify-start flex-wrap flex-row">
          <TabsTrigger value="all">{t('ui.text.all')}</TabsTrigger>
          <TabsTrigger value="byYou">{t('ui.text.byYou')}</TabsTrigger>
          <TabsTrigger value="byOthers">{t('ui.text.byOthers')}</TabsTrigger>
        </TabsList>
        <TabsContent value="all" className="space-y-4">
          <PinnedWorkflowsContent
            workflows={pinnedWorkflows}
            organizationId={organizationId}
            handleFavorite={updateFavorite}
            handleDuplicate={handleDuplicate}
            handleDelete={handleDelete}
          />
          <ByYouWorkflowsContent
            workflows={byYouWorkflows}
            currentUserId={currentUserId}
            handleFavorite={updateFavorite}
            handleEdit={handleEdit}
            handleDuplicate={handleDuplicate}
            handleDelete={handleDelete}
          />
          <ByOthersWorkflowsContent
            workflows={byOthersWorkflows}
            currentUserId={currentUserId}
            handleFavorite={updateFavorite}
            handleEdit={handleEdit}
            handleDuplicate={handleDuplicate}
            handleDelete={handleDelete}
          />
        </TabsContent>
        <TabsContent value="byYou">
          <ByYouWorkflowsContent
            workflows={byYouWorkflows}
            currentUserId={currentUserId}
            handleFavorite={updateFavorite}
            handleEdit={handleEdit}
            handleDuplicate={handleDuplicate}
            handleDelete={handleDelete}
          />
        </TabsContent>
        <TabsContent value="byOthers">
          <ByOthersWorkflowsContent
            workflows={byOthersWorkflows}
            currentUserId={currentUserId}
            handleFavorite={updateFavorite}
            handleEdit={handleEdit}
            handleDuplicate={handleDuplicate}
            handleDelete={handleDelete}
          />
        </TabsContent>
      </Tabs>

      <EditDialog
        open={isEditedDialogOpen}
        onOpenChange={closeEditDialog}
        selectedWorkflow={selectedWorkflow}
        setSelectedWorkflow={setSelectedWorkflow}
        onConfirmEdit={confirmEdit}
      />

      <DuplicateDialog
        open={isDuplicateDialogOpen}
        onOpenChange={closeDuplicateDialog}
        selectedWorkflow={selectedWorkflow}
        setSelectedWorkflow={setSelectedWorkflow}
        onConfirmDuplicate={confirmDuplicate}
      />

      <DeleteDialog
        open={isDeleteDialogOpen}
        onOpenChange={closeDeleteDialog}
        onConfirmDelete={confirmDelete}
      />
    </div>
  )
}

export default OrganizationDetailPage
