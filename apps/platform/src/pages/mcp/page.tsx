import { useCallback, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Loading } from '@/components/shared/loading'
import { SearchInput } from '@/components/shared/search-input'
import { AppPagination } from '@/components/shared/app-pagination'
import { useMcpServers } from '@/features/mcp'
import { McpServerGrid } from '@/features/mcp/components/mcp-server-grid'
import { McpServerDetailDialog } from '@/features/mcp/components/mcp-server-detail-dialog'
import type { McpServer } from '@/features/mcp/services/mcp'

const matches = (value: string | undefined, q: string) =>
  Boolean(value && value.toLowerCase().includes(q))

const filterServers = (servers: McpServer[], query: string): McpServer[] => {
  const q = query.trim().toLowerCase()
  if (!q) return servers
  return servers.filter(
    (server) =>
      matches(server.name, q) ||
      matches(server.description, q) ||
      matches(server.id, q),
  )
}

const LIMIT = 10

const McpPage = () => {
  const { t } = useTranslation()
  const [query, setQuery] = useState('')
  const [page, setPage] = useState(1)
  const [selectedServer, setSelectedServer] = useState<McpServer | null>(null)
  const [dialogOpen, setDialogOpen] = useState(false)

  const { servers, isLoading } = useMcpServers()

  const handleSearch = useCallback((value: string) => {
    setQuery(value)
    setPage(1)
  }, [])

  const handleCardClick = useCallback((server: McpServer) => {
    setSelectedServer(server)
    setDialogOpen(true)
  }, [])

  const filtered = useMemo(
    () => filterServers(servers, query),
    [servers, query],
  )

  const offset = (page - 1) * LIMIT
  const paginated = useMemo(
    () => filtered.slice(offset, offset + LIMIT),
    [filtered, offset],
  )

  const isInitialLoading = isLoading && servers.length === 0

  return (
    <div className="w-full h-full p-6 space-y-6">
      <div className="flex items-center justify-between gap-4 flex-wrap">
        <h1 className="text-3xl font-bold">{t('ui.text.mcpServers')}</h1>
        <div className="w-full sm:w-80">
          <SearchInput
            placeholder={t('ui.text.mcpSearchPlaceholder')}
            onSearch={handleSearch}
          />
        </div>
      </div>

      {isInitialLoading ? (
        <Loading />
      ) : servers.length === 0 ? (
        <div className="flex justify-center items-center w-full text-muted-foreground gap-2 p-10">
          {t('ui.text.mcpNoServersAvailable')}
        </div>
      ) : filtered.length === 0 ? (
        <div className="flex justify-center items-center w-full text-muted-foreground gap-2 p-10">
          {t('ui.text.mcpNoSearchResults')}
        </div>
      ) : (
        <>
          <McpServerGrid servers={paginated} onCardClick={handleCardClick} />
          <AppPagination
            pagination={{
              total: filtered.length,
              page,
              limit: LIMIT,
              offset,
              hasPrev: page > 1,
              hasNext: offset + LIMIT < filtered.length,
            }}
            onPageChange={setPage}
          />
        </>
      )}

      <McpServerDetailDialog
        server={selectedServer}
        open={dialogOpen}
        onOpenChange={setDialogOpen}
      />
    </div>
  )
}

export default McpPage
