import { useCallback } from 'react'
import { useNavigate } from 'react-router'
import { useSWRConfig } from 'swr'
import { useOrganizationStore } from '@/stores/organization'

export function useAuthRedirect() {
  const navigate = useNavigate()
  const { mutate } = useSWRConfig()

  return useCallback(
    async (target: string) => {
      useOrganizationStore.getState().reset()
      await mutate(() => true, undefined, { revalidate: false })
      navigate(target, { replace: true })
    },
    [mutate, navigate],
  )
}
