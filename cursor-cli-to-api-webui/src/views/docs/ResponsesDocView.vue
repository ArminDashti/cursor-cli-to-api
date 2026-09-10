<script setup lang="ts">
import Card from '@/components/ui/card/Card.vue'
</script>

<template>
  <div class="space-y-4">
    <h1 class="text-2xl font-semibold tracking-tight">POST /v1/responses</h1>
    <Card class="p-4 space-y-3 text-sm">
      <p>Create:</p>
      <pre class="overflow-auto rounded-md bg-slate-950 text-slate-50 p-3 text-xs">{
  "model": "auto",
  "input": "Hello"
}</pre>
      <p>Response object fields include <code>id</code>, <code>object: "response"</code>, <code>status</code>, <code>output[]</code>, <code>usage</code>.</p>
      <p>Fetch again with <code>GET /v1/responses/{id}</code> (stored in Postgres).</p>
    </Card>
  </div>
</template>
