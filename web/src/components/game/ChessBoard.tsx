import { useState, useCallback } from 'react'
import { motion } from 'framer-motion'

// ─── Unicode chess pieces ──────────────────────────────────────────────────

const PIECE_UNICODE: Record<string, Record<string, string>> = {
  white: { king: '♔', queen: '♕', rook: '♖', bishop: '♗', knight: '♘', pawn: '♙' },
  black: { king: '♚', queen: '♛', rook: '♜', bishop: '♝', knight: '♞', pawn: '♟' },
}

// FEN → piece lookup
function fenToBoard(fen: string): string[][] {
  const rows = fen.split(' ')[0].split('/')
  return rows.map((row) => {
    const pieces: string[] = []
    for (const ch of row) {
      if (/\d/.test(ch)) {
        for (let i = 0; i < parseInt(ch); i++) pieces.push('')
      } else {
        pieces.push(ch)
      }
    }
    return pieces
  })
}

function parsePiece(fenChar: string): { color: string; piece: string } | null {
  if (!fenChar) return null
  const color = fenChar === fenChar.toUpperCase() ? 'white' : 'black'
  const piece = fenChar.toLowerCase()
  const map: Record<string, string> = {
    k: 'king', q: 'queen', r: 'rook', b: 'bishop', n: 'knight', p: 'pawn',
  }
  return map[piece] ? { color, piece: map[piece] } : null
}

// ─── Types ──────────────────────────────────────────────────────────────────

interface LegalMove {
  from: string
  to: string
}

interface ChessBoardProps {
  boardFEN: string
  legalMoves: LegalMove[]
  currentTurn: string
  myColor: string
  onMove: (from: string, to: string) => void
}

// ─── Component ──────────────────────────────────────────────────────────────

export function ChessBoard({ boardFEN, legalMoves, currentTurn, myColor, onMove }: ChessBoardProps) {
  const [selected, setSelected] = useState<string | null>(null)
  const [lastMove, setLastMove] = useState<{ from: string; to: string } | null>(null)

  const board = fenToBoard(boardFEN)
  const isMyTurn = currentTurn === myColor
  const files = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h']
  const ranks = ['8', '7', '6', '5', '4', '3', '2', '1']

  const getTargetsForSelected = useCallback((from: string) => {
    return legalMoves.filter((m) => m.from === from).map((m) => m.to)
  }, [legalMoves])

  const handleClick = useCallback((row: number, col: number) => {
    const sq = `${files[col]}${ranks[row]}`
    const fenPiece = board[row]?.[col] || ''
    const parsed = parsePiece(fenPiece)

    if (selected) {
      const targets = getTargetsForSelected(selected)
      if (targets.includes(sq)) {
        setLastMove({ from: selected, to: sq })
        onMove(selected, sq)
        setSelected(null)
        return
      }
    }

    if (parsed && parsed.color === myColor && isMyTurn) {
      setSelected(sq)
    } else {
      setSelected(null)
    }
  }, [selected, board, myColor, isMyTurn, legalMoves, onMove, getTargetsForSelected])

  const targets = selected ? getTargetsForSelected(selected) : []

  return (
    <div className="interactive-board">
      {/* Rank labels */}
      <div className="board-labels-left">
        {ranks.map((r, i) => (
          <div key={r} className="board-label-cell">{r}</div>
        ))}
      </div>

      {/* Board */}
      <div className="board-grid" style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(8, 1fr)',
        aspectRatio: '1',
        borderRadius: 'var(--radius-sm)',
        overflow: 'hidden',
        width: '100%',
        maxWidth: 480,
      }}>
        {Array.from({ length: 64 }, (_, i) => {
          const row = Math.floor(i / 8)
          const col = i % 8
          const isLight = (row + col) % 2 === 0
          const sq = `${files[col]}${ranks[row]}`
          const fenPiece = board[row]?.[col] || ''
          const parsed = parsePiece(fenPiece)
          const isSelected = sq === selected
          const isLastFrom = lastMove?.from === sq
          const isLastTo = lastMove?.to === sq
          const isTarget = targets.includes(sq)

          return (
            <div
              key={i}
              onClick={() => handleClick(row, col)}
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                position: 'relative',
                cursor: isMyTurn && (isTarget || (parsed && parsed.color === myColor)) ? 'pointer' : 'default',
                background: isSelected
                  ? 'var(--gold-light-bg)'
                  : isLastFrom || isLastTo
                    ? 'rgba(212, 168, 67, 0.25)'
                    : isLight
                      ? '#1e1e3a'
                      : '#2a2a4a',
                transition: 'background 0.15s ease',
              }}
            >
              {/* Piece */}
              {parsed && (
                <motion.span
                  key={`${sq}-${fenPiece}`}
                  initial={false}
                  style={{
                    fontSize: 'clamp(24px, 5vw, 42px)',
                    lineHeight: 1,
                    userSelect: 'none',
                    filter: parsed.color === 'white'
                      ? 'drop-shadow(0 1px 2px rgba(0,0,0,0.5))'
                      : 'drop-shadow(0 1px 2px rgba(0,0,0,0.3))',
                    color: parsed.color === 'white' ? '#fff' : '#1a1a2e',
                  }}
                >
                  {PIECE_UNICODE[parsed.color]?.[parsed.piece] || fenPiece}
                </motion.span>
              )}

              {/* Legal move dot */}
              {isTarget && !parsed && (
                <div style={{
                  position: 'absolute',
                  width: '28%',
                  height: '28%',
                  borderRadius: '50%',
                  background: 'rgba(212, 168, 67, 0.5)',
                }} />
              )}

              {/* Legal move ring (capture) */}
              {isTarget && parsed && (
                <div style={{
                  position: 'absolute',
                  inset: 4,
                  borderRadius: '50%',
                  border: '3px solid rgba(212, 168, 67, 0.5)',
                }} />
              )}
            </div>
          )
        })}
      </div>

      {/* File labels */}
      <div className="board-labels-bottom">
        {files.map((f) => (
          <div key={f} className="board-label-cell">{f}</div>
        ))}
      </div>
    </div>
  )
}
