/**
 * Guards the signed-in area. Redirects to /login and remembers where the user
 * was headed so they land back there after signing in.
 */
export default defineNuxtRouteMiddleware(async (to) => {
  const { isAuthenticated, user, fetchMe, clearSession } = useAuth()

  if (!isAuthenticated.value) {
    return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
  }

  if (!user.value) {
    try {
      const me = await fetchMe()
      if (!me) {
        return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
      }
    } catch {
      clearSession()
      return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
    }
  }
})
