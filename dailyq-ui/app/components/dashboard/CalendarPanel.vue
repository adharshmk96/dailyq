<script setup lang="ts">
import { CalendarDate, getLocalTimeZone, today } from '@internationalized/date'
import type { DateValue } from '@internationalized/date'

const { datesWithItems } = useJournalDates()
const { activeTagId } = useTagFilter()
const { tags } = useJournalTags()

const selected = shallowRef<CalendarDate>(today(getLocalTimeZone()))
const calendarOpen = ref(false)

const tagSummary = computed(() => {
  if (activeTagId.value) {
    const name = tags.value.find(tag => tag.id === activeTagId.value)?.name
    return name ? `Filtering: ${name}` : 'Filter active'
  }
  if (!tags.value.length) return 'No tags'
  return `${tags.value.length} tag${tags.value.length === 1 ? '' : 's'}`
})

const selectedIso = computed(() => {
  const d = selected.value ?? today(getLocalTimeZone())
  return `${d.year}-${String(d.month).padStart(2, '0')}-${String(d.day).padStart(2, '0')}`
})

const formattedDate = computed(() => {
  const [y, m, day] = selectedIso.value.split('-').map(Number)
  return new Date(y, m - 1, day).toLocaleDateString(undefined, {
    weekday: 'long',
    month: 'long',
    day: 'numeric'
  })
})

const compactDate = computed(() => {
  const [y, m, day] = selectedIso.value.split('-').map(Number)
  return new Date(y, m - 1, day).toLocaleDateString(undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric'
  })
})

function dateKey(date: DateValue) {
  return `${date.year}-${String(date.month).padStart(2, '0')}-${String(date.day).padStart(2, '0')}`
}

function hasItems(date: DateValue) {
  return datesWithItems.value.has(dateKey(date))
}

function shiftDay(delta: number) {
  selected.value = selected.value.add({ days: delta })
}

function goToToday() {
  selected.value = today(getLocalTimeZone())
}
</script>

<template>
  <div class="space-y-4 lg:space-y-5">
    <div class="hidden sm:block">
      <h1 class="text-2xl font-semibold tracking-tight text-highlighted">
        Calendar
      </h1>
      <p class="mt-1 text-sm text-muted">
        Time-bound tasks and notes for a specific day.
      </p>
    </div>

    <div class="grid gap-4 lg:grid-cols-[auto_1fr] lg:items-start lg:gap-6">
      <div class="space-y-3 lg:space-y-4">
        <!-- Mobile: compact date bar when calendar is collapsed -->
        <div
          v-if="!calendarOpen"
          class="flex items-center gap-1 rounded-xl border border-default bg-default p-1.5 lg:hidden"
        >
          <UButton
            icon="i-lucide-chevron-left"
            color="neutral"
            variant="ghost"
            size="sm"
            aria-label="Previous day"
            @click="shiftDay(-1)"
          />
          <button
            type="button"
            class="min-w-0 flex-1 rounded-lg px-2 py-1.5 text-center transition hover:bg-elevated/50"
            @click="calendarOpen = true"
          >
            <span class="block truncate text-sm font-semibold text-highlighted">
              {{ compactDate }}
            </span>
            <span
              v-if="hasItems(selected)"
              class="mt-0.5 inline-flex items-center gap-1 text-xs text-muted"
            >
              <span class="size-1.5 rounded-full bg-primary" />
              Has items
            </span>
          </button>
          <UButton
            icon="i-lucide-chevron-right"
            color="neutral"
            variant="ghost"
            size="sm"
            aria-label="Next day"
            @click="shiftDay(1)"
          />
          <UButton
            icon="i-lucide-calendar-days"
            color="neutral"
            variant="ghost"
            size="sm"
            aria-label="Open calendar"
            @click="calendarOpen = true"
          />
        </div>

        <!-- Mobile: expanded calendar -->
        <div
          v-if="calendarOpen"
          class="overflow-hidden rounded-xl border border-default bg-default lg:hidden"
        >
          <div class="flex items-center justify-between gap-2 border-b border-default px-3 py-2">
            <span class="text-sm font-semibold text-highlighted">Calendar</span>
            <div class="flex items-center gap-1">
              <UButton
                label="Today"
                color="neutral"
                variant="ghost"
                size="xs"
                @click="goToToday"
              />
              <UButton
                icon="i-lucide-chevron-up"
                color="neutral"
                variant="ghost"
                size="xs"
                aria-label="Collapse calendar"
                @click="calendarOpen = false"
              />
            </div>
          </div>
          <div class="flex justify-center p-3">
            <UCalendar
              v-model="selected"
              size="md"
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
          </div>
        </div>

        <!-- Desktop: always visible calendar -->
        <UCard
          class="hidden lg:block"
          :ui="{ body: 'flex justify-center p-4 sm:p-5' }"
        >
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

        <DashboardCollapsibleSection
          title="Tags"
          :summary="tagSummary"
          icon="i-lucide-tags"
          :default-collapsed="true"
        >
          <DashboardTagPanel embedded />
        </DashboardCollapsibleSection>
      </div>

      <DashboardItemBoard
        :title="formattedDate"
        description="Add and manage items for the selected date."
        :date="selectedIso"
        :active-tag-id="activeTagId"
        mobile-priority
      />
    </div>
  </div>
</template>
