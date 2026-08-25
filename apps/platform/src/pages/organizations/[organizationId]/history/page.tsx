import { Link, useParams, useSearchParams } from 'react-router'
import { Button } from '@/components/ui/button'
import { Loading } from '@/components/shared/loading'
import { useCallback, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Plus, History, PlayCircle, Wrench } from 'lucide-react'
import {
  useExecutions,
  useExecutionActions,
  useToolInvocations,
} from '@/features/workflows'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ToolInvocationsTable } from '@/features/organizations/components/history/tool-invocations-table'
import { ToolInvocationDetailDialog } from '@/features/organizations/components/history/tool-invocation-detail-dialog'

import ConfirmationDialog from '@/components/shared/confirmation-dialog'
import { AppPagination } from '@/components/shared/app-pagination'
import { SearchInput } from '@/components/shared/search-input'
import { SearchEmptyState } from '@/components/shared/search-empty-state'
import { HistoryTable } from '@/features/organizations/components/history/history-table'
import { HistoryFilters } from '@/features/organizations/components/history/history-filters'
import { useOrganizationEvents } from '@/hooks/use-organization-events'

function OrganizationHistoryPage() {
  const { t } = useTranslation()
  const { organizationId } = useParams()
  const [searchParams, setSearchParams] = useSearchParams()
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState('')
  const [selectedToolInvocationId, setSelectedToolInvocationId] = useState<
    string | null
  >(null)
  const limit = 10
  const offset = (page - 1) * limit

  const activeTab = searchParams.get('tab') === 'tools' ? 'tools' : 'runs'

  const setActiveTab = (value: string) => {
    setPage(1)
    setQuery('')
    setSearchParams((prev) => {
      prev.set('tab', value)
      return prev
    })
  }

  const handleSearch = useCallback((value: string) => {
    setQuery(value)
    setPage(1)
  }, [])

  const {
    executions: history,
    pagination,
    isLoading,
    mutate,
  } = useExecutions(organizationId, { offset, limit, query })
  const { remove, cancel, rerunAll, rerunFailed } =
    useExecutionActions(organizationId)
  const {
    toolInvocations,
    pagination: toolPagination,
    isLoading: toolsLoading,
    mutate: mutateToolInvocations,
  } = useToolInvocations(organizationId, { offset, limit, query })

  const [open, setOpen] = useState(false)
  const [confirmData, setConfirmData] = useState<
    | {
        description: string
        action: () => Promise<void> | void
        label: string
      }
    | undefined
  >()

  const EmptyState = () => (
    <div className="w-full h-full flex flex-col items-center justify-center text-center">
      <div className="p-4 rounded-full w-16 h-16 mx-auto flex items-center justify-center">
        <History className="size-10 text-muted-foreground/80" />
      </div>
      <h3 className="text-xl font-semibold mb-2">
        {t('ui.text.noRunHistoryYet')}
      </h3>
      <p className="text-muted-foreground mb-6 max-w-md mx-auto">
        {t('ui.text.noRunHistoryDescription')}
      </p>
      <Button className="gap-2 pl-2!" variant={'outline'} asChild>
        <Link to={`/organizations/${organizationId}/create`}>
          <Plus className="h-4 w-4" />
          {t('ui.text.runFirstFlow')}
        </Link>
      </Button>
    </div>
  )

  const ToolCallsEmptyState = () => (
    <div className="w-full h-full flex flex-col items-center justify-center text-center">
      <div className="p-4 rounded-full w-16 h-16 mx-auto flex items-center justify-center">
        <Wrench className="size-10 text-muted-foreground/80" />
      </div>
      <h3 className="text-xl font-semibold mb-2">
        {t('ui.text.noToolCallsYet')}
      </h3>
      <p className="text-muted-foreground mb-6 max-w-md mx-auto">
        {t('ui.text.noToolCallsDescription')}
      </p>
    </div>
  )

  const confirmDelete = (id: string) => {
    const h = history.find((h) => h.id === id)
    if (!h) {
      return
    }
    setConfirmData({
      description: t('ui.text.deleteHistoryConfirmation'),
      action: () => remove(h.id),
      label: t('ui.text.delete'),
    })
    setTimeout(() => setOpen(true), 150)
  }

  const confirmCancel = (id: string) => {
    const h = history.find((h) => h.id === id)
    if (!h) {
      return
    }
    setConfirmData({
      description: t('ui.text.cancelExecutionConfirmation'),
      action: () => cancel({ organizationId, id: h.id }).then(() => mutate()),
      label: t('ui.text.terminate'),
    })
    setTimeout(() => setOpen(true), 150)
  }

  const rerunAllExecution = useCallback(
    async (h: { id: string }) => {
      await rerunAll({ organizationId, id: h.id })
      await mutate()
    },
    [rerunAll, organizationId, mutate],
  )

  const rerunFailedExecution = useCallback(
    async (h: { id: string }) => {
      await rerunFailed({ organizationId, id: h.id })
      await mutate()
    },
    [rerunFailed, organizationId, mutate],
  )

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const handleEvent = useCallback(() => {
    if (debounceRef.current) {
      clearTimeout(debounceRef.current)
    }
    debounceRef.current = setTimeout(() => {
      mutate()
    }, 300)
  }, [mutate])

  useOrganizationEvents(handleEvent, { type: 'task' })

  const toolDebounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const handleToolEvent = useCallback(() => {
    if (toolDebounceRef.current) {
      clearTimeout(toolDebounceRef.current)
    }
    toolDebounceRef.current = setTimeout(() => {
      mutateToolInvocations()
    }, 300)
  }, [mutateToolInvocations])

  useOrganizationEvents(handleToolEvent, { type: 'tool_invocation' })

  const isInitialLoading = isLoading && pagination === undefined
  const isEmptyNoQuery = !isInitialLoading && history.length === 0 && !query
  const isEmptySearch = !isInitialLoading && history.length === 0 && !!query

  const activePagination = activeTab === 'tools' ? toolPagination : pagination

  return (
    <div className="w-full h-full flex px-6 pt-3 flex-col">
      <HistoryFilters />

      <Tabs value={activeTab} onValueChange={setActiveTab} className="mt-3">
        <TabsList>
          <TabsTrigger value="runs">
            <PlayCircle className="size-4" />
            {t('ui.text.runs')}
          </TabsTrigger>
          <TabsTrigger value="tools">
            <Wrench className="size-4" />
            {t('ui.text.toolCalls')}
          </TabsTrigger>
        </TabsList>

        <div className="mt-3 max-w-sm">
          <SearchInput
            placeholder={t('ui.text.searchPlaceholder')}
            onSearch={handleSearch}
          />
        </div>

        <TabsContent value="runs">
          {isInitialLoading ? (
            <Loading />
          ) : isEmptyNoQuery ? (
            <EmptyState />
          ) : isEmptySearch ? (
            <SearchEmptyState />
          ) : (
            <div className="mt-6">
              <HistoryTable
                organizationId={organizationId}
                filteredHistory={history}
                onRerunAll={rerunAllExecution}
                onRerunFailed={rerunFailedExecution}
                onConfirmCancel={confirmCancel}
                onConfirmDelete={confirmDelete}
              />
            </div>
          )}
        </TabsContent>

        <TabsContent value="tools">
          {toolsLoading && toolPagination === undefined ? (
            <Loading />
          ) : toolInvocations.length === 0 ? (
            query ? (
              <SearchEmptyState />
            ) : (
              <ToolCallsEmptyState />
            )
          ) : (
            <div className="mt-6">
              <ToolInvocationsTable
                toolInvocations={toolInvocations}
                onSelect={setSelectedToolInvocationId}
              />
            </div>
          )}
        </TabsContent>
      </Tabs>

      {activePagination && (
        <AppPagination
          pagination={{
            total: activePagination.total,
            page,
            limit,
            offset,
            hasPrev: activePagination.hasPrev,
            hasNext: activePagination.hasNext,
          }}
          onPageChange={setPage}
        />
      )}

      <ToolInvocationDetailDialog
        organizationId={organizationId}
        toolInvocationId={selectedToolInvocationId}
        onOpenChange={(open) => {
          if (!open) {
            setSelectedToolInvocationId(null)
          }
        }}
      />

      <ConfirmationDialog
        open={open}
        onOpenChange={setOpen}
        title={t('ui.text.areYouSure')}
        description={confirmData?.description}
        confirmText={confirmData?.label}
        onConfirm={confirmData?.action}
        variant="destructive"
      />
    </div>
  )
}

export default OrganizationHistoryPage
