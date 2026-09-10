<script setup lang="ts">
import Card from '@/components/ui/card/Card.vue'
</script>

<template>
  <div class="space-y-4">
    <h1 class="text-2xl font-semibold tracking-tight">POST /v1/chat/completions</h1>
    <Card class="p-4 space-y-3 text-sm">
      <p>Request body (OpenAI chat completions):</p>
      <pre class="overflow-auto rounded-md bg-slate-950 text-slate-50 p-3 text-xs">{
  "model": "auto",
  "messages": [
    { "role": "user", "content": "Hello" }
  ],
  "stream": false
}</pre>
      <p>Response object fields include <code>id</code>, <code>object: "chat.completion"</code>, <code>created</code>, <code>model</code>, <code>choices[].message</code>, <code>usage</code>.</p>
      <p>Set <code>stream: true</code> for SSE chunks ending with <code>data: [DONE]</code>.</p>
    </Card>
  </div>
</template>
