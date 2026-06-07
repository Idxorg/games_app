import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'

interface CountdownProps {
  targetDate: string
  label?: string
}

interface TimeLeft {
  days: number
  hours: number
  minutes: number
  seconds: number
}

function calcTimeLeft(target: string): TimeLeft {
  const diff = new Date(target).getTime() - Date.now()
  if (diff <= 0) return { days: 0, hours: 0, minutes: 0, seconds: 0 }
  return {
    days: Math.floor(diff / (1000 * 60 * 60 * 24)),
    hours: Math.floor((diff / (1000 * 60 * 60)) % 24),
    minutes: Math.floor((diff / (1000 * 60)) % 60),
    seconds: Math.floor((diff / 1000) % 60),
  }
}

function FlipDigit({ value, label }: { value: number; label: string }) {
  const padded = String(value).padStart(2, '0')
  return (
    <div className="flex flex-col items-center gap-1">
      <div
        className="flex items-center justify-center font-mono font-bold text-2xl rounded-lg"
        style={{
          width: 56,
          height: 56,
          background: 'var(--bg-glass)',
          border: '1px solid var(--bg-glass-border)',
          color: 'var(--text-primary)',
          perspective: '200px',
        }}
      >
        <motion.span
          key={padded}
          initial={{ rotateX: -90, opacity: 0 }}
          animate={{ rotateX: 0, opacity: 1 }}
          transition={{ duration: 0.3 }}
        >
          {padded}
        </motion.span>
      </div>
      <span className="text-xs font-medium" style={{ color: 'var(--text-muted)' }}>{label}</span>
    </div>
  )
}

export function Countdown({ targetDate, label }: CountdownProps) {
  const [time, setTime] = useState(calcTimeLeft(targetDate))

  useEffect(() => {
    const timer = setInterval(() => setTime(calcTimeLeft(targetDate)), 1000)
    return () => clearInterval(timer)
  }, [targetDate])

  const isExpired = time.days === 0 && time.hours === 0 && time.minutes === 0 && time.seconds === 0

  if (isExpired) {
    return (
      <div className="text-sm font-medium" style={{ color: 'var(--success)' }}>
        Турнир завершён
      </div>
    )
  }

  return (
    <div>
      {label && <div className="text-sm mb-2" style={{ color: 'var(--text-secondary)' }}>{label}</div>}
      <div className="flex gap-3">
        <FlipDigit value={time.days} label="дней" />
        <div className="flex items-center text-xl font-bold" style={{ color: 'var(--text-muted)', marginTop: '-16px' }}>:</div>
        <FlipDigit value={time.hours} label="часов" />
        <div className="flex items-center text-xl font-bold" style={{ color: 'var(--text-muted)', marginTop: '-16px' }}>:</div>
        <FlipDigit value={time.minutes} label="минут" />
        <div className="flex items-center text-xl font-bold" style={{ color: 'var(--text-muted)', marginTop: '-16px' }}>:</div>
        <FlipDigit value={time.seconds} label="секунд" />
      </div>
    </div>
  )
}
