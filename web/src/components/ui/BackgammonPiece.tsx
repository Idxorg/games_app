import React from 'react'

interface BackgammonPieceProps {
  size?: number
  className?: string
}

export function BackgammonPiece({ size = 200, className = '' }: BackgammonPieceProps) {
  return (
    <svg width={size} height={size * 0.7} viewBox="0 0 200 140" className={className}>
      <defs>
        <linearGradient id="bg-board" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#5d4e37" />
          <stop offset="100%" stopColor="#3e2f1f" />
        </linearGradient>
        <linearGradient id="bg-tri-dark" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stopColor="#8b0000" />
          <stop offset="100%" stopColor="#4a0000" />
        </linearGradient>
        <linearGradient id="bg-tri-light" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stopColor="#f5deb3" />
          <stop offset="100%" stopColor="#d2b48c" />
        </linearGradient>
        <filter id="bg-shadow">
          <feDropShadow dx="0" dy="3" stdDeviation="4" floodOpacity="0.5" />
        </filter>
      </defs>
      {/* Board border */}
      <rect x="0" y="0" width="200" height="140" rx="8" fill="#2a1f10" filter="url(#bg-shadow)" />
      <rect x="2" y="2" width="196" height="136" rx="6" fill="url(#bg-board)" />
      {/* Center divider */}
      <rect x="96" y="2" width="8" height="136" fill="#2a1f10" rx="2" />
      {/* Top triangles (left) */}
      <polygon points="15,5 25,5 20,65" fill="url(#bg-tri-dark)" />
      <polygon points="30,5 40,5 35,65" fill="url(#bg-tri-light)" />
      <polygon points="45,5 55,5 50,65" fill="url(#bg-tri-dark)" />
      <polygon points="60,5 70,5 65,65" fill="url(#bg-tri-light)" />
      <polygon points="75,5 85,5 80,65" fill="url(#bg-tri-dark)" />
      <polygon points="90,5 95,5 92,65" fill="url(#bg-tri-light)" />
      {/* Top triangles (right) */}
      <polygon points="105,5 115,5 110,65" fill="url(#bg-tri-dark)" />
      <polygon points="120,5 130,5 125,65" fill="url(#bg-tri-light)" />
      <polygon points="135,5 145,5 140,65" fill="url(#bg-tri-dark)" />
      <polygon points="150,5 160,5 155,65" fill="url(#bg-tri-light)" />
      <polygon points="165,5 175,5 170,65" fill="url(#bg-tri-dark)" />
      <polygon points="180,5 190,5 185,65" fill="url(#bg-tri-light)" />
      {/* Bottom triangles (mirrored) */}
      <polygon points="15,135 25,135 20,75" fill="url(#bg-tri-light)" />
      <polygon points="30,135 40,135 35,75" fill="url(#bg-tri-dark)" />
      <polygon points="45,135 55,135 50,75" fill="url(#bg-tri-light)" />
      <polygon points="60,135 70,135 65,75" fill="url(#bg-tri-dark)" />
      <polygon points="75,135 85,135 80,75" fill="url(#bg-tri-light)" />
      <polygon points="90,135 95,135 92,75" fill="url(#bg-tri-dark)" />
      <polygon points="105,135 115,135 110,75" fill="url(#bg-tri-light)" />
      <polygon points="120,135 130,135 125,75" fill="url(#bg-tri-dark)" />
      <polygon points="135,135 145,135 140,75" fill="url(#bg-tri-light)" />
      <polygon points="150,135 160,135 155,75" fill="url(#bg-tri-dark)" />
      <polygon points="165,135 175,135 170,75" fill="url(#bg-tri-light)" />
      <polygon points="180,135 190,135 185,75" fill="url(#bg-tri-dark)" />
      {/* Sample pieces */}
      <circle cx="20" cy="20" r="5" fill="#f5f5dc" stroke="#d2b48c" strokeWidth="0.8" />
      <circle cx="20" cy="32" r="5" fill="#f5f5dc" stroke="#d2b48c" strokeWidth="0.8" />
      <circle cx="20" cy="44" r="5" fill="#f5f5dc" stroke="#d2b48c" strokeWidth="0.8" />
      <circle cx="35" cy="20" r="5" fill="#1a1a1a" stroke="#333" strokeWidth="0.8" />
      <circle cx="35" cy="32" r="5" fill="#1a1a1a" stroke="#333" strokeWidth="0.8" />
      <circle cx="125" cy="20" r="5" fill="#1a1a1a" stroke="#333" strokeWidth="0.8" />
      <circle cx="110" cy="120" r="5" fill="#f5f5dc" stroke="#d2b48c" strokeWidth="0.8" />
      <circle cx="110" cy="108" r="5" fill="#f5f5dc" stroke="#d2b48c" strokeWidth="0.8" />
      <circle cx="140" cy="120" r="5" fill="#1a1a1a" stroke="#333" strokeWidth="0.8" />
    </svg>
  )
}
