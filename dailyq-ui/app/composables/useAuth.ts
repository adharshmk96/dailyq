import type { FetchError } from 'ofetch'

export const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export interface AuthUser {
  id: string
  email: string
  name: string
  created_at: string
}

interface AuthResponse {
  token: string
  expires_at: string
  session_id: string
  user: AuthUser
}

interface ApiError {
  code: string
  message: string
}

const TOKEN_COOKIE = 'dailyq_token'

/**
 * Turns any failure from $fetch into a human-readable message, preferring the
 * API's own `{ error: { code, message } }` envelope.
 */
export function apiErrorMessage(err: unknown, fallback = 'Something went wrong. Please try again.'): string {
  const data = (err as FetchError<{ error?: ApiError }>)?.data
  return data?.error?.message || fallback
}

export function useAuthToken() {
  return useCookie<string | null>(TOKEN_COOKIE, {
    default: () => null,
    sameSite: 'lax',
    secure: import.meta.client && window.location.protocol === 'https:',
    path: '/',
    maxAge: 60 * 60 * 24 * 30
  })
}

/**
 * Thin $fetch wrapper that targets the API base and attaches the bearer token.
 */
export function useApiFetch() {
  const config = useRuntimeConfig()
  const token = useAuthToken()

  return <T>(path: string, options: Parameters<typeof $fetch>[1] = {}) => {
    const headers = new Headers(options?.headers as HeadersInit)
    if (token.value) {
      headers.set('Authorization', `Bearer ${token.value}`)
    }

    return $fetch<T>(path, {
      ...options,
      baseURL: config.public.apiBase,
      headers
    })
  }
}

export function useAuth() {
  const token = useAuthToken()
  const user = useState<AuthUser | null>('auth.user', () => null)
  const pending = useState<boolean>('auth.pending', () => false)
  const api = useApiFetch()

  const isAuthenticated = computed(() => Boolean(token.value))

  function setSession(res: AuthResponse) {
    token.value = res.token
    user.value = res.user
  }

  function clearSession() {
    token.value = null
    user.value = null
  }

  async function withPending<T>(fn: () => Promise<T>): Promise<T> {
    pending.value = true
    try {
      return await fn()
    } finally {
      pending.value = false
    }
  }

  async function register(payload: { name: string, email: string, password: string }) {
    return withPending(async () => {
      const res = await api<AuthResponse>('/auth/register', { method: 'POST', body: payload })
      setSession(res)
      return res
    })
  }

  async function login(payload: { email: string, password: string }) {
    return withPending(async () => {
      const res = await api<AuthResponse>('/auth/login', { method: 'POST', body: payload })
      setSession(res)
      return res
    })
  }

  async function logout() {
    return withPending(async () => {
      try {
        if (token.value) {
          await api('/auth/logout', { method: 'POST' })
        }
      } finally {
        // A failed revoke must not strand the user in a signed-in shell.
        clearSession()
      }
    })
  }

  async function forgotPassword(email: string) {
    return withPending(() => api('/auth/forgot-password', { method: 'POST', body: { email } }))
  }

  async function resetPassword(resetToken: string, password: string) {
    return withPending(() => api('/auth/reset-password', {
      method: 'POST',
      params: { token: resetToken },
      body: { password }
    }))
  }

  async function changePassword(payload: { current_password: string, new_password: string }) {
    return withPending(async () => {
      await api('/auth/change-password', { method: 'POST', body: payload })
      // The API revokes every other session but keeps this one alive.
    })
  }

  async function updateProfile(payload: { name: string, email: string }) {
    return withPending(async () => {
      const updated = await api<AuthUser>('/auth/me', { method: 'PATCH', body: payload })
      user.value = updated
      return updated
    })
  }

  /** Loads the current user, clearing the session if the token is no longer valid. */
  async function fetchMe(): Promise<AuthUser | null> {
    if (!token.value) {
      user.value = null
      return null
    }

    try {
      user.value = await api<AuthUser>('/auth/me')
      return user.value
    } catch (err) {
      if ((err as FetchError)?.statusCode === 401) {
        clearSession()
        return null
      }
      throw err
    }
  }

  return {
    token,
    user,
    pending,
    isAuthenticated,
    register,
    login,
    logout,
    forgotPassword,
    resetPassword,
    changePassword,
    updateProfile,
    fetchMe,
    clearSession
  }
}
