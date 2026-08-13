<script setup lang="ts">
import type { TagColor } from '~/types/journal'

const {
  totalTasksCount,
  completedTodayCount,
  openGeneralCount,
  notesThisWeekCount,
  todayTasks,
  recentNotes,
  toggleTaskDone,
  getTagById
} = usePlaceholderJournal()

const badgeColorMap: Record<TagColor, 'primary' | 'success' | 'warning' | 'info' | 'error' | 'neutral'> = {
  primary: 'primary',
  success: 'success',
  warning: 'warning',
  info: 'info',
  error: 'error',
  neutral: 'neutral'
}

function resolveTagBadges(tagIds: string[]) {
  return tagIds
    .map(id => getTagById(id))
    .filter((tag): tag is NonNullable<typeof tag> => !!tag)
}

const stats = computed(() => [
  {
    label: 'Total tasks',
    value: totalTasksCount.value,
    icon: 'i-lucide-list-checks',
    color: 'text-primary'
  },
  {
    label: 'Done today',
    value: completedTodayCount.value,
    icon: 'i-lucide-circle-check',
    color: 'text-success'
  },
  {
    label: 'Open general',
    value: openGeneralCount.value,
    icon: 'i-lucide-layers',
    color: 'text-warning'
  },
  {
    label: 'Notes this week',
    value: notesThisWeekCount.value,
    icon: 'i-lucide-notebook-pen',
    color: 'text-info'
  }
])
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight text-highlighted">
        Overview
      </h1>
      <p class="mt-1 text-sm text-muted">
        A quick look at your notes and tasks.
      </p>
    </div>

    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <UCard
        v-for="stat in stats"
        :key="stat.label"
        :ui="{ body: 'p-4 sm:p-5' }"
      >
        <div class="flex items-start justify-between gap-3">
          <div>
            <p class="text-sm text-muted">
              {{ stat.label }}
            </p>
            <p class="mt-1 text-3xl font-semibold tabular-nums text-highlighted">
              {{ stat.value }}
            </p>
          </div>
          <span
            class="flex size-10 items-center justify-center rounded-xl bg-elevated"
            :class="stat.color"
          >
            <UIcon
              :name="stat.icon"
              class="size-5"
            />
          </span>
        </div>
      </UCard>
    </div>

    <div class="grid gap-4 lg:grid-cols-2">
      <UCard>
        <template #header>
          <div class="flex items-center gap-2">
            <UIcon
              name="i-lucide-sun"
              class="size-4 text-primary"
            />
            <h2 class="font-medium text-highlighted">
              Today's tasks
            </h2>
            <UBadge
              color="neutral"
              variant="subtle"
              size="sm"
            >
              {{ todayTasks.length }}
            </UBadge>
          </div>
        </template>

        <ul
          v-if="todayTasks.length"
          class="divide-y divide-default"
        >
          <li
            v-for="task in todayTasks"
            :key="task.id"
            class="flex items-center gap-3 py-3 first:pt-0 last:pb-0"
          >
            <UCheckbox
              :model-value="task.done"
              :aria-label="`Toggle ${task.title}`"
              @update:model-value="toggleTaskDone(task.id)"
            />
            <div class="min-w-0 flex-1 space-y-1">
              <span
                class="block text-sm"
                :class="task.done ? 'text-muted line-through' : 'text-default'"
              >
                {{ task.title }}
              </span>
              <div
                v-if="task.tagIds.length"
                class="flex flex-wrap gap-1"
              >
                <UBadge
                  v-for="tag in resolveTagBadges(task.tagIds)"
                  :key="tag.id"
                  :color="badgeColorMap[tag.color]"
                  variant="subtle"
                  size="xs"
                >
                  {{ tag.name }}
                </UBadge>
              </div>
            </div>
          </li>
        </ul>
        <UEmpty
          v-else
          icon="i-lucide-circle-check"
          title="All clear"
          description="No tasks scheduled for today."
          :ui="{ root: 'py-6' }"
        />
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center gap-2">
            <UIcon
              name="i-lucide-sticky-note"
              class="size-4 text-primary"
            />
            <h2 class="font-medium text-highlighted">
              Recent notes
            </h2>
          </div>
        </template>

        <ul
          v-if="recentNotes.length"
          class="divide-y divide-default"
        >
          <li
            v-for="note in recentNotes"
            :key="note.id"
            class="flex items-start gap-3 py-3 first:pt-0 last:pb-0"
          >
            <UIcon
              name="i-lucide-file-text"
              class="mt-0.5 size-4 shrink-0 text-muted"
            />
            <div class="min-w-0 flex-1">
              <p class="text-sm text-default">
                {{ note.body }}
              </p>
              <p class="mt-1 text-xs text-muted">
                {{ note.date ?? 'General' }}
              </p>
            </div>
          </li>
        </ul>
        <UEmpty
          v-else
          icon="i-lucide-notebook"
          title="No notes yet"
          description="Capture a thought to see it here."
          :ui="{ root: 'py-6' }"
        />
      </UCard>
    </div>
  </div>
</template>
