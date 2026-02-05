import { DEFAULT_API_URL, type ApiResponse } from 'shared-configs'

export async function fetchApi<T>(url: string): Promise<ApiResponse<T>> {
  const headers: HeadersInit = {}
  const apiKey = process.env.NEXT_PUBLIC_API_KEY
  if (apiKey) {
    headers['X-API-Key'] = apiKey
  }

  const response = await fetch(url, { headers })
  let payload: ApiResponse<T> | undefined
  try {
    payload = (await response.json()) as ApiResponse<T>
  } catch {
    payload = undefined
  }

  if (!response.ok) {
    return {
      error: {
        code: payload?.error?.code || 'HTTP_ERROR',
        message: payload?.error?.message || 'Request failed',
      },
    }
  }

  return payload || {}
}

export function getApiUrl(): string {
  return process.env.NEXT_PUBLIC_API_URL || DEFAULT_API_URL
}
