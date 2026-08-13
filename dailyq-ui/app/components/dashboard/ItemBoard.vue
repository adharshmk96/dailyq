<script setup lang="ts">
import type { TabsItem } from '@nuxt/ui'
import type { ItemKind, JournalNote, JournalTask, TagColor } from '~/types/journal'

const props = withDefaults(defineProps<{
  title: string
  description?: string
  /** ISO date for the calendar board; null for the undated General board. */
  date?: string | null
  activeTagId?: string | null
}>(), {
  description: undefined,
  date: null,
  activeTagId: null
})

const boardDate = computed(() => props.date ?? null)

const { tags, getTagById } = useJournalTags()

const {
  tasks: allTasks,
  notes: allNotes,
  pending,
  addTask,
  addNote,
  updateTask,
  updateNote,
  deleteTask,
  deleteNote,
  toggleTaskDone
} = useJournalBoard(boardDate)

const tasks = computed(() => filterItemsByTag(allTasks.value, props.activeTagId))
const notes = computed(() => filterItemsByTag(allNotes.value, props.activeTagId))

const kind = ref<ItemKind>('tasks')
const draft = ref('')
const editOpen = ref(false)
const editId = ref<string | null>(null)
const editText = ref('')
const editTagIds = ref<string[]>([])

const kindItems = computed<TabsItem[]>(() => [
  {
    label: 'Tasks',
    icon: 'i-lucide-list-todo',
    value: 'tasks'
  },
  {
    label: 'Notes',
    icon: 'i-lucide-sticky-note',
    value: 'notes'
  }
])

const isTasks = computed(() => kind.value === 'tasks')

const placeholder = computed(() =>
  isTasks.value ? 'Add a task…' : 'Write a note…'
)

const addLabel = computed(() => (isTasks.value ? 'Add task' : 'Add note'))

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

function addTagIdsForCreate() {
  return props.activeTagId ? [props.activeTagId] : []
}

const saving = ref(false)

async function onAdd() {
  if (saving.value || !draft.value.trim()) return

  const tagIds = addTagIdsForCreate()
  saving.value = true
  try {
    const created = isTasks.value
      ? await addTask(draft.value, tagIds)
      : await addNote(draft.value, tagIds)
    if (created) {
      draft.value = ''
    }
  } finally {
    saving.value = false
  }
}

function openEditTask(task: JournalTask) {
  kind.value = 'tasks'
  editId.value = task.id
  editText.value = task.title
  editTagIds.value = [...task.tagIds]
  editOpen.value = true
}

function openEditNote(note: JournalNote) {
  kind.value = 'notes'
  editId.value = note.id
  editText.value = note.body
  editTagIds.value = [...note.tagIds]
  editOpen.value = true
}

function toggleEditTag(tagId: string) {
  if (editTagIds.value.includes(tagId)) {
    editTagIds.value = editTagIds.value.filter(id => id !== tagId)
    return
  }
  editTagIds.value = [...editTagIds.value, tagId]
}

async function saveEdit() {
  if (!editId.value || saving.value) return

  saving.value = true
  try {
    const ok = isTasks.value
      ? await updateTask(editId.value, editText.value, editTagIds.value)
      : await updateNote(editId.value, editText.value, editTagIds.value)
    if (ok) {
      editOpen.value = false
      editId.value = null
      editText.value = ''
      editTagIds.value = []
    }
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="space-y-5">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight text-highlighted">
          {{ title }}
        </h1>
        <p
          v-if="description"
          class="mt-1 text-sm text-muted"
        >
          {{ description }}
        </p>
      </div>

      <UTabs
        v-model="kind"
        :items="kindItems"
        :content="false"
        color="neutral"
        variant="pill"
        size="sm"
        class="w-full sm:w-auto"
      />
    </div>

    <UCard :ui="{ body: 'space-y-4 p-4 sm:p-5' }">
      <form
        class="flex gap-2"
        @submit.prevent="onAdd"
      >
        <UInput
          v-model="draft"
          :placeholder="placeholder"
          icon="i-lucide-plus"
          class="flex-1"
          size="lg"
        />
        <UButton
          type="submit"
          :label="addLabel"
          icon="i-lucide-plus"
          size="lg"
          :loading="saving"
          :disabled="!draft.trim()"
        />
      </form>

      <USeparator />

      <div
        v-if="pending"
        class="space-y-3 py-2"
      >
        <USkeleton
          v-for="i in 3"
          :key="i"
          class="h-6 w-full"
        />
      </div>

      <template v-else-if="isTasks">
        <ul
          v-if="tasks.length"
          class="divide-y divide-default"
        >
          <li
            v-for="task in tasks"
            :key="task.id"
            class="group flex items-center gap-3 py-3 first:pt-0 last:pb-0"
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
            <div class="flex shrink-0 opacity-100 transition sm:opacity-0 sm:group-hover:opacity-100">
              <UButton
                icon="i-lucide-pencil"
                color="neutral"
                variant="ghost"
                size="xs"
                aria-label="Edit task"
                @click="openEditTask(task)"
              />
              <UButton
                icon="i-lucide-trash-2"
                color="error"
                variant="ghost"
                size="xs"
                aria-label="Delete task"
                @click="deleteTask(task.id)"
              />
            </div>
          </li>
        </ul>
        <UEmpty
          v-else
          icon="i-lucide-list-todo"
          title="No tasks"
          description="Add your first task above."
          :ui="{ root: 'py-8' }"
        />
      </template>

      <template v-else>
        <ul
          v-if="notes.length"
          class="divide-y divide-default"
        >
          <li
            v-for="note in notes"
            :key="note.id"
            class="group flex items-start gap-3 py-3 first:pt-0 last:pb-0"
          >
            <UIcon
              name="i-lucide-file-text"
              class="mt-1 size-4 shrink-0 text-muted"
            />
            <div class="min-w-0 flex-1 space-y-1">
              <p class="text-sm text-default">
                {{ note.body }}
              </p>
              <div
                v-if="note.tagIds.length"
                class="flex flex-wrap gap-1"
              >
                <UBadge
                  v-for="tag in resolveTagBadges(note.tagIds)"
                  :key="tag.id"
                  :color="badgeColorMap[tag.color]"
                  variant="subtle"
                  size="xs"
                >
                  {{ tag.name }}
                </UBadge>
              </div>
            </div>
            <div class="flex shrink-0 opacity-100 transition sm:opacity-0 sm:group-hover:opacity-100">
              <UButton
                icon="i-lucide-pencil"
                color="neutral"
                variant="ghost"
                size="xs"
                aria-label="Edit note"
                @click="openEditNote(note)"
              />
              <UButton
                icon="i-lucide-trash-2"
                color="error"
                variant="ghost"
                size="xs"
                aria-label="Delete note"
                @click="deleteNote(note.id)"
              />
            </div>
          </li>
        </ul>
        <UEmpty
          v-else
          icon="i-lucide-sticky-note"
          title="No notes"
          description="Capture a note above."
          :ui="{ root: 'py-8' }"
        />
      </template>
    </UCard>

    <UModal
      v-model:open="editOpen"
      :title="isTasks ? 'Edit task' : 'Edit note'"
    >
      <template #body>
        <div class="space-y-4">
          <UInput
            v-model="editText"
            :placeholder="placeholder"
            autofocus
            size="lg"
            class="w-full"
            @keydown.enter.prevent="saveEdit"
          />

          <UFormField
            v-if="tags.length"
            label="Tags"
          >
            <div class="flex flex-wrap gap-2">
              <label
                v-for="tag in tags"
                :key="tag.id"
                class="inline-flex cursor-pointer items-center gap-2 rounded-lg border border-default px-2.5 py-1.5 text-sm transition"
                :class="editTagIds.includes(tag.id) ? 'bg-elevated' : 'hover:bg-elevated/50'"
              >
                <UCheckbox
                  :model-value="editTagIds.includes(tag.id)"
                  :aria-label="`Toggle ${tag.name} tag`"
                  @update:model-value="toggleEditTag(tag.id)"
                />
                <span>{{ tag.name }}</span>
              </label>
            </div>
          </UFormField>
        </div>
      </template>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UButton
            label="Cancel"
            color="neutral"
            variant="ghost"
            @click="editOpen = false"
          />
          <UButton
            label="Save"
            icon="i-lucide-check"
            :loading="saving"
            @click="saveEdit"
          />
        </div>
      </template>
    </UModal>
  </div>
</template>
