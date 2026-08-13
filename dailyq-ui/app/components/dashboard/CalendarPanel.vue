<script setup lang="ts">
import { CalendarDate, getLocalTimeZone, today } from '@internationalized/date'
import type { DateValue } from '@internationalized/date'

const journal = usePlaceholderJournal()
const { datesWithItems, getTasksForDate, getNotesForDate, filterItemsByTag } = journal
const { activeTagId } = useTagFilter()

const selected = shallowRef<CalendarDate>(today(getLocalTimeZone()))

const selectedIso = computed(() => {
  const d = selected.value
  if (!d) return journal.today.value
  return `${d.year}-${String(d.month).padStart(2, '0')}-${String(d.day).padStart(2, '0')}`
})

const dayTasks = computed(() =>
  filterItemsByTag(getTasksForDate(selectedIso.value), activeTagId.value)
)
const dayNotes = computed(() =>
  filterItemsByTag(getNotesForDate(selectedIso.value), activeTagId.value)
)

const formattedDate = computed(() => {
  const [y, m, day] = selectedIso.value.split('-').map(Number)
  return new Date(y, m - 1, day).toLocaleDateString(undefined, {
    weekday: 'long',
    month: 'long',
    day: 'numeric'
  })
})

function dateKey(date: DateValue) {
  return `${date.year}-${String(date.month).padStart(2, '0')}-${String(date.day).padStart(2, '0')}`
}

function hasItems(date: DateValue) {
  return datesWithItems.value.has(dateKey(date))
}
</script>

<template>
  <div class="space-y-5">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight text-highlighted">
        Calendar
      </h1>
      <p class="mt-1 text-sm text-muted">
        Time-bound tasks and notes for a specific day.
      </p>
    </div>

    <div class="grid gap-6 lg:grid-cols-[auto_1fr] lg:items-start">
      <div class="space-y-4">
        <UCard :ui="{ body: 'flex justify-center p-4 sm:p-5' }">
          <UCalendar
            v-model="selected"
            size="lg"
            color="primary"
            variant="soft"
          >
            <template #day="{ day }">
              <span class="relative flex size-full items-center justify-center">
                {{ day.day }}
                <span
                  v-if="hasItems(day)"
                  class="absolute bottom-0.5 size-1 rounded-full bg-primary"
                />
              </span>
            </template>
          </UCalendar>
        </UCard>

        <DashboardTagPanel />
      </div>

      <DashboardItemBoard
        :title="formattedDate"
        description="Add and manage items for the selected date."
        :tasks="dayTasks"
        :notes="dayNotes"
        :date="selectedIso"
        :active-tag-id="activeTagId"
      />
    </div>
  </div>
</template>
