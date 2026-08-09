import type { Dispatch, SetStateAction } from 'react'
import type { Node, Edge } from '@xyflow/react'

export interface SchemaField {
  key: string
  type?: string
  title?: string
  description?: string
  default?: unknown
  enum?: Array<string | number>
  format?: string
  minimum?: number
  maximum?: number
  items?: unknown
  hidden?: boolean
  readOnly?: boolean
  writeOnly?: boolean
  required?: boolean
}

export interface NodeOutputType {
  key: string
  description?: string
  isEditable?: boolean
}

export interface NodeIcon {
  light?: string
  dark?: string
}

export interface NodeEngineNode {
  id?: string
  kind?: string
  name?: string
  description?: string
  icon: NodeIcon
  inputSchema?: unknown
  outputSchema?: unknown
  categories?: string[]
  subCategories?: string[]
  tags?: string[]
  supportedCredentials?: string[]
  disabled?: boolean
  isComingSoon?: boolean
  hasNaturalLanguage?: boolean
}

export interface NodeHandle {
  key: string
  label?: string
}

export type IconSource = { brand?: string; glyph?: string } | null | undefined

export interface ResolvedNode {
  id: string
  kind?: string
  name: string
  description: string
  icon: IconSource
  inputs: NodeHandle[]
  outputs: NodeHandle[]
  category: string
  subCategory: string
  provider: string
  tags: string[]
  supportedCredentials?: string[]
  schema: SchemaField[]
  outputTypes: NodeOutputType[]
  disabled?: boolean
  isComingSoon?: boolean
  hasNaturalLanguage?: boolean
}

export interface ContextMenuItem {
  label: string
  value: string
  isEditable?: boolean
}

export interface NodeSettings {
  maxRetries: number
  retryDelay: number
  timeout: number
  continueOnError: boolean
  disabled: boolean
}

export interface FlowNodeData extends Record<string, unknown> {
  id?: string
  title?: string
  description?: string
  category?: string
  subCategory?: string
  tags?: string[]
  contextMenu?: ContextMenuItem[]
  note?: string
}

export type FlowNode = Node<FlowNodeData> & {
  nodeId?: string
  instruction?: string
  parameters?: Record<string, unknown>
  settings?: Partial<NodeSettings>
  origin?: [number, number]
}
export type FlowEdge = Edge

export interface FlowModel {
  id?: string
  organizationId?: string
  title?: string
  description?: string
  isPinned?: boolean
  nodes?: FlowNode[]
  edges?: FlowEdge[]
}

export interface RendererNode {
  id: string
  label: string
  description: string
  type: string
  icon?: IconSource
  inputs?: NodeHandle[]
  outputs?: NodeHandle[]
  subCategory?: string
  actions: string[]
  tags?: string[]
  isComingSoon: boolean
}

export type RendererNodeSubMap = Record<string, RendererNode[]>
export type RendererNodeMap = Record<string, RendererNodeSubMap>

export type SetFlowNodes = Dispatch<SetStateAction<FlowNode[]>>
export type SetFlowEdges = Dispatch<SetStateAction<FlowEdge[]>>
