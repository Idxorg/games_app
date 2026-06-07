import { create } from 'zustand'
import { players, type Player } from '../data/players'

type GameFilter = 'all' | 'chess' | 'checkers' | 'backgammon' | 'trivia'
type PeriodFilter = 'all' | 'week' | 'month' | 'year'

interface LeaderboardStore {
  players: Player[]
  gameFilter: GameFilter
  periodFilter: PeriodFilter
  page: number
  perPage: number
  setGameFilter: (f: GameFilter) => void
  setPeriodFilter: (f: PeriodFilter) => void
  setPage: (p: number) => void
  getFiltered: () => Player[]
  getPaged: () => Player[]
  getTotalPages: () => number
}

export const useLeaderboardStore = create<LeaderboardStore>((set, get) => ({
  players,
  gameFilter: 'all',
  periodFilter: 'all',
  page: 1,
  perPage: 10,
  setGameFilter: (gameFilter) => set({ gameFilter, page: 1 }),
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
