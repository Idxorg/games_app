import { motion } from 'framer-motion'
import { Clock, TrendingUp, TrendingDown, Minus, Trophy } from 'lucide-react'
import { players } from '../data/players'
import { matches, Match } from '../data/matches'
import { useUserStore } from '../stores/userStore'

const resultConfig = {
  win: { label: 'Победа', badgeClass: 'result-badge-win', icon: TrendingUp, colorClass: 'text-success' },
  loss: { label: 'Поражение', badgeClass: 'result-badge-loss', icon: TrendingDown, colorClass: 'text-danger' },
  draw: { label: 'Ничья', badgeClass: 'result-badge-draw', icon: Minus, colorClass: 'text-muted' },
}

const gameNames: Record<string, string> = {
  chess: 'Шахматы',
  checkers: 'Шашки',
  backgammon: 'Нарды',
  trivia: 'Викторины',
}

function MatchEntry({ match, index, currentUserId }: { match: Match; index: number; currentUserId: string }) {
  const opponent = match.player1Id === currentUserId ? match.player2Id : match.player1Id
  const opponentPlayer = players.find((p) => p.sid === opponent)
  const config = resultConfig[match.result]

  return (
    <motion.div
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: index * 0.05 }}
      className="glass-card p-4 mb-2 flex items-center gap-4"
    >
      {/* Result Badge */}
      <div className={`result-badge ${config.badgeClass}`}>
        <config.icon size={12} />
        {config.label}
      </div>

      {/* Game */}
      <div className="game-badge">
        {gameNames[match.game] || match.game}
      </div>

      {/* Opponent */}
      <div className="flex items-center gap-2 flex-grow min-w-0">
        <div className="profile-avatar-sm">
          {opponentPlayer?.initials || '?'}
        </div>
        <div className="min-w-0">
          <div className="text-sm font-semibold truncate">{opponentPlayer?.name || opponent}</div>
          <div className="text-xs text-muted">{opponentPlayer?.department}</div>
        </div>
      </div>

      {/* ELO Change */}
      <motion.div
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ type: 'spring', delay: index * 0.05 + 0.3 }}
        className={`flex-shrink-0 font-mono font-bold text-sm ${config.colorClass}`}
      >
        {match.eloChange > 0 ? '+' : ''}{match.eloChange}
      </motion.div>

      {/* Duration */}
      <div className="hidden sm:flex items-center gap-1 text-xs flex-shrink-0 text-muted">
        <Clock size={12} />
        {match.duration}
      </div>
    </motion.div>
  )
}

export function MatchHistory() {
  const currentUserId = useUserStore((s) => s.currentUserId)

  // Group matches by date
  const grouped = matches.reduce<Record<string, Match[]>>((acc, match) => {
    if (!acc[match.date]) acc[match.date] = []
    acc[match.date].push(match)
    return acc
  }, {})

  const sortedDates = Object.keys(grouped).sort((a, b) => b.localeCompare(a))

  const winsCount = matches.filter((m) => m.result === 'win').length
  const lossesCount = matches.filter((m) => m.result === 'loss').length
  const drawsCount = matches.filter((m) => m.result === 'draw').length

  return (
    <div className="py-8">
      <div className="container">
        <h1 className="text-3xl font-bold mb-2">История матчей</h1>
        <p className="mb-6 text-secondary">
          Хронология всех ваших партий с изменением рейтинга
        </p>

        {/* Summary */}
        <div className="grid grid-cols-3 gap-4 mb-8">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0 }}
            className="glass-card p-4 text-center"
          >
            <div className="text-2xl font-bold text-success">{winsCount}</div>
            <div className="text-xs text-muted">Побед</div>
          </motion.div>
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="glass-card p-4 text-center"
          >
            <div className="text-2xl font-bold text-danger">{lossesCount}</div>
            <div className="text-xs text-muted">Поражений</div>
          </motion.div>
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="glass-card p-4 text-center"
          >
            <div className="text-2xl font-bold text-muted">{drawsCount}</div>
            <div className="text-xs text-muted">Ничьих</div>
          </motion.div>
        </div>

        {/* Timeline */}
        {sortedDates.map((date) => (
          <div key={date} className="mb-8">
            <div className="flex items-center gap-3 mb-3">
              <Trophy size={16} color="var(--gold)" />
              <h2 className="text-sm font-semibold text-secondary">
                {new Date(date).toLocaleDateString('ru-RU', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
              </h2>
              <div className="flex-grow h-px" style={{ background: 'var(--bg-glass-border)' }} />
            </div>
            {grouped[date].map((match, i) => (
              <MatchEntry key={match.id} match={match} index={i} currentUserId={currentUserId} />
            ))}
          </div>
        ))}
      </div>
    </div>
  )
}
