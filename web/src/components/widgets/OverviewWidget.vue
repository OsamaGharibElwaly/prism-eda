<script setup lang="ts">
import { computed } from 'vue'
import type { Profile } from '../../types'
import { num, timeAgo } from '../../format'
import Icon from '../Icon.vue'

const { profile } = defineProps<{ profile: Profile }>()

const missingCells = computed(() => profile.columns.reduce((sum, c) => sum + c.missing, 0))
const stats = computed(() => [
  { label: 'Rows', value: num(profile.rows, 0), icon: 'database', tone: 'cyan' },
  { label: 'Columns', value: num(profile.column_count, 0), icon: 'grid', tone: 'violet' },
  { label: 'Numeric features', value: num(profile.numeric_columns.length, 0), icon: 'chart', tone: 'lime' },
  { label: 'Missing cells', value: num(missingCells.value, 0), icon: 'alert', tone: 'amber' },
])
</script>

<template>
  <div class="overview">
    <div class="dataset-line">
      <span class="file-badge"><Icon name="database" :size="15" /></span>
      <span class="file-name" :title="profile.name">{{ profile.name }}</span>
      <span class="file-time"><Icon name="clock" :size="12" /> {{ timeAgo(profile.created_at) }}</span>
    </div>
    <div class="stats">
      <div v-for="stat in stats" :key="stat.label" class="stat" :class="stat.tone">
        <span class="stat-icon"><Icon :name="stat.icon" :size="15" /></span>
        <span class="stat-label">{{ stat.label }}</span>
        <strong>{{ stat.value }}</strong>
      </div>
    </div>
    <div class="foot">
      <span><Icon name="bolt" :size="13" /> Profiled in {{ profile.analysis_ms }} ms</span>
      <span>ID {{ profile.id.slice(0, 8) }}</span>
    </div>
  </div>
</template>

<style scoped>
.dataset-line { display: flex; align-items: center; gap: 9px; margin-bottom: 14px; font-size: 12px; }
.file-badge { color: var(--accent); display: flex; }
.file-name { font-weight: 600; color: var(--text); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.file-time { display: flex; align-items: center; gap: 4px; margin-left: auto; color: #6c7d98; font-size: 11px; white-space: nowrap; }
.stats { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 10px; }
.stat { position: relative; padding: 12px 13px; border-radius: 13px; background: rgba(255, 255, 255, .026); border: 1px solid var(--border); }
.stat-icon { position: absolute; top: 11px; right: 11px; opacity: .75; }
.stat-label { display: block; font-size: 10px; letter-spacing: .08em; text-transform: uppercase; color: var(--muted); }
.stat strong { display: block; margin-top: 5px; font-size: 23px; font-weight: 650; letter-spacing: -.03em; font-variant-numeric: tabular-nums; }
.cyan { color: var(--accent); } .violet { color: var(--accent2); } .lime { color: #a3e635; } .amber { color: var(--warn); }
.stat strong { color: var(--text); }
.foot { display: flex; justify-content: space-between; gap: 8px; margin-top: 13px; font-size: 10px; color: #667793; }
.foot span { display: flex; align-items: center; gap: 5px; }
</style>
