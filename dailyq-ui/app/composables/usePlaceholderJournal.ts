import type { JournalNote, JournalTag, JournalTask, TagColor } from '~/types/journal'

function createId() {
  return crypto.randomUUID()
}

function todayIso() {
  return new Date().toISOString().slice(0, 10)
}

function daysAgoIso(days: number) {
  const d = new Date()
  d.setDate(d.getDate() - days)
  return d.toISOString().slice(0, 10)
}

function daysFromNowIso(days: number) {
  const d = new Date()
  d.setDate(d.getDate() + days)
  return d.toISOString().slice(0, 10)
}

function buildSeedTags(): JournalTag[] {
  return [
    { id: 'tag-work', name: 'Work', color: 'primary' },
    { id: 'tag-personal', name: 'Personal', color: 'success' },
    { id: 'tag-ideas', name: 'Ideas', color: 'warning' },
    { id: 'tag-health', name: 'Health', color: 'info' }
  ]
}

function buildSeedTasks(): JournalTask[] {
  return [
    { id: 'task-inbox-1', title: 'Review weekly goals', done: false, date: null, tagIds: ['tag-work'] },
    { id: 'task-inbox-2', title: 'Read design notes', done: true, date: null, tagIds: ['tag-ideas'] },
    { id: 'task-inbox-3', title: 'Clear open general items', done: false, date: null, tagIds: [] },
    { id: 'task-today-1', title: 'Morning stretch', done: true, date: todayIso(), tagIds: ['tag-health'] },
    { id: 'task-today-2', title: 'Write journal entry', done: false, date: todayIso(), tagIds: ['tag-personal'] },
    { id: 'task-today-3', title: 'Plan tomorrow focus', done: false, date: todayIso(), tagIds: ['tag-work'] },
    { id: 'task-soon-1', title: 'Team standup prep', done: false, date: daysFromNowIso(1), tagIds: ['tag-work'] },
    { id: 'task-soon-2', title: 'Grocery run', done: false, date: daysFromNowIso(2), tagIds: ['tag-personal'] },
    { id: 'task-past-1', title: 'Ship dashboard UI', done: true, date: daysAgoIso(1), tagIds: ['tag-work'] }
  ]
}

function buildSeedNotes(): JournalNote[] {
  return [
    { id: 'note-inbox-1', body: 'Ideas for the habit tracker — keep it simple.', date: null, tagIds: ['tag-ideas'] },
    { id: 'note-inbox-2', body: 'Quote: small steps compound.', date: null, tagIds: ['tag-personal'] },
    { id: 'note-today-1', body: 'Felt focused after the morning walk.', date: todayIso(), tagIds: ['tag-health'] },
    { id: 'note-today-2', body: 'Sketch calendar chips for days with items.', date: todayIso(), tagIds: ['tag-ideas'] },
    { id: 'note-soon-1', body: 'Tomorrow: try a shorter standup agenda.', date: daysFromNowIso(1), tagIds: ['tag-work'] },
    { id: 'note-past-1', body: 'Reflection: shipping UI before API was the right call.', date: daysAgoIso(2), tagIds: ['tag-work'] }
  ]
}

export function usePlaceholderJournal() {
  const tags = useState<JournalTag[]>('placeholder-tags', () => buildSeedTags())
  const tasks = useState<JournalTask[]>('placeholder-tasks', () => buildSeedTasks())
  const notes = useState<JournalNote[]>('placeholder-notes', () => buildSeedNotes())

  const today = computed(() => todayIso())

  const generalTasks = computed(() => tasks.value.filter(t => t.date === null))
  const generalNotes = computed(() => notes.value.filter(n => n.date === null))

  const todayTasks = computed(() => tasks.value.filter(t => t.date === today.value))
  const todayNotes = computed(() => notes.value.filter(n => n.date === today.value))

  const openGeneralCount = computed(() => generalTasks.value.filter(t => !t.done).length)
  const completedTodayCount = computed(() => todayTasks.value.filter(t => t.done).length)
  const totalTasksCount = computed(() => tasks.value.length)

  const notesThisWeekCount = computed(() => {
    const weekStart = new Date()
    weekStart.setDate(weekStart.getDate() - 6)
    const startIso = weekStart.toISOString().slice(0, 10)
    return notes.value.filter(n => n.date && n.date >= startIso).length
  })

  const recentNotes = computed(() =>
    [...notes.value]
      .sort((a, b) => (b.date ?? '').localeCompare(a.date ?? ''))
      .slice(0, 5)
  )

  function getTasksForDate(date: string) {
    return tasks.value.filter(t => t.date === date)
  }

  function getNotesForDate(date: string) {
    return notes.value.filter(n => n.date === date)
  }

  function filterItemsByTag<T extends { tagIds: string[] }>(items: T[], tagId: string | null) {
    if (!tagId) return items
    return items.filter(item => item.tagIds.includes(tagId))
  }

  function getTagById(id: string) {
    return tags.value.find(t => t.id === id)
  }

  const datesWithItems = computed(() => {
    const set = new Set<string>()
    for (const t of tasks.value) {
      if (t.date) set.add(t.date)
    }
    for (const n of notes.value) {
      if (n.date) set.add(n.date)
    }
    return set
  })

  function addTag(name: string, color: TagColor) {
    const trimmed = name.trim()
    if (!trimmed) return null
    const tag: JournalTag = {
      id: createId(),
      name: trimmed,
      color
    }
    tags.value = [...tags.value, tag]
    return tag
  }

  function updateTag(id: string, updates: { name?: string, color?: TagColor }) {
    const existing = tags.value.find(t => t.id === id)
    if (!existing) return false

    const name = updates.name !== undefined ? updates.name.trim() : existing.name
    if (!name) return false

    tags.value = tags.value.map(t =>
      t.id === id
        ? { ...t, name, color: updates.color ?? t.color }
        : t
    )
    return true
  }

  function deleteTag(id: string) {
    tags.value = tags.value.filter(t => t.id !== id)
    tasks.value = tasks.value.map(t => ({
      ...t,
      tagIds: t.tagIds.filter(tagId => tagId !== id)
    }))
    notes.value = notes.value.map(n => ({
      ...n,
      tagIds: n.tagIds.filter(tagId => tagId !== id)
    }))
  }

  function addTask(title: string, date: string | null = null, tagIds: string[] = []) {
    const trimmed = title.trim()
    if (!trimmed) return null
    const task: JournalTask = {
      id: createId(),
      title: trimmed,
      done: false,
      date,
      tagIds: [...tagIds]
    }
    tasks.value = [task, ...tasks.value]
    return task
  }

  function addNote(body: string, date: string | null = null, tagIds: string[] = []) {
    const trimmed = body.trim()
    if (!trimmed) return null
    const note: JournalNote = {
      id: createId(),
      body: trimmed,
      date,
      tagIds: [...tagIds]
    }
    notes.value = [note, ...notes.value]
    return note
  }

  function updateTask(id: string, title: string, tagIds?: string[]) {
    const trimmed = title.trim()
    if (!trimmed) return false
    tasks.value = tasks.value.map(t =>
      t.id === id
        ? { ...t, title: trimmed, tagIds: tagIds ?? t.tagIds }
        : t
    )
    return true
  }

  function updateNote(id: string, body: string, tagIds?: string[]) {
    const trimmed = body.trim()
    if (!trimmed) return false
    notes.value = notes.value.map(n =>
      n.id === id
        ? { ...n, body: trimmed, tagIds: tagIds ?? n.tagIds }
        : n
    )
    return true
  }

  function deleteTask(id: string) {
    tasks.value = tasks.value.filter(t => t.id !== id)
  }

  function deleteNote(id: string) {
    notes.value = notes.value.filter(n => n.id !== id)
  }

  function toggleTaskDone(id: string) {
    tasks.value = tasks.value.map(t => (t.id === id ? { ...t, done: !t.done } : t))
  }

  return {
    tags,
    tasks,
    notes,
    today,
    generalTasks,
    generalNotes,
    todayTasks,
    todayNotes,
    openGeneralCount,
    completedTodayCount,
    totalTasksCount,
    notesThisWeekCount,
    recentNotes,
    getTasksForDate,
    getNotesForDate,
    filterItemsByTag,
    getTagById,
    datesWithItems,
    addTag,
    updateTag,
    deleteTag,
    addTask,
    addNote,
    updateTask,
    updateNote,
    deleteTask,
    deleteNote,
    toggleTaskDone
  }
}
