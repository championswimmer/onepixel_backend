import { ref, computed } from 'vue'
import type { UserResponse, LoginRequest } from '../types'
import { loginUser } from '../api/users'

const token = ref<string | null>(localStorage.getItem('token'))
const user = ref<{ id: number; email: string } | null>(
  JSON.parse(localStorage.getItem('user') || 'null')
)

export function useAuth() {
  const isAuthenticated = computed(() => !!token.value)

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
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    token,
    user,
    isAuthenticated,
    login,
    logout,
  }
}
