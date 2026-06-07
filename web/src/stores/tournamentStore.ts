import { create } from 'zustand'
import { tournaments, type Tournament } from '../data/tournaments'

type TournamentTab = 'active' | 'upcoming' | 'completed'

interface TournamentStore {
  tournaments: Tournament[]
  activeTab: TournamentTab
  setActiveTab: (tab: TournamentTab) => void
}

export const useTournamentStore = create<TournamentStore>((set) => ({
  tournaments,
  activeTab: 'active',
  setActiveTab: (activeTab) => set({ activeTab }),
}))
