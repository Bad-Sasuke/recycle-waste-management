<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue';
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
const activeTab = ref('info');
const editMode = ref(false);
const isSubmitting = ref(false);
const showDeleteModal = ref(false);
const modalType = ref<'success' | 'error'>('success');
const modalMessage = ref('');

// Walk-in Modal States
const showWalkInModal = ref(false);
const walkInCustomerName = ref('');
const walkInCustomerPhone = ref('');
const isCreatingWalkIn = ref(false);

// Loading states
const loadingStocks = ref(false);
const loadingReceipts = ref(false);

// Data สำหรับ Stock และ Receipts
// Data สำหรับ Stock และ Receipts
interface StockItem {
  id: string;
  updated_at: string;
  category: string;
  name: string;
  quantity: number;
  purchase_price: number;
  current_price: number;
  profit: number;
}

interface Receipt {
  id: string;
  created_at: string;
  customer_name: string;
  items_count: number;
  net_total: number;
  status: string;
}

const stocks = ref<StockItem[]>([]);
const receipts = ref<Receipt[]>([]);

// Pagination for Stock
const stockPage = ref(1);
const stockPageSize = 5;
const totalStockPages = computed(() => Math.ceil(stocks.value.length / stockPageSize));
const paginatedStocks = computed(() => stocks.value); // ใช้ข้อมูลจาก server โดยตรง

// Pagination for Receipts
const receiptPage = ref(1);
const receiptPageSize = 5;
const totalReceiptPages = computed(() => Math.ceil(receipts.value.length / receiptPageSize));
const paginatedReceipts = computed(() => receipts.value); // ใช้ข้อมูลจาก server โดยตรง

// Helper to format currency
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('th-TH', {
    style: 'currency',
    currency: 'THB',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(value);
};

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

// Fetch stocks data
const fetchStocks = async () => {
  if (!shopStore.shop?.shop_id) return;

  loadingStocks.value = true;
  try {
    const { fetchStocksByShopID } = await import('@/services/manageShop');
    const response = await fetchStocksByShopID(shopStore.shop.shop_id, stockPage.value, stockPageSize);
    console.log('fetchStocks response:', response);
    if (response.success && response.data) {
      stocks.value = response.data;
      // Update total pages from server response if available
      if (response.total_pages) {
        // totalStockPages will be recomputed automatically from stocks.value.length
      }
    } else {
      console.error('Stock fetch failed:', response.message);
      stocks.value = [];
    }
  } catch (error) {
    console.error('Error fetching stocks:', error);
    stocks.value = [];
  } finally {
    loadingStocks.value = false;
  }
};

// Fetch receipts data
const fetchReceipts = async () => {
  if (!shopStore.shop?.shop_id) return;

  loadingReceipts.value = true;
  try {
    const { fetchReceiptsByShopID } = await import('@/services/manageShop');
    const response = await fetchReceiptsByShopID(shopStore.shop.shop_id, receiptPage.value, receiptPageSize);
    console.log('fetchReceipts response:', response);
    if (response.success && response.data) {
      receipts.value = response.data;
    } else {
      console.error('Receipt fetch failed:', response.message);
      receipts.value = [];
    }
  } catch (error) {
    console.error('Error fetching receipts:', error);
    receipts.value = [];
  } finally {
    loadingReceipts.value = false;
  }
};

// Watch activeTab to fetch data when switching tabs
watch(activeTab, (newTab) => {
  if (newTab === 'stock') {
    fetchStocks();
  } else if (newTab === 'receipts') {
    fetchReceipts();
  }
});

const goToWastePurchase = () => {
  // Reset form and open modal
  walkInCustomerName.value = '';
  walkInCustomerPhone.value = '';
  showWalkInModal.value = true;
};

const handleCreateWalkIn = async () => {
  if (!shopStore.shop?.shop_id) return;

  isCreatingWalkIn.value = true;
  try {
    // Dynamic import to avoid circular dependencies if any, or just good practice
    const { createWalkInRequest } = await import('@/services/customer_request');

    const response = await createWalkInRequest(walkInCustomerName.value, walkInCustomerPhone.value);

    if (response.success && response.customer_request_id) {
      showWalkInModal.value = false;
      router.push({
        path: '/waste-purchase',
        query: { customer_request_id: response.customer_request_id }
      });
    } else {
      showModal('error', response.message || 'เกิดข้อผิดพลาดในการสร้างรายการ');
    }
  } catch (error) {
    console.error('Error creating walk-in request:', error);
    const errorMessage = error instanceof Error ? error.message : 'เกิดข้อผิดพลาดในการสร้างรายการ';
    showModal('error', errorMessage);
  } finally {
    isCreatingWalkIn.value = false;
  }
};

// Watch pagination changes
watch(stockPage, () => {
  if (activeTab.value === 'stock') {
    fetchStocks();
  }
});

watch(receiptPage, () => {
  if (activeTab.value === 'receipts') {
    fetchReceipts();
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
    <div v-else class="max-w-6xl mx-auto">
      <!-- Shop Header -->
      <div class="bg-white rounded-2xl shadow-xl overflow-hidden mb-6">
        <div class="bg-gradient-to-r from-green-500 to-emerald-600 p-6 text-white">
          <div class="flex items-center gap-4">
            <div class="bg-white/20 p-3 rounded-xl">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" stroke-width="2"
                stroke="currentColor" fill="none">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                <path d="M3 21l18 0"></path>
                <path d="M3 7v1a3 3 0 0 0 6 0v-1m0 1a3 3 0 0 0 6 0v-1m0 1a3 3 0 0 0 6 0v-1"></path>
                <path d="M5 21v-14a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v14"></path>
              </svg>
            </div>
            <div>
              <h1 class="text-2xl font-bold">{{ shopStore.shopName || 'จัดการร้าน' }}</h1>
              <p class="opacity-90">{{ shopStore.shopAddress || 'จัดการข้อมูลร้านของคุณ' }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="bg-white rounded-2xl shadow-xl overflow-hidden">
        <div role="tablist" class="tabs tabs-lifted tabs-lg border-b">
          <button role="tab" :class="['tab', { 'tab-active': activeTab === 'info' }]" @click="activeTab = 'info'">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 mr-2" viewBox="0 0 24 24" stroke-width="2"
              stroke="currentColor" fill="none">
              <path d="M12 12h-6v4h6m0 0h6v-4h-6m-6-4h12v-4h-12z"></path>
            </svg>
            ข้อมูลร้าน
          </button>
          <button role="tab" :class="['tab', { 'tab-active': activeTab === 'stock' }]" @click="activeTab = 'stock'">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 mr-2" viewBox="0 0 24 24" stroke-width="2"
              stroke="currentColor" fill="none">
              <path d="M3 3h18v18h-18z M9 8v8M15 8v8"></path>
            </svg>
            สต็อก
          </button>
          <button role="tab" :class="['tab', { 'tab-active': activeTab === 'receipts' }]"
            @click="activeTab = 'receipts'">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 mr-2" viewBox="0 0 24 24" stroke-width="2"
              stroke="currentColor" fill="none">
              <path d="M9 5h-2a2 2 0 0 0 -2 2v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2 -2v-12a2 2 0 0 0 -2 -2h-2"></path>
              <path d="M9 3m0 2a2 2 0 0 1 2 -2h2a2 2 0 0 1 2 2v0a2 2 0 0 1 -2 2h-2a2 2 0 0 1 -2 -2z"></path>
            </svg>
            ใบเสร็จ
          </button>
        </div>

        <div class="p-6 md:p-8">
          <!-- Tab 1: Shop Info -->
          <div v-show="activeTab === 'info'">
            <div v-if="!editMode">
              <div class="flex justify-between items-center mb-6">
                <h2 class="text-xl font-bold">ข้อมูลร้าน</h2>
                <button @click="editMode = true" class="btn btn-primary btn-sm">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none"
                    stroke="currentColor">
                    <path d="M7 7h-1a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-1"></path>
                    <path d="M20.385 6.585a2.1 2.1 0 0 0 -2.97 -2.97l-8.415 8.385v3h3l8.385 -8.415z"></path>
                  </svg>
                  แก้ไข
                </button>
              </div>

              <div class="grid md:grid-cols-2 gap-6">
                <div class="col-span-2 flex justify-center mb-4">
                  <img :src="shopStore.shop?.image_url || 'https://placehold.co/400x300'"
                    class="w-64 h-64 object-cover rounded-xl shadow-md" />
                </div>
                <div>
                  <h3 class="font-semibold text-gray-600">ชื่อร้าน</h3>
                  <p class="text-lg">{{ shopStore.shop?.name || '-' }}</p>
                </div>
                <div>
                  <h3 class="font-semibold text-gray-600">อีเมล</h3>
                  <p class="text-lg">{{ shopStore.shop?.email || '-' }}</p>
                </div>
                <div>
                  <h3 class="font-semibold text-gray-600">เบอร์โทร</h3>
                  <p class="text-lg">{{ shopStore.shop?.phone || '-' }}</p>
                </div>
                <div class="col-span-2">
                  <h3 class="font-semibold text-gray-600">ที่อยู่</h3>
                  <p class="text-lg">{{ shopStore.shop?.address || '-' }}</p>
                </div>
                <div>
                  <h3 class="font-semibold text-gray-600">เวลาเปิด</h3>
                  <p class="text-lg">{{ shopStore.shop?.opening_time || '-' }}</p>
                </div>
                <div>
                  <h3 class="font-semibold text-gray-600">เวลาปิด</h3>
                  <p class="text-lg">{{ shopStore.shop?.closing_time || '-' }}</p>
                </div>
              </div>

              <div class="flex justify-end mt-8 pt-6 border-t">
                <button @click="confirmDelete" class="btn btn-error btn-sm">ลบร้าน</button>
              </div>
            </div>

            <!-- Edit Form (existing form code shortened for brevity) -->
            <div v-else>
              <form @submit.prevent="updateShop" class="space-y-4">
                <div class="flex justify-center mb-4">
                  <div class="form-control">
                    <img :src="previewImage || shopStore.shop?.image_url || 'https://placehold.co/200'"
                      class="w-48 h-48 object-cover rounded-xl mb-2" />
                    <label class="btn btn-outline btn-sm">
                      อัพโหลดรูป
                      <input type="file" class="hidden" accept="image/*" @change="handleImageChange" />
                    </label>
                  </div>
                </div>

                <div class="grid md:grid-cols-2 gap-4">
                  <div class="form-control"><label class="label"><span class="label-text">ชื่อร้าน
                        *</span></label><input v-model="editShopData.name" type="text" class="input input-bordered"
                      required /></div>
                  <div class="form-control"><label class="label"><span class="label-text">อีเมล</span></label><input
                      v-model="editShopData.email" type="email" class="input input-bordered" /></div>
                  <div class="form-control"><label class="label"><span class="label-text">เบอร์โทร</span></label><input
                      v-model="editShopData.phone" type="tel" class="input input-bordered" /></div>
                  <div class="form-control col-span-2"><label class="label"><span class="label-text">ที่อยู่
                        *</span></label><input v-model="editShopData.address" type="text" class="input input-bordered"
                      required /></div>
                  <div class="form-control"><label class="label"><span class="label-text">เวลาเปิด</span></label><input
                      v-model="editShopData.opening_time" type="time" class="input input-bordered" /></div>
                  <div class="form-control"><label class="label"><span class="label-text">เวลาปิด</span></label><input
                      v-model="editShopData.closing_time" type="time" class="input input-bordered" /></div>
                  <div class="form-control col-span-2"><label class="label"><span
                        class="label-text">คำอธิบาย</span></label><textarea v-model="editShopData.description"
                      class="textarea textarea-bordered" rows="3"></textarea></div>
                </div>

                <div class="form-control">
                  <label class="label"><span class="label-text">ตำแหน่งร้าน</span></label>
                  <div id="edit-map" class="w-full h-64 rounded-lg border"></div>
                </div>

                <div class="flex gap-2">
                  <button type="submit" :disabled="isSubmitting" class="btn btn-primary flex-1">บันทึก</button>
                  <button type="button" @click="editMode = false" class="btn btn-outline flex-1">ยกเลิก</button>
                </div>
              </form>
            </div>
          </div>

          <!-- Tab 2: Stock -->
          <div v-show="activeTab === 'stock'">
            <div class="flex justify-between items-center mb-6">
              <h2 class="text-xl font-bold">วัสดุรีไซเคิลในคลัง</h2>
              <!-- Desktop Action Buttons (Hidden on mobile) -->
              <div class="hidden md:flex gap-2">
                <button class="btn btn-primary btn-sm text-white" @click="goToWastePurchase">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="w-4 h-4 mr-1">
                    <path stroke-linecap="round" stroke-linejoin="round"
                      d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  ชั่งชำระเงิน
                </button>
                <button class="btn btn-error btn-sm text-white">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="w-4 h-4 mr-1">
                    <path stroke-linecap="round" stroke-linejoin="round"
                      d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
                  </svg>
                  นำออกคลัง
                </button>
              </div>
            </div>

            <!-- Mobile Action Buttons (Visible only on mobile) -->
            <div class="grid grid-cols-1 gap-3 mb-6 md:hidden">
              <button class="btn btn-primary w-full text-white" @click="goToWastePurchase">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                  stroke="currentColor" class="w-5 h-5 mr-2">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                ชั่งชำระเงิน
              </button>
              <button class="btn btn-error w-full text-white">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                  stroke="currentColor" class="w-5 h-5 mr-2">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
                </svg>
                นำออกคลัง
              </button>
            </div>
            <!-- TODO: เชื่อม API GET /api/stocks/shop/:shop_id -->
            <div class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr class="bg-base-200">
                    <th>ลำดับ</th>
                    <th class="text-center">อัปเดตล่าสุด</th>
                    <th>หมวดหมู่</th>
                    <th>ชื่อวัสดุ</th>
                    <th class="text-right">จำนวน (กก.)</th>
                    <th class="text-right">ราคา/กก.</th>
                    <th class="text-right">ส่วนต่าง</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(item, index) in paginatedStocks" :key="item.id">
                    <td>{{ (stockPage - 1) * stockPageSize + index + 1 }}</td>
                    <td class="text-center text-sm text-gray-500">
                      {{ new Date(item.updated_at).toLocaleString('th-TH', {
                        dateStyle: 'short', timeStyle: 'short'
                      }) }}
                    </td>
                    <td><span class="badge badge-primary">{{ item.category }}</span></td>
                    <td class="font-semibold">{{ item.name }}</td>
                    <td class="text-right">{{ item.quantity }}</td>
                    <td class="text-right">
                      <div class="flex flex-col items-end">
                        <span class="font-bold text-red-500">ซื้อเข้ามา: {{ formatCurrency(item.purchase_price)
                          }}</span>
                        <span class="font-bold text-green-600">ปัจจุบัน: {{ formatCurrency(item.current_price) }}</span>
                      </div>
                    </td>
                    <td class="text-right font-bold" :class="item.profit >= 0 ? 'text-green-600' : 'text-red-600'">
                      {{ item.profit > 0 ? '+' : '' }}{{ formatCurrency(item.profit) }}
                    </td>
                  </tr>
                  <tr v-if="stocks.length === 0">
                    <td colspan="6" class="text-center text-gray-500 py-8">
                      <span v-if="loadingStocks">กำลังโหลด...</span>
                      <span v-else>ไม่มีข้อมูลสต็อก</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination for Stock -->
            <div v-if="totalStockPages > 1" class="flex justify-center mt-6">
              <div class="join">
                <button class="join-item btn btn-sm" :disabled="stockPage === 1" @click="stockPage--">«</button>
                <button v-for="page in totalStockPages" :key="page" class="join-item btn btn-sm"
                  :class="{ 'btn-active': stockPage === page }" @click="stockPage = page">
                  {{ page }}
                </button>
                <button class="join-item btn btn-sm" :disabled="stockPage === totalStockPages"
                  @click="stockPage++">»</button>
              </div>
            </div>
          </div>

          <!-- Tab 3: Receipts -->
          <div v-show="activeTab === 'receipts'">
            <h2 class="text-xl font-bold mb-6">ประวัติการทำรายการ</h2>
            <!-- TODO: เชื่อม API GET /api/receipts/shop/:shop_id -->
            <div class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr class="bg-base-200">
                    <th>เลขที่ใบเสร็จ</th>
                    <th>วันที่-เวลา</th>
                    <th>ลูกค้า</th>
                    <th class="text-center">จำนวนรายการ</th>
                    <th class="text-right">ยอดรวม</th>
                    <th class="text-center">สถานะ</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="receipt in paginatedReceipts" :key="receipt.id">
                    <td class="font-mono font-semibold text-xs">#{{ receipt.id.substring(0, 8) }}...</td>
                    <td class="text-sm">{{ new Date(receipt.created_at).toLocaleString('th-TH') }}</td>
                    <td>{{ receipt.customer_name }}</td>
                    <td class="text-center">{{ receipt.items_count }}</td>
                    <td class="text-right font-bold text-green-600">{{ formatCurrency(receipt.net_total) }}</td>
                    <td class="text-center">
                      <span class="badge badge-success badge-sm">{{ receipt.status }}</span>
                    </td>
                  </tr>
                  <tr v-if="receipts.length === 0">
                    <td colspan="6" class="text-center text-gray-500 py-8">
                      <span v-if="loadingReceipts">กำลังโหลด...</span>
                      <span v-else>ไม่มีข้อมูลใบเสร็จ</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination for Receipts -->
            <div v-if="totalReceiptPages > 1" class="flex justify-center mt-6">
              <div class="join">
                <button class="join-item btn btn-sm" :disabled="receiptPage === 1" @click="receiptPage--">«</button>
                <button v-for="page in totalReceiptPages" :key="page" class="join-item btn btn-sm"
                  :class="{ 'btn-active': receiptPage === page }" @click="receiptPage = page">
                  {{ page }}
                </button>
                <button class="join-item btn btn-sm" :disabled="receiptPage === totalReceiptPages"
                  @click="receiptPage++">»</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Walk-in Customer Modal -->
    <dialog id="walk-in-modal" class="modal modal-bottom sm:modal-middle" :class="{ 'modal-open': showWalkInModal }">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">ลูกค้า Walk-in</h3>
        <p class="text-sm text-gray-500 mb-4">กรอกข้อมูลลูกค้าที่นำขยะมารีไซเคิลที่ร้าน</p>
        <form @submit.prevent="handleCreateWalkIn">
          <div class="form-control w-full mb-4">
            <label class="label">
              <span class="label-text">ชื่อลูกค้า (ไม่บังคับ)</span>
            </label>
            <input type="text" v-model="walkInCustomerName" placeholder="ระบุชื่อลูกค้า"
              class="input input-bordered w-full" />
          </div>
          <div class="form-control w-full mb-6">
            <label class="label">
              <span class="label-text">เบอร์โทรศัพท์ (ไม่บังคับ)</span>
            </label>
            <input type="tel" v-model="walkInCustomerPhone" placeholder="ระบุเบอร์โทรศัพท์"
              class="input input-bordered w-full" />
          </div>
          <div class="modal-action">
            <button type="button" class="btn" @click="showWalkInModal = false"
              :disabled="isCreatingWalkIn">ยกเลิก</button>
            <button type="submit" class="btn btn-primary" :disabled="isCreatingWalkIn">
              <span v-if="isCreatingWalkIn" class="loading loading-spinner"></span>
              สร้างรายการ
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="showWalkInModal = false">close</button>
      </form>
    </dialog>

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
