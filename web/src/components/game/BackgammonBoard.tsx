import { useState, useCallback } from 'react'
import { Dice5 } from 'lucide-react'

// ─── Types ──────────────────────────────────────────────────────────────────

interface PointData {
  color: 'white' | 'black'
  count: number
}

interface BarData {
  white: number
  black: number
}

interface BorneOffData {
  white: number
  black: number
}

interface LegalMove {
  from: number
  to: number
  die: number
}

interface BackgammonBoardProps {
  points: PointData[]
  bar: BarData
  borneOff: BorneOffData
  dice: number[]
  legalMoves: LegalMove[]
  currentTurn: string
  myColor: string
  onMove: (from: number, to: number, die: number) => void
  onRollDice: () => void
}

// ─── Dice display ──────────────────────────────────────────────────────────

const DIE_DOTS: Record<number, number[][]> = {
  1: [[1, 1]],
  2: [[0, 0], [2, 2]],
  3: [[0, 0], [1, 1], [2, 2]],
  4: [[0, 0], [0, 2], [2, 0], [2, 2]],
  5: [[0, 0], [0, 2], [1, 1], [2, 0], [2, 2]],
  6: [[0, 0], [0, 2], [1, 0], [1, 2], [2, 0], [2, 2]],
}

function DieFace({ value }: { value: number }) {
  const dots = DIE_DOTS[value] || []
  return (
    <div style={{
      width: 36,
      height: 36,
      borderRadius: 6,
      background: 'var(--bg-tertiary)',
      border: '1px solid var(--bg-glass-border)',
      display: 'grid',
      gridTemplateColumns: 'repeat(3, 1fr)',
      gridTemplateRows: 'repeat(3, 1fr)',
      padding: 5,
      flexShrink: 0,
    }}>
      {dots.map(([r, c], i) => (
        <div key={i} style={{
          gridRow: r + 1,
          gridColumn: c + 1,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}>
          <div style={{
            width: 5,
            height: 5,
            borderRadius: '50%',
            background: 'var(--gold)',
          }} />
        </div>
      ))}
    </div>
  )
}

// ─── Checker piece ─────────────────────────────────────────────────────────

function Checker({ color, isTop }: { color: 'white' | 'black'; isTop: boolean }) {
  return (
    <div style={{
      width: '80%',
      aspectRatio: '1',
      borderRadius: '50%',
      background: color === 'white'
        ? 'linear-gradient(135deg, #e8c56a, #d4a843)'
        : 'linear-gradient(135deg, #4a4a6a, #3a3a5a)',
      border: color === 'white' ? '2px solid #b8922e' : '2px solid #8888a0',
      boxShadow: '0 1px 3px rgba(0,0,0,0.3)',
      flexShrink: 0,
    }} />
  )
}

// ─── Component ──────────────────────────────────────────────────────────────

export function BackgammonBoard({
  points,
  bar,
  borneOff,
  dice,
  legalMoves,
  currentTurn,
  myColor,
  onMove,
  onRollDice,
}: BackgammonBoardProps) {
  const [selectedPoint, setSelectedPoint] = useState<number | null>(null)
  const [selectedBar, setSelectedBar] = useState<boolean>(false)

  const isMyTurn = currentTurn === myColor
  const needsRoll = isMyTurn && dice.length === 0

  const getTargets = useCallback((from: number) => {
    return legalMoves.filter((m) => m.from === from)
  }, [legalMoves])

  const getBarTargets = useCallback(() => {
    return legalMoves.filter((m) => m.from === -1)
  }, [legalMoves])

  const handlePointClick = useCallback((ptIdx: number) => {
    if (!isMyTurn) return

    // If something is selected and clicking a target
    if (selectedPoint !== null) {
      const targets = getTargets(selectedPoint)
      const match = targets.find((t) => t.to === ptIdx)
      if (match) {
        onMove(selectedPoint, ptIdx, match.die)
        setSelectedPoint(null)
        return
      }
    }

    if (selectedBar) {
      const barMoves = getBarTargets()
      const match = barMoves.find((t) => t.to === ptIdx)
      if (match) {
        onMove(-1, ptIdx, match.die)
        setSelectedBar(false)
        return
      }
    }

    // Select a point with my checkers
    const pt = points[ptIdx]
    if (pt && pt.color === myColor && pt.count > 0) {
      const moves = getTargets(ptIdx)
      if (moves.length > 0) {
        setSelectedPoint(ptIdx)
        setSelectedBar(false)
      } else {
        setSelectedPoint(null)
      }
    } else {
      setSelectedPoint(null)
    }
  }, [isMyTurn, selectedPoint, selectedBar, points, myColor, legalMoves, onMove, getTargets, getBarTargets])

  const handleBarClick = useCallback(() => {
    if (!isMyTurn) return
    const barCount = myColor === 'white' ? bar.white : bar.black
    if (barCount > 0) {
      const barMoves = getBarTargets()
      if (barMoves.length > 0) {
        setSelectedBar(true)
        setSelectedPoint(null)
      }
    }
  }, [isMyTurn, myColor, bar, getBarTargets])

  const allTargets = [
    ...(selectedPoint !== null ? getTargets(selectedPoint) : []),
    ...(selectedBar ? getBarTargets() : []),
  ]

  const targetSet = new Set(allTargets.map((t) => t.to))

  const pointIndex = (quadrant: number, pos: number): number => {
    // Top row: points 13-24 left to right → indices 12-23
    // Bottom row: points 1-12 right to left → indices 0-11
    if (quadrant === 0) return 12 + pos      // top-left: 13-18
    if (quadrant === 1) return 18 + pos      // top-right: 19-24
    if (quadrant === 2) return 11 - pos       // bottom-right: 12-7
    return 5 - pos                            // bottom-left: 6-1
  }

  const renderPoint = (quadrant: number, pos: number, isTop: boolean) => {
    const idx = pointIndex(quadrant, pos)
    const pt = points[idx]
    if (!pt) return null
    const isTarget = targetSet.has(idx)
    const isSelected = selectedPoint === idx

    return (
      <div
        key={`q${quadrant}-p${pos}`}
        onClick={() => handlePointClick(idx)}
        style={{
          flex: 1,
          display: 'flex',
          flexDirection: isTop ? 'column' : 'column-reverse',
          alignItems: 'center',
          position: 'relative',
          cursor: isMyTurn && (isTarget || (pt.count > 0 && pt.color === myColor)) ? 'pointer' : 'default',
          padding: '0 1px',
        }}
      >
        {/* Triangle */}
        <div style={{
          width: '100%',
          flex: 1,
          maxWidth: 48,
          clipPath: isTop
            ? 'polygon(0% 100%, 50% 0%, 100% 100%)'
            : 'polygon(0% 0%, 50% 100%, 100% 0%)',
          background: isSelected
            ? 'rgba(212, 168, 67, 0.4)'
            : isTarget
              ? 'rgba(212, 168, 67, 0.2)'
              : idx % 2 === 0
                ? '#2a2a4a'
                : '#1e1e3a',
          transition: 'background 0.15s ease',
        }} />

        {/* Target indicator */}
        {isTarget && (
          <div style={{
            position: 'absolute',
            bottom: isTop ? undefined : 4,
            top: isTop ? 4 : undefined,
            width: 20,
            height: 20,
            borderRadius: '50%',
            background: 'rgba(212, 168, 67, 0.5)',
            zIndex: 2,
          }} />
        )}

        {/* Stacked checkers */}
        {pt.count > 0 && (
          <div style={{
            position: 'absolute',
            bottom: isTop ? '30%' : 'auto',
            top: isTop ? 'auto' : '30%',
            display: 'flex',
            flexDirection: isTop ? 'column-reverse' : 'column',
            gap: 1,
            zIndex: 1,
          }}>
            {Array.from({ length: Math.min(pt.count, 5) }, (_, i) => (
              <Checker key={i} color={pt.color} isTop={isTop} />
            ))}
            {pt.count > 5 && (
              <div style={{
                position: 'absolute',
                inset: 0,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: 11,
                fontWeight: 700,
                color: pt.color === 'white' ? '#0a0a0f' : '#e8e8f0',
              }}>
                {pt.count}
              </div>
            )}
          </div>
        )}
      </div>
    )
  }

  return (
    <div style={{ width: '100%', maxWidth: 600 }}>
      {/* Dice + Controls */}
      <div className="flex items-center justify-center gap-3 mb-3">
        {dice.map((d, i) => <DieFace key={i} value={d} />)}
        {needsRoll && (
          <button className="btn btn-primary btn-sm" onClick={onRollDice}>
            <Dice5 size={14} />
            Бросить кубики
          </button>
        )}
      </div>

      {/* Board */}
      <div style={{
        display: 'flex',
        flexDirection: 'column',
        background: 'var(--bg-tertiary)',
        borderRadius: 'var(--radius-md)',
        border: '1px solid var(--bg-glass-border)',
        overflow: 'hidden',
      }}>
        {/* Top row: points 13-24 */}
        <div style={{ display: 'flex', height: 180 }}>
          <div style={{ display: 'flex', flex: 1 }}>
            {Array.from({ length: 6 }, (_, i) => renderPoint(0, i, true))}
          </div>

          {/* Bar */}
          <div
            onClick={handleBarClick}
            style={{
              width: 40,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              gap: 4,
              background: 'var(--bg-secondary)',
              borderLeft: '1px solid var(--bg-glass-border)',
              borderRight: '1px solid var(--bg-glass-border)',
              cursor: isMyTurn && ((myColor === 'black' ? bar.white : bar.black) > 0) ? 'pointer' : 'default',
              ...(selectedBar ? { background: 'var(--gold-light-bg)' } : {}),
            }}
          >
            {bar.black > 0 && (
              <div className="text-xs font-bold text-secondary">{bar.black}</div>
            )}
          </div>

          <div style={{ display: 'flex', flex: 1 }}>
            {Array.from({ length: 6 }, (_, i) => renderPoint(1, i, true))}
          </div>

          {/* Bear-off tray */}
          <div style={{
            width: 36,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'flex-end',
            padding: '4px 2px',
            gap: 2,
            background: 'var(--bg-secondary)',
            borderLeft: '1px solid var(--bg-glass-border)',
          }}>
            <span className="text-xs text-muted">{borneOff.black}</span>
          </div>
        </div>

        {/* Bottom row: points 12-1 */}
        <div style={{ display: 'flex', height: 180 }}>
          <div style={{ display: 'flex', flex: 1 }}>
            {Array.from({ length: 6 }, (_, i) => renderPoint(3, i, false))}
          </div>

          {/* Bar (bottom) */}
          <div
            onClick={handleBarClick}
            style={{
              width: 40,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              gap: 4,
              background: 'var(--bg-secondary)',
              borderLeft: '1px solid var(--bg-glass-border)',
              borderRight: '1px solid var(--bg-glass-border)',
              cursor: isMyTurn && ((myColor === 'white' ? bar.white : bar.black) > 0) ? 'pointer' : 'default',
              ...(selectedBar ? { background: 'var(--gold-light-bg)' } : {}),
            }}
          >
            {bar.white > 0 && (
              <div className="text-xs font-bold text-secondary">{bar.white}</div>
            )}
          </div>

          <div style={{ display: 'flex', flex: 1 }}>
            {Array.from({ length: 6 }, (_, i) => renderPoint(2, i, false))}
          </div>

          {/* Bear-off tray */}
          <div style={{
            width: 36,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'flex-start',
            padding: '4px 2px',
            gap: 2,
            background: 'var(--bg-secondary)',
            borderLeft: '1px solid var(--bg-glass-border)',
          }}>
            <span className="text-xs text-muted">{borneOff.white}</span>
          </div>
        </div>
      </div>
    </div>
  )
}
