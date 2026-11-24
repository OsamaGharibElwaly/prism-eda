<script setup lang="ts">
import { computed } from 'vue'
import type { EChartsCoreOption } from '../../echarts'
import type { Profile } from '../../types'
import { num } from '../../format'
import { axisLabel, splitLine, tooltip } from '../../charts'
import Chart from '../Chart.vue'
import Icon from '../Icon.vue'

const { profile } = defineProps<{ profile: Profile }>()

const rows = computed(() => profile.columns.filter((c) => c.missing > 0).sort((a, b) => b.missing_pct - a.missing_pct))
const totalMissing = computed(() => profile.columns.reduce((sum, c) => sum + c.missing, 0))

const option = computed<EChartsCoreOption>(() => ({
  animationDuration: 450,
  grid: { left: 104, right: 46, top: 8, bottom: 30 },
  tooltip: {
    ...tooltip,
    formatter: (params: unknown) => {
      const p = (params as { name: string; value: number }[])[0]
      const column = rows.value.find((c) => c.name === p.name)
      return `<b>${p.name}</b><br/>${num(column?.missing ?? 0, 0)} missing · ${p.value}%`
    },
  },
  xAxis: { type: 'value', max: 100, axisLabel: { ...axisLabel, formatter: '{value}%' }, splitLine },
  yAxis: { type: 'category', data: rows.value.map((c) => c.name), axisLabel, axisLine: { show: false }, axisTick: { show: false } },
  series: [{
    type: 'bar',
    data: rows.value.map((c) => c.missing_pct),
    barMaxWidth: 18,
    itemStyle: { color: '#fbbf24', borderRadius: [0, 5, 5, 0], opacity: .9 },
    label: { show: true, position: 'right', color: '#e8eefb', fontSize: 10, formatter: '{c}%' },
  }],
} as EChartsCoreOption))
</script>

<template>
  <div v-if="rows.length" class="missing">
    <p class="hint">{{ num(totalMissing, 0) }} blank cells across {{ rows.length }} column{{ rows.length === 1 ? '' : 's' }}</p>
    <Chart :option="option" :height="`${Math.max(140, rows.length * 34 + 44)}px`" />
  </div>
  <div v-else class="complete">
    <span><Icon name="check" :size="22" /></span>
    <strong>No missing values</strong>
    <small>Every cell in this dataset is populated.</small>
  </div>
</template>

<style scoped>
.hint { margin: 0 0 6px; font-size: 11px; color: var(--muted); }
.complete { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px; min-height: 190px; text-align: center; }
.complete span { display: grid; place-items: center; width: 48px; height: 48px; border-radius: 14px; color: var(--good); background: rgba(52, 211, 153, .12); border: 1px solid rgba(52, 211, 153, .25); }
.complete strong { font-size: 13px; color: var(--text); }
.complete small { font-size: 11px; color: var(--muted); }
</style>
