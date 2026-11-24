<script setup lang="ts">
const { name = 'grid', size = 18 } = defineProps<{ name: string; size?: number }>()

const paths: Record<string, string> = {
  upload: 'M12 16V4m0 0L7 9m5-5 5 5M4 16v3a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-3',
  chart: 'M4 19V9m5 10V5m5 14v-7m5 7V8',
  grid: 'M4 4h7v7H4zM13 4h7v7h-7zM4 13h7v7H4zM13 13h7v7h-7z',
  table: 'M3 5h18v14H3zM3 10h18M9 10v9M15 10v9',
  target: 'M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zm0 5a4 4 0 1 0 0 8 4 4 0 0 0 0-8zm0 3a1 1 0 1 0 0 2 1 1 0 0 0 0-2z',
  alert: 'M12 3 2 20h20L12 3zm0 6v5m0 3h.01',
  trash: 'M4 7h16M9 7V4h6v3m-9 0 1 13h10l1-13M10 11v6m4-6v6',
  grip: 'M9 6h.01M9 12h.01M9 18h.01M15 6h.01M15 12h.01M15 18h.01',
  check: 'M5 13l4 4L19 7',
  bolt: 'M13 2 4 14h7l-1 8 9-12h-7l1-8z',
  database: 'M4 6c0-2 4-3 8-3s8 1 8 3-4 3-8 3-8-1-8-3zm0 0v12c0 2 4 3 8 3s8-1 8-3V6M4 12c0 2 4 3 8 3s8-1 8-3',
  refresh: 'M20 7v5h-5M4 17v-5h5M6 9a7 7 0 0 1 12-2l2 2M18 15a7 7 0 0 1-12 2l-2-2',
  sparkles: 'M12 2l2.2 6.8L21 11l-6.8 2.2L12 20l-2.2-6.8L3 11l6.8-2.2L12 2z',
  close: 'M6 6l12 12M18 6 6 18',
  layers: 'M12 3 3 8l9 5 9-5-9-5zM3 13l9 5 9-5',
  info: 'M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zm0 8v5m0-8h.01',
  arrows: 'M7 7h10l-3-3m3 3-3 3M17 17H7l3 3m-3-3 3-3',
  clock: 'M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zm0 4v5l4 2',
}
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"
    stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
    <path :d="paths[name] ?? paths.grid" />
  </svg>
</template>
