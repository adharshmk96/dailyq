import type { ImportResult } from '~/types/journal'
import { DATES_KEY, OVERVIEW_KEY, refreshAllJournalEntries, TAGS_KEY } from '~/composables/useJournal'

export function useJournalDataTransfer() {
  const api = useApiFetch()
  const token = useAuthToken()
  const config = useRuntimeConfig()
  const { fail } = useJournalFeedback()

  const exportPending = ref(false)
  const importPending = ref(false)

  async function exportCsv(): Promise<void> {
    exportPending.value = true
    try {
      const headers: HeadersInit = {}
      if (token.value) {
        headers.Authorization = `Bearer ${token.value}`
      }

      const response = await fetch(`${config.public.apiBase}/journal/export`, { headers })
      if (!response.ok) {
        const body = await response.json().catch(() => null)
        throw { data: body }
      }

      const blob = await response.blob()
      const url = URL.createObjectURL(blob)
      const link = document.createElement('a')
      const date = new Date().toISOString().slice(0, 10)
      link.href = url
      link.download = `dailyq-export-${date}.csv`
      link.click()
      URL.revokeObjectURL(url)
    } catch (err) {
      fail(err, 'Could not export data')
      throw err
    } finally {
      exportPending.value = false
    }
  }

  async function importCsv(file: File): Promise<ImportResult | null> {
    importPending.value = true
    try {
      const form = new FormData()
      form.append('file', file)

      const result = await api<ImportResult>('/journal/import', {
        method: 'POST',
        body: form
      })

      await Promise.all([
        refreshNuxtData([TAGS_KEY, DATES_KEY, OVERVIEW_KEY]),
        refreshAllJournalEntries()
      ])
      return result
    } catch (err) {
      fail(err, 'Could not import data')
      return null
    } finally {
      importPending.value = false
    }
  }

  return {
    exportCsv,
    importCsv,
    exportPending,
    importPending
  }
}
