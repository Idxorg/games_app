export interface Player {
  sid: string
  name: string
  department: string
  initials: string
  elo: Record<string, number>
  wins: number
  losses: number
  draws: number
  trend: 'up' | 'down' | 'same'
}

export const players: Player[] = [
  { sid: 'p1', name: 'Алексей Петров', department: 'Разработка', initials: 'АП', elo: { chess: 2150, checkers: 1800, backgammon: 1650, trivia: 1200 }, wins: 142, losses: 58, draws: 30, trend: 'up' },
  { sid: 'p2', name: 'Мария Иванова', department: 'Дизайн', initials: 'МИ', elo: { chess: 2080, checkers: 1750, backgammon: 1720, trivia: 1400 }, wins: 128, losses: 72, draws: 25, trend: 'up' },
  { sid: 'p3', name: 'Дмитрий Козлов', department: 'Маркетинг', initials: 'ДК', elo: { chess: 1950, checkers: 1900, backgammon: 1580, trivia: 1300 }, wins: 110, losses: 85, draws: 20, trend: 'same' },
  { sid: 'p4', name: 'Елена Смирнова', department: 'HR', initials: 'ЕС', elo: { chess: 1870, checkers: 1650, backgammon: 1490, trivia: 1500 }, wins: 95, losses: 90, draws: 35, trend: 'down' },
  { sid: 'p5', name: 'Сергей Волков', department: 'Разработка', initials: 'СВ', elo: { chess: 2020, checkers: 1700, backgammon: 1800, trivia: 1100 }, wins: 135, losses: 65, draws: 28, trend: 'up' },
  { sid: 'p6', name: 'Анна Новикова', department: 'Финансы', initials: 'АН', elo: { chess: 1780, checkers: 1850, backgammon: 1600, trivia: 1350 }, wins: 88, losses: 95, draws: 22, trend: 'same' },
  { sid: 'p7', name: 'Игорь Морозов', department: 'Разработка', initials: 'ИМ', elo: { chess: 1900, checkers: 1600, backgammon: 1750, trivia: 1250 }, wins: 102, losses: 78, draws: 30, trend: 'down' },
  { sid: 'p8', name: 'Ольга Кузнецова', department: 'Продажи', initials: 'ОК', elo: { chess: 1720, checkers: 1780, backgammon: 1550, trivia: 1450 }, wins: 80, losses: 100, draws: 18, trend: 'down' },
  { sid: 'p9', name: 'Павел Соколов', department: 'Поддержка', initials: 'ПС', elo: { chess: 1850, checkers: 1720, backgammon: 1680, trivia: 1150 }, wins: 98, losses: 82, draws: 25, trend: 'same' },
  { sid: 'p10', name: 'Наталья Попова', department: 'Дизайн', initials: 'НП', elo: { chess: 1680, checkers: 1820, backgammon: 1450, trivia: 1550 }, wins: 75, losses: 105, draws: 20, trend: 'up' },
  { sid: 'p11', name: 'Андрей Лебедев', department: 'Разработка', initials: 'АЛ', elo: { chess: 2050, checkers: 1680, backgammon: 1850, trivia: 1050 }, wins: 140, losses: 60, draws: 22, trend: 'up' },
  { sid: 'p12', name: 'Татьяна Егорова', department: 'Финансы', initials: 'ТЕ', elo: { chess: 1750, checkers: 1750, backgammon: 1520, trivia: 1280 }, wins: 85, losses: 92, draws: 28, trend: 'same' },
]
