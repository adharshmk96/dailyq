<script setup lang="ts">
defineProps<{
  title: string
  description: string
  icon?: string
}>()

const colorMode = useColorMode()

const isDark = computed({
  get: () => colorMode.value === 'dark',
  set: (value: boolean) => {
    colorMode.preference = value ? 'dark' : 'light'
  }
})
</script>

<template>
  <div class="relative min-h-dvh overflow-hidden bg-default">
    <div
      class="pointer-events-none absolute inset-x-0 -top-40 h-[32rem] bg-[radial-gradient(60%_60%_at_50%_50%,var(--ui-primary)_0%,transparent_70%)] opacity-[0.09]"
      aria-hidden="true"
    />
    <div
      class="pointer-events-none absolute inset-0 [mask-image:radial-gradient(70%_50%_at_50%_0%,black,transparent)] bg-[linear-gradient(to_right,var(--ui-border)_1px,transparent_1px),linear-gradient(to_bottom,var(--ui-border)_1px,transparent_1px)] bg-[size:56px_56px] opacity-60"
      aria-hidden="true"
    />

    <header class="relative mx-auto flex h-16 max-w-6xl items-center px-4 sm:px-6">
      <NuxtLink
        to="/"
        class="flex shrink-0 items-center gap-2.5"
      >
        <span class="flex size-9 items-center justify-center rounded-xl bg-primary text-inverted shadow-sm">
          <UIcon
            name="i-lucide-notebook-pen"
            class="size-5"
          />
        </span>
        <span class="text-base font-semibold tracking-tight text-highlighted">DailyQ</span>
      </NuxtLink>

      <div class="ml-auto">
        <ClientOnly>
          <UButton
            :icon="isDark ? 'i-lucide-moon' : 'i-lucide-sun'"
            variant="ghost"
            color="neutral"
            aria-label="Toggle theme"
            @click="isDark = !isDark"
          />
          <template #fallback>
            <div class="size-8" />
          </template>
        </ClientOnly>
      </div>
    </header>

    <main class="relative mx-auto flex w-full max-w-md flex-col px-4 pb-20 pt-8 sm:px-6 sm:pt-14">
      <div class="text-center">
        <span
          v-if="icon"
          class="mx-auto flex size-11 items-center justify-center rounded-xl bg-primary/10 text-primary"
        >
          <UIcon
            :name="icon"
            class="size-5"
          />
        </span>
        <h1
          class="text-balance text-2xl font-semibold tracking-tight text-highlighted sm:text-3xl"
          :class="icon ? 'mt-4' : ''"
        >
          {{ title }}
        </h1>
        <p class="mx-auto mt-2 max-w-sm text-pretty text-sm leading-relaxed text-muted">
          {{ description }}
        </p>
      </div>

      <UCard
        class="mt-8"
        :ui="{ body: 'p-5 sm:p-6' }"
      >
        <slot />
      </UCard>

      <div
        v-if="$slots.footer"
        class="mt-6 text-center text-sm text-muted"
      >
        <slot name="footer" />
      </div>
    </main>
  </div>
</template>
