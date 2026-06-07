import { useEffect } from 'react'
import { useState } from 'react'
import { motion } from 'framer-motion'
import { Building, Trophy, Target, Flame, Star, Settings, RefreshCw } from 'lucide-react'
import { useUserStore } from '../stores/userStore'
import { EloBar } from '../components/ui/EloBar'
import { getLeaderboard } from '../api/ratings'
import { type PlayerRating } from '../api/client'

const achievements = [
  { name: 'Первые шаги', description: 'Сыграйте свою первую партию', icon: Target, unlocked: true },
  { name: 'Шахматный стратег', description: 'Выиграйте 10 партий в шахматы', icon: Trophy, unlocked: true },
  { name: 'На высоте', description: 'Достигните рейтинга 2000+', icon: Star, unlocked: true },
  { name: 'Неудержимый', description: 'Серия из 5 побед подряд', icon: Flame, unlocked: false },
  { name: 'Мастер всех игр', description: 'Сыграйте во все 7 игр', icon: Settings, unlocked: false },
]

function ProfileSkeleton() {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div className="lg:col-span-1">
        <div className="glass-card p-6 text-center mb-4">
          <div className="skeleton skeleton-circle mx-auto mb-4" style={{ width: 96, height: 96 }} />
          <div className="skeleton skeleton-title mx-auto" style={{ width: '50%' }} />
        </div>
        <div className="glass-card p-6">
          <div className="skeleton skeleton-text" />
          <div className="skeleton skeleton-text" />
          <div className="skeleton skeleton-text" />
        </div>
      </div>
      <div className="lg:col-span-2">
        <div className="glass-card p-6">
          <div className="skeleton skeleton-title" />
          <div className="flex items-end gap-1 mt-4" style={{ height: 120 }}>
            {Array.from({ length: 14 }).map((_, i) => (
              <div key={i} className="skeleton flex-1" style={{ height: '60%' }} />
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

export function Profile() {
  const user = useUserStore((s) => s.getCurrentUser())
  const getInitials = useUserStore((s) => s.getInitials)
  const fullUser = useUserStore((s) => s.user)
  const sid = useUserStore((s) => s.sid)
  const isAuthenticated = useUserStore((s) => s.isAuthenticated)
  const isEmbed = useUserStore((s) => s.isEmbed)
  const fetchProfile = useUserStore((s) => s.fetchProfile)

  const [ratings, setRatings] = useState<PlayerRating[]>([])
  const [ratingsLoading, setRatingsLoading] = useState(false)
  const [ratingsError, setRatingsError] = useState<string | null>(null)

  useEffect(() => {
    if (sid && isAuthenticated) {
      setRatingsLoading(true)
      // Fetch ratings for all game types
      Promise.allSettled([
        getLeaderboard('chess').catch(() => []),
        getLeaderboard('checkers').catch(() => []),
        getLeaderboard('backgammon').catch(() => []),
        getLeaderboard('trivia').catch(() => []),
      ])
        .then((results) => {
          const all: PlayerRating[] = []
          for (const r of results) {
            if (r.status === 'fulfilled') all.push(...r.value)
          }
          setRatings(all)
        })
        .catch(() => setRatingsError('Не удалось загрузить рейтинг'))
        .finally(() => setRatingsLoading(false))
    }
  }, [sid, isAuthenticated])

  if (isEmbed && !isAuthenticated) {
    return (
      <div className="py-8">
        <div className="container text-center">
          <p className="text-muted">Для доступа войдите через портал</p>
        </div>
      </div>
    )
  }

  if (!user) {
    return (
      <div className="py-8">
        <div className="container">
          <ProfileSkeleton />
        </div>
      </div>
    )
  }

  // Build ELO map from ratings
  const myRatings = ratings.filter((r) => r.sid === sid)
  const eloMap: Record<string, number> = {}
  for (const r of myRatings) {
    eloMap[r.game_type] = r.elo
  }

  // Aggregate stats from ratings
  const totalGames = myRatings.reduce((s, r) => s + r.games_played, 0)
  const totalWins = myRatings.reduce((s, r) => s + r.wins, 0)
  const totalLosses = myRatings.reduce((s, r) => s + r.losses, 0)
  const totalDraws = myRatings.reduce((s, r) => s + r.draws, 0)
  const winRate = totalGames > 0 ? Math.round((totalWins / totalGames) * 100) : 0

  const displayName = fullUser?.name || user.name || sid
  const displayDepartment = fullUser?.department || user.department || ''

  return (
    <div className="py-8">
      <div className="container">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Sidebar */}
          <div className="lg:col-span-1">
            <div className="glass-card p-6 text-center mb-4">
              <div className="profile-avatar">
                {getInitials()}
              </div>
              <h2 className="text-xl font-bold mb-1">{displayName}</h2>
              {displayDepartment && (
                <div className="flex items-center gap-2 justify-center text-sm mb-4 text-muted">
                  <Building size={14} />
                  {displayDepartment}
                </div>
              )}
            </div>

            {/* Stats */}
            <div className="glass-card p-6">
              <h3 className="text-sm font-semibold mb-4 text-secondary">Статистика</h3>
              <div className="flex flex-col gap-4">
                <div className="flex justify-between text-sm">
                  <span className="text-muted">Всего партий</span>
                  <span className="font-bold">{totalGames || '—'}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-success">Побед</span>
                  <span className="font-bold text-success">{totalWins || '—'}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-danger">Поражений</span>
                  <span className="font-bold text-danger">{totalLosses || '—'}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted">Ничьих</span>
                  <span className="font-bold">{totalDraws || '—'}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted">Винрейт</span>
                  <span className="font-bold text-accent">
                    {totalGames > 0 ? `${winRate}%` : '—'}
                  </span>
                </div>
              </div>
            </div>

            {/* ELO per Game */}
            <div className="glass-card p-6 mt-4">
              <h3 className="text-sm font-semibold mb-4 text-secondary">Рейтинг по играм</h3>
              {ratingsLoading && (
                <div className="flex flex-col gap-3">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="skeleton skeleton-text" style={{ width: '60%' }} />
                  ))}
                </div>
              )}
              {ratingsError && (
                <div className="text-center">
                  <p className="text-xs text-muted mb-2">{ratingsError}</p>
                  <button
                    className="btn btn-ghost text-xs"
                    onClick={() => fetchProfile()}
                  >
                    <RefreshCw size={12} /> Повторить
                  </button>
                </div>
              )}
              {!ratingsLoading && !ratingsError && Object.keys(eloMap).length === 0 && (
                <p className="text-xs text-muted text-center">Рейтинг пока пуст</p>
              )}
              {!ratingsLoading && !ratingsError && Object.keys(eloMap).length > 0 && (
                <div className="flex flex-col gap-3">
                  {Object.entries(eloMap).map(([game, elo]) => (
                    <div key={game}>
                      <div className="text-xs mb-1 text-muted">
                        {game === 'chess' ? 'Шахматы' : game === 'checkers' ? 'Шашки' : game === 'backgammon' ? 'Нарды' : 'Викторины'}
                      </div>
                      <EloBar value={elo} max={2500} />
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Right Content */}
          <div className="lg:col-span-2">
            {/* ELO Chart placeholder */}
            <div className="glass-card p-6 mb-4">
              <h3 className="text-sm font-semibold mb-4 text-secondary">
                Динамика рейтинга (шахматы)
              </h3>
              {eloMap.chess != null ? (
                <>
                  <div className="flex items-center justify-center py-8">
                    <span className="text-lg font-bold">{eloMap.chess}</span>
                  </div>
                  <div className="text-center text-xs text-muted">
                    Текущий рейтинг
                  </div>
                </>
              ) : (
                <div className="flex items-center justify-center py-8">
                  <span className="text-muted">Данные появятся после первой партии</span>
                </div>
              )}
            </div>

            {/* Achievements */}
            <div className="glass-card p-6">
              <h3 className="text-sm font-semibold mb-4 text-secondary">Достижения</h3>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {achievements.map((ach, i) => (
                  <motion.div
                    key={ach.name}
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: i * 0.1 }}
                    className={`achievement-card ${ach.unlocked ? 'achievement-card-unlocked' : 'achievement-card-locked'}`}
                  >
                    <div className={`achievement-icon ${ach.unlocked ? 'achievement-icon-unlocked' : ''}`}>
                      <ach.icon size={20} color={ach.unlocked ? 'var(--gold)' : 'var(--text-muted)'} />
                    </div>
                    <div>
                      <div className="text-sm font-semibold">{ach.name}</div>
                      <div className="text-xs text-muted">{ach.description}</div>
                    </div>
                  </motion.div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
