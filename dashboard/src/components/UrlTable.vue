<script setup lang="ts">
import type { UrlResponse } from '../types'

defineProps<{ urls: UrlResponse[] }>()

function extractPath(shortUrl: string): string {
  try {
    const url = new URL(shortUrl)
    return url.pathname.slice(1)
  } catch {
    return shortUrl
  }
}

function extractGroup(path: string): string | null {
  const parts = path.split('/')
  return parts.length > 1 ? parts[0] : null
}

function extractShortcode(path: string): string {
  const parts = path.split('/')
  return parts[parts.length - 1]
}

function detailLink(shortUrl: string): string {
  const path = extractPath(shortUrl)
  return `/urls/${path}`
}

function truncate(str: string, len: number): string {
  return str.length > len ? str.slice(0, len) + '…' : str
}
</script>

<template>
  <div class="table-responsive">
    <table class="table table-hover align-middle">
      <thead>
        <tr>
          <th>Short URL</th>
          <th>Destination</th>
          <th class="text-end" style="width: 80px;"></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="url in urls" :key="url.short_url">
          <td>
            <router-link
              :to="detailLink(url.short_url)"
              class="font-monospace text-decoration-none"
            >
              <template v-if="extractGroup(extractPath(url.short_url))">
                <span class="badge bg-secondary">{{ extractGroup(extractPath(url.short_url)) }}</span>
                <span class="text-body-secondary mx-1">/</span>
              </template>
              {{ extractShortcode(extractPath(url.short_url)) }}
            </router-link>
          </td>
          <td>
            <a
              :href="url.long_url"
              target="_blank"
              rel="noopener"
              class="text-body-secondary text-decoration-none"
              :title="url.long_url"
            >
              {{ truncate(url.long_url, 60) }}
            </a>
          </td>
          <td class="text-end">
            <router-link
              :to="detailLink(url.short_url)"
              class="btn btn-outline-secondary btn-sm"
              title="View details"
            >
              <i class="bi bi-arrow-right"></i>
            </router-link>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
