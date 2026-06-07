import { create } from 'zustand'
import { tournaments as staticTournaments, type Tournament } from '../data/tournaments'
import { listTournaments } from '../api/tournaments'
import { type Tournament as ApiTournament } from '../api/client'

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
  fetchTournaments: () => Promise<void>
  setActiveTab: (tab: TournamentTab) => void
}

export const useTournamentStore = create<TournamentStore>((set) => ({
  tournaments: staticTournaments,
  activeTab: 'active',
  loading: false,
  fetchTournaments: async () => {
    set({ loading: true })
    try {
      const apiTournaments = await listTournaments()
      const mapped = apiTournaments.map(mapApiTournament)
      set({ tournaments: mapped })
    } catch {
      // Keep static fallback data
    } finally {
      set({ loading: false })
    }
  },
  setActiveTab: (activeTab) => set({ activeTab }),
}))
