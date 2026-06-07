// ─── Embed Handoff Module (G0) ─────────────────────────────────────────────
// Handles the postMessage handshake between the Portal Shell (parent) and
// this Games app loaded in an iframe.
//
// Flow:
//  1. Portal loads /games/?embedInPortal=1
//  2. On mount we post "erlink_games_embed_ready" to parent
//  3. Parent posts back "erlink_games_embed_session" with credentials + secret
//  4. We exchange those for a JWT via POST /api/v1/auth/embed
//  5. JWT is stored in localStorage under "gz_token"

import { embedAuth } from './api/auth'

// ─── Types ────────────────────────────────────────────────────────────────

export type InviteEventHandler = (inviteId: string) => void

export interface GamesEmbedSessionPayload {
  type: 'erlink_games_embed_session'
  sid: string
  email: string
  name: string
  department: string
  embed_handoff_secret: string
  theme?: 'light' | 'dark'
}

interface EmbedAuthResult {
  token: string
  sid: string
  valid: boolean
}

// ─── Embed mode detection ────────────────────────────────────────────────

const EMBED_PARAM = 'embedInPortal'

export function isEmbedMode(): boolean {
  if (typeof window === 'undefined') return false
  const params = new URLSearchParams(window.location.search)
  return params.get(EMBED_PARAM) === '1'
}

// ─── Theme helper ──────────────────────────────────────────────────────────

let _handoffTheme: 'light' | 'dark' | null = null

export function getHandoffTheme(): 'light' | 'dark' | null {
  return _handoffTheme
}

// ─── PostMessage handshake ────────────────────────────────────────────────

let _handoffResolve:
  | ((value: GamesEmbedSessionPayload) => void)
  | null = null

function onPostMessage(event: MessageEvent) {
  const data = event.data
  if (
    data &&
    data.type === 'erlink_games_embed_session' &&
    data.sid &&
    data.embed_handoff_secret
  ) {
    _handoffTheme = data.theme ?? null
    if (_handoffResolve) {
      _handoffResolve(data as GamesEmbedSessionPayload)
      _handoffResolve = null
    }
  }
}

/**
 * Signal readiness to the Portal Shell and wait for session credentials.
 * Resolves with the session payload from the parent window.
 */
export function initEmbedHandoff(): Promise<GamesEmbedSessionPayload> {
  return new Promise<GamesEmbedSessionPayload>((resolve) => {
    _handoffResolve = resolve

    window.addEventListener('message', onPostMessage)

    // Tell parent we're ready
    window.parent.postMessage({ type: 'erlink_games_embed_ready' }, '*')

    // Timeout after 15s — Portal Shell should respond quickly
    setTimeout(() => {
      if (_handoffResolve) {
        window.removeEventListener('message', onPostMessage)
        _handoffResolve = null
        console.warn('Embed handoff timed out — no session received from Portal Shell')
      }
    }, 15_000)
  })
}

// ─── Token exchange ────────────────────────────────────────────────────────

/**
 * Exchange embed session credentials for a JWT token via the backend.
 */
export async function exchangeEmbedToken(
  session: GamesEmbedSessionPayload,
): Promise<EmbedAuthResult> {
  const result = await embedAuth(
    {
      sid: session.sid,
      email: session.email,
      name: session.name,
      department: session.department,
    },
    session.embed_handoff_secret,
  )
  return result
}

// ─── PostMessage invite event listeners ────────────────────────────────────

let _onInviteAccept: InviteEventHandler | null = null
let _onInviteDecline: InviteEventHandler | null = null

function onInvitePostMessage(event: MessageEvent) {
  const data = event.data
  if (!data || typeof data !== 'object') return

  if (data.type === 'erlink_games_invite_accept' && data.invite_id) {
    _onInviteAccept?.(String(data.invite_id))
  }
  if (data.type === 'erlink_games_invite_decline' && data.invite_id) {
    _onInviteDecline?.(String(data.invite_id))
  }
}

/**
 * Register callbacks for invite postMessage events from the parent iframe.
 * Call these during app initialization to handle parent-driven invite actions.
 */
export function onInviteAccept(handler: InviteEventHandler): void {
  _onInviteAccept = handler
  ensureInviteListener()
}

export function onInviteDecline(handler: InviteEventHandler): void {
  _onInviteDecline = handler
  ensureInviteListener()
}

let _inviteListenerAttached = false

function ensureInviteListener(): void {
  if (_inviteListenerAttached) return
  window.addEventListener('message', onInvitePostMessage)
  _inviteListenerAttached = true
}
