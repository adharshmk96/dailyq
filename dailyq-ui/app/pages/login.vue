<script setup lang="ts">
const { pending, login } = useAuth()
const { fail, succeed } = useAuthFeedback()
const route = useRoute()

definePageMeta({
  middleware: 'guest'
})

const email = ref('')
const password = ref('')

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

  try {
    const res = await login({ email: email.value, password: password.value })
    succeed('Signed in', `Welcome back, ${res.user.name}.`)

    const redirect = route.query.redirect
    await navigateTo(typeof redirect === 'string' && redirect.startsWith('/') ? redirect : '/dashboard/overview')
  } catch (err) {
    fail('Sign in failed', apiErrorMessage(err, 'Invalid email or password.'))
  }
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

      <div class="flex items-center justify-end gap-3">
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
