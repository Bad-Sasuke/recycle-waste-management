<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { IconSettings, IconBell, IconLock, IconPalette, IconCheck, IconX } from '@tabler/icons-vue'
import { useUsersStore } from '@/stores/users'
import { useI18nStore } from '@/stores/i18n'

const usersStore = useUsersStore()
const i18nStore = useI18nStore()

// Settings state
const settings = ref({
  // Account settings
  emailNotifications: true,
  pushNotifications: false,
  marketplaceUpdates: true,
  shopUpdates: true,

  // Privacy settings
  profileVisibility: 'public',
  showEmail: false,
  showPhone: false,

  // Appearance settings
  theme: 'light',
  language: 'th'
})

const isLoading = ref(false)
const success = ref('')
const error = ref('')

// Theme options
const themeOptions = [
  { value: 'light', label: 'Light' },
  { value: 'dark', label: 'Dark' },
  { value: 'auto', label: 'Auto' }
]

// Language options
const languageOptions = [
  { value: 'th', label: 'ภาษาไทย' },
  { value: 'en', label: 'English' }
]

// Profile visibility options
const visibilityOptions = [
  { value: 'public', labelKey: 'Settings.privacy.public' },
  { value: 'private', labelKey: 'Settings.privacy.private' },
  { value: 'friends', labelKey: 'Settings.privacy.friends' }
]

// Save settings
const saveSettings = async () => {
  isLoading.value = true
  success.value = ''
  error.value = ''

  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))

    // Save to localStorage
    localStorage.setItem('userSettings', JSON.stringify(settings.value))

    // Update language if changed
    if (settings.value.language !== i18nStore.locale) {
      i18nStore.setLocale(settings.value.language)
    }

    success.value = 'Settings saved successfully!'

    // Clear success message after 3 seconds
    setTimeout(() => {
      success.value = ''
    }, 3000)
  } catch (err: unknown) {
    error.value = 'Failed to save settings. Please try again : ' + err
  } finally {
    isLoading.value = false
  }
}

// Load settings from localStorage
const loadSettings = () => {
  const saved = localStorage.getItem('userSettings')
  if (saved) {
    settings.value = { ...settings.value, ...JSON.parse(saved) }
  }
  // Set current language
  settings.value.language = i18nStore.locale
}

onMounted(() => {
  if (!usersStore.isLogin) {
    error.value = 'Please login to access settings'
    return
  }
  loadSettings()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="container mx-auto px-4 max-w-4xl">
      <!-- Header -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-6">
        <div class="flex items-center gap-4">
          <div class="avatar placeholder">
            <div class="w-16 h-16 rounded-full bg-green-100">
              <IconSettings size="32" class="text-green-600" stroke-width="1.5" />
            </div>
          </div>
          <div>
            <h1 class="text-3xl font-bold text-gray-800">{{ $t('Settings.title') }}</h1>
            <p class="text-gray-600">{{ $t('Settings.subtitle') }}</p>
          </div>
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

      <!-- Settings Content -->
      <div class="space-y-6">
        <!-- Notifications Settings -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center gap-3 mb-6">
            <IconBell size="24" class="text-gray-700" />
            <h2 class="text-xl font-semibold text-gray-800">{{ $t('Settings.notifications.title') }}</h2>
          </div>

          <div class="space-y-4">
            <!-- Email Notifications -->
            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-4">
                <input type="checkbox" class="toggle toggle-success" v-model="settings.emailNotifications" />
                <div>
                  <span class="label-text font-medium">{{ $t('Settings.notifications.emailNotifications') }}</span>
                  <p class="text-sm text-gray-500">{{ $t('Settings.notifications.emailNotificationsDesc') }}</p>
                </div>
              </label>
            </div>

            <!-- Push Notifications -->
            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-4">
                <input type="checkbox" class="toggle toggle-success" v-model="settings.pushNotifications" />
                <div>
                  <span class="label-text font-medium">{{ $t('Settings.notifications.pushNotifications') }}</span>
                  <p class="text-sm text-gray-500">{{ $t('Settings.notifications.pushNotificationsDesc') }}</p>
                </div>
              </label>
            </div>

            <!-- Marketplace Updates -->
            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-4">
                <input type="checkbox" class="toggle toggle-success" v-model="settings.marketplaceUpdates" />
                <div>
                  <span class="label-text font-medium">{{ $t('Settings.notifications.marketplaceUpdates') }}</span>
                  <p class="text-sm text-gray-500">{{ $t('Settings.notifications.marketplaceUpdatesDesc') }}</p>
                </div>
              </label>
            </div>

            <!-- Shop Updates -->
            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-4">
                <input type="checkbox" class="toggle toggle-success" v-model="settings.shopUpdates" />
                <div>
                  <span class="label-text font-medium">{{ $t('Settings.notifications.shopUpdates') }}</span>
                  <p class="text-sm text-gray-500">{{ $t('Settings.notifications.shopUpdatesDesc') }}</p>
                </div>
              </label>
            </div>
          </div>
        </div>

        <!-- Privacy Settings -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center gap-3 mb-6">
            <IconLock size="24" class="text-gray-700" />
            <h2 class="text-xl font-semibold text-gray-800">{{ $t('Settings.privacy.title') }}</h2>
          </div>

          <div class="space-y-4">
            <!-- Profile Visibility -->
            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">{{ $t('Settings.privacy.profileVisibility') }}</span>
              </label>
              <select class="select select-bordered w-full" v-model="settings.profileVisibility">
                <option v-for="option in visibilityOptions" :key="option.value" :value="option.value">
                  {{ $t(option.labelKey) }}
                </option>
              </select>
              <label class="label">
                <span class="label-text-alt text-gray-500">{{ $t('Settings.privacy.profileVisibilityDesc') }}</span>
              </label>
            </div>

            <!-- Show Email -->
            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-4">
                <input type="checkbox" class="toggle toggle-success" v-model="settings.showEmail" />
                <div>
                  <span class="label-text font-medium">{{ $t('Settings.privacy.showEmail') }}</span>
                  <p class="text-sm text-gray-500">{{ $t('Settings.privacy.showEmailDesc') }}</p>
                </div>
              </label>
            </div>

            <!-- Show Phone -->
            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-4">
                <input type="checkbox" class="toggle toggle-success" v-model="settings.showPhone" />
                <div>
                  <span class="label-text font-medium">{{ $t('Settings.privacy.showPhone') }}</span>
                  <p class="text-sm text-gray-500">{{ $t('Settings.privacy.showPhoneDesc') }}</p>
                </div>
              </label>
            </div>
          </div>
        </div>

        <!-- Appearance Settings -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center gap-3 mb-6">
            <IconPalette size="24" class="text-gray-700" />
            <h2 class="text-xl font-semibold text-gray-800">{{ $t('Settings.appearance.title') }}</h2>
          </div>

          <div class="space-y-4">
            <!-- Theme -->
            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">{{ $t('Settings.appearance.theme') }}</span>
              </label>
              <select class="select select-bordered w-full" v-model="settings.theme">
                <option v-for="option in themeOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>

            <!-- Language -->
            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">{{ $t('Settings.appearance.language') }}</span>
              </label>
              <select class="select select-bordered w-full" v-model="settings.language">
                <option v-for="option in languageOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>
          </div>
        </div>

        <!-- Save Button -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <button @click="saveSettings" :disabled="isLoading" class="btn btn-success w-full">
            <IconCheck v-if="!isLoading" size="20" />
            <span v-if="isLoading" class="loading loading-spinner"></span>
            {{ isLoading ? $t('Settings.saving') : $t('Settings.saveSettings') }}
          </button>
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
