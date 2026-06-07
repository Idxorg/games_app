import { motion, AnimatePresence } from 'framer-motion'
import { X } from 'lucide-react'

// ─── Types ──────────────────────────────────────────────────────────────────

interface ResignDialogProps {
  isOpen: boolean
  onClose: () => void
  onResign: () => void
}

// ─── Component ──────────────────────────────────────────────────────────────

export function ResignDialog({ isOpen, onClose, onResign }: ResignDialogProps) {
  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="modal-backdrop"
            onClick={onClose}
          />

          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            transition={{ type: 'spring', damping: 25, stiffness: 300 }}
            className="modal-overlay"
          >
            <div className="glass-card p-6 w-full max-w-sm mx-4">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-bold">Сдаться</h2>
                <button className="btn-icon" onClick={onClose}>
                  <X size={18} />
                </button>
              </div>

              <p className="text-sm text-secondary mb-6">
                Вы уверены что хотите сдаться?
              </p>

              <div className="flex gap-3">
                <button className="btn btn-secondary flex-1 justify-center" onClick={onClose}>
                  Отмена
                </button>
                <button className="btn flex-1 justify-center" onClick={onResign} style={{
                  background: 'var(--danger)',
                  color: '#fff',
                }}>
                  Сдаться
                </button>
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}
