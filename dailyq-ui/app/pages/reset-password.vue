<script setup lang="ts">
const route = useRoute()
const { pending, resetPassword } = useAuth()
const { fail, succeed } = useAuthFeedback()

const password = ref('')
const confirmPassword = ref('')
const done = ref(false)

// The reset link carries the token; the API validates it on submit.
const token = computed(() => {
  const value = route.query.token
  return typeof value === 'string' ? value : ''
})

useHead({
  title: 'Set a new password — DailyQ'
})

async function onSubmit() {
  if (!token.value) {
    fail('Missing reset token', 'Open this page from the link in your reset email.')
    return
  }

  if (!password.value || !confirmPassword.value) {
    fail('Missing fields', 'Enter and confirm your new password.')
    return
  }

  if (password.value.length < 8) {
    fail('Password too short', 'Password must be at least 8 characters.')
    return
  }

  if (password.value !== confirmPassword.value) {
    fail('Passwords do not match', 'New password and confirmation must match.')
    return
  }

  try {
    await resetPassword(token.value, password.value)
    succeed('Password updated', 'Sign in with your new password.')
    done.value = true
  } catch (err) {
    fail('Could not reset password', apiErrorMessage(err, 'This reset link is invalid or has expired.'))
  }
}
</script>

<template>
  <AuthShell
    icon="i-lucide-lock-keyhole"
    title="Set a new password"
    description="Choose something you'll remember. At least 8 characters."
  >
    <div
      v-if="done"
      class="space-y-5 text-center"
    >
      <span class="mx-auto flex size-11 items-center justify-center rounded-full bg-primary/10 text-primary">
        <UIcon
          name="i-lucide-check"
          class="size-5"
        />
      </span>
      <div>
        <p class="text-sm font-medium text-highlighted">
          Your password is set
        </p>
        <p class="mt-1 text-sm leading-relaxed text-muted">
          You can sign in with your new password now.
        </p>
      </div>
      <UButton
        to="/login"
        block
        label="Go to sign in"
        trailing-icon="i-lucide-arrow-right"
      />
    </div>

    <form
      v-else
      class="space-y-5"
      @submit.prevent="onSubmit"
    >
      <UAlert
        v-if="!token"
        color="neutral"
        variant="subtle"
        icon="i-lucide-info"
        title="No reset token"
        description="Open this page from the link in your reset email to set a new password."
      />

      <AuthPasswordField
        v-model="password"
        label="New password"
        hint="8+ characters"
        placeholder="••••••••"
        autocomplete="new-password"
      />

      <AuthPasswordField
        v-model="confirmPassword"
        label="Confirm new password"
        placeholder="••••••••"
        autocomplete="new-password"
      />

      <UButton
        type="submit"
        block
        size="lg"
        label="Update password"
        :loading="pending"
      />
    </form>

    <template #footer>
      <ULink
        to="/login"
        class="font-medium text-primary"
      >
        Back to sign in
      </ULink>
    </template>
  </AuthShell>
</template>
