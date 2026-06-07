import React from 'react'

interface CheckersPieceProps {
  color: 'red' | 'black'
  size?: number
  isKing?: boolean
  className?: string
}

export function CheckersPiece({ color, size = 40, isKing = false, className = '' }: CheckersPieceProps) {
  const fill = color === 'red' ? '#c0392b' : '#2c3e50'
  const highlight = color === 'red' ? '#e74c3c' : '#34495e'
  const innerFill = color === 'red' ? '#e74c3c' : '#34495e'

  return (
    <svg width={size} height={size} viewBox="0 0 40 40" className={className}>
      <defs>
        <linearGradient id={`checker-${color}`} x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor={highlight} />
          <stop offset="100%" stopColor={fill} />
        </linearGradient>
        <filter id={`shadow-${color}`}>
          <feDropShadow dx="0" dy="2" stdDeviation="2" floodOpacity="0.4" />
        </filter>
      </defs>
      <ellipse cx="20" cy="22" rx="16" ry="14" fill="rgba(0,0,0,0.3)" />
      <ellipse cx="20" cy="20" rx="16" ry="14" fill={`url(#checker-${color})`} stroke={color === 'red' ? '#a93226' : '#1a252f'} strokeWidth="1.5" />
      <ellipse cx="20" cy="20" rx="11" ry="9" fill="none" stroke={innerFill} strokeWidth="1" opacity="0.5" />
      <ellipse cx="20" cy="19" rx="6" ry="4" fill="rgba(255,255,255,0.1)" />
      {isKing && (
        <>
          <circle cx="20" cy="16" r="4" fill={color === 'red' ? '#f9e79f' : '#f9e79f'} stroke={color === 'red' ? '#d4a843' : '#d4a843'} strokeWidth="1" />
          <text x="20" y="18" textAnchor="middle" fontSize="7" fontWeight="bold" fill="#b8922e">K</text>
        </>
      )}
    </svg>
  )
}
