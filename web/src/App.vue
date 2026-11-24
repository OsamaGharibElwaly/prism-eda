<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { Component } from 'vue'
import { getHealth, loadSample, uploadDataset } from './api'
import { num } from './format'
import type { Health, Profile, WidgetInstance, WidgetType } from './types'
import Icon from './components/Icon.vue'
import DropZone from './components/DropZone.vue'
import WidgetFrame from './components/WidgetFrame.vue'
import OverviewWidget from './components/widgets/OverviewWidget.vue'
import DistributionWidget from './components/widgets/DistributionWidget.vue'
import CorrelationWidget from './components/widgets/CorrelationWidget.vue'
import MissingWidget from './components/widgets/MissingWidget.vue'
import ScatterWidget from './components/widgets/ScatterWidget.vue'
import PreviewWidget from './components/widgets/PreviewWidget.vue'
import InsightsWidget from './components/widgets/InsightsWidget.vue'

const palette: { type: WidgetType; title: string; detail: string; icon: string }[] = [
  { type: 'overview', title: 'Overview', detail: 'Dataset at a glance', icon: 'target' },
  { type: 'distribution', title: 'Distribution', detail: 'Histogram + moments', icon: 'chart' },
  { type: 'correlation', title: 'Correlation', detail: 'Pearson heatmap', icon: 'grid' },
  { type: 'missing', title: 'Missing values', detail: 'Completeness by column', icon: 'alert' },
  { type: 'scatter', title: 'Scatter explorer', detail: 'Two columns + trend', icon: 'arrows' },
  { type: 'preview', title: 'Data preview', detail: 'First rows of the table', icon: 'table' },
  { type: 'insights', title: 'Insights', detail: 'Automated findings', icon: 'sparkles' },
]
const meta = Object.fromEntries(palette.map((p) => [p.type, p])) as Record<WidgetType, typeof palette[number]>

const widgets: Record<WidgetType, Component> = {
  overview: OverviewWidget,
  distribution: DistributionWidget,
  correlation: CorrelationWidget,
  missing: MissingWidget,
  scatter: ScatterWidget,
  preview: PreviewWidget,
  insights: InsightsWidget,
}

const profile = ref<Profile | null>(null)
const cards = ref<WidgetInstance[]>([])
const health = ref<Health | null>(null)
const busy = ref(false)
const error = ref('')
const toast = ref('')
const toastTone = ref<'ok' | 'error'>('ok')
const canvasOver = ref(false)
const movingIndex = ref<number | null>(null)
let toastTimer: number | undefined

const qualityTone = computed(() => {
  const score = profile.value?.quality.score ?? 0
  return score >= 85 ? 'var(--good)' : score >= 65 ? 'var(--warn)' : 'var(--danger)'
})

onMounted(async () => {
  try { health.value = await getHealth() } catch { /* backend may still be starting */ }
  if (new URLSearchParams(window.location.search).has('sample')) await analyse(loadSample)
})

function notify(message: string, tone: 'ok' | 'error' = 'ok') {
  toast.value = message
  toastTone.value = tone
  window.clearTimeout(toastTimer)
  toastTimer = window.setTimeout(() => { toast.value = '' }, 3500)
}

async function analyse(task: () => Promise<Profile>) {
  busy.value = true
  error.value = ''
  try {
    const result = await task()
    profile.value = result
    cards.value = (['overview', 'distribution', 'correlation', 'missing', 'scatter'] as WidgetType[])
      .map((type) => ({ id: crypto.randomUUID(), type }))
    notify(`${result.name} profiled in ${result.analysis_ms} ms`)
    try { health.value = await getHealth() } catch { /* non-critical */ }
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Analysis failed'
    notify(error.value, 'error')
  } finally {
    busy.value = false
  }
}

function reset() {
  profile.value = null
  cards.value = []
  error.value = ''
  notify('Workspace cleared')
}

function onPaletteDragStart(event: DragEvent, type: WidgetType) {
  if (!profile.value) {
    event.preventDefault()
    notify('Load a dataset before adding widgets', 'error')
    return
  }
  event.dataTransfer?.setData('application/x-prism-widget', type)
  if (event.dataTransfer) event.dataTransfer.effectAllowed = 'copy'
}

function onCanvasDrop(event: DragEvent) {
  canvasOver.value = false
  const type = event.dataTransfer?.getData('application/x-prism-widget') as WidgetType | undefined
  if (!type || !meta[type]) return
  cards.value.push({ id: crypto.randomUUID(), type })
  notify(`${meta[type].title} added to workspace`)
}

function onCardDragStart(event: DragEvent, index: number) {
  event.dataTransfer?.setData('application/x-prism-move', String(index))
  if (event.dataTransfer) event.dataTransfer.effectAllowed = 'move'
  movingIndex.value = index
}

function onCardDrop(event: DragEvent, index: number) {
  const raw = event.dataTransfer?.getData('application/x-prism-move')
  if (raw !== undefined && raw !== '') {
    const from = Number(raw)
    if (Number.isInteger(from) && from !== index && cards.value[from]) {
      const [moved] = cards.value.splice(from, 1)
      cards.value.splice(index, 0, moved)
    }
  }
  movingIndex.value = null
}
</script>

<template>
  <div class="app">
    <header class="topbar">
      <div class="brand">
        <span class="brand-mark">
          <svg viewBox="0 0 40 40" width="30" height="30" aria-hidden="true">
            <defs><linearGradient id="brandGradient" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#22d3ee" /><stop offset="1" stop-color="#a78bfa" /></linearGradient></defs>
            <path d="M20 4 37 34H3L20 4Z" fill="none" stroke="url(#brandGradient)" stroke-width="2.4" stroke-linejoin="round" />
            <path d="M20 4v30M3 34l17-17 17 17" fill="none" stroke="url(#brandGradient)" stroke-width="1.4" opacity=".65" />
          </svg>
        </span>
        <div>
          <h1>Prism<span> EDA</span></h1>
          <p>Smart data analysis, drag &amp; drop</p>
        </div>
      </div>
      <div class="top-actions">
        <span class="health"><i :class="{ down: !health }" />{{ health ? `Go ${health.version} · online` : 'connecting…' }}</span>
        <button v-if="profile" class="btn" type="button" @click="reset"><Icon name="refresh" :size="15" /> New</button>
        <button class="btn primary" type="button" :disabled="busy" @click="analyse(loadSample)">
          <Icon name="sparkles" :size="15" /> Load sample
        </button>
      </div>
    </header>

    <div v-if="error" class="error-banner"><Icon name="alert" :size="16" /> {{ error }}</div>

    <main class="workspace">
      <aside class="sidebar">
        <section class="panel">
          <div class="panel-head"><Icon name="database" :size="16" /><h2>Dataset</h2></div>
          <DropZone @file="(file) => analyse(() => uploadDataset(file))" />
          <div v-if="profile" class="dataset-card">
            <div class="dataset-top"><Icon name="check" :size="14" /> Ready to explore</div>
            <strong :title="profile.name">{{ profile.name }}</strong>
            <span>{{ profile.rows.toLocaleString() }} rows × {{ profile.column_count }} columns</span>
          </div>
        </section>

        <section class="panel palette-panel">
          <div class="panel-head"><Icon name="layers" :size="16" /><h2>Widgets</h2><small>DRAG</small></div>
          <p class="panel-desc">Drag a card onto the canvas to add it to your report.</p>
          <div class="palette">
            <div v-for="item in palette" :key="item.type" class="palette-card" draggable="true"
              @dragstart="onPaletteDragStart($event, item.type)">
              <span class="palette-icon"><Icon :name="item.icon" :size="16" /></span>
              <span class="palette-text"><strong>{{ item.title }}</strong><small>{{ item.detail }}</small></span>
              <span class="palette-grip"><Icon name="grip" :size="14" /></span>
            </div>
          </div>
          <div class="tip"><Icon name="info" :size="13" /> Drag a widget header to reorder cards.</div>
        </section>
      </aside>

      <section class="canvas-column">
        <div class="canvas-title">
          <div>
            <h2>Analysis workspace</h2>
            <p>{{ profile ? `${profile.name} · ${profile.rows.toLocaleString()} rows profiled` : 'Your report comes together here' }}</p>
          </div>
          <span v-if="profile" class="pill"><Icon name="grid" :size="13" /> {{ cards.length }} widgets</span>
        </div>

        <div class="canvas" :class="{ over: canvasOver }" @dragenter.prevent="canvasOver = true"
          @dragleave.self.prevent="canvasOver = false" @dragover.prevent @drop.prevent="onCanvasDrop">
          <div v-if="busy" class="busy"><span class="spinner" /><strong>Profiling dataset…</strong><small>Reading CSV · inferring types · computing statistics</small></div>

          <template v-if="profile && cards.length && !busy">
            <WidgetFrame v-for="(card, index) in cards" :key="card.id" :title="meta[card.type].title"
              :icon="meta[card.type].icon" :index="index" :moving="movingIndex === index"
              @remove="cards.splice(index, 1)" @dragstart="onCardDragStart($event, index)"
              @dragend="movingIndex = null" @drop="onCardDrop($event, index)">
              <component :is="widgets[card.type]" :profile="profile" />
            </WidgetFrame>
          </template>

          <div v-else-if="!profile && !busy" class="welcome">
            <div class="welcome-mark"><Icon name="sparkles" :size="30" /></div>
            <h2>Turn raw CSV into clear decisions</h2>
            <p>Drop a dataset on the left — or start with the bundled housing sample — then drag analysis cards into this workspace to compose your report.</p>
            <div class="steps"><span><b>01</b> Drop a CSV</span><i /><span><b>02</b> Go profiles it</span><i /><span><b>03</b> Drag widgets here</span></div>
            <button class="btn primary" type="button" @click="analyse(loadSample)"><Icon name="sparkles" :size="15" /> Try the sample dataset</button>
          </div>
        </div>
      </section>

      <aside class="insight-column">
        <section v-if="profile" class="panel quality-panel">
          <div class="panel-head"><Icon name="target" :size="16" /><h2>Data quality</h2></div>
          <div class="quality">
            <div class="ring-wrap">
              <svg class="ring" viewBox="0 0 120 120" aria-hidden="true">
                <defs><linearGradient id="ringGradient" x1="0" y1="0" x2="1" y2="1"><stop stop-color="#22d3ee" /><stop offset="1" stop-color="#a78bfa" /></linearGradient></defs>
                <circle class="ring-track" cx="60" cy="60" r="50" />
                <circle class="ring-value" cx="60" cy="60" r="50" :stroke-dasharray="`${profile.quality.score * 3.1416} 314.16`" :stroke="qualityTone" />
              </svg>
              <div class="ring-center"><strong>{{ num(profile.quality.score, 0) }}</strong><small>score</small></div>
            </div>
            <div class="quality-meta">
              <div><span>Completeness</span><b>{{ profile.quality.completeness }}%</b></div>
              <div><span>Duplicate rows</span><b>{{ profile.quality.duplicate_rows }}</b></div>
              <div><span>Column issues</span><b>{{ profile.quality.issue_count }}</b></div>
            </div>
          </div>
        </section>
        <section class="panel insights-panel">
          <div class="panel-head"><Icon name="sparkles" :size="16" /><h2>Smart insights</h2><small>{{ profile?.insights.length ?? 0 }}</small></div>
          <InsightsWidget v-if="profile" :profile="profile" />
          <p v-else class="panel-desc">Load a dataset to see automatic findings about missing values, outliers and relationships.</p>
        </section>
        <div class="built-with">Go engine · Vue 3 · Apache ECharts</div>
      </aside>
    </main>

    <Transition name="toast">
      <div v-if="toast" class="toast" :class="toastTone"><Icon :name="toastTone === 'ok' ? 'check' : 'alert'" :size="15" /> {{ toast }}</div>
    </Transition>
  </div>
</template>

<style scoped>
.app { position: relative; z-index: 1; min-height: 100vh; display: flex; flex-direction: column; }
.topbar {
  display: flex; align-items: center; justify-content: space-between; gap: 16px;
  height: 72px; padding: 0 24px; border-bottom: 1px solid var(--border);
  background: rgba(7, 11, 21, .72); backdrop-filter: blur(16px);
}
.brand { display: flex; align-items: center; gap: 12px; }
.brand-mark { display: flex; filter: drop-shadow(0 0 12px rgba(34, 211, 238, .35)); }
.brand h1 { margin: 0; font-size: 18px; font-weight: 750; letter-spacing: -.02em; }
.brand h1 span { color: var(--accent); font-weight: 600; }
.brand p { margin: 2px 0 0; font-size: 10px; letter-spacing: .04em; color: var(--muted); }
.top-actions { display: flex; align-items: center; gap: 10px; }
.health { display: flex; align-items: center; gap: 7px; font-size: 10px; color: var(--muted); padding-right: 6px; }
.health i { width: 7px; height: 7px; border-radius: 50%; background: var(--good); box-shadow: 0 0 10px var(--good); }
.health i.down { background: var(--warn); box-shadow: 0 0 10px var(--warn); }
.error-banner {
  display: flex; align-items: center; gap: 9px; margin: 14px 24px 0; padding: 11px 14px;
  font-size: 12px; color: #ffd7df; border-radius: 11px;
  border: 1px solid rgba(251, 113, 133, .35); background: rgba(251, 113, 133, .1);
}
.workspace {
  display: grid; grid-template-columns: 272px minmax(0, 1fr) 306px; gap: 16px;
  padding: 18px 24px 24px; flex: 1; align-items: start;
}
.sidebar, .insight-column { display: flex; flex-direction: column; gap: 16px; min-width: 0; }
.palette-panel { flex: 1; }
.palette { display: flex; flex-direction: column; gap: 8px; }
.palette-card {
  display: flex; align-items: center; gap: 10px; padding: 10px; cursor: grab;
  border: 1px solid var(--border); border-radius: 12px; background: rgba(255, 255, 255, .022);
  transition: border-color .16s, transform .16s, background .16s;
}
.palette-card:hover { border-color: rgba(34, 211, 238, .4); background: rgba(34, 211, 238, .06); transform: translateX(3px); }
.palette-card:active { cursor: grabbing; }
.palette-icon { display: flex; color: var(--accent); }
.palette-text { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.palette-text strong { font-size: 11px; font-weight: 600; color: #dbe6f7; }
.palette-text small { font-size: 10px; color: #71829d; }
.palette-grip { margin-left: auto; color: #53627c; display: flex; }
.tip { display: flex; gap: 7px; margin-top: 14px; font-size: 10px; line-height: 1.5; color: #6f809b; }
.tip svg { flex: none; color: var(--accent2); }
.dataset-card { display: flex; flex-direction: column; gap: 4px; margin-top: 12px; padding: 11px; border-radius: 11px; border: 1px solid rgba(52, 211, 153, .2); background: rgba(52, 211, 153, .055); }
.dataset-top { display: flex; align-items: center; gap: 5px; font-size: 10px; color: var(--good); }
.dataset-card strong { font-size: 11px; color: var(--text); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.dataset-card span { font-size: 10px; color: var(--muted); }
.canvas-column { min-width: 0; }
.canvas-title { display: flex; align-items: flex-end; justify-content: space-between; gap: 12px; margin-bottom: 12px; }
.canvas-title h2 { margin: 0; font-size: 15px; font-weight: 650; letter-spacing: -.01em; }
.canvas-title p { margin: 4px 0 0; font-size: 11px; color: var(--muted); }
.canvas {
  display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 14px;
  align-content: start; min-height: calc(100vh - 150px); padding: 16px;
  border: 1px dashed rgba(148, 163, 184, .2); border-radius: 20px;
  background: rgba(9, 14, 26, .4); transition: background .18s, border-color .18s;
}
.canvas.over { border-color: var(--accent); background: rgba(34, 211, 238, .075); }
.busy, .welcome { grid-column: 1 / -1; display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 62vh; text-align: center; }
.busy strong { margin-top: 16px; font-size: 14px; }
.busy small { margin-top: 6px; font-size: 11px; color: var(--muted); }
.spinner { width: 34px; height: 34px; border-radius: 50%; border: 3px solid rgba(34, 211, 238, .2); border-top-color: var(--accent); animation: spin .8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.welcome-mark { display: grid; place-items: center; width: 76px; height: 76px; border-radius: 22px; color: var(--accent); border: 1px solid rgba(34, 211, 238, .28); background: radial-gradient(circle at 50% 20%, rgba(34, 211, 238, .18), rgba(34, 211, 238, .04)); box-shadow: 0 0 60px rgba(34, 211, 238, .14); }
.welcome h2 { margin: 22px 0 0; font-size: 23px; letter-spacing: -.025em; }
.welcome p { max-width: 470px; margin: 10px 0 0; font-size: 12px; line-height: 1.7; color: var(--muted); }
.steps { display: flex; align-items: center; gap: 12px; margin: 24px 0; font-size: 11px; color: #a7b6ce; }
.steps span { display: flex; align-items: center; gap: 6px; }
.steps b { color: var(--accent); font-size: 10px; letter-spacing: .06em; }
.steps i { width: 22px; height: 1px; background: rgba(148, 163, 184, .3); }
.quality { display: flex; align-items: center; gap: 14px; }
.ring-wrap { position: relative; flex: none; width: 104px; height: 104px; }
.ring { width: 104px; height: 104px; transform: rotate(-90deg); }
.ring circle { fill: none; stroke-width: 9; }
.ring-track { stroke: rgba(148, 163, 184, .13); }
.ring-value { stroke-linecap: round; transition: stroke-dasharray .5s ease; }
.ring-center { position: absolute; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; }
.ring-center strong { font-size: 24px; font-weight: 700; letter-spacing: -.04em; }
.ring-center small { font-size: 9px; letter-spacing: .1em; text-transform: uppercase; color: var(--muted); }
.quality-meta { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 8px; }
.quality-meta div { display: flex; justify-content: space-between; gap: 8px; font-size: 10px; color: var(--muted); }
.quality-meta b { color: #dce6f7; font-variant-numeric: tabular-nums; }
.insights-panel { min-height: 220px; }
.built-with { margin-top: auto; padding: 4px 2px; font-size: 10px; letter-spacing: .06em; color: #5c6d88; }
.toast { position: fixed; right: 26px; bottom: 24px; z-index: 10; display: flex; align-items: center; gap: 9px; padding: 12px 16px; font-size: 12px; color: var(--text); border-radius: 12px; border: 1px solid rgba(52, 211, 153, .3); background: rgba(13, 30, 32, .95); box-shadow: 0 18px 40px rgba(0, 0, 0, .4); }
.toast svg { color: var(--good); }
.toast.error { border-color: rgba(251, 113, 133, .35); background: rgba(40, 16, 25, .95); }
.toast.error svg { color: var(--danger); }
.toast-enter-active, .toast-leave-active { transition: opacity .22s, transform .22s; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(10px); }

@media (max-width: 1300px) {
  .workspace { grid-template-columns: 250px minmax(0, 1fr); }
  .insight-column { display: none; }
}
@media (max-width: 850px) {
  .workspace { grid-template-columns: minmax(0, 1fr); padding: 14px; }
  .sidebar { flex-direction: row; }
  .sidebar .panel { flex: 1; }
  .health { display: none; }
  .canvas { grid-template-columns: minmax(0, 1fr); }
}
</style>
