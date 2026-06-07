import { create } from 'zustand'
import { players } from '../data/players'

interface UserStore {
  currentUserId: string
  isAuthenticated: boolean
  login: () => void
  logout: () => void
  getCurrentUser: () => typeof players[0] | undefined
}

export const useUserStore = create<UserStore>((set, get) => ({
  currentUserId: 'p1',
  isAuthenticated: true,
  login: () => set({ isAuthenticated: true }),
  logout: () => set({ isAuthenticated: false }),
  getCurrentUser: () => {
    return players.find((p) => p.sid === get().currentUserId)
  },
}))
