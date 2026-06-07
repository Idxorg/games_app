import { get, post, type Match, type ApiGame } from './client'

interface AvailableGamesResponse {
  games: ApiGame[]
  count: number
}

export async function availableGames(): Promise<ApiGame[]> {
  const res = await get<AvailableGamesResponse>('/games/available')
  return res.games
}

interface StartMatchData {
  game_type: string
  player2_sid: string
  tournament_id?: string
  livekit_room_id?: string
}

export async function startMatch(data: StartMatchData): Promise<Match> {
  return post<Match>('/games/match/start', data)
}

interface CompleteMatchData {
  match_id: string
  winner_sid?: string
  score: string
  moves?: unknown
}

interface CompleteMatchResponse {
  message: string
  match_id: string
}

export async function completeMatch(data: CompleteMatchData): Promise<CompleteMatchResponse> {
  return post<CompleteMatchResponse>('/games/match/complete', data)
}

interface MatchHistoryResponse {
  matches: Match[]
  count: number
}

export async function getMatchHistory(gameType?: string): Promise<Match[]> {
  const qs = gameType ? `?game_type=${gameType}` : ''
  const res = await get<MatchHistoryResponse>(`/games/match/history${qs}`)
  return res.matches
}
