<script setup lang="ts">
const { pending, submit, fail } = usePlaceholderAuth()

const email = ref('')
const password = ref('')
const remember = ref(true)

useHead({
  title: 'Sign in — DailyQ'
})

async function onSubmit() {
  if (!email.value || !password.value) {
    fail('Missing fields', 'Enter your email and password to continue.')
    return
  }

  if (!EMAIL_PATTERN.test(email.value)) {
    fail('Invalid email', 'Enter a valid email address.')
    return
  }

  await submit({
    title: 'Signed in',
    description: 'Placeholder only — no account was checked.',
    redirectTo: '/dashboard/overview'
  })
}
</script>

<template>
  <AuthShell
    icon="i-lucide-log-in"
    title="Welcome back"
    description="Sign in to pick up today's page where you left it."
  >
    <form
      class="space-y-5"
      @submit.prevent="onSubmit"
    >
      <UFormField
        label="Email"
        class="w-full"
      >
        <UInput
          v-model="email"
          type="email"
          placeholder="you@example.com"
          autocomplete="email"
          class="w-full"
        />
      </UFormField>

      <AuthPasswordField
        v-model="password"
        label="Password"
        placeholder="••••••••"
        autocomplete="current-password"
      />

      <div class="flex items-center justify-between gap-3">
        <UCheckbox
          v-model="remember"
          label="Remember me"
        />
        <ULink
          to="/forgot-password"
          class="text-sm font-medium text-primary"
        >
          Forgot password?
        </ULink>
      </div>

      <UButton
        type="submit"
        block
        size="lg"
        label="Sign in"
        :loading="pending"
      />
    </form>

    <template #footer>
      New here?
      <ULink
        to="/register"
        class="font-medium text-primary"
      >
        Create an account
      </ULink>
    </template>
  </AuthShell>
</template>
