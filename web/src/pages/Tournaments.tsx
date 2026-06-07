import { motion } from 'framer-motion'
import { Trophy, Users, Calendar, Gift, ChevronDown, ChevronUp } from 'lucide-react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { TournamentBracket } from '../components/ui/TournamentBracket'
import { Countdown } from '../components/ui/Countdown'
import { LiveIndicator } from '../components/ui/LiveIndicator'
import { useTournamentStore } from '../stores/tournamentStore'
import { players } from '../data/players'

const statusLabels = {
  active: { text: 'Активный', badgeClass: 'badge-success', icon: LiveIndicator },
  upcoming: { text: 'Предстоящий', badgeClass: 'badge-gold', icon: Calendar },
  completed: { text: 'Завершён', badgeClass: 'badge-neutral', icon: Trophy },
}

function TournamentCard({ tournamentId, index }: { tournamentId: string; index: number }) {
  const { tournaments } = useTournamentStore()
  const tournament = tournaments.find((t) => t.id === tournamentId)
  const [expanded, setExpanded] = useState(false)
  const navigate = useNavigate()

  if (!tournament) return null

  const status = statusLabels[tournament.status]

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.1 }}
      className="glass-card overflow-hidden cursor-pointer"
      onClick={() => navigate(`/tournaments/${tournament.id}`)}
    >
      {/* Header */}
      <div className="p-6">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-4">
          <div className="flex items-center gap-3">
            <Trophy size={24} color="var(--gold)" />
            <div>
              <h3 className="text-lg font-bold">{tournament.name}</h3>
              <p className="text-sm text-muted">
                Шахматы — {tournament.startDate} — {tournament.endDate}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <span className={`badge ${status.badgeClass}`}>{status.text}</span>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-6 mb-4">
          <div className="flex items-center gap-2 text-sm text-secondary">
            <Users size={16} />
            <span>{tournament.participants.length} участников</span>
          </div>
          <div className="flex items-center gap-2 text-sm text-secondary">
            <Gift size={16} />
            <span>{tournament.prize}</span>
          </div>
        </div>

        {/* Participants */}
        {tournament.participants.length > 0 && (
          <div className="flex items-center gap-2 mb-4">
            <div className="avatar-stack">
              {tournament.participants.slice(0, 6).map((pid) => {
                const p = players.find((pl) => pl.sid === pid)
                return (
                  <div
                    key={pid}
                    className="avatar-sm"
                    title={p?.name}
                  >
                    {p?.initials}
                  </div>
                )
              })}
              {tournament.participants.length > 6 && (
                <div className="avatar-sm avatar-sm-neutral">
                  +{tournament.participants.length - 6}
                </div>
              )}
            </div>
          </div>
        )}

        {/* Countdown for upcoming */}
        {tournament.status === 'upcoming' && (
          <div className="mb-4">
            <Countdown targetDate={tournament.startDate} label="До начала турнира:" />
          </div>
        )}

        {/* Actions */}
        <div className="flex items-center gap-3" onClick={(e) => e.stopPropagation()}>
          {tournament.status === 'upcoming' && (
            <button className="btn btn-primary">Зарегистрироваться</button>
          )}
          {tournament.status === 'active' && (
            <button className="btn btn-secondary">Наблюдать</button>
          )}
          {tournament.rounds.length > 0 && (
            <button
              onClick={() => setExpanded(!expanded)}
              className="btn btn-ghost"
            >
              {expanded ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
              {expanded ? 'Скрыть сетку' : 'Показать сетку'}
            </button>
          )}
        </div>
      </div>

      {/* Bracket */}
      {expanded && tournament.rounds.length > 0 && (
        <motion.div
          initial={{ height: 0, opacity: 0 }}
          animate={{ height: 'auto', opacity: 1 }}
          transition={{ duration: 0.3 }}
          className="p-6 divider"
        >
          <h4 className="text-sm font-semibold mb-4 text-secondary">Турнирная сетка</h4>
          <TournamentBracket rounds={tournament.rounds} />
        </motion.div>
      )}
    </motion.div>
  )
}

function TournamentSkeleton() {
  return (
    <div className="glass-card p-6">
      <div className="flex items-center gap-3 mb-4">
        <div className="skeleton skeleton-circle" style={{ width: 24, height: 24 }} />
        <div className="skeleton skeleton-title" style={{ width: '40%' }} />
      </div>
      <div className="flex gap-6 mb-4">
        <div className="skeleton skeleton-text" style={{ width: 120 }} />
        <div className="skeleton skeleton-text" style={{ width: 80 }} />
      </div>
      <div className="flex gap-3">
        <div className="skeleton" style={{ width: 160, height: 36, borderRadius: 'var(--radius-sm)' }} />
      </div>
    </div>
  )
}

export function Tournaments() {
  const { activeTab, setActiveTab, tournaments } = useTournamentStore()
  const filtered = tournaments.filter((t) => t.status === activeTab)

  return (
    <div className="py-8">
      <div className="container">
        <h1 className="text-3xl font-bold mb-2">Турниры</h1>
        <p className="mb-6 text-secondary">
          Участвуйте в корпоративных турнирах и выигрывайте призы
        </p>

        {/* Tabs */}
        <div className="tabs mb-6 inline-flex">
          {(['active', 'upcoming', 'completed'] as const).map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`tab ${activeTab === tab ? 'active' : ''}`}
            >
              {statusLabels[tab].text}
            </button>
          ))}
        </div>

        {/* Tournament List */}
        <div className="flex flex-col gap-4">
          {filtered.length > 0 ? (
            filtered.map((tournament, i) => (
              <TournamentCard key={tournament.id} tournamentId={tournament.id} index={i} />
            ))
          ) : (
            <div className="glass-card p-12 text-center">
              <Trophy size={48} color="var(--text-muted)" className="mx-auto mb-4" />
              <p className="text-muted">Нет турниров в этой категории</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
