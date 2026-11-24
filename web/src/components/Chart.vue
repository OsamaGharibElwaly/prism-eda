<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { echarts, type EChartsCoreOption } from '../echarts'

const { option, height = '250px' } = defineProps<{ option: EChartsCoreOption; height?: string }>()

const host = ref<HTMLDivElement | null>(null)
let chart: echarts.ECharts | undefined
let observer: ResizeObserver | undefined

onMounted(() => {
  if (!host.value) return
  chart = echarts.init(host.value, undefined, { renderer: 'canvas' })
  chart.setOption(option)
  observer = new ResizeObserver(() => chart?.resize())
  observer.observe(host.value)
})

watch(() => option, (next) => chart?.setOption(next, true), { deep: true })

onBeforeUnmount(() => {
  observer?.disconnect()
  chart?.dispose()
})
</script>

<template>
  <div ref="host" class="chart-host" :style="{ height }" />
</template>

<style scoped>
.chart-host { width: 100%; }
</style>
