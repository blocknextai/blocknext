import { useState, useCallback, useRef, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { workflowsService } from '@/features/workflows'
import type {
  FlowEdge,
  FlowModel,
  FlowNode,
  SetFlowEdges,
  SetFlowNodes,
} from '@/features/flow-editor/types'

interface FlowData {
  title: string
  description: string
}

interface UseFlowSaveOptions {
  initialFlow: FlowModel | null
  nodes: FlowNode[]
  edges: FlowEdge[]
  setNodes: SetFlowNodes
  setEdges: SetFlowEdges
}

export function useFlowSave({
  initialFlow,
  nodes,
  edges,
  setNodes,
  setEdges,
}: UseFlowSaveOptions) {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { organizationId } = useParams()

  const savePayloadRef = useRef<{ nodes: FlowNode[]; edges: FlowEdge[] }>({
    nodes: [],
    edges: [],
  })
  const nodesRef = useRef(nodes)
  const edgesRef = useRef(edges)

  useEffect(() => {
    nodesRef.current = nodes
  }, [nodes])
  useEffect(() => {
    edgesRef.current = edges
  }, [edges])

  const [flowData, setFlowData] = useState<FlowData>({
    title: initialFlow?.title || t('ui.text.untitledWorkflow'),
    description: initialFlow?.description || '',
  })
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false)
  const [open, setOpen] = useState(false)
  const autosaveTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const updateFlowData = useCallback((key: string, value: string) => {
    setFlowData((prev) => ({ ...prev, [key]: value }))
  }, [])

  function preSave(data: FlowNode[], isNode: true): FlowNode[]
  function preSave(data: FlowEdge[], isNode: false): FlowEdge[]
  function preSave(
    data: FlowNode[] | FlowEdge[],
    isNode: boolean,
  ): FlowNode[] | FlowEdge[] {
    const ps = { ...savePayloadRef.current }
    if (isNode) {
      ps.nodes = data as FlowNode[]
      ps.edges = ps.edges.length > 0 ? ps.edges : edgesRef.current
    } else {
      ps.edges = data as FlowEdge[]
      ps.nodes = ps.nodes.length > 0 ? ps.nodes : nodesRef.current
    }
    savePayloadRef.current = ps

    if (autosaveTimeoutRef.current) {
      clearTimeout(autosaveTimeoutRef.current)
    }

    autosaveTimeoutRef.current = setTimeout(() => {
      setHasUnsavedChanges(false)
    }, 5000)

    setHasUnsavedChanges(true)

    return data
  }

  const stripNodesForSave = (nodeList: FlowNode[]) =>
    nodeList.map(({ data: _data, ...rest }) => rest)

  const createOrUpdateFlow = async () => {
    if (autosaveTimeoutRef.current) {
      clearTimeout(autosaveTimeoutRef.current)
    }
    setHasUnsavedChanges(false)

    const strippedNodes = stripNodesForSave(nodes)

    if (initialFlow?.id) {
      await workflowsService.update(
        initialFlow.organizationId,
        initialFlow.id,
        {
          title: flowData.title || initialFlow.title,
          description: flowData.description || initialFlow.description,
          isPinned: initialFlow.isPinned,
          nodes: strippedNodes,
          edges,
        },
      )
      return initialFlow
    } else {
      const data = {
        title: flowData.title,
        description: flowData.description,
        nodes: strippedNodes,
        edges,
      }
      const response = await workflowsService.create(organizationId, data)
      return response?.data
    }
  }

  const saveFlow = async () => {
    const flow = await createOrUpdateFlow()
    if (!initialFlow?.id && flow?.id) {
      navigate(`/organizations/${organizationId}/edit/${flow.id}`)
    }
  }

  const runFlow = async () => {
    const needsSave = !initialFlow?.id || hasUnsavedChanges
    const flow = needsSave ? await createOrUpdateFlow() : initialFlow
    const flowId = initialFlow?.id || flow?.id
    if (flowId) {
      navigate(`/organizations/${organizationId}/run/${flowId}`)
    }
  }

  const deleteLastSaved = () => {
    setNodes(initialFlow?.nodes || [])
    setEdges(initialFlow?.edges || [])
  }

  useEffect(() => {
    return () => {
      if (autosaveTimeoutRef.current) {
        clearTimeout(autosaveTimeoutRef.current)
      }
    }
  }, [])

  return {
    flowData,
    updateFlowData,
    saveFlow,
    runFlow,
    deleteLastSaved,
    preSave,
    hasUnsavedChanges,
    setHasUnsavedChanges,
    open,
    setOpen,
  }
}
