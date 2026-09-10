<script setup lang="ts">
import { ref } from 'vue'
import Card from '@/components/ui/card/Card.vue'
import Label from '@/components/ui/label/Label.vue'
import Input from '@/components/ui/input/Input.vue'
import Textarea from '@/components/ui/textarea/Textarea.vue'
import Button from '@/components/ui/button/Button.vue'
import { apiFetch } from '@/lib/utils'

const model = ref('auto')
const message = ref('Say hello in one short sentence.')
const loading = ref(false)
const output = ref('')

async function run() {
  loading.value = true
  output.value = ''
  try {
    const { res, text } = await apiFetch('/v1/chat/completions', {
      method: 'POST',
      body: JSON.stringify({
        model: model.value,
        messages: [{ role: 'user', content: message.value }],
        stream: false,
      }),
    })
    output.value = text
    if (!res.ok) {
      /* still show body */
    }
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
      <h1 class="text-2xl font-semibold tracking-tight">Playground — chat completions</h1>
      <p class="text-sm text-muted-foreground">POST /v1/chat/completions → Cursor CLI → OpenAI-shaped JSON</p>
    </div>
    <Card class="p-4 space-y-3">
      <div class="space-y-1">
        <Label>Model</Label>
        <Input v-model="model" placeholder="auto" />
      </div>
      <div class="space-y-1">
        <Label>User message</Label>
        <Textarea v-model="message" :rows="5" />
      </div>
      <Button :disabled="loading" @click="run">{{ loading ? 'Running…' : 'Send' }}</Button>
    </Card>
    <Card class="p-4">
      <h2 class="text-sm font-medium mb-2">Raw response</h2>
      <pre class="overflow-auto rounded-md bg-slate-950 text-slate-50 p-3 text-xs whitespace-pre-wrap">{{ output || '—' }}</pre>
    </Card>
  </div>
</template>
