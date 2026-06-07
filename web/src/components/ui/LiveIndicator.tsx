import { motion } from 'framer-motion'

export function LiveIndicator() {
  return (
    <div className="flex items-center gap-1.5">
      <motion.div
        animate={{ scale: [1, 1.3, 1], opacity: [1, 0.7, 1] }}
        transition={{ repeat: Infinity, duration: 2 }}
        style={{
          width: 8,
          height: 8,
          borderRadius: '50%',
          background: 'var(--success)',
          boxShadow: '0 0 8px rgba(74,222,128,0.5)',
        }}
      />
      <span className="text-xs font-semibold" style={{ color: 'var(--success)' }}>LIVE</span>
    </div>
  )
}
