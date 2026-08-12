export interface JournalTask {
  id: string
  title: string
  done: boolean
  /** ISO date (YYYY-MM-DD) for Calendar; null for General */
  date: string | null
}

export interface JournalNote {
  id: string
  body: string
  /** ISO date (YYYY-MM-DD) for Calendar; null for General */
  date: string | null
}

export type ItemKind = 'tasks' | 'notes'
