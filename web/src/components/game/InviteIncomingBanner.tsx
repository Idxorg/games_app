import { useEffect, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Bell, Swords, X, Check, Loader2 } from 'lucide-react'
import { useInviteStore } from '../../stores/inviteStore'
import { useToastStore } from '../ui/Toast'
import { isEmbedMode } from '../../embedHandoff'

// ─── Game type labels ─────────────────────────────────────────────────────

const GAME_TYPE_LABELS: Record<string, string> = {
  chess: 'шахматы',
  checkers: 'шашки',
  backgammon: 'нарды',
  'tic-tac-toe': 'крестики-нолики',
  trivia: 'викторину',
  memory: 'мемори',
  puzzle: 'головоломку',
}

// ─── Component ──────────────────────────────────────────────────────────────

export function InviteIncomingBanner() {
  const { pendingInvites, loading, accept, decline } = useInviteStore()
  const addToast = useToastStore((s) => s.addToast)
  const [actionInProgress, setActionInProgress] = useState<string | null>(null)

  // Auto-fetch on mount if authenticated
  useEffect(() => {
    useInviteStore.getState().fetchPending()
  }, [])

  // Poll every 30s for new invites
  useEffect(() => {
    const interval = setInterval(() => {
      useInviteStore.getState().fetchPending()
    }, 30_000)
    return () => clearInterval(interval)
  }, [])

  if (loading || pendingInvites.length === 0) return null

  const handleAccept = async (id: string) => {
    setActionInProgress(id)
    try {
      await accept(id)
      addToast('Приглашение принято', 'success')
      // Notify parent iframe
      if (isEmbedMode()) {
        window.parent.postMessage(
          { type: 'erlink_games_invite_accept', invite_id: id },
          '*',
        )
      }
    } catch {
      // Error toast already shown by API layer
    } finally {
      setActionInProgress(null)
    }
  }

  const handleDecline = async (id: string) => {
    setActionInProgress(id)
    try {
      await decline(id)
      // Notify parent iframe
      if (isEmbedMode()) {
        window.parent.postMessage(
          { type: 'erlink_games_invite_decline', invite_id: id },
          '*',
        )
      }
    } catch {
      // Error toast already shown by API layer
    } finally {
      setActionInProgress(null)
    }
  }

  // Show only the first invite (most recent)
  const invite = pendingInvites[0]

  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: -20 }}
        transition={{ type: 'spring', damping: 25, stiffness: 300 }}
        className="glass-card gold-border flex items-center gap-3 px-4 py-3 mx-auto mb-4"
        style={{ maxWidth: 600 }}
      >
        <Bell size={18} color="var(--gold)" className="flex-shrink-0" />

        <div className="flex-1 min-w-0">
          <span className="text-sm font-medium truncate block">
            {invite.inviter_name} хочет сыграть в{' '}
            <span className="text-accent">
              {GAME_TYPE_LABELS[invite.game_type] || invite.game_type}
            </span>
          </span>
        </div>

        <div className="flex items-center gap-2 flex-shrink-0">
          <button
            className="btn btn-primary btn-sm"
            onClick={() => handleAccept(invite.id)}
            disabled={actionInProgress === invite.id}
          >
            {actionInProgress === invite.id ? (
              <Loader2 size={14} className="spin" />
            ) : (
              <Check size={14} />
            )}
            Принять
          </button>

          <button
            className="btn btn-ghost btn-sm"
            onClick={() => handleDecline(invite.id)}
            disabled={actionInProgress === invite.id}
          >
            {actionInProgress === invite.id ? (
              <Loader2 size={14} className="spin" />
            ) : (
              <X size={14} />
            )}
            Отклонить
          </button>
        </div>

        {pendingInvites.length > 1 && (
          <span className="text-xs text-muted flex-shrink-0">
            +{pendingInvites.length - 1}
          </span>
        )}
      </motion.div>
    </AnimatePresence>
  )
}
