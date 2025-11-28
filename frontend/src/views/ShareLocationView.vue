<script setup lang="ts">
import { ref, shallowRef, computed, onMounted, onUnmounted } from 'vue';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

// Fix Leaflet marker icons
// @ts-ignore
import iconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
// @ts-ignore
import iconUrl from 'leaflet/dist/images/marker-icon.png';
// @ts-ignore
import shadowUrl from 'leaflet/dist/images/marker-shadow.png';

delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl,
  iconUrl,
  shadowUrl,
});

// Form Data
const formData = ref({
  latitude: 0,
  longitude: 0,
  wasteAmount: '',
  wasteTypes: [] as string[],
  notes: '',
  phone: '',
});

// Map
const map = shallowRef<L.Map | null>(null);
const marker = shallowRef<L.Marker | null>(null);

// Request States
const currentRequest = ref<any>(null);
const isSubmitting = ref(false);
const isCanceling = ref(false);

// Chat
const showChat = ref(false);
const chatMessages = ref<any[]>([]);
const newMessage = ref('');

// Mock current request (remove when integrating with backend)
// currentRequest.value = {
//   id: '1',
//   latitude: 13.7563,
//   longitude: 100.5018,
//   wasteAmount: 'medium',
//   wasteTypes: ['plastic', 'paper'],
//   notes: 'มีขยะพลาสติกและกระดาษจำนวนมาก',
//   phone: '081-234-5678',
//   status: 'accepted', // pending, accepted, rejected
//   shopResponse: {
//     shopName: 'ร้านรับซื้อของเก่า กรีนไลฟ์',
//     message: 'เราจะไปรับภายใน 1 ชั่วโมง ขอบคุณครับ'
//   }
// };

// Computed
const isFormValid = computed(() => {
  return formData.value.latitude !== 0
    && formData.value.longitude !== 0
    && formData.value.wasteAmount
    && formData.value.wasteTypes.length > 0
    && formData.value.phone;
});

const statusBadgeClass = computed(() => {
  switch (currentRequest.value?.status) {
    case 'pending':
      return 'badge-warning';
    case 'accepted':
      return 'badge-success';
    case 'rejected':
      return 'badge-error';
    default:
      return 'badge-ghost';
  }
});

const statusText = computed(() => {
  switch (currentRequest.value?.status) {
    case 'pending':
      return 'รอร้านคัดรับ';
    case 'accepted':
      return 'ร้านยอมรับ';
    case 'rejected':
      return 'ร้านปฏิเสธ';
    default:
      return '-';
  }
});

// Methods
const initMap = () => {
  const defaultLat = 13.7563;
  const defaultLng = 100.5018;

  map.value = L.map('location-map').setView([defaultLat, defaultLng], 13);

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  }).addTo(map.value);

  // Click to set location
  map.value.on('click', (e: any) => {
    formData.value.latitude = e.latlng.lat;
    formData.value.longitude = e.latlng.lng;

    // Remove old marker
    if (marker.value) {
      marker.value.remove();
    }

    // Add new marker
    marker.value = L.marker([e.latlng.lat, e.latlng.lng])
      .addTo(map.value!)
      .bindPopup('ตำแหน่งที่คุณเลือก')
      .openPopup();
  });

  // Try to get user's current location
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition((position) => {
      const lat = position.coords.latitude;
      const lng = position.coords.longitude;
      map.value?.setView([lat, lng], 15);

      formData.value.latitude = lat;
      formData.value.longitude = lng;

      marker.value = L.marker([lat, lng])
        .addTo(map.value!)
        .bindPopup('ตำแหน่งปัจจุบันของคุณ')
        .openPopup();
    });
  }
};

const submitLocationRequest = async () => {
  isSubmitting.value = true;

  try {
    // TODO: API call to submit location request
    await new Promise(resolve => setTimeout(resolve, 1500));

    // Mock response
    currentRequest.value = {
      id: Date.now().toString(),
      ...formData.value,
      status: 'pending',
      createdAt: new Date(),
    };

    // Reset form
    formData.value = {
      latitude: 0,
      longitude: 0,
      wasteAmount: '',
      wasteTypes: [],
      notes: '',
      phone: '',
    };

    if (marker.value) {
      marker.value.remove();
      marker.value = null;
    }
  } catch (error) {
    console.error('Error submitting request:', error);
  } finally {
    isSubmitting.value = false;
  }
};

const cancelRequest = async () => {
  if (!confirm('คุณแน่ใจหรือไม่ว่าต้องการยกเลิกคำขอนี้?')) {
    return;
  }

  isCanceling.value = true;

  try {
    // TODO: API call to cancel request
    await new Promise(resolve => setTimeout(resolve, 1000));
    currentRequest.value = null;
  } catch (error) {
    console.error('Error canceling request:', error);
  } finally {
    isCanceling.value = false;
  }
};

const openChat = () => {
  showChat.value = true;
};

const sendMessage = () => {
  if (!newMessage.value.trim()) return;

  chatMessages.value.push({
    sender: 'me',
    text: newMessage.value,
    timestamp: new Date(),
  });

  newMessage.value = '';

  // Mock shop response
  setTimeout(() => {
    chatMessages.value.push({
      sender: 'shop',
      text: 'ขอบคุณครับ เรากำลังเดินทางไปหาคุณ',
      timestamp: new Date(),
    });
  }, 1000);
};

const formatTime = (date: Date) => {
  return new Date(date).toLocaleTimeString('th-TH', {
    hour: '2-digit',
    minute: '2-digit',
  });
};

const getWasteAmountText = (amount: string) => {
  const amounts: Record<string, string> = {
    small: 'น้อย (ถุงเล็ก 1-2 ถุง)',
    medium: 'ปานกลาง (ถุงใหญ่ 1-2 ถุง)',
    large: 'มาก (ถุงใหญ่ 3-5 ถุง)',
    very_large: 'มากที่สุด (มากกว่า 5 ถุง)',
  };
  return amounts[amount] || amount;
};

onMounted(() => {
  initMap();
});

onUnmounted(() => {
  if (map.value) {
    map.value.remove();
  }
});
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-green-50 via-white to-emerald-50">
    <div class="max-w-7xl mx-auto px-4 py-8">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-green-800 mb-2">
          {{ $t('ShareLocation.title') }}
        </h1>
        <p class="text-gray-600">
          {{ $t('ShareLocation.subtitle') }}
        </p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Left: Form & Current Status -->
        <div class="space-y-6">
          <!-- Share Location Form -->
          <div v-if="!currentRequest" class="card bg-base-100 shadow-lg border border-green-100">
            <div class="card-body">
              <h2 class="card-title text-green-700 mb-4">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                {{ $t('ShareLocation.shareYourLocation') }}
              </h2>

              <!-- Map for location selection -->
              <div class="mb-4">
                <label class="label">
                  <span class="label-text font-semibold">{{ $t('ShareLocation.selectLocation')
                  }}</span>
                </label>
                <div id="location-map" class="w-full h-[300px] rounded-lg border-2 border-green-200">
                </div>
                <p class="text-sm text-gray-500 mt-2">
                  {{ $t('ShareLocation.clickMapHelp') }}
                </p>
              </div>

              <!-- Waste Amount -->
              <div class="form-control">
                <label class="label">
                  <span class="label-text font-semibold">{{ $t('ShareLocation.wasteAmount') }}</span>
                </label>
                <select v-model="formData.wasteAmount" class="select select-bordered w-full">
                  <option value="">{{ $t('ShareLocation.selectWasteAmount') }}</option>
                  <option value="small">{{ $t('ShareLocation.wasteSmall') }}</option>
                  <option value="medium">{{ $t('ShareLocation.wasteMedium') }}</option>
                  <option value="large">{{ $t('ShareLocation.wasteLarge') }}</option>
                  <option value="very_large">{{ $t('ShareLocation.wasteVeryLarge') }}</option>
                </select>
              </div>

              <!-- Waste Types -->
              <div class="form-control">
                <label class="label">
                  <span class="label-text font-semibold">{{ $t('ShareLocation.wasteTypes') }}</span>
                </label>
                <div class="grid grid-cols-2 gap-2">
                  <label class="cursor-pointer label justify-start gap-2">
                    <input type="checkbox" v-model="formData.wasteTypes" value="plastic"
                      class="checkbox checkbox-success checkbox-sm" />
                    <span class="label-text">{{ $t('ShareLocation.plastic') }}</span>
                  </label>
                  <label class="cursor-pointer label justify-start gap-2">
                    <input type="checkbox" v-model="formData.wasteTypes" value="paper"
                      class="checkbox checkbox-success checkbox-sm" />
                    <span class="label-text">{{ $t('ShareLocation.paper') }}</span>
                  </label>
                  <label class="cursor-pointer label justify-start gap-2">
                    <input type="checkbox" v-model="formData.wasteTypes" value="metal"
                      class="checkbox checkbox-success checkbox-sm" />
                    <span class="label-text">{{ $t('ShareLocation.metal') }}</span>
                  </label>
                  <label class="cursor-pointer label justify-start gap-2">
                    <input type="checkbox" v-model="formData.wasteTypes" value="glass"
                      class="checkbox checkbox-success checkbox-sm" />
                    <span class="label-text">{{ $t('ShareLocation.glass') }}</span>
                  </label>
                </div>
              </div>

              <!-- Additional Notes -->
              <div class="form-control">
                <label class="label">
                  <span class="label-text font-semibold">{{ $t('ShareLocation.additionalNotes')
                  }}</span>
                </label>
                <textarea v-model="formData.notes" class="textarea textarea-bordered h-24"
                  :placeholder="$t('ShareLocation.notesPlaceholder')"></textarea>
              </div>

              <!-- Contact Info -->
              <div class="form-control">
                <label class="label">
                  <span class="label-text font-semibold">{{ $t('ShareLocation.contactPhone') }}</span>
                </label>
                <input v-model="formData.phone" type="tel" class="input input-bordered"
                  :placeholder="$t('ShareLocation.phonePlaceholder')" />
              </div>

              <!-- Submit Button -->
              <div class="card-actions justify-end mt-4">
                <button @click="submitLocationRequest" class="btn btn-primary text-white"
                  :disabled="!isFormValid || isSubmitting">
                  <svg v-if="!isSubmitting" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                  </svg>
                  <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
                  {{ isSubmitting ? $t('ShareLocation.submitting') : $t('ShareLocation.shareLocation')
                  }}
                </button>
              </div>
            </div>
          </div>

          <!-- Current Request Status -->
          <div v-else class="card bg-base-100 shadow-lg border-2 border-green-200">
            <div class="card-body">
              <div class="flex items-center justify-between mb-4">
                <h2 class="card-title text-green-700">
                  {{ $t('ShareLocation.currentRequest') }}
                </h2>
                <div class="badge badge-lg" :class="statusBadgeClass">
                  {{ statusText }}
                </div>
              </div>

              <!-- Request Details -->
              <div class="space-y-3">
                <div class="flex items-start gap-3">
                  <svg class="w-5 h-5 text-gray-500 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  </svg>
                  <div>
                    <p class="font-semibold">{{ $t('ShareLocation.location') }}</p>
                    <p class="text-sm text-gray-600">{{ currentRequest.latitude }}, {{
                      currentRequest.longitude }}</p>
                  </div>
                </div>

                <div class="flex items-start gap-3">
                  <svg class="w-5 h-5 text-gray-500 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                  </svg>
                  <div>
                    <p class="font-semibold">{{ $t('ShareLocation.wasteAmount') }}</p>
                    <p class="text-sm">{{ getWasteAmountText(currentRequest.wasteAmount) }}</p>
                  </div>
                </div>

                <div class="flex items-start gap-3">
                  <svg class="w-5 h-5 text-gray-500 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                  </svg>
                  <div>
                    <p class="font-semibold">{{ $t('ShareLocation.wasteTypes') }}</p>
                    <p class="text-sm">{{ currentRequest.wasteTypes.join(', ') }}</p>
                  </div>
                </div>

                <div v-if="currentRequest.notes" class="flex items-start gap-3">
                  <svg class="w-5 h-5 text-gray-500 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
                  </svg>
                  <div>
                    <p class="font-semibold">{{ $t('ShareLocation.notes') }}</p>
                    <p class="text-sm text-gray-600">{{ currentRequest.notes }}</p>
                  </div>
                </div>

                <!-- Shop Response -->
                <div v-if="currentRequest.shopResponse" class="mt-4 p-4 rounded-lg"
                  :class="currentRequest.status === 'accepted' ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'">
                  <div class="flex items-start gap-3">
                    <svg v-if="currentRequest.status === 'accepted'" class="w-6 h-6 text-green-600" fill="none"
                      stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <svg v-else class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <div class="flex-1">
                      <p class="font-semibold"
                        :class="currentRequest.status === 'accepted' ? 'text-green-800' : 'text-red-800'">
                        {{ currentRequest.shopResponse.shopName }}
                      </p>
                      <p v-if="currentRequest.shopResponse.message" class="text-sm mt-1">
                        {{ currentRequest.shopResponse.message }}
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Actions -->
              <div class="card-actions justify-end mt-4">
                <button v-if="currentRequest.status === 'pending'" @click="cancelRequest"
                  class="btn btn-outline btn-error" :disabled="isCanceling">
                  <span v-if="isCanceling" class="loading loading-spinner loading-sm"></span>
                  {{ isCanceling ? $t('ShareLocation.canceling') : $t('ShareLocation.cancelRequest')
                  }}
                </button>
                <button v-if="currentRequest.status === 'accepted'" @click="openChat"
                  class="btn btn-success text-white">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                  {{ $t('ShareLocation.chat') }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Live Chat -->
        <div v-if="currentRequest && currentRequest.status === 'accepted' && showChat"
          class="card bg-base-100 shadow-lg border border-green-100">
          <div class="card-body p-0 flex flex-col h-[600px]">
            <!-- Chat Header -->
            <div class="p-4 border-b border-gray-200 bg-gradient-to-r from-green-500 to-emerald-600">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="avatar placeholder">
                    <div class="bg-white text-green-600 rounded-full w-10">
                      <span class="text-lg">🏪</span>
                    </div>
                  </div>
                  <div class="text-white">
                    <h3 class="font-bold">{{ currentRequest.shopResponse?.shopName }}</h3>
                    <p class="text-xs opacity-90">{{ $t('ShareLocation.online') }}</p>
                  </div>
                </div>
                <button @click="showChat = false" class="btn btn-ghost btn-sm btn-circle text-white">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>

            <!-- Chat Messages -->
            <div class="flex-1 overflow-y-auto p-4 space-y-3 bg-gray-50">
              <div v-for="(message, index) in chatMessages" :key="index" class="flex"
                :class="message.sender === 'me' ? 'justify-end' : 'justify-start'">
                <div class="max-w-[70%]">
                  <div class="px-4 py-2 rounded-2xl" :class="message.sender === 'me'
                    ? 'bg-green-500 text-white rounded-br-sm'
                    : 'bg-white border border-gray-200 rounded-bl-sm'">
                    <p class="text-sm">{{ message.text }}</p>
                    <p class="text-xs mt-1 opacity-70">{{ formatTime(message.timestamp) }}</p>
                  </div>
                </div>
              </div>

              <!-- Empty State -->
              <div v-if="chatMessages.length === 0" class="text-center py-12">
                <svg class="w-16 h-16 mx-auto text-gray-300 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
                <p class="text-gray-400 text-sm">{{ $t('ShareLocation.noMessages') }}</p>
              </div>
            </div>

            <!-- Chat Input -->
            <div class="p-4 border-t border-gray-200 bg-white">
              <div class="flex gap-2">
                <input v-model="newMessage" type="text" class="input input-bordered flex-1"
                  :placeholder="$t('ShareLocation.typeMessage')" @keyup.enter="sendMessage" />
                <button @click="sendMessage" class="btn btn-primary text-white" :disabled="!newMessage.trim()">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                  </svg>
                </button>
              </div>
              <p class="text-xs text-gray-400 mt-2">{{ $t('ShareLocation.chatNotSaved') }}</p>
            </div>
          </div>
        </div>

        <!-- Right: Info Panel (when no chat) -->
        <div v-else class="space-y-6">
          <!-- How it Works -->
          <div class="card bg-gradient-to-br from-green-500 to-emerald-600 text-white shadow-lg">
            <div class="card-body">
              <h2 class="card-title mb-4">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ $t('ShareLocation.howItWorks') }}
              </h2>
              <div class="space-y-4">
                <div class="flex gap-3">
                  <div
                    class="flex-shrink-0 w-8 h-8 bg-white/20 rounded-full flex items-center justify-center font-bold">
                    1</div>
                  <p class="text-sm">{{ $t('ShareLocation.step1') }}</p>
                </div>
                <div class="flex gap-3">
                  <div
                    class="flex-shrink-0 w-8 h-8 bg-white/20 rounded-full flex items-center justify-center font-bold">
                    2</div>
                  <p class="text-sm">{{ $t('ShareLocation.step2') }}</p>
                </div>
                <div class="flex gap-3">
                  <div
                    class="flex-shrink-0 w-8 h-8 bg-white/20 rounded-full flex items-center justify-center font-bold">
                    3</div>
                  <p class="text-sm">{{ $t('ShareLocation.step3') }}</p>
                </div>
                <div class="flex gap-3">
                  <div
                    class="flex-shrink-0 w-8 h-8 bg-white/20 rounded-full flex items-center justify-center font-bold">
                    4</div>
                  <p class="text-sm">{{ $t('ShareLocation.step4') }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Tips -->
          <div class="card bg-yellow-50 border border-yellow-200 shadow-md">
            <div class="card-body">
              <h3 class="font-bold text-yellow-800 flex items-center gap-2">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                </svg>
                {{ $t('ShareLocation.tips') }}
              </h3>
              <ul class="text-sm text-gray-700 space-y-2 mt-2">
                <li class="flex gap-2">
                  <span>•</span>
                  <span>{{ $t('ShareLocation.tip1') }}</span>
                </li>
                <li class="flex gap-2">
                  <span>•</span>
                  <span>{{ $t('ShareLocation.tip2') }}</span>
                </li>
                <li class="flex gap-2">
                  <span>•</span>
                  <span>{{ $t('ShareLocation.tip3') }}</span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom styles if needed */
</style>
