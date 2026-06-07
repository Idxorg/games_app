import { useEffect } from 'react'
import { motion } from 'framer-motion'
import { ChevronLeft, ChevronRight, RefreshCw, Trophy } from 'lucide-react'
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

function LeaderboardSkeleton() {
  return (
    <div className="flex flex-col gap-2">
      {Array.from({ length: 10 }).map((_, i) => (
        <div key={i} className="glass-card p-4 flex items-center gap-4">
          <div className="skeleton skeleton-circle" style={{ width: 32, height: 32 }} />
          <div className="skeleton skeleton-circle" style={{ width: 40, height: 40 }} />
          <div className="flex-grow">
            <div className="skeleton skeleton-text" style={{ width: '40%' }} />
            <div className="skeleton skeleton-text" style={{ width: '25%' }} />
          </div>
        </div>
      ))}
    </div>
  )
}

export function Leaderboard() {
  const store = useLeaderboardStore()
  const { loading, error, fetchLeaderboard, gameFilter } = store
  const paged = store.getPaged()
  const totalPages = store.getTotalPages()

  useEffect(() => {
    fetchLeaderboard()
  }, [fetchLeaderboard])

  return (
    <div className="py-8">
      <div className="container">
        <h1 className="text-3xl font-bold mb-2">Рейтинг</h1>
        <p className="mb-6 text-secondary">
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
              >
                {opt.label}
              </button>
            ))}
          </div>

          <div className="divider-vertical" />

          <div className="flex flex-wrap gap-2">
            {periodOptions.map((opt) => (
              <button
                key={opt.value}
                onClick={() => store.setPeriodFilter(opt.value as any)}
                className={`tab ${store.periodFilter === opt.value ? 'active' : ''}`}
              >
                {opt.label}
              </button>
            ))}
          </div>
        </div>

        {/* Loading */}
        {loading && <LeaderboardSkeleton />}

        {/* Error */}
        {!loading && error && (
          <div className="glass-card p-12 text-center">
            <Trophy size={48} color="var(--text-muted)" className="mx-auto mb-4" />
            <p className="text-muted mb-4">{error}</p>
            <button className="btn btn-secondary" onClick={() => fetchLeaderboard(gameFilter === 'all' ? undefined : gameFilter)}>
              <RefreshCw size={14} /> Повторить
            </button>
          </div>
        )}

        {/* Empty */}
        {!loading && !error && paged.length === 0 && (
          <div className="glass-card p-12 text-center">
            <Trophy size={48} color="var(--text-muted)" className="mx-auto mb-4" />
            <p className="text-muted">Рейтинг пока пуст</p>
          </div>
        )}

        {/* Rows */}
        {!loading && !error && paged.length > 0 && (
          <>
            {/* Table Header */}
            <div className="hidden md:flex items-center gap-4 px-4 py-2 mb-2 text-xs font-semibold text-muted">
              <div className="w-10">Место</div>
              <div className="w-10" />
              <div className="flex-grow">Игрок</div>
              {gameFilter !== 'all' && <div className="w-32">Рейтинг</div>}
              {gameFilter !== 'all' && <div className="w-16 text-right">ELO</div>}
              <div className="w-8" />
              <div className="flex gap-3">В / П / Н</div>
            </div>

            {paged.map((player, i) => (
              <motion.div
                key={player.sid}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: (i % store.perPage) * 0.05 }}
              >
                <LeaderboardRow
                  player={player}
                  rank={(store.page - 1) * store.perPage + i + 1}
                  gameFilter={store.gameFilter}
                />
              </motion.div>
            ))}

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-3 mt-6">
                <button
                  onClick={() => store.setPage(Math.max(1, store.page - 1))}
                  disabled={store.page === 1}
                  className={`btn btn-secondary ${store.page === 1 ? 'opacity-40' : ''}`}
                >
                  <ChevronLeft size={16} />
                </button>
                <span className="text-sm text-muted">
                  {store.page} из {totalPages}
                </span>
                <button
                  onClick={() => store.setPage(Math.min(totalPages, store.page + 1))}
                  disabled={store.page === totalPages}
                  className={`btn btn-secondary ${store.page === totalPages ? 'opacity-40' : ''}`}
                >
                  <ChevronRight size={16} />
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
