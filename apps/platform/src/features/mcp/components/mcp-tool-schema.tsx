import { ChevronRight } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import { cn } from '@/lib/utils'
import type {
  McpToolSchema,
  McpToolSchemaProperty,
} from '@/features/mcp/services/mcp'

const propertyType = (prop: McpToolSchemaProperty): string => {
  if (prop.type === 'array') {
    const itemType =
      prop.items && typeof prop.items === 'object'
        ? (prop.items as { type?: string }).type
        : undefined
    return itemType ? `${itemType}[]` : 'array'
  }
  return prop.type ?? 'any'
}

type SchemaShape = {
  properties?: Record<string, McpToolSchemaProperty>
  required?: string[]
}

// Resolve the object that actually holds `properties` for a given node:
// objects expose them directly, arrays nest them under `items`.
const resolveShape = (
  node: McpToolSchema | McpToolSchemaProperty | undefined,
): SchemaShape | null => {
  if (!node || typeof node !== 'object') return null
  if (node.type === 'array' && node.items && typeof node.items === 'object') {
    return resolveShape(node.items as McpToolSchemaProperty)
  }
  if (node.properties && Object.keys(node.properties).length > 0) {
    return {
      properties: node.properties as Record<string, McpToolSchemaProperty>,
      required: node.required,
    }
  }
  return null
}

type SchemaFieldsProps = {
  shape: SchemaShape
  depth?: number
}

const SchemaFields = ({ shape, depth = 0 }: SchemaFieldsProps) => {
  const required = new Set(shape.required ?? [])
  const entries = Object.entries(shape.properties ?? {})

  return (
    <ul
      className={cn(
        'flex flex-col gap-1.5',
        depth > 0 && 'mt-1 border-l border-border pl-3',
      )}
    >
      {entries.map(([name, prop]) => {
        const nested = resolveShape(prop)
        return (
          <li key={name} className="flex flex-col gap-0.5">
            <div className="flex flex-wrap items-center gap-1.5">
              <span className="font-mono text-xs font-medium">{name}</span>
              {required.has(name) && (
                <span
                  aria-hidden
                  className="text-xs font-semibold text-destructive"
                >
                  *
                </span>
              )}
              <Badge variant="outline" className="font-mono text-[10px]">
                {propertyType(prop)}
              </Badge>
            </div>
            {prop.description && (
              <p className="text-[11px] leading-snug text-muted-foreground">
                {prop.description}
              </p>
            )}
            {nested && <SchemaFields shape={nested} depth={depth + 1} />}
          </li>
        )
      })}
    </ul>
  )
}

type McpToolSchemaSectionProps = {
  label: string
  schema: McpToolSchema
}

const McpToolSchemaSection = ({ label, schema }: McpToolSchemaSectionProps) => {
  const shape = resolveShape(schema)
  const count = shape ? Object.keys(shape.properties ?? {}).length : 0

  if (count === 0) return null

  return (
    <Collapsible className="rounded-md border bg-muted/30">
      <CollapsibleTrigger className="group/sch flex w-full items-center justify-between gap-2 px-2.5 py-1.5 text-left">
        <span className="flex items-center gap-1.5 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground">
          {label}
          <Badge variant="secondary" className="text-[10px]">
            {count}
          </Badge>
        </span>
        <ChevronRight className="size-3.5 shrink-0 text-muted-foreground transition-transform group-data-[state=open]/sch:rotate-90" />
      </CollapsibleTrigger>
      <CollapsibleContent>
        <div className="px-2.5 pt-0.5 pb-2.5">
          <SchemaFields shape={shape as SchemaShape} />
        </div>
      </CollapsibleContent>
    </Collapsible>
  )
}

McpToolSchemaSection.displayName = 'McpToolSchemaSection'

export { McpToolSchemaSection }
