export type ColumnKind = 'numeric' | 'categorical' | 'boolean'

export interface Bin {
  start: number
  end: number
  count: number
}

export interface NumericStats {
  min: number
  max: number
  mean: number
  median: number
  std_dev: number
  sum: number
  p25: number
  p75: number
  outliers: number
  histogram: Bin[]
}

export interface ValueCount {
  value: string
  count: number
  percent: number
}

export interface CategoricalStats {
  top_values: ValueCount[]
  entropy: number
}

export interface ColumnProfile {
  name: string
  index: number
  kind: ColumnKind
  count: number
  missing: number
  missing_pct: number
  distinct: number
  numeric?: NumericStats
  categorical?: CategoricalStats
}

export interface Correlation {
  x: string
  y: string
  coefficient: number
  strength: string
  direction: string
}

export interface CorrelationMatrix {
  labels: string[]
  values: number[][]
}

export interface Insight {
  kind: string
  severity: 'warning' | 'opportunity' | 'info'
  title: string
  detail: string
  column?: string
}

export interface Quality {
  score: number
  completeness: number
  duplicate_rows: number
  invalid_rows: number
  issue_count: number
}

export interface Profile {
  id: string
  name: string
  rows: number
  column_count: number
  created_at: string
  analysis_ms: number
  numeric_columns: string[]
  columns: ColumnProfile[]
  correlations: Correlation[]
  matrix: CorrelationMatrix
  insights: Insight[]
  quality: Quality
}

export interface Preview {
  headers: string[]
  rows: string[][]
  total: number
}

export interface Point {
  x: number
  y: number
}

export interface Health {
  status: string
  version: string
  go: string
  uptime_s: number
  datasets: number
}

export type WidgetType = 'overview' | 'distribution' | 'correlation' | 'missing' | 'scatter' | 'preview' | 'insights'

export interface WidgetInstance {
  id: string
  type: WidgetType
}
