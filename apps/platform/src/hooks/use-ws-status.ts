import { useSyncExternalStore } from 'react'
import wsManager, { type WsStatus } from '@/lib/ws-manager'

export function useWsStatus(): WsStatus {
  return useSyncExternalStore(
    (listener) => wsManager.subscribeStatus(listener),
    () => wsManager.getStatus(),
    () => 'offline',
  )
}
