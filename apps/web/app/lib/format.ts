export function formatPrice(value?: number | null): string {
  if (value === null || value === undefined || Number.isNaN(value)) return '—'
  return `¥${value.toFixed(0)}`
}

export function formatDateLabel(date: string): string {
  const parsed = new Date(date)
  if (Number.isNaN(parsed.getTime())) return date
  return parsed.toLocaleDateString('ja-JP', { month: '2-digit', day: '2-digit' })
}

export function formatDistance(meters: number | null): string {
  if (meters === null) return '—'
  const km = meters / 1000
  return `${km.toFixed(1)} km`
}
