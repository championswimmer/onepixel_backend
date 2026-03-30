<script setup lang="ts">
import { ref } from 'vue'
import { createRandomUrl, createCustomUrl } from '../api/urls'
import type { UrlResponse } from '../types'

const longUrl = ref('')
const customShortcode = ref('')
const useCustom = ref(false)
const loading = ref(false)
const error = ref('')
const result = ref<UrlResponse | null>(null)
const copied = ref(false)

async function handleSubmit() {
  error.value = ''
  result.value = null
  loading.value = true

  try {
    if (useCustom.value && customShortcode.value) {
      result.value = await createCustomUrl(customShortcode.value, { long_url: longUrl.value })
    } else {
      result.value = await createRandomUrl({ long_url: longUrl.value })
    }
  } catch (e: any) {
    error.value = e.message || 'Failed to create URL'
  } finally {
    loading.value = false
  }
}

async function copyToClipboard() {
  if (!result.value) return
  try {
    await navigator.clipboard.writeText(result.value.short_url)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    // Fallback: select text
  }
}

function resetForm() {
  longUrl.value = ''
  customShortcode.value = ''
  useCustom.value = false
  result.value = null
  error.value = ''
}
</script>

<template>
  <div>
    <h4 class="mb-4">Create Short URL</h4>

    <div v-if="result" class="card border-success mb-4">
      <div class="card-body">
        <h6 class="card-title text-success">
          <i class="bi bi-check-circle me-1"></i> URL created
        </h6>
        <div class="input-group">
          <input
            type="text"
            class="form-control font-monospace"
            :value="result.short_url"
            readonly
          />
          <button class="btn btn-outline-secondary" @click="copyToClipboard">
            <i :class="copied ? 'bi bi-check-lg' : 'bi bi-clipboard'"></i>
            {{ copied ? 'Copied' : 'Copy' }}
          </button>
        </div>
        <div class="mt-2">
          <span class="text-body-secondary small">→ {{ result.long_url }}</span>
        </div>
        <button class="btn btn-outline-dark btn-sm mt-3" @click="resetForm">
          Create another
        </button>
      </div>
    </div>

    <form v-else @submit.prevent="handleSubmit">
      <div v-if="error" class="alert alert-danger">{{ error }}</div>

      <div class="mb-3">
        <label for="longUrl" class="form-label">Long URL</label>
        <input
          id="longUrl"
          v-model="longUrl"
          type="url"
          class="form-control"
          placeholder="https://example.com/very/long/url"
          required
        />
      </div>

      <div class="mb-3 form-check">
        <input
          id="useCustom"
          v-model="useCustom"
          type="checkbox"
          class="form-check-input"
        />
        <label for="useCustom" class="form-check-label">Use custom shortcode</label>
      </div>

      <div v-if="useCustom" class="mb-3">
        <label for="shortcode" class="form-label">Custom shortcode</label>
        <input
          id="shortcode"
          v-model="customShortcode"
          type="text"
          class="form-control font-monospace"
          placeholder="my-link"
          :required="useCustom"
        />
      </div>

      <button type="submit" class="btn btn-dark" :disabled="loading">
        <span v-if="loading" class="spinner-border spinner-border-sm me-1" role="status"></span>
        Shorten URL
      </button>
    </form>
  </div>
</template>
