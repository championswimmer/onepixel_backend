import { ref, computed } from 'vue'
import type { UserResponse, LoginRequest } from '../types'
import { loginUser } from '../api/users'

const token = ref<string | null>(localStorage.getItem('token'))
const user = ref<{ id: number; email: string } | null>(
  JSON.parse(localStorage.getItem('user') || 'null')
)
const adminKey = ref<string | null>(localStorage.getItem('adminKey'))
const adminKeyEnabled = ref<boolean>(localStorage.getItem('adminKeyEnabled') !== 'false')

export function useAuth() {
  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => !!adminKey.value)
  const isAdminKeyActive = computed(() => !!adminKey.value && adminKeyEnabled.value)

  async function login(credentials: LoginRequest): Promise<void> {
    const response: UserResponse = await loginUser(credentials)
    token.value = response.token
    user.value = { id: response.id, email: response.email }
    localStorage.setItem('token', response.token)
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  function logout(): void {
    token.value = null
    user.value = null
    adminKey.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('adminKey')
  }

  function setAdminKey(key: string): void {
    adminKey.value = key
    adminKeyEnabled.value = true
    localStorage.setItem('adminKey', key)
    localStorage.setItem('adminKeyEnabled', 'true')
  }

  function clearAdminKey(): void {
    adminKey.value = null
    adminKeyEnabled.value = true
    localStorage.removeItem('adminKey')
    localStorage.removeItem('adminKeyEnabled')
  }

  function toggleAdminKey(enabled: boolean): void {
    adminKeyEnabled.value = enabled
    localStorage.setItem('adminKeyEnabled', enabled ? 'true' : 'false')
  }

  return {
    token,
    user,
    adminKey,
    adminKeyEnabled,
    isAuthenticated,
    isAdmin,
    isAdminKeyActive,
    login,
    logout,
    setAdminKey,
    clearAdminKey,
    toggleAdminKey,
  }
}
