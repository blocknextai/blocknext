import { useCallback, useState } from 'react'
import { useAuthSessions, useAuthActions } from '@/features/auth'
import Sessions from '@/features/sessions/components/sessions'

const SESSIONS_LIMIT = 10

const PreferencesSessionsPage = () => {
  const [sessionsPage, setSessionsPage] = useState(1)
  const sessionsOffset = (sessionsPage - 1) * SESSIONS_LIMIT

  const {
    sessions,
    pagination: sessionsPagination,
    isLoading: sessionsLoading,
  } = useAuthSessions(sessionsOffset, SESSIONS_LIMIT)
  const { revokeSession, revokeAllSessions } = useAuthActions()

  const [revokingSessionId, setRevokingSessionId] = useState<string | null>(
    null,
  )
  const [revokingAll, setRevokingAll] = useState(false)

  const handleRevokeSession = useCallback(
    async (sessionId: string) => {
      setRevokingSessionId(sessionId)
      try {
        await revokeSession(sessionId)
      } finally {
        setRevokingSessionId(null)
      }
    },
    [revokeSession],
  )

  const handleRevokeAllSessions = useCallback(async () => {
    setRevokingAll(true)
    try {
      await revokeAllSessions()
    } finally {
      setRevokingAll(false)
    }
  }, [revokeAllSessions])

  return (
    <Sessions
      sessions={sessions}
      loading={sessionsLoading}
      revokingSessionId={revokingSessionId}
      revokingAll={revokingAll}
      onRevokeSession={handleRevokeSession}
      onRevokeAllSessions={handleRevokeAllSessions}
      pagination={sessionsPagination}
      page={sessionsPage}
      limit={SESSIONS_LIMIT}
      offset={sessionsOffset}
      onPageChange={setSessionsPage}
    />
  )
}

export default PreferencesSessionsPage
