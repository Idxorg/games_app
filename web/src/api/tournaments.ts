import { get, post, type Tournament, type TournamentPlayer } from './client'

export interface TournamentFilters {
  game_type?: string
  status?: string
  search?: string
  created_by?: string
  limit?: number
  offset?: number
}

interface ListTournamentsResponse {
  tournaments: Tournament[]
  count: number
}

export async function listTournaments(
  filters?: TournamentFilters,
): Promise<Tournament[]> {
  const params = new URLSearchParams()
  if (filters) {
    for (const [key, value] of Object.entries(filters)) {
      if (value !== undefined && value !== '') {
        params.append(key, String(value))
      }
    }
  }
  const qs = params.toString()
  const path = `/tournaments${qs ? `?${qs}` : ''}`
  const res = await get<ListTournamentsResponse>(path)
  return res.tournaments
}

interface CreateTournamentData {
  name: string
  game_type: string
  max_players: number
  start_date?: string
  end_date?: string
  prize_pool?: string
  description?: string
  logo_url?: string
  requires_group?: string
}

export async function createTournament(data: CreateTournamentData): Promise<Tournament> {
  return post<Tournament>('/tournaments', data)
}

export async function getTournament(id: string): Promise<Tournament> {
  return get<Tournament>(`/tournaments/${id}`)
}

interface MessageResponse {
  message: string
}

export async function joinTournament(id: string): Promise<MessageResponse> {
  return post<MessageResponse>(`/tournaments/${id}/join`)
}

export async function leaveTournament(id: string): Promise<MessageResponse> {
  return post<MessageResponse>(`/tournaments/${id}/leave`)
}

interface PlayersResponse {
  players: TournamentPlayer[]
  count: number
}

export async function getPlayers(id: string): Promise<TournamentPlayer[]> {
  const res = await get<PlayersResponse>(`/tournaments/${id}/players`)
  return res.players
}
