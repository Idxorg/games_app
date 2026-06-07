export interface Game {
  id: string
  name: string
  description: string
  iconType: string
  playerCount: string
  category: string
  isLive: boolean
  route: string
}

export const games: Game[] = [
  {
    id: 'chess',
    name: 'Шахматы',
    description: 'Классические шахматы с рейтингом ELO и анализом партий',
    iconType: 'chess',
    playerCount: '1 на 1',
    category: 'Стратегия',
    isLive: true,
    route: '/game/chess',
  },
  {
    id: 'checkers',
    name: 'Шашки',
    description: 'Русские шашки с турнирной поддержкой',
    iconType: 'checkers',
    playerCount: '1 на 1',
    category: 'Настольные',
    isLive: true,
    route: '/game/checkers',
  },
  {
    id: 'backgammon',
    name: 'Нарды',
    description: 'Короткие и длинные нарды онлайн',
    iconType: 'backgammon',
    playerCount: '1 на 1',
    category: 'Настольные',
    isLive: false,
    route: '/game/backgammon',
  },
  {
    id: 'tic-tac-toe',
    name: 'Крестики-нолики',
    description: 'Быстрая игра на перерыве',
    iconType: 'grid',
    playerCount: '1 на 1',
    category: 'Казуальные',
    isLive: true,
    route: '/game/tic-tac-toe',
  },
  {
    id: 'puzzle',
    name: 'Головоломки',
    description: 'Судоку, кроссворды и другие задачи',
    iconType: 'puzzle',
    playerCount: '1 игрок',
    category: 'Казуальные',
    isLive: false,
    route: '/game/puzzle',
  },
  {
    id: 'trivia',
    name: 'Викторины',
    description: 'Многопользовательские викторины на знания',
    iconType: 'brain',
    playerCount: '2-8 игроков',
    category: 'Многопользовательские',
    isLive: true,
    route: '/game/trivia',
  },
  {
    id: 'memory',
    name: 'Мемори',
    description: 'Тренируйте память, находя пары карт',
    iconType: 'cards',
    playerCount: '1 на 1',
    category: 'Казуальные',
    isLive: false,
    route: '/game/memory',
  },
]
