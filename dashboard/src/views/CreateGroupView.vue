<script setup lang="ts">
import { ref } from 'vue'
import { useAuth } from '../composables/useAuth'
import { createUrlGroup } from '../api/urls'
import type { UrlGroupResponse } from '../types'

const { user } = useAuth()

const shortPath = ref('')
const loading = ref(false)
const error = ref('')
const result = ref<UrlGroupResponse | null>(null)

async function handleSubmit() {
  error.value = ''
  result.value = null
  loading.value = true

  try {
    result.value = await createUrlGroup({
      short_path: shortPath.value.trim(),
      creator_id: user.value!.id,
    })
  } catch (e: any) {
    error.value = e.message || 'Failed to create group'
  } finally {
    loading.value = false
  }
}

function resetForm() {
  shortPath.value = ''
  result.value = null
  error.value = ''
}
</script>

<template>
  <div>
    <h4 class="mb-4">Create URL Group</h4>

    <div v-if="result" class="card border-success mb-4" style="max-width: 500px;">
      <div class="card-body">
        <h6 class="card-title text-success">
          <i class="bi bi-check-circle me-1"></i> Group created
        </h6>
        <p class="font-monospace mb-1">{{ result.short_path }}</p>
        <p class="text-body-secondary small mb-3">
          You can now create URLs within this group from the Create URL page.
        </p>
        <div class="d-flex gap-2">
          <button class="btn btn-outline-dark btn-sm" @click="resetForm">
            Create another
          </button>
          <router-link to="/urls/new" class="btn btn-dark btn-sm">
            Create URL in group
          </router-link>
        </div>
      </div>
    </div>

    <form v-else @submit.prevent="handleSubmit" style="max-width: 500px;">
      <div v-if="error" class="alert alert-danger">{{ error }}</div>

      <div class="mb-3">
        <label for="shortPath" class="form-label">Group name</label>
        <input
          id="shortPath"
          v-model="shortPath"
          type="text"
          class="form-control font-monospace"
          placeholder="my-team"
          required
          maxlength="10"
        />
        <div class="form-text">
          Max 10 characters. URLs in this group will be accessible at <code>/group-name/shortcode</code>.
        </div>
      </div>

      <button type="submit" class="btn btn-dark" :disabled="loading">
        <span v-if="loading" class="spinner-border spinner-border-sm me-1" role="status"></span>
        Create Group
      </button>
    </form>
  </div>
</template>
