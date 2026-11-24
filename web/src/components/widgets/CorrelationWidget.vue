<script setup lang="ts">
import { computed } from 'vue'
import type { EChartsCoreOption } from '../../echarts'
import type { Profile } from '../../types'
import { axisLabel } from '../../charts'
import Chart from '../Chart.vue'

const { profile } = defineProps<{ profile: Profile }>()

const labels = computed(() => profile.matrix.labels)
const hasMatrix = computed(() => labels.value.length > 1)

const option = computed<EChartsCoreOption>(() => {
  const names = labels.value
  const values = profile.matrix.values
  const data: { value: [number, number, number] }[] = []
  for (let row = 0; row < names.length; row++) {
    for (let col = 0; col < names.length; col++) {
      data.push({ value: [col, row, values[row]?.[col] ?? 0] })
    }
  }
  return {
    animationDuration: 450,
    grid: { left: 94, right: 20, top: 12, bottom: 78 },
    tooltip: {
      position: 'top',
      formatter: (params: unknown) => {
        const p = (params as { value: [number, number, number] }).value
        return `<b>${names[p[1]]} ↔ ${names[p[0]]}</b><br/>r = ${p[2].toFixed(3)}`
      },
    },
    xAxis: {
      type: 'category', data: names, splitArea: { show: false },
      axisLabel: { ...axisLabel, rotate: 38 },
      axisLine: { show: false }, axisTick: { show: false },
    },
    yAxis: {
      type: 'category', data: names, splitArea: { show: false },
      axisLabel, axisLine: { show: false }, axisTick: { show: false },
    },
    visualMap: {
      min: -1, max: 1, calculable: false, orient: 'horizontal', left: 'center', bottom: 2,
      itemWidth: 11, itemHeight: 130, text: ['+1', '−1'], textStyle: { color: '#8798b5', fontSize: 10 },
      inRange: { color: ['#f472b6', '#243049', '#22d3ee'] },
    },
    series: [{
      type: 'heatmap', data,
      label: { show: names.length <= 8, color: '#e8eefb', fontSize: 10, formatter: (p: unknown) => (p as { value: number[] }).value[2].toFixed(2) },
      itemStyle: { borderColor: '#0c1220', borderWidth: 2, borderRadius: 4 },
      emphasis: { itemStyle: { shadowBlur: 12, shadowColor: 'rgba(34,211,238,.6)' } },
    }],
  } as EChartsCoreOption
})
</script>

<template>
  <div v-if="hasMatrix">
    <p class="hint">Pearson correlation across {{ labels.length }} numeric features</p>
    <Chart :option="option" height="270px" />
  </div>
  <div v-else class="empty-mini">Need at least two numeric columns for a correlation matrix.</div>
</template>

<style scoped>
.hint { margin: 0 0 8px; font-size: 11px; color: var(--muted); }
.empty-mini { padding: 40px 0; text-align: center; color: var(--muted); font-size: 12px; }
</style>
