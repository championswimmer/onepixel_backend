<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const route = useRoute()
const { user, isAdmin, isAdminKeyActive, adminKeyEnabled, logout, toggleAdminKey } = useAuth()

function handleLogout() {
  logout()
  window.location.hash = '#/login'
}
</script>

<template>
  <nav v-if="route.name !== 'login'" class="navbar navbar-expand-md bg-body-tertiary border-bottom mb-4">
    <div class="container">
      <router-link to="/" class="navbar-brand font-monospace fw-bold">1px.li</router-link>

      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarNav"
      >
        <span class="navbar-toggler-icon"></span>
      </button>

      <div id="navbarNav" class="collapse navbar-collapse">
        <ul class="navbar-nav me-auto">
          <li class="nav-item">
            <router-link to="/" class="nav-link" active-class="active" :exact="true">
              <i class="bi bi-speedometer2 me-1"></i> Dashboard
            </router-link>
          </li>
          <li class="nav-item">
            <router-link to="/urls" class="nav-link" active-class="active">
              <i class="bi bi-link-45deg me-1"></i> URLs
            </router-link>
          </li>
          <li class="nav-item">
            <router-link to="/urls/new" class="nav-link" active-class="active">
              <i class="bi bi-plus-circle me-1"></i> Create
            </router-link>
          </li>
          <li v-if="isAdmin" class="nav-item">
            <router-link to="/groups/new" class="nav-link" active-class="active">
              <i class="bi bi-folder-plus me-1"></i> New Group
            </router-link>
          </li>
        </ul>

        <ul class="navbar-nav">
          <li v-if="isAdmin" class="nav-item d-flex align-items-center me-3">
            <div class="form-check form-switch mb-0">
              <input
                id="adminKeyToggle"
                class="form-check-input"
                type="checkbox"
                role="switch"
                :checked="adminKeyEnabled"
                @change="toggleAdminKey(!adminKeyEnabled)"
              />
              <label
                class="form-check-label small"
                for="adminKeyToggle"
                :class="isAdminKeyActive ? 'text-warning' : 'text-body-secondary'"
              >
                <i class="bi bi-key-fill me-1"></i>Admin Key
              </label>
            </div>
          </li>
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
              <i class="bi bi-person-circle me-1"></i> {{ user?.email }}
              <span v-if="isAdmin" class="badge bg-warning text-dark ms-1">Admin</span>
            </a>
            <ul class="dropdown-menu dropdown-menu-end">
              <li>
                <router-link to="/account" class="dropdown-item">
                  <i class="bi bi-gear me-1"></i> Account
                </router-link>
              </li>
              <li><hr class="dropdown-divider" /></li>
              <li>
                <a class="dropdown-item" href="#" @click.prevent="handleLogout">
                  <i class="bi bi-box-arrow-right me-1"></i> Sign out
                </a>
              </li>
            </ul>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>
