import coreNode from '@/features/flow-editor/nodes/core-node'
import coreEdge from '@/features/flow-editor/nodes/core-edge'
import coreGroup from '@/features/flow-editor/nodes/core-group'
import { useBreadcrumbStore } from '@/features/flow-editor/stores/breadcrumb-store'
import { useThemeStore } from '@/stores/theme-store'
import { useCallback, useMemo, useRef, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import WelcomeTour from '@/features/tour/components/welcome-tour'

import {
  ReactFlow,
  useNodesState,
  useEdgesState,
  addEdge,
  useReactFlow,
  useNodesInitialized,
  ReactFlowProvider,
  Background,
  BackgroundVariant,
  applyNodeChanges,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { Loading } from '@/components/shared/loading'
import { EdgeGradients } from '@/features/flow-editor/nodes/core/edge-graidents'
import { Annotation } from '@/features/flow-editor/nodes/annotation'
import { useParams } from 'react-router'
import AIChatSheet from '@/features/generation/components/ai-chat-sheet'

import { useFlowDrag } from '@/features/flow-editor/hooks/use-flow-drag'
import {
  useFlowNodes,
  useNodeEngineNodes,
  useTriggerVariables,
} from '@/features/flow-editor/hooks/use-flow-nodes'
import { useFlowSave } from '@/features/flow-editor/hooks/use-flow-save'
import { useNavigationGuard } from '@/features/flow-editor/hooks/use-navigation-guard'

import { FlowSidebar } from '@/features/flow-editor/components/flow-sidebar'
import { FlowToolbar } from '@/features/flow-editor/components/flow-toolbar'
import { FlowSaveDialog } from '@/features/flow-editor/components/flow-save-dialog'
import { FlowNavigationGuard } from '@/features/flow-editor/components/flow-navigation-guard'
import { FlowDragPreview } from '@/features/flow-editor/components/flow-drag-preview'
import { NodeSettingsPanel } from '@/features/flow-editor/components/node-settings-panel'
import { FlowViewControls } from '@/features/flow-editor/components/flow-view-controls'
import { FlowApiSheet } from '@/features/flow-editor/components/flow-api-sheet'
import { FlowNodesProvider } from '@/features/flow-editor/contexts/flow-nodes-context'

const defaultNodes = [
  {
    id: '0',
    position: { x: 0, y: 0 },
    data: {
      id: '0',
      title: 'ui.text.starterNode',
      description: 'ui.text.starterNodeDescription',
      hide_input: true,
      tags: ['starter'],
      category: 'system',
      inputs: [],
      outputs: [{ key: 'out' }],
    },
    type: 'core',
    nodeId: 'system_starter',
    deletable: false,
  },
]

const nodeTypes = {
  core: coreNode,
  starter: coreNode,
  group: coreGroup,
  annotation: Annotation,
}
const edgeTypes = {
  core: coreEdge,
}

const defaultEdgeOptions = {
  type: 'core',
  markerEnd: 'edge-circle',
}

const nodeOrigin = [0.5, 0]

const fitViewOptions = { padding: 2 }

const FlowCanvas = ({
  initialFlow,
  defaultSidebarOpen = false,
  previewMode = false,
  generatedFlow = null,
  defaultChatOpen = false,
}) => {
  const { setFlowName } = useBreadcrumbStore()
  const mode = useThemeStore((s) => s.mode)
  const themeMode = useMemo(() => {
    if (mode === 'system') {
      return window.matchMedia('(prefers-color-scheme: dark)').matches
        ? 'dark'
        : 'light'
    }
    return mode
  }, [mode])

  const flowWrapper = useRef<HTMLDivElement | null>(null)
  const idRef = useRef(1)
  const nextStepRef = useRef<(() => void) | null>(null)
  const getId = useCallback(() => `${idRef.current++}`, [])

  // Initialize with empty list when loading an existing flow; setFunctions
  // populates nodes with hydrated `data` once apiNodes are ready. Without this,
  // raw backend nodes (no `data` field) would render and crash CoreNode.
  const [nodes, setNodes] = useNodesState(initialFlow ? [] : defaultNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState([])
  const [selectedCategory, setSelectedCategory] = useState<any>(null)
  const reactFlow = useReactFlow()
  const [locked, setLocked] = useState(false)
  const [framed, setFramed] = useState(!initialFlow?.nodes?.length)
  const [sidebarOpen, setSidebarOpen] = useState(defaultSidebarOpen)
  const [chatOpen, setChatOpen] = useState(defaultChatOpen)
  const [apiSheetOpen, setApiSheetOpen] = useState(false)
  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null)
  const { organizationId } = useParams()
  const { t } = useTranslation()

  // Hooks
  const { apiNodes, providerList, isLoading } = useNodeEngineNodes()
  const { triggerVariables } = useTriggerVariables()
  const { createContextNode } = useFlowNodes(apiNodes)

  const {
    flowData,
    updateFlowData,
    saveFlow,
    runFlow,
    preSave,
    hasUnsavedChanges,
    setHasUnsavedChanges,
    open,
    setOpen,
  } = useFlowSave({
    initialFlow,
    nodes,
    edges,
  })

  const {
    showNavigationDialog,
    setShowNavigationDialog,
    handleNavigationConfirm,
    handleNavigationCancel,
  } = useNavigationGuard({
    hasUnsavedChanges,
    previewMode,
  })

  const { dragging, previewPos, previewData, startDrag, addNodeById } =
    useFlowDrag({
      apiNodes,
      setNodes,
      preSave,
      getId,
      nextStep: () => nextStepRef.current?.(),
    })

  // Search handler
  const handleSearch = useCallback(
    (e) => {
      const query = e.target.value.toLowerCase()
      const filtered = apiNodes
        .filter((node) => {
          const tags = node.tags.join(' ')
          const name = t(node.name)
          return (
            tags.toLowerCase().includes(query) ||
            name.toLowerCase().includes(query)
          )
        })
        .map((node) => ({
          id: node.id,
          label: node.name,
          description: node.description,
          type: node.category,
          icon: node.icon,
          subCategory: node.subCategory,
          tags: node.tags,
          isComingSoon: node.isComingSoon || false,
        }))
      setSelectedCategory(
        query
          ? {
              label: t('ui.text.searchResults'),
              nodes: filtered,
              isSearch: true,
            }
          : null,
      )
    },
    [apiNodes, t],
  )

  const buildNodeData = useCallback(
    (node: any) => {
      const apiNode = apiNodes.find((a) => a.id === node.nodeId)
      return {
        id: node.id,
        title: node.title || apiNode?.name || node.nodeId,
        catalogTitle: apiNode?.name || node.nodeId,
        description: apiNode?.description || '',
        tags: apiNode?.tags || [],
        category: apiNode?.category || '',
        icon: apiNode?.icon,
        inputs: apiNode?.inputs,
        outputs: apiNode?.outputs,
        handleLayout: node.handleLayout,
        note: node.parameters?.note,
      }
    },
    [apiNodes],
  )

  const setFunctions = useCallback(() => {
    const sourceNodes = initialFlow.nodes || []
    const eds = initialFlow.edges || []
    const lastNode = sourceNodes[sourceNodes.length - 1]
    if (lastNode) {
      idRef.current = parseInt(lastNode.id) + 1
    }
    const nds = sourceNodes.map((n) => {
      const data = buildNodeData(n)
      data.contextMenu = createContextNode(n, {
        nodes: sourceNodes,
        edges: eds,
      })
      return { ...n, data }
    })
    setNodes(nds)
    setEdges(eds)
  }, [
    initialFlow,
    apiNodes,
    buildNodeData,
    createContextNode,
    setNodes,
    setEdges,
  ])

  //MARK: flow func
  const onNodesChange = useCallback(
    (changes) => {
      setNodes((nds) => {
        const updatedNodes = applyNodeChanges(changes, nds)

        const hasSignificantChange = changes.some(
          (change) =>
            change.type === 'position' ||
            change.type === 'add' ||
            change.type === 'remove',
        )

        if (hasSignificantChange) {
          return preSave(updatedNodes, true)
        }
        return updatedNodes
      })
    },
    [setNodes, preSave],
  )

  const onConnect = useCallback(
    (params) => setEdges((eds) => preSave(addEdge(params, eds), false)),
    [setEdges, preSave],
  )

  const onConnectEnd = useCallback(
    (_event, connectionState) => {
      if (!connectionState.isValid) {
        return
      }
      const { fromNode, toNode } = connectionState
      setNodes((nodes) =>
        nodes.map((n) => {
          if (n.id === toNode.id) {
            const cm = createContextNode(n, undefined, fromNode)
            return {
              ...n,
              data: {
                ...n.data,
                contextMenu: cm,
              },
            }
          }
          return n
        }),
      )
    },
    [setNodes, createContextNode],
  )

  // Effects
  useEffect(() => {
    setFlowName(initialFlow?.title || t('ui.text.untitledWorkflow'))
    idRef.current = initialFlow?.nodes?.length || 1
  }, [])

  const nodesInitialized = useNodesInitialized()
  const framedRef = useRef(false)
  useEffect(() => {
    if (framedRef.current || !nodesInitialized || nodes.length === 0) {
      return
    }
    framedRef.current = true
    reactFlow.fitView({ padding: 0.2 })
    setFramed(true)
  }, [nodesInitialized, nodes.length, reactFlow])

  useEffect(() => {
    if (initialFlow) {
      setFunctions()
      return
    }
    if (apiNodes.length === 0) {
      return
    }
    setNodes((current) =>
      current.map((n) => ({ ...n, data: { ...n.data, ...buildNodeData(n) } })),
    )
  }, [initialFlow, apiNodes])

  const focusCanvas = useCallback(
    (nodeId?: string) => {
      if (!nodeId) {
        reactFlow.fitView({ padding: 0.2, duration: 300 })
        return
      }

      const node = reactFlow.getNode(nodeId)
      const element = document.querySelector(
        `.react-flow__node[data-id="${nodeId}"]`,
      )
      const pane = document.querySelector('.react-flow')
      if (!node || !element || !pane) {
        return
      }

      const nodeRect = element.getBoundingClientRect()
      const paneRect = pane.getBoundingClientRect()
      const margin = 48
      const isVisible =
        nodeRect.top >= paneRect.top + margin &&
        nodeRect.left >= paneRect.left + margin &&
        nodeRect.bottom <= paneRect.bottom - margin &&
        nodeRect.right <= paneRect.right - margin

      if (isVisible) {
        return
      }

      reactFlow.setCenter(
        node.position.x,
        node.position.y + (node.measured?.height ?? 0) / 2,
        { zoom: reactFlow.getZoom(), duration: 300 },
      )
    },
    [reactFlow],
  )

  const getFlowState = useCallback(
    () => ({ nodes: reactFlow.getNodes(), edges: reactFlow.getEdges() }),
    [reactFlow],
  )

  const openNodeSettings = useCallback((nodeId: string) => {
    setSelectedNodeId(nodeId)
    setChatOpen(false)
  }, [])

  const connectNodes = useCallback(
    (sourceId: string, targetId: string) => {
      const fromNode = reactFlow.getNode(sourceId)
      const toNode = reactFlow.getNode(targetId)
      if (!fromNode || !toNode) {
        return false
      }

      const sourceHandle = (fromNode.data as any)?.outputs?.[0]?.key ?? 'out'
      const targetHandle = (toNode.data as any)?.inputs?.[0]?.key ?? 'in'

      onConnect({
        source: sourceId,
        target: targetId,
        sourceHandle,
        targetHandle,
      })
      setNodes((current) =>
        current.map((n) =>
          n.id === targetId
            ? {
                ...n,
                data: {
                  ...n.data,
                  contextMenu: createContextNode(n, undefined, fromNode),
                },
              }
            : n,
        ),
      )

      return true
    },
    [reactFlow, onConnect, setNodes, createContextNode],
  )

  const selectGroupForNode = useCallback(
    (nodeId: string) => {
      const provider = providerList.find((group: any) =>
        group.nodes?.some((node: any) => node.id === nodeId),
      )
      if (!provider) {
        return false
      }

      setSidebarOpen(true)
      setSelectedCategory({ label: provider.label, nodes: provider.nodes })
      return true
    },
    [providerList, setSidebarOpen, setSelectedCategory],
  )

  const revertFlow = useCallback(() => {
    if (initialFlow) {
      setFunctions()
    } else {
      setNodes(
        defaultNodes.map((n) => ({
          ...n,
          data: { ...n.data, ...buildNodeData(n) },
        })),
      )
      setEdges([])
    }
    setHasUnsavedChanges(false)
  }, [
    initialFlow,
    setFunctions,
    setNodes,
    setEdges,
    buildNodeData,
    setHasUnsavedChanges,
  ])

  const applyGeneratedFlow = useCallback(
    (flow) => {
      if (!flow || apiNodes.length === 0) {
        return
      }
      const eds = flow.edges.map((e) => ({ ...e }))
      const nds = flow.nodes.map((n) => {
        const data = buildNodeData(n)
        data.contextMenu = createContextNode(n, {
          nodes: flow.nodes,
          edges: eds,
        })
        return { ...n, data }
      })
      const lastNode = nds[nds.length - 1]
      if (lastNode) {
        idRef.current = parseInt(lastNode.id) + 1
      }
      setNodes(nds)
      setEdges(eds)
      setHasUnsavedChanges(true)
    },
    [apiNodes, setNodes, setEdges],
  )

  useEffect(() => {
    if (generatedFlow) {
      applyGeneratedFlow(generatedFlow)
    }
  }, [generatedFlow, applyGeneratedFlow])

  const handleApplyWorkflow = useCallback(
    (data) => {
      applyGeneratedFlow(data)
      setChatOpen(false)
    },
    [applyGeneratedFlow],
  )

  // Node settings panel
  const onNodeClick = useCallback((_event, node) => {
    if (node.type === 'annotation') {
      return
    }
    setSelectedNodeId(node.id)
    setChatOpen(false)
  }, [])

  const onPaneClick = useCallback(() => {
    setSelectedNodeId(null)
  }, [])

  const selectedNode = useMemo(
    () => (selectedNodeId ? nodes.find((n) => n.id === selectedNodeId) : null),
    [selectedNodeId, nodes],
  )

  const selectedNodeSchema = useMemo(() => {
    if (!selectedNode) {
      return []
    }
    const apiNode = apiNodes.find((a) => a.id === selectedNode.nodeId)
    return apiNode?.schema || []
  }, [selectedNode, apiNodes])

  const selectedNodeHasNaturalLanguage = useMemo(() => {
    if (!selectedNode) {
      return undefined
    }
    const apiNode = apiNodes.find((a) => a.id === selectedNode.nodeId)
    return apiNode?.hasNaturalLanguage
  }, [selectedNode, apiNodes])

  const handleApplyNodeSettings = useCallback(
    (updates) => {
      setNodes((nds) =>
        preSave(
          nds.map((n) => {
            if (n.id === selectedNodeId) {
              return {
                ...n,
                instruction: updates.instruction,
                parameters: updates.parameters,
                settings: updates.nodeSettings,
              }
            }
            return n
          }),
          true,
        ),
      )
      setSelectedNodeId(null)
    },
    [selectedNodeId, setNodes, preSave],
  )

  useEffect(() => {
    document
      .querySelector('.drawer.w-full.h-full')
      ?.classList.toggle('opened', sidebarOpen)
  }, [sidebarOpen])

  const handleSetChatOpen = useCallback((v) => {
    setChatOpen(v)
    if (v) {
      setSelectedNodeId(null)
    }
  }, [])

  const nextStepInit = (fn) => {
    nextStepRef.current = fn
  }

  if (isLoading) {
    return <Loading />
  }

  return (
    <FlowNodesProvider value={{ setNodes, setEdges }}>
      <div
        className="w-full h-full relative"
        ref={flowWrapper}
        onContextMenu={(e) => e.preventDefault()}
      >
        <ReactFlow
          colorMode={themeMode}
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          nodeTypes={nodeTypes}
          edgeTypes={edgeTypes}
          onConnectEnd={onConnectEnd}
          onNodeClick={onNodeClick}
          onPaneClick={onPaneClick}
          fitView
          fitViewOptions={fitViewOptions}
          defaultEdgeOptions={defaultEdgeOptions}
          nodeOrigin={nodeOrigin}
          className={`w-full h-full`}
          panOnScroll
          selectionOnDrag
          style={{ opacity: framed ? 1 : 0 }}
          nodesDraggable={!locked}
          nodesConnectable={!locked}
          elementsSelectable={!locked}
        >
          <EdgeGradients nodes={nodes} edges={edges} />
          <Background
            size={1}
            variant={BackgroundVariant.Dots}
            className="dark:bg-foreground/5! bg-input/30!"
          />
          <FlowSidebar
            previewMode={previewMode}
            sidebarOpen={sidebarOpen}
            setSidebarOpen={setSidebarOpen}
            selectedCategory={selectedCategory}
            setSelectedCategory={setSelectedCategory}
            providerList={providerList}
            apiNodes={apiNodes}
            dragging={dragging}
            startDrag={startDrag}
            handleSearch={handleSearch}
            chatOpen={chatOpen}
            setChatOpen={handleSetChatOpen}
          />
          <FlowToolbar
            viewControls={
              <FlowViewControls
                locked={locked}
                setLocked={setLocked}
                previewMode={previewMode}
              />
            }
            previewMode={previewMode}
            hasUnsavedChanges={hasUnsavedChanges}
            deleteLastSaved={revertFlow}
            runFlow={runFlow}
            setOpen={setOpen}
            openApiSheet={
              initialFlow?.id ? () => setApiSheetOpen(true) : undefined
            }
            nodes={nodes}
            apiNodes={apiNodes}
            organizationId={organizationId}
          />
          <FlowDragPreview
            dragging={dragging}
            previewData={previewData}
            previewPos={previewPos}
          />
          {!previewMode && (
            <WelcomeTour
              sidebarOpen={sidebarOpen}
              setSidebarOpen={setSidebarOpen}
              addNodeById={addNodeById}
              selectGroupForNode={selectGroupForNode}
              connectNodes={connectNodes}
              focusCanvas={focusCanvas}
              openNodeSettings={openNodeSettings}
              getFlowState={getFlowState}
              setChatOpen={handleSetChatOpen}
              nextStep={nextStepInit}
            />
          )}
        </ReactFlow>

        <FlowSaveDialog
          open={open}
          onOpenChange={setOpen}
          flowData={flowData}
          updateFlowData={updateFlowData}
          saveFlow={saveFlow}
        />

        <FlowNavigationGuard
          open={showNavigationDialog}
          onOpenChange={() => setShowNavigationDialog(open)}
          onConfirm={handleNavigationConfirm}
          onCancel={handleNavigationCancel}
        />

        <AIChatSheet
          organizationId={organizationId}
          isOpen={chatOpen}
          onClose={() => setChatOpen(false)}
          onApplyWorkflow={handleApplyWorkflow}
        />

        <FlowApiSheet
          open={apiSheetOpen}
          onOpenChange={setApiSheetOpen}
          flowId={initialFlow?.id}
          nodes={nodes}
          apiNodes={apiNodes}
        />

        {selectedNode && (
          <NodeSettingsPanel
            node={selectedNode}
            schema={selectedNodeSchema}
            hasNaturalLanguage={selectedNodeHasNaturalLanguage}
            triggerVariables={triggerVariables}
            onApply={handleApplyNodeSettings}
            onCancel={() => setSelectedNodeId(null)}
          />
        )}
      </div>
    </FlowNodesProvider>
  )
}

FlowCanvas.displayName = 'FlowCanvas'

const FlowCanvasWithProvider = (props) => {
  return (
    <ReactFlowProvider>
      <FlowCanvas {...props} />
    </ReactFlowProvider>
  )
}
FlowCanvasWithProvider.displayName = 'FlowCanvasWithProvider'

export default FlowCanvasWithProvider
