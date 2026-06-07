import React from 'react'

interface TrophyProps {
  type: 'gold' | 'silver' | 'bronze'
  size?: number
  className?: string
}

export function Trophy({ type, size = 40, className = '' }: TrophyProps) {
  const colors = {
    gold: { fill: '#d4a843', highlight: '#f0d080', stroke: '#b8922e' },
    silver: { fill: '#a0a0b0', highlight: '#d0d0e0', stroke: '#808090' },
    bronze: { fill: '#cd7f32', highlight: '#e8a050', stroke: '#a06020' },
  }
  const c = colors[type]

  return (
    <svg width={size} height={size} viewBox="0 0 40 44" className={className}>
      <defs>
        <linearGradient id={`trophy-${type}`} x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor={c.highlight} />
          <stop offset="100%" stopColor={c.fill} />
        </linearGradient>
      </defs>
      {/* Cup */}
      <path d="M10 8 L10 22 Q10 30 20 30 Q30 30 30 22 L30 8 Z" fill={`url(#trophy-${type})`} stroke={c.stroke} strokeWidth="1" />
      {/* Left handle */}
      <path d="M10 12 Q2 12 2 18 Q2 24 10 24" fill="none" stroke={c.fill} strokeWidth="2.5" strokeLinecap="round" />
      {/* Right handle */}
      <path d="M30 12 Q38 12 38 18 Q38 24 30 24" fill="none" stroke={c.fill} strokeWidth="2.5" strokeLinecap="round" />
      {/* Stem */}
      <rect x="18" y="30" width="4" height="6" fill={c.fill} />
      {/* Base */}
      <rect x="12" y="36" width="16" height="4" rx="2" fill={`url(#trophy-${type})`} stroke={c.stroke} strokeWidth="0.5" />
      {/* Star */}
      <polygon points="20,12 21.5,16 25.5,16 22.5,18.5 23.5,22.5 20,20 16.5,22.5 17.5,18.5 14.5,16 18.5,16" fill={c.stroke} opacity="0.5" />
    </svg>
  )
}

export function Medal({ type, size = 32, className = '' }: { type: 'gold' | 'silver' | 'bronze'; size?: number; className?: string }) {
  const colors = {
    gold: { fill: '#d4a843', ribbon: '#c0392b', highlight: '#f0d080' },
    silver: { fill: '#a0a0b0', ribbon: '#2980b9', highlight: '#d0d0e0' },
    bronze: { fill: '#cd7f32', ribbon: '#27ae60', highlight: '#e8a050' },
  }
  const c = colors[type]

  return (
    <svg width={size} height={size} viewBox="0 0 32 40" className={className}>
      {/* Ribbon */}
      <polygon points="10,2 12,0 14,2 12,14" fill={c.ribbon} />
      <polygon points="22,2 20,0 18,2 20,14" fill={c.ribbon} />
      {/* Medal circle */}
      <circle cx="16" cy="24" r="12" fill={c.fill} stroke={c.fill} strokeWidth="1" />
      <circle cx="16" cy="24" r="10" fill="none" stroke={c.highlight} strokeWidth="1" opacity="0.5" />
      {/* Number */}
      <text x="16" y="28" textAnchor="middle" fontSize="12" fontWeight="bold" fill={type === 'gold' ? '#3e2f1f' : '#fff'}>
        {type === 'gold' ? '1' : type === 'silver' ? '2' : '3'}
      </text>
    </svg>
  )
}

export function Crown({ size = 32, className = '' }: { size?: number; className?: string }) {
  return (
    <svg width={size} height={size} viewBox="0 0 32 28" className={className}>
      <defs>
        <linearGradient id="crown-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#f0d080" />
          <stop offset="100%" stopColor="#d4a843" />
        </linearGradient>
      </defs>
      <path d="M2 22 L2 10 L8 16 L16 4 L24 16 L30 10 L30 22 Z" fill="url(#crown-grad)" stroke="#b8922e" strokeWidth="1" />
      <rect x="2" y="22" width="28" height="4" rx="1" fill="url(#crown-grad)" stroke="#b8922e" strokeWidth="0.5" />
      <circle cx="16" cy="6" r="2" fill="#d4a843" stroke="#b8922e" strokeWidth="0.5" />
      <circle cx="6" cy="13" r="1.5" fill="#d4a843" stroke="#b8922e" strokeWidth="0.5" />
      <circle cx="26" cy="13" r="1.5" fill="#d4a843" stroke="#b8922e" strokeWidth="0.5" />
    </svg>
  )
}
