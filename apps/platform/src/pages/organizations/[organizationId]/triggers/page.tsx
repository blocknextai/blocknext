import { Link, useParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Loading } from '@/components/shared/loading'
import { useCallback, useEffect, useState } from 'react'
import {
  Plus,
  Zap,
  Play,
  Pause,
  Trash2,
  RefreshCcw,
  ShieldCheck,
} from 'lucide-react'
import { useTriggers, useTriggerActions } from '@/features/workflows'
import { useNodeEngineWebhookSources } from '@/features/flow-editor'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from '@/components/ui/sheet'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import ConfirmationDialog from '@/components/shared/confirmation-dialog'
import { AppPagination } from '@/components/shared/app-pagination'
import { SearchInput } from '@/components/shared/search-input'
import { SearchEmptyState } from '@/components/shared/search-empty-state'
import { CopyField } from '@/components/shared/copy-field'
import { TriggerTable } from '@/features/organizations/components/triggers/trigger-table'
import { TriggerActions } from '@/features/organizations/components/triggers/trigger-actions'
import TimeAgoI18n from '@/components/shared/timeagoi18'

function OrganizationTriggersPage() {
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

  const { triggers, pagination, isLoading } = useTriggers(organizationId, {
    offset,
    limit,
    query,
  })
  const { update, remove, regenerateWebhookToken } =
    useTriggerActions(organizationId)

  const [open, setOpen] = useState(false)
  const [confirmData, setConfirmData] = useState<
    | {
        description: string
        action: () => Promise<void> | void
        label: string
      }
    | undefined
  >()
  const [detailsTrigger, setDetailsTrigger] = useState<any>(null)
  const [editCron, setEditCron] = useState('')
  const [editTimezone, setEditTimezone] = useState('')
  const [isSavingSchedule, setIsSavingSchedule] = useState(false)
  const [isRegenerating, setIsRegenerating] = useState(false)
  const [regeneratedToken, setRegeneratedToken] = useState<string | null>(null)
  const [editWebhookSecret, setEditWebhookSecret] = useState('')
  const [isSavingSecret, setIsSavingSecret] = useState(false)
  const showWebhookDetails = detailsTrigger?.type === 'webhook'
  const isSchedule = detailsTrigger?.type === 'schedule'
  const { sources: webhookSources } = useNodeEngineWebhookSources()

  useEffect(() => {
    if (!detailsTrigger) {
      return
    }
    setEditCron(detailsTrigger.cronPattern ?? '')
    setEditTimezone(detailsTrigger.timezone ?? '')
    setEditWebhookSecret('')
  }, [detailsTrigger?.id])

  const initialCron = detailsTrigger?.cronPattern ?? ''
  const initialTimezone = detailsTrigger?.timezone ?? ''
  const isScheduleDirty =
    isSchedule &&
    (editCron.trim() !== initialCron || editTimezone.trim() !== initialTimezone)
  const canSaveSchedule = isScheduleDirty && editCron.trim().length > 0

  const saveSchedule = async () => {
    if (!detailsTrigger || !canSaveSchedule) {
      return
    }
    setIsSavingSchedule(true)
    try {
      const payload: { cronPattern: string; timezone?: string } = {
        cronPattern: editCron.trim(),
      }
      const tz = editTimezone.trim()
      if (tz) {
        payload.timezone = tz
      }
      await update(detailsTrigger.id, payload)
      setDetailsTrigger(null)
    } finally {
      setIsSavingSchedule(false)
    }
  }

  const saveWebhookSecret = async () => {
    const secret = editWebhookSecret.trim()
    if (!detailsTrigger || !secret) {
      return
    }
    setIsSavingSecret(true)
    try {
      await update(detailsTrigger.id, { webhookSecret: secret })
      setDetailsTrigger(null)
    } finally {
      setIsSavingSecret(false)
    }
  }

  const confirmRemoveWebhookSecret = (triggerId: string) => {
    setConfirmData({
      description: t('ui.text.removeWebhookSecretConfirmation'),
      label: t('ui.text.remove'),
      action: async () => {
        await update(triggerId, { webhookSecret: '' })
        setDetailsTrigger(null)
      },
    })
    setOpen(true)
  }

  const lang = localStorage.getItem('i18nextLng') || 'en'

  const EmptyState = () => (
    <div className="w-full h-full flex flex-col items-center justify-center text-center">
      <div className="p-4 rounded-full w-16 h-16 mx-auto flex items-center justify-center">
        <Zap className="size-10 text-muted-foreground/80" />
      </div>
      <h3 className="text-xl font-semibold mb-2">
        {t('ui.text.noTriggersYet')}
      </h3>
      <p className="text-muted-foreground mb-6 max-w-md mx-auto">
        {t('ui.text.noTriggersDescription')}
      </p>
      <Button className="gap-2 pl-2!" variant={'outline'} asChild>
        <Link to={`/organizations/${organizationId}/create`}>
          <Plus className="h-4 w-4" />
          {t('ui.text.createFirstTrigger')}
        </Link>
      </Button>
    </div>
  )

  const confirmDelete = (id: string) => {
    const trigger = triggers.find((tr) => tr.id === id)
    if (!trigger) {
      return
    }
    setConfirmData({
      description: t('ui.text.deleteTriggerConfirmation'),
      action: () => remove(trigger.id),
      label: t('ui.text.delete'),
    })
    setOpen(true)
  }

  const toggleTrigger = async (trigger: any) => {
    await update(trigger.id, { isActive: !trigger.isActive })
  }

  const confirmRegenerateToken = (triggerId: string) => {
    setConfirmData({
      description: t('ui.text.regenerateTokenConfirmation'),
      label: t('ui.text.regenerateToken'),
      action: async () => {
        setIsRegenerating(true)
        try {
          const result: any = await regenerateWebhookToken(triggerId)
          const token = result?.webhookToken
          setDetailsTrigger(null)
          if (token) {
            setRegeneratedToken(token)
          }
        } finally {
          setIsRegenerating(false)
        }
      },
    })
    setOpen(true)
  }

  const isInitialLoading = isLoading && pagination === undefined
  const isEmptyNoQuery = !isInitialLoading && triggers.length === 0 && !query
  const isEmptySearch = !isInitialLoading && triggers.length === 0 && !!query

  return (
    <div className="w-full h-full flex px-6 pt-3 flex-col">
      <TriggerActions />
      <div className="mt-3 max-w-sm">
        <SearchInput
          placeholder={t('ui.text.searchTriggers')}
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
          <TriggerTable
            filteredTriggers={triggers}
            lang={lang}
            onToggleTrigger={toggleTrigger}
            onConfirmDelete={confirmDelete}
            onDetails={(trigger) => {
              setDetailsTrigger(trigger)
            }}
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

          <Sheet
            open={!!detailsTrigger}
            onOpenChange={(open) => !open && setDetailsTrigger(null)}
          >
            <SheetContent className="sm:max-w-lg overflow-y-auto gap-0">
              <SheetHeader>
                <SheetTitle>{t('ui.text.triggerDetails')}</SheetTitle>
                <SheetDescription className="sr-only">
                  {t('ui.text.triggerDetails')}
                </SheetDescription>
              </SheetHeader>
              {detailsTrigger && (
                <div className="grid gap-5 px-4 pb-6">
                  <div className="grid gap-1.5">
                    <Label className="text-muted-foreground">
                      {t('ui.text.type')}
                    </Label>
                    <Badge variant="outline" className="w-fit capitalize">
                      {detailsTrigger.type}
                    </Badge>
                  </div>
                  <div className="grid gap-1.5">
                    <Label className="text-muted-foreground">
                      {t('ui.text.status')}
                    </Label>
                    <span
                      className={`text-sm ${detailsTrigger.isActive ? 'text-green-600' : 'text-muted-foreground'}`}
                    >
                      {detailsTrigger.isActive
                        ? t('ui.text.active')
                        : t('ui.text.inactive')}
                    </span>
                  </div>
                  <div className="grid gap-1.5">
                    <Label className="text-muted-foreground">
                      {t('ui.text.flow')}
                    </Label>
                    <span className="text-sm">
                      {detailsTrigger.workflow?.title}
                    </span>
                  </div>
                  {isSchedule && (
                    <>
                      <div className="grid gap-1.5">
                        <Label className="text-muted-foreground">
                          {t('ui.text.cronPattern')}
                        </Label>
                        <Select value={editCron} onValueChange={setEditCron}>
                          <SelectTrigger className="w-full">
                            <SelectValue
                              placeholder={t('ui.text.selectTime')}
                            />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="*/15 * * * *">
                              {t('ui.text.every15Minutes')}
                            </SelectItem>
                            <SelectItem value="*/30 * * * *">
                              {t('ui.text.every30Minutes')}
                            </SelectItem>
                            <SelectItem value="0 * * * *">
                              {t('ui.text.everyHour')}
                            </SelectItem>
                            <SelectItem value="0 0 * * *">
                              {t('ui.text.everyDay')}
                            </SelectItem>
                          </SelectContent>
                        </Select>
                        <Input
                          value={editCron}
                          onChange={(e) => setEditCron(e.target.value)}
                          placeholder={t('ui.text.cronFormat')}
                          className="font-mono text-xs"
                        />
                      </div>
                      <div className="grid gap-1.5">
                        <Label className="text-muted-foreground">
                          {t('ui.text.timezone')}
                        </Label>
                        <Input
                          value={editTimezone}
                          onChange={(e) => setEditTimezone(e.target.value)}
                          placeholder="UTC"
                        />
                      </div>
                      <Button
                        size="sm"
                        className="w-fit"
                        disabled={!canSaveSchedule || isSavingSchedule}
                        onClick={saveSchedule}
                      >
                        {isSavingSchedule
                          ? t('ui.text.saving')
                          : t('ui.text.save')}
                      </Button>
                    </>
                  )}
                  <div className="grid gap-1.5">
                    <Label className="text-muted-foreground">
                      {t('ui.text.created')}
                    </Label>
                    <span className="text-sm">
                      <TimeAgoI18n date={detailsTrigger.createdAt} />
                    </span>
                  </div>
                  {showWebhookDetails && (
                    <>
                      {webhookSources.length > 0 && (
                        <div className="grid gap-2">
                          <Label className="text-muted-foreground">
                            {t('ui.text.webhookUrls')}
                          </Label>
                          <div className="grid gap-3">
                            {webhookSources.map((item: any) => (
                              <div key={item.source} className="grid gap-1">
                                <span className="text-xs font-medium flex items-center gap-1">
                                  {item.supportsVerification && (
                                    <Tooltip>
                                      <TooltipTrigger asChild>
                                        <ShieldCheck className="size-3.5 text-muted-foreground shrink-0" />
                                      </TooltipTrigger>
                                      <TooltipContent side="top">
                                        {t('ui.text.webhookSourceVerified')}
                                      </TooltipContent>
                                    </Tooltip>
                                  )}
                                  {t(item.name)}
                                </span>
                                <CopyField
                                  value={item.webhookUrl}
                                  inputClassName="h-8"
                                />
                              </div>
                            ))}
                          </div>
                          <p className="text-xs text-muted-foreground">
                            {t('ui.text.webhookUrlReplaceToken')}
                          </p>
                        </div>
                      )}

                      <div className="grid gap-2">
                        <Label
                          htmlFor="webhook-secret"
                          className="text-muted-foreground"
                        >
                          {t('ui.text.webhookSecret')}
                        </Label>
                        <div className="flex items-center gap-2">
                          <Input
                            id="webhook-secret"
                            type="password"
                            autoComplete="new-password"
                            className="h-8"
                            placeholder={t('ui.text.webhookSecretPlaceholder')}
                            value={editWebhookSecret}
                            onChange={(e) =>
                              setEditWebhookSecret(e.target.value)
                            }
                          />
                          <Button
                            size="sm"
                            disabled={
                              !editWebhookSecret.trim() || isSavingSecret
                            }
                            onClick={saveWebhookSecret}
                          >
                            {isSavingSecret
                              ? t('ui.text.saving')
                              : t('ui.text.save')}
                          </Button>
                          {detailsTrigger.hasWebhookSecret && (
                            <Button
                              size="sm"
                              variant="ghost"
                              className="text-destructive hover:text-destructive hover:bg-destructive/10"
                              onClick={() =>
                                confirmRemoveWebhookSecret(detailsTrigger.id)
                              }
                            >
                              {t('ui.text.remove')}
                            </Button>
                          )}
                        </div>
                        <p className="text-xs text-muted-foreground">
                          {detailsTrigger.hasWebhookSecret
                            ? t('ui.text.webhookSecretConfigured')
                            : t('ui.text.webhookSecretNotConfigured')}
                        </p>
                      </div>
                    </>
                  )}

                  <div className="flex flex-col gap-2 pt-4 border-t">
                    <Button
                      variant="outline"
                      className="justify-start gap-2"
                      onClick={async () => {
                        await toggleTrigger(detailsTrigger)
                        setDetailsTrigger(null)
                      }}
                    >
                      {detailsTrigger.isActive ? (
                        <>
                          <Pause className="size-4" />
                          {t('ui.text.pause')}
                        </>
                      ) : (
                        <>
                          <Play className="size-4" />
                          {t('ui.text.activate')}
                        </>
                      )}
                    </Button>
                    {showWebhookDetails && (
                      <Button
                        variant="outline"
                        className="justify-start gap-2"
                        disabled={isRegenerating}
                        onClick={() =>
                          confirmRegenerateToken(detailsTrigger.id)
                        }
                      >
                        <RefreshCcw className="size-4" />
                        {isRegenerating
                          ? t('ui.text.regenerating')
                          : t('ui.text.regenerateToken')}
                      </Button>
                    )}
                    <Button
                      variant="ghost"
                      className="justify-start gap-2 text-destructive hover:text-destructive hover:bg-destructive/10"
                      onClick={() => {
                        const id = detailsTrigger.id
                        setDetailsTrigger(null)
                        confirmDelete(id)
                      }}
                    >
                      <Trash2 className="size-4" />
                      {t('ui.text.delete')}
                    </Button>
                  </div>
                </div>
              )}
            </SheetContent>
          </Sheet>

          <Dialog
            open={!!regeneratedToken}
            onOpenChange={(o) => {
              if (!o) {
                setRegeneratedToken(null)
              }
            }}
          >
            <DialogContent className="sm:max-w-[550px]">
              <DialogHeader>
                <DialogTitle>
                  {t('ui.text.webhookTokenRegenerated')}
                </DialogTitle>
                <DialogDescription>
                  {t('ui.text.webhookTokenRegeneratedDescription')}
                </DialogDescription>
              </DialogHeader>
              {regeneratedToken && (
                <div className="flex flex-col gap-3">
                  <CopyField value={regeneratedToken} />
                  <p className="text-xs text-muted-foreground">
                    {t('ui.text.webhookTokenWarning')}
                  </p>
                </div>
              )}
              <DialogFooter>
                <Button size="sm" onClick={() => setRegeneratedToken(null)}>
                  {t('ui.text.done')}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      )}
    </div>
  )
}

export default OrganizationTriggersPage
