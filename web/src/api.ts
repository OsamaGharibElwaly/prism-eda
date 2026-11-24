import type { Health, Point, Preview, Profile } from './types'

async function unwrap<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const body = await response.json().catch(() => ({})) as { error?: string }
    throw new Error(body.error ?? `Request failed (${response.status})`)
  }
  return response.json() as Promise<T>
}

export function getHealth(): Promise<Health> {
  return fetch('/api/health').then(unwrap<Health>)
}

export function uploadDataset(file: File): Promise<Profile> {
  const body = new FormData()
  body.append('file', file)
  return fetch('/api/datasets', { method: 'POST', body }).then(unwrap<Profile>)
}

export function loadSample(): Promise<Profile> {
  return fetch('/api/datasets/sample', { method: 'POST' }).then(unwrap<Profile>)
}

export function getPreview(id: string, limit = 12): Promise<Preview> {
  return fetch(`/api/datasets/${id}/preview?limit=${limit}`).then(unwrap<Preview>)
}

export function getScatter(id: string, x: string, y: string): Promise<Point[]> {
  return fetch(`/api/datasets/${id}/scatter?x=${encodeURIComponent(x)}&y=${encodeURIComponent(y)}&limit=1500`)
    .then(unwrap<{ points: Point[] }>)
    .then((body) => body.points)
}
