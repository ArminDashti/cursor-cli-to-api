<script setup lang="ts">
import { cn } from '@/lib/utils'

defineProps<{ modelValue?: string; type?: string; placeholder?: string; disabled?: boolean }>()
defineEmits<{ 'update:modelValue': [string] }>()
</script>

<template>
  <input
    :type="type || 'text'"
    :value="modelValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :class="cn('flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring')"
    @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
  />
</template>
