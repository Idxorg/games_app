import { get, type PlayerRating } from './client'

interface RatingsResponse {
  game_type: string
  ratings: PlayerRating[]
  count: number
  limit?: number
}

export async function getRatings(gameType: string, limit = 50): Promise<PlayerRating[]> {
  const res = await get<RatingsResponse>(`/ratings/${gameType}?limit=${limit}`)
  return res.ratings
}

interface LeaderboardResponse {
  game_type: string
  leaderboard: PlayerRating[]
  count: number
}

export async function getLeaderboard(
  gameType: string,
  limit = 50,
): Promise<PlayerRating[]> {
  const res = await get<LeaderboardResponse>(
    `/ratings/${gameType}/leaderboard?limit=${limit}`,
  )
  return res.leaderboard
}

interface DepartmentLeaderboardResponse {
  game_type: string
  department: string
  ratings: PlayerRating[]
  count: number
}

export async function getLeaderboardByDepartment(
  gameType: string,
  department: string,
): Promise<PlayerRating[]> {
  const res = await get<DepartmentLeaderboardResponse>(
    `/ratings/${gameType}/leaderboard/department?department=${encodeURIComponent(department)}`,
  )
  return res.ratings
}
