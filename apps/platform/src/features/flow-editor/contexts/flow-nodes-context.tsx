import { createContext, useContext } from 'react'
import type { SetFlowNodes, SetFlowEdges } from '@/features/flow-editor/types'

interface FlowContextValue {
  setNodes: SetFlowNodes
  setEdges: SetFlowEdges
}

const FlowNodesContext = createContext<FlowContextValue | null>(null)

export const FlowNodesProvider = FlowNodesContext.Provider

export const useFlowSetNodes = () => {
  const ctx = useContext(FlowNodesContext)
  if (!ctx) {
    throw new Error('useFlowSetNodes must be used within a FlowNodesProvider')
  }
  return ctx.setNodes
}

export const useFlowSetEdges = () => {
  const ctx = useContext(FlowNodesContext)
  if (!ctx) {
    throw new Error('useFlowSetEdges must be used within a FlowNodesProvider')
  }
  return ctx.setEdges
}
