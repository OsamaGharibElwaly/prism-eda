<script setup lang="ts">
import Icon from './Icon.vue'

const { title, icon, index, moving = false } = defineProps<{ title: string; icon: string; index: number; moving?: boolean }>()
const emit = defineEmits<{
  remove: []
  dragstart: [event: DragEvent]
  dragend: []
  drop: [event: DragEvent]
}>()
</script>

<template>
  <article class="widget" :class="{ moving }" @dragover.prevent.stop @drop.stop="emit('drop', $event)">
    <header class="widget-head" draggable="true" @dragstart="emit('dragstart', $event)" @dragend="emit('dragend')">
      <span class="grip"><Icon name="grip" :size="16" /></span>
      <span class="widget-icon"><Icon :name="icon" :size="16" /></span>
      <h3>{{ title }}</h3>
      <span class="widget-index">{{ String(index + 1).padStart(2, '0') }}</span>
      <button class="icon-btn" type="button" aria-label="Remove widget" @click="emit('remove')">
        <Icon name="close" :size="15" />
      </button>
    </header>
    <div class="widget-body"><slot /></div>
  </article>
</template>

<style scoped>
.widget {
  min-width: 0; border-radius: 18px; overflow: hidden;
  border: 1px solid var(--border); background: linear-gradient(160deg, rgba(23, 32, 51, .95), rgba(12, 18, 32, .95));
  box-shadow: 0 18px 44px rgba(0, 0, 0, .28); transition: opacity .18s, border-color .18s, transform .18s;
}
.widget:hover { border-color: rgba(34, 211, 238, .3); }
.widget.moving { opacity: .4; }
.widget-head {
  display: flex; align-items: center; gap: 10px; padding: 12px 14px; cursor: grab;
  border-bottom: 1px solid var(--border); background: rgba(255, 255, 255, .015);
}
.widget-head:active { cursor: grabbing; }
.grip { color: #54657f; display: flex; }
.widget-icon { display: flex; color: var(--accent); }
h3 { margin: 0; font-size: 12px; font-weight: 650; letter-spacing: .015em; color: #dce6f7; flex: 1; }
.widget-index { font-size: 10px; color: #52627c; font-variant-numeric: tabular-nums; }
.icon-btn {
  display: grid; place-items: center; width: 24px; height: 24px; padding: 0; cursor: pointer;
  color: #7c8ca7; border: 0; border-radius: 7px; background: transparent;
}
.icon-btn:hover { color: var(--danger); background: rgba(251, 113, 133, .12); }
.widget-body { padding: 16px; }
</style>
