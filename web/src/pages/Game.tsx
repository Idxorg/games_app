import { useState } from 'react'
import { motion } from 'framer-motion'
import { ArrowLeft, Users, Clock, Search } from 'lucide-react'
import { Link, useParams } from 'react-router-dom'
import { ChessPiece } from '../components/ui/ChessPiece'
import { LiveIndicator } from '../components/ui/LiveIndicator'
import { useGameStore } from '../stores/gameStore'
import { useToastStore } from '../components/ui/Toast'
import { InvitePlayerModal } from '../components/game/InvitePlayerModal'
import { InviteIncomingBanner } from '../components/game/InviteIncomingBanner'

export function Game() {
  const { gameId } = useParams<{ gameId: string }>()
  const { games } = useGameStore()
  const game = games.find((g) => g.id === gameId)
  const addToast = useToastStore((s) => s.addToast)
  const [inviteOpen, setInviteOpen] = useState(false)

  if (!game) {
    return (
      <div className="py-16 text-center">
        <h1 className="text-2xl font-bold mb-2">Игра не найдена</h1>
        <Link to="/" className="btn btn-primary mt-4 inline-block no-underline">
          <ArrowLeft size={16} />
          На главную
        </Link>
      </div>
    )
  }

  const handleFindMatch = () => {
    setInviteOpen(true)
  }

  const handleStartPlay = () => {
    addToast('Запуск игры скоро будет доступен', 'info')
  }

  return (
    <div className="py-8">
      <div className="container">
        <Link to="/" className="btn btn-ghost mb-6 inline-flex no-underline">
          <ArrowLeft size={16} />
          Все игры
        </Link>

        {/* Incoming invite banner */}
        <InviteIncomingBanner />

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Game Info */}
          <div>
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
            >
              <div className="flex items-center gap-3 mb-4">
                <h1 className="text-3xl font-bold">{game.name}</h1>
                {game.isLive && <LiveIndicator />}
              </div>
              <p className="mb-6 text-secondary">{game.description}</p>

              <div className="glass-card p-6 mb-4">
                <h3 className="text-sm font-semibold mb-4 text-secondary">Быстрый матч</h3>
                <div className="flex items-center gap-3 mb-4">
                  <Users size={16} color="var(--text-muted)" />
                  <span className="text-sm text-secondary">
                    Онлайн: {game.playerCount}
                  </span>
                </div>
                <div className="flex items-center gap-3 mb-4">
                  <Clock size={16} color="var(--text-muted)" />
                  <span className="text-sm text-secondary">
                    Среднее время партии: 25 мин
                  </span>
                </div>
                <button className="btn btn-primary w-full justify-center text-base" onClick={handleFindMatch}>
                  <Search size={16} />
                  Найти матч
                </button>
              </div>

              <div className="glass-card p-6">
                <h3 className="text-sm font-semibold mb-4 text-secondary">Правила</h3>
                <p className="text-sm text-muted">
                  {gameId === 'chess'
                    ? 'Классические шахматы по стандартным правилам ФИДЕ. Контроль времени: 10 минут + 5 секунд за ход. Рейтинг ELO рассчитывается по системе Эло.'
                    : 'Выберите игру и прочтите правила в разделе справки.'}
                </p>
              </div>
            </motion.div>
          </div>

          {/* Board Preview */}
          <div>
            {gameId === 'chess' && (
              <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="glass-card p-6"
              >
                <h3 className="text-sm font-semibold mb-4 text-secondary">Доска</h3>
                <div className="flex justify-center">
                  <div className="chess-board-small">
                    {Array.from({ length: 64 }, (_, i) => {
                      const row = Math.floor(i / 8)
                      const col = i % 8
                      const isLight = (row + col) % 2 === 0
                      return (
                        <div
                          key={i}
                          className={`chess-square-sm ${isLight ? 'chess-square-light' : 'chess-square-dark'}`}
                        >
                          {row === 0 && (
                            <ChessPiece
                              type={['rook', 'knight', 'bishop', 'queen', 'king', 'bishop', 'knight', 'rook'][col] as any}
                              color="black"
                              size={30}
                            />
                          )}
                          {row === 1 && <ChessPiece type="pawn" color="black" size={28} />}
                          {row === 6 && <ChessPiece type="pawn" color="white" size={28} />}
                          {row === 7 && (
                            <ChessPiece
                              type={['rook', 'knight', 'bishop', 'queen', 'king', 'bishop', 'knight', 'rook'][col] as any}
                              color="white"
                              size={30}
                            />
                          )}
                        </div>
                      )
                    })}
                  </div>
                </div>
              </motion.div>
            )}

            {gameId !== 'chess' && (
              <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="glass-card p-12 flex flex-col items-center justify-center text-center"
              >
                <div className="game-placeholder-icon mb-4">
                  <span className="text-3xl font-bold text-accent">?</span>
                </div>
                <h3 className="text-lg font-bold mb-2">{game.name}</h3>
                <p className="text-sm text-muted">
                  Визуализация доски будет доступна при запуске игры
                </p>
                <button className="btn btn-primary mt-6" onClick={handleStartPlay}>
                  Начать играть
                </button>
              </motion.div>
            )}
          </div>
        </div>
      </div>

      {/* Invite player modal */}
      <InvitePlayerModal
        isOpen={inviteOpen}
        onClose={() => setInviteOpen(false)}
        gameType={gameId}
      />
    </div>
  )
}
