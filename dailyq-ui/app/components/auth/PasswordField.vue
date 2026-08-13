<script setup lang="ts">
defineProps<{
  label: string
  placeholder?: string
  autocomplete?: string
  hint?: string
}>()

const model = defineModel<string>({ required: true })

const show = ref(false)
</script>

<template>
  <UFormField
    :label="label"
    :hint="hint"
    class="w-full"
  >
    <UInput
      v-model="model"
      :type="show ? 'text' : 'password'"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      class="w-full"
    >
      <template #trailing>
        <UButton
          color="neutral"
          variant="link"
          :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'"
          :aria-label="show ? 'Hide password' : 'Show password'"
          @click="show = !show"
        />
      </template>
    </UInput>
  </UFormField>
</template>
