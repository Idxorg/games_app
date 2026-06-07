import { motion } from 'framer-motion'
import { Clock, TrendingUp, TrendingDown, Minus, Trophy } from 'lucide-react'
import { players } from '../data/players'
import { matches, Match } from '../data/matches'

const resultConfig = {
  win: { label: 'Победа', color: 'var(--success)', bg: 'rgba(74,222,128,0.1)', border: 'rgba(74,222,128,0.2)', icon: TrendingUp },
  loss: { label: 'Поражение', color: 'var(--danger)', bg: 'rgba(248,113,113,0.1)', border: 'rgba(248,113,113,0.2)', icon: TrendingDown },
  draw: { label: 'Ничья', color: 'var(--text-muted)', bg: 'rgba(255,255,255,0.05)', border: 'rgba(255,255,255,0.1)', icon: Minus },
}

const gameNames: Record<string, string> = {
  chess: 'Шахматы',
  checkers: 'Шашки',
  backgammon: 'Нарды',
  trivia: 'Викторины',
}

function MatchEntry({ match, index }: { match: Match; index: number }) {
  const opponent = match.player1Id === 'p1' ? match.player2Id : match.player1Id
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
      <div
        className="flex-shrink-0 flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-bold"
        style={{ background: config.bg, color: config.color, border: `1px solid ${config.border}` }}
      >
        <config.icon size={12} />
        {config.label}
      </div>

      {/* Game */}
      <div className="flex-shrink-0 text-xs px-2 py-1 rounded" style={{ background: 'var(--bg-tertiary)', color: 'var(--text-muted)' }}>
        {gameNames[match.game] || match.game}
      </div>

      {/* Opponent */}
      <div className="flex items-center gap-2 flex-grow min-w-0">
        <div
          className="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0"
          style={{ background: 'linear-gradient(135deg, var(--gold), var(--gold-dark))', color: '#0a0a0f' }}
        >
          {opponentPlayer?.initials || '?'}
        </div>
        <div className="min-w-0">
          <div className="text-sm font-semibold truncate">{opponentPlayer?.name || opponent}</div>
          <div className="text-xs" style={{ color: 'var(--text-muted)' }}>{opponentPlayer?.department}</div>
        </div>
      </div>

      {/* ELO Change */}
      <motion.div
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ type: 'spring', delay: index * 0.05 + 0.3 }}
        className="flex-shrink-0 font-mono font-bold text-sm"
        style={{
          color: match.eloChange > 0 ? 'var(--success)' : match.eloChange < 0 ? 'var(--danger)' : 'var(--text-muted)',
        }}
      >
        {match.eloChange > 0 ? '+' : ''}{match.eloChange}
      </motion.div>

      {/* Duration */}
      <div className="hidden sm:flex items-center gap-1 text-xs flex-shrink-0" style={{ color: 'var(--text-muted)' }}>
        <Clock size={12} />
        {match.duration}
      </div>
    </motion.div>
  )
}

export function MatchHistory() {
  // Group matches by date
  const grouped = matches.reduce<Record<string, Match[]>>((acc, match) => {
    if (!acc[match.date]) acc[match.date] = []
    acc[match.date].push(match)
    return acc
  }, {})

  const sortedDates = Object.keys(grouped).sort((a, b) => b.localeCompare(a))

  return (
    <div className="py-8">
      <div className="container">
        <h1 className="text-3xl font-bold mb-2">История матчей</h1>
        <p className="mb-6" style={{ color: 'var(--text-secondary)' }}>
          Хронология всех ваших партий с изменением рейтинга
        </p>

        {/* Summary */}
        <div className="grid grid-cols-3 gap-4 mb-8">
          <div className="glass-card p-4 text-center">
            <div className="text-2xl font-bold" style={{ color: 'var(--success)' }}>
              {matches.filter((m) => m.result === 'win').length}
            </div>
            <div className="text-xs" style={{ color: 'var(--text-muted)' }}>Побед</div>
          </div>
          <div className="glass-card p-4 text-center">
            <div className="text-2xl font-bold" style={{ color: 'var(--danger)' }}>
              {matches.filter((m) => m.result === 'loss').length}
            </div>
            <div className="text-xs" style={{ color: 'var(--text-muted)' }}>Поражений</div>
          </div>
          <div className="glass-card p-4 text-center">
            <div className="text-2xl font-bold" style={{ color: 'var(--text-muted)' }}>
              {matches.filter((m) => m.result === 'draw').length}
            </div>
            <div className="text-xs" style={{ color: 'var(--text-muted)' }}>Ничьих</div>
          </div>
        </div>

        {/* Timeline */}
        {sortedDates.map((date) => (
          <div key={date} className="mb-8">
            <div className="flex items-center gap-3 mb-3">
              <Trophy size={16} color="var(--gold)" />
              <h2 className="text-sm font-semibold" style={{ color: 'var(--text-secondary)' }}>
                {new Date(date).toLocaleDateString('ru-RU', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
              </h2>
              <div className="flex-grow h-px" style={{ background: 'var(--bg-glass-border)' }} />
            </div>
            {grouped[date].map((match, i) => (
              <MatchEntry key={match.id} match={match} index={i} />
            ))}
          </div>
        ))}
      </div>
    </div>
  )
}
