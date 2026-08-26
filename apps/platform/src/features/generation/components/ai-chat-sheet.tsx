import { useCallback, useEffect, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { Sparkles, ChevronLeft, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import ChatMessage from '@/features/generation/components/chat-message'
import SessionHistory from '@/features/generation/components/session-history'
import ChatInput from '@/features/generation/components/chat-input'
import type { ChatInputHandle } from '@/features/generation/components/chat-input'
import { DeleteSessionDialog } from '@/features/generation/components/session-actions'
import { useChatSession } from '@/features/generation/hooks/use-chat-session'

const AIChatSheet = ({ organizationId, isOpen, onClose, onApplyWorkflow }) => {
  const { t } = useTranslation()
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const scrollContainerRef = useRef<HTMLDivElement>(null)
  const isUserScrolledUp = useRef(false)
  const chatInputRef = useRef<ChatInputHandle>(null)

  const {
    sessionId,
    messages,
    input,
    setInput,
    streaming,
    initializing,
    view,
    setView,
    sessions,
    sessionsLoading,
    searchQuery,
    pagination,
    openHistory,
    handleSearchChange,
    handlePageChange,
    selectSession: baseSelectSession,
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
    startNewChat: baseStartNewChat,
    handleSend: baseHandleSend,
    scrollToBottom,
  } = useChatSession(organizationId)

  useEffect(() => {
    if (isOpen) {
      setView('chat')
      setTimeout(() => chatInputRef.current?.focus(), 150)
    }
  }, [isOpen])

  const handleScroll = useCallback(() => {
    const el = scrollContainerRef.current
    if (!el) {
      return
    }
    const threshold = 50
    isUserScrolledUp.current =
      el.scrollHeight - el.scrollTop - el.clientHeight > threshold
  }, [])

  useEffect(() => {
    if (!isUserScrolledUp.current) {
      scrollToBottom(messagesEndRef)
    }
  }, [messages, scrollToBottom])

  const selectSession = async (session) => {
    await baseSelectSession(session)
    setTimeout(() => chatInputRef.current?.focus(), 150)
  }

  const startNewChat = () => {
    baseStartNewChat()
    setTimeout(() => chatInputRef.current?.focus(), 150)
  }

  const handleSend = async () => {
    chatInputRef.current?.resetHeight()
    isUserScrolledUp.current = false
    await baseHandleSend()
    setTimeout(() => chatInputRef.current?.focus(), 150)
  }

  const handleApply = (data) => {
    onApplyWorkflow(data)
    onClose()
  }

  if (!isOpen) {
    return null
  }

  return (
    <>
      <div
        data-tour="ai-chat"
        className="absolute right-4 top-4 bottom-4 w-[420px] max-w-[90vw] z-[4] bg-background border border-border rounded-xl shadow-lg flex flex-col overflow-hidden animate-in slide-in-from-right duration-200"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between gap-3 px-4 py-3 border-b border-border shrink-0">
          <div className="flex items-center gap-2">
            {view === 'history' ? (
              <>
                {sessionId && (
                  <button
                    onClick={() => setView('chat')}
                    className="hover:bg-accent rounded-md p-1 -ml-1 transition-colors"
                  >
                    <ChevronLeft className="w-4 h-4" />
                  </button>
                )}
                <span className="text-sm font-semibold">
                  {t('ui.text.generation.chatHistory')}
                </span>
              </>
            ) : (
              <>
                <Sparkles className="w-4 h-4 text-primary" />
                <span className="text-sm font-semibold">
                  {t('ui.text.generation.aiFlowBuilder')}
                </span>
              </>
            )}
          </div>
          <Button variant="ghost" size="icon-sm" onClick={onClose}>
            <X className="size-4" />
          </Button>
        </div>

        {view === 'history' ? (
          <SessionHistory
            sessions={sessions}
            sessionsLoading={sessionsLoading}
            searchQuery={searchQuery}
            pagination={pagination}
            sessionId={sessionId}
            initializing={initializing}
            editingSessionId={editingSessionId}
            editingTitle={editingTitle}
            actionLoading={actionLoading}
            onSearchChange={handleSearchChange}
            onPageChange={handlePageChange}
            onSelectSession={selectSession}
            onStartNewChat={startNewChat}
            onStartRename={startRename}
            onCancelRename={cancelRename}
            onConfirmRename={confirmRename}
            onRenameKeyDown={handleRenameKeyDown}
            onEditingTitleChange={setEditingTitle}
            onDeleteRequest={setDeleteConfirmSession}
          />
        ) : (
          <>
            <div
              ref={scrollContainerRef}
              onScroll={handleScroll}
              className="flex-1 overflow-y-auto p-4"
            >
              {messages.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-full text-center px-4">
                  <div className="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center mb-3">
                    <Sparkles className="w-6 h-6 text-primary" />
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {t('ui.text.generation.describeWorkflow')}
                  </p>
                </div>
              ) : (
                <div className="space-y-4">
                  {messages.map((msg) => (
                    <ChatMessage
                      key={msg.id}
                      message={msg}
                      onApplyWorkflow={handleApply}
                    />
                  ))}

                  {streaming &&
                    messages[messages.length - 1]?.content === '' && (
                      <div className="flex gap-3">
                        <div className="bg-muted rounded-2xl rounded-bl-md px-4 py-2.5">
                          <div className="flex gap-1">
                            <span className="w-2 h-2 bg-foreground/30 rounded-full animate-bounce [animation-delay:0ms]" />
                            <span className="w-2 h-2 bg-foreground/30 rounded-full animate-bounce [animation-delay:150ms]" />
                            <span className="w-2 h-2 bg-foreground/30 rounded-full animate-bounce [animation-delay:300ms]" />
                          </div>
                        </div>
                      </div>
                    )}

                  <div ref={messagesEndRef} />
                </div>
              )}
            </div>

            <ChatInput
              ref={chatInputRef}
              input={input}
              streaming={streaming}
              initializing={initializing}
              onInputChange={setInput}
              onSend={handleSend}
              onStartNewChat={startNewChat}
              onOpenHistory={openHistory}
            />
          </>
        )}
      </div>

      <DeleteSessionDialog
        session={deleteConfirmSession}
        actionLoading={actionLoading}
        onOpenChange={() => setDeleteConfirmSession(null)}
        onConfirmDelete={confirmDelete}
      />
    </>
  )
}

AIChatSheet.displayName = 'AIChatSheet'

export default AIChatSheet
