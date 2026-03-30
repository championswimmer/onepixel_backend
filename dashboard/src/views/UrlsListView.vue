<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUrls } from '../api/urls'
import type { UrlResponse } from '../types'
import UrlTable from '../components/UrlTable.vue'

const urls = ref<UrlResponse[]>([])
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    urls.value = await getUrls()
  } catch (e: any) {
    error.value = e.message || 'Failed to load URLs'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h4 class="mb-0">My URLs</h4>
      <router-link to="/urls/new" class="btn btn-dark btn-sm">
        <i class="bi bi-plus-lg me-1"></i> New URL
      </router-link>
    </div>

    <div v-if="error" class="alert alert-danger">{{ error }}</div>

    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <template v-else>
      <p v-if="urls.length === 0" class="text-body-secondary">
        No URLs yet. <router-link to="/urls/new">Create your first short URL</router-link>.
      </p>
      <UrlTable v-else :urls="urls" />
    </template>
  </div>
</template>
