import type { JournalNote, JournalTask } from '~/types/journal'

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

function buildSeedTasks(): JournalTask[] {
  return [
    { id: 'task-inbox-1', title: 'Review weekly goals', done: false, date: null },
    { id: 'task-inbox-2', title: 'Read design notes', done: true, date: null },
    { id: 'task-inbox-3', title: 'Clear open general items', done: false, date: null },
    { id: 'task-today-1', title: 'Morning stretch', done: true, date: todayIso() },
    { id: 'task-today-2', title: 'Write journal entry', done: false, date: todayIso() },
    { id: 'task-today-3', title: 'Plan tomorrow focus', done: false, date: todayIso() },
    { id: 'task-soon-1', title: 'Team standup prep', done: false, date: daysFromNowIso(1) },
    { id: 'task-soon-2', title: 'Grocery run', done: false, date: daysFromNowIso(2) },
    { id: 'task-past-1', title: 'Ship dashboard UI', done: true, date: daysAgoIso(1) }
  ]
}

function buildSeedNotes(): JournalNote[] {
  return [
    { id: 'note-inbox-1', body: 'Ideas for the habit tracker — keep it simple.', date: null },
    { id: 'note-inbox-2', body: 'Quote: small steps compound.', date: null },
    { id: 'note-today-1', body: 'Felt focused after the morning walk.', date: todayIso() },
    { id: 'note-today-2', body: 'Sketch calendar chips for days with items.', date: todayIso() },
    { id: 'note-soon-1', body: 'Tomorrow: try a shorter standup agenda.', date: daysFromNowIso(1) },
    { id: 'note-past-1', body: 'Reflection: shipping UI before API was the right call.', date: daysAgoIso(2) }
  ]
}

export function usePlaceholderJournal() {
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

  function addTask(title: string, date: string | null = null) {
    const trimmed = title.trim()
    if (!trimmed) return null
    const task: JournalTask = {
      id: createId(),
      title: trimmed,
      done: false,
      date
    }
    tasks.value = [task, ...tasks.value]
    return task
  }

  function addNote(body: string, date: string | null = null) {
    const trimmed = body.trim()
    if (!trimmed) return null
    const note: JournalNote = {
      id: createId(),
      body: trimmed,
      date
    }
    notes.value = [note, ...notes.value]
    return note
  }

  function updateTask(id: string, title: string) {
    const trimmed = title.trim()
    if (!trimmed) return false
    tasks.value = tasks.value.map(t => (t.id === id ? { ...t, title: trimmed } : t))
    return true
  }

  function updateNote(id: string, body: string) {
    const trimmed = body.trim()
    if (!trimmed) return false
    notes.value = notes.value.map(n => (n.id === id ? { ...n, body: trimmed } : n))
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
    datesWithItems,
    addTask,
    addNote,
    updateTask,
    updateNote,
    deleteTask,
    deleteNote,
    toggleTaskDone
  }
}
