import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { X, Search, Send, Loader2, Gamepad2, Swords, Crown } from 'lucide-react'
import { useInviteStore } from '../../stores/inviteStore'
import { useToastStore } from '../ui/Toast'

// ─── Game type options ─────────────────────────────────────────────────────

const GAME_TYPES = [
  { value: 'chess', label: 'Шахматы', Icon: Crown },
  { value: 'checkers', label: 'Шашки', Icon: Gamepad2 },
  { value: 'backgammon', label: 'Нарды', Icon: Swords },
] as const

// ─── Props ────────────────────────────────────────────────────────────────

interface InvitePlayerModalProps {
  isOpen: boolean
  onClose: () => void
  gameType?: string
}

// ─── Component ──────────────────────────────────────────────────────────────

export function InvitePlayerModal({ isOpen, onClose, gameType }: InvitePlayerModalProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedGameType, setSelectedGameType] = useState(gameType || 'chess')
  const { sending, create } = useInviteStore()
  const addToast = useToastStore((s) => s.addToast)

  const handleSendInvite = async () => {
    const sid = searchQuery.trim()
    if (!sid) {
      addToast('Введите SID или имя коллеги', 'error')
      return
    }

    try {
      await create(selectedGameType, sid)
      addToast('Приглашение отправлено', 'success')
      setSearchQuery('')
      onClose()
    } catch {
      // Error toast is handled by the API layer via useToastStore
    }
  }

  const handleClose = () => {
    setSearchQuery('')
    onClose()
  }

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="modal-backdrop"
            onClick={handleClose}
          />

          {/* Modal */}
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            transition={{ type: 'spring', damping: 25, stiffness: 300 }}
            className="modal-overlay"
          >
            <div className="glass-card p-6 w-full max-w-md mx-4">
              {/* Header */}
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-lg font-bold">Пригласить игрока</h2>
                <button className="btn-icon" onClick={handleClose}>
                  <X size={18} />
                </button>
              </div>

              {/* Game type selector */}
              <div className="mb-5">
                <label className="block text-sm font-medium text-secondary mb-2">
                  Тип игры
                </label>
                <div className="flex gap-2">
                  {GAME_TYPES.map(({ value, label, Icon }) => (
                    <button
                      key={value}
                      className={`btn flex-1 justify-center ${
                        selectedGameType === value ? 'btn-primary' : 'btn-secondary'
                      }`}
                      onClick={() => setSelectedGameType(value)}
                    >
                      <Icon size={14} />
                      {label}
                    </button>
                  ))}
                </div>
              </div>

              {/* Search input */}
              <div className="mb-5">
                <label className="block text-sm font-medium text-secondary mb-2">
                  Игрок
                </label>
                <div className="glass-card flex items-center gap-2 px-3 py-2">
                  <Search size={16} color="var(--text-muted)" />
                  <input
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="SID или имя коллеги"
                    className="flex-1 bg-transparent border-none outline-none text-sm"
                    style={{ color: 'var(--text-primary)' }}
                    onKeyDown={(e) => e.key === 'Enter' && handleSendInvite()}
                    autoFocus
                  />
                </div>
                <p className="text-xs text-muted mt-1">
                  Введите корпоративный SID коллеги, которого хотите пригласить
                </p>
              </div>

              {/* Submit */}
              <button
                className="btn btn-primary w-full justify-center"
                onClick={handleSendInvite}
                disabled={sending || !searchQuery.trim()}
              >
                {sending ? (
                  <Loader2 size={16} className="spin" />
                ) : (
                  <Send size={16} />
                )}
                {sending ? 'Отправка...' : 'Отправить приглашение'}
              </button>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}
