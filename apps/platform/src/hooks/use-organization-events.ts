import { useEffect } from 'react'
import wsManager, { type OrganizationEvent } from '@/lib/ws-manager'

type EventFilter = {
  type?: 'task' | 'node'
  executionId?: string
}

export function useOrganizationEvents(
  callback: (event: OrganizationEvent) => void,
  filter?: EventFilter,
) {
  useEffect(() => {
    const unsubscribe = wsManager.subscribe((event) => {
      if (filter?.type && event.type !== filter.type) return
      if (filter?.executionId && event.id !== filter.executionId) return
      callback(event)
    })

    return unsubscribe
  }, [callback, filter?.type, filter?.executionId])
}
