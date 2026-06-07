import { useEffect } from 'react'
import { BrowserRouter, Routes, Route, useLocation, useNavigate, Link } from 'react-router-dom'
import { AnimatePresence, motion } from 'framer-motion'
import { Home, ArrowLeft } from 'lucide-react'
import { Layout } from './components/layout/Layout'
import { Home as HomePage } from './pages/Home'
import { Tournaments } from './pages/Tournaments'
import { Leaderboard } from './pages/Leaderboard'
import { MatchHistory } from './pages/MatchHistory'
import { Profile } from './pages/Profile'
import { Game } from './pages/Game'
import { ToastContainer } from './components/ui/Toast'
import { useUserStore } from './stores/userStore'
import { useInviteStore } from './stores/inviteStore'
import { isEmbedMode, getHandoffTheme, onInviteAccept, onInviteDecline } from './embedHandoff'

function NotFound() {
  return (
    <div className="not-found">
      <div className="not-found-code">404</div>
      <h1 className="text-xl font-bold mb-2">Страница не найдена</h1>
      <p className="text-secondary mb-6">
        Запрашиваемая страница не существует или была перемещена
      </p>
      <Link to="/" className="btn btn-primary">
        <Home size={16} />
        На главную
      </Link>
    </div>
  )
}

function AnimatedRoutes() {
  const location = useLocation()

  return (
    <AnimatePresence mode="wait">
      <motion.div
        key={location.pathname}
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: -10 }}
        transition={{ duration: 0.2 }}
      >
        <Routes location={location}>
          <Route element={<Layout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/tournaments" element={<Tournaments />} />
            <Route path="/leaderboard" element={<Leaderboard />} />
            <Route path="/matches" element={<MatchHistory />} />
            <Route path="/profile" element={<Profile />} />
            <Route path="/game/:gameId" element={<Game />} />
            <Route path="/game/:gameId/:matchId" element={<Game />} />
            <Route path="*" element={<NotFound />} />
          </Route>
        </Routes>
      </motion.div>
    </AnimatePresence>
  )
}

/** Wire postMessage invite handlers inside Router context (needs useNavigate) */
function InviteHandler() {
  const navigate = useNavigate()

  useEffect(() => {
    onInviteAccept(async (inviteId: string) => {
      const store = useInviteStore.getState()
      // Look up game_type from pending invites before accepting
      const invite = store.pendingInvites.find((i) => i.id === inviteId)
      const gameType = invite?.game_type || 'chess'
      try {
        await store.accept(inviteId)
        // Navigate to the game lobby after accept; once a matchId is available
        // the Game component will auto-connect via WebSocket
        navigate(`/game/${gameType}`)
      } catch (err) {
        console.error('Failed to accept invite:', err)
      }
    })

    onInviteDecline(async (inviteId: string) => {
      try {
        await useInviteStore.getState().decline(inviteId)
      } catch (err) {
        console.error('Failed to decline invite:', err)
      }
    })
  }, [navigate])

  return null
}

export default function App() {
  const initialize = useUserStore((s) => s.initialize)
  const loading = useUserStore((s) => s.loading)
  const theme = useUserStore((s) => s.theme)

  useEffect(() => {
    initialize()
  }, [initialize])

  // Apply theme class + data-theme attribute to document
  useEffect(() => {
    // Prefer handoff theme from postMessage, fall back to userStore theme
    const effectiveTheme = getHandoffTheme() || theme
    if (effectiveTheme === 'light') {
      document.documentElement.classList.add('theme-light')
      document.documentElement.classList.remove('theme-dark')
      document.documentElement.setAttribute('data-theme', 'light')
    } else if (effectiveTheme === 'dark') {
      document.documentElement.classList.add('theme-dark')
      document.documentElement.classList.remove('theme-light')
      document.documentElement.setAttribute('data-theme', 'dark')
    } else {
      document.documentElement.removeAttribute('data-theme')
    }
  }, [theme])

  if (loading) {
    return (
      <div className="loading-screen">
        <div className="skeleton w-48 h-8 mb-4" />
        <div className="skeleton w-64 h-6" />
      </div>
    )
  }

  return (
    <BrowserRouter basename={isEmbedMode() ? '/games' : undefined}>
      <InviteHandler />
      <AnimatedRoutes />
      <ToastContainer />
    </BrowserRouter>
  )
}
