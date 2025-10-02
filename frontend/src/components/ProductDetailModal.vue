<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { GroupedRecyclableItem, ShopInfo } from '../types/recycle_waste'

const showModal = ref(false)
const selectedProduct = ref<GroupedRecyclableItem | null>(null)
const modalId = ref('product-detail-modal')

// Enhanced shops with images
const enhancedShops = ref<ShopInfo[]>([])

// Function to open the modal
const openModal = (product: GroupedRecyclableItem) => {
  selectedProduct.value = product
  showModal.value = true
  // Use the shop data as provided in the grouped response
  enhancedShops.value = product.shops || []
}

// Function to close the modal
const closeModal = () => {
  showModal.value = false
  selectedProduct.value = null
}

// Handle Escape key to close modal
const handleEscKey = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && showModal.value) {
    closeModal()
  }
}

// Close modal when clicking outside
const closeOnOutsideClick = (e: Event) => {
  if (e.target === e.currentTarget) {
    closeModal()
  }
}

// Format date and time to Thai format (e.g., 1 มกราคม 2567 เวลา 10:30 น.)
const formatDateThai = (dateString: string | undefined): string => {
  if (!dateString) return '';
  
  try {
    const date = new Date(dateString);
    
    // Check if the date is valid
    if (isNaN(date.getTime())) {
      return dateString; // Return original if parsing fails
    }
    
    const day = date.getDate();
    const months = [
      'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน',
      'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม'
    ];
    const month = months[date.getMonth()];
    const year = date.getFullYear() + 543; // Convert to Buddhist Era
    
    // Format time
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    
    return `${day} ${month} ${year} เวลา ${hours}:${minutes} น.`;
  } catch (error) {
    console.error('Error formatting date:', error);
    return dateString;
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleEscKey)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscKey)
})

// Expose methods to parent components
defineExpose({
  openModal,
  closeModal
})
</script>

<template>
  <!-- DaisyUI Modal -->
  <dialog :id="modalId" class="modal" :class="{ 'modal-open': showModal }" @click="closeOnOutsideClick">
    <div
      class="modal-box max-w-4xl max-h-[90vh] overflow-y-auto bg-white rounded-xl shadow-2xl border border-green-200">
      <!-- Modal header -->
      <div class="flex flex-col md:flex-row md:items-start gap-4 mb-6 pb-4 border-b border-green-200">
        <div class="flex-shrink-0">
          <img v-if="selectedProduct?.url" :src="selectedProduct.url" alt="Product Image"
            class="w-32 h-32 object-cover rounded-lg border-2 border-green-200" />
          <div v-else
            class="bg-gray-200 border-2 border-dashed rounded-xl w-32 h-32 flex items-center justify-center text-gray-500">
            No Image
          </div>
        </div>

        <div class="flex-1">
          <div class="flex flex-wrap items-start justify-between gap-2">
            <h3 class="text-2xl font-bold text-green-800 flex-1">
              {{ selectedProduct?.name }}
            </h3>
            <button @click="closeModal" class="btn btn-sm btn-circle btn-ghost text-green-800 hover:bg-green-100"
              aria-label="Close">
              ✕
            </button>
          </div>

          <div class="mt-2 flex flex-wrap items-center gap-3">
            <div class="badge badge-outline border-green-700 text-green-700 bg-green-50 px-3 py-1">
              {{ selectedProduct?.category }}
            </div>
            <div class="text-lg font-semibold text-green-700">
              {{ selectedProduct?.price }} บาท/กก.
            </div>
          </div>

          <div class="mt-2 text-sm text-gray-600">
            อัปเดตล่าสุด: {{ formatDateThai(selectedProduct?.last_update) }}
          </div>
        </div>
      </div>

      <!-- Shops offering this product -->
      <div class="mb-6">
        <h4 class="text-lg font-semibold text-green-800 mb-4 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z"
              clip-rule="evenodd" />
          </svg>
          ร้านค้าที่ให้ราคา
        </h4>

        <div v-if="enhancedShops && enhancedShops.length > 0" class="space-y-3">
          <div v-for="(shop, index) in enhancedShops" :key="index"
            class="flex items-center justify-between p-4 bg-green-50 rounded-lg border border-green-200 hover:bg-green-100 transition-colors">
            <div class="flex items-center gap-3">
              <div class="avatar">
                <div
                  class="w-12 h-12 rounded-full bg-green-200 flex items-center justify-center text-green-800 font-bold overflow-hidden">
                  <img v-if="shop.shop_image_url" :src="shop.shop_image_url" :alt="shop.shop_name"
                    class="w-full h-full object-cover" />
                  <span v-else class="flex items-center justify-center w-full h-full">
                    {{ shop.shop_name ? shop.shop_name.charAt(0).toUpperCase() : '?' }}
                  </span>
                </div>
              </div>
              <div>
                <h5 class="font-medium text-gray-900">{{ shop.shop_name || 'ร้านค้าไม่ระบุ' }}</h5>
              </div>
            </div>
            <div class="text-right">
              <div class="font-bold text-green-700">{{ selectedProduct?.price }} บาท/กก.</div>
              <div class="text-xs text-gray-500">ราคาเดียวกันทุกร้าน</div>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-8 text-gray-500">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-400" fill="none"
            viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="mt-2">ไม่พบข้อมูลร้านค้า</p>
        </div>
      </div>

      <!-- Additional info -->
      <div class="text-sm text-gray-600 border-t border-green-200 pt-4">
        <h5 class="font-medium text-green-800 mb-2">หมายเหตุ</h5>
        <p class="mb-1">• ราคานี้เป็นราคาต่อหน่วยน้ำหนัก (กิโลกรัม)</p>
        <p>• กรุณาตรวจสอบรายละเอียดเพิ่มเติมกับร้านค้าที่คุณเลือก</p>
      </div>
    </div>

    <!-- Close button outside modal box -->
    <form class="modal-backdrop" @click="closeModal">
      <button type="button">close</button>
    </form>
  </dialog>
</template>
