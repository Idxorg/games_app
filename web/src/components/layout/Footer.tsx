import { Link } from 'react-router-dom'
import { Gamepad2, Heart, Code } from 'lucide-react'

export function Footer() {
  return (
    <footer className="footer-border">
      <div className="container">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* Brand */}
          <div>
            <div className="flex items-center gap-2 mb-3">
              <Gamepad2 size={20} color="var(--gold)" />
              <span className="font-bold text-accent">Game Portal</span>
            </div>
            <p className="text-sm text-muted">
              Корпоративная игровая платформа для отдыха и соревнований
            </p>
          </div>

          {/* Links */}
          <div>
            <h4 className="text-sm font-semibold text-secondary mb-3">Навигация</h4>
            <div className="flex flex-col gap-2">
              <Link to="/" className="footer-link">Главная</Link>
              <Link to="/tournaments" className="footer-link">Турниры</Link>
              <Link to="/leaderboard" className="footer-link">Рейтинг</Link>
              <Link to="/matches" className="footer-link">Матчи</Link>
            </div>
          </div>

          {/* Info */}
          <div>
            <h4 className="text-sm font-semibold text-secondary mb-3">Информация</h4>
            <div className="flex flex-col gap-2">
              <div className="flex items-center gap-2 text-sm text-muted">
                <Code size={14} />
                <span>Версия 1.0.0</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-muted">
                <Heart size={14} />
                <span>Корпоративный портал</span>
              </div>
            </div>
          </div>
        </div>

        <div className="mt-8 pt-4 text-center text-xs text-muted" style={{ borderTop: '1px solid var(--bg-glass-border)' }}>
          Game Portal &mdash; Корпоративная игровая платформа &mdash; 2026
        </div>
      </div>
    </footer>
  )
}
