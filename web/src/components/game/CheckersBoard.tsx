import { useState, useCallback } from 'react'

// ─── Types ──────────────────────────────────────────────────────────────────

interface CheckersPiece {
  color: 'white' | 'black'
  king: boolean
}

interface LegalMove {
  from: [number, number]
  to: [number, number]
}

interface CheckersBoardProps {
  board: (CheckersPiece | null)[][]
  legalMoves: LegalMove[]
  currentTurn: string
  myColor: string
  onMove: (from: [number, number], to: [number, number]) => void
}

// ─── Component ──────────────────────────────────────────────────────────────

export function CheckersBoard({ board, legalMoves, currentTurn, myColor, onMove }: CheckersBoardProps) {
  const [selected, setSelected] = useState<[number, number] | null>(null)

  const isMyTurn = currentTurn === myColor

  const getTargetsForSelected = useCallback((from: [number, number]) => {
    return legalMoves.filter((m) => m.from[0] === from[0] && m.from[1] === from[1])
  }, [legalMoves])

  const handleClick = useCallback((row: number, col: number) => {
    // Only dark squares are playable
    if ((row + col) % 2 === 0) return

    const piece = board[row]?.[col]

    if (selected) {
      const moves = getTargetsForSelected(selected)
      const targetMove = moves.find((m) => m.to[0] === row && m.to[1] === col)
      if (targetMove) {
        onMove(selected, targetMove.to)
        // Multi-jump: if the target is also a "from" in remaining moves, auto-select
        const afterMoves = legalMoves.filter(
          (m) => m.from[0] === targetMove.to[0] && m.from[1] === targetMove.to[1]
        )
        if (afterMoves.length > 0) {
          setSelected(targetMove.to)
        } else {
          setSelected(null)
        }
        return
      }
    }

    if (piece && piece.color === myColor && isMyTurn) {
      const moves = legalMoves.filter((m) => m.from[0] === row && m.from[1] === col)
      if (moves.length > 0) {
        setSelected([row, col])
      } else {
        setSelected(null)
      }
    } else {
      setSelected(null)
    }
  }, [selected, board, myColor, isMyTurn, legalMoves, onMove, getTargetsForSelected])

  const targets = selected
    ? getTargetsForSelected(selected).map((m) => m.to as unknown as string)
    : []

  return (
    <div className="interactive-board">
      <div style={{
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
          const piece = board[row]?.[col]
          const isSelected = selected?.[0] === row && selected?.[1] === col
          const isTarget = targets.some(
            (t) => (t as unknown as [number, number])[0] === row && (t as unknown as [number, number])[1] === col
          )
          const isPlayable = !isLight

          return (
            <div
              key={i}
              onClick={() => handleClick(row, col)}
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                position: 'relative',
                cursor: isMyTurn && isPlayable && (isTarget || (piece && piece.color === myColor)) ? 'pointer' : 'default',
                background: isSelected
                  ? 'var(--gold-light-bg)'
                  : isLight
                    ? '#1e1e3a'
                    : '#2a2a4a',
                transition: 'background 0.15s ease',
              }}
            >
              {/* Piece */}
              {piece && (
                <div style={{
                  width: '70%',
                  aspectRatio: '1',
                  borderRadius: '50%',
                  background: piece.color === 'white'
                    ? 'linear-gradient(135deg, #e8c56a, #d4a843)'
                    : 'linear-gradient(135deg, #4a4a6a, #3a3a5a)',
                  border: piece.color === 'white'
                    ? '2px solid #b8922e'
                    : '2px solid #8888a0',
                  boxShadow: piece.king
                    ? `inset 0 0 0 3px ${piece.color === 'white' ? '#b8922e' : '#8888a0'}, 0 2px 4px rgba(0,0,0,0.3)`
                    : '0 2px 4px rgba(0,0,0,0.3)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: 'clamp(10px, 2vw, 16px)',
                  fontWeight: 700,
                  color: piece.color === 'white' ? '#0a0a0f' : '#e8e8f0',
                }}>
                  {piece.king && '♛'}
                </div>
              )}

              {/* Legal move dot */}
              {isTarget && !piece && (
                <div style={{
                  position: 'absolute',
                  width: '28%',
                  height: '28%',
                  borderRadius: '50%',
                  background: 'rgba(212, 168, 67, 0.5)',
                }} />
              )}

              {/* Legal move ring (capture) */}
              {isTarget && piece && (
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
    </div>
  )
}
