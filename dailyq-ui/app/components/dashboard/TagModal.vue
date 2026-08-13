<script setup lang="ts">
import { TAG_COLORS, type TagColor } from '~/types/journal'

const props = defineProps<{
  open: boolean
  editingId?: string | null
  initialName?: string
  initialColor?: TagColor
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  save: [payload: { name: string, color: TagColor }]
}>()

const name = ref('')
const color = ref<TagColor>('primary')

const isEditing = computed(() => !!props.editingId)
const title = computed(() => (isEditing.value ? 'Edit tag' : 'New tag'))

watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen) return
    name.value = props.initialName ?? ''
    color.value = props.initialColor ?? 'primary'
  }
)

function close() {
  emit('update:open', false)
}

function save() {
  const trimmed = name.value.trim()
  if (!trimmed) return
  emit('save', { name: trimmed, color: color.value })
  close()
}

const colorClassMap: Record<TagColor, string> = {
  primary: 'bg-primary',
  success: 'bg-success',
  warning: 'bg-warning',
  info: 'bg-info',
  error: 'bg-error',
  neutral: 'bg-neutral'
}
</script>

<template>
  <UModal
    :open="open"
    :title="title"
    @update:open="emit('update:open', $event)"
  >
    <template #body>
      <div class="space-y-4">
        <UFormField label="Tag name">
          <UInput
            v-model="name"
            placeholder="e.g. Work, Personal…"
            autofocus
            size="lg"
            class="w-full"
            @keydown.enter.prevent="save"
          />
        </UFormField>

        <UFormField label="Color">
          <div class="flex flex-wrap gap-2">
            <button
              v-for="c in TAG_COLORS"
              :key="c"
              type="button"
              class="flex size-9 items-center justify-center rounded-full transition ring-offset-2 ring-offset-default focus:outline-none focus-visible:ring-2 focus-visible:ring-primary"
              :class="[
                colorClassMap[c],
                color === c ? 'ring-2 ring-highlighted scale-110' : 'opacity-70 hover:opacity-100'
              ]"
              :aria-label="`Select ${c} color`"
              :aria-pressed="color === c"
              @click="color = c"
            >
              <UIcon
                v-if="color === c"
                name="i-lucide-check"
                class="size-4 text-white"
              />
            </button>
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
          @click="close"
        />
        <UButton
          label="Save"
          icon="i-lucide-check"
          :disabled="!name.trim()"
          @click="save"
        />
      </div>
    </template>
  </UModal>
</template>
