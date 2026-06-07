import { useState, useEffect, useCallback, useRef } from 'react'
import { motion } from 'framer-motion'
import { ArrowLeft, Users, Clock, Search, Flag, Handshake, Loader2 } from 'lucide-react'
import { Link, useParams, useNavigate } from 'react-router-dom'
import { ChessPiece } from '../components/ui/ChessPiece'
import { LiveIndicator } from '../components/ui/LiveIndicator'
import { useGameStore } from '../stores/gameStore'
import { useToastStore } from '../components/ui/Toast'
import { useUserStore } from '../stores/userStore'
import { getToken } from '../api/client'
import { InvitePlayerModal } from '../components/game/InvitePlayerModal'
import { InviteIncomingBanner } from '../components/game/InviteIncomingBanner'
import { ChessBoard } from '../components/game/ChessBoard'
import { CheckersBoard } from '../components/game/CheckersBoard'
import { BackgammonBoard } from '../components/game/BackgammonBoard'
import { GameClock } from '../components/game/GameClock'
import { MoveHistory } from '../components/game/MoveHistory'
import { ResignDialog } from '../components/game/ResignDialog'
import { GameResultModal } from '../components/game/GameResultModal'

// ─── WS message types ──────────────────────────────────────────────────────

interface WSStateMsg {
  type: 'state'
  board: string
  turn: string
  legal_moves: Array<{ from: string | [number, number]; to: string | [number, number]; die?: number }>
  clock: { white_ms: number; black_ms: number }
}

interface WSMoveAppliedMsg {
  type: 'move_applied'
  board: string
  turn: string
  legal_moves: Array<{ from: string | [number, number]; to: string | [number, number]; die?: number }>
  clock: { white_ms: number; black_ms: number }
  move: { from: string | number | [number, number]; to: string | number | [number, number]; color: string }
}

interface WSGameOverMsg {
  type: 'game_over'
  reason: string
  winner_sid: string
  winner_name: string
  score: string
}

interface WSInitMsg {
  type: 'init'
  color: string
}

type WSIncoming = WSStateMsg | WSMoveAppliedMsg | WSGameOverMsg | WSInitMsg

// ─── Game Component ─────────────────────────────────────────────────────────

export function Game() {
  const { gameId, matchId } = useParams<{ gameId: string; matchId: string }>()
  const { games } = useGameStore()
  const game = games.find((g) => g.id === gameId)
  const addToast = useToastStore((s) => s.addToast)
  const token = useUserStore((s) => s.token)
  const mySid = useUserStore((s) => s.sid)
  const navigate = useNavigate()

  const [inviteOpen, setInviteOpen] = useState(false)

  // ─── Active match state ──────────────────────────────────────────────────
  const [wsConnected, setWsConnected] = useState(false)
  const [connecting, setConnecting] = useState(!!matchId)
  const [myColor, setMyColor] = useState<string>('white')
  const [boardFEN, setBoardFEN] = useState<string>(
    'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1'
  )
  const [checkersBoard, setCheckersBoard] = useState<
    Array<Array<{ color: 'white' | 'black'; king: boolean } | null>>
  >(Array.from({ length: 8 }, () => Array.from({ length: 8 }, () => null)))
  const [bgPoints, setBgPoints] = useState<
    Array<{ color: 'white' | 'black'; count: number }>
  >(Array.from({ length: 24 }, () => ({ color: 'white', count: 0 })))
  const [bgBar, setBgBar] = useState({ white: 0, black: 0 })
  const [bgBorneOff, setBgBorneOff] = useState({ white: 0, black: 0 })
  const [bgDice, setBgDice] = useState<number[]>([])

  const [currentTurn, setCurrentTurn] = useState<string>('white')
  const [legalMoves, setLegalMoves] = useState<Array<{ from: string | [number, number]; to: string | [number, number]; die?: number }>>([])
  const [whiteMs, setWhiteMs] = useState(600000)
  const [blackMs, setBlackMs] = useState(600000)
  const [moveHistory, setMoveHistory] = useState<Array<{ from: string | number | [number, number]; to: string | number | [number, number]; color: string }>>([])
  const [resignOpen, setResignOpen] = useState(false)
  const [gameOverResult, setGameOverResult] = useState<{
    reason: string
    winner_sid: string
    winner_name: string
    score: string
  } | null>(null)

  const wsRef = useRef<WebSocket | null>(null)

  // ─── WebSocket connection ──────────────────────────────────────────────

  useEffect(() => {
    if (!matchId || !token) {
      setConnecting(false)
      return
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const wsUrl = `${protocol}//${host}/ws/game/${matchId}?token=${token}`

    setConnecting(true)
    const ws = new WebSocket(wsUrl)
    wsRef.current = ws

    ws.onopen = () => {
      setWsConnected(true)
      setConnecting(false)
      ws.send(JSON.stringify({ type: 'join' }))
    }

    ws.onmessage = (event) => {
      try {
        const msg: WSIncoming = JSON.parse(event.data)

        if (msg.type === 'init') {
          setMyColor(msg.color)
        }

        if (msg.type === 'state') {
          setBoardFEN(msg.board)
          setCurrentTurn(msg.turn)
          setLegalMoves(msg.legal_moves)
          setWhiteMs(msg.clock.white_ms)
          setBlackMs(msg.clock.black_ms)
        }

        if (msg.type === 'move_applied') {
          setBoardFEN(msg.board)
          setCurrentTurn(msg.turn)
          setLegalMoves(msg.legal_moves)
          setWhiteMs(msg.clock.white_ms)
          setBlackMs(msg.clock.black_ms)
          if (msg.move) {
            setMoveHistory((prev) => [...prev, msg.move])
          }
        }

        if (msg.type === 'game_over') {
          setGameOverResult({
            reason: msg.reason,
            winner_sid: msg.winner_sid,
            winner_name: msg.winner_name,
            score: msg.score,
          })
        }
      } catch {
        // Ignore malformed messages
      }
    }

    ws.onclose = () => {
      setWsConnected(false)
      setConnecting(false)
      addToast('Соединение с сервером закрыто', 'info')
    }

    ws.onerror = () => {
      addToast('Ошибка подключения к серверу', 'error')
    }

    return () => {
      ws.close()
      wsRef.current = null
    }
  }, [matchId, token, addToast])

  // ─── WS send helper ──────────────────────────────────────────────────────

  const sendWs = useCallback((data: unknown) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data))
    }
  }, [])

  // ─── Move handlers ───────────────────────────────────────────────────────

  const handleChessMove = useCallback((from: string, to: string) => {
    sendWs({ type: 'move', from, to })
  }, [sendWs])

  const handleCheckersMove = useCallback((from: [number, number], to: [number, number]) => {
    sendWs({ type: 'move', from, to })
  }, [sendWs])

  const handleBgMove = useCallback((from: number, to: number, die: number) => {
    sendWs({ type: 'move', from, to, die })
  }, [sendWs])

  const handleBgRollDice = useCallback(() => {
    sendWs({ type: 'roll_dice' })
  }, [sendWs])

  const handleResign = useCallback(() => {
    sendWs({ type: 'resign' })
    setResignOpen(false)
  }, [sendWs])

  // ─── Lobby view (no active match) ─────────────────────────────────────────

  if (!matchId) {
    if (!game) {
      return (
        <div className="py-16 text-center">
          <h1 className="text-2xl font-bold mb-2">Игра не найдена</h1>
          <Link to="/" className="btn btn-primary mt-4 inline-flex no-underline">
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

          <InviteIncomingBanner />

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div>
              <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
                <div className="flex items-center gap-3 mb-4">
                  <h1 className="text-3xl font-bold">{game.name}</h1>
                  {game.isLive && <LiveIndicator />}
                </div>
                <p className="mb-6 text-secondary">{game.description}</p>

                <div className="glass-card p-6 mb-4">
                  <h3 className="text-sm font-semibold mb-4 text-secondary">Быстрый матч</h3>
                  <div className="flex items-center gap-3 mb-4">
                    <Users size={16} color="var(--text-muted)" />
                    <span className="text-sm text-secondary">Онлайн: {game.playerCount}</span>
                  </div>
                  <div className="flex items-center gap-3 mb-4">
                    <Clock size={16} color="var(--text-muted)" />
                    <span className="text-sm text-secondary">Среднее время партии: 25 мин</span>
                  </div>
                  <button className="btn btn-primary w-full justify-center text-base" onClick={() => setInviteOpen(true)}>
                    <Search size={16} />
                    Найти матч
                  </button>
                </div>

                <div className="glass-card p-6">
                  <h3 className="text-sm font-semibold mb-4 text-secondary">Правила</h3>
                  <p className="text-sm text-muted">
                    {gameId === 'chess'
                      ? 'Классические шахматы по стандартным правилам ФИДЕ. Контроль времени: 10 минут + 5 секунд за ход. Рейтинг ELO рассчитывается по системе Эло.'
                      : gameId === 'checkers'
                        ? 'Классические шашки на доске 8x8. Обязательное взятие, многократные прыжки. Победа — capture всех фигур или блокировка.'
                        : gameId === 'backgammon'
                          ? 'Классические нарды (длинные). 15 шашек каждый, движение по 24 пунктам. Бросайте кубики и перемещайте шашки к дому.'
                          : 'Выберите игру и прочтите правила в разделе справки.'}
                  </p>
                </div>
              </motion.div>
            </div>

            <div>
              {gameId === 'chess' && (
                <motion.div initial={{ opacity: 0, scale: 0.95 }} animate={{ opacity: 1, scale: 1 }} className="glass-card p-6">
                  <h3 className="text-sm font-semibold mb-4 text-secondary">Доска</h3>
                  <div className="flex justify-center">
                    <div className="chess-board-small">
                      {Array.from({ length: 64 }, (_, i) => {
                        const row = Math.floor(i / 8)
                        const col = i % 8
                        const isLight = (row + col) % 2 === 0
                        return (
                          <div key={i} className={`chess-square-sm ${isLight ? 'chess-square-light' : 'chess-square-dark'}`}>
                            {row === 0 && <ChessPiece type={['rook', 'knight', 'bishop', 'queen', 'king', 'bishop', 'knight', 'rook'][col] as 'rook'} color="black" size={30} />}
                            {row === 1 && <ChessPiece type="pawn" color="black" size={28} />}
                            {row === 6 && <ChessPiece type="pawn" color="white" size={28} />}
                            {row === 7 && <ChessPiece type={['rook', 'knight', 'bishop', 'queen', 'king', 'bishop', 'knight', 'rook'][col] as 'rook'} color="white" size={30} />}
                          </div>
                        )
                      })}
                    </div>
                  </div>
                </motion.div>
              )}

              {gameId !== 'chess' && (
                <motion.div initial={{ opacity: 0, scale: 0.95 }} animate={{ opacity: 1, scale: 1 }} className="glass-card p-12 flex flex-col items-center justify-center text-center">
                  <div className="game-placeholder-icon mb-4">
                    <span className="text-3xl font-bold text-accent">?</span>
                  </div>
                  <h3 className="text-lg font-bold mb-2">{game.name}</h3>
                  <p className="text-sm text-muted">Визуализация доски будет доступна при запуске игры</p>
                </motion.div>
              )}
            </div>
          </div>
        </div>

        <InvitePlayerModal isOpen={inviteOpen} onClose={() => setInviteOpen(false)} gameType={gameId} />
      </div>
    )
  }

  // ─── Active match view ──────────────────────────────────────────────────

  if (connecting) {
    return (
      <div className="py-8">
        <div className="container">
          <Link to={`/game/${gameId}`} className="btn btn-ghost mb-6 inline-flex no-underline">
            <ArrowLeft size={16} />
            Назад
          </Link>
          <div className="flex flex-col items-center justify-center" style={{ minHeight: 400 }}>
            <Loader2 size={32} className="spin mb-4" color="var(--gold)" />
            <span className="text-secondary">Подключение к матчу...</span>
          </div>
        </div>
      </div>
    )
  }

  const gameName = gameId === 'chess' ? 'Шахматы' : gameId === 'checkers' ? 'Шашки' : 'Нарды'

  const legalMovesChess = legalMoves.map((m) => ({
    from: String(m.from),
    to: String(m.to),
  }))

  const legalMovesCheckers = legalMoves.map((m) => ({
    from: m.from as [number, number],
    to: m.to as [number, number],
  }))

  const legalMovesBg = legalMoves.map((m) => ({
    from: m.from as unknown as number,
    to: m.to as unknown as number,
    die: m.die || 0,
  }))

  return (
    <div className="py-8">
      <div className="container">
        <div className="flex items-center gap-3 mb-4">
          <Link to={`/game/${gameId}`} className="btn btn-ghost inline-flex no-underline">
            <ArrowLeft size={16} />
          </Link>
          <div className="flex items-center gap-2">
            <h1 className="text-xl font-bold">{gameName}</h1>
            {wsConnected && <LiveIndicator />}
          </div>
          <span className="text-xs text-muted font-mono ml-2">match: {matchId.slice(0, 8)}</span>
        </div>

        <div className="flex items-center justify-center gap-2 mb-4">
          <GameClock whiteMs={whiteMs} blackMs={blackMs} activeColor={currentTurn} />
        </div>

        <div className="game-board-layout">
          {/* Board */}
          <div className="flex items-start justify-center">
            {(gameId === 'chess') && (
              <ChessBoard
                boardFEN={boardFEN}
                legalMoves={legalMovesChess}
                currentTurn={currentTurn}
                myColor={myColor}
                onMove={handleChessMove}
              />
            )}

            {gameId === 'checkers' && (
              <CheckersBoard
                board={checkersBoard}
                legalMoves={legalMovesCheckers}
                currentTurn={currentTurn}
                myColor={myColor}
                onMove={handleCheckersMove}
              />
            )}

            {gameId === 'backgammon' && (
              <BackgammonBoard
                points={bgPoints}
                bar={bgBar}
                borneOff={bgBorneOff}
                dice={bgDice}
                legalMoves={legalMovesBg}
                currentTurn={currentTurn}
                myColor={myColor}
                onMove={handleBgMove}
                onRollDice={handleBgRollDice}
              />
            )}
          </div>

          {/* Side panel */}
          <div className="flex flex-col gap-3">
            <MoveHistory moves={moveHistory} />

            <div className="glass-card p-3 flex flex-col gap-2">
              <div className="flex items-center gap-2 text-xs text-muted">
                <span style={{
                  width: 10,
                  height: 10,
                  borderRadius: '50%',
                  background: myColor === 'white' ? '#fff' : '#1a1a2e',
                  border: myColor === 'white' ? '1px solid #c9b896' : '1px solid #3a3a5a',
                  flexShrink: 0,
                }} />
                Вы: {myColor === 'white' ? 'Белые' : 'Чёрные'}
              </div>
              <div className="text-xs text-muted">
                Ход: {currentTurn === 'white' ? 'Белые' : 'Чёрные'}
              </div>
            </div>

            {!gameOverResult && (
              <div className="flex gap-2">
                <button className="btn btn-ghost btn-sm flex-1 justify-center" onClick={() => setResignOpen(true)}>
                  <Flag size={14} />
                  Сдаться
                </button>
                <button className="btn btn-ghost btn-sm flex-1 justify-center" onClick={() => addToast('Предложение ничьи отправлено', 'info')}>
                  <Handshake size={14} />
                  Ничья
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Modals */}
      <ResignDialog
        isOpen={resignOpen}
        onClose={() => setResignOpen(false)}
        onResign={handleResign}
      />

      <GameResultModal
        isOpen={!!gameOverResult}
        onClose={() => {
          setGameOverResult(null)
          navigate(`/game/${gameId}`)
        }}
        result={gameOverResult ? { ...gameOverResult, my_sid: mySid } : null}
      />
    </div>
  )
}
