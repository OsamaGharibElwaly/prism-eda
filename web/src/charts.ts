export const palette = ['#22d3ee', '#a78bfa', '#f472b6', '#a3e635', '#fbbf24', '#60a5fa']

export const muted = '#8798b5'
export const line = 'rgba(148,163,184,.16)'

export const tooltip = {
  trigger: 'axis',
  backgroundColor: '#0d1425',
  borderColor: 'rgba(148,163,184,.25)',
  borderWidth: 1,
  padding: [10, 12],
  textStyle: { color: '#e8eefb', fontSize: 12 },
  extraCssText: 'border-radius:10px;box-shadow:0 16px 40px rgba(0,0,0,.5)',
}

export const axisLabel = { color: muted, fontSize: 11, hideOverlap: true }

export const splitLine = { show: true, lineStyle: { color: line, type: 'dashed' as const } }
