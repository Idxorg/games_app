import { get, post } from './client'

// ─── Types ────────────────────────────────────────────────────────────────

export interface Invite {
  id: string
  game_type: string
  inviter_sid: string
  inviter_name: string
  recipient_sid: string
  status: 'pending' | 'accepted' | 'declined'
  created_at: string
  updated_at: string
}

interface PendingInvitesResponse {
  invites: Invite[]
  count: number
}

interface CreateInviteResponse {
  invite: Invite
  message: string
}

interface InviteActionResponse {
  message: string
  invite_id: string
}

// ─── API functions ──────────────────────────────────────────────────────────

export async function createInvite(
  gameType: string,
  recipientSid: string,
): Promise<Invite> {
  const res = await post<CreateInviteResponse>('/games/invite', {
    game_type: gameType,
    recipient_sid: recipientSid,
  })
  return res.invite
}

export async function acceptInvite(id: string): Promise<InviteActionResponse> {
  return post<InviteActionResponse>(`/games/invite/${id}/accept`)
}

export async function declineInvite(id: string): Promise<InviteActionResponse> {
  return post<InviteActionResponse>(`/games/invite/${id}/decline`)
}

export async function getPendingInvites(): Promise<Invite[]> {
  const res = await get<PendingInvitesResponse>('/games/invite/pending')
  return res.invites
}
