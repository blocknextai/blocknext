import cookie from '@/lib/cookie'

const ACCESS_TOKEN_KEY = 'accessToken'
const REFRESH_TOKEN_KEY = 'refreshToken'

const getAccessToken = () => {
  return cookie.get(ACCESS_TOKEN_KEY)
}

const getRefreshToken = () => {
  return cookie.get(REFRESH_TOKEN_KEY)
}

const setTokens = (accessToken: string, refreshToken?: string) => {
  cookie.set(ACCESS_TOKEN_KEY, accessToken, 1)
  if (refreshToken) {
    cookie.set(REFRESH_TOKEN_KEY, refreshToken, 30)
  }
}

const clearTokens = () => {
  cookie.remove(ACCESS_TOKEN_KEY)
  cookie.remove(REFRESH_TOKEN_KEY)
}

const decodeAccessToken = <T = Record<string, unknown>>(): T | null => {
  const token = getAccessToken()
  if (!token) {
    return null
  }
  const parts = token.split('.')
  if (parts.length !== 3) {
    return null
  }
  try {
    const base64 = parts[1].replace(/-/g, '+').replace(/_/g, '/')
    const padded = base64 + '='.repeat((4 - (base64.length % 4)) % 4)
    return JSON.parse(atob(padded)) as T
  } catch {
    return null
  }
}

const getUserId = (): string | null => {
  const claims = decodeAccessToken<{ sub?: string; userId?: string }>()
  return claims?.userId ?? claims?.sub ?? null
}

const getSessionId = (): string | null => {
  const claims = decodeAccessToken<{ sessionId?: string }>()
  return claims?.sessionId ?? null
}

export default {
  getAccessToken,
  getRefreshToken,
  setTokens,
  clearTokens,
  decodeAccessToken,
  getUserId,
  getSessionId,
}
