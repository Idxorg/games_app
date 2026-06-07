import { motion } from 'framer-motion'
import { Building, Trophy, Target, Flame, Star, Settings } from 'lucide-react'
import { useUserStore } from '../stores/userStore'
import { EloBar } from '../components/ui/EloBar'

const achievements = [
  { name: 'Первые шаги', description: 'Сыграйте свою первую партию', icon: Target, unlocked: true },
  { name: 'Шахматный стратег', description: 'Выиграйте 10 партий в шахматы', icon: Trophy, unlocked: true },
  { name: 'На высоте', description: 'Достигните рейтинга 2000+', icon: Star, unlocked: true },
  { name: 'Неудержимый', description: 'Серия из 5 побед подряд', icon: Flame, unlocked: false },
  { name: 'Мастер всех игр', description: 'Сыграйте во все 7 игр', icon: Settings, unlocked: false },
]

const eloData = [1800, 1850, 1820, 1900, 1950, 1920, 1980, 2010, 2050, 2080, 2060, 2100, 2120, 2150]

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

  if (!user) {
    return (
      <div className="py-8">
        <div className="container">
          <ProfileSkeleton />
        </div>
      </div>
    )
  }

  return (
    <div className="py-8">
      <div className="container">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Sidebar */}
          <div className="lg:col-span-1">
            <div className="glass-card p-6 text-center mb-4">
              <div className="profile-avatar">
                {user.initials}
              </div>
              <h2 className="text-xl font-bold mb-1">{user.name}</h2>
              <div className="flex items-center gap-2 justify-center text-sm mb-4 text-muted">
                <Building size={14} />
                {user.department}
              </div>
            </div>

            {/* Stats */}
            <div className="glass-card p-6">
              <h3 className="text-sm font-semibold mb-4 text-secondary">Статистика</h3>
              <div className="flex flex-col gap-4">
                <div className="flex justify-between text-sm">
                  <span className="text-muted">Всего партий</span>
                  <span className="font-bold">{user.wins + user.losses + user.draws}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-success">Побед</span>
                  <span className="font-bold text-success">{user.wins}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-danger">Поражений</span>
                  <span className="font-bold text-danger">{user.losses}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted">Ничьих</span>
                  <span className="font-bold">{user.draws}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted">Винрейт</span>
                  <span className="font-bold text-accent">
                    {Math.round((user.wins / (user.wins + user.losses + user.draws)) * 100)}%
                  </span>
                </div>
              </div>
            </div>

            {/* ELO per Game */}
            <div className="glass-card p-6 mt-4">
              <h3 className="text-sm font-semibold mb-4 text-secondary">Рейтинг по играм</h3>
              <div className="flex flex-col gap-3">
                {Object.entries(user.elo).map(([game, elo]) => (
                  <div key={game}>
                    <div className="text-xs mb-1 text-muted">
                      {game === 'chess' ? 'Шахматы' : game === 'checkers' ? 'Шашки' : game === 'backgammon' ? 'Нарды' : 'Викторины'}
                    </div>
                    <EloBar value={elo} max={2500} />
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Right Content */}
          <div className="lg:col-span-2">
            {/* ELO Chart */}
            <div className="glass-card p-6 mb-4">
              <h3 className="text-sm font-semibold mb-4 text-secondary">
                Динамика рейтинга (шахматы)
              </h3>
              <div className="flex items-end gap-1" style={{ height: 120 }}>
                {eloData.map((elo, i) => {
                  const min = Math.min(...eloData) - 50
                  const max = Math.max(...eloData) + 50
                  const height = ((elo - min) / (max - min)) * 100
                  return (
                    <motion.div
                      key={i}
                      initial={{ height: 0 }}
                      animate={{ height: `${Math.max(height, 5)}%` }}
                      transition={{ delay: i * 0.05, duration: 0.3 }}
                      className={`elo-bar ${i === eloData.length - 1 ? 'elo-bar-active' : 'elo-bar-inactive'}`}
                    />
                  )
                })}
              </div>
              <div className="flex justify-between mt-2 text-xs text-muted">
                <span>14 нед. назад</span>
                <span>Сегодня</span>
              </div>
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
