<script setup lang="ts">
import type { DropdownMenuItem, TabsItem } from '@nuxt/ui'

const props = defineProps<{
  modelValue: string
  items: TabsItem[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const toast = useToast()

const activeTab = computed({
  get: () => props.modelValue,
  set: (value: string) => emit('update:modelValue', value)
})

const userMenuItems = computed<DropdownMenuItem[][]>(() => [
  [
    {
      label: 'Alex Rivera',
      type: 'label'
    }
  ],
  [
    {
      label: 'Settings',
      icon: 'i-lucide-settings',
      onSelect() {
        toast.add({
          title: 'Settings',
          description: 'Coming soon',
          icon: 'i-lucide-settings'
        })
      }
    }
  ],
  [
    {
      label: 'Logout',
      icon: 'i-lucide-log-out',
      color: 'error',
      onSelect() {
        toast.add({
          title: 'Logout',
          description: 'Coming soon',
          icon: 'i-lucide-log-out',
          color: 'error'
        })
      }
    }
  ]
])
</script>

<template>
  <header class="sticky top-0 z-40 border-b border-default bg-default/80 backdrop-blur-md">
    <div class="mx-auto flex h-16 max-w-6xl items-center gap-3 px-4 sm:px-6">
      <NuxtLink
        to="/dashboard/overview"
        class="flex shrink-0 items-center gap-2.5"
      >
        <span class="flex size-9 items-center justify-center rounded-xl bg-primary text-inverted shadow-sm">
          <UIcon
            name="i-lucide-book-open"
            class="size-5"
          />
        </span>
        <span class="hidden text-lg font-semibold tracking-tight text-highlighted sm:inline">
          DailyQ
        </span>
      </NuxtLink>

      <div class="flex min-w-0 flex-1 justify-center">
        <UTabs
          v-model="activeTab"
          :items="items"
          :content="false"
          color="primary"
          variant="pill"
          size="sm"
          class="w-auto max-w-full"
          :ui="{
            list: 'bg-elevated/80',
            trigger: 'px-2.5 sm:px-3',
            label: 'hidden sm:inline'
          }"
        />
      </div>

      <div class="flex shrink-0 items-center gap-1.5">
        <UColorModeButton />

        <UDropdownMenu
          :items="userMenuItems"
          :content="{ align: 'end' }"
          :ui="{ content: 'w-48' }"
        >
          <UButton
            color="neutral"
            variant="ghost"
            class="rounded-full p-0.5"
            aria-label="User menu"
          >
            <UAvatar
              src="https://i.pravatar.cc/80?u=dailyq"
              alt="Alex Rivera"
              size="sm"
            />
          </UButton>
        </UDropdownMenu>
      </div>
    </div>
  </header>
</template>
