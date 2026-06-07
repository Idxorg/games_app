import { create } from 'zustand'
import { getMatchHistory } from '../api/games'
import { type Match } from '../api/client'

export interface MatchEntry {
  id: string
  player1Sid: string
  player2Sid: string
  gameType: string
  date: string
  winnerSid: string | null
  score: string
  status: string
}

function mapApiMatch(m: Match): MatchEntry {
  return {
    id: m.id,
    player1Sid: m.player1_sid,
    player2Sid: m.player2_sid,
    gameType: m.game_type,
    date: m.completed_at ?? m.created_at,
    winnerSid: m.winner_sid || null,
    score: m.score,
    status: m.status,
  }
}

interface MatchStore {
  matches: MatchEntry[]
  loading: boolean
  error: string | null
  fetchMatches: (gameType?: string) => Promise<void>
}

export const useMatchStore = create<MatchStore>((set) => ({
  matches: [],
  loading: false,
  error: null,
  fetchMatches: async (gameType?: string) => {
    set({ loading: true, error: null })
    try {
      const apiMatches = await getMatchHistory(gameType)
      const mapped = apiMatches.map(mapApiMatch)
      set({ matches: mapped })
    } catch {
      set({ error: 'Не удалось загрузить историю матчей' })
    } finally {
      set({ loading: false })
    }
  },
}))
