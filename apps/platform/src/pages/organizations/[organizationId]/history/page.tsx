import { Link, useParams } from 'react-router'
import { Button } from '@/components/ui/button'
import { Loading } from '@/components/shared/loading'
import { useCallback, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Plus, History } from 'lucide-react'
import { useExecutions, useExecutionActions } from '@/features/workflows'

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
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState('')
  const limit = 10
  const offset = (page - 1) * limit

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

  const confirmDelete = (id: string) => {
    const h = history.find((h) => h.id === id)
    if (!h) return
    setConfirmData({
      description: t('ui.text.deleteHistoryConfirmation'),
      action: () => remove(h.id),
      label: t('ui.text.delete'),
    })
    setTimeout(() => setOpen(true), 150)
  }

  const confirmCancel = (id: string) => {
    const h = history.find((h) => h.id === id)
    if (!h) return
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
    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => {
      mutate()
    }, 300)
  }, [mutate])

  useOrganizationEvents(handleEvent, { type: 'task' })

  const isInitialLoading = isLoading && pagination === undefined
  const isEmptyNoQuery = !isInitialLoading && history.length === 0 && !query
  const isEmptySearch = !isInitialLoading && history.length === 0 && !!query

  return (
    <div className="w-full h-full flex px-6 pt-3 flex-col">
      <HistoryFilters />
      <div className="mt-3 max-w-sm">
        <SearchInput
          placeholder={t('ui.text.searchPlaceholder')}
          onSearch={handleSearch}
        />
      </div>
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

          {pagination && (
            <AppPagination
              pagination={{
                total: pagination.total,
                page,
                limit,
                offset,
                hasPrev: pagination.hasPrev,
                hasNext: pagination.hasNext,
              }}
              onPageChange={setPage}
            />
          )}
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
      )}
    </div>
  )
}

export default OrganizationHistoryPage
