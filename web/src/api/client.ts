import { useToastStore } from '../components/ui/Toast'

const TOKEN_KEY = 'gz_token'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

const BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

// ─── Shared types matching Go backend models ──────────────────────────────────

export interface User {
  sid: string
  email: string
  name: string
  gender: string
  department: string
  position: string
  photo_url: string
  last_sync: string
  created_at: string
  updated_at: string
}

export interface Tournament {
  id: string
  name: string
  game_type: string
  status: string
  start_date: string
  end_date: string
  max_players: number
  current_players: number
  prize_pool: string
  description: string
  logo_url: string
  created_by: string
  requires_group: string
  created_at: string
}

export interface TournamentPlayer {
  id: number
  tournament_id: string
  sid: string
  rank: number
  points: number
  wins: number
  draws: number
  losses: number
  joined_at: string
}

export interface Match {
  id: string
  tournament_id: string
  game_type: string
  player1_sid: string
  player2_sid: string
  winner_sid: string
  score: string
  moves: unknown
  pgn_url: string
  game_id: string
  livekit_room_id: string
  status: string
  started_at: string | null
  completed_at: string | null
  created_at: string
}

export interface PlayerRating {
  id: number
  sid: string
  game_type: string
  elo: number
  games_played: number
  wins: number
  draws: number
  losses: number
}

export interface ApiGame {
  type: string
  name: string
  description: string
  min_players: number
  max_players: number
}

// ─── API Error ───────────────────────────────────────────────────────────────

export class ApiError extends Error {
  status: number
  body: unknown

  constructor(status: number, body: unknown) {
    const msg =
      typeof body === 'object' && body !== null && 'error' in body
        ? (body as { error: string }).error
        : `Request failed (${status})`
    super(msg)
    this.name = 'ApiError'
    this.status = status
    this.body = body
  }
}

// ─── HTTP helpers ────────────────────────────────────────────────────────────

async function request<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const token = getToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const url = `${BASE_URL}${path}`
  const res = await fetch(url, { ...options, headers })

  // Handle 401 — clear token and redirect (no auto-refresh for JWT)
  if (res.status === 401) {
    removeToken()
    // Avoid infinite loop if already on login-like page
    if (window.location.pathname !== '/') {
      useToastStore.getState().addToast('Сессия истекла. Пожалуйста, войдите снова.', 'error')
    }
    throw new ApiError(401, await res.json().catch(() => null))
  }

  if (!res.ok) {
    const body = await res.json().catch(() => null)
    const err = new ApiError(res.status, body)
    useToastStore.getState().addToast(err.message, 'error')
    throw err
  }

  // 204 No Content
  if (res.status === 204) {
    return undefined as unknown as T
  }

  return res.json()
}

// Convenience wrappers
export function get<T>(path: string): Promise<T> {
  return request<T>(path)
}

export function post<T>(path: string, body?: unknown): Promise<T> {
  return request<T>(path, { method: 'POST', body: body ? JSON.stringify(body) : undefined })
}

export function put<T>(path: string, body?: unknown): Promise<T> {
  return request<T>(path, { method: 'PUT', body: body ? JSON.stringify(body) : undefined })
}

export function del<T>(path: string): Promise<T> {
  return request<T>(path, { method: 'DELETE' })
}
