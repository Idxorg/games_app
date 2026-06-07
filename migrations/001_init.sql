-- Таблица пользователей (с полями из корп портала)
CREATE TABLE users (
    sid VARCHAR(50) PRIMARY KEY,  -- sid от SSO (emp_12345)
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    gender VARCHAR(20),  -- 'male' | 'female' | 'other' (из корп портала)
    department VARCHAR(100),
    position VARCHAR(100),
    photo_url VARCHAR(500),  -- URL в S3
    last_sync TIMESTAMP DEFAULT NOW(),  -- Последняя синхронизация из портала
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Группы пользователей (для контроля доступа)
CREATE TABLE user_groups (
    id SERIAL PRIMARY KEY,
    sid VARCHAR(50) REFERENCES users(sid),
    group_name VARCHAR(100) NOT NULL,  -- 'games', 'tournaments', 'admin'
    granted_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(sid, group_name)
);

-- Рейтинги (по SID)
CREATE TABLE player_ratings (
    id SERIAL PRIMARY KEY,
    sid VARCHAR(50) REFERENCES users(sid),
    game_type VARCHAR(50) NOT NULL,  -- chess, checkers, backgammon
    elo INTEGER DEFAULT 1200,
    games_played INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0,
    draws INTEGER DEFAULT 0,
    losses INTEGER DEFAULT 0,
    UNIQUE(sid, game_type)
);

-- Турниры (с логотипами в S3)
CREATE TABLE tournaments (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    game_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    max_players INTEGER DEFAULT 64,
    current_players INTEGER DEFAULT 0,
    prize_pool VARCHAR(100),
    description TEXT,
    logo_url VARCHAR(500),  -- URL в S3
    created_by VARCHAR(50) REFERENCES users(sid),
    requires_group VARCHAR(100),  -- Требуемая группа для участия (опционально)
    created_at TIMESTAMP DEFAULT NOW()
);

-- Участники турниров
CREATE TABLE tournament_players (
    id SERIAL PRIMARY KEY,
    tournament_id VARCHAR(50) REFERENCES tournaments(id),
    sid VARCHAR(50) REFERENCES users(sid),
    rank INTEGER,
    points INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0,
    draws INTEGER DEFAULT 0,
    losses INTEGER DEFAULT 0,
    joined_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tournament_id, sid)
);

-- Матчи (с записями в S3)
CREATE TABLE matches (
    id VARCHAR(50) PRIMARY KEY,
    tournament_id VARCHAR(50) REFERENCES tournaments(id),
    game_type VARCHAR(50) NOT NULL,
    player1_sid VARCHAR(50) REFERENCES users(sid),
    player2_sid VARCHAR(50) REFERENCES users(sid),
    winner_sid VARCHAR(50),  -- NULL если ничья
    score VARCHAR(20),  -- "1-0", "0-1", "1/2-1/2"
    moves JSONB,  -- Ходы в формате JSON
    pgn_url VARCHAR(500),  -- URL в S3 с полной записью
    game_id VARCHAR(100),  -- WebSocket game_id
    livekit_room_id VARCHAR(100),  -- ID комнаты LiveKit (опционально)
    status VARCHAR(20) DEFAULT 'waiting',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для производительности
CREATE INDEX idx_player_ratings_game ON player_ratings(game_type);
CREATE INDEX idx_tournament_players_tournament ON tournament_players(tournament_id);
CREATE INDEX idx_matches_tournament ON matches(tournament_id);
CREATE INDEX idx_matches_status ON matches(status);
CREATE INDEX idx_matches_player1 ON matches(player1_sid);
CREATE INDEX idx_matches_player2 ON matches(player2_sid);
CREATE INDEX idx_user_groups_group ON user_groups(group_name);
