import config from '@/config'

const getCookie = (name: string): string | null => {
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) return parts.pop()?.split(';').shift() || null
  return null
}

const getHeaders = () => {
  const token = getCookie('token')
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  }
}

const request = async <T>(method: string, url: string, data?: unknown): Promise<{ data: T }> => {
  const apiUrl = config.webAPI || 'http://localhost:8080'
  const fullUrl = `${apiUrl}/api${url}`

  const options: RequestInit = {
    method,
    headers: getHeaders(),
    body: data ? JSON.stringify(data) : undefined,
  }

  const response = await fetch(fullUrl, options)

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.message || `HTTP error! Status: ${response.status}`)
  }

  const responseData = await response.json()
  return { data: responseData }
}

export default {
  get: <T>(url: string, config?: { params?: Record<string, unknown> }) => {
    let query = ''
    if (config?.params) {
      // Convert params to string, handling numbers etc.
      const params = new URLSearchParams()
      Object.entries(config.params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          params.append(key, String(value))
        }
      })
      query = '?' + params.toString()
    }
    return request<T>('GET', url + query)
  },
  post: <T>(url: string, data?: unknown) => request<T>('POST', url, data),
  put: <T>(url: string, data?: unknown) => request<T>('PUT', url, data),
  delete: <T>(url: string) => request<T>('DELETE', url),
}
