<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'
import type { JournalTag, TagColor } from '~/types/journal'

const props = withDefaults(defineProps<{
  /** Hide outer card chrome when nested inside CollapsibleSection. */
  embedded?: boolean
}>(), {
  embedded: false
})

const { tags, addTag, updateTag, deleteTag } = useJournalTags()
const { activeTagId, toggleTag } = useTagFilter()

const modalOpen = ref(false)
const editingTag = ref<JournalTag | null>(null)

const colorDotClass: Record<TagColor, string> = {
  primary: 'bg-primary',
  success: 'bg-success',
  warning: 'bg-warning',
  info: 'bg-info',
  error: 'bg-error',
  neutral: 'bg-neutral'
}

const activeTagName = computed(() => {
  if (!activeTagId.value) return null
  return tags.value.find(tag => tag.id === activeTagId.value)?.name ?? null
})

function openCreate() {
  editingTag.value = null
  modalOpen.value = true
}

function openEdit(tag: JournalTag) {
  editingTag.value = tag
  modalOpen.value = true
}

function onSave(payload: { name: string, color: TagColor }) {
  if (editingTag.value) {
    updateTag(editingTag.value.id, payload)
    return
  }
  addTag(payload.name, payload.color)
}

function onDelete(tag: JournalTag) {
  if (activeTagId.value === tag.id) {
    toggleTag(tag.id)
  }
  deleteTag(tag.id)
}

function menuItems(tag: JournalTag): DropdownMenuItem[][] {
  return [
    [
      {
        label: 'Edit',
        icon: 'i-lucide-pencil',
        onSelect() {
          openEdit(tag)
        }
      },
      {
        label: 'Delete',
        icon: 'i-lucide-trash-2',
        color: 'error',
        onSelect() {
          onDelete(tag)
        }
      }
    ]
  ]
}
</script>

<template>
  <UCard
    :class="embedded ? 'border-0 shadow-none lg:border lg:shadow-sm' : undefined"
    :ui="{ body: 'space-y-3 p-3 sm:p-4 lg:p-5' }"
  >
    <div class="flex items-center justify-between gap-2">
      <div class="min-w-0">
        <h2 class="text-sm font-semibold text-highlighted">
          Tag
        </h2>
        <p
          v-if="activeTagName"
          class="truncate text-xs text-primary"
        >
          Filtering: {{ activeTagName }}
        </p>
      </div>
      <UButton
        icon="i-lucide-plus"
        color="neutral"
        variant="ghost"
        size="xs"
        aria-label="Add tag"
        @click="openCreate"
      />
    </div>

    <ul
      v-if="tags.length"
      class="space-y-1"
    >
      <li
        v-for="tag in tags"
        :key="tag.id"
        class="group flex items-center gap-1 rounded-lg"
      >
        <button
          type="button"
          class="flex min-w-0 flex-1 items-center gap-2 rounded-lg px-2 py-1.5 text-left text-sm transition"
          :class="activeTagId === tag.id
            ? 'bg-primary/10 text-highlighted ring-1 ring-primary/30'
            : 'text-default hover:bg-elevated'"
          @click="toggleTag(tag.id)"
        >
          <span
            class="size-2.5 shrink-0 rounded-full"
            :class="colorDotClass[tag.color]"
          />
          <span class="truncate">{{ tag.name }}</span>
        </button>

        <UDropdownMenu
          :items="menuItems(tag)"
          :content="{ align: 'end' }"
          :ui="{ content: 'w-36' }"
        >
          <UButton
            icon="i-lucide-ellipsis"
            color="neutral"
            variant="ghost"
            size="xs"
            class="opacity-100 sm:opacity-0 sm:group-hover:opacity-100"
            aria-label="Tag actions"
          />
        </UDropdownMenu>
      </li>
    </ul>

    <p
      v-else
      class="text-xs text-muted"
    >
      No tags yet. Click + to create one.
    </p>

    <DashboardTagModal
      v-model:open="modalOpen"
      :editing-id="editingTag?.id"
      :initial-name="editingTag?.name"
      :initial-color="editingTag?.color"
      @save="onSave"
    />
  </UCard>
</template>
