
<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <main class="flex-1">
      <RouterView />
    </main>

    <footer class="footer p-6 bg-base-200 text-base-content border-t">
      <div class="w-full max-w-6xl mx-auto flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div class="flex items-center space-x-3">
          <div class="avatar">
            <div class="w-10 rounded-full bg-primary text-white flex items-center justify-center">dA</div>
          </div>
          <div>
            <div class="font-bold text-lg">dAcademy</div>
            <div class="text-sm opacity-70">Learn. Build. Share.</div>
          </div>
        </div>

        <nav class="flex items-center gap-4">
          <RouterLink to="/" class="link link-hover">Home</RouterLink>
          <RouterLink to="/about" class="link link-hover">About</RouterLink>
        </nav>


      </div>
    </footer>
  </div>
</template>

<style scoped>
.footer nav { flex-wrap: wrap; }
@media (max-width: 640px) {
  .footer { padding: 1.25rem; }
}
</style>

