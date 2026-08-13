<script setup lang="ts">
const { generalTasks, generalNotes, filterItemsByTag } = usePlaceholderJournal()
const { activeTagId } = useTagFilter()

const tasks = computed(() => filterItemsByTag(generalTasks.value, activeTagId.value))
const notes = computed(() => filterItemsByTag(generalNotes.value, activeTagId.value))
</script>

<template>
  <div class="grid gap-6 lg:grid-cols-[240px_1fr] lg:items-start">
    <DashboardTagPanel />

    <DashboardItemBoard
      title="General"
      description="Tasks and notes without a date — capture anytime."
      :tasks="tasks"
      :notes="notes"
      :date="null"
      :active-tag-id="activeTagId"
    />
  </div>
</template>
