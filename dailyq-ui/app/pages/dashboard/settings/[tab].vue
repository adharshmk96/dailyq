<script setup lang="ts">
const VALID_TABS = ['account', 'data', 'notification'] as const

const route = useRoute()

const tab = computed(() => {
  const param = route.params.tab
  return typeof param === 'string' ? param : ''
})

watch(tab, (value) => {
  if (!value || !(VALID_TABS as readonly string[]).includes(value)) {
    navigateTo('/dashboard/settings/account', { replace: true })
  }
}, { immediate: true })

useHead({
  title: 'Settings — DailyQ'
})
</script>

<template>
  <DashboardSettingsPanel />
</template>
