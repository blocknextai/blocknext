import { useCallback, useEffect, useState } from 'react'
import { Loading } from '@/components/shared/loading'
import { Separator } from '@/components/ui/separator'
import { Button } from '@/components/ui/button'
import { useTranslation } from 'react-i18next'
import ConfirmationDialog from '@/components/shared/confirmation-dialog'
import { AppPagination } from '@/components/shared/app-pagination'
import { SearchInput } from '@/components/shared/search-input'
import { SearchEmptyState } from '@/components/shared/search-empty-state'
import { bc } from '@/lib/broadcast'
import { Plus, KeyRound } from 'lucide-react'
import { CredentialTable } from '@/features/credentials/components/credential-table'
import { CredentialFormDialog } from '@/features/credentials/components/credential-form-dialog'
import { CredentialPopover } from '@/features/credentials/components/credential-popover'

const CredentialTab = ({
  isUserMode = false,
  loading,
  secrets,
  credentialsArray,
  pagination,
  hasAccess,
  onFetch,
  getSecrets,
  getSecretDetails,
  onSave,
  onDelete,
  onAuthorizeOAuth2,
  onReauthorizeOAuth2,
}) => {
  const { t } = useTranslation()

  const [selectedSecret, setSelectedSecret] = useState<any>()
  const [open, setOpen] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [secretDetailLoading, setSecretDetailLoading] = useState(false)
  const [usePlatformCredential, setUsePlatformCredential] = useState(false)
  const [query, setQuery] = useState('')

  const handleSearch = useCallback(
    (value: string) => {
      setQuery(value)
      onFetch(0, pagination?.limit ?? 10, value)
    },
    [onFetch, pagination?.limit],
  )

  const handlePageChange = (newPage: number) => {
    const newOffset = (newPage - 1) * pagination.limit
    onFetch(newOffset, pagination.limit, query)
  }

  useEffect(() => {
    onFetch()
  }, [])

  const updateSelectedSecret = (
    s: { key: string; value: string },
    o = false,
  ) => {
    const nsa = Object.assign({}, selectedSecret)
    if (o) {
      if (!nsa.data) {
        nsa.data = {}
      }
      nsa.data[s.key] = s.value
    } else {
      nsa[s.key] = s.value
    }
    setSelectedSecret(nsa)
  }

  const selectSecret = async (s: any, w?: any) => {
    setUsePlatformCredential(false)

    if (s.uiKey && w) {
      const foundSecret = secrets.find((secret: any) => secret.id === s.key)
      const initial = JSON.parse(JSON.stringify(foundSecret))
      initial.ref = s.id
      initial.value = s.uiKey
      initial.label = s.title
      initial.isSupportPlatform =
        s.isSupportPlatform || foundSecret?.isSupportPlatform

      setSelectedSecret(initial)
      setSecretDetailLoading(true)
      w(true)

      const secretDetails = await getSecretDetails(s.id)
      setSecretDetailLoading(false)

      if (secretDetails) {
        setSelectedSecret((prev: any) => {
          if (!prev) return prev
          const next = { ...prev }
          next.secretDetails = secretDetails
          next.sourceType = secretDetails.sourceType
          if (secretDetails.data) next.data = { ...secretDetails.data }
          if (secretDetails.title) next.label = secretDetails.title
          return next
        })
        if (secretDetails.sourceType === 'platform') {
          setUsePlatformCredential(true)
        }
      }
      return
    }

    const ns = s
    ns.isNew = true
    ns.label = t(ns.name)
    const c = JSON.parse(JSON.stringify(ns))
    setSelectedSecret(c)
    if (w) w(true)
    else setOpen(true)
  }

  const saveCredentials = async () => {
    if (!hasAccess) {
      console.error(
        isUserMode ? 'User mode error' : 'Organization not available',
      )
      return
    }

    const ca = selectedSecret
    const payload: any = {
      title:
        (ca?.label ? t(ca.label) : '') ||
        (ca?.title ? t(ca.title) : '') ||
        (ca?.name ? t(ca.name) : '') ||
        '',
      key: ca.id,
      data: usePlatformCredential ? {} : ca.data,
      ref: ca.ref,
    }

    if (usePlatformCredential && ca.isNew) {
      payload.sourceType = 'platform'
    }

    try {
      await onSave(payload)
      getSecrets()
      setOpen(false)
      setUsePlatformCredential(false)
    } catch (error) {
      console.error('Error saving credentials:', error)
    }
  }

  const deleteSecret = async () => {
    if (!hasAccess || !selectedSecret?.ref) {
      console.error('Base path or secret ref not available')
      return
    }

    try {
      await onDelete(selectedSecret.ref)
      getSecrets()
      setOpenDelete(false)
    } catch (error) {
      console.error('Error deleting credential:', error)
    }
  }

  const authorizeOAuth2 = async () => {
    if (!hasAccess) {
      console.error(
        isUserMode ? 'User mode error' : 'Organization not available',
      )
      return
    }

    const ca = selectedSecret
    const payload: any = {
      title: ca.label || ca.title || t(ca.name),
      key: ca.id,
      data: usePlatformCredential ? {} : ca.data || {},
    }

    if (usePlatformCredential) {
      payload.sourceType = 'platform'
    }

    try {
      await onAuthorizeOAuth2(payload)
      setUsePlatformCredential(false)
    } catch (error) {
      console.error('Error during OAuth2 authorization:', error)
    }
  }

  const reauthorizeOAuth2 = async () => {
    if (!selectedSecret?.ref) {
      console.error('Credential id not available')
      return
    }

    try {
      await onReauthorizeOAuth2(selectedSecret.ref)
    } catch (error) {
      console.error('Error during OAuth2 re-authorization:', error)
    }
  }

  useEffect(() => {
    const handleOAuthMessage = (e: MessageEvent) => {
      if (e.data === 'hideAuth') {
        setOpen(false)
        getSecrets()
      }
    }
    bc.addEventListener('message', handleOAuthMessage)
    return () => {
      bc.removeEventListener('message', handleOAuthMessage)
    }
  }, [])

  return (
    <div className="w-full h-full flex flex-col">
      <div className="flex flex-col gap-2 mb-4 p-6 shrink-0">
        <div className="flex justify-between items-center">
          <h2 className="text-2xl font-bold flex items-center gap-2">
            {t('ui.text.credentials')}
          </h2>
          <CredentialPopover
            secrets={secrets}
            getSecrets={getSecrets}
            onSelect={selectSecret}
            isUserMode={isUserMode}
          >
            <Button size={'sm'}>
              <Plus />
              <span>{t('ui.text.new')}</span>
            </Button>
          </CredentialPopover>
        </div>
        <Separator className="my-2 p-0" decorative={true} />
        <div className="max-w-sm">
          <SearchInput
            placeholder={t('ui.text.searchCredentials')}
            onSearch={handleSearch}
          />
        </div>
      </div>

      <div className="flex-1 min-h-0">
        {loading && !pagination?.total ? (
          <Loading />
        ) : !credentialsArray?.length && query ? (
          <SearchEmptyState />
        ) : !credentialsArray?.length ? (
          <div className="flex w-full h-full flex-col items-center justify-center text-center">
            <div className="flex flex-col items-center gap-3">
              <div className="p-3 rounded-full w-12 h-12 flex items-center justify-center">
                <KeyRound className="size-6 text-muted-foreground/80" />
              </div>
              <div className="flex flex-col gap-2">
                <span className="text-md font-semibold">
                  {t('ui.text.noCredentialsYet')}
                </span>
                <span className="text-muted-foreground text-xs max-w-xs">
                  {t('ui.text.noCredentialsDescription')}
                </span>
              </div>
              <CredentialPopover
                secrets={secrets}
                getSecrets={getSecrets}
                onSelect={selectSecret}
                isUserMode={isUserMode}
              >
                <Button variant={'outline'} className="size-10">
                  <Plus />
                </Button>
              </CredentialPopover>
            </div>
          </div>
        ) : (
          <div className="pb-6">
            <CredentialTable
              credentialsArray={credentialsArray}
              onEdit={(c) => selectSecret(c, setOpen)}
              onDelete={(c) => selectSecret(c, setOpenDelete)}
            />
            <AppPagination
              pagination={pagination}
              onPageChange={handlePageChange}
            />
          </div>
        )}
      </div>

      <CredentialFormDialog
        open={open}
        onOpenChange={setOpen}
        selectedSecret={selectedSecret}
        secretDetailLoading={secretDetailLoading}
        usePlatformCredential={usePlatformCredential}
        setUsePlatformCredential={setUsePlatformCredential}
        updateSelectedSecret={updateSelectedSecret}
        saveCredentials={saveCredentials}
        authorizeOAuth2={authorizeOAuth2}
        reauthorizeOAuth2={reauthorizeOAuth2}
      />

      <ConfirmationDialog
        open={openDelete}
        onOpenChange={setOpenDelete}
        title={t('ui.text.areYouAbsolutelySure')}
        description={t('ui.text.deleteSecretDescription')}
        confirmText={t('ui.text.delete')}
        onConfirm={deleteSecret}
        variant="destructive"
      />
    </div>
  )
}

CredentialTab.displayName = 'CredentialTab'

export default CredentialTab
