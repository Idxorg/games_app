import { useEffect, useRef } from 'react'

// ─── Types ──────────────────────────────────────────────────────────────────

interface MoveEntry {
  from: string | number | [number, number]
  to: string | number | [number, number]
  color: string
}

interface MoveHistoryProps {
  moves: MoveEntry[]
}

// ─── Helpers ────────────────────────────────────────────────────────────────

function formatMove(m: MoveEntry, index: number): string {
  const fromStr = Array.isArray(m.from)
    ? `${m.from[0]},${m.from[1]}`
    : String(m.from)
  const toStr = Array.isArray(m.to)
    ? `${m.to[0]},${m.to[1]}`
    : String(m.to)
  return `${fromStr}-${toStr}`
}

// ─── Component ──────────────────────────────────────────────────────────────

export function MoveHistory({ moves }: MoveHistoryProps) {
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [moves.length])

  const pairs: { num: number; white?: MoveEntry; black?: MoveEntry }[] = []
  for (let i = 0; i < moves.length; i += 2) {
    pairs.push({
      num: Math.floor(i / 2) + 1,
      white: moves[i],
      black: moves[i + 1],
    })
  }

  return (
    <div className="glass-card flex flex-col" style={{ height: 300, maxHeight: '60vh' }}>
      <div className="flex items-center justify-between px-3 py-2" style={{
        borderBottom: '1px solid var(--bg-glass-border)',
      }}>
        <span className="text-xs font-semibold text-secondary">Ходы</span>
        <span className="text-xs text-muted">{moves.length} ходов</span>
      </div>

      <div style={{ flex: 1, overflowY: 'auto', padding: '4px 0' }}>
        {pairs.length === 0 && (
          <div className="flex items-center justify-center" style={{ height: '100%' }}>
            <span className="text-sm text-muted">Нет ходов</span>
          </div>
        )}
        {pairs.map((pair) => (
          <div
            key={pair.num}
            className="flex items-center"
            style={{ padding: '2px 8px', fontSize: 13 }}
          >
            <span className="font-mono text-muted" style={{ width: 32, flexShrink: 0 }}>
              {pair.num}.
            </span>
            <span className="font-mono flex-1 truncate" style={{ color: 'var(--text-primary)' }}>
              {pair.white ? formatMove(pair.white, 0) : ''}
            </span>
            <span className="font-mono flex-1 truncate" style={{ color: 'var(--text-secondary)' }}>
              {pair.black ? formatMove(pair.black, 1) : ''}
            </span>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>
    </div>
  )
}
