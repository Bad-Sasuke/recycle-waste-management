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
const isUploadingImage = ref(false)
const uploadProgress = ref(0)
const imageFile = ref(null)

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

// Image upload functions
const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    validateAndUploadFile(file)
  }
}

const validateAndUploadFile = (file) => {
  // Validate file type
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    error.value = 'Please select a valid image file (JPG, PNG, GIF)'
    return
  }

  // Validate file size (max 10MB)
  const maxSize = 10 * 1024 * 1024 // 10MB
  if (file.size > maxSize) {
    error.value = 'File size must be less than 10MB'
    return
  }

  imageFile.value = file
  uploadImage()
}

const uploadImage = async () => {
  if (!imageFile.value) return

  isUploadingImage.value = true
  uploadProgress.value = 0
  error.value = ''
  success.value = ''

  try {
    const token = getCookie('token')
    if (!token) {
      error.value = 'Please login to update profile image'
      return
    }

    if (!apiUrl) {
      error.value = 'API configuration error'
      return
    }

    // Create FormData for file upload
    const formData = new FormData()
    formData.append('image', imageFile.value)

    // Upload with progress tracking
    const xhr = new XMLHttpRequest()

    // Track upload progress
    xhr.upload.addEventListener('progress', (event) => {
      if (event.lengthComputable) {
        uploadProgress.value = Math.round((event.loaded / event.total) * 100)
      }
    })

    // Handle response
    const uploadPromise = new Promise((resolve, reject) => {
      xhr.onload = () => {
        if (xhr.status === 200) {
          resolve(JSON.parse(xhr.responseText))
        } else {
          reject(new Error(`Upload failed: ${xhr.status} ${xhr.statusText}`))
        }
      }
      xhr.onerror = () => reject(new Error('Network error'))
    })

    // Configure and send request
    xhr.open('POST', `${apiUrl}/api/user/update-image-profile`)
    xhr.setRequestHeader('Authorization', `Bearer ${token}`)
    xhr.send(formData)

    // Wait for upload to complete
    const result = await uploadPromise

    // Update user data with new image URL
    userData.value.image_url = result.data.image_url
    success.value = 'Profile image updated successfully!'

    // Reset image error state
    imageLoadError.value = false

  } catch (err) {
    console.error('Image upload error:', err)
    error.value = `Failed to upload image: ${err.message}`
  } finally {
    isUploadingImage.value = false
    uploadProgress.value = 0
    imageFile.value = null

    // Reset file input
    const fileInput = document.getElementById('imageUpload')
    if (fileInput) fileInput.value = ''
  }
}

const triggerFileUpload = () => {
  const fileInput = document.getElementById('imageUpload')
  if (fileInput) {
    fileInput.click()
  }
}

// Drag and drop functionality
const isDragging = ref(false)

const handleDragOver = (event) => {
  event.preventDefault()
  isDragging.value = true
}

const handleDragLeave = (event) => {
  event.preventDefault()
  isDragging.value = false
}

const handleDrop = (event) => {
  event.preventDefault()
  isDragging.value = false

  const files = event.dataTransfer.files
  if (files.length > 0) {
    const file = files[0]
    validateAndUploadFile(file)
  }
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
            <div class="avatar relative group">
              <div
                class="w-16 h-16 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden relative transition-all duration-200"
                :class="{ 'ring-2 ring-blue-500 ring-offset-2': isDragging }"
                @dragover="handleDragOver"
                @dragleave="handleDragLeave"
                @drop="handleDrop"
              >
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

                <!-- Upload overlay -->
                <div class="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200 cursor-pointer rounded-full"
                     @click="triggerFileUpload"
                     v-if="!isUploadingImage">
                  <IconEdit size="24" class="text-white" />
                </div>

                <!-- Upload progress -->
                <div v-if="isUploadingImage" class="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center rounded-full">
                  <div class="text-white text-xs font-semibold">{{ uploadProgress }}%</div>
                </div>

                <!-- Drag overlay -->
                <div v-if="isDragging" class="absolute inset-0 bg-blue-500 bg-opacity-75 flex items-center justify-center rounded-full">
                  <div class="text-white text-xs font-semibold">Drop here</div>
                </div>
              </div>

              <!-- Hidden file input -->
              <input
                id="imageUpload"
                type="file"
                accept="image/*"
                @change="handleFileSelect"
                class="hidden"
              />
            </div>
            <div>
              <h1 class="text-3xl font-bold text-gray-800">{{ $t('Profile.title') }}</h1>
              <p class="text-gray-600">{{ $t('Profile.subtitle') }}</p>
              <p class="text-sm text-gray-500 mt-1">{{ $t('Profile.clickOrDragToUpload') }}</p>
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