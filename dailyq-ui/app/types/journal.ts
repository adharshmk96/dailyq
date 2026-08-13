export type TagColor = 'primary' | 'success' | 'warning' | 'info' | 'error' | 'neutral'

export const TAG_COLORS: TagColor[] = ['primary', 'success', 'warning', 'info', 'error', 'neutral']

export interface JournalTag {
  id: string
  name: string
  color: TagColor
}

export interface JournalTask {
  id: string
  title: string
  done: boolean
  /** ISO date (YYYY-MM-DD) for Calendar; null for General */
  date: string | null
  tagIds: string[]
}

export interface JournalNote {
  id: string
  body: string
  /** ISO date (YYYY-MM-DD) for Calendar; null for General */
  date: string | null
  tagIds: string[]
}

export type ItemKind = 'tasks' | 'notes'
