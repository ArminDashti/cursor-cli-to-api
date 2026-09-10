import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const AUTH_KEY = 'cursor-cli-to-api-auth'

export function getApiBase(): string {
  return (import.meta.env.VITE_API_BASE_URL as string) || ''
}

export function saveAuth(username: string, password: string) {
  localStorage.setItem(AUTH_KEY, JSON.stringify({ username, password }))
}

export function clearAuth() {
  localStorage.removeItem(AUTH_KEY)
}

export function loadAuth(): { username: string; password: string } | null {
  const raw = localStorage.getItem(AUTH_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

export function authHeader(): HeadersInit {
  const auth = loadAuth()
  if (!auth) return {}
  const token = btoa(`${auth.username}:${auth.password}`)
  return { Authorization: `Basic ${token}` }
}

export async function apiFetch(path: string, init: RequestInit = {}) {
  const base = getApiBase().replace(/\/$/, '')
  const headers = new Headers(init.headers || {})
  headers.set('Content-Type', 'application/json')
  const auth = authHeader()
  Object.entries(auth).forEach(([k, v]) => headers.set(k, v as string))
  const res = await fetch(`${base}${path}`, { ...init, headers })
  const text = await res.text()
  let json: unknown = null
  try {
    json = text ? JSON.parse(text) : null
  } catch {
    json = text
  }
  return { res, json, text }
}
