/* ============================================
   GamePortal - Main SPA Client
   Client-side routing, rendering, interactions
   ============================================ */

const Portal = (() => {
    'use strict';

    // =====================
    // SAMPLE DATA
    // =====================

    const GAMES = [
        {
            id: 'chess', name: 'Шахматы', desc: 'Классическая стратегическая игра с глубоким тактическим содержанием',
            category: 'strategy', players: 847, matchesToday: 156, icon: 'chess', elo: '1400-2200'
        },
        {
            id: 'checkers', name: 'Шашки', desc: 'Традиционная настольная игра — простые правила, сложная стратегия',
            category: 'strategy', players: 523, matchesToday: 89, icon: 'checkers', elo: '800-1800'
        },
        {
            id: 'backgammon', name: 'Нарды', desc: 'Игра в кости и стратегия — сочетание удачи и мастерства',
            category: 'strategy', players: 312, matchesToday: 67, icon: 'backgammon', elo: '1000-2000'
        },
        {
            id: 'snake', name: 'Змейка', desc: 'Классическая аркадная игра — кто наберёт больше очков',
            category: 'arcade', players: 1024, matchesToday: 342, icon: 'snake', elo: 'N/A'
        },
        {
            id: 'mines', name: 'Сапёр', desc: 'Логическая головоломка — рассчитывайте каждый ход',
            category: 'arcade', players: 678, matchesToday: 215, icon: 'mines', elo: 'N/A'
        },
        {
            id: 'arena', name: 'Арена', desc: 'Мультиплеерная арена — соревнование в реальном времени',
            category: 'strategy', players: 445, matchesToday: 128, icon: 'arena', elo: '1200-2500'
        },
        {
            id: 'poker', name: 'Покер', desc: 'Корпоративный турнирный покер — Техасский Холдем',
            category: 'card', players: 389, matchesToday: 94, icon: 'poker', elo: '1000-2000'
        }
    ];

    const LEADERBOARD_DATA = [
        { rank: 1, name: 'Петров Алексей', dept: 'Backend', elo: 2180, wins: 342, losses: 48, trend: +12, avatar: 1 },
        { rank: 2, name: 'Сидорова Мария', dept: 'Frontend', elo: 2120, wins: 298, losses: 62, trend: +8, avatar: 2 },
        { rank: 3, name: 'Козлов Дмитрий', dept: 'DevOps', elo: 2050, wins: 265, losses: 75, trend: -3, avatar: 3 },
        { rank: 4, name: 'Новикова Анна', dept: 'QA', elo: 1980, wins: 234, losses: 88, trend: +15, avatar: 4 },
        { rank: 5, name: 'Морозов Игорь', dept: 'Mobile', elo: 1950, wins: 210, losses: 92, trend: +5, avatar: 5 },
        { rank: 6, name: 'Волков Сергей', dept: 'Data Science', elo: 1920, wins: 198, losses: 102, trend: -7, avatar: 6 },
        { rank: 7, name: 'Соколова Елена', dept: 'Design', elo: 1890, wins: 185, losses: 95, trend: +22, avatar: 7 },
        { rank: 8, name: 'Лебедев Андрей', dept: 'Security', elo: 1860, wins: 178, losses: 112, trend: +3, avatar: 8 },
        { rank: 9, name: 'Кузнецова Ольга', dept: 'HR', elo: 1830, wins: 165, losses: 118, trend: -1, avatar: 9 },
        { rank: 10, name: 'Попов Максим', dept: 'Backend', elo: 1800, wins: 156, losses: 125, trend: +10, avatar: 10 },
        { rank: 11, name: 'Васильев Роман', dept: 'Frontend', elo: 1780, wins: 148, losses: 130, trend: -4, avatar: 11 },
        { rank: 12, name: 'Михайлова Юлия', dept: 'Product', elo: 1750, wins: 142, losses: 135, trend: +6, avatar: 12 },
    ];

    const TOURNAMENTS = [
        {
            id: 1, name: 'Кубок Чемпиона — Шахматы', game: 'chess', status: 'live',
            desc: 'Ежеквартальный турнир по шахматам среди топ-50 игроков',
            prize: '50 000 руб.', players: 48, maxPlayers: 64,
            startsAt: new Date(Date.now() - 2 * 3600000), endsAt: new Date(Date.now() + 4 * 3600000)
        },
        {
            id: 2, name: 'Весенний Турнир — Нарды', game: 'backgammon', status: 'live',
            desc: 'Сезонный турнир по нардам. Регистрация открыта',
            prize: '30 000 руб.', players: 24, maxPlayers: 32,
            startsAt: new Date(Date.now() - 1 * 3600000), endsAt: new Date(Date.now() + 6 * 3600000)
        },
        {
            id: 3, name: 'Лига Шашек', game: 'checkers', status: 'upcoming',
            desc: 'Открытая лига по шашкам. Все уровни приветствуются',
            prize: '20 000 руб.', players: 18, maxPlayers: 32,
            startsAt: new Date(Date.now() + 24 * 3600000), endsAt: new Date(Date.now() + 30 * 3600000)
        },
        {
            id: 4, name: 'Арена Покера — Май 2026', game: 'poker', status: 'upcoming',
            desc: 'Месячный покерный турнир с еженедельными этапами',
            prize: '100 000 руб.', players: 56, maxPlayers: 128,
            startsAt: new Date(Date.now() + 48 * 3600000), endsAt: new Date(Date.now() + 720 * 3600000)
        },
        {
            id: 5, name: 'Снайпер Сапёра', game: 'mines', status: 'upcoming',
            desc: 'Соревнование на скорость прохождения в Сапёре',
            prize: '15 000 руб.', players: 34, maxPlayers: 64,
            startsAt: new Date(Date.now() + 72 * 3600000), endsAt: new Date(Date.now() + 80 * 3600000)
        },
        {
            id: 6, name: 'Зимний Кубок — Шахматы', game: 'chess', status: 'completed',
            desc: 'Завершённый зимний турнир. Результаты доступны',
            prize: '40 000 руб.', players: 64, maxPlayers: 64,
            startsAt: new Date(Date.now() - 720 * 3600000), endsAt: new Date(Date.now() - 600 * 3600000)
        }
    ];

    const LIVE_MATCHES = [
        { id: 1, game: 'chess', player1: 'Петров А.', player2: 'Сидорова М.', time: '12:34', move: 'Ход 24' },
        { id: 2, game: 'checkers', player1: 'Козлов Д.', player2: 'Новикова А.', time: '08:12', move: 'Ход 18' },
        { id: 3, game: 'backgammon', player1: 'Морозов И.', player2: 'Волков С.', time: '15:45', move: 'Партия 3' },
        { id: 4, game: 'poker', player1: 'Соколова Е.', player2: 'Лебедев А.', time: '22:10', move: 'Раунд 5' },
        { id: 5, game: 'chess', player1: 'Кузнецова О.', player2: 'Попов М.', time: '31:20', move: 'Ход 32' },
    ];

    const MATCH_HISTORY = [
        { id: 1, game: 'chess', opponent: 'Петров Алексей', result: 'loss', eloChange: -8, date: 'Сегодня', time: '14:32' },
        { id: 2, game: 'chess', opponent: 'Сидорова Мария', result: 'win', eloChange: +12, date: 'Сегодня', time: '13:15' },
        { id: 3, game: 'checkers', opponent: 'Козлов Дмитрий', result: 'win', eloChange: +6, date: 'Сегодня', time: '11:40' },
        { id: 4, game: 'backgammon', opponent: 'Новикова Анна', result: 'draw', eloChange: +1, date: 'Вчера', time: '17:20' },
        { id: 5, game: 'snake', opponent: 'Турнир: Скорость', result: 'win', eloChange: +0, date: 'Вчера', time: '15:00' },
        { id: 6, game: 'chess', opponent: 'Морозов Игорь', result: 'win', eloChange: +15, date: 'Вчера', time: '12:30' },
        { id: 7, game: 'poker', opponent: 'Турнир: Арена', result: 'loss', eloChange: -20, date: 'Вчера', time: '10:00' },
        { id: 8, game: 'mines', opponent: 'Челлендж: Эксперт', result: 'win', eloChange: +0, date: '2 дня назад', time: '19:15' },
        { id: 9, game: 'chess', opponent: 'Волков Сергей', result: 'loss', eloChange: -5, date: '2 дня назад', time: '14:45' },
        { id: 10, game: 'checkers', opponent: 'Соколова Елена', result: 'win', eloChange: +8, date: '2 дня назад', time: '11:00' },
    ];

    const PROFILE = {
        name: 'Иванов Иван',
        role: 'Frontend Developer',
        department: 'Frontend Team',
        elo: 1450,
        wins: 89,
        losses: 67,
        draws: 12,
        gamesPlayed: 168,
        tournaments: 8,
        bestElo: 1520,
        streak: 3,
        rank: 42,
        joinDate: 'Январь 2026',
        favGame: 'Шахматы'
    };

    // =====================
    // STATE
    // =====================

    let currentPage = 'home';
    let currentFilter = 'all';
    let leaderboardTab = 'active';
    let historyPage = 0;
    const historyPerPage = 10;
    let countdownIntervals = [];

    // =====================
    // UTILITY HELPERS
    // =====================

    function $(sel) { return document.querySelector(sel); }
    function $$(sel) { return document.querySelectorAll(sel); }

    function avatarUrl(seed, size = 40) {
        // Generate SVG avatar from seed
        const hue = (seed * 37) % 360;
        return `data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 40 40'%3E%3Crect width='40' height='40' rx='20' fill='hsl(${hue},30%25,20%25)'/%3E%3Ccircle cx='20' cy='15' r='7' fill='hsl(${hue},40%25,45%25)'/%3E%3Cellipse cx='20' cy='35' rx='12' ry='10' fill='hsl(${hue},40%25,45%25)'/%3E%3C/svg%3E`;
    }

    function eloBarClass(elo) {
        if (elo >= 2000) return 'high';
        if (elo >= 1500) return 'mid';
        return 'low';
    }

    function eloBarWidth(elo) {
        return Math.min(((elo - 800) / 1700) * 100, 100);
    }

    function winRate(wins, losses) {
        const total = wins + losses;
        if (total === 0) return '0%';
        return Math.round((wins / total) * 100) + '%';
    }

    function formatCountdown(date) {
        const now = new Date();
        const diff = date - now;
        if (diff <= 0) return 'Идёт';

        const hours = Math.floor(diff / 3600000);
        const minutes = Math.floor((diff % 3600000) / 60000);
        const seconds = Math.floor((diff % 60000) / 1000);

        return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
    }

    // =====================
    // ROUTER
    // =====================

    function navigate(route) {
        const hash = route.startsWith('#') ? route : '#' + route;
        window.location.hash = hash;
    }

    function handleRoute() {
        const hash = window.location.hash.slice(1) || 'home';
        clearCountdowns();

        let pageId;
        let params = {};

        if (hash.startsWith('game:')) {
            pageId = 'game';
            params = { game: hash.split(':')[1] };
        } else {
            pageId = hash;
        }

        // Check if page exists
        const targetPage = $(`#page-${pageId}`);
        if (!targetPage) {
            pageId = 'home';
            navigate('home');
            return;
        }

        const fromPage = $(`.page.active`);
        if (fromPage === targetPage) return;

        currentPage = pageId;

        GameAnimations.transitionPage(fromPage, targetPage).then(() => {
            renderPageContent(pageId, params);
            updateNavigation(pageId);

            if (typeof lucide !== 'undefined') {
                lucide.createIcons();
            }
        });

        // Save last page
        try { localStorage.setItem('gp_last_page', hash); } catch (e) {}
    }

    function updateNavigation(page) {
        $$('.nav-link').forEach(link => {
            link.classList.toggle('active', link.getAttribute('data-route') === page);
        });
    }

    // =====================
    // RENDER FUNCTIONS
    // =====================

    function renderPageContent(page, params = {}) {
        switch (page) {
            case 'home':
                renderGamesGrid();
                renderTournamentsPreview();
                renderTopPlayers();
                renderLiveTicker();
                break;
            case 'tournaments':
                renderTournamentsList();
                break;
            case 'leaderboard':
                renderLeaderboard();
                break;
            case 'history':
                renderHistory();
                break;
            case 'game':
                renderGamePage(params.game);
                break;
            case 'profile':
                renderProfile();
                break;
        }
    }

    // === Games Grid ===
    function renderGamesGrid(filter = 'all') {
        const grid = $('#games-grid');
        if (!grid) return;

        const filtered = filter === 'all' ? GAMES : GAMES.filter(g => g.category === filter);

        grid.innerHTML = filtered.map(game => `
            <div class="game-card gradient-border" data-game="${game.id}" onclick="Portal.navigate('game:${game.id}')">
                <div class="game-card-header">
                    <div class="game-card-icon">
                        ${GameIcons.gameIcons[game.icon]()}
                    </div>
                    <div class="game-card-badge">
                        <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
                        ${game.players}
                    </div>
                </div>
                <div class="game-card-title">${game.name}</div>
                <div class="game-card-desc">${game.desc}</div>
                <div class="game-card-footer">
                    <div class="game-card-meta">
                        <div class="game-meta-item">
                            <span class="game-meta-value">${game.matchesToday}</span>
                            <span class="game-meta-label">Матчей сегодня</span>
                        </div>
                    </div>
                    <div class="game-card-play tooltip" data-tooltip="Играть">
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                    </div>
                </div>
            </div>
        `).join('');

        // Re-init card effects
        setTimeout(() => GameAnimations.initCardEffects(), 100);
    }

    // === Live Ticker ===
    function renderLiveTicker() {
        const container = $('#live-ticker');
        if (!container) return;

        container.innerHTML = LIVE_MATCHES.map(match => `
            <div class="live-match-card pulse-live">
                <div class="live-match-game-icon">
                    ${GameIcons.gameIcons[match.game]()}
                </div>
                <div class="live-match-info">
                    <div class="live-match-players">
                        ${match.player1}
                        <span class="live-match-vs">vs</span>
                        ${match.player2}
                    </div>
                    <div class="live-match-detail">${match.move}</div>
                </div>
                <div class="live-match-time">${match.time}</div>
            </div>
        `).join('');
    }

    // === Tournaments Preview ===
    function renderTournamentsPreview() {
        const container = $('#tournaments-preview');
        if (!container) return;

        const active = TOURNAMENTS.filter(t => t.status === 'live' || t.status === 'upcoming').slice(0, 4);

        container.innerHTML = active.map(t => renderTournamentCard(t)).join('');
        startCountdowns();
    }

    function renderTournamentCard(t) {
        const game = GAMES.find(g => g.id === t.game);
        const statusLabels = { live: 'В реальном времени', upcoming: 'Предстоящий', completed: 'Завершён' };

        return `
            <div class="tournament-card" onclick="Portal.navigate('tournaments')">
                <div class="tournament-card-header">
                    <div class="tournament-card-game">
                        <div class="tournament-card-game-icon">
                            ${game ? GameIcons.gameIcons[game.icon]() : ''}
                        </div>
                        <span class="tournament-card-game-name">${game ? game.name : t.game}</span>
                    </div>
                    <span class="tournament-status ${t.status}">${statusLabels[t.status]}</span>
                </div>
                <div class="tournament-card-title">${t.name}</div>
                <div class="tournament-card-desc">${t.desc}</div>
                <div class="tournament-card-info">
                    <div class="tournament-info-item">
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
                        <span class="tournament-info-value tournament-prize">${t.prize}</span>
                    </div>
                    <div class="tournament-info-item">
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/></svg>
                        <span class="tournament-info-value">${t.players}/${t.maxPlayers}</span>
                    </div>
                    <div class="tournament-info-item countdown-item" data-end="${t.startsAt.toISOString()}">
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                        <span class="tournament-info-value tournament-countdown">${t.status === 'live' ? 'Идёт' : formatCountdown(t.startsAt)}</span>
                    </div>
                </div>
                <div class="tournament-card-footer">
                    <div class="tournament-players">
                        ${Array.from({length: Math.min(t.players, 4)}, (_, i) => `
                            <div class="tournament-player-avatar">
                                <img src="${avatarUrl(t.id * 10 + i)}" alt="">
                            </div>
                        `).join('')}
                        <span class="tournament-player-count">+${t.players - 4 > 0 ? t.players - 4 : 0}</span>
                    </div>
                    <button class="btn btn-sm btn-gold-outline">Подробнее</button>
                </div>
            </div>
        `;
    }

    // === Tournaments Full List ===
    function renderTournamentsList() {
        const container = $('#tournaments-list');
        if (!container) return;

        const activeTab = leaderboardTab;
        const filtered = activeTab === 'all'
            ? TOURNAMENTS
            : TOURNAMENTS.filter(t => t.status === activeTab);

        if (filtered.length === 0) {
            container.innerHTML = `
                <div class="tournament-card" style="text-align:center; padding: 48px;">
                    <p class="text-muted">Нет турниров в данной категории</p>
                </div>
            `;
            return;
        }

        container.innerHTML = filtered.map(t => renderTournamentCard(t)).join('');

        // Add bracket for first live tournament
        const liveTournament = filtered.find(t => t.status === 'live');
        if (liveTournament) {
            container.innerHTML += renderBracket(liveTournament);
        }

        startCountdowns();
    }

    // === Tournament Bracket (SVG) ===
    function renderBracket(tournament) {
        const participants = [
            { name: 'Петров А.' }, { name: 'Сидорова М.' },
            { name: 'Козлов Д.' }, { name: 'Новикова А.' },
            { name: 'Морозов И.' }, { name: 'Волков С.' },
            { name: 'Соколова Е.' }, { name: 'Лебедев А.' }
        ];

        const roundNames = ['Четвертьфинал', 'Полуфинал', 'Финал'];
        const rounds = [
            // Round 1: 4 matches
            [
                { p1: participants[0], p2: participants[1], winner: 0 },
                { p1: participants[2], p2: participants[3], winner: 1 },
                { p1: participants[4], p2: participants[5], winner: 0 },
                { p1: participants[6], p2: participants[7], winner: 1 },
            ],
            // Round 2: 2 matches
            [
                { p1: participants[0], p2: participants[3], winner: 0 },
                { p1: participants[4], p2: participants[7], winner: null },
            ],
            // Round 3: 1 match
            [
                { p1: participants[0], p2: participants[7], winner: null },
            ]
        ];

        const svgWidth = 800;
        const svgHeight = 380;
        const roundGap = svgWidth / (rounds.length + 1);
        const matchHeight = 70;
        const matchGap = 90;

        let svgPaths = '';

        rounds.forEach((round, ri) => {
            const x = (ri + 0.5) * roundGap;
            const totalHeight = round.length * matchGap;
            const startY = (svgHeight - totalHeight) / 2 + matchHeight / 2;

            round.forEach((match, mi) => {
                const y = startY + mi * matchGap;

                // Connector from previous round
                if (ri > 0) {
                    const prevRound = rounds[ri - 1];
                    const prevMatchIdx1 = mi * 2;
                    const prevMatchIdx2 = mi * 2 + 1;
                    if (prevRound[prevMatchIdx1] && prevRound[prevMatchIdx2]) {
                        const prevX = (ri - 0.5) * roundGap;
                        const prevTotalH = prevRound.length * matchGap;
                        const prevStartY = (svgHeight - prevTotalH) / 2 + matchHeight / 2;
                        const y1 = prevStartY + prevMatchIdx1 * matchGap;
                        const y2 = prevStartY + prevMatchIdx2 * matchGap;
                        const midY = (y1 + y2) / 2;

                        svgPaths += `<path d="M ${prevX + 120} ${y1} H ${x - 20} V ${midY} H ${x}" stroke="#2a2a3e" stroke-width="1.5" fill="none"/>`;
                        svgPaths += `<path d="M ${prevX + 120} ${y2} H ${x - 20} V ${midY} H ${x}" stroke="#2a2a3e" stroke-width="1.5" fill="none"/>`;
                    }
                }

                // Match box
                const isWinnerSet = match.winner !== null;
                const boxColor = isWinnerSet ? 'rgba(212,168,67,0.08)' : 'rgba(20,20,31,0.5)';
                const borderColor = isWinnerSet ? 'rgba(212,168,67,0.3)' : 'rgba(42,42,62,0.8)';

                svgPaths += `
                    <rect x="${x}" y="${y - matchHeight/2}" width="120" height="${matchHeight}" rx="8" fill="${boxColor}" stroke="${borderColor}" stroke-width="1"/>
                `;

                // Player 1
                const p1Color = match.winner === 0 ? '#d4a843' : '#e8e8f0';
                const p1Opacity = match.winner === 1 ? '0.4' : '1';
                svgPaths += `
                    <text x="${x + 10}" y="${y - 5}" fill="${p1Color}" font-size="10" font-weight="600" opacity="${p1Opacity}">${match.p1.name}</text>
                `;

                // Player 2
                const p2Color = match.winner === 1 ? '#d4a843' : '#e8e8f0';
                const p2Opacity = match.winner === 0 ? '0.4' : '1';
                svgPaths += `
                    <text x="${x + 10}" y="${y + 14}" fill="${p2Color}" font-size="10" font-weight="600" opacity="${p2Opacity}">${match.p2.name}</text>
                `;

                // Score
                if (isWinnerSet) {
                    svgPaths += `
                        <text x="${x + 110}" y="${y + 5}" fill="#d4a843" font-size="9" text-anchor="end" font-weight="700">W</text>
                    `;
                }
            });
        });

        return `
            <div class="tournament-bracket">
                <div class="tournament-bracket-title">
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>
                    Сетка турнира: ${tournament.name}
                </div>
                <svg class="bracket-svg" viewBox="0 0 ${svgWidth} ${svgHeight}" preserveAspectRatio="xMidYMid meet">
                    ${svgPaths}
                </svg>
            </div>
        `;
    }

    // === Top Players ===
    function renderTopPlayers() {
        const container = $('#top-players');
        if (!container) return;

        container.innerHTML = LEADERBOARD_DATA.slice(0, 6).map((p, i) => {
            const rankClass = i < 3 ? ['gold', 'silver', 'bronze'][i] : 'default';
            return `
                <div class="top-player-card" onclick="Portal.navigate('profile')">
                    <div class="top-player-rank ${rankClass}">${p.rank}</div>
                    <div class="top-player-avatar">
                        <img src="${avatarUrl(p.avatar)}" alt="${p.name}">
                    </div>
                    <div class="top-player-info">
                        <div class="top-player-name">${p.name}</div>
                        <div class="top-player-dept">${p.dept}</div>
                    </div>
                    <div class="top-player-elo">
                        <div class="top-player-elo-value">${p.elo}</div>
                        <div class="top-player-elo-label">ELO</div>
                    </div>
                </div>
            `;
        }).join('');
    }

    // === Leaderboard ===
    function renderLeaderboard() {
        const body = $('#leaderboard-body');
        if (!body) return;

        const gameFilter = $('#leaderboard-game-filter')?.value || 'all';
        const data = LEADERBOARD_DATA;

        body.innerHTML = data.map(p => {
            const rankClass = p.rank === 1 ? 'gold' : p.rank === 2 ? 'silver' : p.rank === 3 ? 'bronze' : 'default';
            const barWidth = eloBarWidth(p.elo);
            const barClass = eloBarClass(p.elo);
            const wr = winRate(p.wins, p.losses);
            const trendClass = p.trend > 0 ? 'trend-up' : p.trend < 0 ? 'trend-down' : 'trend-neutral';
            const trendIcon = p.trend > 0
                ? `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="18 15 12 9 6 15"/></svg> ${p.trend}`
                : p.trend < 0
                    ? `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg> ${Math.abs(p.trend)}`
                    : `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"/></svg> 0`;

            const isHighlighted = p.name === 'Иванов Иван';

            return `
                <tr class="${isHighlighted ? 'highlight' : ''}" style="animation-delay: ${p.rank * 0.03}s">
                    <td>
                        <div class="leaderboard-rank-cell">
                            <span class="rank-badge ${rankClass}">${p.rank}</span>
                        </div>
                    </td>
                    <td>
                        <div class="leaderboard-player-cell">
                            <div class="leaderboard-player-avatar">
                                <img src="${avatarUrl(p.avatar)}" alt="${p.name}">
                            </div>
                            <span class="leaderboard-player-name">${p.name}</span>
                        </div>
                    </td>
                    <td class="text-muted" style="font-size:13px;">${p.dept}</td>
                    <td>
                        <div class="elo-rating-wrapper">
                            <div class="elo-bar">
                                <div class="elo-bar-fill ${barClass}" data-width="${barWidth}"></div>
                            </div>
                            <span class="elo-value">${p.elo}</span>
                        </div>
                    </td>
                    <td style="text-align:center; color:var(--green); font-weight:600;">${p.wins}</td>
                    <td style="text-align:center; color:var(--red); font-weight:600;">${p.losses}</td>
                    <td style="text-align:center; font-weight:600;">${wr}</td>
                    <td style="text-align:center;">
                        <div class="leaderboard-trend ${trendClass}">${trendIcon}</div>
                    </td>
                </tr>
            `;
        }).join('');

        setTimeout(() => GameAnimations.animateEloBars(), 200);
    }

    // === Match History ===
    function renderHistory() {
        const container = $('#history-timeline');
        if (!container) return;

        const gameFilter = $('#history-game-filter')?.value || 'all';
        const filtered = gameFilter === 'all' ? MATCH_HISTORY : MATCH_HISTORY.filter(m => m.game === gameFilter);

        // Group by date
        const groups = {};
        filtered.forEach(m => {
            if (!groups[m.date]) groups[m.date] = [];
            groups[m.date].push(m);
        });

        container.innerHTML = Object.entries(groups).map(([date, matches]) => `
            <div class="history-date-group">
                <div class="history-date-label">${date}</div>
                ${matches.map(m => {
                    const game = GAMES.find(g => g.id === m.game);
                    const resultLabels = { win: 'Победа', loss: 'Поражение', draw: 'Ничья' };
                    return `
                        <div class="history-item">
                            <div class="history-item-icon">
                                ${game ? GameIcons.gameIcons[game.icon]() : ''}
                            </div>
                            <div class="history-item-content">
                                <div class="history-item-title">${m.opponent}</div>
                                <div class="history-item-sub">${game ? game.name : m.game} ${m.eloChange !== 0 ? (m.eloChange > 0 ? `(+${m.eloChange} ELO)` : `(${m.eloChange} ELO)`) : ''}</div>
                            </div>
                            <div class="history-item-result">
                                <span class="history-result-badge ${m.result}">${resultLabels[m.result]}</span>
                                <div class="history-item-time">${m.time}</div>
                            </div>
                        </div>
                    `;
                }).join('')}
            </div>
        `).join('');
    }

    // === Game Page ===
    function renderGamePage(gameId) {
        const container = $('#game-page-content');
        const title = $('#game-page-title');
        if (!container) return;

        const game = GAMES.find(g => g.id === gameId);
        if (!game) {
            container.innerHTML = '<p class="text-muted">Игра не найдена</p>';
            return;
        }

        title.innerHTML = game.name;

        const backBtn = $('#back-to-home');
        if (backBtn) {
            backBtn.onclick = () => navigate('home');
        }

        container.innerHTML = `
            <div style="display:grid; grid-template-columns: 1fr 1fr; gap: 24px;">
                <div class="game-card gradient-border" style="cursor:default; padding: 40px; text-align:center;">
                    <div class="game-card-icon" style="width:100px; height:100px; margin: 0 auto 24px;">
                        ${GameIcons.gameIcons[game.icon]({ size: 60 })}
                    </div>
                    <h2 style="font-size: 28px; font-weight:800; margin-bottom: 8px;">${game.name}</h2>
                    <p style="color:var(--text-secondary); font-size: 15px; margin-bottom: 24px;">${game.desc}</p>
                    <div class="hero-stats" style="justify-content: center; margin-bottom: 32px;">
                        <div class="hero-stat">
                            <span class="hero-stat-number">${game.players}</span>
                            <span class="hero-stat-label">Игроков</span>
                        </div>
                        <div class="hero-stat">
                            <span class="hero-stat-number">${game.matchesToday}</span>
                            <span class="hero-stat-label">Матчей</span>
                        </div>
                    </div>
                    <button class="btn btn-primary" style="padding: 16px 48px; font-size: 16px;" onclick="Portal.showToast('success', 'Поиск соперника для ${game.name}...')">
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                        Начать игру
                    </button>
                </div>
                <div>
                    <div class="game-card gradient-border" style="cursor:default; margin-bottom: 16px;">
                        <h3 style="font-size: 16px; font-weight: 700; margin-bottom: 16px;">Быстрый матч</h3>
                        <p style="color:var(--text-secondary); font-size: 13px; margin-bottom: 16px;">Найдите случайного соперника с похожим рейтингом</p>
                        <div style="display:flex; gap: 12px;">
                            <button class="btn btn-primary" style="flex:1;" onclick="Portal.showToast('info', 'Поиск соперника...')">Найти соперника</button>
                            <button class="btn btn-secondary" onclick="Portal.showToast('info', 'Создание комнаты...')">Создать комнату</button>
                        </div>
                    </div>
                    <div class="game-card gradient-border" style="cursor:default;">
                        <h3 style="font-size: 16px; font-weight: 700; margin-bottom: 16px;">Последние результаты</h3>
                        <div style="display:flex; flex-direction:column; gap: 8px;">
                            ${MATCH_HISTORY.filter(m => m.game === game.id).slice(0, 3).map(m => `
                                <div style="display:flex; align-items:center; justify-content:space-between; padding: 10px; background:var(--bg-secondary); border-radius:var(--radius-sm);">
                                    <div>
                                        <div style="font-weight:600; font-size:13px;">${m.opponent}</div>
                                        <div style="font-size:11px; color:var(--text-muted);">${m.date} ${m.time}</div>
                                    </div>
                                    <span class="history-result-badge ${m.result}">${m.result === 'win' ? 'Победа' : m.result === 'loss' ? 'Поражение' : 'Ничья'}</span>
                                </div>
                            `).join('')}
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    // === Profile ===
    function renderProfile() {
        const sidebar = $('#profile-sidebar');
        const main = $('#profile-main');
        if (!sidebar || !main) return;

        sidebar.innerHTML = `
            <div class="profile-card">
                <div class="profile-avatar-section">
                    <div class="profile-avatar-large">
                        <img src="${avatarUrl(0, 96)}" alt="Аватар">
                    </div>
                    <div class="profile-name">${PROFILE.name}</div>
                    <div class="profile-role">${PROFILE.role}</div>
                </div>
                <div class="profile-stats-grid">
                    <div class="profile-stat-item">
                        <div class="profile-stat-value">${PROFILE.elo}</div>
                        <div class="profile-stat-label">ELO</div>
                    </div>
                    <div class="profile-stat-item">
                        <div class="profile-stat-value">${PROFILE.gamesPlayed}</div>
                        <div class="profile-stat-label">Игр</div>
                    </div>
                    <div class="profile-stat-item">
                        <div class="profile-stat-value">${PROFILE.wins}</div>
                        <div class="profile-stat-label">Побед</div>
                    </div>
                    <div class="profile-stat-item">
                        <div class="profile-stat-value">${PROFILE.tournaments}</div>
                        <div class="profile-stat-label">Турниров</div>
                    </div>
                </div>
            </div>
            <div class="profile-card">
                <div class="profile-nav-list">
                    <a href="#" class="profile-nav-item active">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
                        Обзор
                    </a>
                    <a href="#" class="profile-nav-item">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
                        Активность
                    </a>
                    <a href="#" class="profile-nav-item">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                        Достижения
                    </a>
                </div>
            </div>
            <div class="profile-card">
                <h4 style="font-size:13px; font-weight:600; margin-bottom:12px; color:var(--text-muted);">Информация</h4>
                <div style="font-size:13px; display:flex; flex-direction:column; gap:8px;">
                    <div style="display:flex; justify-content:space-between;"><span style="color:var(--text-muted);">Отдел</span><span>${PROFILE.department}</span></div>
                    <div style="display:flex; justify-content:space-between;"><span style="color:var(--text-muted);">Ранг</span><span>#${PROFILE.rank}</span></div>
                    <div style="display:flex; justify-content:space-between;"><span style="color:var(--text-muted);">Макс. ELO</span><span class="text-gold">${PROFILE.bestElo}</span></div>
                    <div style="display:flex; justify-content:space-between;"><span style="color:var(--text-muted);">Серия</span><span style="color:var(--green);">${PROFILE.streak} побед</span></div>
                    <div style="display:flex; justify-content:space-between;"><span style="color:var(--text-muted);">Любимая игра</span><span>${PROFILE.favGame}</span></div>
                    <div style="display:flex; justify-content:space-between;"><span style="color:var(--text-muted);">Дата регистрации</span><span>${PROFILE.joinDate}</span></div>
                </div>
            </div>
        `;

        main.innerHTML = `
            <div class="game-card gradient-border" style="cursor:default;">
                <h3 style="font-size: 16px; font-weight: 700; margin-bottom: 20px;">Последние матчи</h3>
                ${MATCH_HISTORY.slice(0, 5).map(m => {
                    const game = GAMES.find(g => g.id === m.game);
                    const resultLabels = { win: 'Победа', loss: 'Поражение', draw: 'Ничья' };
                    return `
                        <div class="history-item" style="margin-bottom:8px;">
                            <div class="history-item-icon">${game ? GameIcons.gameIcons[game.icon]() : ''}</div>
                            <div class="history-item-content">
                                <div class="history-item-title">${m.opponent}</div>
                                <div class="history-item-sub">${game ? game.name : m.game} — ${m.date} ${m.time}</div>
                            </div>
                            <div class="history-item-result">
                                <span class="history-result-badge ${m.result}">${resultLabels[m.result]}</span>
                                ${m.eloChange !== 0 ? `<div class="history-item-time" style="color:${m.eloChange > 0 ? 'var(--green)' : 'var(--red)'};">${m.eloChange > 0 ? '+' : ''}${m.eloChange} ELO</div>` : ''}
                            </div>
                        </div>
                    `;
                }).join('')}
            </div>
            <div class="game-card gradient-border" style="cursor:default;">
                <h3 style="font-size: 16px; font-weight: 700; margin-bottom: 20px;">Статистика по играм</h3>
                <div style="display:grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 12px;">
                    ${GAMES.slice(0, 4).map(game => `
                        <div style="padding:16px; background:var(--bg-secondary); border-radius:var(--radius-md);">
                            <div style="display:flex; align-items:center; gap:8px; margin-bottom:10px;">
                                ${GameIcons.gameIcons[game.icon]({ size: 20 })}
                                <span style="font-weight:600; font-size:13px;">${game.name}</span>
                            </div>
                            <div style="display:flex; justify-content:space-between; font-size:12px; color:var(--text-secondary);">
                                <span>${Math.floor(Math.random() * 30 + 10)} игр</span>
                                <span style="color:var(--green);">${Math.floor(Math.random() * 20 + 50)}% побед</span>
                            </div>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }

    // =====================
    // COUNTDOWN TIMERS
    // =====================

    function startCountdowns() {
        clearCountdowns();

        $$('.countdown-item').forEach(el => {
            const endStr = el.getAttribute('data-end');
            const endDate = new Date(endStr);
            const countdownSpan = el.querySelector('.tournament-countdown');
            if (!countdownSpan) return;

            const interval = setInterval(() => {
                countdownSpan.textContent = formatCountdown(endDate);
            }, 1000);

            countdownIntervals.push(interval);
        });
    }

    function clearCountdowns() {
        countdownIntervals.forEach(id => clearInterval(id));
        countdownIntervals = [];
    }

    // =====================
    // TOAST SYSTEM
    // =====================

    function showToast(type, message, duration = 4000) {
        const container = $('#toast-container');
        if (!container) return;

        const icons = {
            success: '<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>',
            error: '<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>',
            info: '<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg>',
            warning: '<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>'
        };

        const toast = document.createElement('div');
        toast.className = 'toast';
        toast.innerHTML = `
            <span class="toast-icon ${type}">${icons[type] || icons.info}</span>
            <span>${message}</span>
        `;

        container.appendChild(toast);

        setTimeout(() => {
            toast.classList.add('exit');
            setTimeout(() => toast.remove(), 300);
        }, duration);
    }

    // =====================
    // SEARCH
    // =====================

    function initSearch() {
        const input = $('#global-search');
        const results = $('#search-results');
        if (!input || !results) return;

        let debounceTimer;

        input.addEventListener('input', () => {
            clearTimeout(debounceTimer);
            const query = input.value.trim().toLowerCase();

            if (query.length < 2) {
                results.classList.remove('visible');
                return;
            }

            debounceTimer = setTimeout(() => {
                // Search games
                const gameResults = GAMES.filter(g =>
                    g.name.toLowerCase().includes(query) || g.desc.toLowerCase().includes(query)
                ).map(g => ({
                    type: 'game',
                    icon: g.icon,
                    title: g.name,
                    sub: g.desc,
                    link: `#game:${g.id}`
                }));

                // Search players
                const playerResults = LEADERBOARD_DATA.filter(p =>
                    p.name.toLowerCase().includes(query) || p.dept.toLowerCase().includes(query)
                ).map(p => ({
                    type: 'player',
                    icon: 'user',
                    title: p.name,
                    sub: `${p.dept} — ELO ${p.elo}`,
                    link: '#profile'
                }));

                const allResults = [...gameResults, ...playerResults];

                if (allResults.length === 0) {
                    results.innerHTML = '<div style="padding:16px; text-align:center; color:var(--text-muted); font-size:13px;">Ничего не найдено</div>';
                } else {
                    results.innerHTML = allResults.map(r => `
                        <div class="search-result-item" onclick="Portal.navigate('${r.link}'); document.getElementById('search-results').classList.remove('visible'); document.getElementById('global-search').value = '';">
                            <div class="search-result-icon" style="color:var(--gold);">
                                ${r.type === 'game' && GameIcons.gameIcons[r.icon] ? GameIcons.gameIcons[r.icon]({ size: 18 }) : '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>'}
                            </div>
                            <div>
                                <div class="search-result-text">${r.title}</div>
                                <div class="search-result-sub">${r.sub}</div>
                            </div>
                        </div>
                    `).join('');
                }

                results.classList.add('visible');
            }, 300);
        });

        // Close search results on click outside
        document.addEventListener('click', (e) => {
            if (!e.target.closest('.search-wrapper')) {
                results.classList.remove('visible');
            }
        });
    }

    // =====================
    // KEYBOARD SHORTCUTS
    // =====================

    function initShortcuts() {
        document.addEventListener('keydown', (e) => {
            // Don't trigger in input/textarea
            if (['INPUT', 'TEXTAREA', 'SELECT'].includes(e.target.tagName)) return;

            switch (e.key) {
                case '/':
                    e.preventDefault();
                    const search = $('#global-search');
                    if (search) search.focus();
                    break;
                case '?':
                    e.preventDefault();
                    toggleShortcutsModal();
                    break;
                case 'Escape':
                    closeShortcutsModal();
                    const results = $('#search-results');
                    if (results) results.classList.remove('visible');
                    const dropdown = $('#user-dropdown');
                    if (dropdown) dropdown.classList.remove('visible');
                    break;
            }
        });
    }

    function toggleShortcutsModal() {
        let overlay = $('.shortcuts-modal-overlay');
        if (overlay) {
            overlay.remove();
            return;
        }

        overlay = document.createElement('div');
        overlay.className = 'shortcuts-modal-overlay visible';
        overlay.onclick = (e) => { if (e.target === overlay) closeShortcutsModal(); };

        overlay.innerHTML = `
            <div class="shortcuts-modal">
                <h3>
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="M6 8h.001M10 8h.001M14 8h.001M18 8h.001M8 12h.001M12 12h.001M16 12h.001M7 16h10"/></svg>
                    Горячие клавиши
                </h3>
                <div class="shortcut-row">
                    <span class="shortcut-action">Поиск</span>
                    <div class="shortcut-keys"><kbd>/</kbd></div>
                </div>
                <div class="shortcut-row">
                    <span class="shortcut-action">Горячие клавиши</span>
                    <div class="shortcut-keys"><kbd>?</kbd></div>
                </div>
                <div class="shortcut-row">
                    <span class="shortcut-action">Закрыть / Назад</span>
                    <div class="shortcut-keys"><kbd>Esc</kbd></div>
                </div>
                <div class="shortcut-row">
                    <span class="shortcut-action">Главная</span>
                    <div class="shortcut-keys"><kbd>Alt</kbd><kbd>1</kbd></div>
                </div>
                <div class="shortcut-row">
                    <span class="shortcut-action">Турниры</span>
                    <div class="shortcut-keys"><kbd>Alt</kbd><kbd>2</kbd></div>
                </div>
                <div class="shortcut-row">
                    <span class="shortcut-action">Рейтинг</span>
                    <div class="shortcut-keys"><kbd>Alt</kbd><kbd>3</kbd></div>
                </div>
            </div>
        `;

        document.body.appendChild(overlay);
    }

    function closeShortcutsModal() {
        const overlay = $('.shortcuts-modal-overlay');
        if (overlay) overlay.remove();
    }

    // =====================
    // HEADER SCROLL EFFECT
    // =====================

    function initHeaderScroll() {
        const header = $('#main-header');
        if (!header) return;

        window.addEventListener('scroll', () => {
            header.classList.toggle('scrolled', window.scrollY > 10);
        }, { passive: true });
    }

    // =====================
    // MOBILE MENU
    // =====================

    function initMobileMenu() {
        const btn = $('#mobile-menu-btn');
        const nav = $('#main-nav');
        const overlay = $('#mobile-nav-overlay');

        if (!btn || !nav) return;

        btn.addEventListener('click', () => {
            nav.classList.toggle('mobile-open');
            if (overlay) overlay.classList.toggle('visible');
        });

        if (overlay) {
            overlay.addEventListener('click', () => {
                nav.classList.remove('mobile-open');
                overlay.classList.remove('visible');
            });
        }

        // Close mobile nav on link click
        $$('.nav-link').forEach(link => {
            link.addEventListener('click', () => {
                nav.classList.remove('mobile-open');
                if (overlay) overlay.classList.remove('visible');
            });
        });
    }

    // =====================
    // USER DROPDOWN
    // =====================

    function initUserDropdown() {
        const avatar = $('#user-avatar');
        const dropdown = $('#user-dropdown');

        if (!avatar || !dropdown) return;

        avatar.addEventListener('click', (e) => {
            e.stopPropagation();
            dropdown.classList.toggle('visible');
        });

        document.addEventListener('click', (e) => {
            if (!e.target.closest('.user-menu')) {
                dropdown.classList.remove('visible');
            }
        });
    }

    // =====================
    // FILTER BUTTONS
    // =====================

    function initFilters() {
        // Game category filters
        $$('.filter-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                $$('.filter-btn').forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                currentFilter = btn.getAttribute('data-filter');
                renderGamesGrid(currentFilter);
            });
        });

        // Tournament tabs
        $$('.tab-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                $$('.tab-btn').forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                leaderboardTab = btn.getAttribute('data-tab');
                renderTournamentsList();
            });
        });

        // Leaderboard filters
        const gameFilter = $('#leaderboard-game-filter');
        if (gameFilter) {
            gameFilter.addEventListener('change', () => renderLeaderboard());
        }

        const periodFilter = $('#leaderboard-period');
        if (periodFilter) {
            periodFilter.addEventListener('change', () => renderLeaderboard());
        }

        // History filter
        const histFilter = $('#history-game-filter');
        if (histFilter) {
            histFilter.addEventListener('change', () => renderHistory());
        }

        // Load more
        const loadMore = $('#load-more-btn');
        if (loadMore) {
            loadMore.addEventListener('click', () => {
                showToast('info', 'Загрузка дополнительных матчей...');
            });
        }
    }

    // =====================
    // API CLIENT (stub with offline queue)
    // =====================

    const ApiClient = {
        queue: [],
        isOnline: navigator.onLine,

        async request(url, options = {}) {
            if (!this.isOnline) {
                this.queue.push({ url, options, timestamp: Date.now() });
                showToast('warning', 'Вы офлайн. Запрос добавлен в очередь.');
                return null;
            }

            try {
                // In production, this would make real API calls
                // For now, simulate a delay and return success
                await new Promise(resolve => setTimeout(resolve, 200));
                return { ok: true, data: null };
            } catch (err) {
                showToast('error', `Ошибка запроса: ${err.message}`);
                return { ok: false, error: err.message };
            }
        },

        async processQueue() {
            if (this.queue.length === 0) return;
            this.isOnline = true;
            const items = [...this.queue];
            this.queue = [];

            for (const item of items) {
                await this.request(item.url, item.options);
            }

            if (items.length > 0) {
                showToast('success', `Обработано ${items.length} отложенных запросов`);
            }
        }
    };

    // Online/offline listeners
    window.addEventListener('online', () => ApiClient.processQueue());
    window.addEventListener('offline', () => { ApiClient.isOnline = false; });

    // =====================
    // LOCALSTORAGE PREFERENCES
    // =====================

    function loadPreferences() {
        try {
            const lastPage = localStorage.getItem('gp_last_page');
            if (lastPage && lastPage !== 'home') {
                // Don't auto-navigate, just remember it
            }
        } catch (e) {}
    }

    // =====================
    // INITIALIZATION
    // =====================

    function init() {
        // Wait for DOM + Lucide
        const checkReady = () => {
            if (typeof lucide !== 'undefined') {
                start();
            } else {
                setTimeout(checkReady, 50);
            }
        };

        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', checkReady);
        } else {
            checkReady();
        }
    }

    function start() {
        // Init Lucide icons
        lucide.createIcons();

        // Init animations
        GameAnimations.init();

        // Hide loading screen after everything is set up
        setTimeout(() => {
            GameAnimations.hideLoadingScreen();
        }, 800);

        // Init interactions
        initSearch();
        initShortcuts();
        initHeaderScroll();
        initMobileMenu();
        initUserDropdown();
        initFilters();
        loadPreferences();

        // Render initial page content
        renderGamesGrid();
        renderTournamentsPreview();
        renderTopPlayers();
        renderLiveTicker();
        startCountdowns();

        // Setup hash routing
        window.addEventListener('hashchange', handleRoute);

        // Handle initial route
        handleRoute();

        // Welcome toast
        setTimeout(() => {
            showToast('info', 'Добро пожаловать на игровой портал!');
        }, 1500);

        // Periodic mock: live match update
        setInterval(() => {
            showToast('success', 'Турнирное обновление: Петров А. выиграл партию!', 5000);
        }, 45000);
    }

    // =====================
    // PUBLIC API
    // =====================

    return {
        init,
        navigate,
        showToast,
        launchConfetti: () => GameAnimations.launchConfetti()
    };
})();

// Auto-init when script loads
Portal.init();

// Make available globally
window.Portal = Portal;
