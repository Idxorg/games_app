import { create } from 'zustand'
import { players } from '../data/players'
import { getProfile } from '../api/auth'
import { getToken, type User } from '../api/client'

interface UserStore {
  currentUserId: string
  isAuthenticated: boolean
  user: User | null
  login: () => void
  logout: () => void
  fetchProfile: () => Promise<void>
  getCurrentUser: () => typeof players[0] | undefined
}

export const useUserStore = create<UserStore>((set, get) => ({
  currentUserId: 'p1',
  isAuthenticated: true,
  user: null,
  login: () => set({ isAuthenticated: true }),
  logout: () => set({ isAuthenticated: false, user: null }),
  fetchProfile: async () => {
    const sid = get().currentUserId
    const token = getToken()
    if (!token) return

    try {
      const user = await getProfile(sid)
      set({ user })
    } catch {
      // Silently fall back to static data — getCurrentUser() still works
    }
  },
  getCurrentUser: () => {
    return players.find((p) => p.sid === get().currentUserId)
  },
}))
