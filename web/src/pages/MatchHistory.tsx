import { useEffect } from 'react'
import { motion } from 'framer-motion'
import { Clock, TrendingUp, TrendingDown, Minus, Trophy, RefreshCw } from 'lucide-react'
import { useUserStore } from '../stores/userStore'
import { useMatchStore } from '../stores/matchStore'

const gameNames: Record<string, string> = {
  chess: 'Шахматы',
  checkers: 'Шашки',
  backgammon: 'Нарды',
  trivia: 'Викторины',
}

function MatchEntry({
  match,
  index,
  currentUserId,
}: {
  match: ReturnType<typeof useMatchStore.getState>['matches'][number]
  index: number
  currentUserId: string
}) {
  const opponent = match.player1Sid === currentUserId ? match.player2Sid : match.player1Sid
  const isWin = match.winnerSid === currentUserId
  const isDraw = !match.winnerSid || match.winnerSid === '' || match.status === 'draw'
  const isLoss = !isWin && !isDraw

  const config = isDraw
    ? { label: 'Ничья', badgeClass: 'result-badge-draw', icon: Minus, colorClass: 'text-muted' }
    : isWin
      ? { label: 'Победа', badgeClass: 'result-badge-win', icon: TrendingUp, colorClass: 'text-success' }
      : { label: 'Поражение', badgeClass: 'result-badge-loss', icon: TrendingDown, colorClass: 'text-danger' }

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
        {gameNames[match.gameType] || match.gameType}
      </div>

      {/* Opponent */}
      <div className="flex items-center gap-2 flex-grow min-w-0">
        <div className="profile-avatar-sm">
          {opponent ? opponent.slice(0, 2).toUpperCase() : '?'}
        </div>
        <div className="min-w-0">
          <div className="text-sm font-semibold truncate">{opponent || '—'}</div>
          <div className="text-xs text-muted">&nbsp;</div>
        </div>
      </div>

      {/* Duration */}
      {match.date && (
        <div className="hidden sm:flex items-center gap-1 text-xs flex-shrink-0 text-muted">
          <Clock size={12} />
          {new Date(match.date).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })}
        </div>
      )}
    </motion.div>
  )
}

function MatchSkeleton() {
  return (
    <div className="flex flex-col gap-2">
      {[1, 2, 3, 4].map((i) => (
        <div key={i} className="glass-card p-4 flex items-center gap-4">
          <div className="skeleton" style={{ width: 80, height: 28, borderRadius: 'var(--radius-sm)' }} />
          <div className="skeleton" style={{ width: 60, height: 24, borderRadius: 'var(--radius-sm)' }} />
          <div className="flex items-center gap-2 flex-grow">
            <div className="skeleton skeleton-circle" style={{ width: 32, height: 32 }} />
            <div className="skeleton skeleton-text" style={{ width: '40%' }} />
          </div>
        </div>
      ))}
    </div>
  )
}

export function MatchHistory() {
  const currentUserId = useUserStore((s) => s.sid)
  const isAuthenticated = useUserStore((s) => s.isAuthenticated)
  const isEmbed = useUserStore((s) => s.isEmbed)
  const { matches, loading, error, fetchMatches } = useMatchStore()

  useEffect(() => {
    if (isAuthenticated) {
      fetchMatches()
    }
  }, [isAuthenticated, fetchMatches])

  if (isEmbed && !isAuthenticated) {
    return (
      <div className="py-8">
        <div className="container text-center">
          <p className="text-muted">Для доступа войдите через портал</p>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return (
      <div className="py-8">
        <div className="container">
          <h1 className="text-3xl font-bold mb-2">История матчей</h1>
          <p className="mb-6 text-secondary">
            Хронология всех ваших партий с изменением рейтинга
          </p>
          <div className="glass-card p-12 text-center">
            <Trophy size={48} color="var(--text-muted)" className="mx-auto mb-4" />
            <p className="text-muted">Войдите, чтобы увидеть историю матчей</p>
          </div>
        </div>
      </div>
    )
  }

  // Calculate summary
  const winsCount = matches.filter((m) => m.winnerSid === currentUserId).length
  const lossesCount = matches.filter(
    (m) => m.winnerSid && m.winnerSid !== currentUserId && m.status !== 'draw',
  ).length
  const drawsCount = matches.filter((m) => !m.winnerSid || m.status === 'draw').length

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

        {/* Loading */}
        {loading && <MatchSkeleton />}

        {/* Error */}
        {!loading && error && (
          <div className="glass-card p-12 text-center">
            <Trophy size={48} color="var(--text-muted)" className="mx-auto mb-4" />
            <p className="text-muted mb-4">{error}</p>
            <button className="btn btn-secondary" onClick={() => fetchMatches()}>
              <RefreshCw size={14} /> Повторить
            </button>
          </div>
        )}

        {/* Empty */}
        {!loading && !error && matches.length === 0 && (
          <div className="glass-card p-12 text-center">
            <Trophy size={48} color="var(--text-muted)" className="mx-auto mb-4" />
            <p className="text-muted">Нет активных матчей</p>
          </div>
        )}

        {/* Match List */}
        {!loading && !error && matches.length > 0 && (
          matches.map((match, i) => (
            <MatchEntry
              key={match.id}
              match={match}
              index={i}
              currentUserId={currentUserId}
            />
          ))
        )}
      </div>
    </div>
  )
}
