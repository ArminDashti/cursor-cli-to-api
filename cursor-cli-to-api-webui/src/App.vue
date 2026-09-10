<script setup lang="ts">
import { useRouter, useRoute, RouterLink, RouterView } from 'vue-router'
import { clearAuth } from '@/lib/utils'
import Button from '@/components/ui/button/Button.vue'

const router = useRouter()
const route = useRoute()

function logout() {
  clearAuth()
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="min-h-full flex flex-col">
    <header v-if="route.name !== 'login'" class="border-b bg-card">
      <div class="mx-auto flex max-w-6xl items-center justify-between gap-4 px-4 py-3">
        <div class="flex items-center gap-6">
          <RouterLink to="/" class="font-semibold tracking-tight">Cursor CLI to API</RouterLink>
          <nav class="flex gap-3 text-sm">
            <RouterLink class="text-muted-foreground hover:text-foreground" active-class="!text-foreground font-medium" to="/playground/chat">Playground</RouterLink>
            <RouterLink class="text-muted-foreground hover:text-foreground" active-class="!text-foreground font-medium" to="/docs/overview">Docs</RouterLink>
          </nav>
        </div>
        <Button variant="outline" @click="logout">Sign out</Button>
      </div>
      <div v-if="route.path.startsWith('/playground')" class="mx-auto flex max-w-6xl gap-3 px-4 pb-3 text-sm">
        <RouterLink class="rounded-md px-2 py-1 hover:bg-accent" active-class="bg-accent font-medium" to="/playground/chat">Chat completions</RouterLink>
        <RouterLink class="rounded-md px-2 py-1 hover:bg-accent" active-class="bg-accent font-medium" to="/playground/responses">Responses</RouterLink>
      </div>
      <div v-else-if="route.path.startsWith('/docs')" class="mx-auto flex max-w-6xl gap-3 px-4 pb-3 text-sm">
        <RouterLink class="rounded-md px-2 py-1 hover:bg-accent" active-class="bg-accent font-medium" to="/docs/overview">Overview</RouterLink>
        <RouterLink class="rounded-md px-2 py-1 hover:bg-accent" active-class="bg-accent font-medium" to="/docs/chat-completions">POST /v1/chat/completions</RouterLink>
        <RouterLink class="rounded-md px-2 py-1 hover:bg-accent" active-class="bg-accent font-medium" to="/docs/responses">POST /v1/responses</RouterLink>
      </div>
    </header>
    <main class="flex-1 mx-auto w-full max-w-6xl px-4 py-6">
      <RouterView />
    </main>
  </div>
</template>
