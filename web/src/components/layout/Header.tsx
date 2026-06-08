import { useState, useCallback, useRef, useEffect } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { Search, ChevronDown, Menu, X, Gamepad2, Home, Trophy, BarChart3, Clock, User } from 'lucide-react'
import { useUserStore } from '../../stores/userStore'
import { useToastStore } from '../ui/Toast'
import { isEmbedMode } from '../../embedHandoff'

const navItems = [
  { path: '/', label: 'Главная', icon: Home },
  { path: '/tournaments', label: 'Турниры', icon: Trophy },
  { path: '/leaderboard', label: 'Рейтинг', icon: BarChart3 },
  { path: '/matches', label: 'Матчи', icon: Clock },
  { path: '/profile', label: 'Профиль', icon: User },
]

export function Header() {
  // Hide entire Header in embed mode
  if (isEmbedMode()) return null

  const [searchOpen, setSearchOpen] = useState(false)
  const [searchValue, setSearchValue] = useState('')
  const [menuOpen, setMenuOpen] = useState(false)
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const location = useLocation()
  const dropdownRef = useRef<HTMLDivElement>(null)
  const addToast = useToastStore((s) => s.addToast)
  const logout = useUserStore((s) => s.logout)
  const getInitials = useUserStore((s) => s.getInitials)
  const currentUser = useUserStore((s) => s.getCurrentUser())

  const initials = currentUser ? getInitials() : '?'

  const handleSearch = useCallback((value: string) => {
    setSearchValue(value)
  }, [])

  // Click-away handler for dropdown
  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setMenuOpen(false)
      }
    }
    if (menuOpen) {
      document.addEventListener('click', handleClickOutside)
    }
    return () => document.removeEventListener('click', handleClickOutside)
  }, [menuOpen])

  // Close mobile menu on route change
  useEffect(() => {
    setMobileMenuOpen(false)
  }, [location.pathname])

  const handleLogout = useCallback(() => {
    setMenuOpen(false)
    logout()
    addToast('Вы вышли из системы', 'info')
  }, [logout, addToast])

  return (
    <header className="glass header-sticky">
      <div className="container flex-between header-inner">
        {/* Logo */}
        <Link to="/" className="flex items-center gap-3 no-underline">
          <Gamepad2 size={28} color="var(--gold)" />
          <span className="text-lg font-bold text-accent">
            Игры · ЭР-Линк
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
                className={`nav-link ${isActive ? 'nav-link-active' : ''}`}
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
                  className="search-input"
                  aria-label="Поиск игр"
                  autoFocus
                />
              </motion.div>
            )}
          </AnimatePresence>
          <button
            onClick={() => setSearchOpen(!searchOpen)}
            className="btn-icon"
            aria-label={searchOpen ? 'Закрыть поиск' : 'Открыть поиск'}
          >
            <Search size={18} />
          </button>

          {/* User Dropdown */}
          <div className="relative hidden md:block" ref={dropdownRef}>
            <button
              onClick={() => setMenuOpen(!menuOpen)}
              className="flex items-center gap-2 px-2 py-1 rounded-lg cursor-pointer"
              style={{ background: 'none', border: 'none' }}
              aria-label="Меню пользователя"
              aria-expanded={menuOpen}
            >
              <div className="profile-avatar-sm">
                {initials}
              </div>
              <ChevronDown size={14} color="var(--text-muted)" />
            </button>
            <AnimatePresence>
              {menuOpen && (
                <motion.div
                  initial={{ opacity: 0, y: -8 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -8 }}
                  className="dropdown-menu absolute right-0 top-full mt-2 w-48"
                  onClick={(e) => e.stopPropagation()}
                >
                  <Link to="/profile" className="dropdown-item">
                    Мой профиль
                  </Link>
                  <Link to="/matches" className="dropdown-item">
                    История матчей
                  </Link>
                  <div className="dropdown-divider" />
                  <button
                    className="dropdown-item"
                    style={{ color: 'var(--danger)' }}
                    onClick={handleLogout}
                  >
                    Выйти
                  </button>
                </motion.div>
              )}
            </AnimatePresence>
          </div>

          {/* Mobile Menu Toggle */}
          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="btn-icon md:hidden"
            aria-label={mobileMenuOpen ? 'Закрыть меню' : 'Открыть меню'}
            aria-expanded={mobileMenuOpen}
          >
            {mobileMenuOpen ? <X size={20} color="var(--text-primary)" /> : <Menu size={20} />}
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
                    className={`nav-link-mobile ${isActive ? 'nav-link-active' : ''}`}
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
