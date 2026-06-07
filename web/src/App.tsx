import { BrowserRouter, Routes, Route, useLocation } from 'react-router-dom'
import { AnimatePresence, motion } from 'framer-motion'
import { Layout } from './components/layout/Layout'
import { Home } from './pages/Home'
import { Tournaments } from './pages/Tournaments'
import { Leaderboard } from './pages/Leaderboard'
import { MatchHistory } from './pages/MatchHistory'
import { Profile } from './pages/Profile'
import { Game } from './pages/Game'
import { ToastContainer } from './components/ui/Toast'

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
            <Route path="/" element={<Home />} />
            <Route path="/tournaments" element={<Tournaments />} />
            <Route path="/leaderboard" element={<Leaderboard />} />
            <Route path="/matches" element={<MatchHistory />} />
            <Route path="/profile" element={<Profile />} />
            <Route path="/game/:gameId" element={<Game />} />
          </Route>
        </Routes>
      </motion.div>
    </AnimatePresence>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AnimatedRoutes />
      <ToastContainer />
    </BrowserRouter>
  )
}
