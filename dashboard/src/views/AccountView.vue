<script setup lang="ts">
import { ref } from 'vue'
import { useAuth } from '../composables/useAuth'
import { updateUser } from '../api/users'

const { user } = useAuth()

const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)

async function handleSubmit() {
  error.value = ''
  success.value = false

  if (newPassword.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }

  if (newPassword.value.length < 6) {
    error.value = 'Password must be at least 6 characters'
    return
  }

  loading.value = true
  try {
    await updateUser(user.value!.id, { password: newPassword.value })
    success.value = true
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (e: any) {
    error.value = e.message || 'Failed to update password'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <h4 class="mb-4">Account</h4>

    <div class="card" style="max-width: 500px;">
      <div class="card-body">
        <div class="mb-3">
          <label class="form-label text-body-secondary">Email</label>
          <div class="fw-bold">{{ user?.email }}</div>
        </div>

        <hr />

        <h6>Change Password</h6>

        <div v-if="error" class="alert alert-danger py-2">{{ error }}</div>
        <div v-if="success" class="alert alert-success py-2">Password updated.</div>

        <form @submit.prevent="handleSubmit">
          <div class="mb-3">
            <label for="newPassword" class="form-label">New password</label>
            <input
              id="newPassword"
              v-model="newPassword"
              type="password"
              class="form-control"
              required
              minlength="6"
            />
          </div>

          <div class="mb-3">
            <label for="confirmPassword" class="form-label">Confirm password</label>
            <input
              id="confirmPassword"
              v-model="confirmPassword"
              type="password"
              class="form-control"
              required
            />
          </div>

          <button type="submit" class="btn btn-dark" :disabled="loading">
            <span v-if="loading" class="spinner-border spinner-border-sm me-1" role="status"></span>
            Update password
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
