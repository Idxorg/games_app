import { post, get, type User } from './client'

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
