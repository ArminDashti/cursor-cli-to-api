<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from '@/components/ui/card/Card.vue'
import Label from '@/components/ui/label/Label.vue'
import Input from '@/components/ui/input/Input.vue'
import Button from '@/components/ui/button/Button.vue'
import { apiFetch, saveAuth } from '@/lib/utils'

const router = useRouter()
const route = useRoute()
const username = ref('armin')
const password = ref('dopadopa123')
const error = ref('')
const loading = ref(false)

async function onSubmit() {
  error.value = ''
  loading.value = true
  saveAuth(username.value, password.value)
  try {
    const { res, json } = await apiFetch('/health')
    // health is public; verify auth against models
    const authCheck = await apiFetch('/v1/models')
    if (!authCheck.res.ok) {
      error.value = typeof authCheck.json === 'object' && authCheck.json && 'error' in (authCheck.json as object)
        ? JSON.stringify((authCheck.json as { error: unknown }).error)
        : 'Login failed'
      loading.value = false
      return
    }
    void res
    void json
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/playground/chat'
    await router.replace(redirect)
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto flex min-h-full max-w-md items-center">
    <Card class="w-full p-6 space-y-4">
      <div>
        <h1 class="text-xl font-semibold">Sign in</h1>
        <p class="text-sm text-muted-foreground mt-1">Default: armin / dopadopa123</p>
      </div>
      <form class="space-y-3" @submit.prevent="onSubmit">
        <div class="space-y-1">
          <Label forId="user">Username</Label>
          <Input id="user" v-model="username" autocomplete="username" />
        </div>
        <div class="space-y-1">
          <Label forId="pass">Password</Label>
          <Input id="pass" v-model="password" type="password" autocomplete="current-password" />
        </div>
        <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
        <Button type="submit" :disabled="loading" class="w-full">{{ loading ? 'Signing in…' : 'Sign in' }}</Button>
      </form>
    </Card>
  </div>
</template>
