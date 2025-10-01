<template>
  <div class="min-h-screen bg-gradient-to-br from-green-50 to-cyan-50 py-8 px-4 sm:px-6">
    <div class="max-w-4xl mx-auto">
      <!-- Shop Header -->
      <div class="bg-white rounded-2xl shadow-xl overflow-hidden mb-6">
        <div class="bg-gradient-to-r from-green-500 to-emerald-600 p-6 text-white">
          <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
            <div class="flex items-center gap-4">
              <div class="bg-white/20 p-3 rounded-xl">
                <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-building-store text-3xl" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                  <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                  <path d="M3 21l18 0"></path>
                  <path d="M3 7v1a3 3 0 0 0 6 0v-1m0 1a3 3 0 0 0 6 0v-1m0 1a3 3 0 0 0 6 0v-1"></path>
                  <path d="M5 21v-14a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v14"></path>
                  <path d="M4 11h7m-7 4h7m-7 4h16"></path>
                </svg>
              </div>
              <div>
                <h1 class="text-2xl font-bold">{{ shopStore.shopName || 'Your Shop' }}</h1>
                <p class="opacity-90">{{ shopStore.shopAddress || 'Shop address' }}</p>
              </div>
            </div>
            <button @click="editMode = !editMode" class="btn btn-outline btn-primary">
              <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-edit" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                <path d="M7 7h-1a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-1"></path>
                <path d="M20.385 6.585a2.1 2.1 0 0 0 -2.97 -2.97l-8.415 8.385v3h3l8.385 -8.415z"></path>
                <path d="M16 5l3 3"></path>
              </svg>
              {{ editMode ? 'Cancel' : 'Edit Shop' }}
            </button>
          </div>
        </div>

        <!-- Shop Image -->
        <div class="p-6 flex justify-center">
          <img 
            :src="shopStore.shop?.image_url || 'https://placehold.co/400x300?text=No+Image'" 
            alt="Shop Image"
            class="w-full max-w-md h-64 object-cover rounded-xl border border-gray-200"
          />
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
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">Shop Name</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.name || 'Not set' }}</p>
                </div>
                
                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">Email</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.email || 'Not provided' }}</p>
                </div>
                
                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">Phone</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.phone || 'Not provided' }}</p>
                </div>
                
                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">Address</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.address || 'Not provided' }}</p>
                </div>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">Opening Time</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.opening_time || 'Not set' }}</p>
                </div>
                
                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">Closing Time</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.closing_time || 'Not set' }}</p>
                </div>
              </div>

              <div>
                <h3 class="text-lg font-semibold text-gray-700 mb-2">Description</h3>
                <p class="text-gray-900">{{ shopStore.shop?.description || 'No description provided' }}</p>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">Latitude</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.latitude || 'Not set' }}</p>
                </div>
                
                <div>
                  <h3 class="text-lg font-semibold text-gray-700 mb-2">Longitude</h3>
                  <p class="text-gray-900">{{ shopStore.shop?.longitude || 'Not set' }}</p>
                </div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex justify-end mt-8 pt-6 border-t border-gray-200">
              <button 
                @click="confirmDelete"
                class="btn btn-error flex items-center gap-2"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-trash" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                  <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                  <path d="M4 7l16 0"></path>
                  <path d="M10 11v6"></path>
                  <path d="M14 11v6"></path>
                  <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12"></path>
                  <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3"></path>
                </svg>
                Delete Shop
              </button>
            </div>
          </div>

          <!-- Edit Mode -->
          <div v-else>
            <form @submit.prevent="updateShop" class="space-y-6">
              <!-- Shop Image Upload -->
              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">Shop Image</span>
                </label>
                <div class="flex flex-col items-center">
                  <div v-if="previewImage" class="mb-4">
                    <img :src="previewImage" alt="Preview" class="w-48 h-48 object-cover rounded-xl border-2 border-gray-200" />
                  </div>
                  <div v-else-if="shopStore.shop?.image_url" class="mb-4">
                    <img :src="shopStore.shop.image_url" alt="Current Shop Image" class="w-48 h-48 object-cover rounded-xl border-2 border-gray-200" />
                  </div>
                  <label class="btn btn-outline w-full max-w-xs cursor-pointer">
                    <span class="flex items-center gap-2">
                      <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-upload" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                        <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                        <path d="M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2 -2v-2"></path>
                        <path d="M7 9l5 -5l5 5"></path>
                        <path d="M12 4v12"></path>
                      </svg>
                      {{ shopImage ? 'Change Image' : 'Upload New Image' }}
                    </span>
                    <input 
                      type="file" 
                      class="hidden" 
                      accept="image/*"
                      @change="handleImageChange"
                    />
                  </label>
                </div>
              </div>

              <!-- Shop Name -->
              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">Shop Name *</span>
                </label>
                <input
                  v-model="editShopData.name"
                  type="text"
                  class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500"
                  required
                />
              </div>

              <!-- Description -->
              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">Description</span>
                </label>
                <textarea
                  v-model="editShopData.description"
                  class="textarea textarea-bordered w-full max-w-md focus:ring-2 focus:ring-green-500"
                  rows="3"
                ></textarea>
              </div>

              <!-- Address -->
              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">Address *</span>
                </label>
                <input
                  v-model="editShopData.address"
                  type="text"
                  class="input input-bordered w-full max-w-md focus:ring-2 focus:ring-green-500"
                  required
                />
              </div>

              <!-- Contact Information -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="form-control w-full">
                  <label class="label">
                    <span class="label-text font-semibold text-gray-700">Phone Number</span>
                  </label>
                  <input
                    v-model="editShopData.phone"
                    type="tel"
                    class="input input-bordered w-full focus:ring-2 focus:ring-green-500"
                  />
                </div>

                <div class="form-control w-full">
                  <label class="label">
                    <span class="label-text font-semibold text-gray-700">Email</span>
                  </label>
                  <input
                    v-model="editShopData.email"
                    type="email"
                    class="input input-bordered w-full focus:ring-2 focus:ring-green-500"
                  />
                </div>
              </div>

              <!-- Business Hours -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="form-control w-full">
                  <label class="label">
                    <span class="label-text font-semibold text-gray-700">Opening Time</span>
                  </label>
                  <input
                    v-model="editShopData.opening_time"
                    type="time"
                    class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500"
                  />
                </div>

                <div class="form-control w-full">
                  <label class="label">
                    <span class="label-text font-semibold text-gray-700">Closing Time</span>
                  </label>
                  <input
                    v-model="editShopData.closing_time"
                    type="time"
                    class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500"
                  />
                </div>
              </div>

              <!-- Location Coordinates -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="form-control w-full">
                  <label class="label">
                    <span class="label-text font-semibold text-gray-700">Latitude</span>
                  </label>
                  <input
                    v-model="editShopData.latitude"
                    type="number"
                    step="any"
                    class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500"
                  />
                </div>

                <div class="form-control w-full">
                  <label class="label">
                    <span class="label-text font-semibold text-gray-700">Longitude</span>
                  </label>
                  <input
                    v-model="editShopData.longitude"
                    type="number"
                    step="any"
                    class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500"
                  />
                </div>
              </div>

              <!-- Form Actions -->
              <div class="flex flex-col sm:flex-row gap-4 pt-4">
                <button
                  type="submit"
                  :disabled="isSubmitting"
                  class="btn btn-primary flex-1"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-check" v-if="!isSubmitting" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                    <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                    <path d="M5 12l5 5l10 -10"></path>
                  </svg>
                  <span class="loading loading-spinner" v-else></span>
                  {{ isSubmitting ? 'Updating Shop...' : 'Update Shop' }}
                </button>

                <button
                  type="button"
                  @click="editMode = false"
                  class="btn btn-outline flex-1"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-x" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                    <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                    <path d="M18 6l-12 12"></path>
                    <path d="M6 6l12 12"></path>
                  </svg>
                  Cancel
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
          <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-alert-triangle text-error" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M12 9v2m0 4v.01"></path>
            <path d="M5 19h14a2 2 0 0 0 1.84 -2.75l-7.1 -12.25a2 2 0 0 0 -3.5 0l-7.1 12.25a2 2 0 0 0 1.75 2.75"></path>
          </svg>
          Confirm Deletion
        </h3>
        <p class="py-4">Are you sure you want to delete your shop? This action cannot be undone.</p>
        <div class="modal-action">
          <button @click="showDeleteModal = false" class="btn btn-outline">Cancel</button>
          <button @click="deleteShop" class="btn btn-error">Delete Shop</button>
        </div>
      </div>
      <label class="modal-backdrop" @click="showDeleteModal = false">Close</label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue';
import { useShopStore } from '@/stores/shop';
import { useRouter } from 'vue-router';
import type { UpdateShopRequest } from '@/types/shop';

const router = useRouter();
const shopStore = useShopStore();

// States
const editMode = ref(false);
const isSubmitting = ref(false);
const showDeleteModal = ref(false);

// Image handling
const shopImage = ref<File | null>(null);
const previewImage = ref<string | null>(null);

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
      // Show success message or notification
      alert('Shop updated successfully!');
    } else {
      alert(`Error: ${result.message}`);
    }
  } catch (error) {
    console.error('Error updating shop:', error);
    alert('An error occurred while updating your shop. Please try again.');
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
      await router.push({ name: 'create-shop' });
    } else {
      alert(`Error: ${result.message}`);
    }
  } catch (error) {
    console.error('Error deleting shop:', error);
    alert('An error occurred while deleting your shop. Please try again.');
  }
};

// Load shop data on component mount
onMounted(async () => {
  if (!shopStore.checked) {
    await shopStore.checkUserShop();
  }
});
</script>