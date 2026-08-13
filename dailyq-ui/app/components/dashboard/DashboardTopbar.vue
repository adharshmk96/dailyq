<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'

interface DashboardNavItem {
  label: string
  icon: string
  value: string
  to: string
}

defineProps<{
  activeSection: string | null
  items: DashboardNavItem[]
}>()

const toast = useToast()

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
        navigateTo('/dashboard/settings/account')
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
        <nav
          class="flex max-w-full items-center gap-1 rounded-lg bg-elevated/80 p-1"
          aria-label="Dashboard sections"
        >
          <NuxtLink
            v-for="item in items"
            :key="item.value"
            :to="item.to"
            class="flex items-center gap-1.5 rounded-md px-2.5 py-1.5 text-sm font-medium transition-colors sm:px-3"
            :class="activeSection === item.value
              ? 'bg-default text-highlighted shadow-sm'
              : 'text-muted hover:bg-default/60 hover:text-highlighted'"
            :aria-current="activeSection === item.value ? 'page' : undefined"
          >
            <UIcon
              :name="item.icon"
              class="size-4 shrink-0"
            />
            <span class="hidden sm:inline">{{ item.label }}</span>
          </NuxtLink>
        </nav>
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
