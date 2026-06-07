import { create } from 'zustand'
import { listTournaments } from '../api/tournaments'
import { type Tournament as ApiTournament } from '../api/client'

export interface Tournament {
  id: string
  name: string
  game: string
  status: 'active' | 'upcoming' | 'completed'
  startDate: string
  endDate: string
  participants: string[]
  rounds: {
    round: string
    matches: {
      player1: string
      score1: number
      player2: string
      score2: number
      winner?: string
    }[]
  }[]
  prize: string
}

function mapApiTournament(t: ApiTournament): Tournament {
  return {
    id: t.id,
    name: t.name,
    game: t.game_type,
    status: t.status as Tournament['status'],
    startDate: t.start_date?.slice(0, 10) ?? '',
    endDate: t.end_date?.slice(0, 10) ?? '',
    participants: [],
    rounds: [],
    prize: t.prize_pool || '',
  }
}

type TournamentTab = 'active' | 'upcoming' | 'completed'

interface TournamentStore {
  tournaments: Tournament[]
  activeTab: TournamentTab
  loading: boolean
  error: string | null
  fetchTournaments: () => Promise<void>
  setActiveTab: (tab: TournamentTab) => void
}

export const useTournamentStore = create<TournamentStore>((set) => ({
  tournaments: [],
  activeTab: 'active',
  loading: false,
  error: null,
  fetchTournaments: async () => {
    set({ loading: true, error: null })
    try {
      const apiTournaments = await listTournaments()
      const mapped = apiTournaments.map(mapApiTournament)
      set({ tournaments: mapped, error: null })
    } catch {
      set({ error: 'Не удалось загрузить турниры' })
    } finally {
      set({ loading: false })
    }
  },
  setActiveTab: (activeTab) => set({ activeTab }),
}))
