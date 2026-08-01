import { Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { DataSourcePopover } from '@/features/flow-editor/components/data-source-popover'
import { useTranslation } from 'react-i18next'

interface SchemaField {
  key: string
  type?: string
  title?: string
  description?: string
  default?: any
  enum?: Array<string | number>
  format?: string
  minimum?: number
  maximum?: number
  items?: any
  hidden?: boolean
  readOnly?: boolean
  writeOnly?: boolean
  required?: boolean
}

interface ContextMenuItem {
  label: string
  value: string
  isEditable?: boolean
}

interface SchemaFieldRendererProps {
  field: SchemaField
  value: any
  onChange: (key: string, value: any) => void
  contextMenu?: ContextMenuItem[]
  triggerVariables?: string[]
}

const isMultilineFormat = (format?: string) =>
  format === 'multiline' || format === 'textarea' || format === 'multi-line'

const inputTypeFor = (field: SchemaField) => {
  if (field.writeOnly) return 'password'
  if (field.format === 'email') return 'email'
  if (field.format === 'uri' || field.format === 'url') return 'url'
  return 'text'
}

const stringifyComplex = (val: any) => {
  if (val === undefined || val === null) return ''
  if (typeof val === 'string') return val
  try {
    return JSON.stringify(val, null, 2)
  } catch {
    return String(val)
  }
}

const PRIMITIVE_ITEM_TYPES = new Set(['string', 'number', 'integer', 'boolean'])

const isPrimitiveItemSchema = (items: any) => {
  if (!items || typeof items !== 'object') return true
  const type = items.type
  if (type === undefined) return true
  return PRIMITIVE_ITEM_TYPES.has(type)
}

const defaultForItemType = (type?: string) => {
  if (type === 'boolean') return false
  if (type === 'number' || type === 'integer') return 0
  return ''
}

interface ArrayItemRowProps {
  index: number
  value: any
  itemSchema: any
  isReadOnly: boolean
  onChange: (next: any) => void
  onRemove: () => void
  contextMenu?: ContextMenuItem[]
  triggerVariables?: string[]
}

function ArrayItemRow({
  index,
  value,
  itemSchema,
  isReadOnly,
  onChange,
  onRemove,
  contextMenu,
  triggerVariables,
}: ArrayItemRowProps) {
  const itemType: string = itemSchema?.type ?? 'string'
  const itemEnum: Array<string | number> | undefined = Array.isArray(
    itemSchema?.enum,
  )
    ? itemSchema.enum
    : undefined
  const itemFormat: string | undefined = itemSchema?.format
  const itemMultiline = isMultilineFormat(itemFormat)
  const showVariablePopover = !isReadOnly && !itemEnum && itemType === 'string'

  const handleInsert = (varValue: string) => {
    const current = typeof value === 'string' ? value : ''
    onChange(`${current} ${varValue}`)
  }

  const renderInput = () => {
    if (itemEnum && itemEnum.length > 0) {
      return (
        <Select
          value={String(value ?? '')}
          onValueChange={(val) => onChange(val)}
          disabled={isReadOnly}
        >
          <SelectTrigger className="w-full">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {itemEnum.map((opt) => {
              const v = String(opt)
              return (
                <SelectItem key={v} value={v}>
                  {v}
                </SelectItem>
              )
            })}
          </SelectContent>
        </Select>
      )
    }

    if (itemType === 'boolean') {
      return (
        <Switch
          checked={!!value}
          onCheckedChange={(val) => onChange(val)}
          disabled={isReadOnly}
        />
      )
    }

    if (itemType === 'number' || itemType === 'integer') {
      return (
        <Input
          type="number"
          value={value ?? ''}
          step={itemType === 'integer' ? 1 : undefined}
          min={itemSchema?.minimum}
          max={itemSchema?.maximum}
          onChange={(e) =>
            onChange(e.target.value === '' ? '' : Number(e.target.value))
          }
          readOnly={isReadOnly}
        />
      )
    }

    if (itemMultiline) {
      return (
        <Textarea
          value={value ?? ''}
          onChange={(e) => onChange(e.target.value)}
          readOnly={isReadOnly}
        />
      )
    }

    return (
      <Input
        type="text"
        value={value ?? ''}
        onChange={(e) => onChange(e.target.value)}
        readOnly={isReadOnly}
      />
    )
  }

  return (
    <div className="flex items-start gap-2">
      <div className="flex-1 min-w-0">{renderInput()}</div>
      {showVariablePopover && (
        <DataSourcePopover
          contextMenu={contextMenu || []}
          triggerVariables={triggerVariables}
          onSelect={handleInsert}
        />
      )}
      {!isReadOnly && (
        <Button
          type="button"
          variant="ghost"
          size="icon-sm"
          className="shrink-0 text-muted-foreground hover:text-destructive"
          onClick={onRemove}
          aria-label={`Remove item ${index + 1}`}
        >
          <Trash2 className="size-4" />
        </Button>
      )}
    </div>
  )
}

export function SchemaFieldRenderer({
  field,
  value,
  onChange,
  contextMenu,
  triggerVariables,
}: SchemaFieldRendererProps) {
  const { t } = useTranslation()

  if (field.hidden) return null

  const labelText = t(field.title || field.key, field.title || field.key)
  const isReadOnly = !!field.readOnly
  const hasEnum = Array.isArray(field.enum) && field.enum.length > 0
  const isComplex = field.type === 'array' || field.type === 'object'
  const currentValue = value ?? field.default ?? ''

  const isTextInsertable =
    !isReadOnly &&
    !hasEnum &&
    field.type !== 'boolean' &&
    field.type !== 'number' &&
    field.type !== 'integer' &&
    !isComplex

  const handleInsertVariable = (varValue: string) => {
    onChange(field.key, `${currentValue} ${varValue}`)
  }

  const renderField = () => {
    if (hasEnum) {
      return (
        <Select
          value={String(currentValue ?? '')}
          onValueChange={(val) => onChange(field.key, val)}
          disabled={isReadOnly}
        >
          <SelectTrigger className="w-full">
            <SelectValue placeholder={labelText} />
          </SelectTrigger>
          <SelectContent>
            {field.enum!.map((opt) => {
              const v = String(opt)
              return (
                <SelectItem key={v} value={v}>
                  {t(v, v)}
                </SelectItem>
              )
            })}
          </SelectContent>
        </Select>
      )
    }

    switch (field.type) {
      case 'boolean':
        return (
          <Switch
            checked={value ?? field.default ?? false}
            onCheckedChange={(val) => onChange(field.key, val)}
            disabled={isReadOnly}
          />
        )

      case 'number':
      case 'integer':
        return (
          <Input
            type="number"
            value={currentValue}
            min={field.minimum}
            max={field.maximum}
            step={field.type === 'integer' ? 1 : undefined}
            onChange={(e) => onChange(field.key, Number(e.target.value))}
            readOnly={isReadOnly}
          />
        )

      case 'array': {
        if (isPrimitiveItemSchema(field.items)) {
          const itemSchema = field.items ?? { type: 'string' }
          const itemType: string = itemSchema?.type ?? 'string'
          const items: any[] = Array.isArray(value)
            ? value
            : Array.isArray(field.default)
              ? field.default
              : []

          const updateItem = (index: number, next: any) => {
            const updated = items.slice()
            updated[index] = next
            onChange(field.key, updated)
          }

          const removeItem = (index: number) => {
            const updated = items.filter((_, i) => i !== index)
            onChange(field.key, updated)
          }

          const addItem = () => {
            onChange(field.key, [...items, defaultForItemType(itemType)])
          }

          return (
            <div className="grid gap-2">
              {items.length > 0 && (
                <div className="grid gap-2">
                  {items.map((item, index) => (
                    <ArrayItemRow
                      key={index}
                      index={index}
                      value={item}
                      itemSchema={itemSchema}
                      isReadOnly={isReadOnly}
                      onChange={(next) => updateItem(index, next)}
                      onRemove={() => removeItem(index)}
                      contextMenu={contextMenu}
                      triggerVariables={triggerVariables}
                    />
                  ))}
                </div>
              )}
              {!isReadOnly && (
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="justify-center"
                  onClick={addItem}
                >
                  <Plus className="size-4" />
                  {t('ui.text.addItem', 'Add Item')}
                </Button>
              )}
            </div>
          )
        }
        return (
          <Textarea
            value={stringifyComplex(value ?? field.default)}
            onChange={(e) => {
              const raw = e.target.value
              try {
                onChange(field.key, raw ? JSON.parse(raw) : undefined)
              } catch {
                onChange(field.key, raw)
              }
            }}
            readOnly={isReadOnly}
            className="font-mono text-xs min-h-32"
          />
        )
      }
      case 'object':
        return (
          <Textarea
            value={stringifyComplex(value ?? field.default)}
            onChange={(e) => {
              const raw = e.target.value
              try {
                onChange(field.key, raw ? JSON.parse(raw) : undefined)
              } catch {
                onChange(field.key, raw)
              }
            }}
            readOnly={isReadOnly}
            className="font-mono text-xs min-h-32"
          />
        )

      default:
        if (isMultilineFormat(field.format)) {
          return (
            <Textarea
              value={currentValue}
              onChange={(e) => onChange(field.key, e.target.value)}
              readOnly={isReadOnly}
            />
          )
        }
        return (
          <Input
            type={inputTypeFor(field)}
            value={currentValue}
            onChange={(e) => onChange(field.key, e.target.value)}
            readOnly={isReadOnly}
          />
        )
    }
  }

  return (
    <div className="grid gap-2">
      <div className="flex items-center justify-between">
        <Label
          htmlFor={field.key}
          className={
            field.type === 'boolean'
              ? 'flex items-center justify-between flex-1'
              : ''
          }
        >
          <span>
            {labelText}
            {field.required && (
              <span className="text-destructive">&nbsp;*</span>
            )}
          </span>
          {field.type === 'boolean' && renderField()}
        </Label>
        {isTextInsertable && (
          <DataSourcePopover
            contextMenu={contextMenu || []}
            triggerVariables={triggerVariables}
            onSelect={handleInsertVariable}
          />
        )}
      </div>
      {field.type !== 'boolean' && renderField()}
    </div>
  )
}
