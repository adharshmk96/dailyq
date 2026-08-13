/** Shared toast helpers for the auth screens. */
export function useAuthFeedback() {
  const toast = useToast()

  function fail(title: string, description: string) {
    toast.add({
      title,
      description,
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
  }

  function succeed(title: string, description: string) {
    toast.add({
      title,
      description,
      icon: 'i-lucide-check-circle'
    })
  }

  return { fail, succeed }
}
