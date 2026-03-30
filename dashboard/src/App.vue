<script setup lang="ts">
import NavBar from './components/NavBar.vue'
</script>

<template>
  <NavBar />
  <main class="container pb-4">
    <router-view />
  </main>
</template>
