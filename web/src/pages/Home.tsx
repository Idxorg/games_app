import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'
import { ArrowRight, Flame, Trophy, BarChart3, Zap } from 'lucide-react'
import { GameCard } from '../components/ui/GameCard'
import { LiveIndicator } from '../components/ui/LiveIndicator'
import { EloBar } from '../components/ui/EloBar'
import { StatsCard } from '../components/ui/StatsCard'
import { ChessPiece } from '../components/ui/ChessPiece'
import { useGameStore } from '../stores/gameStore'
import { players } from '../data/players'
import { tournaments } from '../data/tournaments'
import { matches } from '../data/matches'

function HeroSection() {
  return (
    <section className="py-16 md:py-24 relative overflow-hidden">
      {/* Background gradient */}
      <div
        style={{
          position: 'absolute',
          top: -200,
          left: '50%',
          transform: 'translateX(-50%)',
          width: 600,
          height: 600,
          borderRadius: '50%',
          background: 'radial-gradient(circle, rgba(212,168,67,0.06) 0%, transparent 70%)',
          pointerEvents: 'none',
        }}
      />

      <div className="container relative z-10">
        <div className="flex flex-col lg:flex-row items-center gap-12">
          {/* Text */}
          <div className="flex-1 text-center lg:text-left">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
            >
              <div className="flex items-center gap-2 justify-center lg:justify-start mb-4">
                <LiveIndicator />
                <span className="text-sm" style={{ color: 'var(--text-secondary)' }}>
                  {matches.filter((m) => m.result === 'win').length} активных партий
                </span>
              </div>
              <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-4" style={{ lineHeight: 1.1 }}>
                Корпоративная
                <br />
                <span style={{ color: 'var(--gold)' }}>игровая платформа</span>
              </h1>
              <p className="text-lg mb-8 max-w-lg mx-auto lg:mx-0" style={{ color: 'var(--text-secondary)' }}>
                Соревнуйтесь с коллегами в шахматы, шашки, нарды и другие игры.
                Турниры, рейтинг ELO и мгновенные матчи.
              </p>
              <div className="flex items-center gap-4 justify-center lg:justify-start">
                <Link to="/game/chess" className="btn btn-primary text-base px-6 py-3 no-underline">
                  Начать игру
                  <ArrowRight size={18} />
                </Link>
                <Link to="/tournaments" className="btn btn-secondary text-base px-6 py-3 no-underline">
                  Турниры
                </Link>
              </div>
            </motion.div>
          </div>

          {/* 3D Chess Board */}
          <motion.div
            initial={{ opacity: 0, rotateY: -15 }}
            animate={{ opacity: 1, rotateY: 0 }}
            transition={{ duration: 0.8, delay: 0.2 }}
            style={{ perspective: 800 }}
            className="flex-shrink-0"
          >
            <div className="grid grid-cols-8 gap-0" style={{ transform: 'rotateX(10deg) rotateY(5deg)', boxShadow: '0 20px 60px rgba(0,0,0,0.5), 0 0 40px rgba(212,168,67,0.1)', borderRadius: 8, overflow: 'hidden' }}>
              {Array.from({ length: 64 }, (_, i) => {
                const row = Math.floor(i / 8)
                const col = i % 8
                const isLight = (row + col) % 2 === 0
                return (
                  <motion.div
                    key={i}
                    initial={{ opacity: 0, scale: 0 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.3 + i * 0.015, duration: 0.3 }}
                    style={{
                      width: 44,
                      height: 44,
                      background: isLight ? '#1a1a2e' : 'rgba(212,168,67,0.12)',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                    }}
                  >
                    {/* Sample pieces on first/last rows */}
                    {row === 0 && (
                      <ChessPiece
                        type={['rook', 'knight', 'bishop', 'queen', 'king', 'bishop', 'knight', 'rook'][col] as any}
                        color="black"
                        size={34}
                      />
                    )}
                    {row === 1 && <ChessPiece type="pawn" color="black" size={32} />}
                    {row === 6 && <ChessPiece type="pawn" color="white" size={32} />}
                    {row === 7 && (
                      <ChessPiece
                        type={['rook', 'knight', 'bishop', 'queen', 'king', 'bishop', 'knight', 'rook'][col] as any}
                        color="white"
                        size={34}
                      />
                    )}
                  </motion.div>
                )
              })}
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  )
}

function LiveMatches() {
  const liveMatches = matches.slice(0, 3)
  return (
    <section className="py-8">
      <div className="container">
        <div className="flex items-center gap-3 mb-6">
          <Flame size={20} color="var(--danger)" />
          <h2 className="text-xl font-bold">Активные матчи</h2>
          <LiveIndicator />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
          {liveMatches.map((match, i) => {
            const p1 = players.find((p) => p.sid === match.player1Id)
            const p2 = players.find((p) => p.sid === match.player2Id)
            return (
              <motion.div
                key={match.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.1 }}
                className="glass-card p-4 flex items-center gap-3"
              >
                <div className="flex-grow">
                  <div className="flex items-center gap-2 text-sm">
                    <span className="font-semibold" style={{ color: 'var(--text-primary)' }}>{p1?.initials}</span>
                    <span style={{ color: 'var(--text-muted)' }}>vs</span>
                    <span className="font-semibold" style={{ color: 'var(--text-primary)' }}>{p2?.initials}</span>
                  </div>
                  <div className="text-xs mt-1" style={{ color: 'var(--text-muted)' }}>
                    Шахматы -- {match.duration}
                  </div>
                </div>
                <Zap size={16} color="var(--warning)" />
              </motion.div>
            )
          })}
        </div>
      </div>
    </section>
  )
}

function StatsSection() {
  return (
    <section className="py-8">
      <div className="container">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <StatsCard label="Игроков онлайн" value={247} icon={<Zap size={24} />} color="var(--success)" />
          <StatsCard label="Активных партий" value={58} icon={<Flame size={24} />} color="var(--warning)" />
          <StatsCard label="Турниров" value={12} icon={<Trophy size={24} />} />
          <StatsCard label="Матчей сыграно" value={4820} icon={<BarChart3 size={24} />} color="var(--info)" />
        </div>
      </div>
    </section>
  )
}

function GamesSection() {
  const { games } = useGameStore()
  return (
    <section className="py-8">
      <div className="container">
        <div className="flex items-center justify-between mb-6">
          <h2 className="section-title mb-0" style={{ margin: 0 }}>Игры</h2>
          <Link to="/leaderboard" className="btn btn-ghost text-sm no-underline">
            Все игры <ArrowRight size={14} />
          </Link>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {games.map((game, i) => (
            <motion.div
              key={game.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.1 }}
            >
              <GameCard game={game} />
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  )
}

function TournamentPreview() {
  const activeTournament = tournaments.find((t) => t.status === 'active')
  if (!activeTournament) return null
  return (
    <section className="py-8">
      <div className="container">
        <div className="flex items-center justify-between mb-6">
          <h2 className="section-title mb-0" style={{ margin: 0 }}>Турниры</h2>
          <Link to="/tournaments" className="btn btn-ghost text-sm no-underline">
            Все турниры <ArrowRight size={14} />
          </Link>
        </div>
        <div className="glass-card p-6 gold-border">
          <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
            <div>
              <div className="flex items-center gap-2 mb-2">
                <Trophy size={20} color="var(--gold)" />
                <span className="badge badge-gold">Активный</span>
              </div>
              <h3 className="text-xl font-bold mb-1">{activeTournament.name}</h3>
              <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
                {activeTournament.participants.length} участников -- Приз: {activeTournament.prize}
              </p>
            </div>
            <Link to="/tournaments" className="btn btn-primary no-underline">
              Смотреть сетку
            </Link>
          </div>
        </div>
      </div>
    </section>
  )
}

function TopLeaderboard() {
  const top3 = players.slice(0, 3)
  return (
    <section className="py-8 pb-16">
      <div className="container">
        <div className="flex items-center justify-between mb-6">
          <h2 className="section-title mb-0" style={{ margin: 0 }}>Лучшие игроки</h2>
          <Link to="/leaderboard" className="btn btn-ghost text-sm no-underline">
            Полный рейтинг <ArrowRight size={14} />
          </Link>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {top3.map((player, i) => (
            <motion.div
              key={player.sid}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.15 }}
              className="glass-card p-5"
              style={i === 0 ? { borderColor: 'rgba(212,168,67,0.3)', boxShadow: 'var(--shadow-gold)' } : {}}
            >
              <div className="flex items-center gap-3 mb-3">
                <div className="text-2xl font-bold" style={{ color: i === 0 ? 'var(--gold)' : i === 1 ? '#a0a0b0' : '#cd7f32', width: 32 }}>
                  #{i + 1}
                </div>
                <div
                  className="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold"
                  style={{ background: 'linear-gradient(135deg, var(--gold), var(--gold-dark))', color: '#0a0a0f' }}
                >
                  {player.initials}
                </div>
                <div>
                  <div className="font-semibold text-sm">{player.name}</div>
                  <div className="text-xs" style={{ color: 'var(--text-muted)' }}>{player.department}</div>
                </div>
              </div>
              <div className="mb-2">
                <EloBar value={player.elo.chess} max={2500} />
              </div>
              <div className="flex items-center gap-3 text-xs">
                <span style={{ color: 'var(--success)' }}>{player.wins}В</span>
                <span style={{ color: 'var(--danger)' }}>{player.losses}П</span>
                <span style={{ color: 'var(--text-muted)' }}>{player.draws}Н</span>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  )
}

export function Home() {
  return (
    <>
      <HeroSection />
      <StatsSection />
      <LiveMatches />
      <GamesSection />
      <TournamentPreview />
      <TopLeaderboard />
    </>
  )
}
