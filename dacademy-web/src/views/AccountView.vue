<script setup>
import axios from 'axios'
import { ref, onMounted } from 'vue'
import Nav from '@/views/components/Nav.vue'

const loading = ref(true)
const error = ref(null)
const username = ref('')

const logout = () => {
  document.cookie = "d_session=; path=/; max-age=0";
  location.reload();
}

onMounted(async () => {
  try {
    const res = await axios.get('/api/me')
    username.value = res.data.username || ''
  } catch (e) {
    error.value = 'Failed to load user info'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <Nav />

  <div class="p-6 flex justify-center">
    <div class="card w-full max-w-md bg-base-200 shadow-xl">
      <div class="card-body">
        <h2 class="card-title flex items-center gap-2">
          <span class="text-xl">User Center</span>
          <span>👤</span>
        </h2>

        <div v-if="loading" class="mt-4">
          <span class="loading loading-spinner"></span>
        </div>

        <div v-else-if="error" class="text-error">{{ error }}</div>

        <div v-else class="space-y-4">
          <div class="flex items-center justify-between p-4 bg-base-100 rounded-xl shadow">
            <span class="font-semibold">Username</span>
            <span class="badge badge-primary">{{ username }}</span>
          </div>

          <button class="btn btn-warning w-full" @click="logout">Logout</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
