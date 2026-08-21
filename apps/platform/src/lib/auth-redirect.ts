export const LOGIN_PATH = '/auth/login'

export const sanitizeReturnUrl = (
  value: string | null | undefined,
): string | null => {
  if (!value || !value.startsWith('/')) {
    return null
  }
  if (value.startsWith('//') || value.startsWith('/\\')) {
    return null
  }
  return value
}

export const buildLoginUrl = (currentPath: string): string => {
  const returnTo = sanitizeReturnUrl(currentPath)
  if (!returnTo || returnTo === '/' || returnTo.startsWith(LOGIN_PATH)) {
    return LOGIN_PATH
  }
  return `${LOGIN_PATH}?returnUrl=${encodeURIComponent(returnTo)}`
}

export const getReturnUrl = (fallback = '/'): string => {
  const raw = new URLSearchParams(window.location.search).get('returnUrl')
  return sanitizeReturnUrl(raw) ?? fallback
}
