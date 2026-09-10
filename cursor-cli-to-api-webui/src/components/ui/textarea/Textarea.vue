<script setup lang="ts">
import { cn } from '@/lib/utils'

defineProps<{ modelValue?: string; rows?: number; placeholder?: string; disabled?: boolean }>()
defineEmits<{ 'update:modelValue': [string] }>()
</script>

<template>
  <textarea
    :value="modelValue"
    :rows="rows || 6"
    :placeholder="placeholder"
    :disabled="disabled"
    :class="cn('flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring')"
    @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
  />
</template>
