<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import Nav from '@/views/components/Nav.vue'

const courses = ref([])
const loading = ref(true)
const error = ref(null)

async function fetchCourses() {
  loading.value = true
  try {
    const res = await axios.get('/api/course/list')
    courses.value = res.data.courses || []
  } catch (err) {
    error.value = err.message || 'Failed to load courses'
  } finally {
    loading.value = false
  }
}

async function scanCourses() {
  try {
    const res = await axios.get('/api/course/scan')
    alert(res.data.message || 'Scan complete.')
    // Refetch courses after scanning
    await fetchCourses()
  } catch (err) {
    alert('Scan failed: ' + (err.message || 'Unknown error'))
  }
}

onMounted(fetchCourses)
</script>

<template>
  <Nav />
  <div class="container mx-auto p-6 space-y-6">
    <h1 class="text-3xl font-bold text-center text-primary mb-6">Courses</h1>
    <button
      @click="scanCourses"
      class="btn btn-accent"
    >
      Scan Courses
    </button>

    <div v-if="loading" class="text-center text-lg text-gray-500">Loading...</div>
    <div v-else-if="error" class="text-center text-lg text-red-500">Error: {{ error }}</div>

    <div v-else>
      <div
        v-for="(course, index) in courses"
        :key="index"
        class="card bg-base-100 shadow-xl p-4 space-y-4"
      >
        <h2 class="text-xl font-semibold text-primary">{{ course.name }}</h2>
        <p class="text-gray-700">{{ course.description }}</p>
        <p class="text-sm text-gray-600">Chapters: {{ course.chapter_count }}</p>
        <RouterLink
          :to="`/c/${course.slug}`"
          class="btn btn-primary w-full"
        >
          Start Course
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 800px;
}
.card {
  border-radius: 12px;
  transition: transform 0.3s ease;
}
.card:hover {
  transform: translateY(-5px);
}
</style>
