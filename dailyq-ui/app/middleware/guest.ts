/** Keeps signed-in users out of the login/register pages. */
export default defineNuxtRouteMiddleware(() => {
  const { isAuthenticated } = useAuth()

  if (isAuthenticated.value) {
    return navigateTo('/dashboard/overview')
  }
})
