<script setup>
import Nav from '@/views/components/Nav.vue'
</script>

<template>
  <Nav />
  <div class="container mx-auto p-6 space-y-6">
    <h1 class="text-4xl font-bold text-primary text-center">dAcademy</h1>
    <p class="text-lg text-gray-700 text-center">
      Welcome to dAcademy
    </p>

    <div class="text-center">
      <a
        href="https://github.com/dUnionNet/dAcademy"
        target="_blank"
        class="btn btn-accent mt-4"
      >
        Visit our GitHub Repository
      </a>
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 800px;
}
</style>

