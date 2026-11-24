<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { EChartsCoreOption } from '../../echarts'
import type { Profile } from '../../types'
import { axisNum, num } from '../../format'
import { axisLabel, line, splitLine, tooltip } from '../../charts'
import Chart from '../Chart.vue'

const { profile } = defineProps<{ profile: Profile }>()

const numeric = computed(() => profile.columns.filter((c) => c.kind === 'numeric' && c.numeric))
const selected = ref('')
watch(numeric, (list) => {
  if (!list.some((c) => c.name === selected.value)) selected.value = list[0]?.name ?? ''
}, { immediate: true })

const column = computed(() => numeric.value.find((c) => c.name === selected.value) ?? numeric.value[0])

const option = computed<EChartsCoreOption>(() => {
  const bins = column.value?.numeric?.histogram ?? []
  return {
    animationDuration: 450,
    grid: { left: 46, right: 14, top: 14, bottom: 36 },
    tooltip: {
      ...tooltip,
      formatter: (params: unknown) => {
        const point = (params as { name: string; value: number }[])[0]
        return `<b>${point.name}</b><br/>count: ${num(point.value, 0)}`
      },
    },
    xAxis: {
      type: 'category',
      data: bins.map((b) => axisNum(b.start)),
      axisLabel,
      axisTick: { show: false },
      axisLine: { lineStyle: { color: line } },
    },
    yAxis: { type: 'value', axisLabel, splitLine, axisLine: { show: false } },
    series: [{
      type: 'bar',
      data: bins.map((b) => ({
        value: b.count,
        name: `${axisNum(b.start)} – ${axisNum(b.end)}`,
      })),
      barMaxWidth: 46,
      itemStyle: { color: '#22d3ee', borderRadius: [5, 5, 0, 0], opacity: .88 },
      emphasis: { itemStyle: { color: '#67e8f9' } },
    }],
  } as EChartsCoreOption
})
</script>

<template>
  <div v-if="column?.numeric" class="distribution">
    <div class="toolbar">
      <label for="dist-column">Column</label>
      <select id="dist-column" v-model="selected">
        <option v-for="c in numeric" :key="c.name" :value="c.name">{{ c.name }}</option>
      </select>
      <span class="distinct">{{ num(column.distinct, 0) }} unique</span>
    </div>
    <div class="chips">
      <span>mean <b>{{ num(column.numeric.mean) }}</b></span>
      <span>median <b>{{ num(column.numeric.median) }}</b></span>
      <span>σ <b>{{ num(column.numeric.std_dev) }}</b></span>
      <span>outliers <b>{{ column.numeric.outliers }}</b></span>
    </div>
    <Chart :option="option" height="220px" />
  </div>
  <div v-else class="empty-mini">No numeric columns to plot.</div>
</template>

<style scoped>
.toolbar { display: flex; align-items: center; gap: 9px; margin-bottom: 12px; }
label { font-size: 11px; color: var(--muted); }
select {
  flex: 1; min-width: 0; padding: 7px 10px; color: var(--text); font-size: 12px;
  background: rgba(8, 13, 24, .8); border: 1px solid var(--border); border-radius: 9px; outline: none;
}
select:focus { border-color: var(--accent); }
.distinct { font-size: 10px; color: #71829d; white-space: nowrap; }
.chips { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 8px; }
.chips span { padding: 4px 8px; font-size: 10px; color: var(--muted); border-radius: 7px; background: rgba(255, 255, 255, .035); }
.chips b { color: #dbe6f7; font-weight: 600; font-variant-numeric: tabular-nums; }
.empty-mini { padding: 40px 0; text-align: center; color: var(--muted); font-size: 12px; }
</style>
