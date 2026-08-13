import type { Ref } from 'vue'
import type { JournalNote, JournalTag, JournalTask, TagColor } from '~/types/journal'

/** Wire shapes as the API sends them (snake_case). */
interface ApiTag {
  id: string
  name: string
  color: TagColor
}

interface ApiTask {
  id: string
  title: string
  done: boolean
  date: string | null
  tag_ids: string[]
}

interface ApiNote {
  id: string
  body: string
  date: string | null
  tag_ids: string[]
}

interface ApiEntries {
  date: string | null
  tasks: ApiTask[]
  notes: ApiNote[]
}

export const TAGS_KEY = 'journal-tags'
export const ENTRIES_KEY = 'journal-entries'
export const DATES_KEY = 'journal-dates'
export const OVERVIEW_KEY = 'dashboard-overview'

function toTag(tag: ApiTag): JournalTag {
  return { id: tag.id, name: tag.name, color: tag.color }
}

function toTask(task: ApiTask): JournalTask {
  return { id: task.id, title: task.title, done: task.done, date: task.date, tagIds: task.tag_ids ?? [] }
}

function toNote(note: ApiNote): JournalNote {
  return { id: note.id, body: note.body, date: note.date, tagIds: note.tag_ids ?? [] }
}

/** Refreshes the views a write can invalidate, without blocking the caller. */
function refreshDerived() {
  return refreshNuxtData([DATES_KEY, OVERVIEW_KEY])
}

/**
 * Toasts whatever the API complained about and swallows the error, so a failed
 * write leaves the UI consistent instead of throwing into a template handler.
 */
export function useJournalFeedback() {
  const toast = useToast()

  function fail(err: unknown, fallback: string) {
    toast.add({
      title: fallback,
      description: apiErrorMessage(err, 'Please try again.'),
      icon: 'i-lucide-alert-circle',
      color: 'error'
    })
  }

  return { fail }
}

/**
 * The user's tags, shared across every page through one Nuxt data key.
 */
export function useJournalTags() {
  const api = useApiFetch()
  const { fail } = useJournalFeedback()

  const { data, pending, error, refresh } = useAsyncData<JournalTag[]>(
    TAGS_KEY,
    async () => {
      const res = await api<{ tags: ApiTag[] }>('/journal/tags')
      return (res.tags ?? []).map(toTag)
    },
    { default: () => [], dedupe: 'defer' }
  )

  const tags = computed(() => data.value ?? [])

  const tagsById = computed(() => {
    const map = new Map<string, JournalTag>()
    for (const tag of tags.value) map.set(tag.id, tag)
    return map
  })

  function getTagById(id: string) {
    return tagsById.value.get(id)
  }

  function sortByName(list: JournalTag[]) {
    return [...list].sort((a, b) => a.name.localeCompare(b.name))
  }

  async function addTag(name: string, color: TagColor) {
    const trimmed = name.trim()
    if (!trimmed) return null

    try {
      const created = toTag(await api<ApiTag>('/journal/tags', {
        method: 'POST',
        body: { name: trimmed, color }
      }))
      data.value = sortByName([...tags.value, created])
      return created
    } catch (err) {
      fail(err, 'Could not create tag')
      return null
    }
  }

  async function updateTag(id: string, updates: { name?: string, color?: TagColor }) {
    const existing = getTagById(id)
    if (!existing) return false

    const name = (updates.name ?? existing.name).trim()
    if (!name) return false

    try {
      const saved = toTag(await api<ApiTag>(`/journal/tags/${id}`, {
        method: 'PATCH',
        body: { name, color: updates.color ?? existing.color }
      }))
      data.value = sortByName(tags.value.map(t => (t.id === id ? saved : t)))
      return true
    } catch (err) {
      fail(err, 'Could not update tag')
      return false
    }
  }

  async function deleteTag(id: string) {
    try {
      await api(`/journal/tags/${id}`, { method: 'DELETE' })
      data.value = tags.value.filter(t => t.id !== id)
      // Items still carry the removed id in their local copy.
      await Promise.all([refreshNuxtData(ENTRIES_KEY), refreshDerived()])
      return true
    } catch (err) {
      fail(err, 'Could not delete tag')
      return false
    }
  }

  return { tags, tagsById, getTagById, pending, error, refresh, addTag, updateTag, deleteTag }
}

/**
 * Tasks and notes for one scope: a specific ISO date, or the undated "general"
 * bucket when `date` is null. Writes patch the local list from the API's own
 * response, so the board never drifts from the server.
 */
export function useJournalBoard(date: Ref<string | null>) {
  const api = useApiFetch()
  const { fail } = useJournalFeedback()

  const { data, pending, error, refresh } = useAsyncData<ApiEntries>(
    ENTRIES_KEY,
    () => api<ApiEntries>('/journal/entries', {
      params: date.value ? { date: date.value } : {}
    }),
    { watch: [date], dedupe: 'defer' }
  )

  const tasks = computed(() => (data.value?.tasks ?? []).map(toTask))
  const notes = computed(() => (data.value?.notes ?? []).map(toNote))

  /**
   * useAsyncData hands back a shallow ref, so every local write has to swap the
   * whole object for the board to re-render.
   */
  function patchEntries(changes: Partial<ApiEntries>) {
    if (!data.value) return
    data.value = { ...data.value, ...changes }
  }

  async function addTask(title: string, tagIds: string[] = []) {
    const trimmed = title.trim()
    if (!trimmed) return null

    try {
      const created = await api<ApiTask>('/journal/tasks', {
        method: 'POST',
        body: { title: trimmed, date: date.value, tag_ids: tagIds }
      })
      patchEntries({ tasks: [created, ...(data.value?.tasks ?? [])] })
      await refreshDerived()
      return toTask(created)
    } catch (err) {
      fail(err, 'Could not add task')
      return null
    }
  }

  async function addNote(body: string, tagIds: string[] = []) {
    const trimmed = body.trim()
    if (!trimmed) return null

    try {
      const created = await api<ApiNote>('/journal/notes', {
        method: 'POST',
        body: { body: trimmed, date: date.value, tag_ids: tagIds }
      })
      patchEntries({ notes: [created, ...(data.value?.notes ?? [])] })
      await refreshDerived()
      return toNote(created)
    } catch (err) {
      fail(err, 'Could not add note')
      return null
    }
  }

  async function updateTask(id: string, title: string, tagIds?: string[]) {
    const trimmed = title.trim()
    if (!trimmed) return false

    const current = data.value?.tasks.find(t => t.id === id)
    if (!current) return false

    try {
      const saved = await api<ApiTask>(`/journal/tasks/${id}`, {
        method: 'PATCH',
        body: { title: trimmed, date: current.date, tag_ids: tagIds ?? current.tag_ids }
      })
      patchEntries({ tasks: (data.value?.tasks ?? []).map(t => (t.id === id ? saved : t)) })
      return true
    } catch (err) {
      fail(err, 'Could not update task')
      return false
    }
  }

  async function updateNote(id: string, body: string, tagIds?: string[]) {
    const trimmed = body.trim()
    if (!trimmed) return false

    const current = data.value?.notes.find(n => n.id === id)
    if (!current) return false

    try {
      const saved = await api<ApiNote>(`/journal/notes/${id}`, {
        method: 'PATCH',
        body: { body: trimmed, date: current.date, tag_ids: tagIds ?? current.tag_ids }
      })
      patchEntries({ notes: (data.value?.notes ?? []).map(n => (n.id === id ? saved : n)) })
      return true
    } catch (err) {
      fail(err, 'Could not update note')
      return false
    }
  }

  async function deleteTask(id: string) {
    const previous = data.value?.tasks ?? []
    patchEntries({ tasks: previous.filter(t => t.id !== id) })

    try {
      await api(`/journal/tasks/${id}`, { method: 'DELETE' })
      await refreshDerived()
    } catch (err) {
      patchEntries({ tasks: previous })
      fail(err, 'Could not delete task')
    }
  }

  async function deleteNote(id: string) {
    const previous = data.value?.notes ?? []
    patchEntries({ notes: previous.filter(n => n.id !== id) })

    try {
      await api(`/journal/notes/${id}`, { method: 'DELETE' })
      await refreshDerived()
    } catch (err) {
      patchEntries({ notes: previous })
      fail(err, 'Could not delete note')
    }
  }

  /** Flips the checkbox first and rolls back if the server disagrees. */
  async function toggleTaskDone(id: string) {
    const task = data.value?.tasks.find(t => t.id === id)
    if (!task) return

    const next = !task.done
    const setDone = (done: boolean) =>
      patchEntries({ tasks: (data.value?.tasks ?? []).map(t => (t.id === id ? { ...t, done } : t)) })

    setDone(next)

    try {
      await api(`/journal/tasks/${id}/done`, { method: 'PATCH', body: { done: next } })
      await refreshDerived()
    } catch (err) {
      setDone(!next)
      fail(err, 'Could not update task')
    }
  }

  return {
    tasks,
    notes,
    pending,
    error,
    refresh,
    addTask,
    addNote,
    updateTask,
    updateNote,
    deleteTask,
    deleteNote,
    toggleTaskDone
  }
}

/** The set of dates that have at least one item — the calendar's day dots. */
export function useJournalDates() {
  const api = useApiFetch()

  const { data, refresh } = useAsyncData<string[]>(
    DATES_KEY,
    async () => {
      const res = await api<{ dates: string[] }>('/journal/dates')
      return res.dates ?? []
    },
    { default: () => [] }
  )

  const datesWithItems = computed(() => new Set(data.value ?? []))

  return { datesWithItems, refresh }
}

/** Filters any tag-carrying list down to the active tag. */
export function filterItemsByTag<T extends { tagIds: string[] }>(items: T[], tagId: string | null) {
  if (!tagId) return items
  return items.filter(item => item.tagIds.includes(tagId))
}
