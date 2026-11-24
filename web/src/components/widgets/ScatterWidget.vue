<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { EChartsCoreOption } from '../../echarts'
import type { Point, Profile } from '../../types'
import { getScatter } from '../../api'
import { axisLabel, line, splitLine, tooltip } from '../../charts'
import Chart from '../Chart.vue'

const { profile } = defineProps<{ profile: Profile }>()

const numeric = computed(() => profile.columns.filter((c) => c.kind === 'numeric'))
const x = ref('')
const y = ref('')
const points = ref<Point[]>([])
const loading = ref(false)

watch(numeric, (list) => {
  if (!list.some((c) => c.name === x.value)) x.value = list[0]?.name ?? ''
  if (!list.some((c) => c.name === y.value)) y.value = list[1]?.name ?? list[0]?.name ?? ''
}, { immediate: true })

async function load() {
  if (!x.value || !y.value || x.value === y.value) {
    points.value = []
    return
  }
  loading.value = true
  try {
    points.value = await getScatter(profile.id, x.value, y.value)
  } catch {
    points.value = []
  } finally {
    loading.value = false
  }
}
watch([() => profile.id, x, y], load, { immediate: true })

const relation = computed(() => profile.correlations.find(
  (c) => (c.x === x.value && c.y === y.value) || (c.x === y.value && c.y === x.value),
))

const regression = computed(() => {
  const data = points.value
  if (data.length < 2) return null
  const mx = data.reduce((s, p) => s + p.x, 0) / data.length
  const my = data.reduce((s, p) => s + p.y, 0) / data.length
  let cov = 0
  let variance = 0
  for (const p of data) {
    cov += (p.x - mx) * (p.y - my)
    variance += (p.x - mx) ** 2
  }
  if (variance === 0) return null
  const slope = cov / variance
  const intercept = my - slope * mx
  const minX = Math.min(...data.map((p) => p.x))
  const maxX = Math.max(...data.map((p) => p.x))
  return [[minX, slope * minX + intercept], [maxX, slope * maxX + intercept]]
})

const option = computed<EChartsCoreOption>(() => ({
  animationDuration: 400,
  grid: { left: 62, right: 20, top: 16, bottom: 42 },
  tooltip: {
    ...tooltip,
    trigger: 'item',
    formatter: (params: unknown) => {
      const p = (params as { value: number[] }).value
      return `<b>${x.value}</b> ${p[0]}<br/><b>${y.value}</b> ${p[1]}`
    },
  },
  xAxis: { type: 'value', name: x.value, nameLocation: 'middle', nameGap: 26, nameTextStyle: { color: '#8798b5', fontSize: 10 }, axisLabel, splitLine, axisLine: { lineStyle: { color: line } } },
  yAxis: { type: 'value', name: y.value, nameLocation: 'middle', nameGap: 42, nameTextStyle: { color: '#8798b5', fontSize: 10 }, axisLabel, splitLine, axisLine: { lineStyle: { color: line } } },
  series: [
    { type: 'scatter', data: points.value.map((p) => [p.x, p.y]), symbolSize: 7, itemStyle: { color: '#22d3ee', opacity: .55 } },
    ...(regression.value ? [{
      type: 'line' as const,
      data: regression.value,
      symbol: 'none',
      silent: true,
      lineStyle: { color: '#a78bfa', width: 2, type: 'dashed' as const },
    }] : []),
  ],
} as EChartsCoreOption))
</script>

<template>
  <div class="scatter">
    <div class="toolbar">
      <label for="scatter-x">X</label>
      <select id="scatter-x" v-model="x"><option v-for="c in numeric" :key="c.name" :value="c.name">{{ c.name }}</option></select>
      <label for="scatter-y">Y</label>
      <select id="scatter-y" v-model="y"><option v-for="c in numeric" :key="c.name" :value="c.name">{{ c.name }}</option></select>
    </div>
    <div v-if="loading" class="status">Loading points…</div>
    <template v-else-if="points.length && x !== y">
      <div class="meta">
        <span>{{ points.length }} points</span>
        <span v-if="relation">Pearson r <b>{{ relation.coefficient.toFixed(3) }}</b></span>
        <span class="trend"><i /> regression</span>
      </div>
      <Chart :option="option" height="230px" />
    </template>
    <div v-else class="status">Choose two different numeric columns.</div>
  </div>
</template>

<style scoped>
.toolbar { display: flex; align-items: center; gap: 7px; margin-bottom: 10px; }
label { font-size: 11px; color: var(--muted); }
select {
  flex: 1; min-width: 0; padding: 7px 9px; font-size: 12px; color: var(--text);
  background: rgba(8, 13, 24, .8); border: 1px solid var(--border); border-radius: 9px; outline: none;
}
select:focus { border-color: var(--accent); }
.meta { display: flex; gap: 12px; margin-bottom: 6px; font-size: 10px; color: var(--muted); }
.meta b { color: var(--accent); }
.trend { display: flex; align-items: center; gap: 5px; margin-left: auto; }
.trend i { width: 14px; border-top: 2px dashed var(--accent2); }
.status { display: grid; place-items: center; min-height: 230px; color: var(--muted); font-size: 12px; }
</style>
