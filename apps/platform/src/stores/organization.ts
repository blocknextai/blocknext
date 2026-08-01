import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface OrganizationState {
  organizationId: string | null
  organizations: any[]
  setOrganizations: (value: any[]) => void
  setOrganizationId: (value: string) => void
  reset: () => void
}

export const useOrganizationStore = create<OrganizationState>()(
  persist(
    (set) => ({
      organizationId: null,
      organizations: [],

      setOrganizations: (value) => {
        set({ organizations: value })
      },
      setOrganizationId: (value) => {
        set({ organizationId: value })
      },
      reset: () => {
        set({ organizationId: null, organizations: [] })
      },
    }),
    {
      name: 'organization-store',
      partialize: (state) => ({ organizationId: state.organizationId }),
    },
  ),
)
