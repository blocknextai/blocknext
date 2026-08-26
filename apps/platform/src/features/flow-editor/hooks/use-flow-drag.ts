import { useState, useRef } from 'react'
import { useReactFlow } from '@xyflow/react'
import type {
  FlowNode,
  IconSource,
  NodeHandle,
  NodeOutputType,
  SchemaField,
  SetFlowNodes,
} from '@/features/flow-editor/types'

interface PreviewPos {
  x: number
  y: number
}

interface DragItem {
  id: string
  kind?: string
  name: string
  description: string
  category: string
  icon?: IconSource
  inputs?: NodeHandle[]
  outputs?: NodeHandle[]
  subCategory?: string
  tags?: string[]
  isComingSoon?: boolean
  outputTypes?: NodeOutputType[]
  schema?: SchemaField[]
}

interface UseFlowDragOptions {
  apiNodes: DragItem[]
  setNodes: SetFlowNodes
  preSave: (data: FlowNode[], isNode: true) => FlowNode[]
  getId: () => string
  nextStep: () => void
}

export function useFlowDrag({
  apiNodes,
  setNodes,
  preSave,
  getId,
  nextStep,
}: UseFlowDragOptions) {
  const { screenToFlowPosition } = useReactFlow()
  const [dragging, setDragging] = useState(false)
  const [previewPos, setPreviewPos] = useState<PreviewPos>({ x: 0, y: 0 })
  const [previewData, setPreviewData] = useState<DragItem | null>(null)
  const dragData = useRef<DragItem | null>(null)

  const onMouseMove = (e: MouseEvent) => {
    setPreviewPos({ x: e.clientX, y: e.clientY })
  }

  const onMouseUp = (e: MouseEvent) => {
    setDragging(false)
    setPreviewData(null)
    const id = getId()
    const pane = document.querySelector(
      '.react-flow__pane.draggable',
    ) as HTMLElement
    if (pane) {
      pane.style.cursor = 'grab'
    }
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)

    if (!dragData.current) {
      return
    }
    const position = screenToFlowPosition({
      x: e.clientX + 100,
      y: e.clientY,
    })

    appendNode(dragData.current, id, position)
    nextStep()
  }

  const appendNode = (
    item: DragItem,
    id: string,
    position: { x: number; y: number },
  ) => {
    const node =
      item.kind === 'note'
        ? {
            id,
            data: { id, note: '', category: item.category },
            parameters: {},
            nodeId: item.id,
            type: 'annotation',
            position,
          }
        : {
            id,
            data: {
              id,
              description: item.description,
              tags: item.tags,
              title: item.name,
              category: item.category,
              icon: item.icon,
              inputs: item.inputs,
              outputs: item.outputs,
            },
            instruction: undefined,
            parameters: {},
            origin: [0.5, 0.0],
            nodeId: item.id,
            type: 'core',
            position,
          }

    setNodes((nds: FlowNode[]) => preSave(nds.concat(node as FlowNode), true))
  }

  const addNodeById = (
    nodeId: string,
    flowPosition?: { x: number; y: number },
  ) => {
    const item = apiNodes.find((node) => node.id === nodeId)
    if (!item || item.isComingSoon) {
      return false
    }

    if (flowPosition) {
      appendNode(item, getId(), flowPosition)
      return true
    }

    const pane = document.querySelector('.react-flow') as HTMLElement | null
    const rect = pane?.getBoundingClientRect()
    const position = screenToFlowPosition({
      x: rect ? rect.left + rect.width / 2 : window.innerWidth / 2,
      y: rect ? rect.top + rect.height / 2 : window.innerHeight / 2,
    })

    appendNode(item, getId(), position)
    return true
  }

  const startDrag = (e: React.MouseEvent, id: string) => {
    const indx = apiNodes.findIndex((node) => node.id === id)
    const item = apiNodes[indx]

    // Prevent dragging for coming soon nodes
    if (item.isComingSoon) {
      return
    }

    e.preventDefault()
    setDragging(true)
    dragData.current = item
    setPreviewData(item)
    setPreviewPos({ x: e.clientX, y: e.clientY })
    const pane = document.querySelector(
      '.react-flow__pane.draggable',
    ) as HTMLElement
    if (pane) {
      pane.style.cursor = 'grabbing'
    }
    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
  }

  return {
    dragging,
    previewPos,
    previewData,
    startDrag,
    addNodeById,
  }
}
