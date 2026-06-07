import { create } from 'zustand'
import { games, type Game } from '../data/games'
import { availableGames } from '../api/games'
import { type ApiGame } from '../api/client'

function mapApiGameToGame(g: ApiGame): Game {
  return {
    id: g.type,
    name: g.name,
    description: g.description,
    iconType: g.type,
    playerCount:
      g.min_players === g.max_players
        ? `${g.min_players} игрок`
        : `${g.min_players}–${g.max_players} игроков`,
    category: '',
    isLive: true,
    route: `/game/${g.type}`,
  }
}

interface GameStore {
  games: Game[]
  selectedGame: Game | null
  searchQuery: string
  loading: boolean
  fetchGames: () => Promise<void>
  setSearchQuery: (q: string) => void
  selectGame: (game: Game) => void
}

export const useGameStore = create<GameStore>((set) => ({
  games,
  selectedGame: null,
  searchQuery: '',
  loading: false,
  fetchGames: async () => {
    set({ loading: true })
    try {
      const apiGames = await availableGames()
      const mapped = apiGames.map(mapApiGameToGame)
      set({ games: mapped })
    } catch {
      // Keep static fallback data
    } finally {
      set({ loading: false })
    }
  },
  setSearchQuery: (searchQuery) => set({ searchQuery }),
  selectGame: (selectedGame) => set({ selectedGame }),
}))
