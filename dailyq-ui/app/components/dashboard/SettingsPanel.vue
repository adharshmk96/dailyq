<script setup lang="ts">
const VALID_TABS = ['account', 'data', 'notification'] as const
type SettingsTab = typeof VALID_TABS[number]

const route = useRoute()
const toast = useToast()

function isSettingsTab(value: string): value is SettingsTab {
  return (VALID_TABS as readonly string[]).includes(value)
}

const activeTab = computed<SettingsTab>(() => {
  const param = route.params.tab
  if (typeof param === 'string' && isSettingsTab(param)) {
    return param
  }
  return 'account'
})

const tabItems = [
  {
    label: 'Account',
    icon: 'i-lucide-user',
    value: 'account' as const,
    to: '/dashboard/settings/account'
  },
  {
    label: 'Data',
    icon: 'i-lucide-database',
    value: 'data' as const,
    to: '/dashboard/settings/data'
  },
  {
    label: 'Notification',
    icon: 'i-lucide-bell',
    value: 'notification' as const,
    to: '/dashboard/settings/notification'
  }
]

const { user, pending, changePassword, updateProfile } = useAuth()

const displayName = computed(() => user.value?.name || '')
const email = computed(() => user.value?.email || '')

const profileName = ref('')
const profileEmail = ref('')

watch(user, (current) => {
  if (!current) return
  profileName.value = current.name
  profileEmail.value = current.email
}, { immediate: true })

const isProfileDirty = computed(() =>
  profileName.value.trim() !== displayName.value
  || profileEmail.value.trim().toLowerCase() !== email.value.toLowerCase()
)

const canSubmitProfile = computed(() =>
  profileName.value.trim().length > 0
  && EMAIL_PATTERN.test(profileEmail.value.trim())
  && isProfileDirty.value
)

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

const canSubmitPassword = computed(() =>
  currentPassword.value.length > 0
  && newPassword.value.length >= 8
  && newPassword.value === confirmPassword.value
)

const emailNotifications = ref(true)
const pushNotifications = ref(false)
const weeklyDigest = ref(true)

function comingSoon(title: string) {
  toast.add({
    title,
    description: 'Coming soon',
    icon: 'i-lucide-info'
  })
}

const { exportCsv, importCsv, exportPending, importPending } = useJournalDataTransfer()
const importInputRef = ref<HTMLInputElement | null>(null)

function formatImportSummary(result: ImportResult): string {
  const created = result.created.tags + result.created.tasks + result.created.notes
  const updated = result.updated.tags + result.updated.tasks + result.updated.notes
  const parts: string[] = []
  if (created > 0) parts.push(`${created} created`)
  if (updated > 0) parts.push(`${updated} updated`)
  if (result.skipped > 0) parts.push(`${result.skipped} skipped`)
  return parts.length > 0 ? parts.join(', ') : 'No changes'
}

async function onExportData() {
  try {
    await exportCsv()
    toast.add({
      title: 'Export complete',
      description: 'Your journal data was downloaded as CSV.',
      icon: 'i-lucide-download'
    })
  } catch {
    // useJournalDataTransfer already shows the error toast
  }
}

function onImportClick() {
  importInputRef.value?.click()
}

async function onImportFileSelected(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  const result = await importCsv(file)
  if (!result) return

  const color = result.errors.length > 0 ? 'warning' : 'success'
  toast.add({
    title: 'Import complete',
    description: formatImportSummary(result),
    icon: result.errors.length > 0 ? 'i-lucide-alert-triangle' : 'i-lucide-upload',
    color
  })
}

function resetPasswordForm() {
  currentPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
  showCurrentPassword.value = false
  showNewPassword.value = false
  showConfirmPassword.value = false
}

async function onSaveProfile() {
  const name = profileName.value.trim()
  const emailValue = profileEmail.value.trim()

  if (!name || !emailValue) {
    toast.add({
      title: 'Missing fields',
      description: 'Please fill in your display name and email.',
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
    return
  }

  if (!EMAIL_PATTERN.test(emailValue)) {
    toast.add({
      title: 'Invalid email',
      description: 'Enter a valid email address.',
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
    return
  }

  if (!isProfileDirty.value) {
    return
  }

  try {
    await updateProfile({ name, email: emailValue })

    toast.add({
      title: 'Profile updated',
      description: 'Your account details were saved.',
      icon: 'i-lucide-check-circle'
    })
  } catch (err) {
    toast.add({
      title: 'Could not update profile',
      description: apiErrorMessage(err),
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
  }
}

async function onChangePassword() {
  if (!currentPassword.value || !newPassword.value || !confirmPassword.value) {
    toast.add({
      title: 'Missing fields',
      description: 'Please fill in all password fields.',
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
    return
  }

  if (newPassword.value.length < 8) {
    toast.add({
      title: 'Password too short',
      description: 'Password must be at least 8 characters.',
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
    return
  }

  if (newPassword.value !== confirmPassword.value) {
    toast.add({
      title: 'Passwords do not match',
      description: 'New password and confirmation must match.',
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
    return
  }

  if (newPassword.value === currentPassword.value) {
    toast.add({
      title: 'Same password',
      description: 'New password must be different from your current password.',
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
    return
  }

  try {
    await changePassword({
      current_password: currentPassword.value,
      new_password: newPassword.value
    })

    toast.add({
      title: 'Password updated',
      description: 'Your other sessions have been signed out.',
      icon: 'i-lucide-check-circle'
    })
    resetPasswordForm()
  } catch (err) {
    toast.add({
      title: 'Could not update password',
      description: apiErrorMessage(err),
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
  }
}

function onSaveNotifications() {
  comingSoon('Notifications')
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight text-highlighted">
        Settings
      </h1>
      <p class="mt-1 text-sm text-muted">
        Manage your account, data, and notification preferences.
      </p>
    </div>

    <div class="flex flex-col gap-6 sm:flex-row sm:items-start sm:gap-8">
      <nav
        class="flex w-full shrink-0 flex-col gap-1.5 sm:w-56"
        aria-label="Settings sections"
      >
        <NuxtLink
          v-for="item in tabItems"
          :key="item.value"
          :to="item.to"
          class="flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium transition-colors"
          :class="activeTab === item.value
            ? 'bg-primary/10 text-primary'
            : 'text-muted hover:bg-elevated hover:text-highlighted'"
        >
          <UIcon
            :name="item.icon"
            class="size-5 shrink-0"
          />
          <span>{{ item.label }}</span>
        </NuxtLink>
      </nav>

      <div class="min-w-0 flex-1">
        <div
          v-if="activeTab === 'account'"
          class="space-y-5"
        >
          <UCard :ui="{ body: 'space-y-5 p-4 sm:p-5' }">
            <div class="flex items-center gap-4">
              <UAvatar
                :alt="displayName || email"
                size="lg"
              />
              <div class="min-w-0">
                <p class="truncate text-sm font-medium text-highlighted">
                  {{ displayName }}
                </p>
                <p class="truncate text-sm text-muted">
                  {{ email }}
                </p>
              </div>
            </div>

            <UFormField
              label="Display name"
              class="w-full"
            >
              <UInput
                v-model="profileName"
                autocomplete="name"
                class="w-full"
              />
            </UFormField>

            <UFormField
              label="Email"
              class="w-full"
            >
              <UInput
                v-model="profileEmail"
                type="email"
                autocomplete="email"
                class="w-full"
              />
            </UFormField>

            <div class="flex justify-end">
              <UButton
                label="Save profile"
                icon="i-lucide-user-pen"
                :loading="pending"
                :disabled="!canSubmitProfile"
                @click="onSaveProfile"
              />
            </div>
          </UCard>

          <UCard :ui="{ body: 'space-y-5 p-4 sm:p-5' }">
            <div>
              <h2 class="text-sm font-medium text-highlighted">
                Change password
              </h2>
              <p class="mt-1 text-sm text-muted">
                Update your account password.
              </p>
            </div>

            <UFormField
              label="Current password"
              class="w-full"
            >
              <UInput
                v-model="currentPassword"
                :type="showCurrentPassword ? 'text' : 'password'"
                class="w-full"
                autocomplete="current-password"
              >
                <template #trailing>
                  <UButton
                    color="neutral"
                    variant="link"
                    :icon="showCurrentPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                    :aria-label="showCurrentPassword ? 'Hide password' : 'Show password'"
                    @click="showCurrentPassword = !showCurrentPassword"
                  />
                </template>
              </UInput>
            </UFormField>

            <UFormField
              label="New password"
              class="w-full"
            >
              <UInput
                v-model="newPassword"
                :type="showNewPassword ? 'text' : 'password'"
                class="w-full"
                autocomplete="new-password"
              >
                <template #trailing>
                  <UButton
                    color="neutral"
                    variant="link"
                    :icon="showNewPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                    :aria-label="showNewPassword ? 'Hide password' : 'Show password'"
                    @click="showNewPassword = !showNewPassword"
                  />
                </template>
              </UInput>
            </UFormField>

            <UFormField
              label="Confirm password"
              class="w-full"
            >
              <UInput
                v-model="confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                class="w-full"
                autocomplete="new-password"
              >
                <template #trailing>
                  <UButton
                    color="neutral"
                    variant="link"
                    :icon="showConfirmPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                    :aria-label="showConfirmPassword ? 'Hide password' : 'Show password'"
                    @click="showConfirmPassword = !showConfirmPassword"
                  />
                </template>
              </UInput>
            </UFormField>

            <div class="flex justify-end">
              <UButton
                label="Update password"
                icon="i-lucide-lock"
                :loading="pending"
                :disabled="!canSubmitPassword"
                @click="onChangePassword"
              />
            </div>
          </UCard>
        </div>

        <UCard
          v-else-if="activeTab === 'data'"
          :ui="{ body: 'space-y-5 p-4 sm:p-5' }"
        >
          <div>
            <h2 class="text-sm font-medium text-highlighted">
              Journal data
            </h2>
            <p class="mt-1 text-sm text-muted">
              Export your tasks, notes, and tags as CSV, or import a previously exported file. Re-importing the same file updates existing entries by id — no duplicates.
            </p>
          </div>

          <input
            ref="importInputRef"
            type="file"
            accept=".csv,text/csv"
            class="hidden"
            @change="onImportFileSelected"
          >

          <div class="flex flex-wrap gap-2">
            <UButton
              label="Export data"
              icon="i-lucide-download"
              color="neutral"
              variant="soft"
              :loading="exportPending"
              @click="onExportData"
            />
            <UButton
              label="Import data"
              icon="i-lucide-upload"
              color="neutral"
              variant="soft"
              :loading="importPending"
              @click="onImportClick"
            />
          </div>
        </UCard>

        <UCard
          v-else-if="activeTab === 'notification'"
          :ui="{ body: 'space-y-5 p-4 sm:p-5' }"
        >
          <USwitch
            v-model="emailNotifications"
            label="Email notifications"
            description="Get updates about tasks and notes by email."
          />

          <USwitch
            v-model="pushNotifications"
            label="Push notifications"
            description="Receive alerts on this device."
          />

          <USwitch
            v-model="weeklyDigest"
            label="Weekly digest"
            description="A summary of your week every Monday."
          />

          <div class="flex justify-end">
            <UButton
              label="Save preferences"
              icon="i-lucide-check"
              @click="onSaveNotifications"
            />
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>
