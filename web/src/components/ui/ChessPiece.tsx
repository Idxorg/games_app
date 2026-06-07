import React from 'react'

type PieceType = 'king' | 'queen' | 'rook' | 'bishop' | 'knight' | 'pawn'
type PieceColor = 'white' | 'black'

interface ChessPieceProps {
  type: PieceType
  color?: PieceColor
  size?: number
  className?: string
}

const gradientId = (type: string, color: string) => `gp-${type}-${color}`

export function ChessPiece({ type, color = 'white', size = 40, className = '' }: ChessPieceProps) {
  const fillColor = color === 'white' ? '#f0e6d2' : '#1a1a2e'
  const strokeColor = color === 'white' ? '#c9b896' : '#3a3a5a'

  const pieces: Record<PieceType, React.ReactNode> = {
    king: (
      <>
        <defs>
          <linearGradient id={gradientId(type, color)} x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor={color === 'white' ? '#fff8e7' : '#2a2a4e'} />
            <stop offset="100%" stopColor={fillColor} />
          </linearGradient>
        </defs>
        <rect x="8" y="30" width="24" height="4" rx="1" fill={strokeColor} />
        <rect x="10" y="34" width="20" height="4" rx="1" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <path d="M12 34 L12 18 L14 12 L18 8 L20 12 L20 14 L24 14 L24 12 L26 8 L28 12 L28 18 L28 34" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <line x1="12" y1="18" x2="28" y2="18" stroke={strokeColor} strokeWidth="0.5" />
        <rect x="14" y="38" width="12" height="2" rx="1" fill={strokeColor} />
        <circle cx="20" cy="6" r="2.5" fill={fillColor} stroke={strokeColor} strokeWidth="0.5" />
        <line x1="20" y1="3.5" x2="20" y2="1" stroke={strokeColor} strokeWidth="1" />
        <line x1="17.5" y1="4" x2="16" y2="2" stroke={strokeColor} strokeWidth="0.8" />
        <line x1="22.5" y1="4" x2="24" y2="2" stroke={strokeColor} strokeWidth="0.8" />
      </>
    ),
    queen: (
      <>
        <defs>
          <linearGradient id={gradientId(type, color)} x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor={color === 'white' ? '#fff8e7' : '#2a2a4e'} />
            <stop offset="100%" stopColor={fillColor} />
          </linearGradient>
        </defs>
        <rect x="8" y="30" width="24" height="4" rx="1" fill={strokeColor} />
        <rect x="10" y="34" width="20" height="4" rx="1" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <path d="M12 34 L12 18 C12 14 14 12 16 12 C16 8 20 4 20 2 C20 4 24 8 24 12 C26 12 28 14 28 18 L28 34" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <circle cx="14" cy="12" r="2" fill={fillColor} stroke={strokeColor} strokeWidth="0.5" />
        <circle cx="20" cy="2" r="2" fill={fillColor} stroke={strokeColor} strokeWidth="0.5" />
        <circle cx="26" cy="12" r="2" fill={fillColor} stroke={strokeColor} strokeWidth="0.5" />
        <rect x="14" y="38" width="12" height="2" rx="1" fill={strokeColor} />
      </>
    ),
    rook: (
      <>
        <defs>
          <linearGradient id={gradientId(type, color)} x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor={color === 'white' ? '#fff8e7' : '#2a2a4e'} />
            <stop offset="100%" stopColor={fillColor} />
          </linearGradient>
        </defs>
        <rect x="8" y="30" width="24" height="4" rx="1" fill={strokeColor} />
        <rect x="10" y="34" width="20" height="4" rx="1" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <rect x="12" y="12" width="16" height="22" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <rect x="10" y="8" width="20" height="6" rx="1" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <rect x="11" y="8" width="4" height="4" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.3" />
        <rect x="17" y="8" width="4" height="4" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.3" />
        <rect x="23" y="8" width="4" height="4" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.3" />
        <rect x="14" y="38" width="12" height="2" rx="1" fill={strokeColor} />
      </>
    ),
    bishop: (
      <>
        <defs>
          <linearGradient id={gradientId(type, color)} x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor={color === 'white' ? '#fff8e7' : '#2a2a4e'} />
            <stop offset="100%" stopColor={fillColor} />
          </linearGradient>
        </defs>
        <rect x="8" y="30" width="24" height="4" rx="1" fill={strokeColor} />
        <rect x="10" y="34" width="20" height="4" rx="1" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <path d="M12 34 L12 22 Q12 16 16 14 Q18 12 20 6 Q22 12 24 14 Q28 16 28 22 L28 34" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <line x1="16" y1="26" x2="24" y2="26" stroke={strokeColor} strokeWidth="0.5" />
        <circle cx="20" cy="4" r="1.5" fill={fillColor} stroke={strokeColor} strokeWidth="0.5" />
        <line x1="20" y1="2.5" x2="20" y2="0.5" stroke={strokeColor} strokeWidth="0.8" />
        <rect x="14" y="38" width="12" height="2" rx="1" fill={strokeColor} />
      </>
    ),
    knight: (
      <>
        <defs>
          <linearGradient id={gradientId(type, color)} x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor={color === 'white' ? '#fff8e7' : '#2a2a4e'} />
            <stop offset="100%" stopColor={fillColor} />
          </linearGradient>
        </defs>
        <rect x="8" y="30" width="24" height="4" rx="1" fill={strokeColor} />
        <rect x="10" y="34" width="20" height="4" rx="1" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <path d="M14 34 L14 26 L12 20 L10 16 L12 10 L16 8 L22 8 L26 12 L26 20 L28 22 L28 26 L26 26 L26 28 L14 28 Z" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <circle cx="14" cy="10" r="1" fill={strokeColor} />
        <path d="M12 16 L14 14 L14 18 Z" fill={strokeColor} opacity="0.3" />
        <rect x="14" y="38" width="12" height="2" rx="1" fill={strokeColor} />
      </>
    ),
    pawn: (
      <>
        <defs>
          <linearGradient id={gradientId(type, color)} x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor={color === 'white' ? '#fff8e7' : '#2a2a4e'} />
            <stop offset="100%" stopColor={fillColor} />
          </linearGradient>
        </defs>
        <rect x="8" y="32" width="24" height="3" rx="1" fill={strokeColor} />
        <rect x="10" y="34" width="20" height="4" rx="1" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <path d="M16 34 L16 24 Q16 20 20 20 Q24 20 24 24 L24 34" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <circle cx="20" cy="16" r="5" fill={`url(#${gradientId(type, color)})`} stroke={strokeColor} strokeWidth="0.5" />
        <rect x="14" y="38" width="12" height="2" rx="1" fill={strokeColor} />
      </>
    ),
  }

  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 40 42"
      className={className}
      aria-label={`${color === 'white' ? 'Белый' : 'Чёрный'} ${type}`}
    >
      {pieces[type]}
    </svg>
  )
}
