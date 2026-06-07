import { useState, useCallback } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { Search, Bell, ChevronDown, Menu, X, Gamepad2, Home, Trophy, BarChart3, Clock, User } from 'lucide-react'

const navItems = [
  { path: '/', label: 'Главная', icon: Home },
  { path: '/tournaments', label: 'Турниры', icon: Trophy },
  { path: '/leaderboard', label: 'Рейтинг', icon: BarChart3 },
  { path: '/matches', label: 'Матчи', icon: Clock },
  { path: '/profile', label: 'Профиль', icon: User },
]

export function Header() {
  const [searchOpen, setSearchOpen] = useState(false)
  const [searchValue, setSearchValue] = useState('')
  const [menuOpen, setMenuOpen] = useState(false)
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [notifications] = useState(3)
  const location = useLocation()
  const navigate = useNavigate()

  const handleSearch = useCallback((value: string) => {
    setSearchValue(value)
  }, [])

  return (
    <header
      className="glass"
      style={{
        position: 'sticky',
        top: 0,
        zIndex: 100,
        borderBottom: '1px solid var(--bg-glass-border)',
      }}
    >
      <div className="container flex items-center justify-between h-16">
        {/* Logo */}
        <Link to="/" className="flex items-center gap-3 no-underline">
          <Gamepad2 size={28} color="var(--gold)" />
          <span className="text-lg font-bold" style={{ color: 'var(--gold)' }}>
            Game Portal
          </span>
        </Link>

        {/* Desktop Nav */}
        <nav className="hidden md:flex items-center gap-1">
          {navItems.map((item) => {
            const isActive = location.pathname === item.path
            return (
              <Link
                key={item.path}
                to={item.path}
                className="flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium no-underline transition-all"
                style={{
                  color: isActive ? 'var(--gold)' : 'var(--text-secondary)',
                  background: isActive ? 'rgba(212,168,67,0.08)' : 'transparent',
                }}
              >
                <item.icon size={16} />
                {item.label}
              </Link>
            )
          })}
        </nav>

        {/* Actions */}
        <div className="flex items-center gap-2">
          {/* Search */}
          <AnimatePresence>
            {searchOpen && (
              <motion.div
                initial={{ width: 0, opacity: 0 }}
                animate={{ width: 240, opacity: 1 }}
                exit={{ width: 0, opacity: 0 }}
                className="overflow-hidden"
              >
                <input
                  type="text"
                  placeholder="Поиск игр..."
                  value={searchValue}
                  onChange={(e) => handleSearch(e.target.value)}
                  className="w-full px-3 py-1.5 rounded-lg text-sm"
                  style={{
                    background: 'var(--bg-tertiary)',
                    border: '1px solid var(--bg-glass-border)',
                    color: 'var(--text-primary)',
                    outline: 'none',
                  }}
                  autoFocus
                />
              </motion.div>
            )}
          </AnimatePresence>
          <button
            onClick={() => setSearchOpen(!searchOpen)}
            className="btn-ghost p-2 rounded-lg"
            style={{ background: 'none', border: 'none', cursor: 'pointer' }}
          >
            <Search size={18} color="var(--text-secondary)" />
          </button>

          {/* Notifications */}
          <button
            className="btn-ghost p-2 rounded-lg relative"
            style={{ background: 'none', border: 'none', cursor: 'pointer' }}
          >
            <Bell size={18} color="var(--text-secondary)" />
            {notifications > 0 && (
              <span
                className="absolute -top-0.5 -right-0.5 w-4 h-4 flex items-center justify-center rounded-full text-xs font-bold"
                style={{ background: 'var(--danger)', color: '#fff', fontSize: 9 }}
              >
                {notifications}
              </span>
            )}
          </button>

          {/* User Dropdown */}
          <div className="relative hidden md:block">
            <button
              onClick={() => setMenuOpen(!menuOpen)}
              className="flex items-center gap-2 px-2 py-1 rounded-lg cursor-pointer"
              style={{ background: 'none', border: 'none' }}
            >
              <div
                className="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold"
                style={{ background: 'linear-gradient(135deg, var(--gold), var(--gold-dark))', color: '#0a0a0f' }}
              >
                АП
              </div>
              <ChevronDown size={14} color="var(--text-muted)" />
            </button>
            <AnimatePresence>
              {menuOpen && (
                <motion.div
                  initial={{ opacity: 0, y: -8 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -8 }}
                  className="glass-card absolute right-0 top-full mt-2 w-48 py-2"
                  onClick={() => setMenuOpen(false)}
                >
                  <Link to="/profile" className="block px-4 py-2 text-sm no-underline" style={{ color: 'var(--text-primary)' }}>
                    Мой профиль
                  </Link>
                  <Link to="/matches" className="block px-4 py-2 text-sm no-underline" style={{ color: 'var(--text-primary)' }}>
                    История матчей
                  </Link>
                  <div style={{ borderTop: '1px solid var(--bg-glass-border)', margin: '8px 0' }} />
                  <button className="block w-full text-left px-4 py-2 text-sm" style={{ color: 'var(--danger)', background: 'none', border: 'none', cursor: 'pointer' }}>
                    Выйти
                  </button>
                </motion.div>
              )}
            </AnimatePresence>
          </div>

          {/* Mobile Menu Toggle */}
          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="md:hidden p-2 rounded-lg"
            style={{ background: 'none', border: 'none', cursor: 'pointer' }}
          >
            {mobileMenuOpen ? <X size={20} color="var(--text-primary)" /> : <Menu size={20} color="var(--text-secondary)" />}
          </button>
        </div>
      </div>

      {/* Mobile Menu */}
      <AnimatePresence>
        {mobileMenuOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            className="md:hidden overflow-hidden"
            style={{ borderBottom: '1px solid var(--bg-glass-border)' }}
          >
            <nav className="container py-4 flex flex-col gap-1">
              {navItems.map((item) => {
                const isActive = location.pathname === item.path
                return (
                  <Link
                    key={item.path}
                    to={item.path}
                    onClick={() => setMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium no-underline"
                    style={{
                      color: isActive ? 'var(--gold)' : 'var(--text-secondary)',
                      background: isActive ? 'rgba(212,168,67,0.08)' : 'transparent',
                    }}
                  >
                    <item.icon size={18} />
                    {item.label}
                  </Link>
                )
              })}
            </nav>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  )
}
