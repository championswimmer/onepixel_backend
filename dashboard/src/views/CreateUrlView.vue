<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { createRandomUrl, createCustomUrl, createGroupedRandomUrl, createGroupedCustomUrl, getUrlGroups } from '../api/urls'
import type { UrlResponse, UrlGroupResponse } from '../types'

const longUrl = ref('')
const customShortcode = ref('')
const useCustom = ref(false)
const groupName = ref('')
const useGroup = ref(false)
const loading = ref(false)
const error = ref('')
const result = ref<UrlResponse | null>(null)
const copied = ref(false)

const availableGroups = ref<UrlGroupResponse[]>([])
const groupSearch = ref('')
const showGroupDropdown = ref(false)

const filteredGroups = computed(() => {
  const query = groupSearch.value.toLowerCase()
  if (!query) return availableGroups.value
  return availableGroups.value.filter(g =>
    g.short_path.toLowerCase().includes(query)
  )
})

async function fetchGroups() {
  try {
    availableGroups.value = await getUrlGroups()
  } catch {
    availableGroups.value = []
  }
}

function selectGroup(group: UrlGroupResponse) {
  groupName.value = group.short_path
  groupSearch.value = group.short_path
  showGroupDropdown.value = false
}

function onGroupInput() {
  groupName.value = groupSearch.value
  showGroupDropdown.value = true
}

function onGroupFocus() {
  showGroupDropdown.value = true
}

function onGroupBlur() {
  // Delay to allow click on dropdown item
  setTimeout(() => { showGroupDropdown.value = false }, 200)
}

watch(useGroup, (val) => {
  if (val && availableGroups.value.length === 0) {
    fetchGroups()
  }
})

onMounted(() => {
  fetchGroups()
})

async function handleSubmit() {
  error.value = ''
  result.value = null
  loading.value = true

  const payload = { long_url: longUrl.value }
  const group = groupName.value.trim()
  const shortcode = customShortcode.value.trim()

  try {
    if (useGroup.value && group) {
      if (useCustom.value && shortcode) {
        result.value = await createGroupedCustomUrl(group, shortcode, payload)
      } else {
        result.value = await createGroupedRandomUrl(group, payload)
      }
    } else {
      if (useCustom.value && shortcode) {
        result.value = await createCustomUrl(shortcode, payload)
      } else {
        result.value = await createRandomUrl(payload)
      }
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
  groupName.value = ''
  groupSearch.value = ''
  useGroup.value = false
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
          id="useGroup"
          v-model="useGroup"
          type="checkbox"
          class="form-check-input"
        />
        <label for="useGroup" class="form-check-label">Add to a URL group</label>
      </div>

      <div v-if="useGroup" class="mb-3">
        <label for="groupName" class="form-label">Group name</label>
        <div class="position-relative">
          <input
            id="groupName"
            v-model="groupSearch"
            type="text"
            class="form-control font-monospace"
            placeholder="Search or type group name..."
            autocomplete="off"
            :required="useGroup"
            @input="onGroupInput"
            @focus="onGroupFocus"
            @blur="onGroupBlur"
          />
          <ul
            v-if="showGroupDropdown && filteredGroups.length > 0"
            class="list-group position-absolute w-100 shadow-sm"
            style="z-index: 10; max-height: 200px; overflow-y: auto;"
          >
            <li
              v-for="group in filteredGroups"
              :key="group.short_path"
              class="list-group-item list-group-item-action font-monospace py-2"
              style="cursor: pointer;"
              @mousedown.prevent="selectGroup(group)"
            >
              {{ group.short_path }}
            </li>
          </ul>
          <div v-if="showGroupDropdown && groupSearch && filteredGroups.length === 0" class="position-absolute w-100 shadow-sm" style="z-index: 10;">
            <div class="list-group-item text-body-secondary small">No matching groups found</div>
          </div>
        </div>
        <div class="form-text">Select a group you own. The short URL will be <code>group/shortcode</code>.</div>
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
