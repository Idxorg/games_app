import { Link } from 'react-router-dom'
import { Gamepad2, Github, Keyboard, Heart } from 'lucide-react'

export function Footer() {
  return (
    <footer style={{ borderTop: '1px solid var(--bg-glass-border)', marginTop: 'auto', padding: '32px 0' }}>
      <div className="container">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* Brand */}
          <div>
            <div className="flex items-center gap-2 mb-3">
              <Gamepad2 size={20} color="var(--gold)" />
              <span className="font-bold" style={{ color: 'var(--gold)' }}>Game Portal</span>
            </div>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>
              Корпоративная игровая платформа для отдыха и соревнований
            </p>
          </div>

          {/* Links */}
          <div>
            <h4 className="text-sm font-semibold mb-3" style={{ color: 'var(--text-secondary)' }}>Навигация</h4>
            <div className="flex flex-col gap-2">
              <Link to="/" className="text-sm no-underline" style={{ color: 'var(--text-muted)' }}>Главная</Link>
              <Link to="/tournaments" className="text-sm no-underline" style={{ color: 'var(--text-muted)' }}>Турниры</Link>
              <Link to="/leaderboard" className="text-sm no-underline" style={{ color: 'var(--text-muted)' }}>Рейтинг</Link>
              <Link to="/matches" className="text-sm no-underline" style={{ color: 'var(--text-muted)' }}>Матчи</Link>
            </div>
          </div>

          {/* Info */}
          <div>
            <h4 className="text-sm font-semibold mb-3" style={{ color: 'var(--text-secondary)' }}>Информация</h4>
            <div className="flex flex-col gap-2">
              <div className="flex items-center gap-2 text-sm" style={{ color: 'var(--text-muted)' }}>
                <Keyboard size={14} />
                <span>Ctrl+K -- поиск</span>
              </div>
              <div className="flex items-center gap-2 text-sm" style={{ color: 'var(--text-muted)' }}>
                <Heart size={14} />
                <span>Версия 1.0.0</span>
              </div>
              <div className="flex items-center gap-2 text-sm" style={{ color: 'var(--text-muted)' }}>
                <Github size={14} />
                <span>Корпоративный портал</span>
              </div>
            </div>
          </div>
        </div>

        <div className="mt-8 pt-4 text-center text-xs" style={{ borderTop: '1px solid var(--bg-glass-border)', color: 'var(--text-muted)' }}>
          Game Portal -- Корпоративная игровая платформа -- 2026
        </div>
      </div>
    </footer>
  )
}
