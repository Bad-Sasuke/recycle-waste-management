<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { IconUserCircle, IconEdit, IconCheck, IconX } from '@tabler/icons-vue'
import { useUsersStore } from '@/stores/users'
import { getCookie } from '@/stores/cookie'

const usersStore = useUsersStore()

// User data
const userData = ref({
  user_id: '',
  username: '',
  email: '',
  image_url: '',
  created_at: '',
  last_login: ''
})

// Form states
const isEditing = ref(false)
const isLoading = ref(false)
const error = ref('')
const success = ref('')

// Form data for editing
const editForm = ref({
  username: '',
  email: ''
})

// Image handling
const imageLoadError = ref(false)

// API URL for profile endpoints
const apiUrl = import.meta.env.VITE_WEB_API


const handleImageError = (event) => {
  console.error('Image load error:', {
    src: event.target.src,
    error: event,
    naturalWidth: event.target.naturalWidth,
    naturalHeight: event.target.naturalHeight
  })
  imageLoadError.value = true
  // Hide the broken image
  event.target.style.display = 'none'
}

const handleImageLoad = (event) => {
  console.log('Image loaded successfully:', {
    src: event.target.src,
    naturalWidth: event.target.naturalWidth,
    naturalHeight: event.target.naturalHeight
  })
  imageLoadError.value = false
}

// Fetch user profile data
const fetchUserProfile = async () => {
  isLoading.value = true
  error.value = ''

  try {
    const token = getCookie('token')
    if (!token) {
      error.value = 'Please login to view your profile'
      return
    }

    if (!apiUrl) {
      error.value = 'API configuration error. Please check environment settings.'
      return
    }

    const response = await fetch(`${apiUrl}/api/user/profile`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (response.ok) {
      const result = await response.json()
      userData.value = result.data
      editForm.value = {
        username: result.data.username,
        email: result.data.email
      }
      // Reset image error state when new data is loaded
      imageLoadError.value = false
    } else {
      error.value = `Failed to fetch profile data: ${response.status} ${response.statusText}`
    }
  } catch (err) {
    error.value = `Network error: ${err.message || 'Cannot connect to server'}`
  } finally {
    isLoading.value = false
  }
}

// Update user profile
const updateProfile = async () => {
  isLoading.value = true
  error.value = ''
  success.value = ''

  try {
    const token = getCookie('token')

    if (!apiUrl) {
      error.value = 'API configuration error. Please check environment settings.'
      return
    }

    const response = await fetch(`${apiUrl}/api/user/profile`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(editForm.value)
    })

    if (response.ok) {
      success.value = 'Profile updated successfully!'
      userData.value.username = editForm.value.username
      userData.value.email = editForm.value.email
      isEditing.value = false
    } else {
      error.value = `Failed to update profile: ${response.status} ${response.statusText}`
    }
  } catch (err) {
    error.value = `Network error: ${err.message || 'Cannot connect to server'}`
  } finally {
    isLoading.value = false
  }
}

// Cancel editing
const cancelEdit = () => {
  isEditing.value = false
  editForm.value = {
    username: userData.value.username,
    email: userData.value.email
  }
  error.value = ''
  success.value = ''
}

// Format date
const formatDate = (dateString: string) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  if (!usersStore.isLogin) {
    error.value = 'Please login to view your profile'
    return
  }
  fetchUserProfile()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="container mx-auto px-4 max-w-4xl">
      <!-- Header -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div class="avatar">
              <div class="w-16 h-16 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden">
                <img
                  v-if="userData.image_url && userData.image_url.trim() !== '' && !imageLoadError"
                  :src="userData.image_url"
                  :alt="userData.username || 'Profile'"
                  class="w-full h-full object-cover"
                  @error="handleImageError"
                  @load="handleImageLoad"
                />
                <IconUserCircle
                  v-else
                  size="64"
                  class="text-gray-400"
                />
              </div>
            </div>
            <div>
              <h1 class="text-3xl font-bold text-gray-800">{{ $t('Profile.title') }}</h1>
              <p class="text-gray-600">{{ $t('Profile.subtitle') }}</p>
            </div>
          </div>

          <!-- Edit Button -->
          <button
            v-if="!isEditing"
            @click="isEditing = true"
            class="btn btn-outline btn-primary"
          >
            <IconEdit size="20" />
            {{ $t('Profile.edit') }}
          </button>
        </div>
      </div>

      <!-- Alert Messages -->
      <div v-if="error" class="alert alert-error mb-6">
        <IconX size="20" />
        <span>{{ error }}</span>
      </div>

      <div v-if="success" class="alert alert-success mb-6">
        <IconCheck size="20" />
        <span>{{ success }}</span>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex justify-center items-center py-12">
        <span class="loading loading-spinner loading-lg text-primary"></span>
      </div>

      <!-- Profile Content -->
      <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Basic Information Card -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">{{ $t('Profile.basicInfo') }}</h2>

          <div class="space-y-4">
            <!-- Username -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('Profile.username') }}
              </label>
              <input
                v-if="isEditing"
                v-model="editForm.username"
                type="text"
                class="input input-bordered w-full"
                :placeholder="$t('Profile.usernamePlaceholder')"
              />
              <p v-else class="text-gray-900 font-medium">{{ userData.username || 'N/A' }}</p>
            </div>

            <!-- Email -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('Profile.email') }}
              </label>
              <input
                v-if="isEditing"
                v-model="editForm.email"
                type="email"
                class="input input-bordered w-full"
                :placeholder="$t('Profile.emailPlaceholder')"
              />
              <p v-else class="text-gray-900 font-medium">{{ userData.email || 'N/A' }}</p>
            </div>

            <!-- User ID -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('Profile.userId') }}
              </label>
              <p class="text-gray-900 font-medium font-mono text-sm">{{ userData.user_id || 'N/A' }}</p>
            </div>
          </div>

          <!-- Edit Actions -->
          <div v-if="isEditing" class="flex gap-2 mt-6">
            <button
              @click="updateProfile"
              :disabled="isLoading"
              class="btn btn-primary"
            >
              <IconCheck size="20" />
              {{ $t('Profile.save') }}
            </button>
            <button
              @click="cancelEdit"
              class="btn btn-outline"
            >
              <IconX size="20" />
              {{ $t('Profile.cancel') }}
            </button>
          </div>
        </div>

        <!-- Account Information Card -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">{{ $t('Profile.accountInfo') }}</h2>

          <div class="space-y-4">
            <!-- Created At -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('Profile.memberSince') }}
              </label>
              <p class="text-gray-900">{{ formatDate(userData.created_at) }}</p>
            </div>

            <!-- Last Login -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('Profile.lastLogin') }}
              </label>
              <p class="text-gray-900">{{ formatDate(userData.last_login) }}</p>
            </div>

            <!-- Profile Picture -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ $t('Profile.profilePicture') }}
              </label>
              <div class="flex items-center gap-4">
                <div class="avatar">
                  <div class="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden">
                    <img
                      v-if="userData.image_url && userData.image_url.trim() !== '' && !imageLoadError"
                      :src="userData.image_url"
                      :alt="userData.username || 'Profile'"
                      class="w-full h-full object-cover"
                      @error="handleImageError"
                      @load="handleImageLoad"
                    />
                    <IconUserCircle
                      v-else
                      size="48"
                      class="text-gray-400"
                    />
                  </div>
                </div>
                <p class="text-sm text-gray-600">
                  {{ userData.image_url ? $t('Profile.profilePictureSet') : $t('Profile.profilePictureNotSet') }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>


    </div>
  </div>
</template>

<style scoped>
.alert {
  @apply flex items-center gap-3 p-4 rounded-lg;
}

.alert-error {
  @apply bg-red-50 border border-red-200 text-red-800;
}

.alert-success {
  @apply bg-green-50 border border-green-200 text-green-800;
}
</style>