import { motion } from 'framer-motion'
import { TrendingUp, TrendingDown, Minus } from 'lucide-react'
import { EloBar } from './EloBar'
import type { Player } from '../../data/players'

interface LeaderboardRowProps {
  player: Player
  rank: number
  gameFilter: string
}

export function LeaderboardRow({ player, rank, gameFilter }: LeaderboardRowProps) {
  const elo = gameFilter === 'all' ? 0 : (player.elo[gameFilter] || 0)
  const maxElo = 2500
  const percentile = Math.min((elo / maxElo) * 100, 100)

  const badgeColor = rank === 1 ? 'gold' : rank === 2 ? 'silver' : rank === 3 ? 'bronze' : null
  const badgeStyles: Record<string, { bg: string; color: string; border: string }> = {
    gold: { bg: 'rgba(212,168,67,0.15)', color: '#d4a843', border: 'rgba(212,168,67,0.4)' },
    silver: { bg: 'rgba(160,160,176,0.15)', color: '#a0a0b0', border: 'rgba(160,160,176,0.4)' },
    bronze: { bg: 'rgba(205,127,50,0.15)', color: '#cd7f32', border: 'rgba(205,127,50,0.4)' },
  }

  return (
    <motion.div
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: rank * 0.05 }}
      className="glass-card p-4 mb-2 flex items-center gap-4"
    >
      {/* Rank */}
      <div className="flex-shrink-0 w-10 h-10 flex items-center justify-center rounded-lg font-bold text-sm"
        style={badgeColor ? {
          background: badgeStyles[badgeColor].bg,
          color: badgeStyles[badgeColor].color,
          border: `1px solid ${badgeStyles[badgeColor].border}`,
        } : {
          background: 'var(--bg-tertiary)',
          color: 'var(--text-muted)',
        }}
      >
        {rank}
      </div>

      {/* Avatar */}
      <div className="flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold"
        style={{ background: 'linear-gradient(135deg, var(--gold), var(--gold-dark))', color: '#0a0a0f' }}
      >
        {player.initials}
      </div>

      {/* Name + Department */}
      <div className="flex-grow min-w-0">
        <div className="font-semibold text-sm truncate" style={{ color: 'var(--text-primary)' }}>{player.name}</div>
        <div className="text-xs" style={{ color: 'var(--text-muted)' }}>{player.department}</div>
      </div>

      {/* ELO Bar */}
      {gameFilter !== 'all' && (
        <div className="hidden sm:block flex-shrink-0 w-32">
          <EloBar value={elo} max={maxElo} />
        </div>
      )}

      {/* ELO Number */}
      {gameFilter !== 'all' && (
        <div className="flex-shrink-0 font-mono font-bold text-sm w-16 text-right" style={{ color: 'var(--text-primary)' }}>
          {elo}
        </div>
      )}

      {/* Trend */}
      <div className="flex-shrink-0 w-8">
        {player.trend === 'up' && <TrendingUp size={18} color="var(--success)" />}
        {player.trend === 'down' && <TrendingDown size={18} color="var(--danger)" />}
        {player.trend === 'same' && <Minus size={18} color="var(--text-muted)" />}
      </div>

      {/* W/L/D */}
      <div className="hidden md:flex flex-shrink-0 gap-3 text-xs">
        <span style={{ color: 'var(--success)' }}>{player.wins}В</span>
        <span style={{ color: 'var(--danger)' }}>{player.losses}П</span>
        <span style={{ color: 'var(--text-muted)' }}>{player.draws}Н</span>
      </div>
    </motion.div>
  )
}
