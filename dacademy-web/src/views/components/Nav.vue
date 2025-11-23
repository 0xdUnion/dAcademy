<script setup>
import { ref, reactive, onMounted } from 'vue'
import { RouterLink } from 'vue-router'

const isLoginMode = ref(true)
const isLoading = ref(false)
const errorMessage = ref('')
const currentUser = ref(null)

const formData = reactive({
  username: '',
  password: ''
})

async function sha256(source) {
  const sourceBytes = new TextEncoder().encode(source);
  const digest = await crypto.subtle.digest("SHA-256", sourceBytes);
  return [...new Uint8Array(digest)].map(x => x.toString(16).padStart(2, '0')).join("");
}

const getCookie = (name) => {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
}

onMounted(async () => {
  const session = getCookie('d_session')
  if (session) {
    currentUser.value = "user"
    const response = await fetch("/api/me", {
      method: 'GET',
    })
    const result = await response.json()
    if (!response.ok) {
      logout()

      throw new Error('Failed to fetch account info')
    }
    currentUser.value = result.username

  }
})

const handleSubmit = async () => {
  errorMessage.value = ''

  if (!formData.username || !formData.password) {
    errorMessage.value = "Fill all fields please"
    return
  }

  try {
    isLoading.value = true

    // sha256 encoding
    const hashedPassword = await sha256(formData.password)

    const data = new FormData()
    data.append('username', formData.username)
    data.append('password', hashedPassword)

    const url = isLoginMode.value ? '/api/auth/login' : '/api/auth/signup'

    const response = await fetch(url, {
      method: 'POST',
      body: data
    })

    const result = await response.json()

    if (!response.ok) {
      throw new Error(result.error || (isLoginMode.value ? 'Failed to login' : 'Failed to sign up'))
    }


    if (isLoginMode.value) {
      document.cookie = `d_session=${result.session}; path=/; max-age=8640000; SameSite=Strict`

      currentUser.value = result.username
      alert(`Welcome back, ${result.username}`)
      document.getElementById('auth_modal').close()
      location.reload();

    } else {
      alert('Success')
      isLoginMode.value = true // Switch to login mode
      formData.password = ''   // Clear password field
    }

  } catch (err) {
    errorMessage.value = err.message
  } finally {
    isLoading.value = false
  }
}

const openAuthModal = () => {
  errorMessage.value = ''
  document.getElementById('auth_modal').showModal()
}

const logout = () => {
  document.cookie = "d_session=; path=/; max-age=0";
  location.reload();
}
</script>

<template>
  <nav class="bg-base-200 p-4 shadow-md rounded-md">
    <div class="container mx-auto flex justify-between items-center">
      <RouterLink to="/" class="text-2xl font-semibold text-primary hover:text-accent">
        dAcademy
      </RouterLink>

      <div class="space-x-4">
        <button v-if="!currentUser" class="btn btn-sm btn-primary" @click="openAuthModal">
          Login / Sign up
        </button>

        <div v-else class="flex items-center gap-4">
          <span class="font-bold text-primary">{{ currentUser }}</span>
          <RouterLink to="/account" class="btn btn-sm btn-outline">My Account</RouterLink>
        </div>
      </div>
    </div>
  </nav>

  <dialog id="auth_modal" class="modal">
    <div class="modal-box">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
      </form>

      <div class="flex flex-col items-center gap-2 mb-6">
        <h3 class="font-bold text-2xl">{{ isLoginMode ? 'Welcome back' : 'Join now' }}</h3>
        <div class="text-sm">
          {{ isLoginMode ? "Don't have an account?" : "Already have an account?" }}
          <a class="link link-primary" @click="isLoginMode = !isLoginMode; errorMessage=''">
          {{ isLoginMode ? "Sign Up" : "Sign In" }}
          </a>
        </div>
      </div>

      <div v-if="errorMessage" role="alert" class="alert alert-error text-sm py-2 mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
        <span>{{ errorMessage }}</span>
      </div>

      <div class="flex flex-col gap-4">
        <div class="form-control">
          <label class="label"><span class="label-text">User Name</span></label>
          <input
            v-model="formData.username"
            type="text"
            class="input input-bordered w-full"
            maxlength="8"
          />
        </div>

        <div class="form-control">
          <label class="label"><span class="label-text">Password</span></label>
          <input
            v-model="formData.password"
            type="password"
            class="input input-bordered w-full"
            @keyup.enter="handleSubmit"
          />
        </div>

        <button class="btn btn-primary w-full mt-4" @click="handleSubmit" :disabled="isLoading">
          <span v-if="isLoading" class="loading loading-spinner loading-xs"></span>
          {{ isLoginMode ? 'Login' : 'Register' }}
        </button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
</template>

<style scoped>
nav { border-radius: 12px; }
</style>
