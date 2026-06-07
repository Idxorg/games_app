import { useEffect } from 'react'
import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'
import { ArrowRight, Flame, Trophy, BarChart3, Zap, RefreshCw } from 'lucide-react'
import { GameCard } from '../components/ui/GameCard'
import { LiveIndicator } from '../components/ui/LiveIndicator'
import { EloBar } from '../components/ui/EloBar'
import { StatsCard } from '../components/ui/StatsCard'
import { ChessPiece } from '../components/ui/ChessPiece'
import { useGameStore } from '../stores/gameStore'
import { useTournamentStore } from '../stores/tournamentStore'
import { useLeaderboardStore } from '../stores/leaderboardStore'
import { useMatchStore } from '../stores/matchStore'
import { useUserStore } from '../stores/userStore'

const gameNames: Record<string, string> = {
  chess: 'Шахматы',
  checkers: 'Шашки',
  backgammon: 'Нарды',
  trivia: 'Викторины',
}

function HeroSection({ matchCount }: { matchCount: number | null }) {
  return (
    <section className="py-16 md:py-24 relative overflow-hidden">
      {/* Background gradient */}
      <div className="hero-bg-glow" />

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
                <span className="text-sm text-secondary">
                  {matchCount != null ? `${matchCount} активных партий` : '—'}
                </span>
              </div>
              <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-4" style={{ lineHeight: 1.1 }}>
                Корпоративная
                <br />
                <span className="gold-gradient">игровая платформа</span>
              </h1>
              <p className="text-lg mb-8 max-w-lg mx-auto lg:mx-0 text-secondary">
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

          {/* Chess Board */}
          <motion.div
            initial={{ opacity: 0, rotateY: -15 }}
            animate={{ opacity: 1, rotateY: 0 }}
            transition={{ duration: 0.8, delay: 0.2 }}
            className="chess-board-wrapper flex-shrink-0"
          >
            <div className="chess-board">
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
                    className={`chess-square ${isLight ? 'chess-square-light' : 'chess-square-dark'}`}
                  >
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
  const { matches, loading, error, fetchMatches } = useMatchStore()

  useEffect(() => {
    fetchMatches()
  }, [fetchMatches])

  const liveMatches = matches.filter((m) => m.status === 'in_progress').slice(0, 3)

  return (
    <section className="py-8">
      <div className="container">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="flex items-center gap-3 mb-6"
        >
          <Flame size={20} color="var(--danger)" />
          <h2 className="text-xl font-bold">Активные матчи</h2>
          <LiveIndicator />
        </motion.div>

        {loading && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="glass-card p-4">
                <div className="flex items-center gap-3">
                  <div className="skeleton skeleton-circle" style={{ width: 32, height: 32 }} />
                  <div className="flex-grow">
                    <div className="skeleton skeleton-text" style={{ width: '60%' }} />
                    <div className="skeleton skeleton-text" style={{ width: '40%' }} />
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {error && (
          <div className="glass-card p-8 text-center">
            <p className="text-muted mb-3">{error}</p>
            <button className="btn btn-secondary" onClick={() => fetchMatches()}>
              <RefreshCw size={14} /> Повторить
            </button>
          </div>
        )}

        {!loading && !error && liveMatches.length === 0 && (
          <div className="glass-card p-8 text-center">
            <p className="text-muted">Нет активных матчей прямо сейчас</p>
          </div>
        )}

        {!loading && !error && liveMatches.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            {liveMatches.map((match, i) => (
              <motion.div
                key={match.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.1 }}
                className="glass-card p-4 flex items-center gap-3"
              >
                <div className="flex-grow">
                  <div className="flex items-center gap-2 text-sm">
                    <span className="font-semibold text-primary">{match.player1Sid.slice(0, 6)}</span>
                    <span className="text-muted">vs</span>
                    <span className="font-semibold text-primary">{match.player2Sid.slice(0, 6)}</span>
                  </div>
                  <div className="text-xs mt-1 text-muted">
                    {gameNames[match.gameType] || match.gameType}
                  </div>
                </div>
                <Zap size={16} color="var(--warning)" />
              </motion.div>
            ))}
          </div>
        )}
      </div>
    </section>
  )
}

function StatsSection() {
  const { matches } = useMatchStore()
  const { tournaments } = useTournamentStore()

  const totalMatches = matches.length
  const activeTournaments = tournaments.filter((t) => t.status === 'active').length

  return (
    <section className="py-8">
      <div className="container">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {[
            { label: 'Игроков онлайн', value: 0, icon: <Zap size={24} />, color: 'var(--success)', placeholder: true },
            { label: 'Активных партий', value: totalMatches, icon: <Flame size={24} />, color: 'var(--warning)', placeholder: totalMatches === 0 },
            { label: 'Турниров', value: activeTournaments, icon: <Trophy size={24} />, placeholder: activeTournaments === 0 },
            { label: 'Матчей сыграно', value: totalMatches, icon: <BarChart3 size={24} />, color: 'var(--info)', placeholder: totalMatches === 0 },
          ].map((stat, i) => (
            <motion.div
              key={stat.label}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.1 }}
            >
              {stat.placeholder ? (
                <div className="glass-card p-6 flex flex-col items-center text-center">
                  <div className="mb-3" style={{ color: stat.color }}>{stat.icon}</div>
                  <div className="text-3xl font-bold font-mono" style={{ color: 'var(--text-muted)' }}>—</div>
                  <div className="text-xs mt-1" style={{ color: 'var(--text-muted)' }}>{stat.label}</div>
                </div>
              ) : (
                <StatsCard label={stat.label} value={stat.value} icon={stat.icon} color={stat.color} />
              )}
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  )
}

function GamesSection() {
  const { games, loading, error, fetchGames } = useGameStore()

  useEffect(() => {
    fetchGames()
  }, [fetchGames])

  if (loading) {
    return (
      <section className="py-8">
        <div className="container">
          <h2 className="section-title mb-6">Игры</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="glass-card p-6">
                <div className="skeleton skeleton-circle mb-4" style={{ width: 48, height: 48 }} />
                <div className="skeleton skeleton-title mb-2" />
                <div className="skeleton skeleton-text" style={{ width: '80%' }} />
              </div>
            ))}
          </div>
        </div>
      </section>
    )
  }

  if (error) {
    return (
      <section className="py-8">
        <div className="container">
          <h2 className="section-title mb-6">Игры</h2>
          <div className="glass-card p-8 text-center">
            <p className="text-muted mb-3">{error}</p>
            <button className="btn btn-secondary" onClick={fetchGames}>
              <RefreshCw size={14} /> Повторить
            </button>
          </div>
        </div>
      </section>
    )
  }

  if (games.length === 0) {
    return (
      <section className="py-8">
        <div className="container">
          <h2 className="section-title mb-6">Игры</h2>
          <div className="glass-card p-8 text-center">
            <p className="text-muted">Список игр пуст</p>
          </div>
        </div>
      </section>
    )
  }

  return (
    <section className="py-8">
      <div className="container">
        <div className="flex items-center justify-between mb-6">
          <h2 className="section-title mb-0">Игры</h2>
          <Link to="/games" className="btn btn-ghost text-sm no-underline">
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
  const { tournaments, loading, error, fetchTournaments } = useTournamentStore()

  useEffect(() => {
    fetchTournaments()
  }, [fetchTournaments])

  const activeTournament = tournaments.find((t) => t.status === 'active')

  if (loading) {
    return (
      <section className="py-8">
        <div className="container">
          <div className="glass-card p-6">
            <div className="skeleton skeleton-title mb-3" style={{ width: '40%' }} />
            <div className="skeleton skeleton-text" style={{ width: '60%' }} />
          </div>
        </div>
      </section>
    )
  }

  if (!activeTournament) {
    return (
      <section className="py-8">
        <div className="container">
          <div className="flex items-center justify-between mb-6">
            <h2 className="section-title mb-0">Турниры</h2>
            <Link to="/tournaments" className="btn btn-ghost text-sm no-underline">
              Все турниры <ArrowRight size={14} />
            </Link>
          </div>
          <div className="glass-card p-8 text-center">
            <p className="text-muted">Нет активных турниров</p>
          </div>
        </div>
      </section>
    )
  }

  return (
    <section className="py-8">
      <div className="container">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <div className="flex items-center justify-between mb-6">
            <h2 className="section-title mb-0">Турниры</h2>
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
                <p className="text-sm text-secondary">
                  {activeTournament.prize || 'Нет приза'}
                </p>
              </div>
              <Link to="/tournaments" className="btn btn-primary no-underline">
                Смотреть сетку
              </Link>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  )
}

function TopLeaderboard() {
  const { players, loading, error, fetchLeaderboard } = useLeaderboardStore()

  useEffect(() => {
    fetchLeaderboard('chess')
  }, [fetchLeaderboard])

  const top3 = players.slice(0, 3)
  const rankColors = ['var(--gold)', '#a0a0b0', '#cd7f32']

  if (loading) {
    return (
      <section className="py-8 pb-16">
        <div className="container">
          <h2 className="section-title mb-6">Лучшие игроки</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="glass-card p-5">
                <div className="flex items-center gap-3 mb-3">
                  <div className="skeleton skeleton-circle" style={{ width: 32, height: 32 }} />
                  <div className="skeleton skeleton-circle" style={{ width: 40, height: 40 }} />
                  <div className="skeleton skeleton-text" style={{ width: '50%' }} />
                </div>
                <div className="skeleton skeleton-text" />
              </div>
            ))}
          </div>
        </div>
      </section>
    )
  }

  if (error) {
    return (
      <section className="py-8 pb-16">
        <div className="container">
          <h2 className="section-title mb-6">Лучшие игроки</h2>
          <div className="glass-card p-8 text-center">
            <p className="text-muted mb-3">{error}</p>
            <button className="btn btn-secondary" onClick={() => fetchLeaderboard('chess')}>
              <RefreshCw size={14} /> Повторить
            </button>
          </div>
        </div>
      </section>
    )
  }

  if (top3.length === 0) {
    return (
      <section className="py-8 pb-16">
        <div className="container">
          <div className="flex items-center justify-between mb-6">
            <h2 className="section-title mb-0">Лучшие игроки</h2>
            <Link to="/leaderboard" className="btn btn-ghost text-sm no-underline">
              Полный рейтинг <ArrowRight size={14} />
            </Link>
          </div>
          <div className="glass-card p-8 text-center">
            <p className="text-muted">Рейтинг пока пуст</p>
          </div>
        </div>
      </section>
    )
  }

  return (
    <motion.section
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.4 }}
      className="py-8 pb-16"
    >
      <div className="container">
        <div className="flex items-center justify-between mb-6">
          <h2 className="section-title mb-0">Лучшие игроки</h2>
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
              className={`glass-card p-5 ${i === 0 ? 'gold-border gold-glow' : ''}`}
            >
              <div className="flex items-center gap-3 mb-3">
                <div className="text-2xl font-bold" style={{ color: rankColors[i], width: 32 }}>
                  #{i + 1}
                </div>
                <div className="profile-avatar-md">
                  {player.initials}
                </div>
                <div>
                  <div className="font-semibold text-sm">{player.name}</div>
                  <div className="text-xs text-muted">{player.department}</div>
                </div>
              </div>
              {player.elo.chess != null && (
                <div className="mb-2">
                  <EloBar value={player.elo.chess} max={2500} />
                </div>
              )}
              <div className="flex items-center gap-3 text-xs">
                <span className="text-success">{player.wins}В</span>
                <span className="text-danger">{player.losses}П</span>
                <span className="text-muted">{player.draws}Н</span>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </motion.section>
  )
}

export function Home() {
  const isAuthenticated = useUserStore((s) => s.isAuthenticated)
  const isEmbed = useUserStore((s) => s.isEmbed)

  if (isEmbed && !isAuthenticated) {
    return (
      <div className="py-8">
        <div className="container text-center">
          <p className="text-muted">Для доступа войдите через портал</p>
        </div>
      </div>
    )
  }

  return (
    <>
      <HeroSection matchCount={null} />
      <StatsSection />
      <LiveMatches />
      <GamesSection />
      <TournamentPreview />
      <TopLeaderboard />
    </>
  )
}
