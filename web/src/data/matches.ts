export interface Match {
  id: string
  player1Id: string
  player2Id: string
  game: string
  date: string
  result: 'win' | 'loss' | 'draw'
  eloChange: number
  duration: string
}

export const matches: Match[] = [
  { id: 'm1', player1Id: 'p1', player2Id: 'p2', game: 'chess', date: '2026-06-07', result: 'win', eloChange: 12, duration: '32 мин' },
  { id: 'm2', player1Id: 'p1', player2Id: 'p5', game: 'chess', date: '2026-06-06', result: 'win', eloChange: 8, duration: '28 мин' },
  { id: 'm3', player1Id: 'p3', player2Id: 'p1', game: 'checkers', date: '2026-06-05', result: 'loss', eloChange: -15, duration: '18 мин' },
  { id: 'm4', player1Id: 'p1', player2Id: 'p11', game: 'chess', date: '2026-06-04', result: 'draw', eloChange: 0, duration: '45 мин' },
  { id: 'm5', player1Id: 'p7', player2Id: 'p1', game: 'backgammon', date: '2026-06-03', result: 'win', eloChange: 18, duration: '25 мин' },
  { id: 'm6', player1Id: 'p1', player2Id: 'p9', game: 'chess', date: '2026-06-02', result: 'win', eloChange: 6, duration: '38 мин' },
  { id: 'm7', player1Id: 'p4', player2Id: 'p1', game: 'checkers', date: '2026-06-01', result: 'loss', eloChange: -10, duration: '22 мин' },
  { id: 'm8', player1Id: 'p1', player2Id: 'p2', game: 'backgammon', date: '2026-05-31', result: 'win', eloChange: 14, duration: '30 мин' },
  { id: 'm9', player1Id: 'p1', player2Id: 'p3', game: 'chess', date: '2026-05-30', result: 'draw', eloChange: -2, duration: '50 мин' },
  { id: 'm10', player1Id: 'p11', player2Id: 'p1', game: 'chess', date: '2026-05-29', result: 'win', eloChange: 20, duration: '35 мин' },
]
