import { create } from 'zustand'
import { getLeaderboard } from '../api/ratings'
import { type PlayerRating } from '../api/client'

export interface Player {
  sid: string
  name: string
  department: string
  initials: string
  elo: Record<string, number>
  wins: number
  losses: number
  draws: number
  trend: 'up' | 'down' | 'same'
}

function mapRatingToPlayer(r: PlayerRating): Player {
  return {
    sid: r.sid,
    name: r.sid,
    department: '',
    initials: r.sid.slice(0, 2).toUpperCase(),
    elo: { [r.game_type]: r.elo },
    wins: r.wins,
    losses: r.losses,
    draws: r.draws,
    trend: 'same' as const,
  }
}

type GameFilter = 'all' | 'chess' | 'checkers' | 'backgammon' | 'trivia'
type PeriodFilter = 'all' | 'week' | 'month' | 'year'

const GAME_TYPES: GameFilter[] = ['chess', 'checkers', 'backgammon', 'trivia']

interface LeaderboardStore {
  players: Player[]
  gameFilter: GameFilter
  periodFilter: PeriodFilter
  page: number
  perPage: number
  loading: boolean
  error: string | null
  fetchLeaderboard: (gameType?: string) => Promise<void>
  setGameFilter: (f: GameFilter) => void
  setPeriodFilter: (f: PeriodFilter) => void
  setPage: (p: number) => void
  getFiltered: () => Player[]
  getPaged: () => Player[]
  getTotalPages: () => number
}

export const useLeaderboardStore = create<LeaderboardStore>((set, get) => ({
  players: [],
  gameFilter: 'all',
  periodFilter: 'all',
  page: 1,
  perPage: 10,
  loading: false,
  error: null,
  fetchLeaderboard: async (gameType?: string) => {
    set({ loading: true, error: null })
    try {
      const type = gameType || get().gameFilter
      if (type !== 'all') {
        const ratings = await getLeaderboard(type)
        const mapped = ratings.map(mapRatingToPlayer)
        set({ players: mapped })
      } else {
        // Fetch all game types and merge by sid
        const results = await Promise.allSettled(
          GAME_TYPES.map((gt) => getLeaderboard(gt)),
        )
        const merged = new Map<string, Player>()
        for (const result of results) {
          if (result.status === 'fulfilled') {
            for (const r of result.value) {
              const existing = merged.get(r.sid)
              if (existing) {
                existing.elo[r.game_type] = r.elo
                existing.wins += r.wins
                existing.losses += r.losses
                existing.draws += r.draws
              } else {
                merged.set(r.sid, mapRatingToPlayer(r))
              }
            }
          }
        }
        set({ players: Array.from(merged.values()) })
      }
    } catch {
      set({ error: 'Не удалось загрузить рейтинг' })
    } finally {
      set({ loading: false })
    }
  },
  setGameFilter: (gameFilter) => {
    set({ gameFilter, page: 1 })
    get().fetchLeaderboard(gameFilter)
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
