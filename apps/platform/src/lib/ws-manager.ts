import tokenManager from '@/lib/token-manager'
import { config } from '@/lib/config'

export interface OrganizationEvent {
  id: string
  type: 'task' | 'node'
  organizationId: string
  executionContext: 'workflow'
  contextItemId: string
  status:
    | 'pending'
    | 'scheduled'
    | 'running'
    | 'retrying'
    | 'success'
    | 'failed'
    | 'cancelled'
  error?: string
  duration?: number
  // node-specific fields
  nodeId?: string
  nodeType?: string
  output?: Record<string, any>[]
}

type EventListener = (event: OrganizationEvent) => void

export type WsStatus = 'offline' | 'connecting' | 'online' | 'reconnecting'
type StatusListener = (status: WsStatus) => void

const BASE_URL = config.platformApiUrl

const MAX_RECONNECT_DELAY = 30000
const INITIAL_RECONNECT_DELAY = 1000

class WSManager {
  private ws: WebSocket | null = null
  private organizationId: string | null = null
  private listeners = new Set<EventListener>()
  private statusListeners = new Set<StatusListener>()
  private status: WsStatus = 'offline'
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null
  private reconnectDelay = INITIAL_RECONNECT_DELAY
  private shouldReconnect = false

  connect(organizationId: string) {
    if (
      this.organizationId === organizationId &&
      this.ws?.readyState === WebSocket.OPEN
    ) {
      return
    }

    this.disconnect()
    this.organizationId = organizationId
    this.shouldReconnect = true
    this.reconnectDelay = INITIAL_RECONNECT_DELAY
    this.setStatus('connecting')
    this.createConnection()
  }

  disconnect() {
    this.shouldReconnect = false
    this.organizationId = null

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }

    if (this.ws) {
      this.ws.onclose = null
      this.ws.close()
      this.ws = null
    }

    this.setStatus('offline')
  }

  subscribe(listener: EventListener): () => void {
    this.listeners.add(listener)
    return () => {
      this.listeners.delete(listener)
    }
  }

  getStatus(): WsStatus {
    return this.status
  }

  subscribeStatus(listener: StatusListener): () => void {
    this.statusListeners.add(listener)
    listener(this.status)
    return () => {
      this.statusListeners.delete(listener)
    }
  }

  private setStatus(status: WsStatus) {
    if (this.status === status) {
      return
    }
    this.status = status
    this.statusListeners.forEach((listener) => listener(status))
  }

  private createConnection() {
    const token = tokenManager.getAccessToken()
    if (!token || !this.organizationId) {
      return
    }

    const wsUrl = BASE_URL.replace(/^http/, 'ws')

    const url = `${wsUrl}/organizations/${this.organizationId}/ws?token=${token}`

    this.ws = new WebSocket(url)

    this.ws.onmessage = (event) => {
      try {
        const data: OrganizationEvent = JSON.parse(event.data)
        this.listeners.forEach((listener) => listener(data))
      } catch {
        // ignore malformed messages
      }
    }

    this.ws.onopen = () => {
      this.reconnectDelay = INITIAL_RECONNECT_DELAY
      this.setStatus('online')
    }

    this.ws.onclose = () => {
      this.ws = null
      if (this.shouldReconnect) {
        this.setStatus('reconnecting')
      }
      this.scheduleReconnect()
    }

    this.ws.onerror = () => {
      this.ws?.close()
    }
  }

  private scheduleReconnect() {
    if (!this.shouldReconnect || !this.organizationId) {
      return
    }

    this.reconnectTimer = setTimeout(() => {
      this.createConnection()
    }, this.reconnectDelay)

    this.reconnectDelay = Math.min(this.reconnectDelay * 2, MAX_RECONNECT_DELAY)
  }
}

const wsManager = new WSManager()
export default wsManager
