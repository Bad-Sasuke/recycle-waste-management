<template>
  <div class="min-h-screen bg-gradient-to-br from-green-50 to-cyan-50 py-8 px-4 sm:px-6">
    <div class="max-w-4xl mx-auto">
      <div class="bg-white rounded-2xl shadow-xl overflow-hidden">
        <!-- Header -->
        <div class="bg-gradient-to-r from-green-500 to-emerald-600 p-6 text-white">
          <h1 class="text-3xl font-bold flex items-center gap-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-building-store text-3xl"
              width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
              stroke-linecap="round" stroke-linejoin="round">
              <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
              <path d="M3 21l18 0"></path>
              <path d="M3 7v1a3 3 0 0 0 6 0v-1m0 1a3 3 0 0 0 6 0v-1m0 1a3 3 0 0 0 6 0v-1"></path>
              <path d="M5 21v-14a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v14"></path>
              <path d="M4 11h7m-7 4h7m-7 4h16"></path>
            </svg>
            {{ t('Shop.create.title') }}
          </h1>
          <p class="mt-2 opacity-90">{{ t('Shop.create.subtitle') }}</p>
        </div>

        <!-- Form Section -->
        <div class="p-6 md:p-8">
          <form @submit.prevent="submitShop" class="space-y-6">
            <!-- Shop Image Upload -->
            <div class="form-control w-full">
              <label class="label">
                <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.shopImage') }}</span>
              </label>
              <div class="flex flex-col items-center">
                <div v-if="previewImage" class="mb-4">
                  <img :src="previewImage" alt="Preview"
                    class="w-48 h-48 object-cover rounded-xl border-2 border-gray-200" />
                </div>
                <label class="btn btn-outline w-full max-w-xs cursor-pointer">
                  <span class="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-upload" width="16"
                      height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
                      stroke-linecap="round" stroke-linejoin="round">
                      <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                      <path d="M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2 -2v-2"></path>
                      <path d="M7 9l5 -5l5 5"></path>
                      <path d="M12 4v12"></path>
                    </svg>
                    {{ shopImage ? t('Shop.create.changeImage') : t('Shop.create.uploadImage') }}
                  </span>
                  <input type="file" class="hidden" accept="image/*" @change="handleImageChange" />
                </label>
                <p class="text-sm text-gray-500 mt-2">{{ t('Shop.create.uploadImagePlaceholder') }}</p>
              </div>
            </div>

            <!-- Shop Name -->
            <div class="form-control w-full">
              <label class="label">
                <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.shopName') }} *</span>
              </label>
              <input v-model="shopData.name" type="text" :placeholder="t('Shop.create.shopNamePlaceholder')"
                class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500" required />
            </div>

            <!-- Description -->
            <div class="form-control w-full">
              <label class="label">
                <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.description') }}</span>
              </label>
              <textarea v-model="shopData.description" :placeholder="t('Shop.create.descriptionPlaceholder')"
                class="textarea textarea-bordered w-full max-w-md focus:ring-2 focus:ring-green-500"
                rows="3"></textarea>
            </div>

            <!-- Address -->
            <div class="form-control w-full">
              <label class="label">
                <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.address') }} *</span>
              </label>
              <input v-model="shopData.address" type="text" :placeholder="t('Shop.create.addressPlaceholder')"
                class="input input-bordered w-full max-w-md focus:ring-2 focus:ring-green-500" required />
            </div>

            <!-- Contact Information -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.phone') }}</span>
                </label>
                <input v-model="shopData.phone" type="tel" :placeholder="t('Shop.create.phonePlaceholder')"
                  class="input input-bordered w-full focus:ring-2 focus:ring-green-500" />
              </div>

              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.email') }}</span>
                </label>
                <input v-model="shopData.email" type="email" :placeholder="t('Shop.create.emailPlaceholder')"
                  class="input input-bordered w-full focus:ring-2 focus:ring-green-500" />
              </div>
            </div>

            <!-- Business Hours -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.openingTime') }}</span>
                </label>
                <input v-model="shopData.opening_time" type="time"
                  class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500" />
              </div>

              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.closingTime') }}</span>
                </label>
                <input v-model="shopData.closing_time" type="time"
                  class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500" />
              </div>
            </div>

            <!-- Location Map -->
            <div class="form-control w-full">
              <label class="label">
                <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.location') }}</span>
              </label>
              <div class="space-y-2">
                <p class="text-sm text-gray-600">{{ t('Shop.create.locationHelp') }}</p>
                <div id="map" class="w-full h-96 rounded-lg border-2 border-gray-200 z-0"></div>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-2">
                  <div class="form-control">
                    <label class="label">
                      <span class="label-text text-sm">{{ t('Shop.create.latitude') }}</span>
                    </label>
                    <input v-model="shopData.latitude" type="number" step="any"
                      :placeholder="t('Shop.create.latitudePlaceholder')" class="input input-bordered input-sm w-full"
                      readonly />
                  </div>
                  <div class="form-control">
                    <label class="label">
                      <span class="label-text text-sm">{{ t('Shop.create.longitude') }}</span>
                    </label>
                    <input v-model="shopData.longitude" type="number" step="any"
                      :placeholder="t('Shop.create.longitudePlaceholder')" class="input input-bordered input-sm w-full"
                      readonly />
                  </div>
                </div>
              </div>
            </div>

            <!-- Form Actions -->
            <div class="flex flex-col sm:flex-row gap-4 pt-4">
              <button type="submit" :disabled="isSubmitting" class="btn btn-primary flex-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-check" v-if="!isSubmitting"
                  width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
                  stroke-linecap="round" stroke-linejoin="round">
                  <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                  <path d="M5 12l5 5l10 -10"></path>
                </svg>
                <span class="loading loading-spinner" v-else></span>
                {{ isSubmitting ? t('Shop.create.creating') : t('Shop.create.createShop') }}
              </button>

              <button type="button" @click="cancel" class="btn btn-outline flex-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-x" width="16" height="16"
                  viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
                  stroke-linejoin="round">
                  <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                  <path d="M18 6l-12 12"></path>
                  <path d="M6 6l12 12"></path>
                </svg>
                {{ t('Shop.create.cancel') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Success/Error Modal -->
    <dialog id="create-result-modal" class="modal modal-bottom sm:modal-middle">
      <div class="modal-box">
        <div class="flex flex-col items-center gap-4 py-4">
          <!-- Success Icon -->
          <svg v-if="modalType === 'success'" xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-success"
            fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <!-- Error Icon -->
          <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-error" fill="none" viewBox="0 0 24 24"
            stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>

          <h3 class="font-bold text-xl">
            {{ modalType === 'success' ? t('Global.success') : t('Global.error') }}
          </h3>
          <p class="text-center">{{ modalMessage }}</p>
        </div>
        <div class="modal-action">
          <form method="dialog">
            <button class="btn btn-primary">{{ t('Global.close') }}</button>
          </form>
        </div>
      </div>
    </dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { useShopStore } from '@/stores/shop';
import { useUsersStore } from '@/stores/users';
import { useRouter } from 'vue-router';
import type { CreateShopRequest } from '@/types/shop';
import { useI18n } from 'vue-i18n';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

const { t } = useI18n();
const router = useRouter();
const shopStore = useShopStore();
const usersStore = useUsersStore();

// Loading states
const isSubmitting = ref(false);
const modalType = ref<'success' | 'error'>('success');
const modalMessage = ref('');

// Show modal function
const showModal = (type: 'success' | 'error', message: string) => {
  modalType.value = type;
  modalMessage.value = message;
  const modal = document.getElementById('create-result-modal') as HTMLDialogElement;
  modal?.showModal();
};

// Shop data
const shopData = reactive<CreateShopRequest>({
  name: '',
  description: '',
  address: '',
  phone: '',
  email: '',
  opening_time: '',
  closing_time: '',
  latitude: undefined,
  longitude: undefined
});

// Image handling
const shopImage = ref<File | null>(null);
const previewImage = ref<string | null>(null);

// Map handling
let map: L.Map | null = null;
let marker: L.Marker | null = null;

const handleImageChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    const file = target.files[0];
    shopImage.value = file;

    // Create preview
    const reader = new FileReader();
    reader.onload = (e) => {
      previewImage.value = e.target?.result as string;
    };
    reader.readAsDataURL(file);
  }
};

const initMap = () => {
  // Default center (Thailand - Bangkok)
  const defaultLat = 13.7563;
  const defaultLng = 100.5018;

  // Initialize map
  map = L.map('map').setView([defaultLat, defaultLng], 13);

  // Add OpenStreetMap tile layer
  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  }).addTo(map);

  // Add click event to map
  map.on('click', (e: L.LeafletMouseEvent) => {
    const { lat, lng } = e.latlng;

    // Update shop data
    shopData.latitude = lat;
    shopData.longitude = lng;

    // Remove existing marker if any
    if (marker) {
      map?.removeLayer(marker);
    }

    // Add new marker
    marker = L.marker([lat, lng]).addTo(map!);
    marker.bindPopup(`<b>${t('Shop.create.selectedLocation')}</b><br>Lat: ${lat.toFixed(6)}<br>Lng: ${lng.toFixed(6)}`).openPopup();
  });

  // Try to get user's current location
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords;
        map?.setView([latitude, longitude], 13);

        // Set initial marker at user's location
        shopData.latitude = latitude;
        shopData.longitude = longitude;
        marker = L.marker([latitude, longitude]).addTo(map!);
        marker.bindPopup(`<b>${t('Shop.create.yourLocation')}</b>`).openPopup();
      },
      () => {
        // If geolocation fails, use default location
        console.log('Geolocation not available, using default location');
      }
    );
  }
};

const submitShop = async () => {
  isSubmitting.value = true;

  try {
    const formData = new FormData();

    // Add text fields to form data
    Object.entries(shopData).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        formData.append(key, value.toString());
      }
    });

    // Add image if provided
    if (shopImage.value) {
      formData.append('image', shopImage.value);
    }

    const result = await shopStore.createShop(formData);

    if (result.success) {
      // Check if there's a redirect query parameter to go back to intended destination
      const redirectPath = router.currentRoute.value.query.redirect as string | undefined;
      if (redirectPath) {
        await router.push(redirectPath);
      } else {
        // Redirect to shop management page after successful creation
        await router.push({ name: 'manage-shop' });
      }
    } else {
      showModal('error', result.message);
    }
  } catch (error) {
    console.error('Error creating shop:', error);
    showModal('error', t('Shop.create.createError'));
  } finally {
    isSubmitting.value = false;
  }
};

const cancel = () => {
  router.push({ name: 'home' });
};

onMounted(async () => {
  // Guard: Check if user is moderator and doesn't have a shop yet
  const user = usersStore.user;

  // Check if user role is moderator
  if (!user || user.role !== 'moderator') {
    console.warn('Access denied: User is not a moderator');
    await router.push({ name: 'home' });
    return;
  }

  // Check if user already has a shop using shopStore
  await shopStore.checkUserShop();
  if (shopStore.hasShop) {
    console.warn('Access denied: User already has a shop');
    // Redirect to manage shop instead
    await router.push({ name: 'manage-shop' });
    return;
  }

  // Reset any previous shop data
  Object.keys(shopData).forEach(key => {
    (shopData as Record<string, string | number | undefined>)[key] = key === 'latitude' || key === 'longitude' ? undefined : '';
  });

  // Initialize map after component is mounted
  setTimeout(() => {
    initMap();
  }, 100);
});

onUnmounted(() => {
  // Clean up map instance
  if (map) {
    map.remove();
    map = null;
  }
});
</script>

<style scoped>
/* Fix Leaflet marker icons */
:deep(.leaflet-default-icon-path) {
  background-image: url('https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png');
}

/* Ensure map container has proper z-index */
#map {
  position: relative;
}
</style>
