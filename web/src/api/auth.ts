import { post, get, setToken, type User } from './client'

interface VerifyResponse {
  valid: boolean
  sid?: string
  email?: string
  name?: string
  groups?: string[]
  portal_verified?: boolean
  message?: string
}

export async function verifyToken(token: string): Promise<VerifyResponse> {
  // For verify we pass the token directly in Authorization header
  // The client already attaches it from localStorage, but verify needs explicit token
  const res = await fetch(`${import.meta.env.VITE_API_URL || '/api/v1'}/auth/verify`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  })
  if (!res.ok) {
    const body = await res.json().catch(() => null)
    throw { status: res.status, body }
  }
  return res.json()
}

export async function getProfile(sid: string): Promise<User> {
  return get<User>(`/users/${sid}/profile`)
}

// ─── Embed Auth (G0) ───────────────────────────────────────────────────────

interface EmbedAuthRequest {
  sid: string
  email?: string
  name?: string
  department?: string
}

interface EmbedAuthResponse {
  token: string
  sid: string
  valid: boolean
}

export async function embedAuth(
  request: EmbedAuthRequest,
  embedHandoffSecret: string,
): Promise<EmbedAuthResponse> {
  const baseUrl = import.meta.env.VITE_API_URL || '/api/v1'
  const res = await fetch(`${baseUrl}/auth/embed`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Erlink-Embed-Secret': embedHandoffSecret,
    },
    body: JSON.stringify(request),
  })
  if (!res.ok) {
    const body = await res.json().catch(() => null)
    throw { status: res.status, body }
  }
  const data = await res.json()
  // Store JWT token
  setToken(data.token)
  return data
}
