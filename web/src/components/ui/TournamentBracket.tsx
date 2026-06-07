import React from 'react'
import { players } from '../../data/players'

interface TournamentBracketProps {
  rounds: {
    round: string
    matches: { player1: string; score1: number; player2: string; score2: number; winner?: string }[]
  }[]
}

function getPlayerName(sid: string): string {
  if (sid === 'TBD') return 'TBD'
  const p = players.find((pl) => pl.sid === sid)
  return p ? p.name : sid
}

export function TournamentBracket({ rounds }: TournamentBracketProps) {
  if (!rounds || rounds.length === 0) {
    return (
      <div className="text-center py-12" style={{ color: 'var(--text-muted)' }}>
        Сетка турнира будет доступна после начала
      </div>
    )
  }

  const roundLabels = ['Четвертьфинал', 'Полуфинал', 'Финал']
  const colWidth = 200
  const rowHeight = 50
  const gap = 40

  return (
    <div style={{ overflowX: 'auto', paddingBottom: '16px' }}>
      <svg
        width={rounds.length * (colWidth + gap) + 40}
        height={Math.max(4, rounds[0]?.matches.length || 1) * rowHeight * 2 + 60}
        viewBox={`0 0 ${rounds.length * (colWidth + gap) + 40} ${Math.max(4, rounds[0]?.matches.length || 1) * rowHeight * 2 + 60}`}
      >
        {rounds.map((round, roundIdx) => {
          const matchesInRound = round.matches.length
          const totalHeight = matchesInRound * rowHeight * 2
          const startY = 30

          return round.matches.map((match, matchIdx) => {
            const x = 20 + roundIdx * (colWidth + gap)
            const y = startY + matchIdx * (totalHeight / matchesInRound)
            const p1Name = getPlayerName(match.player1)
            const p2Name = getPlayerName(match.player2)
            const isWinner1 = match.winner === match.player1
            const isWinner2 = match.winner === match.player2

            return (
              <g key={`r${roundIdx}-m${matchIdx}`}>
                {/* Connector line to next round */}
                {roundIdx < rounds.length - 1 && (
                  <line
                    x1={x + colWidth}
                    y1={y + rowHeight}
                    x2={x + colWidth + gap}
                    y2={y + rowHeight + (rowHeight * 2) / (round.matches.length * 2 || 1)}
                    stroke="rgba(212,168,67,0.3)"
                    strokeWidth="1.5"
                    strokeDasharray="4,4"
                  />
                )}
                {/* Match box */}
                <rect
                  x={x}
                  y={y}
                  width={colWidth}
                  height={rowHeight * 2}
                  rx={8}
                  fill="rgba(255,255,255,0.03)"
                  stroke="rgba(255,255,255,0.08)"
                  strokeWidth="1"
                />
                {match.winner && (
                  <rect
                    x={x}
                    y={y}
                    width={colWidth}
                    height={rowHeight * 2}
                    rx={8}
                    fill="rgba(212,168,67,0.05)"
                    stroke="rgba(212,168,67,0.2)"
                    strokeWidth="1"
                  />
                )}
                {/* Player 1 */}
                <text
                  x={x + 10}
                  y={y + 22}
                  fill={isWinner1 ? '#d4a843' : 'var(--text-primary)'}
                  fontSize="11"
                  fontWeight={isWinner1 ? 700 : 400}
                >
                  {p1Name.length > 16 ? p1Name.slice(0, 16) + '...' : p1Name}
                </text>
                <text x={x + colWidth - 25} y={y + 22} fill={isWinner1 ? '#d4a843' : 'var(--text-muted)'} fontSize="11" fontWeight="600" textAnchor="end">
                  {match.score1}
                </text>
                {/* Divider */}
                <line x1={x + 10} y1={y + rowHeight} x2={x + colWidth - 10} y2={y + rowHeight} stroke="rgba(255,255,255,0.06)" strokeWidth="0.5" />
                {/* Player 2 */}
                <text
                  x={x + 10}
                  y={y + rowHeight + 22}
                  fill={isWinner2 ? '#d4a843' : 'var(--text-primary)'}
                  fontSize="11"
                  fontWeight={isWinner2 ? 700 : 400}
                >
                  {p2Name.length > 16 ? p2Name.slice(0, 16) + '...' : p2Name}
                </text>
                <text x={x + colWidth - 25} y={y + rowHeight + 22} fill={isWinner2 ? '#d4a843' : 'var(--text-muted)'} fontSize="11" fontWeight="600" textAnchor="end">
                  {match.score2}
                </text>
              </g>
            )
          })
        })}

        {/* Round labels */}
        {roundLabels.map((label, i) => {
          if (i >= rounds.length) return null
          const x = 20 + i * (colWidth + gap)
          return (
            <text key={label} x={x + colWidth / 2} y={15} textAnchor="middle" fill="var(--text-muted)" fontSize="10" fontWeight="600">
              {rounds[i]?.round || label}
            </text>
          )
        })}
      </svg>
    </div>
  )
}
