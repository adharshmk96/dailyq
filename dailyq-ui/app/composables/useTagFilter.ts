export function useTagFilter() {
  const route = useRoute()
  const journal = usePlaceholderJournal()

  const activeTagId = computed(() => {
    const param = route.query.tag
    if (typeof param !== 'string' || !param) return null
    const exists = journal.tags.value.some(t => t.id === param)
    return exists ? param : null
  })

  watch(activeTagId, (value) => {
    const param = route.query.tag
    if (typeof param === 'string' && param && !value) {
      const { tag: _tag, ...rest } = route.query
      navigateTo({ path: route.path, query: rest }, { replace: true })
    }
  }, { immediate: true })

  function setActiveTag(id: string | null) {
    const query = { ...route.query }
    if (id) {
      query.tag = id
    } else {
      delete query.tag
    }
    navigateTo({ path: route.path, query }, { replace: true })
  }

  function toggleTag(id: string) {
    setActiveTag(activeTagId.value === id ? null : id)
  }

  return {
    activeTagId,
    setActiveTag,
    toggleTag
  }
}
