/* ============================================
   GamePortal - Custom SVG Icon Library
   Hand-crafted vector icons for all games
   ============================================ */

const GameIcons = (() => {
    'use strict';

    // Helper: create SVG string
    function svg(viewBox, content, opts = {}) {
        const size = opts.size || 30;
        const cls = opts.class || '';
        const fill = opts.fill || 'currentColor';
        return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="${viewBox}" width="${size}" height="${size}" class="${cls}" fill="${fill}">${content}</svg>`;
    }

    // =====================
    // CHESS PIECES (30x30)
    // =====================

    const chess = {
        // White King
        king() {
            return svg('0 0 30 30', `
                <defs>
                    <linearGradient id="kg" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="#f0e6d0"/>
                        <stop offset="100%" stop-color="#c8b898"/>
                    </linearGradient>
                </defs>
                <path d="M15 3l-1.5 3h3L15 3z" fill="#d4a843"/>
                <path d="M15 3l-0.5 1.5L15 6l0.5-1.5L15 3z" fill="#e8c36a"/>
                <rect x="13" y="6" width="4" height="2" rx="0.5" fill="url(#kg)"/>
                <rect x="12" y="8" width="6" height="3" rx="0.8" fill="url(#kg)"/>
                <path d="M9 12c0-1.5 2.5-2 3.5-1l2.5 2.5 2.5-2.5c1-1 3.5-0.5 3.5 1v2H9v-2z" fill="url(#kg)"/>
                <rect x="9" y="14" width="12" height="3" rx="0.8" fill="url(#kg)"/>
                <rect x="8.5" y="17" width="13" height="3" rx="0.5" fill="url(#kg)"/>
                <rect x="7" y="20" width="16" height="2" rx="0.5" fill="url(#kg)"/>
                <rect x="6" y="22" width="18" height="3" rx="1" fill="url(#kg)"/>
                <rect x="5" y="25.5" width="20" height="2.5" rx="1" fill="url(#kg)" stroke="#a09070" stroke-width="0.3"/>
                <circle cx="12" cy="15.5" r="1" fill="#a09070" opacity="0.3"/>
                <circle cx="18" cy="15.5" r="1" fill="#a09070" opacity="0.3"/>
            `);
        },

        // White Queen
        queen() {
            return svg('0 0 30 30', `
                <defs>
                    <linearGradient id="qn" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="#f0e6d0"/>
                        <stop offset="100%" stop-color="#c8b898"/>
                    </linearGradient>
                </defs>
                <circle cx="15" cy="4.5" r="1.8" fill="#d4a843"/>
                <path d="M10 8l2 5-4 4h16l-4-4 2-5c-1.5-1.5 1.5-4 1.5-4s-3 1-3.5 2c0 0-1.5-2-1.5-2s0 2.5-1.5 4z" fill="url(#qn)"/>
                <circle cx="8.5" cy="8" r="1.3" fill="url(#qn)"/>
                <circle cx="15" cy="5.5" r="1.3" fill="url(#qn)"/>
                <circle cx="21.5" cy="8" r="1.3" fill="url(#qn)"/>
                <rect x="8" y="17" width="14" height="3" rx="0.8" fill="url(#qn)"/>
                <rect x="7.5" y="20" width="15" height="3" rx="0.5" fill="url(#qn)"/>
                <rect x="5" y="25.5" width="20" height="2.5" rx="1" fill="url(#qn)" stroke="#a09070" stroke-width="0.3"/>
            `);
        },

        // White Rook
        rook() {
            return svg('0 0 30 30', `
                <defs>
                    <linearGradient id="rk" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="#f0e6d0"/>
                        <stop offset="100%" stop-color="#c8b898"/>
                    </linearGradient>
                </defs>
                <rect x="7" y="4" width="3.5" height="3" rx="0.5" fill="url(#rk)"/>
                <rect x="13.2" y="4" width="3.5" height="3" rx="0.5" fill="url(#rk)"/>
                <rect x="19.5" y="4" width="3.5" height="3" rx="0.5" fill="url(#rk)"/>
                <rect x="6" y="7" width="18" height="2" rx="0.5" fill="url(#rk)"/>
                <rect x="8" y="9" width="14" height="7" rx="0.8" fill="url(#rk)"/>
                <rect x="7" y="16" width="16" height="3" rx="0.8" fill="url(#rk)"/>
                <rect x="6.5" y="19" width="17" height="3" rx="0.5" fill="url(#rk)"/>
                <rect x="5" y="25.5" width="20" height="2.5" rx="1" fill="url(#rk)" stroke="#a09070" stroke-width="0.3"/>
                <rect x="11" y="10" width="8" height="1.5" rx="0.5" fill="#a09070" opacity="0.25"/>
            `);
        },

        // White Bishop
        bishop() {
            return svg('0 0 30 30', `
                <defs>
                    <linearGradient id="bp" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="#f0e6d0"/>
                        <stop offset="100%" stop-color="#c8b898"/>
                    </linearGradient>
                </defs>
                <path d="M15 3l1.5 2.5c0 0 0.5 1-0.5 2l-1 2.5c0 0 3 2.5 3 5v3H12v-3c0-2.5 3-5 3-5l-1-2.5c-1-1-0.5-2-0.5-2L15 3z" fill="url(#bp)"/>
                <circle cx="15" cy="3.2" r="1.3" fill="#d4a843"/>
                <path d="M12 18h6l1 3H11l1-3z" fill="url(#bp)"/>
                <rect x="10.5" y="21" width="9" height="3" rx="0.5" fill="url(#bp)"/>
                <rect x="5" y="25.5" width="20" height="2.5" rx="1" fill="url(#bp)" stroke="#a09070" stroke-width="0.3"/>
            `);
        },

        // White Knight
        knight() {
            return svg('0 0 30 30', `
                <defs>
                    <linearGradient id="kn" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="#f0e6d0"/>
                        <stop offset="100%" stop-color="#c8b898"/>
                    </linearGradient>
                </defs>
                <path d="M10 24l1.5-5c-0.5-1.5 0-3 1.5-4l-1-3c0-3 2-6 4-7 1-0.5 2.5-0.3 3 0.5 1 1.5 0.5 3-0.5 4 -1 1-2 1.5-1.5 3l0.5 2c0.5 2 3 3.5 4 5l0.5 4.5H10z" fill="url(#kn)"/>
                <circle cx="17.5" cy="6" r="1.2" fill="#a09070"/>
                <path d="M15 8c0.5 0 1 0.2 1.2 0.8" stroke="#a09070" stroke-width="0.5" fill="none"/>
                <rect x="5" y="25.5" width="20" height="2.5" rx="1" fill="url(#kn)" stroke="#a09070" stroke-width="0.3"/>
            `);
        },

        // White Pawn
        pawn() {
            return svg('0 0 30 30', `
                <defs>
                    <linearGradient id="pw" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="#f0e6d0"/>
                        <stop offset="100%" stop-color="#c8b898"/>
                    </linearGradient>
                </defs>
                <circle cx="15" cy="7" r="4" fill="url(#pw)"/>
                <path d="M11 11c0-1 4-3 8 0v4h-8v-4z" fill="url(#pw)"/>
                <rect x="10" y="15" width="10" height="3" rx="0.8" fill="url(#pw)"/>
                <rect x="9.5" y="18" width="11" height="3" rx="0.5" fill="url(#pw)"/>
                <rect x="5" y="25.5" width="20" height="2.5" rx="1" fill="url(#pw)" stroke="#a09070" stroke-width="0.3"/>
            `);
        }
    };

    // =====================
    // CHECKERS PIECES
    // =====================

    const checkers = {
        red(opts = {}) {
            return svg('0 0 30 30', `
                <defs>
                    <linearGradient id="cr" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="#e85050"/>
                        <stop offset="100%" stop-color="#b03030"/>
                    </linearGradient>
                    <radialGradient id="crs" cx="0.4" cy="0.35" r="0.6">
                        <stop offset="0%" stop-color="#ff7070" stop-opacity="0.5"/>
                        <stop offset="100%" stop-color="#b03030" stop-opacity="0"/>
                    </radialGradient>
                </defs>
                <ellipse cx="15" cy="22" rx="12" ry="3" fill="rgba(0,0,0,0.3)"/>
                <circle cx="15" cy="15" r="12" fill="url(#cr)"/>
                <circle cx="15" cy="15" r="12" fill="url(#crs)"/>
                <circle cx="15" cy="15" r="8" fill="none" stroke="rgba(255,255,255,0.15)" stroke-width="1"/>
                <circle cx="15" cy="15" r="12" fill="none" stroke="rgba(0,0,0,0.2)" stroke-width="0.5"/>
            `, opts);
        },

        black(opts = {}) {
            return svg('0 0 30 30', `
                <defs>
                    <linearGradient id="cb" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="#4a4555"/>
                        <stop offset="100%" stop-color="#2a2535"/>
                    </linearGradient>
                    <radialGradient id="cbs" cx="0.4" cy="0.35" r="0.6">
                        <stop offset="0%" stop-color="#7a7585" stop-opacity="0.4"/>
                        <stop offset="100%" stop-color="#2a2535" stop-opacity="0"/>
                    </radialGradient>
                </defs>
                <ellipse cx="15" cy="22" rx="12" ry="3" fill="rgba(0,0,0,0.4)"/>
                <circle cx="15" cy="15" r="12" fill="url(#cb)"/>
                <circle cx="15" cy="15" r="12" fill="url(#cbs)"/>
                <circle cx="15" cy="15" r="8" fill="none" stroke="rgba(255,255,255,0.08)" stroke-width="1"/>
                <circle cx="15" cy="15" r="12" fill="none" stroke="rgba(0,0,0,0.3)" stroke-width="0.5"/>
            `, opts);
        },

        redKing(opts = {}) {
            return svg('0 0 30 30', `
                <defs>
                    <linearGradient id="ckr" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="#e85050"/>
                        <stop offset="100%" stop-color="#b03030"/>
                    </linearGradient>
                </defs>
                <ellipse cx="15" cy="22" rx="12" ry="3" fill="rgba(0,0,0,0.3)"/>
                <circle cx="15" cy="15" r="12" fill="url(#ckr)"/>
                <circle cx="15" cy="15" r="8" fill="none" stroke="rgba(255,255,255,0.15)" stroke-width="1"/>
                <path d="M12 10l3-4 3 4-3 1.5z" fill="#d4a843"/>
                <circle cx="15" cy="6.5" r="1" fill="#e8c36a"/>
            `, opts);
        }
    };

    // =====================
    // BACKGAMMON
    // =====================

    const backgammon = {
        board() {
            return svg('0 0 30 30', `
                <defs>
                    <linearGradient id="bg-board" x1="0" y1="0" x2="1" y2="1">
                        <stop offset="0%" stop-color="#5c4a32"/>
                        <stop offset="100%" stop-color="#3a2a1a"/>
                    </linearGradient>
                </defs>
                <rect x="2" y="2" width="26" height="26" rx="3" fill="url(#bg-board)" stroke="#7a6a52" stroke-width="0.5"/>
                <rect x="4" y="4" width="22" height="22" rx="2" fill="none" stroke="#8a7a62" stroke-width="0.3"/>
                <line x1="15" y1="4" x2="15" y2="26" stroke="#8a7a62" stroke-width="0.5"/>
                <circle cx="6" cy="4.5" r="0.6" fill="#d4a843"/>
                <circle cx="24" cy="4.5" r="0.6" fill="#d4a843"/>
                <circle cx="6" cy="25.5" r="0.6" fill="#d4a843"/>
                <circle cx="24" cy="25.5" r="0.6" fill="#d4a843"/>
                <path d="M5 4v8" stroke="#4a3a28" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M9 4v6" stroke="#c8a878" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M13 4v8" stroke="#4a3a28" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M17 4v6" stroke="#c8a878" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M21 4v8" stroke="#4a3a28" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M25 4v6" stroke="#c8a878" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M5 26v-8" stroke="#c8a878" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M9 26v-6" stroke="#4a3a28" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M13 26v-8" stroke="#c8a878" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M17 26v-6" stroke="#4a3a28" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M21 26v-8" stroke="#c8a878" stroke-width="2.5" stroke-linecap="round"/>
                <path d="M25 26v-6" stroke="#4a3a28" stroke-width="2.5" stroke-linecap="round"/>
                <circle cx="7" cy="10" r="2.5" fill="#e8e0d0" stroke="#b8a888" stroke-width="0.3"/>
                <circle cx="7" cy="16" r="2.5" fill="#e8e0d0" stroke="#b8a888" stroke-width="0.3"/>
                <circle cx="23" cy="18" r="2.5" fill="#4a4555" stroke="#6a6575" stroke-width="0.3"/>
                <circle cx="23" cy="24" r="2.5" fill="#4a4555" stroke="#6a6575" stroke-width="0.3"/>
            `);
        },

        piece(opts = {}) {
            const color = opts.color || 'white';
            const fill = color === 'white' ? '#e8e0d0' : '#3a3530';
            const stroke = color === 'white' ? '#b8a888' : '#5a5550';
            return svg('0 0 30 30', `
                <ellipse cx="15" cy="22" rx="10" ry="2.5" fill="rgba(0,0,0,0.25)"/>
                <circle cx="15" cy="15" r="10" fill="${fill}" stroke="${stroke}" stroke-width="0.5"/>
                <circle cx="13" cy="13" r="3" fill="rgba(255,255,255,0.08)"/>
            `, opts);
        }
    };

    // =====================
    // UTILITY ICONS
    // =====================

    const trophy = (opts = {}) => svg('0 0 30 30', `
        <defs>
            <linearGradient id="tp" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#ffd700"/>
                <stop offset="100%" stop-color="#b8922f"/>
            </linearGradient>
        </defs>
        <path d="M8 5h14v10c0 4-3 7-7 7s-7-3-7-7V5z" fill="url(#tp)"/>
        <path d="M8 7H5c0 4 1.5 6 3 7V7z" fill="url(#tp)"/>
        <path d="M22 7h3c0 4-1.5 6-3 7V7z" fill="url(#tp)"/>
        <rect x="13" y="22" width="4" height="3" rx="0.5" fill="#b8922f"/>
        <rect x="10" y="25" width="10" height="2" rx="1" fill="#b8922f"/>
        <path d="M13 8l2 3 2-3" stroke="#b8922f" stroke-width="0.5" fill="none"/>
    `, opts);

    const medal = (opts = {}) => svg('0 0 30 30', `
        <defs>
            <linearGradient id="md" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#ffd700"/>
                <stop offset="100%" stop-color="#daa520"/>
            </linearGradient>
        </defs>
        <path d="M8 3l7 10 7-10" stroke="#b8922f" stroke-width="1.5" fill="none"/>
        <path d="M11 3l4 6 4-6" fill="#d4a843"/>
        <circle cx="15" cy="20" r="8" fill="url(#md)" stroke="#b8922f" stroke-width="0.5"/>
        <circle cx="15" cy="20" r="5" fill="none" stroke="#b8922f" stroke-width="0.5"/>
        <text x="15" y="22.5" text-anchor="middle" font-size="7" font-weight="bold" fill="#8a6a10">1</text>
    `, opts);

    const crown = (opts = {}) => svg('0 0 30 30', `
        <defs>
            <linearGradient id="cr" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#ffd700"/>
                <stop offset="100%" stop-color="#b8922f"/>
            </linearGradient>
        </defs>
        <path d="M5 22l2-10 5 4 3-9 3 9 5-4 2 10H5z" fill="url(#cr)"/>
        <rect x="5" y="22" width="20" height="3" rx="1" fill="url(#cr)"/>
        <circle cx="7" cy="11.5" r="1" fill="#e8c36a"/>
        <circle cx="15" cy="7" r="1" fill="#e8c36a"/>
        <circle cx="23" cy="11.5" r="1" fill="#e8c36a"/>
    `, opts);

    const star = (opts = {}) => svg('0 0 30 30', `
        <polygon points="15,3 18.5,11 27,11.5 20.5,17 22.5,26 15,21.5 7.5,26 9.5,17 3,11.5 11.5,11" fill="${(opts && opts.fill) || '#d4a843'}"/>
    `, opts);

    const snake = (opts = {}) => svg('0 0 30 30', `
        <defs>
            <linearGradient id="sn" x1="0" y1="0" x2="1" y2="1">
                <stop offset="0%" stop-color="#4ade80"/>
                <stop offset="100%" stop-color="#22c55e"/>
            </linearGradient>
        </defs>
        <path d="M7 24h4v-4h4v-4h4v-4h4v4h-4v4h-4v4h-4v4H7v-4z" fill="url(#sn)" rx="1"/>
        <path d="M7 24h4" stroke="#166534" stroke-width="0.5" fill="none"/>
        <circle cx="22" cy="11" r="1" fill="#166534"/>
        <circle cx="21.5" cy="10.5" r="0.4" fill="white"/>
        <rect x="19" y="19" width="4" height="4" rx="0.5" fill="url(#sn)" opacity="0.7"/>
    `, opts);

    const mines = (opts = {}) => svg('0 0 30 30', `
        <defs>
            <radialGradient id="mb" cx="0.4" cy="0.35" r="0.6">
                <stop offset="0%" stop-color="#606070"/>
                <stop offset="100%" stop-color="#303040"/>
            </radialGradient>
        </defs>
        <circle cx="15" cy="15" r="10" fill="url(#mb)"/>
        <circle cx="15" cy="15" r="10" fill="none" stroke="#505060" stroke-width="0.5"/>
        <ellipse cx="15" cy="24" rx="10" ry="2" fill="rgba(0,0,0,0.2)"/>
        <rect x="14" y="6" width="2" height="4" rx="1" fill="#505060"/>
        <circle cx="15" cy="5" r="1.5" fill="rgba(255,200,100,0.6)"/>
        <circle cx="8" cy="15" width="4" height="4" rx="1" fill="#505060"/>
        <circle cx="22" cy="15" width="4" height="4" rx="1" fill="#505060"/>
        <circle cx="15" cy="8" r="3" fill="#505060"/>
        <circle cx="15" cy="22" r="3" fill="#505060"/>
        <circle cx="9" cy="10" r="2.5" fill="#505060"/>
        <circle cx="21" cy="10" r="2.5" fill="#505060"/>
        <circle cx="9" cy="20" r="2.5" fill="#505060"/>
        <circle cx="21" cy="20" r="2.5" fill="#505060"/>
        <path d="M15 12v2l-1.5 1.5M15 14l1.5 1.5M15 14v3" stroke="#303040" stroke-width="1" fill="none" stroke-linecap="round"/>
    `, opts);

    const arena = (opts = {}) => svg('0 0 30 30', `
        <defs>
            <linearGradient id="ar" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#a78bfa"/>
                <stop offset="100%" stop-color="#7c3aed"/>
            </linearGradient>
        </defs>
        <circle cx="15" cy="15" r="12" fill="none" stroke="url(#ar)" stroke-width="1.5"/>
        <circle cx="15" cy="15" r="8" fill="none" stroke="url(#ar)" stroke-width="1"/>
        <circle cx="15" cy="15" r="4" fill="none" stroke="url(#ar)" stroke-width="0.8"/>
        <circle cx="15" cy="15" r="2" fill="url(#ar)"/>
        <path d="M15 3l2 5h-4l2-5z" fill="url(#ar)"/>
        <path d="M15 27l-2-5h4l-2 5z" fill="url(#ar)"/>
        <path d="M3 15l5-2v4l-5-2z" fill="url(#ar)"/>
        <path d="M27 15l-5 2v-4l5 2z" fill="url(#ar)"/>
        <path d="M6.5 6.5l5.5 2.5-2.8 2.8-2.7-5.3z" fill="url(#ar)"/>
        <path d="M23.5 23.5l-5.5-2.5 2.8-2.8 2.7 5.3z" fill="url(#ar)"/>
        <path d="M6.5 23.5l2.7-5.3 2.8 2.8-5.5 2.5z" fill="url(#ar)"/>
        <path d="M23.5 6.5l-2.7 5.3-2.8-2.8 5.5-2.5z" fill="url(#ar)"/>
    `, opts);

    const poker = (opts = {}) => svg('0 0 30 30', `
        <defs>
            <linearGradient id="pk" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#e8e0d0"/>
                <stop offset="100%" stop-color="#c8b898"/>
            </linearGradient>
        </defs>
        <rect x="3" y="6" width="14" height="20" rx="2" fill="url(#pk)" stroke="#a09070" stroke-width="0.5" transform="rotate(-5 10 16)"/>
        <rect x="13" y="4" width="14" height="20" rx="2" fill="url(#pk)" stroke="#a09070" stroke-width="0.5" transform="rotate(5 20 14)"/>
        <text x="10" y="17" text-anchor="middle" font-size="8" font-weight="bold" fill="#d4a843" transform="rotate(-5 10 16)">A</text>
        <text x="20" y="15" text-anchor="middle" font-size="8" font-weight="bold" fill="#e85050" transform="rotate(5 20 14)">K</text>
        <path d="M7 10c1 1 2 1 3 0c-0.5 2-1.5 2-3 0z" fill="#d4a843" transform="rotate(-5 10 16)"/>
        <path d="M17 8c1 1 2 1 3 0c-0.5 2-1.5 2-3 0z" fill="#e85050" transform="rotate(5 20 14)"/>
    `, opts);

    // =====================
    // GAME CARD ICONS (larger, colored)
    // =====================

    const gameIcons = {
        chess: () => chess.king(),
        checkers: () => checkers.red({ size: 30 }),
        backgammon: () => backgammon.board(),
        snake: () => snake(),
        mines: () => mines(),
        arena: () => arena(),
        poker: () => poker()
    };

    return {
        chess,
        checkers,
        backgammon,
        trophy,
        medal,
        crown,
        star,
        snake,
        mines,
        arena,
        poker,
        gameIcons
    };
})();

// Make available globally
window.GameIcons = GameIcons;
