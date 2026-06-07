export interface TournamentMatch {
  player1: string
  score1: number
  player2: string
  score2: number
  winner?: string
}

export interface Tournament {
  id: string
  name: string
  game: string
  status: 'active' | 'upcoming' | 'completed'
  startDate: string
  endDate: string
  participants: string[]
  rounds: {
    round: string
    matches: TournamentMatch[]
  }[]
  prize: string
}

export const tournaments: Tournament[] = [
  {
    id: 't1',
    name: 'Кубок компании по шахматам',
    game: 'chess',
    status: 'active',
    startDate: '2026-06-01',
    endDate: '2026-06-30',
    participants: ['p1', 'p2', 'p3', 'p5', 'p7', 'p11', 'p9', 'p4'],
    prize: 'Премиум подписка на обучение',
    rounds: [
      { round: 'Четвертьфинал', matches: [
        { player1: 'p1', score1: 2, player2: 'p9', score2: 1, winner: 'p1' },
        { player1: 'p2', score1: 2, player2: 'p4', score2: 0, winner: 'p2' },
        { player1: 'p5', score1: 1, player2: 'p7', score2: 2, winner: 'p7' },
        { player1: 'p11', score1: 2, player2: 'p3', score2: 1, winner: 'p11' },
      ]},
      { round: 'Полуфинал', matches: [
        { player1: 'p1', score1: 0, player2: 'p2', score2: 0, winner: undefined },
        { player1: 'p7', score1: 0, player2: 'p11', score2: 0, winner: undefined },
      ]},
      { round: 'Финал', matches: [
        { player1: 'TBD', score1: 0, player2: 'TBD', score2: 0, winner: undefined },
      ]},
    ],
  },
  {
    id: 't2',
    name: 'Турнир по шашкам',
    game: 'checkers',
    status: 'upcoming',
    startDate: '2026-07-15',
    endDate: '2026-07-30',
    participants: ['p1', 'p3', 'p6', 'p8', 'p10', 'p12', 'p9', 'p5'],
    prize: 'Сертификат в магазин',
    rounds: [
      { round: 'Четвертьфинал', matches: [
        { player1: 'p1', score1: 0, player2: 'p9', score2: 0, winner: undefined },
        { player1: 'p3', score1: 0, player2: 'p8', score2: 0, winner: undefined },
        { player1: 'p6', score1: 0, player2: 'p5', score2: 0, winner: undefined },
        { player1: 'p10', score1: 0, player2: 'p12', score2: 0, winner: undefined },
      ]},
      { round: 'Полуфинал', matches: [
        { player1: 'TBD', score1: 0, player2: 'TBD', score2: 0, winner: undefined },
        { player1: 'TBD', score1: 0, player2: 'TBD', score2: 0, winner: undefined },
      ]},
      { round: 'Финал', matches: [
        { player1: 'TBD', score1: 0, player2: 'TBD', score2: 0, winner: undefined },
      ]},
    ],
  },
  {
    id: 't3',
    name: 'Первенство по нардам',
    game: 'backgammon',
    status: 'completed',
    startDate: '2026-05-01',
    endDate: '2026-05-20',
    participants: ['p1', 'p5', 'p7', 'p9', 'p3', 'p11', 'p6', 'p12'],
    prize: 'Подарочный набор настольных игр',
    rounds: [
      { round: 'Четвертьфинал', matches: [
        { player1: 'p1', score1: 3, player2: 'p9', score2: 2, winner: 'p1' },
        { player1: 'p5', score1: 3, player2: 'p3', score2: 1, winner: 'p5' },
        { player1: 'p7', score1: 3, player2: 'p6', score2: 2, winner: 'p7' },
        { player1: 'p11', score1: 3, player2: 'p12', score2: 0, winner: 'p11' },
      ]},
      { round: 'Полуфинал', matches: [
        { player1: 'p1', score1: 3, player2: 'p5', score2: 2, winner: 'p1' },
        { player1: 'p7', score1: 2, player2: 'p11', score2: 3, winner: 'p11' },
      ]},
      { round: 'Финал', matches: [
        { player1: 'p1', score1: 3, player2: 'p11', score2: 1, winner: 'p1' },
      ]},
    ],
  },
  {
    id: 't4',
    name: 'Лига викторин',
    game: 'trivia',
    status: 'active',
    startDate: '2026-06-01',
    endDate: '2026-06-15',
    participants: ['p2', 'p4', 'p6', 'p8', 'p10', 'p12'],
    prize: 'Корпоративный мерч',
    rounds: [
      { round: 'Группа A', matches: [
        { player1: 'p2', score1: 850, player2: 'p4', score2: 720, winner: 'p2' },
        { player1: 'p6', score1: 0, player2: 'p8', score2: 0, winner: undefined },
      ]},
      { round: 'Группа B', matches: [
        { player1: 'p10', score1: 0, player2: 'p12', score2: 0, winner: undefined },
      ]},
      { round: 'Финал', matches: [
        { player1: 'TBD', score1: 0, player2: 'TBD', score2: 0, winner: undefined },
      ]},
    ],
  },
  {
    id: 't5',
    name: 'Осенний турнир по шахматам',
    game: 'chess',
    status: 'upcoming',
    startDate: '2026-09-01',
    endDate: '2026-09-30',
    participants: [],
    prize: 'Главный приз сезона',
    rounds: [],
  },
  {
    id: 't6',
    name: 'Рапид-чемпионат',
    game: 'chess',
    status: 'completed',
    startDate: '2026-04-10',
    endDate: '2026-04-12',
    participants: ['p1', 'p2', 'p3', 'p5', 'p7', 'p11'],
    prize: 'Медаль чемпиона',
    rounds: [
      { round: 'Полуфинал', matches: [
        { player1: 'p1', score1: 2, player2: 'p5', score2: 1, winner: 'p1' },
        { player1: 'p11', score1: 2, player2: 'p2', score2: 0, winner: 'p11' },
      ]},
      { round: 'Финал', matches: [
        { player1: 'p1', score1: 3, player2: 'p11', score2: 2, winner: 'p1' },
      ]},
    ],
  },
]
