// Базовый JS для игрового портала
const API_BASE = '/api/v1';
let authToken = null;
let currentUser = null;

// Инициализация при загрузке
document.addEventListener('DOMContentLoaded', () => {
    initAuth();
    loadGames();
    loadTournaments();
    loadLeaderboard();
});

// Инициализация аутентификации
async function initAuth() {
    // Получаем токен из localStorage (будет установлено корп порталом)
    authToken = localStorage.getItem('auth_token');
    
    if (authToken) {
        await verifyToken();
    }
}

// Проверка токена
async function verifyToken() {
    try {
        const response = await fetch(`${API_BASE}/auth/verify`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json'
            }
        });
        
        if (response.ok) {
            const data = await response.json();
            console.log('Token verified:', data);
        } else {
            console.error('Token verification failed');
            localStorage.removeItem('auth_token');
            authToken = null;
        }
    } catch (error) {
        console.error('Token verification error:', error);
    }
}

// Загрузка доступных игр
async function loadGames() {
    try {
        const response = await fetch(`${API_BASE}/games/available`, {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });
        
        if (response.ok) {
            const games = await response.json();
            console.log('Available games:', games);
        }
    } catch (error) {
        console.error('Failed to load games:', error);
    }
}

// Загрузка турниров
async function loadTournaments() {
    try {
        const response = await fetch(`${API_BASE}/tournaments`, {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });
        
        if (response.ok) {
            const tournaments = await response.json();
            displayTournaments(tournaments);
        }
    } catch (error) {
        console.error('Failed to load tournaments:', error);
    }
}

// Отображение турниров
function displayTournaments(tournaments) {
    const container = document.getElementById('tournaments-list');
    container.innerHTML = '';
    
    tournaments.forEach(tournament => {
        const div = document.createElement('div');
        div.className = 'tournament-item';
        div.innerHTML = `
            <h4>${tournament.name}</h4>
            <p>Игра: ${tournament.game_type}</p>
            <p>Дата: ${new Date(tournament.start_date).toLocaleDateString()} - ${new Date(tournament.end_date).toLocaleDateString()}</p>
            <p>Участников: ${tournament.current_players}/${tournament.max_players}</p>
            <button onclick="joinTournament('${tournament.id}')">Записаться</button>
        `;
        container.appendChild(div);
    });
}

// Запись на турнир
async function joinTournament(tournamentId) {
    try {
        const response = await fetch(`${API_BASE}/tournaments/${tournamentId}/join`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json'
            }
        });
        
        if (response.ok) {
            alert('Вы успешно записались на турнир!');
            loadTournaments();
        } else {
            alert('Ошибка при записи на турнир');
        }
    } catch (error) {
        console.error('Failed to join tournament:', error);
        alert('Ошибка при записи на турнир');
    }
}

// Загрузка лидерборда
async function loadLeaderboard() {
    try {
        const response = await fetch(`${API_BASE}/ratings/chess`, {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });
        
        if (response.ok) {
            const ratings = await response.json();
            displayLeaderboard(ratings);
        }
    } catch (error) {
        console.error('Failed to load leaderboard:', error);
    }
}

// Отображение лидерборда
function displayLeaderboard(ratings) {
    const container = document.getElementById('leaderboard-content');
    container.innerHTML = '';
    
    ratings.slice(0, 10).forEach((player, index) => {
        const div = document.createElement('div');
        div.className = 'leaderboard-item';
        div.innerHTML = `
            <h4>#${index + 1} ${player.name}</h4>
            <p>Отдел: ${player.department}</p>
            <p>Рейтинг: ${player.elo}</p>
            <p>Побед: ${player.wins} | Поражений: ${player.losses}</p>
        `;
        container.appendChild(div);
    });
}

// Запуск игры
function startGame(gameType) {
    console.log(`Starting game: ${gameType}`);
    // TODO: Реализовать запуск игры через WebSocket
    alert(`Игра ${gameType} будет запущена в ближайшее время!`);
}
