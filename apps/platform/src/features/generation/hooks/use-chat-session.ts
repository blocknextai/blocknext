import { useState, useRef, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import generation from '@/features/generation/services/generation'

export interface ChatMessage {
  id: number | string
  role: 'user' | 'assistant'
  content: string
}

export interface Session {
  id: string
  title: string
  updatedAt: string
}

export interface PaginationState {
  total: number
  page: number
  limit: number
  offset: number
  hasNext: boolean
  hasPrev: boolean
}

export type ViewMode = 'chat' | 'history'

export function useChatSession(organizationId: string) {
  const { t } = useTranslation()
  const [sessionId, setSessionId] = useState<string | null>(null)
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [input, setInput] = useState('')
  const [streaming, setStreaming] = useState(false)
  const [initializing, setInitializing] = useState(false)
  const abortRef = useRef(false)

  // Session history state
  const [view, setView] = useState<ViewMode>('chat')
  const [sessions, setSessions] = useState<Session[]>([])
  const [sessionsLoading, setSessionsLoading] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const [pagination, setPagination] = useState<PaginationState>({
    total: 0,
    page: 1,
    limit: 10,
    offset: 0,
    hasNext: false,
    hasPrev: false,
  })
  const searchTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  // Rename/Delete state
  const [editingSessionId, setEditingSessionId] = useState<string | null>(null)
  const [editingTitle, setEditingTitle] = useState('')
  const [deleteConfirmSession, setDeleteConfirmSession] =
    useState<Session | null>(null)
  const [actionLoading, setActionLoading] = useState(false)

  // Session history
  const fetchSessions = async (query = searchQuery, offset = 0) => {
    setSessionsLoading(true)
    try {
      const response = await generation.getSessions(organizationId, {
        query,
        offset,
        limit: pagination.limit,
      })
      if (response?.isSuccess && response.data && response.meta) {
        setSessions(response.data)
        setPagination(response.meta.pagination)
      }
    } catch {
      // api service handles errors
    } finally {
      setSessionsLoading(false)
    }
  }

  const openHistory = () => {
    setView('history')
    fetchSessions('', 0)
    setSearchQuery('')
  }

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setSearchQuery(value)
    if (searchTimeoutRef.current) {
      clearTimeout(searchTimeoutRef.current)
    }
    searchTimeoutRef.current = setTimeout(() => {
      fetchSessions(value, 0)
    }, 300)
  }

  const handlePageChange = (newPage: number) => {
    const newOffset = (newPage - 1) * pagination.limit
    fetchSessions(searchQuery, newOffset)
  }

  const selectSession = async (session: Session) => {
    if (editingSessionId) {
      return
    } // Don't select while editing
    setSessionId(session.id)
    setMessages([])
    setView('chat')
    try {
      const msgs = await generation.getSessionMessages(
        organizationId,
        session.id,
      )
      setMessages(msgs.data)
    } catch {
      // api service handles errors
    }
  }

  const startRename = (session: Session, e?: React.MouseEvent) => {
    e?.stopPropagation()
    setEditingSessionId(session.id)
    setEditingTitle(session.title)
  }

  const cancelRename = () => {
    setEditingSessionId(null)
    setEditingTitle('')
  }

  const confirmRename = async () => {
    if (!editingTitle.trim() || !editingSessionId) {
      return
    }
    setActionLoading(true)
    try {
      await generation.updateSession(
        organizationId,
        editingSessionId,
        editingTitle.trim(),
      )
      setSessions((prev) =>
        prev.map((s) =>
          s.id === editingSessionId ? { ...s, title: editingTitle.trim() } : s,
        ),
      )
      setEditingSessionId(null)
      setEditingTitle('')
    } catch {
      // api service handles errors
    } finally {
      setActionLoading(false)
    }
  }

  const handleRenameKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      confirmRename()
    } else if (e.key === 'Escape') {
      cancelRename()
    }
  }

  const confirmDelete = async () => {
    if (!deleteConfirmSession) {
      return
    }
    setActionLoading(true)
    try {
      await generation.deleteSession(organizationId, deleteConfirmSession.id)
      setSessions((prev) =>
        prev.filter((s) => s.id !== deleteConfirmSession.id),
      )
      if (sessionId === deleteConfirmSession.id) {
        setSessionId(null)
        setMessages([])
      }
      setDeleteConfirmSession(null)
    } catch {
      // api service handles errors
    } finally {
      setActionLoading(false)
    }
  }

  const startNewChat = () => {
    setSessionId(null)
    setMessages([])
    setView('chat')
  }

  const handleSend = async () => {
    const text = input.trim()
    if (!text || streaming) {
      return
    }

    const userMessage: ChatMessage = {
      role: 'user',
      content: text,
      id: Date.now(),
    }
    setMessages((prev) => [...prev, userMessage])
    setInput('')
    setStreaming(true)
    abortRef.current = false

    let currentSessionId = sessionId

    // Create session on first message if not exists
    if (!currentSessionId) {
      setInitializing(true)
      try {
        const session = await generation.createSession(
          organizationId,
          t('ui.text.generation.newChat'),
        )
        const newId = session?.data?.id ?? session?.id
        if (newId) {
          currentSessionId = newId
          setSessionId(newId)
        }
      } catch {
        setStreaming(false)
        setInitializing(false)
        return
      } finally {
        setInitializing(false)
      }
    }

    if (!currentSessionId) {
      setStreaming(false)
      return
    }

    const assistantId = Date.now() + 1
    setMessages((prev) => [
      ...prev,
      { role: 'assistant', content: '', id: assistantId },
    ])

    await generation.sendMessageStream(organizationId, currentSessionId, text, {
      onChunk: (chunk: string) => {
        if (abortRef.current) {
          return
        }
        setMessages((prev) =>
          prev.map((m) =>
            m.id === assistantId ? { ...m, content: m.content + chunk } : m,
          ),
        )
      },
      onDone: () => {
        setStreaming(false)
      },
      onError: (err: string) => {
        setStreaming(false)
        if (!abortRef.current) {
          setMessages((prev) =>
            prev.map((m) =>
              m.id === assistantId
                ? {
                    ...m,
                    content: m.content || `Error: ${err}`,
                  }
                : m,
            ),
          )
        }
      },
    })
  }

  const scrollToBottom = useCallback(
    (ref: React.RefObject<HTMLDivElement | null>) => {
      ref.current?.scrollIntoView({ behavior: 'smooth' })
    },
    [],
  )

  return {
    // Session state
    sessionId,
    setSessionId,
    messages,
    setMessages,
    input,
    setInput,
    streaming,
    initializing,

    // View state
    view,
    setView,

    // Session history
    sessions,
    sessionsLoading,
    searchQuery,
    pagination,
    fetchSessions,
    openHistory,
    handleSearchChange,
    handlePageChange,
    selectSession,

    // Rename/Delete
    editingSessionId,
    editingTitle,
    setEditingTitle,
    deleteConfirmSession,
    setDeleteConfirmSession,
    actionLoading,
    startRename,
    cancelRename,
    confirmRename,
    handleRenameKeyDown,
    confirmDelete,

    // Chat actions
    startNewChat,
    handleSend,
    scrollToBottom,
  }
}
