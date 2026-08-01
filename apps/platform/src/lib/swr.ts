import type { SWRConfiguration } from 'swr'
import type { ApiResponse } from '@/lib/api-client'

export const swrConfig: SWRConfiguration = {
  revalidateOnFocus: false,
  shouldRetryOnError: false,
  errorRetryCount: 0,
  dedupingInterval: 2000,
}

export async function unwrap<T>(promise: Promise<ApiResponse<T>>): Promise<T> {
  const res = await promise
  if (!res.isSuccess) {
    throw new Error(res.message || 'Request failed')
  }
  return res.data
}
