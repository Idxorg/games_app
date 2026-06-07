import { create } from 'zustand'
import { games, type Game } from '../data/games'

interface GameStore {
  games: Game[]
  selectedGame: Game | null
  searchQuery: string
  setSearchQuery: (q: string) => void
  selectGame: (game: Game) => void
}

export const useGameStore = create<GameStore>((set) => ({
  games,
  selectedGame: null,
  searchQuery: '',
  setSearchQuery: (searchQuery) => set({ searchQuery }),
  selectGame: (selectedGame) => set({ selectedGame }),
}))
