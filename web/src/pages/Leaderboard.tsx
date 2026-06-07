import { motion } from 'framer-motion'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { LeaderboardRow } from '../components/ui/LeaderboardRow'
import { useLeaderboardStore } from '../stores/leaderboardStore'

const gameOptions = [
  { value: 'all', label: 'Все игры' },
  { value: 'chess', label: 'Шахматы' },
  { value: 'checkers', label: 'Шашки' },
  { value: 'backgammon', label: 'Нарды' },
  { value: 'trivia', label: 'Викторины' },
]

const periodOptions = [
  { value: 'all', label: 'Всё время' },
  { value: 'week', label: 'За неделю' },
  { value: 'month', label: 'За месяц' },
  { value: 'year', label: 'За год' },
]

export function Leaderboard() {
  const store = useLeaderboardStore()
  const paged = store.getPaged()
  const totalPages = store.getTotalPages()

  return (
    <div className="py-8">
      <div className="container">
        <h1 className="text-3xl font-bold mb-2">Рейтинг</h1>
        <p className="mb-6" style={{ color: 'var(--text-secondary)' }}>
          Таблица лидеров с рейтингом ELO по всем играм
        </p>

        {/* Filters */}
        <div className="flex flex-wrap items-center gap-3 mb-6">
          <div className="flex flex-wrap gap-2">
            {gameOptions.map((opt) => (
              <button
                key={opt.value}
                onClick={() => store.setGameFilter(opt.value as any)}
                className={`tab ${store.gameFilter === opt.value ? 'active' : ''}`}
                style={store.gameFilter !== opt.value ? {
                  background: 'var(--bg-glass)',
                  border: '1px solid var(--bg-glass-border)',
                  borderRadius: 'var(--radius-sm)',
                  padding: '8px 16px',
                  fontSize: '0.875rem',
                  fontWeight: 500,
                  color: 'var(--text-secondary)',
                  cursor: 'pointer',
                  fontFamily: 'var(--font-sans)',
                } : undefined}
              >
                {opt.label}
              </button>
            ))}
          </div>

          <div style={{ borderLeft: '1px solid var(--bg-glass-border)', height: 24, margin: '0 8px' }} />

          <div className="flex gap-2">
            {periodOptions.map((opt) => (
              <button
                key={opt.value}
                onClick={() => store.setPeriodFilter(opt.value as any)}
                className={`tab ${store.periodFilter === opt.value ? 'active' : ''}`}
                style={store.periodFilter !== opt.value ? {
                  background: 'var(--bg-glass)',
                  border: '1px solid var(--bg-glass-border)',
                  borderRadius: 'var(--radius-sm)',
                  padding: '8px 16px',
                  fontSize: '0.875rem',
                  fontWeight: 500,
                  color: 'var(--text-secondary)',
                  cursor: 'pointer',
                  fontFamily: 'var(--font-sans)',
                } : undefined}
              >
                {opt.label}
              </button>
            ))}
          </div>
        </div>

        {/* Table Header */}
        <div className="hidden md:flex items-center gap-4 px-4 py-2 mb-2 text-xs font-semibold" style={{ color: 'var(--text-muted)' }}>
          <div className="w-10">Место</div>
          <div className="w-10" />
          <div className="flex-grow">Игрок</div>
          {store.gameFilter !== 'all' && <div className="w-32">Рейтинг</div>}
          {store.gameFilter !== 'all' && <div className="w-16 text-right">ELO</div>}
          <div className="w-8" />
          <div className="flex gap-3">В / П / Н</div>
        </div>

        {/* Rows */}
        {paged.map((player, i) => (
          <LeaderboardRow
            key={player.sid}
            player={player}
            rank={(store.page - 1) * store.perPage + i + 1}
            gameFilter={store.gameFilter}
          />
        ))}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-center gap-3 mt-6">
            <button
              onClick={() => store.setPage(Math.max(1, store.page - 1))}
              disabled={store.page === 1}
              className="btn btn-secondary"
              style={{ opacity: store.page === 1 ? 0.4 : 1 }}
            >
              <ChevronLeft size={16} />
            </button>
            <span className="text-sm" style={{ color: 'var(--text-muted)' }}>
              {store.page} из {totalPages}
            </span>
            <button
              onClick={() => store.setPage(Math.min(totalPages, store.page + 1))}
              disabled={store.page === totalPages}
              className="btn btn-secondary"
              style={{ opacity: store.page === totalPages ? 0.4 : 1 }}
            >
              <ChevronRight size={16} />
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
