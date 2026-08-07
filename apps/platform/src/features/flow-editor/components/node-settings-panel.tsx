import { useState, useEffect, useMemo, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { usePlatformFeatures } from '@/features/platform'
import { useIconResolver } from '@/features/flow-editor/icons'
import { AdvancedDebouncedTextarea } from '@/components/shared/advanced-debounced-textarea'
import { SchemaFieldRenderer } from '@/features/flow-editor/components/schema-field-renderer'
import { DataSourcePopover } from '@/features/flow-editor/components/data-source-popover'
import { flowCategoryPreferences } from '@/lib/flow-categories'
import type {
  FlowNode,
  IconSource,
  NodeSettings,
  SchemaField,
} from '@/features/flow-editor/types'

const DEFAULT_SETTINGS: NodeSettings = {
  maxRetries: 0,
  retryDelay: 1000,
  timeout: 30000,
  continueOnError: false,
  disabled: false,
}

interface NodeSettingsPanelProps {
  node: FlowNode
  schema: SchemaField[]
  hasNaturalLanguage?: boolean
  triggerVariables?: string[]
  onApply: (updates: {
    instruction: string
    parameters: Record<string, unknown>
    nodeSettings: NodeSettings
  }) => void
  onCancel: () => void
}

export function NodeSettingsPanel({
  node,
  schema,
  hasNaturalLanguage,
  triggerVariables,
  onApply,
  onCancel,
}: NodeSettingsPanelProps) {
  const { t } = useTranslation()

  const [instruction, setInstruction] = useState(node.instruction || '')
  const [parameters, setParameters] = useState<Record<string, any>>(() => {
    const initial: Record<string, any> = {}
    const existing = node.parameters || {}
    for (const field of schema || []) {
      if (field.hidden) {
        continue
      }
      initial[field.key] = existing[field.key] ?? field.default ?? ''
    }
    return initial
  })
  const [nodeSettings, setNodeSettings] = useState<NodeSettings>(() => ({
    ...DEFAULT_SETTINGS,
    ...(node.settings || {}),
  }))

  // Reset draft when switching to a different node
  useEffect(() => {
    setInstruction(node.instruction || '')
    const initial: Record<string, any> = {}
    const existing = node.parameters || {}
    for (const field of schema || []) {
      if (field.hidden) {
        continue
      }
      initial[field.key] = existing[field.key] ?? field.default ?? ''
    }
    setParameters(initial)
    setNodeSettings({ ...DEFAULT_SETTINGS, ...(node.settings || {}) })
  }, [node.id])

  const handleSchemaChange = useCallback((key: string, value: any) => {
    setParameters((prev) => ({ ...prev, [key]: value }))
  }, [])

  const updateSetting = useCallback(
    <K extends keyof NodeSettings>(key: K, value: NodeSettings[K]) => {
      setNodeSettings((prev) => ({ ...prev, [key]: value }))
    },
    [],
  )

  const handleApply = () => {
    onApply({ instruction, parameters, nodeSettings })
  }

  const prefs = useMemo(() => {
    const category = node.data.category as
      keyof typeof flowCategoryPreferences | undefined
    if (!category || !flowCategoryPreferences[category]) {
      return { color: '#ffffff', icon: null, label: '' }
    }
    return flowCategoryPreferences[category]
  }, [node.data.category])

  const resolveIcon = useIconResolver()
  const Icon = resolveIcon(node.data?.icon as IconSource) ?? prefs.icon

  const visibleSchema = useMemo(
    () => (schema || []).filter((f) => !f.hidden),
    [schema],
  )

  const { functionCallingEnabled } = usePlatformFeatures()
  const showBasicTab = hasNaturalLanguage !== false && functionCallingEnabled
  const showAdvancedTab = visibleSchema.length > 0
  const defaultTab = showBasicTab
    ? 'basic'
    : showAdvancedTab
      ? 'advanced'
      : 'settings'

  return (
    <div
      className="absolute right-4 top-4 bottom-4 w-[400px] z-[4] bg-background border border-border rounded-xl shadow-lg flex flex-col overflow-hidden animate-in slide-in-from-right duration-200"
      onClick={(e) => e.stopPropagation()}
    >
      {/* Header */}
      <div className="flex items-center gap-3 px-4 py-3 border-b border-border shrink-0">
        <div className="bg-muted ring-border flex size-8 shrink-0 items-center justify-center rounded-lg ring-1">
          {Icon && <Icon className="size-5" />}
        </div>
        <div className="flex-1 min-w-0">
          <div className="text-sm font-medium truncate">
            {t(node.data.title ?? '')}
          </div>
          <div className="text-xs text-muted-foreground capitalize">
            {prefs.labelKey
              ? t(prefs.labelKey)
              : String(node.data.category || '')}
          </div>
        </div>
        <Button variant="ghost" size="icon-sm" onClick={onCancel}>
          <X className="size-4" />
        </Button>
      </div>

      {/* Body */}
      <div className="flex-1 overflow-y-auto px-4 py-3 flex flex-col">
        <Tabs
          defaultValue={defaultTab}
          key={defaultTab}
          className="w-full flex-1 flex flex-col"
        >
          <TabsList className="w-full">
            {showBasicTab && (
              <TabsTrigger value="basic" className="flex-1">
                {t('ui.text.basic')}
              </TabsTrigger>
            )}
            {showAdvancedTab && (
              <TabsTrigger value="advanced" className="flex-1">
                {t('ui.text.advanced')}
              </TabsTrigger>
            )}
            <TabsTrigger value="settings" className="flex-1">
              {t('ui.text.settings')}
            </TabsTrigger>
          </TabsList>

          {showBasicTab && (
            <TabsContent value="basic" className="mt-4 flex-1 flex flex-col">
              <div className="flex flex-col gap-2 flex-1">
                <div className="flex items-center justify-between">
                  <label className="text-sm font-medium">
                    {t('ui.text.manageWithNaturalLanguage')}
                  </label>
                  <DataSourcePopover
                    contextMenu={node.data.contextMenu ?? []}
                    triggerVariables={triggerVariables}
                    onSelect={(val) =>
                      setInstruction((prev) => `${prev || ''} ${val}`)
                    }
                  />
                </div>
                <div className="bg-muted/50 rounded-lg border border-border overflow-hidden flex-1">
                  <AdvancedDebouncedTextarea
                    className="h-full text-foreground placeholder:text-muted-foreground"
                    value={instruction}
                    onChange={setInstruction}
                    mapping={node.data.contextMenu}
                    placeholder={t('ui.text.writeYourInstructions')}
                  />
                </div>
              </div>
            </TabsContent>
          )}

          {showAdvancedTab && (
            <TabsContent value="advanced" className="mt-4">
              <div className="grid gap-4">
                {visibleSchema.map((field) => (
                  <SchemaFieldRenderer
                    key={field.key}
                    field={field}
                    value={parameters[field.key]}
                    onChange={handleSchemaChange}
                    contextMenu={node.data.contextMenu}
                    triggerVariables={triggerVariables}
                  />
                ))}
              </div>
            </TabsContent>
          )}

          <TabsContent value="settings" className="mt-4">
            <div className="grid gap-5">
              {/* Max Retries */}
              <div className="grid gap-1.5">
                <Label htmlFor="maxRetries">{t('ui.text.maxRetries')}</Label>
                <Input
                  id="maxRetries"
                  type="number"
                  min={0}
                  value={nodeSettings.maxRetries}
                  onChange={(e) =>
                    updateSetting('maxRetries', Number(e.target.value))
                  }
                />
                <p className="text-xs text-muted-foreground">
                  {t('ui.text.maxRetriesDescription')}
                </p>
              </div>

              {/* Retry Delay */}
              <div className="grid gap-1.5">
                <Label htmlFor="retryDelay">{t('ui.text.retryDelay')}</Label>
                <Input
                  id="retryDelay"
                  type="number"
                  min={0}
                  step={100}
                  value={nodeSettings.retryDelay}
                  onChange={(e) =>
                    updateSetting('retryDelay', Number(e.target.value))
                  }
                />
                <p className="text-xs text-muted-foreground">
                  {t('ui.text.retryDelayDescription')}
                </p>
              </div>

              {/* Timeout */}
              <div className="grid gap-1.5">
                <Label htmlFor="timeout">{t('ui.text.timeout')}</Label>
                <Input
                  id="timeout"
                  type="number"
                  min={0}
                  step={1000}
                  value={nodeSettings.timeout}
                  onChange={(e) =>
                    updateSetting('timeout', Number(e.target.value))
                  }
                />
                <p className="text-xs text-muted-foreground">
                  {t('ui.text.timeoutDescription')}
                </p>
              </div>

              {/* Continue on Error */}
              <div className="flex items-center justify-between gap-3 rounded-lg border border-border p-3">
                <div className="grid gap-0.5">
                  <Label htmlFor="continueOnError">
                    {t('ui.text.continueOnError')}
                  </Label>
                  <p className="text-xs text-muted-foreground">
                    {t('ui.text.continueOnErrorDescription')}
                  </p>
                </div>
                <Switch
                  id="continueOnError"
                  checked={nodeSettings.continueOnError}
                  onCheckedChange={(val) =>
                    updateSetting('continueOnError', val)
                  }
                />
              </div>

              {/* Disabled */}
              <div className="flex items-center justify-between gap-3 rounded-lg border border-border p-3">
                <div className="grid gap-0.5">
                  <Label htmlFor="disabled">{t('ui.text.disabled')}</Label>
                  <p className="text-xs text-muted-foreground">
                    {t('ui.text.disabledDescription')}
                  </p>
                </div>
                <Switch
                  id="disabled"
                  checked={nodeSettings.disabled}
                  onCheckedChange={(val) => updateSetting('disabled', val)}
                />
              </div>
            </div>
          </TabsContent>
        </Tabs>
      </div>

      {/* Footer */}
      <div className="flex items-center justify-end gap-2 px-4 py-3 border-t border-border shrink-0">
        <Button variant="outline" size="sm" onClick={onCancel}>
          {t('ui.text.cancel')}
        </Button>
        <Button size="sm" onClick={handleApply}>
          {t('ui.text.apply')}
        </Button>
      </div>
    </div>
  )
}
