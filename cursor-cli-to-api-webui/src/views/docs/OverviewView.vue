<script setup lang="ts">
import Card from '@/components/ui/card/Card.vue'
</script>

<template>
  <div class="space-y-4 prose prose-slate max-w-none">
    <h1 class="text-2xl font-semibold tracking-tight">Docs — overview</h1>
    <Card class="p-4 space-y-3 text-sm leading-relaxed">
      <p>
        This API is an OpenAI-compatible front door. Point any client
        <code class="rounded bg-muted px-1">base_url</code> at the API host
        (default <code class="rounded bg-muted px-1">http://127.0.0.1:8201/v1</code>).
      </p>
      <p>
        Each request is turned into a Cursor CLI prompt
        (<code class="rounded bg-muted px-1">agent -p --mode ask</code>) and the CLI text is wrapped in OpenAI-shaped JSON.
      </p>
      <ul class="list-disc pl-5 space-y-1">
        <li>Auth: HTTP Basic <code>armin:dopadopa123</code> or Bearer <code>dopadopa123</code></li>
        <li>Prerequisite: Cursor Agent logged in on the host</li>
        <li>Set <code>CURSOR_AGENT_BIN</code> if PATH <code>agent</code> is not Cursor</li>
      </ul>
    </Card>
  </div>
</template>
