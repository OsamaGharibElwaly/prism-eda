<script setup lang="ts">
import { ref } from 'vue'
import Icon from './Icon.vue'

const emit = defineEmits<{ file: [file: File] }>()
const over = ref(false)
const input = ref<HTMLInputElement | null>(null)

function take(files: FileList | null) {
  const file = files?.[0]
  if (file) emit('file', file)
}

function onDrop(event: DragEvent) {
  over.value = false
  take(event.dataTransfer?.files ?? null)
}

function onChange(event: Event) {
  const target = event.target as HTMLInputElement
  take(target.files)
  target.value = ''
}
</script>

<template>
  <div class="dropzone" :class="{ over }" role="button" tabindex="0" aria-label="Drop a CSV file or click to browse"
    @click="input?.click()" @keydown.enter="input?.click()" @keydown.space.prevent="input?.click()"
    @dragenter.prevent="over = true" @dragover.prevent="over = true" @dragleave.prevent="over = false"
    @drop.prevent="onDrop">
    <input ref="input" type="file" accept=".csv,.tsv,.txt,text/csv" hidden @change="onChange" />
    <div class="drop-icon"><Icon name="upload" :size="22" /></div>
    <strong>Drop your dataset</strong>
    <span>or <em>browse files</em></span>
    <small>CSV · TSV · up to 25 MB</small>
  </div>
</template>

<style scoped>
.dropzone {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  gap: 6px; min-height: 168px; padding: 18px; cursor: pointer; text-align: center;
  border: 1px dashed rgba(34, 211, 238, .38); border-radius: 16px;
  background: radial-gradient(circle at 50% 0%, rgba(34, 211, 238, .09), transparent 70%), rgba(15, 23, 42, .45);
  transition: border-color .18s, background .18s, transform .18s;
}
.dropzone:hover, .dropzone.over {
  border-color: var(--accent); background: rgba(34, 211, 238, .12); transform: translateY(-2px);
}
.dropzone:focus-visible { outline: 2px solid var(--accent); outline-offset: 3px; }
.drop-icon {
  display: grid; place-items: center; width: 44px; height: 44px; margin-bottom: 4px;
  border-radius: 13px; color: var(--accent); background: rgba(34, 211, 238, .12);
  border: 1px solid rgba(34, 211, 238, .22);
}
strong { font-size: 13px; color: var(--text); }
span { font-size: 12px; color: var(--muted); }
em { color: var(--accent); font-style: normal; }
small { margin-top: 6px; font-size: 10px; letter-spacing: .06em; text-transform: uppercase; color: #61728e; }
</style>
