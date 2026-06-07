import { useEffect, useState } from 'react'
import { motion } from 'framer-motion'

interface StatsCardProps {
  label: string
  value: number
  suffix?: string
  icon: React.ReactNode
  color?: string
}

export function StatsCard({ label, value, suffix = '', icon, color = 'var(--gold)' }: StatsCardProps) {
  const [display, setDisplay] = useState(0)

  useEffect(() => {
    const duration = 1500
    const steps = 60
    const increment = value / steps
    let current = 0
    const timer = setInterval(() => {
      current += increment
      if (current >= value) {
        setDisplay(value)
        clearInterval(timer)
      } else {
        setDisplay(Math.floor(current))
      }
    }, duration / steps)
    return () => clearInterval(timer)
  }, [value])

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="glass-card p-6 flex flex-col items-center text-center"
    >
      <div className="mb-3" style={{ color }}>{icon}</div>
      <div className="text-3xl font-bold font-mono" style={{ color: 'var(--text-primary)' }}>
        {display.toLocaleString('ru-RU')}{suffix}
      </div>
      <div className="text-xs mt-1" style={{ color: 'var(--text-muted)' }}>{label}</div>
    </motion.div>
  )
}
