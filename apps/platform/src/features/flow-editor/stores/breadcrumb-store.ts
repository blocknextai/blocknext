import { create } from 'zustand'

interface BreadcrumbState {
  flowName?: string
  setFlowName: (name: string) => void
}

export const useBreadcrumbStore = create<BreadcrumbState>((set) => ({
  flowName: undefined,
  setFlowName: (name) => set({ flowName: name }),
}))
