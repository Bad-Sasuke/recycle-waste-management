<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { IconUserCircle, IconEdit, IconCheck, IconX } from '@tabler/icons-vue'
import MiniAvatar3D from '@/components/MiniAvatar3D.vue'
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
const imageFile = ref<File | null>(null)

// API URL for profile endpoints
const apiUrl = import.meta.env.VITE_WEB_API


const handleImageError = (event: Event) => {
  const target = event.target as HTMLImageElement
  console.error('Image load error:', {
    src: target.src,
    error: event,
    naturalWidth: target.naturalWidth,
    naturalHeight: target.naturalHeight
  })
  imageLoadError.value = true
  // Hide the broken image
  target.style.display = 'none'
}

const handleImageLoad = (event: Event) => {
  const target = event.target as HTMLImageElement
  console.log('Image loaded successfully:', {
    src: target.src,
    naturalWidth: target.naturalWidth,
    naturalHeight: target.naturalHeight
  })
  imageLoadError.value = false
}

// Image upload functions
const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    validateAndUploadFile(file)
  }
}

const validateAndUploadFile = (file: File) => {
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
    // Wait for upload to complete
    const result = await uploadPromise as { data: { image_url: string } }

    // Update user data with new image URL
    userData.value.image_url = result.data.image_url
    success.value = 'Profile image updated successfully!'

    // Reset image error state
    imageLoadError.value = false

  } catch (err: unknown) {
    console.error('Image upload error:', err)
    const message = err instanceof Error ? err.message : 'Unknown error'
    error.value = `Failed to upload image: ${message}`
  } finally {
    isUploadingImage.value = false
    uploadProgress.value = 0
    imageFile.value = null

    // Reset file input
    const fileInput = document.getElementById('imageUpload')
    if (fileInput) (fileInput as HTMLInputElement).value = ''
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

const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
  isDragging.value = true
}

const handleDragLeave = (event: DragEvent) => {
  event.preventDefault()
  isDragging.value = false
}

const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  isDragging.value = false

  const files = event.dataTransfer?.files
  if (files && files.length > 0) {
    const file = files[0]
    validateAndUploadFile(file)
  }
}

// Fetch user profile data
const fetchUserProfile = async () => {
  // Check if this is an employee
  if (usersStore.isEmployee && usersStore.employee) {
    // For employees, use employee data instead
    userData.value = {
      user_id: usersStore.employee.employee_id || '',
      username: usersStore.employee.username || '',
      email: '', // Employees don't have email
      image_url: '', // Employees don't have profile image
      created_at: usersStore.employee.created_at || '',
      last_login: ''
    }
    editForm.value = {
      username: usersStore.employee.username || '',
      email: ''
    }
    imageLoadError.value = true // Don't show image for employees
    return
  }

  // Regular user profile fetch
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
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Cannot connect to server'
    error.value = `Network error: ${message}`
  } finally {
    isLoading.value = false
  }
}

// Update user profile
const updateProfile = async () => {
  // Employees cannot update profile through this form
  if (usersStore.isEmployee) {
    error.value = 'Employees cannot update profile here'
    return
  }

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
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Cannot connect to server'
    error.value = `Network error: ${message}`
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
  <div class="min-h-screen bg-gradient-to-br from-green-50 via-emerald-50 to-teal-50 py-12 px-4 sm:px-6">
    <div class="container mx-auto max-w-5xl">

      <!-- Profile Header Card -->
      <div class="bg-white rounded-3xl shadow-lg border border-white/60 p-0 mb-8 overflow-hidden relative group">
        <!-- Cover Gradient -->
        <div class="h-32 bg-gradient-to-r from-green-400 to-emerald-600 relative overflow-visible">
          <div class="absolute inset-0 bg-[url('https://www.transparenttextures.com/patterns/cubes.png')] opacity-10">
          </div>

          <!-- Mini 3D Avatar Walking -->
          <div class="absolute right-4 bottom-0 w-28 h-32 z-10 pointer-events-none hidden md:block">
            <MiniAvatar3D :outfitId="usersStore.user?.avatar_config" />
          </div>
        </div>

        <div class="px-8 pb-8 flex flex-col md:flex-row items-end md:items-end gap-6 -mt-12 relative z-10">
          <!-- Avatar Section -->
          <div class="relative">
            <div
              class="w-32 h-32 rounded-full border-4 border-white shadow-md bg-white flex items-center justify-center overflow-hidden relative group/avatar cursor-pointer"
              :class="{ 'ring-4 ring-green-400/50 ring-offset-2': isDragging }" @dragover="handleDragOver"
              @dragleave="handleDragLeave" @drop="handleDrop" @click="triggerFileUpload">

              <img v-if="userData.image_url && userData.image_url.trim() !== '' && !imageLoadError"
                :src="userData.image_url" :alt="userData.username || 'Profile'"
                class="w-full h-full object-cover transition-transform duration-500 group-hover/avatar:scale-110"
                @error="handleImageError" @load="handleImageLoad" />
              <div v-else class="w-full h-full bg-green-50 flex items-center justify-center">
                <IconUserCircle size="80" class="text-green-300" />
              </div>

              <!-- Upload overlay -->
              <div
                class="absolute inset-0 bg-black/40 flex flex-col items-center justify-center opacity-0 group-hover/avatar:opacity-100 transition-all duration-300 backdrop-blur-[2px]">
                <IconEdit size="24" class="text-white mb-1" />
                <span class="text-white text-[10px] font-medium uppercase tracking-wider">Change</span>
              </div>

              <!-- Upload progress -->
              <div v-if="isUploadingImage" class="absolute inset-0 bg-black/60 flex items-center justify-center">
                <div class="radial-progress text-green-400 text-xs font-bold" :style="{ '--value': uploadProgress }"
                  role="progressbar">{{ uploadProgress }}%</div>
              </div>
            </div>

            <!-- Online Indicator (Example) -->
            <div class="absolute bottom-2 right-2 w-5 h-5 bg-green-500 border-4 border-white rounded-full"></div>

            <!-- Hidden Input -->
            <input id="imageUpload" type="file" accept="image/*" @change="handleFileSelect" class="hidden" />
          </div>

          <!-- User Info -->
          <div class="flex-1 pb-4 text-center md:text-left">
            <h1 class="text-3xl font-bold text-gray-900 tracking-tight pt-12">{{ userData.username || 'Loading...' }}
            </h1>
            <div class="flex flex-wrap items-center justify-center md:justify-start gap-3 mt-1 text-gray-500 text-sm">
              <span
                class="flex items-center gap-1 bg-green-50 text-green-700 px-3 py-1 rounded-full font-medium border border-green-100">
                <span class="w-2 h-2 rounded-full bg-green-500"></span>
                {{ usersStore.user?.role?.toUpperCase() || 'USER' }}
              </span>
              <span>{{ userData.email }}</span>
              <span class="text-gray-300 hidden md:inline">•</span>
              <span>ID: {{ userData.user_id?.substring(0, 8) }}...</span>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-3 mb-4 w-full md:w-auto justify-center md:justify-end">
            <button v-if="!isEditing" @click="isEditing = true"
              class="btn btn-primary px-6 rounded-xl shadow-lg shadow-green-200 hover:shadow-green-300 transition-all min-w-[120px]">
              <IconEdit size="18" /> {{ $t('Profile.edit') }}
            </button>
            <div v-else class="flex items-center gap-2">
              <button @click="updateProfile" :disabled="isLoading"
                class="btn btn-success text-white px-6 rounded-xl shadow-lg hover:shadow-green-200 transition-all min-w-[120px]">
                <IconCheck size="18" /> {{ $t('Profile.save') }}
              </button>
              <button @click="cancelEdit"
                class="btn btn-ghost bg-gray-50 text-gray-500 rounded-xl hover:bg-gray-100 hover:text-red-500 w-12 px-0 shadow-sm border border-gray-100">
                <IconX size="20" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Alert Messages -->
      <Transition name="fade">
        <div v-if="error" class="alert alert-error mb-6 shadow-sm rounded-2xl border-none">
          <IconX size="20" /> <span>{{ error }}</span>
        </div>
      </Transition>
      <Transition name="fade">
        <div v-if="success" class="alert alert-success mb-6 shadow-sm rounded-2xl border-none">
          <IconCheck size="20" /> <span>{{ success }}</span>
        </div>
      </Transition>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column: Details -->
        <div class="lg:col-span-2 space-y-6">

          <!-- Basic Info Card -->
          <div class="bg-white rounded-3xl shadow-sm border border-gray-100 p-8">
            <div class="flex items-center gap-3 mb-6 pb-4 border-b border-gray-50">
              <div class="p-2 bg-blue-50 text-blue-600 rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <h2 class="text-xl font-bold text-gray-800">{{ $t('Profile.basicInfo') }}</h2>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Username -->
              <div class="form-control">
                <label class="label text-sm font-semibold text-gray-500 mb-1 pl-1">{{ $t('Profile.username') }}</label>
                <input v-if="isEditing" v-model="editForm.username" type="text"
                  class="input input-bordered rounded-xl bg-gray-50 focus:bg-white transition-all w-full"
                  :placeholder="$t('Profile.usernamePlaceholder')" />
                <div v-else class="p-3 bg-gray-50 rounded-xl text-gray-800 font-medium border border-transparent">{{
                  userData.username || 'N/A' }}</div>
              </div>

              <!-- Email -->
              <div class="form-control">
                <label class="label text-sm font-semibold text-gray-500 mb-1 pl-1">{{ $t('Profile.email') }}</label>
                <input v-if="isEditing" v-model="editForm.email" type="email"
                  class="input input-bordered rounded-xl bg-gray-50 focus:bg-white transition-all w-full"
                  :placeholder="$t('Profile.emailPlaceholder')" />
                <div v-else
                  class="p-3 bg-gray-50 rounded-xl text-gray-800 font-medium border border-transparent flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-400" viewBox="0 0 20 20"
                    fill="currentColor">
                    <path d="M2.003 5.884L10 9.882l7.997-3.998A2 2 0 0016 4H4a2 2 0 00-1.997 1.884z" />
                    <path d="M18 8.118l-8 4-8-4V14a2 2 0 002 2h12a2 2 0 002-2V8.118z" />
                  </svg>
                  {{ userData.email || 'N/A' }}
                </div>
              </div>
            </div>
          </div>

          <!-- Account Info Card -->
          <div class="bg-white rounded-3xl shadow-sm border border-gray-100 p-8">
            <div class="flex items-center gap-3 mb-6 pb-4 border-b border-gray-50">
              <div class="p-2 bg-purple-50 text-purple-600 rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <h2 class="text-xl font-bold text-gray-800">{{ $t('Profile.accountInfo') }}</h2>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="label text-sm font-semibold text-gray-500 mb-1 pl-1">{{ $t('Profile.memberSince')
                }}</label>
                <div class="p-3 bg-white border border-gray-100 rounded-xl text-gray-700 flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  {{ formatDate(userData.created_at) }}
                </div>
              </div>
              <div>
                <label class="label text-sm font-semibold text-gray-500 mb-1 pl-1">{{ $t('Profile.lastLogin') }}</label>
                <div class="p-3 bg-white border border-gray-100 rounded-xl text-gray-700 flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  {{ formatDate(userData.last_login) }}
                </div>
              </div>
            </div>
          </div>

        </div>

        <!-- Right Column: Activities / Status -->
        <div class="space-y-6">
          <!-- Membership Card (Example) -->
          <div
            class="bg-gradient-to-br from-gray-900 to-gray-800 text-white rounded-3xl shadow-xl p-6 relative overflow-hidden">
            <div class="absolute top-0 right-0 p-8 opacity-10">
              <IconUserCircle size="120" />
            </div>
            <div class="relative z-10">
              <h3 class="text-lg font-medium text-gray-300 mb-1">Membership Status</h3>
              <div
                class="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-yellow-200 to-amber-400">
                {{ usersStore.user?.role === 'admin' ? 'Admin' : 'Standard Member' }}
              </div>
              <div class="mt-8 flex items-center gap-2 text-sm text-gray-400">
                <div class="w-full bg-gray-700 rounded-full h-2 flex-1">
                  <div class="bg-green-500 h-2 rounded-full" style="width: 70%"></div>
                </div>
                <span>Good</span>
              </div>
              <p class="text-xs text-gray-500 mt-2">Active participation score</p>
            </div>
          </div>

          <!-- Activities Menu -->
          <div v-if="usersStore.user?.role !== 'employee'"
            class="bg-white rounded-3xl shadow-sm border border-gray-100 p-6">
            <h3 class="text-lg font-bold text-gray-800 mb-4 px-2">Your Activities</h3>
            <div class="space-y-3">
              <router-link to="/analytics"
                class="flex items-center gap-4 p-4 rounded-2xl bg-gray-50 hover:bg-green-50 hover:text-green-700 transition-all group">
                <div
                  class="w-10 h-10 rounded-xl bg-white flex items-center justify-center shadow-sm text-green-500 group-hover:scale-110 transition-transform">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                </div>
                <div>
                  <div class="font-bold">Recycling Analytics</div>
                  <div class="text-xs text-gray-500 group-hover:text-green-600/70">View dashboard</div>
                </div>
                <svg xmlns="http://www.w3.org/2000/svg"
                  class="h-5 w-5 ml-auto text-gray-300 group-hover:text-green-500 group-hover:translate-x-1 transition-all"
                  fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </router-link>

              <!-- Placeholder for History -->
              <div class="flex items-center gap-4 p-4 rounded-2xl bg-gray-50 opacity-60 grayscale cursor-not-allowed">
                <div class="w-10 h-10 rounded-xl bg-white flex items-center justify-center shadow-sm text-gray-400">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div>
                  <div class="font-bold text-gray-600">Purchase History</div>
                  <div class="text-xs text-gray-400">Coming soon</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.radial-progress {
  --size: 3rem;
  --thickness: 3px;
}
</style>
