import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Loader2, Plus } from 'lucide-react'
import { AppPagination } from '@/components/shared/app-pagination'
import { SessionItemActions } from '@/features/generation/components/session-actions'

const SessionHistory = ({
  sessions,
  sessionsLoading,
  searchQuery,
  pagination,
  sessionId,
  initializing,
  editingSessionId,
  editingTitle,
  actionLoading,
  onPageChange,
  onSelectSession,
  onStartNewChat,
  onStartRename,
  onCancelRename,
  onConfirmRename,
  onRenameKeyDown,
  onEditingTitleChange,
  onDeleteRequest,
}) => {
  const { t } = useTranslation()

  return (
    <div className="flex-1 flex flex-col overflow-hidden">
      <div className="px-4 py-3 border-b flex justify-end">
        <Button
          size="sm"
          className="h-9"
          onClick={onStartNewChat}
          disabled={initializing}
        >
          {initializing ? (
            <Loader2 className="w-4 h-4 animate-spin" />
          ) : (
            <Plus className="w-4 h-4" />
          )}
        </Button>
      </div>

      <div className="flex-1 overflow-y-auto">
        {sessionsLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="w-5 h-5 animate-spin text-muted-foreground" />
          </div>
        ) : sessions.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 text-center px-4">
            <p className="text-sm text-muted-foreground">
              {searchQuery
                ? t('ui.text.generation.noChatsFound')
                : t('ui.text.generation.noChatHistoryYet')}
            </p>
          </div>
        ) : (
          <div className="divide-y">
            {sessions.map((session) => (
              <div
                key={session.id}
                className={`group relative flex items-center px-4 py-3 hover:bg-accent/50 transition-colors ${session.id === sessionId ? 'bg-accent' : ''}`}
              >
                <SessionItemActions
                  session={session}
                  isEditing={editingSessionId === session.id}
                  editingTitle={editingTitle}
                  actionLoading={actionLoading}
                  onEditingTitleChange={onEditingTitleChange}
                  onStartRename={onStartRename}
                  onCancelRename={onCancelRename}
                  onConfirmRename={onConfirmRename}
                  onRenameKeyDown={onRenameKeyDown}
                  onSelectSession={onSelectSession}
                  onDeleteRequest={onDeleteRequest}
                />
              </div>
            ))}
          </div>
        )}
      </div>

      {pagination.total > 0 && (
        <div className="border-t">
          <AppPagination
            pagination={pagination}
            onPageChange={onPageChange}
            showSummary={false}
          />
        </div>
      )}
    </div>
  )
}

SessionHistory.displayName = 'SessionHistory'

export default SessionHistory
