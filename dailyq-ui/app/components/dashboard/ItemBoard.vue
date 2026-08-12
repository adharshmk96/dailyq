<script setup lang="ts">
import type { TabsItem } from '@nuxt/ui'
import type { ItemKind, JournalNote, JournalTask } from '~/types/journal'

const props = withDefaults(defineProps<{
  title: string
  description?: string
  tasks: JournalTask[]
  notes: JournalNote[]
  date?: string | null
}>(), {
  description: undefined,
  date: null
})

const {
  addTask,
  addNote,
  updateTask,
  updateNote,
  deleteTask,
  deleteNote,
  toggleTaskDone
} = usePlaceholderJournal()

const kind = ref<ItemKind>('tasks')
const draft = ref('')
const editOpen = ref(false)
const editId = ref<string | null>(null)
const editText = ref('')

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

function onAdd() {
  if (isTasks.value) {
    if (addTask(draft.value, props.date)) {
      draft.value = ''
    }
    return
  }
  if (addNote(draft.value, props.date)) {
    draft.value = ''
  }
}

function openEditTask(task: JournalTask) {
  kind.value = 'tasks'
  editId.value = task.id
  editText.value = task.title
  editOpen.value = true
}

function openEditNote(note: JournalNote) {
  kind.value = 'notes'
  editId.value = note.id
  editText.value = note.body
  editOpen.value = true
}

function saveEdit() {
  if (!editId.value) return
  const ok = isTasks.value
    ? updateTask(editId.value, editText.value)
    : updateNote(editId.value, editText.value)
  if (ok) {
    editOpen.value = false
    editId.value = null
    editText.value = ''
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
        />
      </form>

      <USeparator />

      <template v-if="isTasks">
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
            <span
              class="min-w-0 flex-1 text-sm"
              :class="task.done ? 'text-muted line-through' : 'text-default'"
            >
              {{ task.title }}
            </span>
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
            <p class="min-w-0 flex-1 text-sm text-default">
              {{ note.body }}
            </p>
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
        <UInput
          v-model="editText"
          :placeholder="placeholder"
          autofocus
          size="lg"
          class="w-full"
          @keydown.enter.prevent="saveEdit"
        />
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
            @click="saveEdit"
          />
        </div>
      </template>
    </UModal>
  </div>
</template>
