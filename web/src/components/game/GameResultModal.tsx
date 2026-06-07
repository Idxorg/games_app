import { motion, AnimatePresence } from 'framer-motion'
import { X, Trophy } from 'lucide-react'

// ─── Types ──────────────────────────────────────────────────────────────────

interface GameResult {
  reason: string
  winner_sid: string
  winner_name: string
  score: string
  my_sid: string
}

interface GameResultModalProps {
  isOpen: boolean
  onClose: () => void
  result: GameResult | null
}

// ─── Helpers ────────────────────────────────────────────────────────────────

const REASON_LABELS: Record<string, string> = {
  checkmate: 'Мат',
  resign: 'Сдача',
  draw: 'Ничья',
  timeout: 'Время вышло',
  stalemate: 'Пат',
}

// ─── Component ──────────────────────────────────────────────────────────────

export function GameResultModal({ isOpen, onClose, result }: GameResultModalProps) {
  if (!result) return null

  const isWinner = result.winner_sid === result.my_sid
  const isDraw = result.reason === 'draw' || result.reason === 'stalemate'

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
            <div className="glass-card p-6 w-full max-w-sm mx-4 text-center">
              <div className="flex items-center justify-between mb-4">
                <div />
                <button className="btn-icon" onClick={onClose}>
                  <X size={18} />
                </button>
              </div>

              <div style={{ marginBottom: 16 }}>
                {isWinner && (
                  <Trophy size={40} color="var(--gold)" style={{ marginBottom: 8 }} />
                )}
                <h2 className="text-xl font-bold">
                  {isDraw
                    ? 'Ничья'
                    : isWinner
                      ? 'Победа'
                      : 'Поражение'}
                </h2>
              </div>

              <div className="glass-card p-4 mb-4" style={{ cursor: 'default' }}>
                {isDraw ? (
                  <span className="text-sm text-secondary">
                    Причина: {REASON_LABELS[result.reason] || result.reason}
                  </span>
                ) : (
                  <span className="text-sm text-secondary">
                    {isWinner ? 'Вы' : result.winner_name || 'Оппонент'} победили
                  </span>
                )}
              </div>

              <div className="flex items-center justify-center gap-3 mb-6">
                <span className="text-xs text-muted">Счёт:</span>
                <span className="font-mono font-bold text-accent">
                  {result.score || '-'}
                </span>
              </div>

              <button className="btn btn-primary w-full justify-center" onClick={onClose}>
                Закрыть
              </button>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}
