<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUrls } from '../api/urls'
import type { UrlResponse } from '../types'
import UrlTable from '../components/UrlTable.vue'

const urls = ref<UrlResponse[]>([])
const loading = ref(true)
const error = ref('')

const urlCount = ref(0)

onMounted(async () => {
  try {
    urls.value = await getUrls()
    urlCount.value = urls.value.length
  } catch (e: any) {
    error.value = e.message || 'Failed to load data'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <h4 class="mb-4">Dashboard</h4>

    <div v-if="error" class="alert alert-danger">{{ error }}</div>

    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <template v-else>
      <div class="row mb-4">
        <div class="col-md-4 mb-3">
          <div class="card">
            <div class="card-body">
              <div class="text-body-secondary mb-1">Total URLs</div>
              <div class="fs-3 fw-bold">{{ urlCount }}</div>
            </div>
          </div>
        </div>
      </div>

      <h5 class="mb-3">Recent URLs</h5>
      <UrlTable :urls="urls.slice(0, 10)" />

      <router-link v-if="urls.length > 10" to="/urls" class="btn btn-outline-secondary btn-sm mt-2">
        View all URLs →
      </router-link>
    </template>
  </div>
</template>
