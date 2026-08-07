import platformApi from '@/lib/platform-api'
import tokenManager from '@/lib/token-manager'
import { config } from '@/lib/config'
import i18next from 'i18next'

const BASE_URL = config.platformApiUrl

function parseSSEEvents(raw) {
  const events = []
  // Split by double newline to separate SSE events
  const blocks = raw.split('\n\n')
  for (const block of blocks) {
    if (!block.trim()) {
      continue
    }

    const lines = block.split('\n')
    let eventType = 'message'
    const dataLines = []

    for (const line of lines) {
      if (line.startsWith('event: ')) {
        eventType = line.slice(7).trim()
      } else if (line.startsWith('data: ')) {
        dataLines.push(line.slice(6))
      } else if (line.startsWith('data:')) {
        // "data:" with no space means empty line
        dataLines.push(line.slice(5))
      }
    }

    // SSE spec: multiple data lines are joined with newlines
    const data = dataLines.join('\n')
    events.push({ eventType, data })
  }
  return events
}

export const generation = {
  getSessions: async (
    organizationId,
    { query = '', offset = 0, limit = 10 } = {},
  ) => {
    const params = new URLSearchParams()
    if (query) {
      params.set('query', query)
    }
    params.set('offset', offset)
    params.set('limit', limit)
    return await platformApi.get(
      `/organizations/${organizationId}/workflows/generation/sessions?${params.toString()}`,
    )
  },

  createSession: async (organizationId, title) => {
    return await platformApi.post(
      `/organizations/${organizationId}/workflows/generation/sessions`,
      {
        organizationId,
        title,
      },
    )
  },

  getSessionMessages: async (
    organizationId,
    sessionId,
    { offset = 0, limit = 100 } = {},
  ) => {
    const params = new URLSearchParams()
    params.set('offset', String(offset))
    params.set('limit', String(limit))
    return await platformApi.get(
      `/organizations/${organizationId}/workflows/generation/sessions/${sessionId}/messages?${params.toString()}`,
    )
  },

  updateSession: async (organizationId, sessionId, title) => {
    return await platformApi.patch(
      `/organizations/${organizationId}/workflows/generation/sessions/${sessionId}`,
      { title },
    )
  },

  deleteSession: async (organizationId, sessionId) => {
    return await platformApi.delete(
      `/organizations/${organizationId}/workflows/generation/sessions/${sessionId}`,
    )
  },

  sendMessageStream: async (
    organizationId,
    sessionId,
    message,
    { onChunk, onDone, onError },
  ) => {
    const token = tokenManager.getAccessToken()
    const headers = {
      'Content-Type': 'application/json',
    }
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }

    try {
      const response = await fetch(
        `${BASE_URL}/organizations/${organizationId}/workflows/generation/sessions/${sessionId}/messages`,
        {
          method: 'POST',
          headers,
          body: JSON.stringify({ message }),
        },
      )

      if (!response.ok) {
        const errorData = await response.json().catch(() => null)
        onError?.(errorData?.message ?? `HTTP ${response.status}`)
        return
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) {
          break
        }

        buffer += decoder.decode(value, { stream: true })

        // Only process complete SSE events (terminated by \n\n)
        // Keep incomplete data in buffer for next iteration
        const lastDoubleNewline = buffer.lastIndexOf('\n\n')
        if (lastDoubleNewline === -1) {
          continue
        }

        const complete = buffer.slice(0, lastDoubleNewline + 2)
        buffer = buffer.slice(lastDoubleNewline + 2)

        const events = parseSSEEvents(complete)
        for (const { eventType, data } of events) {
          switch (eventType) {
            case 'message':
              onChunk?.(data)
              break
            case 'done':
              onDone?.()
              break
            case 'error':
              onError?.(data)
              break
          }
        }
      }

      // Process any remaining buffer
      if (buffer.trim()) {
        const events = parseSSEEvents(buffer)
        for (const { eventType, data } of events) {
          switch (eventType) {
            case 'done':
              onDone?.()
              break
            case 'message':
              if (data) {
                onChunk?.(data)
              }
              break
            case 'error':
              onError?.(data)
              break
          }
        }
      }
    } catch (err) {
      onError?.(err.message || i18next.t('ui.text.connectionFailed'))
    }
  },
}

export default generation
