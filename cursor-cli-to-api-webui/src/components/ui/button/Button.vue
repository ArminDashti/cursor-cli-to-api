<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    variant?: 'default' | 'secondary' | 'outline'
    type?: 'button' | 'submit' | 'reset'
    disabled?: boolean
  }>(),
  { variant: 'default', type: 'button', disabled: false },
)

const classes = computed(() =>
  cn(
    'inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors h-10 px-4 py-2 disabled:opacity-50',
    props.variant === 'default' && 'bg-primary text-primary-foreground hover:opacity-90',
    props.variant === 'secondary' && 'bg-secondary text-secondary-foreground hover:opacity-90',
    props.variant === 'outline' && 'border border-input bg-background hover:bg-accent',
  ),
)
</script>

<template>
  <button :type="type" :disabled="disabled" :class="classes">
    <slot />
  </button>
</template>
