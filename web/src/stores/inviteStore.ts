import { create } from 'zustand'
import {
  createInvite,
  acceptInvite,
  declineInvite,
  getPendingInvites,
  type Invite,
  type AcceptInviteResponse,
} from '../api/invites'
import { useUserStore } from './userStore'

// ─── Store interface ──────────────────────────────────────────────────────

interface InviteState {
  pendingInvites: Invite[]
  loading: boolean
  error: string | null
  sending: boolean
}

interface InviteActions {
  fetchPending: () => Promise<void>
  accept: (id: string) => Promise<AcceptInviteResponse>
  decline: (id: string) => Promise<void>
  create: (gameType: string, recipientSid: string) => Promise<void>
  clearError: () => void
}

type InviteStore = InviteState & InviteActions

// ─── Zustand store ─────────────────────────────────────────────────────────

export const useInviteStore = create<InviteStore>((set, get) => ({
  // State
  pendingInvites: [],
  loading: false,
  error: null,
  sending: false,

  // Actions

  fetchPending: async () => {
    const isAuthenticated = useUserStore.getState().isAuthenticated
    if (!isAuthenticated) return

    set({ loading: true })
    try {
      const invites = await getPendingInvites()
      set({ pendingInvites: invites, loading: false, error: null })
    } catch (err) {
      set({
        loading: false,
        error: err instanceof Error ? err.message : 'Failed to fetch invites',
      })
    }
  },

  accept: async (id: string) => {
    try {
      const res = await acceptInvite(id)
      await get().fetchPending()
      return res
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : 'Failed to accept invite',
      })
      throw err
    }
  },

  decline: async (id: string) => {
    try {
      await declineInvite(id)
      await get().fetchPending()
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : 'Failed to decline invite',
      })
      throw err
    }
  },

  create: async (gameType: string, recipientSid: string) => {
    set({ sending: true, error: null })
    try {
      await createInvite(gameType, recipientSid)
      set({ sending: false })
    } catch (err) {
      set({
        sending: false,
        error: err instanceof Error ? err.message : 'Failed to send invite',
      })
      throw err
    }
  },

  clearError: () => set({ error: null }),
}))
