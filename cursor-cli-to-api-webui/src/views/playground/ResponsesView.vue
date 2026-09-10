<script setup lang="ts">
import { ref } from 'vue'
import Card from '@/components/ui/card/Card.vue'
import Label from '@/components/ui/label/Label.vue'
import Input from '@/components/ui/input/Input.vue'
import Textarea from '@/components/ui/textarea/Textarea.vue'
import Button from '@/components/ui/button/Button.vue'
import { apiFetch } from '@/lib/utils'

const model = ref('auto')
const input = ref('Explain what this proxy does in two sentences.')
const responseId = ref('')
const loading = ref(false)
const output = ref('')

async function create() {
  loading.value = true
  output.value = ''
  try {
    const { text, json } = await apiFetch('/v1/responses', {
      method: 'POST',
      body: JSON.stringify({
        model: model.value,
        input: input.value,
      }),
    })
    output.value = text
    if (json && typeof json === 'object' && json !== null && 'id' in json) {
      responseId.value = String((json as { id: string }).id)
    }
  } catch (e) {
    output.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

async function fetchById() {
  if (!responseId.value) return
  loading.value = true
  try {
    const { text } = await apiFetch(`/v1/responses/${encodeURIComponent(responseId.value)}`)
    output.value = text
  } catch (e) {
    output.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="space-y-4">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">Playground — responses</h1>
      <p class="text-sm text-muted-foreground">POST /v1/responses and GET /v1/responses/{id}</p>
    </div>
    <Card class="p-4 space-y-3">
      <div class="space-y-1">
        <Label>Model</Label>
        <Input v-model="model" placeholder="auto" />
      </div>
      <div class="space-y-1">
        <Label>Input</Label>
        <Textarea v-model="input" :rows="5" />
      </div>
      <div class="flex flex-wrap gap-2">
        <Button :disabled="loading" @click="create">{{ loading ? 'Running…' : 'Create response' }}</Button>
      </div>
      <div class="space-y-1">
        <Label>Response id</Label>
        <div class="flex gap-2">
          <Input v-model="responseId" placeholder="resp_…" />
          <Button variant="secondary" :disabled="loading || !responseId" @click="fetchById">GET</Button>
        </div>
      </div>
    </Card>
    <Card class="p-4">
      <h2 class="text-sm font-medium mb-2">Raw response</h2>
      <pre class="overflow-auto rounded-md bg-slate-950 text-slate-50 p-3 text-xs whitespace-pre-wrap">{{ output || '—' }}</pre>
    </Card>
  </div>
</template>
