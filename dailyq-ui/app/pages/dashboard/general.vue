<script setup lang="ts">
const { activeTagId } = useTagFilter()
const { tags } = useJournalTags()

const tagSummary = computed(() => {
  if (activeTagId.value) {
    const name = tags.value.find(tag => tag.id === activeTagId.value)?.name
    return name ? `Filtering: ${name}` : 'Filter active'
  }
  if (!tags.value.length) return 'No tags'
  return `${tags.value.length} tag${tags.value.length === 1 ? '' : 's'}`
})
</script>

<template>
  <div class="grid gap-4 lg:grid-cols-[240px_1fr] lg:items-start lg:gap-6">
    <DashboardCollapsibleSection
      title="Tags"
      :summary="tagSummary"
      icon="i-lucide-tags"
      :default-collapsed="true"
    >
      <DashboardTagPanel embedded />
    </DashboardCollapsibleSection>

    <DashboardItemBoard
      title="General"
      description="Tasks and notes without a date — capture anytime."
      :date="null"
      :active-tag-id="activeTagId"
      mobile-priority
    />
  </div>
</template>
