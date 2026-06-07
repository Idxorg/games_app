import { motion } from 'framer-motion'

interface EloBarProps {
  value: number
  max: number
}

function getEloColor(value: number): string {
  if (value >= 2000) return '#4ade80'
  if (value >= 1500) return '#fbbf24'
  return '#f87171'
}

function getEloGradient(value: number): string {
  if (value >= 2000) return 'linear-gradient(90deg, #22c55e, #4ade80)'
  if (value >= 1500) return 'linear-gradient(90deg, #f59e0b, #fbbf24)'
  return 'linear-gradient(90deg, #ef4444, #f87171)'
}

export function EloBar({ value, max }: EloBarProps) {
  const percentage = Math.min((value / max) * 100, 100)
  const color = getEloColor(value)

  return (
    <div className="flex items-center gap-2">
      <div className="flex-grow h-2 rounded-full overflow-hidden" style={{ background: 'var(--bg-tertiary)' }}>
        <motion.div
          initial={{ width: 0 }}
          animate={{ width: `${percentage}%` }}
          transition={{ duration: 1, ease: 'easeOut', delay: 0.3 }}
          className="h-full rounded-full"
          style={{ background: getEloGradient(value) }}
        />
      </div>
      <span className="text-xs font-mono font-bold" style={{ color, minWidth: '36px', textAlign: 'right' }}>
        {value}
      </span>
    </div>
  )
}
