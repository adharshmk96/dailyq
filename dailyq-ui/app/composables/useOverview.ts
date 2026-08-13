import type { TagColor } from '~/types/journal'

export interface OverviewTag {
  id: string
  name: string
  color: TagColor
}

export interface OverviewTask {
  id: string
  title: string
  done: boolean
  date: string | null
  tag_ids: string[]
}

export interface OverviewNote {
  id: string
  body: string
  date: string | null
  tag_ids: string[]
}

export interface OverviewStats {
  total_tasks: number
  completed_today: number
  open_general: number
  notes_this_week: number
}

export interface Overview {
  today: string
  stats: OverviewStats
  tags: OverviewTag[]
  today_tasks: OverviewTask[]
  recent_notes: OverviewNote[]
}

/**
 * Loads the dashboard overview from the API. One request backs the whole page:
 * stat cards, today's tasks and recent notes.
 */
export function useOverview() {
  const api = useApiFetch()

  const { data, pending, error, refresh } = useAsyncData<Overview>(
    OVERVIEW_KEY,
    () => api<Overview>('/journal/overview'),
    { dedupe: 'defer' }
  )

  const tagsById = computed(() => {
    const map = new Map<string, OverviewTag>()
    for (const tag of data.value?.tags ?? []) {
      map.set(tag.id, tag)
    }
    return map
  })

  function getTagById(id: string) {
    return tagsById.value.get(id)
  }

  /**
   * Flips a task optimistically and reconciles with the server response so the
   * stat cards stay in step with the checkbox.
   */
  async function toggleTaskDone(id: string) {
    const task = data.value?.today_tasks.find(t => t.id === id)
    if (!task || !data.value) return

    const next = !task.done
    // useAsyncData's ref is shallow, so swap the object to flip the checkbox.
    const setDone = (done: boolean) => {
      if (!data.value) return
      data.value = {
        ...data.value,
        today_tasks: data.value.today_tasks.map(t => (t.id === id ? { ...t, done } : t))
      }
    }

    setDone(next)

    try {
      await api(`/journal/tasks/${id}/done`, { method: 'PATCH', body: { done: next } })
      await refresh()
    } catch (err) {
      setDone(!next)
      throw err
    }
  }

  return {
    overview: data,
    pending,
    error,
    refresh,
    getTagById,
    toggleTaskDone
  }
}
