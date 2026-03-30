<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUrlInfo, getGroupedUrlInfo } from '../api/urls'
import type { UrlInfoResponse } from '../types'

const props = defineProps<{ shortcode: string; group?: string }>()

const urlInfo = ref<UrlInfoResponse | null>(null)
const loading = ref(true)
const error = ref('')

const displayPath = props.group ? `/${props.group}/${props.shortcode}` : `/${props.shortcode}`

onMounted(async () => {
  try {
    if (props.group) {
      urlInfo.value = await getGroupedUrlInfo(props.group, props.shortcode)
    } else {
      urlInfo.value = await getUrlInfo(props.shortcode)
    }
  } catch (e: any) {
    error.value = e.message || 'Failed to load URL info'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div class="mb-4">
      <router-link to="/urls" class="text-body-secondary text-decoration-none">
        ← Back to URLs
      </router-link>
    </div>

    <h4 class="mb-4 font-monospace">{{ displayPath }}</h4>

    <div v-if="error" class="alert alert-danger">{{ error }}</div>

    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <template v-else-if="urlInfo">
      <div class="row">
        <div class="col-md-6 mb-3">
          <div class="card">
            <div class="card-body">
              <div class="text-body-secondary mb-1">Destination</div>
              <a :href="urlInfo.long_url" target="_blank" rel="noopener" class="text-break">
                {{ urlInfo.long_url }}
              </a>
            </div>
          </div>
        </div>
        <div class="col-md-3 mb-3">
          <div class="card">
            <div class="card-body">
              <div class="text-body-secondary mb-1">Total Hits</div>
              <div class="fs-3 fw-bold">{{ urlInfo.hit_count }}</div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
