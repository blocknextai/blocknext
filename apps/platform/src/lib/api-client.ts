import { toast } from 'sonner'
import tokenManager from '@/lib/token-manager'
import { buildLoginUrl, LOGIN_PATH } from '@/lib/auth-redirect'
import { config } from '@/lib/config'

export interface ApiPaginationMetaResponse {
  total: number
  page: number
  totalPages: number
  limit: number
  offset: number
  hasNext: boolean
  hasPrev: boolean
}

export interface ApiMetaResponse {
  pagination?: ApiPaginationMetaResponse
}

export interface ApiResponse<T = any> {
  isSuccess: boolean
  data: T
  message: string
  meta?: ApiMetaResponse
}

let isRefreshing = false
let refreshPromise: Promise<boolean> | null = null

const getHeaders = (): Headers => {
  const headers = new Headers()
  headers.set('Content-Type', 'application/json')

  const accessToken = tokenManager.getAccessToken()
  if (accessToken) {
    headers.set('Authorization', `Bearer ${accessToken}`)
  }

  return headers
}

const attemptTokenRefresh = async (): Promise<boolean> => {
  const refreshToken = tokenManager.getRefreshToken()
  if (!refreshToken) {
    return false
  }

  if (isRefreshing && refreshPromise) {
    return refreshPromise
  }

  isRefreshing = true
  refreshPromise = (async () => {
    try {
      const response = await fetch(
        `${config.platformApiUrl}/auth/token/refresh`,
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ refreshToken }),
        },
      )

      if (!response.ok) {
        return false
      }

      const result: ApiResponse<{ accessToken: string; refreshToken: string }> =
        await response.json()

      if (result.isSuccess && result.data) {
        tokenManager.setTokens(
          result.data.accessToken,
          result.data.refreshToken,
        )
        return true
      }

      return false
    } catch {
      return false
    } finally {
      isRefreshing = false
      refreshPromise = null
    }
  })()

  return refreshPromise
}

const redirectToLogin = () => {
  tokenManager.clearTokens()
  const { pathname, search } = window.location

  if (pathname.startsWith(LOGIN_PATH)) {
    return
  }

  window.location.href = buildLoginUrl(pathname + search)
}

export const createApiClient = (baseURL: string) => {
  const request = async <T = any>(
    url: string,
    options: RequestInit = {},
  ): Promise<ApiResponse<T>> => {
    const response = await fetch(`${baseURL}${url}`, {
      headers: getHeaders(),
      ...options,
    })

    if (response.status === 401 && !url.startsWith('/auth')) {
      const refreshed = await attemptTokenRefresh()

      if (refreshed) {
        const retryResponse = await fetch(`${baseURL}${url}`, {
          headers: getHeaders(),
          ...options,
        })

        const retryResult: ApiResponse<T> = await retryResponse.json()

        if (retryResult?.message) {
          if (retryResult?.isSuccess) {
            toast.success(retryResult.message)
          } else {
            toast.error(retryResult.message)
          }
        }

        return retryResult
      }

      redirectToLogin()
      return { isSuccess: false, data: undefined as T, message: '' }
    }

    const result: ApiResponse<T> = await response.json()

    if (result?.message) {
      if (result?.isSuccess) {
        toast.success(result.message)
      } else {
        toast.error(result.message)
      }
    }

    return result
  }

  const get = async <T = any>(url: string): Promise<ApiResponse<T>> => {
    return await request<T>(url, {
      method: 'GET',
    })
  }

  const post = async <T = any>(
    url: string,
    body = {},
  ): Promise<ApiResponse<T>> => {
    return await request<T>(url, {
      method: 'POST',
      body: JSON.stringify(body),
    })
  }

  const put = async <T = any>(
    url: string,
    body = {},
  ): Promise<ApiResponse<T>> => {
    return await request<T>(url, {
      method: 'PUT',
      body: JSON.stringify(body),
    })
  }

  const patch = async <T = any>(
    url: string,
    body = {},
  ): Promise<ApiResponse<T>> => {
    return await request<T>(url, {
      method: 'PATCH',
      body: JSON.stringify(body),
    })
  }

  const del = async <T = any>(
    url: string,
    body = {},
  ): Promise<ApiResponse<T>> => {
    return await request<T>(url, {
      method: 'DELETE',
      body: JSON.stringify(body),
    })
  }

  return {
    get,
    post,
    put,
    patch,
    delete: del,
  }
}
