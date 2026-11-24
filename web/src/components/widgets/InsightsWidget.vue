<script setup lang="ts">
import type { Profile } from '../../types'
import Icon from '../Icon.vue'

const { profile } = defineProps<{ profile: Profile }>()

const icons: Record<string, string> = {
  warning: 'alert', opportunity: 'sparkles', info: 'info',
}
</script>

<template>
  <ul class="insights">
    <li v-for="(insight, index) in profile.insights" :key="index" :class="insight.severity">
      <span class="insight-icon"><Icon :name="icons[insight.severity] ?? 'info'" :size="15" /></span>
      <div>
        <strong>{{ insight.title }}</strong>
        <p>{{ insight.detail }}</p>
      </div>
    </li>
    <li v-if="!profile.insights.length" class="none">Nothing unusual detected — this dataset looks clean.</li>
  </ul>
</template>

<style scoped>
.insights { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 8px; max-height: 330px; overflow: auto; }
.insights li { display: flex; gap: 10px; padding: 10px 11px; border-radius: 11px; background: rgba(255, 255, 255, .025); border: 1px solid var(--border); }
.insight-icon { flex: none; display: flex; padding-top: 2px; }
.warning .insight-icon { color: var(--warn); }
.opportunity .insight-icon { color: var(--accent); }
.info .insight-icon { color: var(--accent2); }
strong { display: block; font-size: 11px; font-weight: 600; color: #dfe8f8; }
p { margin: 4px 0 0; font-size: 10px; line-height: 1.5; color: var(--muted); }
.none { display: block !important; color: var(--muted); font-size: 11px; }
</style>
