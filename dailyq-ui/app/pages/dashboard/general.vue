<script setup lang="ts">
const { generalTasks, generalNotes } = usePlaceholderJournal()
</script>

<template>
  <DashboardItemBoard
    title="General"
    description="Tasks and notes without a date — capture anytime."
    :tasks="generalTasks"
    :notes="generalNotes"
    :date="null"
  />
</template>
