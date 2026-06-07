import { create } from 'zustand'
import { availableGames } from '../api/games'
import { type ApiGame } from '../api/client'

export interface Game {
  id: string
  name: string
  description: string
  iconType: string
  playerCount: string
  category: string
  isLive: boolean
  route: string
}

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
  error: string | null
  fetchGames: () => Promise<void>
  setSearchQuery: (q: string) => void
  selectGame: (game: Game) => void
}

export const useGameStore = create<GameStore>((set) => ({
  games: [],
  selectedGame: null,
  searchQuery: '',
  loading: false,
  error: null,
  fetchGames: async () => {
    set({ loading: true, error: null })
    try {
      const apiGames = await availableGames()
      const mapped = apiGames.map(mapApiGameToGame)
      set({ games: mapped })
    } catch {
      set({ error: 'Не удалось загрузить список игр' })
    } finally {
      set({ loading: false })
    }
  },
  setSearchQuery: (searchQuery) => set({ searchQuery }),
  selectGame: (selectedGame) => set({ selectedGame }),
}))
