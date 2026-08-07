import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Loading } from '@/components/shared/loading'
import { Separator } from '@/components/ui/separator'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import ConfirmationDialog from '@/components/shared/confirmation-dialog'
import { ApiKeyTable } from '@/features/api-keys/components/api-key-table'
import { ApiKeyDetailSheet } from '@/features/api-keys/components/api-key-detail-sheet'
import { AppPagination } from '@/components/shared/app-pagination'
import { SearchInput } from '@/components/shared/search-input'
import { SearchEmptyState } from '@/components/shared/search-empty-state'
import { CopyField } from '@/components/shared/copy-field'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Plus, Key } from 'lucide-react'

const APIKeyTab = ({
  loading,
  apiKeys,
  pagination,
  scopes = [],
  loadingScopes,
  onFetch,
  onCreate,
  onRegenerate,
  onDelete,
}) => {
  const { t } = useTranslation()
  const [createOpen, setCreateOpen] = useState(false)
  const [newKeyName, setNewKeyName] = useState('')
  const [selectedScopes, setSelectedScopes] = useState<string[]>([])
  const [createdKey, setCreatedKey] = useState<any>(null)
  const [detailsKey, setDetailsKey] = useState<any>(null)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [confirmData, setConfirmData] = useState<any>(null)
  const [query, setQuery] = useState('')

  const resetCreateForm = useCallback(() => {
    setNewKeyName('')
    setSelectedScopes([])
  }, [])

  const toggleScope = useCallback((scope: string, checked: boolean) => {
    setSelectedScopes((prev) =>
      checked ? [...prev, scope] : prev.filter((s) => s !== scope),
    )
  }, [])

  const scopeDescriptions: Record<string, string> = {
    'workflows:trigger': t(
      'ui.text.scopeWorkflowsTrigger',
      'Trigger workflows programmatically',
    ),
    'mcp:invoke': t('ui.text.scopeMcpInvoke', 'Invoke MCP tools'),
  }

  const handleSearch = useCallback(
    (value: string) => {
      setQuery(value)
      onFetch(0, pagination?.limit ?? 10, value)
    },
    [onFetch, pagination?.limit],
  )

  const handleCreate = async () => {
    if (!selectedScopes.length) {
      return
    }
    const created = await onCreate({
      name: newKeyName || t('ui.text.apiKey', 'API Key'),
      scopes: selectedScopes,
    })
    if (created) {
      setCreatedKey(created)
      resetCreateForm()
      setCreateOpen(false)
    }
  }

  const confirmDelete = (apiKey: any) => {
    setConfirmData({
      description: t('ui.text.deleteApiKeyConfirmation'),
      action: async () => {
        await onDelete(apiKey.id)
        setDetailsKey(null)
      },
      label: t('ui.text.delete'),
    })
    setConfirmOpen(true)
  }

  const confirmRegenerate = (apiKey: any) => {
    setConfirmData({
      description: t(
        'ui.text.revokeAndRegenerateConfirmation',
        'This will invalidate the current key. Any integrations using it will stop working until updated with the new key.',
      ),
      action: async () => {
        const regenerated = await onRegenerate(apiKey.id, {
          name: apiKey.name,
        })
        if (regenerated) {
          setDetailsKey(null)
          setCreatedKey(regenerated)
        }
      },
      label: t('ui.text.revokeAndRegenerate', 'Revoke & Regenerate'),
    })
    setConfirmOpen(true)
  }

  const handlePageChange = (newPage: number) => {
    const newOffset = (newPage - 1) * pagination.limit
    onFetch(newOffset, pagination.limit, query)
  }

  useEffect(() => {
    onFetch()
  }, [])

  return (
    <div className="w-full h-full flex flex-col">
      <div className="flex flex-col gap-2 mb-4 p-6 shrink-0">
        <div className="flex justify-between items-center">
          <h2 className="text-2xl font-bold flex items-center gap-2">
            {t('ui.text.apiKeys')}
          </h2>
          <Button size="sm" onClick={() => setCreateOpen(true)}>
            <Plus />
            <span>{t('ui.text.new')}</span>
          </Button>
        </div>
        <Separator className="my-2 p-0" decorative={true} />
        <div className="max-w-sm">
          <SearchInput
            placeholder={t('ui.text.searchApiKeys')}
            onSearch={handleSearch}
          />
        </div>
      </div>

      <div className="flex-1 min-h-0">
        {loading && !pagination?.total ? (
          <Loading />
        ) : !apiKeys?.length && query ? (
          <SearchEmptyState />
        ) : !apiKeys?.length ? (
          <div className="flex w-full h-full flex-col items-center justify-center text-center">
            <div className="flex flex-col items-center gap-3">
              <div className="p-3 rounded-full w-12 h-12 flex items-center justify-center">
                <Key className="size-6 text-muted-foreground/80" />
              </div>
              <div className="flex flex-col gap-2">
                <span className="text-md font-semibold">
                  {t('ui.text.noApiKeys')}
                </span>
                <span className="text-muted-foreground text-xs max-w-xs">
                  {t('ui.text.noApiKeysDescription')}
                </span>
              </div>
              <Button
                variant="outline"
                className="size-10"
                onClick={() => setCreateOpen(true)}
              >
                <Plus />
              </Button>
            </div>
          </div>
        ) : (
          <div className="pb-6">
            <ApiKeyTable
              apiKeys={apiKeys}
              onDetails={setDetailsKey}
              onConfirmRegenerate={confirmRegenerate}
              onConfirmDelete={confirmDelete}
            />
            <AppPagination
              pagination={pagination}
              onPageChange={handlePageChange}
            />
          </div>
        )}
      </div>

      <Dialog
        open={createOpen}
        onOpenChange={(open) => {
          setCreateOpen(open)
          if (!open) {
            resetCreateForm()
          }
        }}
      >
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{t('ui.text.createApiKey')}</DialogTitle>
            <DialogDescription>
              {t('ui.text.createApiKeyDescription')}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-2">
            <Label htmlFor="keyName">{t('ui.text.name')}</Label>
            <Input
              id="keyName"
              value={newKeyName}
              onChange={(e) => setNewKeyName(e.target.value)}
              placeholder={t('ui.text.myApiKey')}
            />
          </div>
          <div className="grid gap-2">
            <Label>{t('ui.text.scopes', 'Scopes')}</Label>
            <p className="text-xs text-muted-foreground">
              {t(
                'ui.text.scopesDescription',
                'Select the permissions this API key will have.',
              )}
            </p>
            {loadingScopes ? (
              <Loading />
            ) : (
              <div className="flex flex-col gap-1.5 rounded-md border p-1">
                {scopes.map((scope: { key: string }) => {
                  const checked = selectedScopes.includes(scope.key)
                  return (
                    <Label
                      key={scope.key}
                      className="flex items-start gap-3 rounded-sm p-2 font-normal hover:bg-muted/50"
                    >
                      <Checkbox
                        className="mt-0.5"
                        checked={checked}
                        onCheckedChange={(value) =>
                          toggleScope(scope.key, value === true)
                        }
                      />
                      <span className="flex flex-col gap-0.5">
                        <span className="font-mono text-sm">{scope.key}</span>
                        {scopeDescriptions[scope.key] && (
                          <span className="text-xs text-muted-foreground">
                            {scopeDescriptions[scope.key]}
                          </span>
                        )}
                      </span>
                    </Label>
                  )
                })}
              </div>
            )}
          </div>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline" size="sm">
                {t('ui.text.cancel')}
              </Button>
            </DialogClose>
            <Button
              size="sm"
              onClick={handleCreate}
              disabled={!selectedScopes.length}
            >
              {t('ui.text.create')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ApiKeyDetailSheet
        apiKey={detailsKey}
        onOpenChange={(open) => !open && setDetailsKey(null)}
        onRegenerate={confirmRegenerate}
        onDelete={confirmDelete}
      />

      <Dialog open={!!createdKey} onOpenChange={() => setCreatedKey(null)}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>{t('ui.text.apiKeyCreated')}</DialogTitle>
            <DialogDescription>
              {t('ui.text.apiKeyCreatedDescription')}
            </DialogDescription>
          </DialogHeader>
          {createdKey && <CopyField value={createdKey.key} />}
          <DialogFooter>
            <Button size="sm" onClick={() => setCreatedKey(null)}>
              {t('ui.text.done')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ConfirmationDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title={t('ui.text.areYouSure')}
        description={confirmData?.description}
        confirmText={confirmData?.label}
        onConfirm={confirmData?.action}
        variant="destructive"
      />
    </div>
  )
}

APIKeyTab.displayName = 'APIKeyTab'

export default APIKeyTab
