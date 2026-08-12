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

const displayName = ref('Alex Rivera')
const email = ref('alex.rivera@example.com')

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

function onSaveAccount() {
  comingSoon('Account')
}

function onExportData() {
  comingSoon('Export data')
}

function onImportData() {
  comingSoon('Import data')
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
        <UCard
          v-if="activeTab === 'account'"
          :ui="{ body: 'space-y-5 p-4 sm:p-5' }"
        >
          <div class="flex items-center gap-4">
            <UAvatar
              src="https://i.pravatar.cc/80?u=dailyq"
              alt="Alex Rivera"
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
              v-model="displayName"
              class="w-full"
              autocomplete="name"
            />
          </UFormField>

          <UFormField
            label="Email"
            class="w-full"
          >
            <UInput
              v-model="email"
              type="email"
              class="w-full"
              autocomplete="email"
            />
          </UFormField>

          <div class="flex justify-end">
            <UButton
              label="Save changes"
              icon="i-lucide-check"
              @click="onSaveAccount"
            />
          </div>
        </UCard>

        <UCard
          v-else-if="activeTab === 'data'"
          :ui="{ body: 'space-y-5 p-4 sm:p-5' }"
        >
          <div>
            <h2 class="text-sm font-medium text-highlighted">
              Local data
            </h2>
            <p class="mt-1 text-sm text-muted">
              Journal data is stored as placeholder data in this browser session only. Export and import will be available later.
            </p>
          </div>

          <div class="flex flex-wrap gap-2">
            <UButton
              label="Export data"
              icon="i-lucide-download"
              color="neutral"
              variant="soft"
              @click="onExportData"
            />
            <UButton
              label="Import data"
              icon="i-lucide-upload"
              color="neutral"
              variant="soft"
              @click="onImportData"
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
