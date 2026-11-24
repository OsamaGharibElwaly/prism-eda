<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Preview, Profile } from '../../types'
import { getPreview } from '../../api'

const { profile } = defineProps<{ profile: Profile }>()
const preview = ref<Preview | null>(null)
const error = ref('')

watch(() => profile.id, async () => {
  try {
    preview.value = await getPreview(profile.id, 12)
    error.value = ''
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Could not load preview'
  }
}, { immediate: true })
</script>

<template>
  <div class="preview">
    <p v-if="error" class="error">{{ error }}</p>
    <template v-else-if="preview">
      <div class="preview-meta">First {{ preview.rows.length }} of {{ preview.total }} rows</div>
      <div class="table-wrap">
        <table>
          <thead><tr><th>#</th><th v-for="header in preview.headers" :key="header">{{ header }}</th></tr></thead>
          <tbody>
            <tr v-for="(row, index) in preview.rows" :key="index">
              <td class="row-num">{{ index + 1 }}</td>
              <td v-for="(cell, i) in row" :key="i" :class="{ empty: cell === '' }">{{ cell || '—' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
    <div v-else class="loading">Loading rows…</div>
  </div>
</template>

<style scoped>
.preview-meta { margin-bottom: 8px; font-size: 10px; color: var(--muted); }
.table-wrap { max-height: 292px; overflow: auto; border: 1px solid var(--border); border-radius: 10px; }
table { width: 100%; border-collapse: collapse; font-size: 11px; font-variant-numeric: tabular-nums; }
th {
  position: sticky; top: 0; z-index: 1; padding: 8px 10px; text-align: left; white-space: nowrap;
  font-weight: 600; color: #b9c7dd; background: #151d2e; border-bottom: 1px solid var(--border);
}
td { padding: 7px 10px; white-space: nowrap; color: #c7d3e6; border-bottom: 1px solid rgba(148, 163, 184, .07); }
tr:last-child td { border-bottom: 0; }
tbody tr:hover td { background: rgba(34, 211, 238, .045); }
.row-num { color: #5f708c; font-size: 10px; }
.empty { color: #61728e; font-style: italic; }
.loading, .error { padding: 40px 0; text-align: center; font-size: 12px; color: var(--muted); }
.error { color: var(--danger); }
</style>
