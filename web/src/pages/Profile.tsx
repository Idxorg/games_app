import { motion } from 'framer-motion'
import { Settings, Mail, Building, Trophy, Target, Flame, Star } from 'lucide-react'
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

export function Profile() {
  const user = useUserStore((s) => s.getCurrentUser())

  if (!user) return null

  return (
    <div className="py-8">
      <div className="container">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Sidebar */}
          <div className="lg:col-span-1">
            <div className="glass-card p-6 text-center mb-4">
              <div
                className="w-24 h-24 rounded-full flex items-center justify-center text-3xl font-bold mx-auto mb-4"
                style={{ background: 'linear-gradient(135deg, var(--gold), var(--gold-dark))', color: '#0a0a0f' }}
              >
                {user.initials}
              </div>
              <h2 className="text-xl font-bold mb-1">{user.name}</h2>
              <div className="flex items-center gap-2 justify-center text-sm mb-4" style={{ color: 'var(--text-muted)' }}>
                <Building size={14} />
                {user.department}
              </div>
              <div className="flex items-center gap-2 justify-center text-sm mb-1" style={{ color: 'var(--text-muted)' }}>
                <Mail size={14} />
                {user.sid}@company.ru
              </div>
            </div>

            {/* Stats */}
            <div className="glass-card p-6">
              <h3 className="text-sm font-semibold mb-4" style={{ color: 'var(--text-secondary)' }}>Статистика</h3>
              <div className="flex flex-col gap-4">
                <div className="flex justify-between text-sm">
                  <span style={{ color: 'var(--text-muted)' }}>Всего партий</span>
                  <span className="font-bold">{user.wins + user.losses + user.draws}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span style={{ color: 'var(--success)' }}>Побед</span>
                  <span className="font-bold" style={{ color: 'var(--success)' }}>{user.wins}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span style={{ color: 'var(--danger)' }}>Поражений</span>
                  <span className="font-bold" style={{ color: 'var(--danger)' }}>{user.losses}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span style={{ color: 'var(--text-muted)' }}>Ничьих</span>
                  <span className="font-bold">{user.draws}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span style={{ color: 'var(--text-muted)' }}>Винрейт</span>
                  <span className="font-bold" style={{ color: 'var(--gold)' }}>
                    {Math.round((user.wins / (user.wins + user.losses + user.draws)) * 100)}%
                  </span>
                </div>
              </div>
            </div>

            {/* ELO per Game */}
            <div className="glass-card p-6 mt-4">
              <h3 className="text-sm font-semibold mb-4" style={{ color: 'var(--text-secondary)' }}>Рейтинг по играм</h3>
              <div className="flex flex-col gap-3">
                {Object.entries(user.elo).map(([game, elo]) => (
                  <div key={game}>
                    <div className="text-xs mb-1" style={{ color: 'var(--text-muted)' }}>
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
            {/* ELO Chart (simple bar chart) */}
            <div className="glass-card p-6 mb-4">
              <h3 className="text-sm font-semibold mb-4" style={{ color: 'var(--text-secondary)' }}>
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
                      className="flex-1 rounded-t"
                      style={{
                        background: i === eloData.length - 1
                          ? 'linear-gradient(180deg, var(--gold), var(--gold-dark))'
                          : 'rgba(212,168,67,0.15)',
                        minWidth: 4,
                      }}
                    />
                  )
                })}
              </div>
              <div className="flex justify-between mt-2 text-xs" style={{ color: 'var(--text-muted)' }}>
                <span>14 нед. назад</span>
                <span>Сегодня</span>
              </div>
            </div>

            {/* Achievements */}
            <div className="glass-card p-6">
              <h3 className="text-sm font-semibold mb-4" style={{ color: 'var(--text-secondary)' }}>Достижения</h3>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {achievements.map((ach, i) => (
                  <motion.div
                    key={ach.name}
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: i * 0.1 }}
                    className="flex items-center gap-3 p-3 rounded-lg"
                    style={{
                      background: ach.unlocked ? 'rgba(212,168,67,0.08)' : 'rgba(255,255,255,0.02)',
                      border: `1px solid ${ach.unlocked ? 'rgba(212,168,67,0.2)' : 'var(--bg-glass-border)'}`,
                      opacity: ach.unlocked ? 1 : 0.5,
                    }}
                  >
                    <div
                      className="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0"
                      style={{
                        background: ach.unlocked ? 'rgba(212,168,67,0.15)' : 'var(--bg-tertiary)',
                      }}
                    >
                      <ach.icon size={20} color={ach.unlocked ? 'var(--gold)' : 'var(--text-muted)'} />
                    </div>
                    <div>
                      <div className="text-sm font-semibold">{ach.name}</div>
                      <div className="text-xs" style={{ color: 'var(--text-muted)' }}>{ach.description}</div>
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
