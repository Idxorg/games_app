import { motion } from 'framer-motion'

// ─── Types ──────────────────────────────────────────────────────────────────

interface GameClockProps {
  whiteMs: number
  blackMs: number
  activeColor: string
}

// ─── Helpers ────────────────────────────────────────────────────────────────

function formatTime(ms: number): string {
  const totalSec = Math.max(0, Math.floor(ms / 1000))
  const min = Math.floor(totalSec / 60)
  const sec = totalSec % 60
  return `${min}:${sec.toString().padStart(2, '0')}`
}

// ─── Component ──────────────────────────────────────────────────────────────

export function GameClock({ whiteMs, blackMs, activeColor }: GameClockProps) {
  const whiteActive = activeColor === 'white'
  const blackActive = activeColor === 'black'

  return (
    <div className="glass-card flex items-center justify-between p-3" style={{ maxWidth: 320 }}>
      {/* Black clock */}
      <div className="flex items-center gap-2 px-3 py-2 rounded" style={{
        opacity: blackActive ? 1 : 0.4,
        background: blackActive ? 'rgba(212, 168, 67, 0.1)' : 'transparent',
        border: blackActive ? '1px solid rgba(212, 168, 67, 0.3)' : '1px solid transparent',
        minWidth: 90,
        justifyContent: 'center',
      }}>
        <div style={{
          width: 10,
          height: 10,
          borderRadius: '50%',
          background: '#1a1a2e',
          border: '1px solid #3a3a5a',
          flexShrink: 0,
        }} />
        {blackActive && blackMs < 30000 && (
          <motion.span
            animate={{ opacity: [1, 0.5, 1] }}
            transition={{ repeat: Infinity, duration: 1 }}
            className="font-mono text-base font-bold"
            style={{ color: 'var(--danger)' }}
          >
            {formatTime(blackMs)}
          </motion.span>
        )}
        {(!blackActive || blackMs >= 30000) && (
          <span className="font-mono text-base font-bold">{formatTime(blackMs)}</span>
        )}
      </div>

      <span className="text-xs text-muted">VS</span>

      {/* White clock */}
      <div className="flex items-center gap-2 px-3 py-2 rounded" style={{
        opacity: whiteActive ? 1 : 0.4,
        background: whiteActive ? 'rgba(212, 168, 67, 0.1)' : 'transparent',
        border: whiteActive ? '1px solid rgba(212, 168, 67, 0.3)' : '1px solid transparent',
        minWidth: 90,
        justifyContent: 'center',
      }}>
        <div style={{
          width: 10,
          height: 10,
          borderRadius: '50%',
          background: '#fff',
          border: '1px solid #c9b896',
          flexShrink: 0,
        }} />
        {whiteActive && whiteMs < 30000 && (
          <motion.span
            animate={{ opacity: [1, 0.5, 1] }}
            transition={{ repeat: Infinity, duration: 1 }}
            className="font-mono text-base font-bold"
            style={{ color: 'var(--danger)' }}
          >
            {formatTime(whiteMs)}
          </motion.span>
        )}
        {(!whiteActive || whiteMs >= 30000) && (
          <span className="font-mono text-base font-bold">{formatTime(whiteMs)}</span>
        )}
      </div>
    </div>
  )
}
