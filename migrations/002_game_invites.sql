CREATE TABLE IF NOT EXISTS game_invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_type VARCHAR(20) NOT NULL CHECK (game_type IN ('chess','checkers','backgammon')),
    inviter_sid VARCHAR(100) NOT NULL REFERENCES users(sid),
    recipient_sid VARCHAR(100) NOT NULL REFERENCES users(sid),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','accepted','declined','expired')),
    match_id UUID REFERENCES matches(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '5 minutes'
);

CREATE INDEX idx_invites_recipient ON game_invites(recipient_sid, status);
CREATE INDEX idx_invites_inviter ON game_invites(inviter_sid);
