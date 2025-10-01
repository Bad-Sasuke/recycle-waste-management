<template>
  <div class="min-h-screen bg-gradient-to-br from-green-50 to-cyan-50 py-8 px-4 sm:px-6">
    <div class="max-w-4xl mx-auto">
      <div class="bg-white rounded-2xl shadow-xl overflow-hidden">
        <!-- Header -->
        <div class="bg-gradient-to-r from-green-500 to-emerald-600 p-6 text-white">
          <h1 class="text-3xl font-bold flex items-center gap-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-building-store text-3xl" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
               <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
               <path d="M3 21l18 0"></path>
               <path d="M3 7v1a3 3 0 0 0 6 0v-1m0 1a3 3 0 0 0 6 0v-1m0 1a3 3 0 0 0 6 0v-1"></path>
               <path d="M5 21v-14a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v14"></path>
               <path d="M4 11h7m-7 4h7m-7 4h16"></path>
            </svg>
            Create Your Recycling Shop
          </h1>
          <p class="mt-2 opacity-90">Fill in your shop details to start accepting recyclable materials</p>
        </div>

        <!-- Form Section -->
        <div class="p-6 md:p-8">
          <form @submit.prevent="submitShop" class="space-y-6">
            <!-- Shop Image Upload -->
            <div class="form-control w-full">
              <label class="label">
                <span class="label-text font-semibold text-gray-700">Shop Image</span>
              </label>
              <div class="flex flex-col items-center">
                <div v-if="previewImage" class="mb-4">
                  <img :src="previewImage" alt="Preview" class="w-48 h-48 object-cover rounded-xl border-2 border-gray-200" />
                </div>
                <label class="btn btn-outline w-full max-w-xs cursor-pointer">
                  <span class="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-upload" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                      <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                      <path d="M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2 -2v-2"></path>
                      <path d="M7 9l5 -5l5 5"></path>
                      <path d="M12 4v12"></path>
                    </svg>
                    {{ shopImage ? 'Change Image' : 'Upload Shop Image' }}
                  </span>
                  <input 
                    type="file" 
                    class="hidden" 
                    accept="image/*"
                    @change="handleImageChange"
                  />
                </label>
                <p class="text-sm text-gray-500 mt-2">Upload a JPG, PNG, or WebP image (max 5MB)</p>
              </div>
            </div>

            <!-- Shop Name -->
            <div class="form-control w-full">
              <label class="label">
                <span class="label-text font-semibold text-gray-700">Shop Name *</span>
              </label>
              <input
                v-model="shopData.name"
                type="text"
                placeholder="Enter your shop name"
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
                v-model="shopData.description"
                placeholder="Describe your shop and what materials you accept"
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
                v-model="shopData.address"
                type="text"
                placeholder="Enter your shop address"
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
                  v-model="shopData.phone"
                  type="tel"
                  placeholder="Enter your phone number"
                  class="input input-bordered w-full focus:ring-2 focus:ring-green-500"
                />
              </div>

              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">Email</span>
              </label>
                <input
                  v-model="shopData.email"
                  type="email"
                  placeholder="Enter your email address"
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
                  v-model="shopData.opening_time"
                  type="time"
                  class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500"
                />
              </div>

              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">Closing Time</span>
                </label>
                <input
                  v-model="shopData.closing_time"
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
                  v-model="shopData.latitude"
                  type="number"
                  step="any"
                  placeholder="Enter latitude coordinates"
                  class="input input-bordered w-full max-w-xs focus:ring-2 focus:ring-green-500"
                />
              </div>

              <div class="form-control w-full">
                <label class="label">
                  <span class="label-text font-semibold text-gray-700">Longitude</span>
                </label>
                <input
                  v-model="shopData.longitude"
                  type="number"
                  step="any"
                  placeholder="Enter longitude coordinates"
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
                {{ isSubmitting ? 'Creating Shop...' : 'Create Shop' }}
              </button>

              <button
                type="button"
                @click="cancel"
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
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useShopStore } from '@/stores/shop';
import { useRouter } from 'vue-router';
import type { CreateShopRequest } from '@/types/shop';

const router = useRouter();
const shopStore = useShopStore();

// Loading states
const isSubmitting = ref(false);

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
      // Redirect to shop management page after successful creation
      await router.push({ name: 'manage-shop' });
    } else {
      alert(`Error: ${result.message}`);
    }
  } catch (error) {
    console.error('Error creating shop:', error);
    alert('An error occurred while creating your shop. Please try again.');
  } finally {
    isSubmitting.value = false;
  }
};

const cancel = () => {
  router.push({ name: 'home' });
};

onMounted(() => {
  // Reset any previous shop data
  Object.keys(shopData).forEach(key => {
    (shopData as any)[key] = key === 'latitude' || key === 'longitude' ? undefined : '';
  });
});
</script>