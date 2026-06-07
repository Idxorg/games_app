import { useState, useEffect, useRef, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { X, Search, Send, Loader2, Gamepad2, Swords, Crown, User, ChevronDown } from 'lucide-react'
import { useInviteStore } from '../../stores/inviteStore'
import { useToastStore } from '../ui/Toast'
import { directorySearch, type DirectoryEntry } from '../../api/client'

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
  const [manualSid, setManualSid] = useState('')
  const [showManualInput, setShowManualInput] = useState(false)

  // Directory search state
  const [directoryResults, setDirectoryResults] = useState<DirectoryEntry[]>([])
  const [searching, setSearching] = useState(false)
  const [selectedEntry, setSelectedEntry] = useState<DirectoryEntry | null>(null)
  const [showDropdown, setShowDropdown] = useState(false)

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const { sending, create } = useInviteStore()
  const addToast = useToastStore((s) => s.addToast)

  // Debounced directory search (300ms)
  const doSearch = useCallback(async (query: string) => {
    if (!query.trim()) {
      setDirectoryResults([])
      setShowDropdown(false)
      return
    }
    setSearching(true)
    try {
      const results = await directorySearch(query)
      setDirectoryResults(results)
      setShowDropdown(results.length > 0)
    } catch {
      setDirectoryResults([])
      setShowDropdown(false)
    } finally {
      setSearching(false)
    }
  }, [])

  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current)
    if (!showManualInput) {
      debounceRef.current = setTimeout(() => doSearch(searchQuery), 300)
    }
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current)
    }
  }, [searchQuery, showManualInput, doSearch])

  // Close dropdown on outside click
  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setShowDropdown(false)
      }
    }
    if (showDropdown) {
      document.addEventListener('mousedown', handleClick)
      return () => document.removeEventListener('mousedown', handleClick)
    }
  }, [showDropdown])

  // Reset on open/close
  useEffect(() => {
    if (isOpen) {
      setSearchQuery('')
      setManualSid('')
      setSelectedEntry(null)
      setDirectoryResults([])
      setShowDropdown(false)
      setSearching(false)
      setShowManualInput(false)
    }
  }, [isOpen])

  const handleSelectEntry = (entry: DirectoryEntry) => {
    setSelectedEntry(entry)
    setSearchQuery(entry.name)
    setShowDropdown(false)
  }

  const handleClearSelection = () => {
    setSelectedEntry(null)
    setSearchQuery('')
  }

  const handleSendInvite = async () => {
    const recipientSid = selectedEntry ? selectedEntry.sid : (showManualInput ? manualSid.trim() : '')
    if (!recipientSid) {
      addToast('Введите SID или выберите коллегу из каталога', 'error')
      return
    }

    try {
      await create(selectedGameType, recipientSid)
      addToast('Приглашение отправлено', 'success')
      onClose()
    } catch {
      // Error toast is handled by the API layer via useToastStore
    }
  }

  const handleClose = () => {
    onClose()
  }

  const canSend = selectedEntry
    ? true
    : showManualInput
      ? manualSid.trim().length > 0
      : false

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

              {/* Directory search / player input */}
              <div className="mb-5">
                <label className="block text-sm font-medium text-secondary mb-2">
                  Игрок
                </label>

                {!selectedEntry ? (
                  <>
                    {/* Search input with dropdown */}
                    <div className="relative" ref={dropdownRef}>
                      <div className="glass-card flex items-center gap-2 px-3 py-2">
                        <Search size={16} color="var(--text-muted)" />
                        <input
                          type="text"
                          value={showManualInput ? '' : searchQuery}
                          onChange={(e) => {
                            setSearchQuery(e.target.value)
                            setSelectedEntry(null)
                          }}
                          placeholder={showManualInput ? '' : 'Поиск по имени в каталоге...'}
                          className="flex-1 bg-transparent border-none outline-none text-sm"
                          style={{ color: 'var(--text-primary)' }}
                          disabled={showManualInput}
                          autoFocus
                        />
                        {searching && <Loader2 size={14} className="spin text-muted" />}
                      </div>

                      {/* Dropdown results */}
                      {showDropdown && directoryResults.length > 0 && (
                        <div className="absolute z-50 mt-1 w-full glass-card border border-border/50 rounded-lg overflow-hidden shadow-xl">
                          {directoryResults.map((entry) => (
                            <button
                              key={entry.sid}
                              className="w-full flex items-center gap-3 px-3 py-2 text-left text-sm hover:bg-white/5 transition-colors"
                              style={{ color: 'var(--text-primary)' }}
                              onClick={() => handleSelectEntry(entry)}
                            >
                              {entry.photo_url ? (
                                <img
                                  src={entry.photo_url}
                                  alt={entry.name}
                                  className="w-8 h-8 rounded-full object-cover"
                                />
                              ) : (
                                <div
                                  className="w-8 h-8 rounded-full flex items-center justify-center text-xs font-semibold"
                                  style={{ background: 'var(--primary)', color: '#fff' }}
                                >
                                  {entry.name.charAt(0).toUpperCase()}
                                </div>
                              )}
                              <div className="flex-1 min-w-0">
                                <div className="font-medium truncate">{entry.name}</div>
                                <div className="text-xs text-muted truncate">
                                  {entry.department}
                                </div>
                              </div>
                              <span className="text-xs text-muted">{entry.sid}</span>
                            </button>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* Manual SID toggle */}
                    <div className="flex items-center justify-between mt-2">
                      <button
                        className="text-xs text-muted hover:text-primary transition-colors"
                        onClick={() => setShowManualInput(!showManualInput)}
                      >
                        {showManualInput
                          ? '← Поиск в каталоге'
                          : 'Ввести SID вручную'}
                      </button>
                    </div>

                    {/* Manual SID input */}
                    {showManualInput && (
                      <div className="glass-card flex items-center gap-2 px-3 py-2 mt-2">
                        <User size={16} color="var(--text-muted)" />
                        <input
                          type="text"
                          value={manualSid}
                          onChange={(e) => setManualSid(e.target.value)}
                          placeholder="corp-username"
                          className="flex-1 bg-transparent border-none outline-none text-sm"
                          style={{ color: 'var(--text-primary)' }}
                          autoFocus
                          onKeyDown={(e) => e.key === 'Enter' && handleSendInvite()}
                        />
                      </div>
                    )}

                    {!showManualInput && (
                      <p className="text-xs text-muted mt-1">
                        Введите имя для поиска или введите SID вручную
                      </p>
                    )}
                  </>
                ) : (
                  /* Selected entry card */
                  <div className="glass-card flex items-center gap-3 px-3 py-2">
                    {selectedEntry.photo_url ? (
                      <img
                        src={selectedEntry.photo_url}
                        alt={selectedEntry.name}
                        className="w-8 h-8 rounded-full object-cover"
                      />
                    ) : (
                      <div
                        className="w-8 h-8 rounded-full flex items-center justify-center text-xs font-semibold"
                        style={{ background: 'var(--primary)', color: '#fff' }}
                      >
                        {selectedEntry.name.charAt(0).toUpperCase()}
                      </div>
                    )}
                    <div className="flex-1 min-w-0">
                      <div className="font-medium text-sm truncate">{selectedEntry.name}</div>
                      <div className="text-xs text-muted truncate">
                        {selectedEntry.department} · {selectedEntry.sid}
                      </div>
                    </div>
                    <button
                      className="btn-icon"
                      onClick={handleClearSelection}
                      title="Убрать"
                    >
                      <X size={14} />
                    </button>
                  </div>
                )}
              </div>

              {/* Submit */}
              <button
                className="btn btn-primary w-full justify-center"
                onClick={handleSendInvite}
                disabled={sending || !canSend}
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
