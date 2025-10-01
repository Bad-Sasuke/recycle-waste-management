<template>
  <div class="min-h-screen bg-gradient-to-br from-green-50 to-cyan-50 py-8 px-4 sm:px-6">
    <!-- Loading State -->
    <div v-if="shopStore.isLoading" class="flex justify-center items-center min-h-screen">
      <div class="text-center">
        <span class="loading loading-spinner loading-lg text-primary"></span>
        <p class="mt-4 text-gray-600">{{ t('Global.loading') }}</p>
      </div>
    </div>

    <!-- No Shop State -->
    <div v-else-if="!shopStore.hasShop || !shopStore.shop" class="flex justify-center items-center min-h-screen">
      <div class="text-center max-w-md">
        <svg xmlns="http://www.w3.org/2000/svg" class="mx-auto h-24 w-24 text-gray-400" fill="none" viewBox="0 0 24 24"
          stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
        </svg>
        <h2 class="mt-4 text-2xl font-bold text-gray-900">{{ t('Shop.manage.noShop') }}</h2>
        <p class="mt-2 text-gray-600">{{ t('Shop.manage.noShopDescription') }}</p>
        <button @click="$router.push('/create-shop')" class="btn btn-primary mt-6">
          {{ t('Shop.create.title') }}
        </button>
      </div>
    </div>

    <!-- Shop Content -->
    <div v-else class="max-w-4xl mx-auto">
      <!-- Shop Header -->
      <div class="bg-white rounded-2xl shadow-xl overflow-hidden mb-6">
        <div class="bg-gradient-to-r from-green-500 to-emerald-600 p-6 text-white">
          <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
            <div class="flex items-center gap-4">
              <div class="bg-white/20 p-3 rounded-xl">
                <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-building-store text-3xl"
                  width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
                  stroke-linecap="round" stroke-linejoin="round">
                  <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                  <path d="M3 21l18 0"></path>
                  <path d="M3 7v1a3 3 0 0 0 6 0v-1m0 1a3 3 0 0 0 6 0v-1m0 1a3 3 0 0 0 6 0v-1"></path>
                  <path d="M5 21v-14a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v14"></path>
                  <path d="M4 11h7m-7 4h7m-7 4h16"></path>
                </svg>
              </div>
              <div>
                <h1 class="text-2xl font-bold">{{ shopStore.shopName || t('Shop.manage.title') }}</h1>
                <p class="opacity-90">{{ shopStore.shopAddress || t('Shop.manage.shopDetails') }}</p>
              </div>
            </div>
            <button @click="editMode = !editMode" class="btn btn-primary">
              <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-edit" width="16" height="16"
                viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
                stroke-linejoin="round">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                <path d="M7 7h-1a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-1"></path>
                <path d="M20.385 6.585a2.1 2.1 0 0 0 -2.97 -2.97l-8.415 8.385v3h3l8.385 -8.415z"></path>
                <path d="M16 5l3 3"></path>
              </svg>
              {{ editMode ? t('Shop.manage.cancelEdit') : t('Shop.manage.editMode') }}
            </button>
          </div>
        </div>

        <!-- Shop Image -->
        <div class="p-6 flex justify-center">
          <img :src="shopStore.shop?.image_url || 'https://placehold.co/400x300?text=No+Image'" alt="Shop Image"
            class="w-full max-w-md h-64 object-cover rounded-xl border border-gray-200" />
        </div>
      </div>

      <!-- Shop Details -->
      <div class="bg-white rounded-2xl shadow-xl overflow-hidden">
        <div class="p-6 md:p-8">
          <div v-if="!editMode">
            <!-- Shop Info Display -->
            <div class="space-y-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">{{ t('Shop.manage.shopName') }}</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.name || t('Global.comingSoon') }}</p>
                </div>

                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">{{ t('Shop.manage.email') }}</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.email || t('Global.comingSoon') }}</p>
                </div>

                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">{{ t('Shop.manage.phone') }}</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.phone || t('Global.comingSoon') }}</p>
                </div>

                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">{{ t('Shop.manage.address') }}</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.address || t('Global.comingSoon') }}</p>
                </div>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">{{ t('Shop.manage.openingTime') }}</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.opening_time || t('Global.comingSoon') }}</p>
                </div>

                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">{{ t('Shop.manage.closingTime') }}</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.closing_time || t('Global.comingSoon') }}</p>
                </div>
              </div>

              <div>
                <h3 class="text-lg font-semibold text-gray-700 mb-2">{{ t('Shop.manage.description') }}</h3>
                <p class="text-gray-900">{{ shopStore.shop?.description || t('Global.comingSoon') }}</p>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">{{ t('Shop.manage.latitude') }}</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.latitude || t('Global.comingSoon') }}</p>
                </div>

                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">{{ t('Shop.manage.longitude') }}</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.longitude || t('Global.comingSoon') }}</p>
                </div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex justify-end mt-8 pt-6 border-t border-gray-200">
              <button @click="confirmDelete" class="btn btn-error text-white flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-trash" width="16"
                  height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
                  stroke-linecap="round" stroke-linejoin="round">
                  <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                  <path d="M4 7l16 0"></path>
                  <path d="M10 11v6"></path>
                  <path d="M14 11v6"></path>
                  <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12"></path>
                  <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3"></path>
                </svg>
                {{ t('Shop.manage.deleteShop') }}
              </button>
            </div>
          </div>

          <!-- Edit Mode -->
          <div v-else>
            <form @submit.prevent="updateShop" class="space-y-6">
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
                  <div v-else-if="shopStore.shop?.image_url" class="mb-4">
                    <img :src="shopStore.shop.image_url" alt="Current Shop Image"
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
                </div>
              </div>

              <!-- Shop Name -->
              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.shopName') }} *</span>
                </label>
                <input v-model="editShopData.name" type="text"
                  class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500" required />
              </div>

              <!-- Description -->
              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.description') }}</span>
                </label>
                <textarea v-model="editShopData.description"
                  class="textarea textarea-bordered w-full max-w-md focus:ring-2 focus:ring-green-500"
                  rows="3"></textarea>
              </div>

              <!-- Address -->
              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.address') }} *</span>
                </label>
                <input v-model="editShopData.address" type="text"
                  class="input input-bordered w-full max-w-md focus:ring-2 focus:ring-green-500" required />
              </div>

              <!-- Contact Information -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="form-control w-full">
                  <label class="label">
                    <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.phone') }}</span>
                  </label>
                  <input v-model="editShopData.phone" type="tel"
                    class="input input-bordered w-full focus:ring-2 focus:ring-green-500" />
                </div>

                <div class="form-control w-full">
                  <label class="label">
                    <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.email') }}</span>
                  </label>
                  <input v-model="editShopData.email" type="email"
                    class="input input-bordered w-full focus:ring-2 focus:ring-green-500" />
                </div>
              </div>

              <!-- Business Hours -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="form-control w-full">
                  <label class="label">
                    <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.openingTime') }}</span>
                  </label>
                  <input v-model="editShopData.opening_time" type="time"
                    class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500" />
                </div>

                <div class="form-control w-full">
                  <label class="label">
                    <span class="label-text font-semibold text-gray-700">{{ t('Shop.create.closingTime') }}</span>
                  </label>
                  <input v-model="editShopData.closing_time" type="time"
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
                  <div id="edit-map" class="w-full h-96 rounded-lg border-2 border-gray-200 z-0"></div>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-2">
                    <div class="form-control">
                      <label class="label">
                        <span class="label-text text-sm">{{ t('Shop.create.latitude') }}</span>
                      </label>
                      <input v-model="editShopData.latitude" type="number" step="any"
                        class="input input-bordered input-sm w-full" readonly />
                    </div>
                    <div class="form-control">
                      <label class="label">
                        <span class="label-text text-sm">{{ t('Shop.create.longitude') }}</span>
                      </label>
                      <input v-model="editShopData.longitude" type="number" step="any"
                        class="input input-bordered input-sm w-full" readonly />
                    </div>
                  </div>
                </div>
              </div>

              <!-- Form Actions -->
              <div class="flex flex-col sm:flex-row gap-4 pt-4">
                <button type="submit" :disabled="isSubmitting" class="btn btn-primary flex-1">
                  <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-check"
                    v-if="!isSubmitting" width="16" height="16" viewBox="0 0 24 24" stroke-width="2"
                    stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                    <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                    <path d="M5 12l5 5l10 -10"></path>
                  </svg>
                  <span class="loading loading-spinner" v-else></span>
                  {{ isSubmitting ? t('Shop.manage.updating') : t('Shop.manage.updateShop') }}
                </button>

                <button type="button" @click="editMode = false" class="btn btn-outline flex-1">
                  <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-x" width="16" height="16"
                    viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                    <path d="M18 6l-12 12"></path>
                    <path d="M6 6l12 12"></path>
                  </svg>
                  {{ t('Shop.manage.cancelEdit') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <input type="checkbox" id="delete-modal" class="modal-toggle" v-model="showDeleteModal" />
    <div class="modal" role="dialog">
      <div class="modal-box">
        <h3 class="font-bold text-lg flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-alert-triangle text-error"
            width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
            stroke-linecap="round" stroke-linejoin="round">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M12 9v2m0 4v.01"></path>
            <path d="M5 19h14a2 2 0 0 0 1.84 -2.75l-7.1 -12.25a2 2 0 0 0 -3.5 0l-7.1 12.25a2 2 0 0 0 1.75 2.75"></path>
          </svg>
          {{ t('Shop.manage.confirmDelete') }}
        </h3>
        <p class="py-4">{{ t('Shop.manage.deleteConfirmation') }}</p>
        <div class="modal-action">
          <button @click="showDeleteModal = false" class="btn btn-outline">{{ t('Shop.manage.cancelDelete') }}</button>
          <button @click="deleteShop" class="btn btn-error text-white">{{ t('Shop.manage.deleteShop') }}</button>
        </div>
      </div>
      <label class="modal-backdrop" @click="showDeleteModal = false">Close</label>
    </div>

    <!-- Success/Error Modal -->
    <dialog id="result-modal" class="modal modal-bottom sm:modal-middle">
      <div class="modal-box">
        <div class="flex flex-col items-center gap-4 py-4">
          <!-- Success Icon -->
          <svg v-if="modalType === 'success'" xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <!-- Error Icon -->
          <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-error" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
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
import { ref, reactive, onMounted, onUnmounted, watch } from 'vue';
import { useShopStore } from '@/stores/shop';
import { useRouter } from 'vue-router';
import type { UpdateShopRequest } from '@/types/shop';
import { useI18n } from 'vue-i18n';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

const { t } = useI18n();
const router = useRouter();
const shopStore = useShopStore();

// States
const editMode = ref(false);
const isSubmitting = ref(false);
const showDeleteModal = ref(false);
const modalType = ref<'success' | 'error'>('success');
const modalMessage = ref('');

// Image handling
const shopImage = ref<File | null>(null);
const previewImage = ref<string | null>(null);

// Map handling
let map: L.Map | null = null;
let marker: L.Marker | null = null;

// Show modal function
const showModal = (type: 'success' | 'error', message: string) => {
  modalType.value = type;
  modalMessage.value = message;
  const modal = document.getElementById('result-modal') as HTMLDialogElement;
  modal?.showModal();
};

// Edit shop data - initialize with shop store data
const editShopData = reactive<UpdateShopRequest>({
  name: shopStore.shop?.name || '',
  description: shopStore.shop?.description || '',
  address: shopStore.shop?.address || '',
  phone: shopStore.shop?.phone || '',
  email: shopStore.shop?.email || '',
  opening_time: shopStore.shop?.opening_time || '',
  closing_time: shopStore.shop?.closing_time || '',
  latitude: shopStore.shop?.latitude,
  longitude: shopStore.shop?.longitude
});

// Watch for changes in shop store to update edit form
watch(() => shopStore.shop, (newShop) => {
  if (newShop) {
    editShopData.name = newShop.name || '';
    editShopData.description = newShop.description || '';
    editShopData.address = newShop.address || '';
    editShopData.phone = newShop.phone || '';
    editShopData.email = newShop.email || '';
    editShopData.opening_time = newShop.opening_time || '';
    editShopData.closing_time = newShop.closing_time || '';
    editShopData.latitude = newShop.latitude;
    editShopData.longitude = newShop.longitude;
  }
}, { immediate: true });

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
  // Get shop location or use default (Thailand - Bangkok)
  const lat = editShopData.latitude || 13.7563;
  const lng = editShopData.longitude || 100.5018;

  // Initialize map
  map = L.map('edit-map').setView([lat, lng], 13);

  // Add OpenStreetMap tile layer
  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  }).addTo(map);

  // Add existing marker if shop has location
  if (editShopData.latitude && editShopData.longitude) {
    marker = L.marker([editShopData.latitude, editShopData.longitude]).addTo(map);
    marker.bindPopup(`<b>${t('Shop.create.selectedLocation')}</b><br>Lat: ${editShopData.latitude.toFixed(6)}<br>Lng: ${editShopData.longitude.toFixed(6)}`).openPopup();
  }

  // Add click event to map
  map.on('click', (e: L.LeafletMouseEvent) => {
    const { lat, lng } = e.latlng;

    // Update shop data
    editShopData.latitude = lat;
    editShopData.longitude = lng;

    // Remove existing marker if any
    if (marker) {
      map?.removeLayer(marker);
    }

    // Add new marker
    marker = L.marker([lat, lng]).addTo(map!);
    marker.bindPopup(`<b>${t('Shop.create.selectedLocation')}</b><br>Lat: ${lat.toFixed(6)}<br>Lng: ${lng.toFixed(6)}`).openPopup();
  });
};

const cleanupMap = () => {
  if (map) {
    map.remove();
    map = null;
    marker = null;
  }
};

const updateShop = async () => {
  if (!shopStore.shop?.shop_id) return;

  isSubmitting.value = true;

  try {
    const formData = new FormData();

    // Add text fields to form data
    Object.entries(editShopData).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        formData.append(key, value.toString());
      }
    });

    // Add image if provided
    if (shopImage.value) {
      formData.append('image', shopImage.value);
    }

    const result = await shopStore.updateShop(shopStore.shop.shop_id, formData);

    if (result.success) {
      editMode.value = false;
      showModal('success', t('Shop.manage.updateSuccess'));
    } else {
      showModal('error', result.message);
    }
  } catch (error) {
    console.error('Error updating shop:', error);
    showModal('error', t('Shop.manage.updateError'));
  } finally {
    isSubmitting.value = false;
  }
};

const confirmDelete = () => {
  showDeleteModal.value = true;
};

const deleteShop = async () => {
  if (!shopStore.shop?.shop_id) return;

  try {
    const result = await shopStore.deleteShop(shopStore.shop.shop_id);

    if (result.success) {
      showDeleteModal.value = false;
      // Redirect to create shop page since user no longer has a shop
      // But check if there's a redirect query parameter to go back to intended destination
      const redirectPath = router.currentRoute.value.query.redirect as string | undefined;
      if (redirectPath) {
        await router.push(redirectPath);
      } else {
        await router.push({ name: 'create-shop' });
      }
    } else {
      showModal('error', result.message);
    }
  } catch (error) {
    console.error('Error deleting shop:', error);
    showModal('error', t('Shop.manage.deleteError'));
  }
};

// Watch edit mode to initialize/cleanup map
watch(editMode, (newValue) => {
  if (newValue) {
    // When entering edit mode, initialize map after a short delay
    setTimeout(() => {
      initMap();
    }, 100);
  } else {
    // When exiting edit mode, cleanup map
    cleanupMap();
  }
});

// Load shop data on component mount
onMounted(async () => {
  console.log('ManageShopView mounted');
  console.log('shopStore.checked:', shopStore.checked);
  console.log('shopStore.hasShop:', shopStore.hasShop);
  console.log('shopStore.shop:', shopStore.shop);

  if (!shopStore.checked) {
    console.log('Checking user shop...');
    await shopStore.checkUserShop();
    console.log('After check - hasShop:', shopStore.hasShop);
    console.log('After check - shop:', shopStore.shop);
  }
});

onUnmounted(() => {
  // Clean up map instance on component unmount
  cleanupMap();
});
</script>

<style scoped>
/* Fix Leaflet marker icons */
:deep(.leaflet-default-icon-path) {
  background-image: url('https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png');
}

/* Ensure map container has proper z-index */
#edit-map {
  position: relative;
}
</style>
