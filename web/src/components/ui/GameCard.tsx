import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'
import { Users, Crown } from 'lucide-react'
import { ChessPiece } from './ChessPiece'
import { CheckersPiece } from './CheckersPiece'
import { BackgammonPiece } from './BackgammonPiece'
import { Brain, Grid3X3, Puzzle, Gamepad2, Dice5, Swords } from 'lucide-react'
import type { Game } from '../../data/games'
import { LiveIndicator } from './LiveIndicator'

interface GameCardProps {
  game: Game
}

function GameIcon({ type, size = 48 }: { type: string; size?: number }) {
  switch (type) {
    case 'chess':
      return <ChessPiece type="knight" color="white" size={size} />
    case 'checkers':
      return <CheckersPiece color="red" size={size} />
    case 'backgammon':
      return <BackgammonPiece size={size * 2} />
    case 'brain':
      return <Brain size={size} color="var(--gold)" />
    case 'grid':
      return <Grid3X3 size={size} color="var(--gold)" />
    case 'puzzle':
      return <Puzzle size={size} color="var(--gold)" />
    case 'cards':
      return <Swords size={size} color="var(--gold)" />
    default:
      return <Gamepad2 size={size} color="var(--gold)" />
  }
}

export function GameCard({ game }: GameCardProps) {
  return (
    <motion.div
      whileHover={{ y: -8, rotateX: 2, rotateY: -2 }}
      transition={{ type: 'spring', stiffness: 300, damping: 20 }}
      style={{ perspective: 1000 }}
    >
      <Link to={game.route} className="block">
        <div className="glass-card p-6 h-full flex flex-col cursor-pointer">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center justify-center w-16 h-16 rounded-xl" style={{ background: 'rgba(212,168,67,0.08)' }}>
              <GameIcon type={game.iconType} />
            </div>
            {game.isLive && <LiveIndicator />}
          </div>
          <h3 className="text-lg font-bold mb-2" style={{ color: 'var(--text-primary)' }}>{game.name}</h3>
          <p className="text-sm mb-4 flex-grow" style={{ color: 'var(--text-secondary)' }}>{game.description}</p>
          <div className="flex items-center justify-between pt-3" style={{ borderTop: '1px solid var(--bg-glass-border)' }}>
            <div className="flex items-center gap-1.5 text-xs" style={{ color: 'var(--text-muted)' }}>
              <Users size={14} />
              <span>{game.playerCount}</span>
            </div>
            <span className="badge badge-gold">
              <Crown size={10} />
              {game.category}
            </span>
          </div>
        </div>
      </Link>
    </motion.div>
  )
}
