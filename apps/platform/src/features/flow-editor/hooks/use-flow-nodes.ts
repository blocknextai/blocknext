import { useCallback, useEffect, useMemo, useRef } from 'react'
import useSWR from 'swr'
import nodeEngineService from '@/features/flow-editor/services/node-engine'
import { useTranslation } from 'react-i18next'
import type {
  ContextMenuItem,
  FlowEdge,
  FlowNode,
  IconSource,
  NodeEngineNode,
  NodeOutputType,
  RendererNode,
  RendererNodeMap,
  ResolvedNode,
  SchemaField,
} from '@/features/flow-editor/types'

const normalizeKey = (value?: string) =>
  (value ?? '').toLowerCase().replace(/\s+/g, '')

const inputSchemaToFields = (schema: any): SchemaField[] => {
  if (!schema?.properties) {
    return []
  }
  const required = new Set<string>(schema.required ?? [])
  return Object.entries(schema.properties).map(
    ([key, value]: [string, any]) => ({
      key,
      type: value.type,
      title: value.title,
      description: value.description,
      default: value.default,
      enum: value.enum,
      format: value.format,
      minimum: value.minimum,
      maximum: value.maximum,
      items: value.items,
      hidden: !!value.hidden,
      readOnly: !!value.readOnly,
      writeOnly: !!value.writeOnly,
      required: required.has(key),
    }),
  )
}

const collectOutputPaths = (
  node: any,
  prefix: string,
  acc: NodeOutputType[],
) => {
  if (!node || typeof node !== 'object') {
    return
  }

  if (node.type === 'object' && node.properties) {
    for (const [key, value] of Object.entries<any>(node.properties)) {
      const path = prefix ? `${prefix}.${key}` : key
      collectOutputPaths(value, path, acc)
    }
    return
  }

  if (node.type === 'array') {
    const items = node.items
    if (
      items &&
      items.type === 'object' &&
      items.properties &&
      Object.keys(items.properties).length > 0
    ) {
      for (const [key, value] of Object.entries<any>(items.properties)) {
        const path = prefix ? `${prefix}.${key}` : key
        collectOutputPaths(value, path, acc)
      }
      return
    }
    acc.push({ key: prefix, description: node.description })
    return
  }

  acc.push({ key: prefix, description: node.description })
}

const outputSchemaToTypes = (schema: any): NodeOutputType[] => {
  if (!schema || typeof schema !== 'object') {
    return []
  }
  const acc: NodeOutputType[] = []
  collectOutputPaths(schema, '', acc)
  return acc
}

const buildNodeList = (nodes: ResolvedNode[]): RendererNodeMap => {
  const r: RendererNodeMap = {}
  for (const item of nodes) {
    const subCategory = item.subCategory ?? ''
    if (!r[item.category]) {
      r[item.category] = {}
    }
    if (!r[item.category][subCategory]) {
      r[item.category][subCategory] = []
    }
    const labels: string[] = []
    if (item.schema && item.schema.length > 0) {
      for (const field of item.schema) {
        if (Array.isArray(field.enum)) {
          for (const opt of field.enum) {
            labels.push(String(opt))
          }
        }
      }
    }
    const entry: RendererNode = {
      id: item.id,
      label: item.name,
      description: item.description,
      type: item.category,
      icon: item.icon,
      inputs: item.inputs,
      outputs: item.outputs,
      subCategory: item.subCategory,
      actions: labels,
      tags: item.tags,
      isComingSoon: item.isComingSoon ?? false,
    }
    r[item.category][subCategory].push(entry)
  }
  return r
}

const buildProviderList = (nodes: ResolvedNode[]) => {
  const groups = new Map<
    string,
    { key: string; label: string; icon: IconSource; nodes: RendererNode[] }
  >()

  for (const item of nodes) {
    const label = item.provider || item.subCategory || ''
    const key = item.subCategory || normalizeKey(label)
    if (!groups.has(key)) {
      groups.set(key, {
        key,
        label,
        icon: { brand: item.icon?.brand },
        nodes: [],
      })
    }
    groups.get(key)!.nodes.push({
      id: item.id,
      label: item.name,
      description: item.description,
      type: item.category,
      icon: item.icon,
      inputs: item.inputs,
      outputs: item.outputs,
      subCategory: item.subCategory,
      actions: [],
      tags: item.tags,
      isComingSoon: item.isComingSoon ?? false,
    })
  }

  return [...groups.values()].sort((a, b) => a.label.localeCompare(b.label))
}

const NODES_KEY = 'node-engine-nodes'
const TRIGGER_VARIABLES_KEY = 'node-engine-trigger-variables'

export function useNodeEngineNodes() {
  const { data, isLoading } = useSWR<NodeEngineNode[]>(NODES_KEY, async () => {
    const response = await nodeEngineService.getNodes()
    return (response.data as NodeEngineNode[]) || []
  })

  const apiNodes = useMemo<ResolvedNode[]>(
    () =>
      (data ?? []).map((item) => ({
        id: item.id ?? '',
        kind: item.kind,
        name: item.name ?? '',
        description: item.description ?? '',
        icon: item.icon,
        inputs: item.inputs,
        outputs: item.outputs,
        category: normalizeKey(item.categories?.[0]),
        subCategory: normalizeKey(item.subCategories?.[0]),
        provider: item.subCategories?.[0] ?? '',
        tags: item.tags ?? [],
        supportedCredentials: item.supportedCredentials,
        schema: inputSchemaToFields(item.inputSchema),
        outputTypes: outputSchemaToTypes(item.outputSchema),
        disabled: item.disabled,
        isComingSoon: item.isComingSoon,
        hasNaturalLanguage: item.hasNaturalLanguage,
      })),
    [data],
  )
  const nodeList = useMemo(() => buildNodeList(apiNodes), [apiNodes])
  const providerList = useMemo(() => buildProviderList(apiNodes), [apiNodes])

  return { apiNodes, nodeList, providerList, isLoading }
}

export function useTriggerVariables() {
  const { data, isLoading } = useSWR<string[]>(
    TRIGGER_VARIABLES_KEY,
    async () => {
      const response = await nodeEngineService.getTriggerVariables()
      if (response.isSuccess && response.data) {
        return response.data
      }
      return []
    },
  )
  return { triggerVariables: data ?? [], isLoading }
}

interface SourceRef {
  nodeId?: string
  id: string
  data?: {
    title?: string
    contextMenu?: ContextMenuItem[]
  }
}

interface GraphSource {
  nodes: Array<FlowNode & { nodeId?: string }>
  edges: FlowEdge[]
}

export function useFlowNodes(apiNodes: ResolvedNode[] = []) {
  const { t } = useTranslation()
  const apiOriginRef = useRef<ResolvedNode[]>(apiNodes)

  useEffect(() => {
    apiOriginRef.current = apiNodes
  }, [apiNodes])

  const createContextNode = useCallback(
    (
      node: FlowNode & { nodeId?: string },
      source?: GraphSource,
      from?: SourceRef,
      apiNodesArg?: ResolvedNode[],
    ): ContextMenuItem[] | undefined => {
      if (node.id === '0') {
        return undefined
      }
      let contextMenu: ContextMenuItem[] = []
      const ref = node.id

      if (from) {
        const sourceNode = from
        const apiNode = apiOriginRef.current.find(
          (n) => n.id === sourceNode.nodeId,
        )
        if (apiNode === undefined) {
          return contextMenu
        }

        const outputs = apiNode.outputTypes ?? []
        for (let j = 0; j < outputs.length; j++) {
          const outputType = outputs[j]
          contextMenu.push({
            label: `${t(sourceNode.data?.title || apiNode.name)} ${outputType.key}`,
            value: `$${apiNode.id}_${sourceNode.id}.${outputType.key}`,
            isEditable: outputType.isEditable ?? false,
          })
        }
        if (
          sourceNode.data?.contextMenu &&
          sourceNode.data.contextMenu.length > 0
        ) {
          contextMenu = contextMenu.concat(sourceNode.data.contextMenu)
        }
        if (node.data?.contextMenu) {
          contextMenu = [...node.data.contextMenu, ...contextMenu]
        }
        const cNew: ContextMenuItem[] = []
        const seenValues = new Set<string>()
        contextMenu.forEach((item) => {
          if (!seenValues.has(item.value)) {
            seenValues.add(item.value)
            cNew.push(item)
          }
        })
        contextMenu = cNew
      } else if (source) {
        const currentApiNodes = apiNodesArg || apiOriginRef.current
        const nodeName = (
          node: { title?: string; data?: { title?: string } },
          fallback: string,
        ) => node.title || node.data?.title || fallback

        const directParents = source.edges.filter((e) => e.target === ref)
        if (directParents.length === 1 && directParents[0].source !== '0') {
          const parentNode = source.nodes.find(
            (n) => n.id === directParents[0].source,
          )
          const parentApiNode = currentApiNodes?.find(
            (n) => n.id === parentNode?.nodeId,
          )
          for (const outputType of parentApiNode?.outputTypes ?? []) {
            contextMenu.push({
              label: `${t('ui.text.input', 'Input')} ${outputType.key}`,
              value: `$input.${outputType.key}`,
              isEditable: outputType.isEditable ?? false,
            })
          }
        }

        const ancestors: string[] = []
        const seenNodes = new Set<string>([ref])
        const queue = [ref]

        while (queue.length > 0) {
          const current = queue.shift()!
          for (const edge of source.edges) {
            if (edge.target !== current || seenNodes.has(edge.source)) {
              continue
            }
            seenNodes.add(edge.source)
            queue.push(edge.source)
            if (edge.source !== '0') {
              ancestors.push(edge.source)
            }
          }
        }

        for (const ancestorID of ancestors) {
          const sourceNode = source.nodes.find((n) => n.id === ancestorID)
          if (!sourceNode) {
            continue
          }
          const apiNode = currentApiNodes?.find(
            (n) => n.id === sourceNode.nodeId,
          )
          if (apiNode === undefined) {
            continue
          }

          const outputs = apiNode.outputTypes ?? []
          for (let j = 0; j < outputs.length; j++) {
            const outputType = outputs[j]
            contextMenu.push({
              label: `${t(nodeName(sourceNode, apiNode.name))} ${outputType.key}`,
              value: `$${apiNode.id}_${sourceNode.id}.${outputType.key}`,
              isEditable: outputType.isEditable ?? false,
            })
          }
        }
      }

      const seenValues = new Set<string>()
      const cNew: ContextMenuItem[] = []
      contextMenu.forEach((item) => {
        if (!seenValues.has(item.value)) {
          seenValues.add(item.value)
          cNew.push(item)
        }
      })
      return cNew
    },
    [t],
  )

  return { createContextNode }
}
