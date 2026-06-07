/* ============================================
   GamePortal - Animations Library
   CSS + JS animations and visual effects
   ============================================ */

const GameAnimations = (() => {
    'use strict';

    // =====================
    // CHESS BOARD PIECES ANIMATION
    // =====================

    function animateChessPieces() {
        const pieces = document.querySelectorAll('.chess-square svg');
        pieces.forEach((piece, index) => {
            piece.style.animationDelay = `${index * 0.04 + 0.5}s`;
        });
    }

    // =====================
    // COUNTING ANIMATION
    // =====================

    function animateCountUp(element, target, duration = 1500) {
        const start = 0;
        const startTime = performance.now();
        const text = element.textContent;

        function update(currentTime) {
            const elapsed = currentTime - startTime;
            const progress = Math.min(elapsed / duration, 1);

            // Ease out cubic
            const eased = 1 - Math.pow(1 - progress, 3);
            const current = Math.round(start + (target - start) * eased);

            element.textContent = current.toLocaleString('ru-RU');

            if (progress < 1) {
                requestAnimationFrame(update);
            }
        }

        requestAnimationFrame(update);
    }

    function initCountUp() {
        document.querySelectorAll('[data-count]').forEach(el => {
            const target = parseInt(el.getAttribute('data-count'), 10);
            if (!isNaN(target)) {
                // Use IntersectionObserver
                const observer = new IntersectionObserver((entries) => {
                    entries.forEach(entry => {
                        if (entry.isIntersecting) {
                            animateCountUp(el, target);
                            observer.unobserve(el);
                        }
                    });
                }, { threshold: 0.3 });
                observer.observe(el);
            }
        });
    }

    // =====================
    // CONFETTI EFFECT
    // =====================

    function launchConfetti(duration = 3000) {
        const canvas = document.createElement('canvas');
        canvas.className = 'confetti-canvas';
        document.body.appendChild(canvas);

        const ctx = canvas.getContext('2d');
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;

        const colors = ['#d4a843', '#e8c36a', '#ffd700', '#b8922f', '#fff', '#4ade80', '#60a5fa'];
        const particles = [];

        for (let i = 0; i < 150; i++) {
            particles.push({
                x: Math.random() * canvas.width,
                y: Math.random() * canvas.height - canvas.height,
                w: Math.random() * 8 + 4,
                h: Math.random() * 6 + 2,
                color: colors[Math.floor(Math.random() * colors.length)],
                vx: (Math.random() - 0.5) * 4,
                vy: Math.random() * 3 + 2,
                rotation: Math.random() * 360,
                rotSpeed: (Math.random() - 0.5) * 10,
                opacity: 1
            });
        }

        const startTime = performance.now();

        function draw(currentTime) {
            const elapsed = currentTime - startTime;
            ctx.clearRect(0, 0, canvas.width, canvas.height);

            if (elapsed > duration) {
                const fadeOut = 1 - ((elapsed - duration) / 500);
                if (fadeOut <= 0) {
                    canvas.remove();
                    return;
                }
                ctx.globalAlpha = fadeOut;
            }

            particles.forEach(p => {
                p.x += p.vx;
                p.y += p.vy;
                p.vy += 0.05;
                p.rotation += p.rotSpeed;

                ctx.save();
                ctx.translate(p.x, p.y);
                ctx.rotate((p.rotation * Math.PI) / 180);
                ctx.fillStyle = p.color;
                ctx.fillRect(-p.w / 2, -p.h / 2, p.w, p.h);
                ctx.restore();
            });

            requestAnimationFrame(draw);
        }

        requestAnimationFrame(draw);

        window.addEventListener('resize', () => {
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;
        });
    }

    // =====================
    // RATING CHANGE ANIMATION
    // =====================

    function animateRatingChange(element, change) {
        const el = document.createElement('div');
        el.className = `rating-change ${change > 0 ? 'positive' : 'negative'}`;
        el.textContent = `${change > 0 ? '+' : ''}${change}`;

        element.style.position = 'relative';
        element.appendChild(el);

        setTimeout(() => el.remove(), 1600);
    }

    // =====================
    // PARTICLE BACKGROUND
    // =====================

    function initParticles() {
        const canvas = document.getElementById('particles-canvas');
        if (!canvas) return;

        const ctx = canvas.getContext('2d');
        let particles = [];
        let animId;

        function resize() {
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;
        }

        resize();
        window.addEventListener('resize', resize);

        // Create chess-piece shaped particles (simplified as small geometric shapes)
        const shapes = ['circle', 'diamond', 'square', 'triangle'];

        for (let i = 0; i < 30; i++) {
            particles.push({
                x: Math.random() * canvas.width,
                y: Math.random() * canvas.height,
                size: Math.random() * 3 + 1,
                shape: shapes[Math.floor(Math.random() * shapes.length)],
                opacity: Math.random() * 0.3 + 0.1,
                vx: (Math.random() - 0.5) * 0.3,
                vy: (Math.random() - 0.5) * 0.3,
                color: Math.random() > 0.7 ? '#d4a843' : '#ffffff'
            });
        }

        function drawShape(p) {
            ctx.globalAlpha = p.opacity;
            ctx.fillStyle = p.color;

            switch (p.shape) {
                case 'circle':
                    ctx.beginPath();
                    ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2);
                    ctx.fill();
                    break;
                case 'diamond':
                    ctx.beginPath();
                    ctx.moveTo(p.x, p.y - p.size);
                    ctx.lineTo(p.x + p.size, p.y);
                    ctx.lineTo(p.x, p.y + p.size);
                    ctx.lineTo(p.x - p.size, p.y);
                    ctx.closePath();
                    ctx.fill();
                    break;
                case 'square':
                    ctx.fillRect(p.x - p.size, p.y - p.size, p.size * 2, p.size * 2);
                    break;
                case 'triangle':
                    ctx.beginPath();
                    ctx.moveTo(p.x, p.y - p.size);
                    ctx.lineTo(p.x + p.size, p.y + p.size);
                    ctx.lineTo(p.x - p.size, p.y + p.size);
                    ctx.closePath();
                    ctx.fill();
                    break;
            }
        }

        function animate() {
            ctx.clearRect(0, 0, canvas.width, canvas.height);

            particles.forEach(p => {
                p.x += p.vx;
                p.y += p.vy;

                // Wrap around
                if (p.x < -10) p.x = canvas.width + 10;
                if (p.x > canvas.width + 10) p.x = -10;
                if (p.y < -10) p.y = canvas.height + 10;
                if (p.y > canvas.height + 10) p.y = -10;

                drawShape(p);
            });

            ctx.globalAlpha = 1;
            animId = requestAnimationFrame(animate);
        }

        // Respect reduced motion preference
        if (!window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
            animate();
        }
    }

    // =====================
    // PARALLAX EFFECT
    // =====================

    function initParallax() {
        const hero = document.getElementById('hero-section');
        const heroVisual = document.getElementById('hero-visual');
        if (!hero || !heroVisual) return;

        window.addEventListener('scroll', () => {
            const scrolled = window.pageYOffset;
            const rate = scrolled * 0.15;
            if (heroVisual) {
                heroVisual.style.transform = `translateY(${rate}px)`;
            }
        }, { passive: true });
    }

    // =====================
    // CARD HOVER EFFECTS
    // =====================

    function initCardEffects() {
        document.querySelectorAll('.game-card').forEach(card => {
            card.addEventListener('mousemove', (e) => {
                const rect = card.getBoundingClientRect();
                const x = e.clientX - rect.left;
                const y = e.clientY - rect.top;
                const centerX = rect.width / 2;
                const centerY = rect.height / 2;

                const rotateX = ((y - centerY) / centerY) * -3;
                const rotateY = ((x - centerX) / centerX) * 3;

                card.style.transform = `perspective(800px) rotateX(${rotateX}deg) rotateY(${rotateY}deg) translateY(-4px)`;
            });

            card.addEventListener('mouseleave', () => {
                card.style.transform = '';
            });
        });
    }

    // =====================
    // LOADING SCREEN
    // =====================

    function hideLoadingScreen() {
        const loading = document.getElementById('loading-screen');
        const app = document.querySelector('.app-wrapper');

        if (loading) {
            loading.classList.add('hidden');
            setTimeout(() => loading.remove(), 600);
        }

        if (app) {
            setTimeout(() => app.classList.add('loaded'), 100);
        }
    }

    // =====================
    // PAGE TRANSITION
    // =====================

    function transitionPage(fromPage, toPage) {
        if (fromPage) {
            fromPage.style.animation = 'page-exit 0.3s ease forwards';
        }

        return new Promise(resolve => {
            setTimeout(() => {
                if (fromPage) {
                    fromPage.classList.remove('active');
                    fromPage.style.animation = '';
                }
                toPage.classList.add('active');
                toPage.style.animation = 'page-enter 0.4s cubic-bezier(0.16, 1, 0.3, 1)';

                // Scroll to top
                window.scrollTo({ top: 0, behavior: 'smooth' });

                setTimeout(resolve, 400);
            }, fromPage ? 200 : 0);
        });
    }

    // =====================
    // ELO BAR ANIMATION
    // =====================

    function animateEloBars() {
        document.querySelectorAll('.elo-bar-fill').forEach(bar => {
            const width = bar.getAttribute('data-width');
            if (width) {
                bar.style.width = '0%';
                setTimeout(() => {
                    bar.style.width = width + '%';
                }, 200);
            }
        });
    }

    // =====================
    // HERO CHESS BOARD BUILDER
    // =====================

    function buildHeroChessBoard() {
        const container = document.getElementById('hero-visual');
        if (!container) return;

        // Starting position (row by row, top to bottom: black pieces, empty, white pawns, empty, white pieces)
        const backRow = ['rook', 'knight', 'bishop', 'queen', 'king', 'bishop', 'knight', 'rook'];
        const pieces = [
            // Black back rank
            ...backRow.map(p => ({ type: p, color: 'black' })),
            // Black pawns
            ...Array(8).fill(null).map(() => ({ type: 'pawn', color: 'black' })),
            // Empty rows
            ...Array(16).fill(null),
            // White pawns
            ...Array(8).fill(null).map(() => ({ type: 'pawn', color: 'white' })),
            // White back rank
            ...backRow.map(p => ({ type: p, color: 'white' })),
        ];

        const board = document.createElement('div');
        board.className = 'hero-chess-board';

        for (let i = 0; i < 64; i++) {
            const row = Math.floor(i / 8);
            const col = i % 8;
            const isLight = (row + col) % 2 === 0;
            const square = document.createElement('div');
            square.className = `chess-square ${isLight ? 'light' : 'dark'}`;

            const piece = pieces[i];
            if (piece && GameIcons.chess[piece.type]) {
                const pieceDiv = document.createElement('div');
                pieceDiv.className = `chess-piece-${piece.color}`;
                pieceDiv.innerHTML = GameIcons.chess[piece.type]();
                square.appendChild(pieceDiv);
            }

            board.appendChild(square);
        }

        container.appendChild(board);
        animateChessPieces();
    }

    // =====================
    // TOAST ANIMATION HELPERS
    // =====================

    function animateToastIn(toast) {
        toast.style.animation = 'toast-enter 0.3s cubic-bezier(0.16, 1, 0.3, 1)';
    }

    function animateToastOut(toast) {
        toast.classList.add('exit');
    }

    // =====================
    // SKELETON LOADING
    // =====================

    function showSkeleton(container, count = 3) {
        container.innerHTML = '';
        for (let i = 0; i < count; i++) {
            const skeleton = document.createElement('div');
            skeleton.className = 'skeleton skeleton-card';
            skeleton.innerHTML = `
                <div class="skeleton skeleton-avatar" style="margin: 16px;"></div>
                <div class="skeleton skeleton-title" style="margin: 0 16px;"></div>
                <div class="skeleton skeleton-text" style="margin: 0 16px;"></div>
                <div class="skeleton skeleton-text short" style="margin: 0 16px;"></div>
            `;
            skeleton.style.height = '120px';
            skeleton.style.marginBottom = '16px';
            skeleton.style.padding = '16px';
            skeleton.style.borderRadius = '16px';
            container.appendChild(skeleton);
        }
    }

    // =====================
    // NUMBER SCROLL ANIMATION (for stats)
    // =====================

    function initStatAnimations() {
        document.querySelectorAll('.profile-stat-value, .hero-stat-number').forEach(el => {
            const target = parseInt(el.textContent.replace(/[^\d]/g, ''), 10);
            if (!isNaN(target) && target > 0) {
                el.setAttribute('data-count', target);
                el.textContent = '0';
            }
        });
    }

    // =====================
    // INITIALIZATION
    // =====================

    function init() {
        buildHeroChessBoard();
        initCountUp();
        initParticles();
        initParallax();

        // Delay card effects to let DOM settle
        setTimeout(() => {
            initCardEffects();
            animateEloBars();
        }, 500);
    }

    return {
        init,
        animateChessPieces,
        animateCountUp,
        launchConfetti,
        animateRatingChange,
        initParticles,
        initParallax,
        initCardEffects,
        hideLoadingScreen,
        transitionPage,
        animateEloBars,
        animateToastIn,
        animateToastOut,
        showSkeleton,
        buildHeroChessBoard
    };
})();

// Make available globally
window.GameAnimations = GameAnimations;
