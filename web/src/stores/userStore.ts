import { create } from 'zustand'
import { setToken, removeToken, type User } from '../api/client'
import { isEmbedMode, initEmbedHandoff, exchangeEmbedToken, type GamesEmbedSessionPayload } from '../embedHandoff'
import { getProfile } from '../api/auth'

// ─── JWT helpers ──────────────────────────────────────────────────────────────

function parseJWT(token: string): Record<string, unknown> | null {
  try {
    const base64 = token.split('.')[1]
    const json = atob(base64.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(json)
  } catch {
    return null
  }
}

function isTokenExpired(token: string): boolean {
  const claims = parseJWT(token)
  if (!claims || typeof claims.exp !== 'number') return true
  // exp is seconds since epoch; give 10s buffer
  return Date.now() / 1000 > claims.exp - 10
}

function getAuthHeaders(token: string | null): Record<string, string> {
  if (!token) return {}
  return { Authorization: `Bearer ${token}` }
}

function getInitials(name: string): string {
  if (!name) return '?'
  const parts = name.trim().split(/\s+/)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
  }
  return name.slice(0, 2).toUpperCase()
}

function decodeJWTUser(token: string): {
  sid: string
  email: string
  name: string
  department: string
} | null {
  const claims = parseJWT(token)
  if (!claims) return null
  return {
    sid: String(claims.sid || ''),
    email: String(claims.email || ''),
    name: String(claims.name || ''),
    department: String(claims.department || ''),
  }
}

// ─── Store interface ──────────────────────────────────────────────────────────

interface UserState {
  sid: string
  email: string
  name: string
  department: string
  isAuthenticated: boolean
  isEmbed: boolean
  currentUserId: string
  theme: 'light' | 'dark' | null
  user: User | null        // full profile from API
  token: string | null
  loading: boolean
  error: string | null
}

interface UserActions {
  initialize: () => Promise<void>
  loginFromEmbed: (session: GamesEmbedSessionPayload) => Promise<void>
  loginFromToken: (token: string) => void
  logout: () => void
  fetchProfile: () => Promise<void>
  getInitials: () => string
  setTheme: (theme: 'light' | 'dark' | null) => void
  setError: (error: string | null) => void
  // Backward compatibility
  getCurrentUser: () => UserState | null
}

type UserStore = UserState & UserActions
const TOKEN_KEY = 'gz_token'

function readStoredToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

// ─── Zustand store ────────────────────────────────────────────────────────────

export const useUserStore = create<UserStore>((set, get) => ({
  // State
  sid: '',
  email: '',
  name: '',
  department: '',
  isAuthenticated: false,
  isEmbed: isEmbedMode(),
  currentUserId: '',
  theme: null,
  user: null,
  token: null,
  loading: true,
  error: null,

  // Actions

  initialize: async () => {
    const isEmbed = isEmbedMode()

    // 1. Check for existing JWT in localStorage
    const stored = readStoredToken()
    if (stored) {
      const jwtUser = decodeJWTUser(stored)
      if (jwtUser && jwtUser.sid && !isTokenExpired(stored)) {
        set({
          sid: jwtUser.sid,
          email: jwtUser.email,
          name: jwtUser.name,
          department: jwtUser.department,
          isAuthenticated: true,
          isEmbed,
          token: stored,
          loading: false,
        })
        get().fetchProfile()
        return
      }
      // Token is missing, invalid, or expired — clear it
      removeToken()
    }

    // 2. Embed flow: handshake with Portal Shell
    if (isEmbed) {
      try {
        const session = await initEmbedHandoff()
        await get().loginFromEmbed(session)
        return
      } catch (err) {
        set({
          loading: false,
          error: err instanceof Error ? err.message : 'Embed auth failed',
        })
        return
      }
    }

    // 3. Standalone: no auth available
    set({ loading: false })
  },

  loginFromEmbed: async (session: GamesEmbedSessionPayload) => {
    try {
      const result = await exchangeEmbedToken(session)
      const jwtUser = decodeJWTUser(result.token)
      const theme = session.theme ?? null
      set({
        sid: result.sid,
        email: jwtUser?.email || session.email || '',
        name: jwtUser?.name || session.name || '',
        department: jwtUser?.department || session.department || '',
        isAuthenticated: true,
        isEmbed: true,
        theme,
        token: result.token,
        loading: false,
        error: null,
      })
      setToken(result.token)
      get().fetchProfile()
    } catch (err) {
      set({
        loading: false,
        error: err instanceof Error ? err.message : 'Embed token exchange failed',
      })
    }
  },

  loginFromToken: (token: string) => {
    const jwtUser = decodeJWTUser(token)
    if (!jwtUser || !jwtUser.sid) {
      set({ error: 'Invalid token: cannot decode claims' })
      return
    }
    if (isTokenExpired(token)) {
      set({ error: 'Token is expired' })
      return
    }
    setToken(token)
    set({
      sid: jwtUser.sid,
      email: jwtUser.email,
      name: jwtUser.name,
      department: jwtUser.department,
      isAuthenticated: true,
      token,
      loading: false,
      error: null,
    })
    get().fetchProfile()
  },

  logout: () => {
    removeToken()
    set({
      sid: '',
      email: '',
      name: '',
      department: '',
      isAuthenticated: false,
      user: null,
      token: null,
      error: null,
    })
    if (get().isEmbed) {
      window.parent.postMessage({ type: 'erlink_games_logout' }, '*')
    }
  },

  fetchProfile: async () => {
    const { sid, token } = get()
    if (!sid || !token) return
    try {
      const user = await getProfile(sid)
      set({ user, error: null })
    } catch {
      // Profile fetch failed — user still authenticated, just no profile data
      set({ error: 'Failed to load profile' })
    }
  },

  getInitials: () => {
    return getInitials(get().name || get().sid)
  },

  setTheme: (theme) => set({ theme }),

  setError: (error) => set({ error }),

  getCurrentUser: () => {
    const s = get()
    return s.isAuthenticated ? s : null
  },
}))

// ─── Exported utility helpers ─────────────────────────────────────────────────

export { parseJWT, isTokenExpired, getAuthHeaders, getInitials }
