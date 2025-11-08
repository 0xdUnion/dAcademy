<script setup>
import { RouterLink } from 'vue-router'
</script>

<template>
  <nav class="bg-base-200 p-4 shadow-md rounded-md">
    <div class="container mx-auto flex justify-between items-center">
      <RouterLink to="/" class="text-2xl font-semibold text-primary hover:text-accent">
        Home
      </RouterLink>
      <div class="space-x-4">
        <RouterLink
          to="/about"
          class="text-lg font-medium text-gray-700 hover:text-primary transition-colors"
        >
          About
        </RouterLink>
      </div>
    </div>
  </nav>
</template>

<style scoped>
nav {
  border-radius: 12px; /* Add border-radius for rounded corners */
}
</style>
