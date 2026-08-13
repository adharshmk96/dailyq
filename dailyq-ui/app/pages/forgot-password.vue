<script setup lang="ts">
const { pending, forgotPassword } = useAuth()
const { fail, succeed } = useAuthFeedback()

const email = ref('')
const sent = ref(false)

useHead({
  title: 'Reset your password — DailyQ'
})

async function onSubmit() {
  if (!email.value) {
    fail('Missing email', 'Enter the email you signed up with.')
    return
  }

  if (!EMAIL_PATTERN.test(email.value)) {
    fail('Invalid email', 'Enter a valid email address.')
    return
  }

  try {
    await forgotPassword(email.value)
    succeed('Reset link sent', 'If that email is registered, a reset link is on its way.')
    sent.value = true
  } catch (err) {
    fail('Could not send reset link', apiErrorMessage(err))
  }
}
</script>

<template>
  <AuthShell
    icon="i-lucide-key-round"
    title="Forgot your password?"
    description="Enter your email and we'll send you a link to set a new one."
  >
    <div
      v-if="sent"
      class="space-y-5 text-center"
    >
      <span class="mx-auto flex size-11 items-center justify-center rounded-full bg-primary/10 text-primary">
        <UIcon
          name="i-lucide-mail-check"
          class="size-5"
        />
      </span>
      <div>
        <p class="text-sm font-medium text-highlighted">
          Check your inbox
        </p>
        <p class="mt-1 text-sm leading-relaxed text-muted">
          If an account exists for
          <span class="font-medium text-highlighted">{{ email }}</span>,
          a reset link is on its way.
        </p>
      </div>

      <div class="flex flex-col gap-2">
        <UButton
          to="/login"
          block
          label="Back to sign in"
          trailing-icon="i-lucide-arrow-right"
        />
        <UButton
          block
          variant="ghost"
          color="neutral"
          label="Use a different email"
          @click="sent = false"
        />
      </div>
    </div>

    <form
      v-else
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

      <UButton
        type="submit"
        block
        size="lg"
        label="Send reset link"
        :loading="pending"
      />
    </form>

    <template #footer>
      Remembered it?
      <ULink
        to="/login"
        class="font-medium text-primary"
      >
        Back to sign in
      </ULink>
    </template>
  </AuthShell>
</template>
