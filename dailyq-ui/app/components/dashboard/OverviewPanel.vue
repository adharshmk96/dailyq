<script setup lang="ts">
import type { TagColor } from '~/types/journal'

const { overview, pending, error, refresh, getTagById, toggleTaskDone } = useOverview()

const toast = useToast()

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

async function onToggle(id: string) {
  try {
    await toggleTaskDone(id)
  } catch (err) {
    toast.add({
      title: apiErrorMessage(err, 'Could not update the task.'),
      color: 'error'
    })
  }
}

const todayTasks = computed(() => overview.value?.today_tasks ?? [])
const recentNotes = computed(() => overview.value?.recent_notes ?? [])

const stats = computed(() => {
  const s = overview.value?.stats
  return [
    {
      label: 'Total tasks',
      value: s?.total_tasks ?? 0,
      icon: 'i-lucide-list-checks',
      color: 'text-primary'
    },
    {
      label: 'Done today',
      value: s?.completed_today ?? 0,
      icon: 'i-lucide-circle-check',
      color: 'text-success'
    },
    {
      label: 'Open general',
      value: s?.open_general ?? 0,
      icon: 'i-lucide-layers',
      color: 'text-warning'
    },
    {
      label: 'Notes this week',
      value: s?.notes_this_week ?? 0,
      icon: 'i-lucide-notebook-pen',
      color: 'text-info'
    }
  ]
})

// Only the very first load shows skeletons; refreshes keep the current data.
const loading = computed(() => pending.value && !overview.value)
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

    <UAlert
      v-if="error"
      icon="i-lucide-triangle-alert"
      color="error"
      variant="subtle"
      title="Couldn't load your overview"
      :description="apiErrorMessage(error)"
      :actions="[{ label: 'Retry', color: 'error', variant: 'outline', onClick: () => refresh() }]"
    />

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
            <USkeleton
              v-if="loading"
              class="mt-2 h-8 w-12"
            />
            <p
              v-else
              class="mt-1 text-3xl font-semibold tabular-nums text-highlighted"
            >
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
              v-if="!loading"
              color="neutral"
              variant="subtle"
              size="sm"
            >
              {{ todayTasks.length }}
            </UBadge>
          </div>
        </template>

        <div
          v-if="loading"
          class="space-y-3"
        >
          <USkeleton
            v-for="n in 3"
            :key="n"
            class="h-6 w-full"
          />
        </div>
        <ul
          v-else-if="todayTasks.length"
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
              @update:model-value="onToggle(task.id)"
            />
            <div class="min-w-0 flex-1 space-y-1">
              <span
                class="block text-sm"
                :class="task.done ? 'text-muted line-through' : 'text-default'"
              >
                {{ task.title }}
              </span>
              <div
                v-if="task.tag_ids.length"
                class="flex flex-wrap gap-1"
              >
                <UBadge
                  v-for="tag in resolveTagBadges(task.tag_ids)"
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

        <div
          v-if="loading"
          class="space-y-3"
        >
          <USkeleton
            v-for="n in 3"
            :key="n"
            class="h-6 w-full"
          />
        </div>
        <ul
          v-else-if="recentNotes.length"
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
