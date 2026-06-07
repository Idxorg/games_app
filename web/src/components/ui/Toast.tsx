import { motion, AnimatePresence } from 'framer-motion'
import { X, CheckCircle, AlertCircle, Info } from 'lucide-react'
import { create } from 'zustand'

interface ToastItem {
  id: string
  message: string
  type: 'success' | 'error' | 'info'
}

interface ToastStore {
  toasts: ToastItem[]
  addToast: (message: string, type?: ToastItem['type']) => void
  removeToast: (id: string) => void
}

export const useToastStore = create<ToastStore>((set) => ({
  toasts: [],
  addToast: (message, type = 'info') => {
    const id = Math.random().toString(36).slice(2)
    set((s) => ({ toasts: [...s.toasts, { id, message, type }] }))
    setTimeout(() => {
      set((s) => ({ toasts: s.toasts.filter((t) => t.id !== id) }))
    }, 4000)
  },
  removeToast: (id) => set((s) => ({ toasts: s.toasts.filter((t) => t.id !== id) })),
}))

const icons = {
  success: <CheckCircle size={18} color="var(--success)" />,
  error: <AlertCircle size={18} color="var(--danger)" />,
  info: <Info size={18} color="var(--info)" />,
}

const bgColors = {
  success: 'rgba(74,222,128,0.1)',
  error: 'rgba(248,113,113,0.1)',
  info: 'rgba(96,165,250,0.1)',
}

const borderColors = {
  success: 'rgba(74,222,128,0.2)',
  error: 'rgba(248,113,113,0.2)',
  info: 'rgba(96,165,250,0.2)',
}

export function ToastContainer() {
  const { toasts, removeToast } = useToastStore()

  return (
    <div style={{ position: 'fixed', top: 80, right: 24, zIndex: 1000, display: 'flex', flexDirection: 'column', gap: 8 }}>
      <AnimatePresence>
        {toasts.map((toast) => (
          <motion.div
            key={toast.id}
            initial={{ opacity: 0, x: 100, scale: 0.9 }}
            animate={{ opacity: 1, x: 0, scale: 1 }}
            exit={{ opacity: 0, x: 100, scale: 0.9 }}
            transition={{ type: 'spring', damping: 25, stiffness: 300 }}
            className="glass-card flex items-center gap-3 px-4 py-3 cursor-pointer"
            style={{
              background: bgColors[toast.type],
              borderColor: borderColors[toast.type],
              minWidth: 280,
              maxWidth: 400,
            }}
            onClick={() => removeToast(toast.id)}
          >
            {icons[toast.type]}
            <span className="text-sm flex-grow" style={{ color: 'var(--text-primary)' }}>{toast.message}</span>
            <X size={14} color="var(--text-muted)" />
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  )
}
