<script setup lang="ts">
const props = withDefaults(defineProps<{
  title: string
  summary?: string
  icon?: string
  /** Collapsed by default on mobile. Desktop is always expanded. */
  defaultCollapsed?: boolean
}>(), {
  summary: undefined,
  icon: undefined,
  defaultCollapsed: true
})

const open = ref(!props.defaultCollapsed)

function toggle() {
  open.value = !open.value
}
</script>

<template>
  <div class="overflow-hidden rounded-xl border border-default bg-default lg:border-0 lg:bg-transparent lg:overflow-visible">
    <button
      type="button"
      class="flex w-full items-center gap-2 px-3 py-2.5 text-left transition hover:bg-elevated/50 lg:hidden"
      :aria-expanded="open"
      @click="toggle"
    >
      <UIcon
        v-if="icon"
        :name="icon"
        class="size-4 shrink-0 text-muted"
      />
      <span class="min-w-0 flex-1">
        <span class="block text-sm font-semibold text-highlighted">{{ title }}</span>
        <span
          v-if="summary && !open"
          class="block truncate text-xs text-muted"
        >
          {{ summary }}
        </span>
      </span>
      <UIcon
        :name="open ? 'i-lucide-chevron-up' : 'i-lucide-chevron-down'"
        class="size-4 shrink-0 text-muted"
      />
    </button>

    <div
      class="lg:block"
      :class="open ? 'block' : 'hidden lg:block'"
    >
      <slot />
    </div>
  </div>
</template>
