# CHANGELOG — Игры · ЭР-Линк (games_app)

## [2026-06-08] — G0 → G1a → G1b → G2 → G4 Complete

### New Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NOTIFICATIONS_API_URL` | Notifications API endpoint for game events | _(empty, noop)_ |
| `NOTIFICATIONS_INTERNAL_TOKEN` | Bearer token for notifications API | _(empty)_ |
| `GAMES_EMBED_HANDOFF_SECRET` | Shared secret for iframe embed auth | _(required in prod)_ |

### Breaking Changes

- **WS path changed**: `/api/v1/ws/game/:match_id` → `/ws/game/:match_id`
- **Invite API paths**: `/api/v1/games/invite` → `/games/invite` (no `/api/v1` prefix in protected group)
- **Accept invite response**: now includes `match` object with `{id, game_type, player1_sid, player2_sid, status}`
- **Accept invite frontend**: navigates to `/game/:gameType/:matchId` (opens live board + WS)

### Required Migrations

- `002_game_invites.sql` — creates `game_invites` table with 5-minute TTL for pending invites

### Portal Integration

- **postMessage types received**: `erlink_games_invite_accept`, `erlink_games_invite_decline`
- **postMessage types sent**: `erlink_games_handoff` (on load), `erlink_games_ready`
- **Embed flow**: Portal sends `erlink_games_handoff` → SPA stores JWT from `POST /api/v1/auth/embed`
- **Notifications events**: `game.invite`, `game.invite.accepted`, `game.invite.declined` (ADR 029 payload)

### Copy Tree for corp_mes

Source: `games_app` → Destination: `corp_mes/services/erlink-games-api/`

```
cmd/server/main.go          → server entrypoint
internal/                   → all Go packages
migrations/                 → PostgreSQL migrations (001, 002)
web/dist/                   → built SPA (React 19 + Vite 6)
Dockerfile                  → multi-stage build
go.mod, go.sum              → Go module
```

### Platform-side MR (corp_mes)

1. **`apps/erlink-portal`**: Handle inbox click → `postMessage('erlink_games_invite_accept', {invite_id})` to games iframe
2. **`services/erlink-notifications-api`**: Accept `game.invite*` events with `source_subsystem: "games"`
3. **Smoke test**: `https://localhost:3443/portal/?section=games`

### Coverage

- `internal/handler`: 81.1%
- `internal/websocket`: 81.2%
- `internal/service`: 70.0% (s3_client + portal_api external deps excluded from target)
- Total: 15 packages, all pass

### E2E Tests

Run with: `go test -tags=e2e -v ./tests/e2e/...`

- `TestChessParty_E2E`: invite → accept → match verification
- `TestChessParty_DeclineInvite_E2E`: invite → decline verification

### Backlog (Not in This Release)

- G3: Tournaments + Admin panel (wave 2)
- Redis room persistence for WS (pilot-safe, Issue 5)
- Backlog games: snake, mines, poker, trivia (not shown in prod UI)
- LiveKit voice/video (deferred — not needed for board games)
- Avatar upload to S3 (UI exists, needs external S3 endpoint)
