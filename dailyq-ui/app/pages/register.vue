<script setup lang="ts">
const { pending, register } = useAuth()
const { fail, succeed } = useAuthFeedback()

definePageMeta({
  middleware: 'guest'
})

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const acceptedTerms = ref(false)

useHead({
  title: 'Create your account — DailyQ'
})

async function onSubmit() {
  if (!name.value || !email.value || !password.value || !confirmPassword.value) {
    fail('Missing fields', 'Please fill in every field to continue.')
    return
  }

  if (!EMAIL_PATTERN.test(email.value)) {
    fail('Invalid email', 'Enter a valid email address.')
    return
  }

  if (password.value.length < 8) {
    fail('Password too short', 'Password must be at least 8 characters.')
    return
  }

  if (password.value !== confirmPassword.value) {
    fail('Passwords do not match', 'Password and confirmation must match.')
    return
  }

  if (!acceptedTerms.value) {
    fail('Terms not accepted', 'Please accept the terms to create an account.')
    return
  }

  try {
    await register({ name: name.value, email: email.value, password: password.value })
    succeed('Account created', 'You are signed in and ready to go.')
    await navigateTo('/dashboard/overview')
  } catch (err) {
    fail('Could not create account', apiErrorMessage(err))
  }
}
</script>

<template>
  <AuthShell
    icon="i-lucide-user-plus"
    title="Create your account"
    description="One page a day. It takes about a minute to get started."
  >
    <form
      class="space-y-5"
      @submit.prevent="onSubmit"
    >
      <UFormField
        label="Display name"
        class="w-full"
      >
        <UInput
          v-model="name"
          placeholder="Alex Rivera"
          autocomplete="name"
          class="w-full"
        />
      </UFormField>

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
        hint="8+ characters"
        placeholder="••••••••"
        autocomplete="new-password"
      />

      <AuthPasswordField
        v-model="confirmPassword"
        label="Confirm password"
        placeholder="••••••••"
        autocomplete="new-password"
      />

      <UCheckbox v-model="acceptedTerms">
        <template #label>
          <span class="text-sm text-muted">
            I agree to the terms of service and privacy policy.
          </span>
        </template>
      </UCheckbox>

      <UButton
        type="submit"
        block
        size="lg"
        label="Create account"
        :loading="pending"
      />
    </form>

    <template #footer>
      Already have an account?
      <ULink
        to="/login"
        class="font-medium text-primary"
      >
        Sign in
      </ULink>
    </template>
  </AuthShell>
</template>
