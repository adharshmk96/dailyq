<script setup lang="ts">
const { activeTagId } = useTagFilter()
</script>

<template>
  <div class="grid gap-6 lg:grid-cols-[240px_1fr] lg:items-start">
    <DashboardTagPanel />

    <DashboardItemBoard
      title="General"
      description="Tasks and notes without a date — capture anytime."
      :date="null"
      :active-tag-id="activeTagId"
    />
  </div>
</template>
