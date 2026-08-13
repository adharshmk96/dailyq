<script setup lang="ts">
const colorMode = useColorMode()

const isDark = computed({
  get: () => colorMode.value === 'dark',
  set: (value: boolean) => {
    colorMode.preference = value ? 'dark' : 'light'
  }
})

const features = [
  {
    icon: 'i-lucide-feather',
    title: 'Write, then get on with it',
    description: 'A single quiet page. No projects, no boards to maintain, no setup ritual before your first line.'
  },
  {
    icon: 'i-lucide-list-checks',
    title: 'Tasks and notes together',
    description: 'A thought and a to-do live in the same place, because that is how they arrive.'
  },
  {
    icon: 'i-lucide-tags',
    title: 'Tags, not folders',
    description: 'Colour-coded tags keep things findable without forcing you to file anything away.'
  },
  {
    icon: 'i-lucide-calendar-days',
    title: 'A day at a time',
    description: 'Each day gets its own page. Yesterday stays where you left it, ready when you look back.'
  },
  {
    icon: 'i-lucide-keyboard',
    title: 'Keyboard first',
    description: 'A rich editor with slash commands and shortcuts. Your hands never have to leave home row.'
  },
  {
    icon: 'i-lucide-moon-star',
    title: 'Easy on the eyes',
    description: 'Considered light and dark themes for early mornings and late, quieter hours.'
  }
]

const steps = [
  {
    number: '01',
    title: 'Open the day',
    description: 'DailyQ starts on today. Empty, patient, waiting for whatever is on your mind.'
  },
  {
    number: '02',
    title: 'Empty your head',
    description: 'Type tasks and notes as they come. Tag them if it helps. Do not organise, just write.'
  },
  {
    number: '03',
    title: 'Close the loop',
    description: 'Tick things off through the day. What is left rolls forward. Nothing is lost.'
  }
]

const previewTasks = [
  { label: 'Review the Q3 handover doc', done: true, tag: 'work', color: 'primary' as const },
  { label: 'Call the dentist back', done: true, tag: 'errand', color: 'warning' as const },
  { label: 'Sketch onboarding flow', done: false, tag: 'work', color: 'primary' as const },
  { label: 'Water the fig tree', done: false, tag: 'home', color: 'success' as const }
]

useHead({
  title: 'DailyQ — a quiet place for today\'s tasks and notes'
})

useSeoMeta({
  description: 'DailyQ is a minimal daily journal for tasks and notes. One page a day, tags instead of folders, and nothing else in the way.'
})
</script>

<template>
  <div class="min-h-dvh bg-default">
    <!-- Nav -->
    <header class="sticky top-0 z-40 border-b border-default/60 bg-default/80 backdrop-blur-md">
      <div class="mx-auto flex h-16 max-w-6xl items-center gap-3 px-4 sm:px-6">
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

        <nav class="ml-6 hidden items-center gap-1 md:flex">
          <UButton
            to="#features"
            variant="ghost"
            color="neutral"
            label="Features"
          />
          <UButton
            to="#how"
            variant="ghost"
            color="neutral"
            label="How it works"
          />
        </nav>

        <div class="ml-auto flex items-center gap-2">
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
          <UButton
            to="/login"
            variant="ghost"
            color="neutral"
            label="Sign in"
          />
          <UButton
            to="/register"
            label="Get started"
            trailing-icon="i-lucide-arrow-right"
          />
        </div>
      </div>
    </header>

    <main>
      <!-- Hero -->
      <section class="relative overflow-hidden">
        <div
          class="pointer-events-none absolute inset-x-0 -top-40 h-[32rem] bg-[radial-gradient(60%_60%_at_50%_50%,var(--ui-primary)_0%,transparent_70%)] opacity-[0.09]"
          aria-hidden="true"
        />
        <div
          class="pointer-events-none absolute inset-0 [mask-image:radial-gradient(70%_50%_at_50%_0%,black,transparent)] bg-[linear-gradient(to_right,var(--ui-border)_1px,transparent_1px),linear-gradient(to_bottom,var(--ui-border)_1px,transparent_1px)] bg-[size:56px_56px] opacity-60"
          aria-hidden="true"
        />

        <div class="relative mx-auto max-w-6xl px-4 pb-16 pt-20 sm:px-6 sm:pb-24 sm:pt-28">
          <div class="mx-auto max-w-2xl text-center">
            <UBadge
              variant="subtle"
              color="primary"
              class="rounded-full px-3 py-1"
            >
              <UIcon
                name="i-lucide-sparkles"
                class="mr-1.5 size-3.5"
              />
              One page a day. That's the whole idea.
            </UBadge>

            <h1 class="mt-6 text-balance text-4xl font-semibold tracking-tight text-highlighted sm:text-6xl">
              A quiet place for
              <span class="text-primary">today's</span>
              tasks and notes.
            </h1>

            <p class="mx-auto mt-5 max-w-xl text-pretty text-base leading-relaxed text-muted sm:text-lg">
              DailyQ gives each day a single blank page. Write what you need to do,
              write what you're thinking, and let the day carry the rest.
            </p>

            <div class="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
              <UButton
                to="/register"
                size="lg"
                label="Start writing"
                trailing-icon="i-lucide-arrow-right"
              />
              <UButton
                to="#how"
                size="lg"
                variant="subtle"
                color="neutral"
                label="See how it works"
              />
            </div>

            <p class="mt-5 text-xs text-dimmed">
              Free to use · No setup · Your day, not your data
            </p>
          </div>

          <!-- App preview -->
          <div class="relative mx-auto mt-16 max-w-3xl">
            <div class="rounded-2xl border border-default bg-default shadow-2xl shadow-black/5 ring-1 ring-black/[0.03] dark:shadow-black/40">
              <div class="flex items-center gap-2 border-b border-default px-4 py-3">
                <span class="size-2.5 rounded-full bg-elevated" />
                <span class="size-2.5 rounded-full bg-elevated" />
                <span class="size-2.5 rounded-full bg-elevated" />
                <span class="ml-3 text-xs text-dimmed">Today · Thursday</span>
              </div>

              <div class="space-y-5 p-5 sm:p-7">
                <div>
                  <h2 class="text-lg font-semibold tracking-tight text-highlighted">
                    Today
                  </h2>
                  <p class="mt-0.5 text-sm text-muted">
                    Two down, two to go.
                  </p>
                </div>

                <ul class="space-y-1">
                  <li
                    v-for="task in previewTasks"
                    :key="task.label"
                    class="flex items-center gap-3 rounded-lg px-2 py-2 transition-colors hover:bg-elevated/60"
                  >
                    <UIcon
                      :name="task.done ? 'i-lucide-circle-check' : 'i-lucide-circle'"
                      class="size-4 shrink-0"
                      :class="task.done ? 'text-primary' : 'text-dimmed'"
                    />
                    <span
                      class="text-sm"
                      :class="task.done ? 'text-dimmed line-through' : 'text-default'"
                    >
                      {{ task.label }}
                    </span>
                    <UBadge
                      :color="task.color"
                      variant="subtle"
                      size="sm"
                      class="ml-auto rounded-full"
                      :label="task.tag"
                    />
                  </li>
                </ul>

                <div class="rounded-xl border border-dashed border-default bg-muted/40 p-4">
                  <p class="text-xs font-medium uppercase tracking-wide text-dimmed">
                    Note
                  </p>
                  <p class="mt-2 text-sm leading-relaxed text-muted">
                    The onboarding flow feels heavy — three screens before anyone
                    writes a word. Try collapsing it into one and see what breaks.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Features -->
      <section
        id="features"
        class="border-t border-default/60 bg-muted/40 py-20 sm:py-28"
      >
        <div class="mx-auto max-w-6xl px-4 sm:px-6">
          <div class="max-w-2xl">
            <p class="text-sm font-medium text-primary">
              Deliberately small
            </p>
            <h2 class="mt-2 text-balance text-3xl font-semibold tracking-tight text-highlighted sm:text-4xl">
              Everything you need, and pointedly nothing more.
            </h2>
            <p class="mt-4 text-pretty text-muted">
              Most tools ask you to build a system before you can use them.
              DailyQ asks you to start typing.
            </p>
          </div>

          <div class="mt-12 grid gap-px overflow-hidden rounded-2xl border border-default bg-border sm:grid-cols-2 lg:grid-cols-3">
            <div
              v-for="feature in features"
              :key="feature.title"
              class="group bg-default p-6 transition-colors hover:bg-elevated/40 sm:p-7"
            >
              <span class="flex size-10 items-center justify-center rounded-xl bg-primary/10 text-primary transition-transform group-hover:scale-105">
                <UIcon
                  :name="feature.icon"
                  class="size-5"
                />
              </span>
              <h3 class="mt-4 text-base font-semibold text-highlighted">
                {{ feature.title }}
              </h3>
              <p class="mt-2 text-sm leading-relaxed text-muted">
                {{ feature.description }}
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- How it works -->
      <section
        id="how"
        class="py-20 sm:py-28"
      >
        <div class="mx-auto max-w-6xl px-4 sm:px-6">
          <div class="mx-auto max-w-2xl text-center">
            <h2 class="text-balance text-3xl font-semibold tracking-tight text-highlighted sm:text-4xl">
              Three steps, and two of them are just writing.
            </h2>
          </div>

          <ol class="mt-14 grid gap-8 md:grid-cols-3">
            <li
              v-for="step in steps"
              :key="step.number"
              class="relative"
            >
              <span class="text-sm font-semibold tabular-nums text-primary/70">
                {{ step.number }}
              </span>
              <div class="mt-3 h-px w-full bg-border" />
              <h3 class="mt-4 text-lg font-semibold text-highlighted">
                {{ step.title }}
              </h3>
              <p class="mt-2 text-sm leading-relaxed text-muted">
                {{ step.description }}
              </p>
            </li>
          </ol>
        </div>
      </section>

      <!-- Quote -->
      <section class="border-y border-default/60 bg-muted/40 py-20 sm:py-24">
        <div class="mx-auto max-w-3xl px-4 text-center sm:px-6">
          <UIcon
            name="i-lucide-quote"
            class="size-8 text-primary/30"
          />
          <blockquote class="mt-6 text-balance text-xl font-medium leading-relaxed text-highlighted sm:text-2xl">
            The point isn't to capture everything. It's to put the day down
            somewhere so your head can stop holding it.
          </blockquote>
          <p class="mt-6 text-sm text-dimmed">
            The idea behind DailyQ
          </p>
        </div>
      </section>

      <!-- CTA -->
      <section class="py-20 sm:py-28">
        <div class="mx-auto max-w-6xl px-4 sm:px-6">
          <div class="relative overflow-hidden rounded-3xl border border-default bg-default px-6 py-14 text-center shadow-sm sm:px-12 sm:py-20">
            <div
              class="pointer-events-none absolute inset-0 bg-[radial-gradient(50%_80%_at_50%_100%,var(--ui-primary)_0%,transparent_70%)] opacity-[0.08]"
              aria-hidden="true"
            />
            <div class="relative">
              <h2 class="text-balance text-3xl font-semibold tracking-tight text-highlighted sm:text-4xl">
                Today's page is still blank.
              </h2>
              <p class="mx-auto mt-4 max-w-lg text-pretty text-muted">
                Give it a minute of your morning and see what the rest of the day feels like.
              </p>
              <UButton
                to="/register"
                size="lg"
                class="mt-8"
                label="Open DailyQ"
                trailing-icon="i-lucide-arrow-right"
              />
            </div>
          </div>
        </div>
      </section>
    </main>

    <footer class="border-t border-default/60 py-10">
      <div class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-4 sm:flex-row sm:px-6">
        <div class="flex items-center gap-2.5">
          <span class="flex size-7 items-center justify-center rounded-lg bg-primary text-inverted">
            <UIcon
              name="i-lucide-notebook-pen"
              class="size-4"
            />
          </span>
          <span class="text-sm font-medium text-highlighted">DailyQ</span>
        </div>
        <p class="text-sm text-dimmed">
          One page a day.
        </p>
      </div>
    </footer>
  </div>
</template>
