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

    if (dragData.current.kind === 'note') {
      const noteNode = {
        id,
        data: {
          id,
          note: '',
          category: dragData.current.category,
        },
        parameters: {},
        nodeId: dragData.current.id,
        type: 'annotation',
        position,
      }
      setNodes((nds: FlowNode[]) =>
        preSave(nds.concat(noteNode as FlowNode), true),
      )
      return
    }

    const newNode = {
      id,
      data: {
        id,
        description: dragData.current.description,
        tags: dragData.current.tags,
        title: dragData.current.name,
        category: dragData.current.category,
        icon: dragData.current.icon,
        inputs: dragData.current.inputs,
        outputs: dragData.current.outputs,
      },
      instruction: undefined,
      parameters: {},
      origin: [0.5, 0.0],
      nodeId: dragData.current.id,
      type: 'core',
      position,
    }
    setNodes((nds: FlowNode[]) =>
      preSave(nds.concat(newNode as FlowNode), true),
    )
    nextStep()
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
  }
}
