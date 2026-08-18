export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  authStore.checkAuth()
  // Cek apakah route membutuhkan auth
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return navigateTo('/auth/signin')
  }
})
