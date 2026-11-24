const compact = new Intl.NumberFormat('en-US', { notation: 'compact', maximumFractionDigits: 2 })

export function num(value: number, digits = 2): string {
  if (!Number.isFinite(value)) return '—'
  if (Math.abs(value) >= 1_000_000) return compact.format(value)
  return value.toLocaleString('en-US', { maximumFractionDigits: digits, minimumFractionDigits: 0 })
}

export function axisNum(value: number): string {
  if (!Number.isFinite(value)) return ''
  if (Math.abs(value) >= 10_000) return compact.format(value)
  return String(Math.round(value * 100) / 100)
}

export function timeAgo(iso: string): string {
  const seconds = Math.max(0, Math.round((Date.now() - new Date(iso).getTime()) / 1000))
  if (seconds < 60) return 'just now'
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
  return `${Math.floor(seconds / 3600)}h ago`
}
