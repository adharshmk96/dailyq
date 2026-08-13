<script setup lang="ts">
const route = useRoute()

definePageMeta({
  middleware: 'auth'
})

const tabItems = [
  {
    label: 'Overview',
    icon: 'i-lucide-layout-dashboard',
    value: 'overview',
    to: '/dashboard/overview'
  },
  {
    label: 'General',
    icon: 'i-lucide-layers',
    value: 'general',
    to: '/dashboard/general'
  },
  {
    label: 'Calendar',
    icon: 'i-lucide-calendar-days',
    value: 'calendar',
    to: '/dashboard/calendar'
  }
]

const activeSection = computed(() => {
  const segment = route.path.split('/')[2]
  if (segment === 'general' || segment === 'calendar' || segment === 'overview') {
    return segment
  }

  return null
})

useHead({
  title: 'Dashboard — DailyQ'
})
</script>

<template>
  <div class="min-h-dvh bg-muted">
    <DashboardTopbar
      :active-section="activeSection"
      :items="tabItems"
    />

    <main class="mx-auto max-w-6xl px-4 py-6 sm:px-6 sm:py-8">
      <NuxtPage />
    </main>
  </div>
</template>
