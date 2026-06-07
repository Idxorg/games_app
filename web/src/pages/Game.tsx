import { motion } from 'framer-motion'
import { ArrowLeft, Users, Clock } from 'lucide-react'
import { Link, useParams } from 'react-router-dom'
import { ChessPiece } from '../components/ui/ChessPiece'
import { LiveIndicator } from '../components/ui/LiveIndicator'
import { games } from '../data/games'

export function Game() {
  const { gameId } = useParams<{ gameId: string }>()
  const game = games.find((g) => g.id === gameId)

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

  return (
    <div className="py-8">
      <div className="container">
        <Link to="/" className="btn btn-ghost mb-6 inline-flex no-underline">
          <ArrowLeft size={16} />
          Все игры
        </Link>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Game Info */}
          <div>
            <div className="flex items-center gap-3 mb-4">
              <h1 className="text-3xl font-bold">{game.name}</h1>
              {game.isLive && <LiveIndicator />}
            </div>
            <p className="mb-6" style={{ color: 'var(--text-secondary)' }}>{game.description}</p>

            <div className="glass-card p-6 mb-4">
              <h3 className="text-sm font-semibold mb-4" style={{ color: 'var(--text-secondary)' }}>Быстрый матч</h3>
              <div className="flex items-center gap-3 mb-4">
                <Users size={16} color="var(--text-muted)" />
                <span className="text-sm" style={{ color: 'var(--text-secondary)' }}>
                  Онлайн: {Math.floor(Math.random() * 50) + 10} игроков
                </span>
              </div>
              <div className="flex items-center gap-3 mb-4">
                <Clock size={16} color="var(--text-muted)" />
                <span className="text-sm" style={{ color: 'var(--text-secondary)' }}>
                  Среднее время партии: 25 мин
                </span>
              </div>
              <button className="btn btn-primary w-full justify-center text-base">
                Найти матч
              </button>
            </div>

            <div className="glass-card p-6">
              <h3 className="text-sm font-semibold mb-4" style={{ color: 'var(--text-secondary)' }}>Правила</h3>
              <p className="text-sm" style={{ color: 'var(--text-muted)' }}>
                {gameId === 'chess'
                  ? 'Классические шахматы по стандартным правилам ФИДЕ. Контроль времени: 10 минут + 5 секунд за ход. Рейтинг ELO рассчитывается по системе Эло.'
                  : 'Выберите игру и прочтите правила в разделе справки.'}
              </p>
            </div>
          </div>

          {/* Board Preview */}
          <div>
            {gameId === 'chess' && (
              <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="glass-card p-6"
              >
                <h3 className="text-sm font-semibold mb-4" style={{ color: 'var(--text-secondary)' }}>Доска</h3>
                <div className="flex justify-center">
                  <div className="grid grid-cols-8 gap-0" style={{ borderRadius: 8, overflow: 'hidden', boxShadow: '0 10px 40px rgba(0,0,0,0.3)' }}>
                    {Array.from({ length: 64 }, (_, i) => {
                      const row = Math.floor(i / 8)
                      const col = i % 8
                      const isLight = (row + col) % 2 === 0
                      return (
                        <div
                          key={i}
                          style={{
                            width: 40,
                            height: 40,
                            background: isLight ? '#1a1a2e' : 'rgba(212,168,67,0.12)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                          }}
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
              <div className="glass-card p-12 flex flex-col items-center justify-center text-center">
                <div className="w-20 h-20 rounded-2xl flex items-center justify-center mb-4" style={{ background: 'rgba(212,168,67,0.08)' }}>
                  <span className="text-3xl font-bold" style={{ color: 'var(--gold)' }}>?</span>
                </div>
                <h3 className="text-lg font-bold mb-2">{game.name}</h3>
                <p className="text-sm" style={{ color: 'var(--text-muted)' }}>
                  Визуализация доски будет доступна при запуске игры
                </p>
                <button className="btn btn-primary mt-6">
                  Начать играть
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
