<script setup>
import { useRoute } from 'vue-router'
import axios from 'axios'
import { onMounted, ref } from 'vue'

import Nav from '@/views/components/Nav.vue'

const route = useRoute()
const courseSlug = route.params.slug

const loading = ref(true)
const error = ref(null)
const courseDetail = ref([])
const chapters = ref([])

async function fetchCourseDetail() {
  loading.value = true
  try {
    const res = await axios.get(`/api/course/${courseSlug}`)
    courseDetail.value = res.data.course || []
    chapters.value = res.data.chapters || []
  } catch (err) {
    error.value = err.message || 'Failed to load course detail'
  } finally {
    loading.value = false
  }
}

onMounted(fetchCourseDetail)
</script>

<template>
  <Nav />
  <div class="container mx-auto p-6 space-y-6">
    <div v-if="loading" class="text-center text-lg text-gray-500">Loading...</div>
    <div v-else-if="error" class="text-center text-lg text-red-500">Error: {{ error }}</div>
    <div v-else>
      <div class="card bg-base-100 shadow-xl p-6 mb-6">
        <h2 class="text-3xl font-bold text-primary">{{ courseDetail.name }}</h2>
        <p class="text-lg text-gray-700 mt-2">{{ courseDetail.description }}</p>
        <p class="text-sm text-gray-600 mt-4">Chapters: {{ courseDetail.chapter_count }}</p>
      </div>

      <div v-for="(chapter, index) in chapters" :key="index" class="card bg-base-100 shadow-md p-4 mb-4">
        <h3 class="text-xl font-semibold text-primary">{{ chapter.title }}</h3>
        <p class="text-gray-600">Chapter ID: {{ chapter.id }}</p>
        <RouterLink
          :to="`/c/${courseDetail.slug}/chapter/${chapter.id}`"
          class="btn btn-primary"
        >
          Enter
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 1024px;
}

.card {
  border-radius: 12px;
  transition: transform 0.3s ease;
}

.card:hover {
  transform: translateY(-5px);
}
</style>
