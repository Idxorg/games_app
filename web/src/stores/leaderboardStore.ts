import { create } from 'zustand'
import { players as staticPlayers, type Player } from '../data/players'
import { getLeaderboard } from '../api/ratings'
import { type PlayerRating } from '../api/client'

function mapRatingToPlayer(r: PlayerRating): Player {
  // We may not have the full player info from the API, so build a minimal one
  const staticMatch = staticPlayers.find((p) => p.sid === r.sid)
  return {
    sid: r.sid,
    name: staticMatch?.name ?? r.sid,
    department: staticMatch?.department ?? '',
    initials: staticMatch?.initials ?? r.sid.slice(0, 2).toUpperCase(),
    elo: { [r.game_type]: r.elo },
    wins: r.wins,
    losses: r.losses,
    draws: r.draws,
    trend: 'same' as const,
  }
}

type GameFilter = 'all' | 'chess' | 'checkers' | 'backgammon' | 'trivia'
type PeriodFilter = 'all' | 'week' | 'month' | 'year'

interface LeaderboardStore {
  players: Player[]
  gameFilter: GameFilter
  periodFilter: PeriodFilter
  page: number
  perPage: number
  loading: boolean
  fetchLeaderboard: (gameType?: string) => Promise<void>
  setGameFilter: (f: GameFilter) => void
  setPeriodFilter: (f: PeriodFilter) => void
  setPage: (p: number) => void
  getFiltered: () => Player[]
  getPaged: () => Player[]
  getTotalPages: () => number
}

export const useLeaderboardStore = create<LeaderboardStore>((set, get) => ({
  players: staticPlayers,
  gameFilter: 'all',
  periodFilter: 'all',
  page: 1,
  perPage: 10,
  loading: false,
  fetchLeaderboard: async (gameType?: string) => {
    set({ loading: true })
    try {
      // If a specific game type is selected, fetch its leaderboard
      const type = gameType || get().gameFilter
      if (type !== 'all') {
        const ratings = await getLeaderboard(type)
        const mapped = ratings.map(mapRatingToPlayer)
        // Merge with static data to preserve names/departments we already know
        set({ players: mapped })
      }
      // For 'all', keep the static fallback (no combined endpoint exists)
    } catch {
      // Keep static fallback data
    } finally {
      set({ loading: false })
    }
  },
  setGameFilter: (gameFilter) => {
    set({ gameFilter, page: 1 })
    // Auto-fetch when changing game filter
    if (gameFilter !== 'all') {
      get().fetchLeaderboard(gameFilter)
    } else {
      set({ players: staticPlayers })
    }
  },
  setPeriodFilter: (periodFilter) => set({ periodFilter, page: 1 }),
  setPage: (page) => set({ page }),
  getFiltered: () => {
    const { players, gameFilter } = get()
    if (gameFilter === 'all') return [...players]
    return [...players].sort((a, b) => (b.elo[gameFilter] || 0) - (a.elo[gameFilter] || 0))
  },
  getPaged: () => {
    const filtered = get().getFiltered()
    const { page, perPage } = get()
    return filtered.slice((page - 1) * perPage, page * perPage)
  },
  getTotalPages: () => {
    return Math.ceil(get().getFiltered().length / get().perPage)
  },
}))
