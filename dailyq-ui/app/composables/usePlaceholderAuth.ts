/**
 * Placeholder auth actions. No API is wired up yet — every call just waits a
 * beat so the UI can show its loading state, then resolves.
 */
export function usePlaceholderAuth() {
  const toast = useToast()
  const pending = ref(false)

  function fail(title: string, description: string) {
    toast.add({
      title,
      description,
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
  }

  async function submit(options: { title: string, description: string, redirectTo?: string }) {
    if (pending.value) {
      return
    }

    pending.value = true
    await new Promise(resolve => setTimeout(resolve, 700))
    pending.value = false

    toast.add({
      title: options.title,
      description: options.description,
      icon: 'i-lucide-info'
    })

    if (options.redirectTo) {
      await navigateTo(options.redirectTo)
    }
  }

  return { pending, submit, fail }
}

export const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
