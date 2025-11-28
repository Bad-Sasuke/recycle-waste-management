<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue';
import { useRoute } from 'vue-router';
import {
  IconScale,
  IconX,
  IconCheck,
  IconSearch,
  IconFileText,
  IconBox,
  IconBottle,
  IconGlass,
  IconBarrel,
  IconTool,
  IconRecycle,
  IconCash,
  IconPencil,
  IconTrash
} from '@tabler/icons-vue';
import { useShopStore } from '../stores/shop';
import { useWastesStore } from '../stores/wastes';
import { useWastePurchaseStore } from '../stores/wastePurchase';
import { useUsersStore } from '../stores/users';
import { storeToRefs } from 'pinia';
import config from '../config';
import gsap from 'gsap';
import { useRouter } from 'vue-router';

const route = useRoute();
const shopStore = useShopStore();
const router = useRouter();
const showSuccessModal = ref(false);
const wastesStore = useWastesStore();
const wastePurchaseStore = useWastePurchaseStore();
const usersStore = useUsersStore();
const { wastes, pagination } = storeToRefs(wastesStore);
const { items, isBottomSheetExpanded, grandTotalWeight } = storeToRefs(wastePurchaseStore);

const webAPI = config.webAPI;
const vatRate = config.vatRate;
const showPaymentModal = ref(false);
const isProcessingPayment = ref(false);

// Animate totals when expanding/collapsing
watch(isBottomSheetExpanded, async (isExpanded) => {
  await nextTick();
  const totalsEl = document.querySelector('.header-totals');
  if (totalsEl) {
    if (isExpanded) {
      // Animate out (hide)
      gsap.to(totalsEl, {
        opacity: 0,
        x: 20,
        duration: 0.3,
        ease: 'power2.out'
      });
    } else {
      // Animate in (show)
      gsap.fromTo(totalsEl,
        { opacity: 0, x: -20 },
        { opacity: 1, x: 0, duration: 0.3, ease: 'power2.out' }
      );
    }
  }
});

// Format number with commas and 2 decimal places
const formatCurrency = (value: number): string => {
  return value.toLocaleString('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
};


// Price calculations
const totalPrice = computed(() => {
  return items.value.reduce((sum, item) => {
    const itemPrice = (item.price || 0) * item.totalWeight;
    return sum + itemPrice;
  }, 0);
});

const vatAmount = computed(() => {
  return totalPrice.value * vatRate;
});

const grandTotal = computed(() => {
  return totalPrice.value + vatAmount.value;
});

const openPaymentModal = () => {
  showPaymentModal.value = true;
};

const closePaymentModal = () => {
  showPaymentModal.value = false;
};

const closeSuccessModal = () => {
  showSuccessModal.value = false;
  router.push('/nearby-customers');
};



// Edit Modal State
const showEditModal = ref(false);
const editingItem = ref<any>(null);
const reduceWeight = ref('');

const remainingWeight = computed(() => {
  if (!editingItem.value) return 0;
  const reduce = calculateTotal(reduceWeight.value);
  return Math.max(0, editingItem.value.totalWeight - reduce);
});

const openEditModal = (item: any) => {
  editingItem.value = item;
  reduceWeight.value = '';
  showEditModal.value = true;
};

const closeEditModal = () => {
  showEditModal.value = false;
  editingItem.value = null;
  reduceWeight.value = '';
};

const confirmEdit = () => {
  if (!editingItem.value) return;

  const reduce = calculateTotal(reduceWeight.value);
  if (reduce <= 0) {
    closeEditModal();
    return;
  }

  const newWeight = remainingWeight.value;
  if (newWeight <= 0) {
    // Remove item if weight becomes 0
    wastePurchaseStore.removeItem(editingItem.value.id);
  } else {
    // Update weight and expression
    wastePurchaseStore.updateItemWeight(
      editingItem.value.id,
      newWeight,
      `${editingItem.value.expression} - ${reduceWeight.value}`
    );
  }

  closeEditModal();
};


// Helper to map category/name to icon and color
const getWasteStyle = (name: string, category: string) => {
  const n = name.toLowerCase();
  const c = category?.toLowerCase() || '';

  if (n.includes('กระดาษ') || c.includes('paper')) {
    return { color: 'bg-amber-100 text-amber-800', iconComponent: IconFileText };
  } else if (n.includes('ลัง') || c.includes('cardboard')) {
    return { color: 'bg-orange-100 text-orange-800', iconComponent: IconBox };
  } else if (n.includes('pet') || n.includes('ขวดใส') || c.includes('plastic')) {
    return { color: 'bg-blue-100 text-blue-800', iconComponent: IconBottle };
  } else if (n.includes('แก้ว') || c.includes('glass')) {
    return { color: 'bg-emerald-100 text-emerald-800', iconComponent: IconGlass };
  } else if (n.includes('กระป๋อง') || c.includes('can') || c.includes('aluminium')) {
    return { color: 'bg-gray-100 text-gray-800', iconComponent: IconBarrel };
  } else if (n.includes('เหล็ก') || c.includes('metal') || c.includes('iron')) {
    return { color: 'bg-red-100 text-red-800', iconComponent: IconTool };
  } else {
    return { color: 'bg-green-100 text-green-800', iconComponent: IconRecycle };
  }
};

const searchQuery = ref('');
const selectedCategory = ref<string | null>(null);
const currentPage = ref(1);
const itemsPerPage = 5;
const allCategories = ref<string[]>([]); // Store all categories to prevent them from disappearing

const wasteTypes = computed(() => {
  return wastes.value.map(w => {
    const style = getWasteStyle(w.name || '', w.category || '');
    return {
      id: w.waste_id,
      name: w.name,
      price: w.price,
      category: w.category,
      ...style
    };
  });
});

// Use allCategories instead of computing from filtered data
const categories = computed(() => allCategories.value);

// Client-side filtering for search only (category is now handled by API)
const filteredWasteTypes = computed(() => {
  let result = wasteTypes.value;

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(w =>
      w.name?.toLowerCase().includes(query) ||
      w.category?.toLowerCase().includes(query)
    );
  }

  return result;
});

// Animate cards when they appear
const animateCards = async () => {
  await nextTick();

  const cards = document.querySelectorAll('.waste-card');
  if (cards.length === 0) return;

  // Reset and add animation class to each card with stagger
  cards.forEach((card, index) => {
    const element = card as HTMLElement;
    // Reset initial state
    element.style.opacity = '0';
    element.style.transform = 'translateY(30px) scale(0.95)';

    // Trigger animation with stagger delay
    setTimeout(() => {
      element.style.transition = 'all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1)';
      element.style.opacity = '1';
      element.style.transform = 'translateY(0) scale(1)';
    }, index * 80 + 100);
  });
};

// Fetch all categories once (without any filters)
const fetchAllCategories = async () => {
  if (!shopStore.shop?.shop_id) return;

  try {
    // Fetch all items without category filter to get full category list
    // Temporarily store current wastes and pagination to restore after fetching all categories
    const currentWastes = wastesStore.wastes;
    const currentPagination = wastesStore.pagination;

    await wastesStore.fetchWastes(1, 1000, shopStore.shop.shop_id); // Fetch a large number to get all categories

    // Extract unique categories
    const cats = new Set(wastesStore.wastes.map(w => w.category).filter(Boolean) as string[]);
    allCategories.value = Array.from(cats);

    // Restore original wastes and pagination
    wastesStore.wastes = currentWastes;
    wastesStore.pagination = currentPagination;

  } catch (error) {
    console.error('Error fetching categories:', error);
  }
};

// Fetch data from API
const loadWastes = async (page: number = 1) => {
  if (!shopStore.shop?.shop_id) return;

  await wastesStore.fetchWastes(
    page,
    itemsPerPage,
    shopStore.shop.shop_id,
    selectedCategory.value || undefined
  );
  currentPage.value = pagination.value.page;

  // Trigger animation after data loads
  await animateCards();
};

// Reset to page 1 when filters change
watch([searchQuery, selectedCategory], () => {
  loadWastes(1);
});

const selectedType = ref<any>(null);
const weightInput = ref('');

onMounted(async () => {
  // Check for request_id in query params
  const requestId = (route.query.request_id || route.query.customer_request_id) as string;
  if (requestId) {
    wastePurchaseStore.setRequestId(requestId);
  }

  // Ensure we have shop info (only for shop owners/admins, not regular users)
  if (!shopStore.shop && usersStore.user?.role !== 'user') {
    await shopStore.checkUserShop();
  }

  // Fetch all categories first
  await fetchAllCategories();

  // Load first page
  await loadWastes(1);
});

const openWeighing = (type: any) => {
  selectedType.value = type;
  weightInput.value = '';
};

const closeWeighing = () => {
  selectedType.value = null;
};

const calculateTotal = (expression: string) => {
  if (!expression) return 0;
  try {
    // Safe evaluation for simple addition
    const cleanExpr = expression.replace(/[^0-9+\.]/g, '');
    const parts = cleanExpr.split('+').map(p => parseFloat(p.trim()));
    if (parts.some(p => isNaN(p))) return 0;
    return parts.reduce((a, b) => a + b, 0);
  } catch {
    return 0;
  }
};

const currentTotal = computed(() => calculateTotal(weightInput.value));
const currentPrice = computed(() => {
  if (!selectedType.value || !selectedType.value.price) return 0;
  return currentTotal.value * selectedType.value.price;
});

const addItem = () => {
  if (!selectedType.value || currentTotal.value <= 0) return;

  wastePurchaseStore.addItem({
    ...selectedType.value,
    waste_id: selectedType.value.id,
    expression: weightInput.value,
    totalWeight: currentTotal.value,
    id: Date.now().toString()
  });

  closeWeighing();
};

const confirmPayment = async () => {
  if (!shopStore.shop?.shop_id) return;

  isProcessingPayment.value = true;
  try {
    const payload = {
      shop_id: shopStore.shop.shop_id,
      payment_method: 'cash',
      customer_request_id: (route.query.request_id || route.query.customer_request_id) as string || '',
      items: items.value.map(item => ({
        shop_id: shopStore.shop?.shop_id || '',
        waste_id: item.waste_id,
        name: item.name,
        category: item.category || '',
        weight: item.totalWeight,
        unit_price: item.price || 0,
        price: (item.price || 0) * item.totalWeight
      }))
    };

    const response = await fetch(`${webAPI}/api/receipts`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${usersStore.jwt}`
      },
      body: JSON.stringify(payload)
    });

    if (!response.ok) throw new Error('Failed to create receipt');

    // Clear items and close modal
    wastePurchaseStore.clearItems();
    closePaymentModal();

    // Show success modal
    showSuccessModal.value = true;

  } catch (error) {
    console.error(error);
    alert('เกิดข้อผิดพลาดในการบันทึกรายการ');
  } finally {
    isProcessingPayment.value = false;
  }
};

const removeItem = (id: string) => {
  wastePurchaseStore.removeItem(id);
};

// Pagination controls
const nextPage = async () => {
  if (currentPage.value < pagination.value.total_pages) {
    await loadWastes(currentPage.value + 1);
  }
};

const prevPage = async () => {
  if (currentPage.value > 1) {
    await loadWastes(currentPage.value - 1);
  }
};

const goToPage = async (page: number) => {
  await loadWastes(page);
};

</script>

<template>
  <div class="min-h-screen bg-base-200 pb-20 font-sans">
    <!-- Header -->
    <div class="bg-white shadow-sm px-4 py-3 sticky top-0 z-20 flex items-center justify-between">
      <h1 class="text-xl font-bold text-gray-800 flex items-center gap-2">
        <IconScale class="text-green-600" />
        รับซื้อขยะ
      </h1>
      <div class="badge badge-primary badge-lg font-bold text-white shadow-sm">
        {{ grandTotalWeight }} kg
      </div>
    </div>

    <!-- Search & Filters -->
    <div class="px-4 pt-4 pb-2 space-y-3">
      <!-- Search Input -->
      <div class="relative">
        <input v-model="searchQuery" type="text" placeholder="ค้นหาประเภทขยะ..."
          class="input input-bordered w-full pl-10 bg-white shadow-sm focus:ring-2 focus:ring-green-500 focus:border-green-500" />
        <div class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
          <IconSearch size="20" />
        </div>
        <button v-if="searchQuery" @click="searchQuery = ''"
          class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
          <IconX size="16" />
        </button>
      </div>

      <!-- Category Chips -->
      <div class="flex gap-2 overflow-x-auto pb-2 no-scrollbar">
        <button @click="selectedCategory = null" class="btn btn-sm rounded-full whitespace-nowrap"
          :class="!selectedCategory ? 'btn-neutral text-white' : 'btn-ghost bg-white'">
          ทั้งหมด
        </button>
        <button v-for="cat in categories" :key="cat" @click="selectedCategory = cat"
          class="btn btn-sm rounded-full whitespace-nowrap"
          :class="selectedCategory === cat ? 'btn-neutral text-white' : 'btn-ghost bg-white'">
          {{ cat }}
        </button>
      </div>
    </div>

    <!-- Waste Grid -->
    <div class="p-4 grid grid-cols-2 gap-4 min-h-[300px] content-start">
      <button v-for="type in filteredWasteTypes" :key="type.id" @click="openWeighing(type)"
        class="waste-card relative overflow-hidden bg-gradient-to-br from-white to-gray-50 rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 active:scale-[0.98] border border-gray-100 group">
        <!-- Decorative background pattern -->
        <div class="absolute inset-0 opacity-[0.03] group-hover:opacity-[0.06] transition-opacity">
          <div class="absolute inset-0" :class="type.color.split(' ')[0]"></div>
        </div>

        <div class="relative p-5 flex flex-col items-center">
          <!-- Icon with glass effect -->
          <div class="relative mb-4">
            <div
              :class="`w-20 h-20 rounded-2xl flex items-center justify-center shadow-lg ${type.color} group-hover:scale-110 group-hover:rotate-3 transition-all duration-300`">
              <component :is="type.iconComponent" :size="48" :stroke-width="1.5" class="filter drop-shadow-sm" />
            </div>
          </div>

          <!-- Item name -->
          <h3 class="font-bold text-base text-gray-800 mb-3 group-hover:text-green-600 transition-colors">
            {{ type.name }}
          </h3>

          <!-- Hover indicator -->
          <div class="opacity-0 group-hover:opacity-100 transition-opacity">
            <div class="flex items-center gap-1 text-xs text-green-600 font-semibold">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
              </svg>
              <span>เพิ่มรายการ</span>
            </div>
          </div>
        </div>

        <!-- Bottom accent line -->
        <div
          class="absolute bottom-0 left-0 right-0 h-1 bg-gradient-to-r from-transparent via-green-500 to-transparent opacity-0 group-hover:opacity-100 transition-opacity">
        </div>
      </button>
    </div>

    <!-- Pagination Controls -->
    <div v-if="pagination.total_pages > 1" class="flex justify-center items-center gap-2 px-4 pb-4 mb-20">
      <button @click="prevPage" class="btn btn-sm btn-circle btn-ghost bg-white shadow-sm"
        :disabled="currentPage === 1">
        «
      </button>

      <div class="flex gap-1">
        <button v-for="page in pagination.total_pages" :key="page" @click="goToPage(page)" class="btn btn-sm btn-circle"
          :class="currentPage === page ? 'btn-primary text-white' : 'btn-ghost bg-white shadow-sm'"
          v-show="Math.abs(page - currentPage) < 2 || page === 1 || page === pagination.total_pages">
          {{ page }}
        </button>
      </div>

      <button @click="nextPage" class="btn btn-sm btn-circle btn-ghost bg-white shadow-sm"
        :disabled="currentPage === pagination.total_pages">
        »
      </button>
    </div>

    <!-- Backdrop Overlay -->
    <div v-if="items.length > 0 && isBottomSheetExpanded" @click="isBottomSheetExpanded = false"
      class="fixed inset-0 bg-black/40 backdrop-blur-sm transition-all duration-300 z-40">
    </div>

    <!-- Floating Bottom Sheet Summary -->
    <div v-if="items.length > 0"
      class="fixed bottom-0 left-0 right-0 bg-white rounded-t-3xl shadow-2xl transition-all duration-300 z-50"
      :class="isBottomSheetExpanded ? 'h-[70vh]' : 'h-24'">

      <!-- Handle Bar -->
      <div @click="isBottomSheetExpanded = !isBottomSheetExpanded"
        class="w-full py-3 flex flex-col items-center cursor-pointer border-b border-gray-100">
        <div class="w-12 h-1.5 bg-gray-300 rounded-full mb-2"></div>
        <div class="flex items-center justify-between w-full px-4">
          <div class="flex items-center gap-2">
            <h3 class="text-sm font-bold text-gray-700">รายการล่าสุด</h3>
            <span class="bg-green-500 text-white text-xs font-bold px-2 py-0.5 rounded-full">{{ items.length }}</span>
          </div>
          <div class="flex items-center gap-3">
            <!-- Show totals only when collapsed -->
            <template v-if="!isBottomSheetExpanded">
              <div class="header-totals flex items-center gap-3">
                <div class="text-right">
                  <div class="text-xs text-gray-500">น้ำหนักรวม</div>
                  <div class="text-base font-bold text-green-600">{{ grandTotalWeight }} kg</div>
                </div>
                <div class="w-px h-8 bg-gray-300"></div>
                <div class="text-right">
                  <div class="text-xs text-gray-500">ยอดชำระ</div>
                  <div class="text-base font-bold text-blue-600">฿{{ formatCurrency(grandTotal) }}</div>
                </div>
              </div>
            </template>
            <svg class="w-5 h-5 text-gray-400 transition-transform duration-300"
              :class="isBottomSheetExpanded ? 'rotate-180' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Items List (Scrollable) -->
      <div class="overflow-y-auto px-4 py-3 space-y-2"
        :class="isBottomSheetExpanded ? 'h-[calc(70vh-80px)]' : 'h-0 opacity-0'">
        <div v-for="item in items" :key="item.id" @click="openEditModal(item)"
          class="bg-gradient-to-r from-white to-gray-50 p-3 rounded-xl shadow-sm flex justify-between items-center border-l-4 hover:shadow-md transition-shadow cursor-pointer select-none active:scale-[0.98] active:bg-gray-100"
          :class="getWasteStyle(item.name, item.category || '').color.replace('bg-', 'border-').split(' ')[0]">
          <div class="flex items-center gap-3">
            <div
              :class="`w-12 h-12 rounded-xl flex items-center justify-center shadow-md ${getWasteStyle(item.name, item.category || '').color}`">
              <component :is="getWasteStyle(item.name, item.category || '').iconComponent" :size="24"
                :stroke-width="2" />
            </div>
            <div>
              <div class="font-bold text-gray-800">{{ item.name }}</div>
              <div class="text-xs text-gray-500 font-mono bg-gray-100 px-2 py-0.5 rounded mt-1">{{ item.expression }}
              </div>
              <div class="text-xs text-green-600 mt-1 font-semibold">
                ฿{{ formatCurrency(item.price || 0) }}/kg
              </div>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <div class="text-right">
              <div>
                <span class="font-bold text-lg text-gray-700">{{ item.totalWeight }}</span>
                <span class="text-xs text-gray-400 ml-1">kg</span>
              </div>
              <div class="text-sm font-bold text-blue-600">
                ฿{{ formatCurrency((item.price || 0) * item.totalWeight) }}
              </div>
            </div>
            <div class="text-gray-300">
              <IconPencil :size="20" />
            </div>
          </div>
        </div>
      </div>

      <!-- Payment Button Footer with Price Summary -->
      <div v-if="isBottomSheetExpanded" class="border-t border-gray-100 bg-white sticky bottom-0">
        <!-- Price Summary -->
        <div class="px-4 pt-3 pb-2 bg-gray-50 space-y-2">
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">ยอดรวม:</span>
            <span class="font-bold">฿{{ formatCurrency(totalPrice) }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">VAT {{ (vatRate * 100).toFixed(0) }}%:</span>
            <span class="font-bold">฿{{ formatCurrency(vatAmount) }}</span>
          </div>
          <div class="flex justify-between text-base border-t border-gray-200 pt-2">
            <span class="font-bold text-gray-800">ยอดชำระทั้งหมด:</span>
            <span class="font-bold text-blue-600 text-lg">฿{{ formatCurrency(grandTotal) }}</span>
          </div>
        </div>
        <!-- Payment Button -->
        <div class="p-4">
          <button v-if="usersStore.user?.role === 'moderator' || usersStore.user?.role === 'admin'"
            @click="openPaymentModal" class="btn btn-primary w-full text-white text-lg shadow-lg gap-2">
            <IconCash :size="24" />
            ชำระเงิน
          </button>
          <div v-else class="alert alert-warning">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
              viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <span>เฉพาะเจ้าของร้านเท่านั้นที่สามารถสร้างใบเสร็จได้</span>
          </div>
        </div>
        <div class="h-4"></div>
      </div>
    </div>

    <!-- Empty State (when no items and bottom sheet hidden) -->
    <div v-else class="text-center text-gray-400 mt-8 mb-20">
      <p>ยังไม่มีรายการ</p>
      <p class="text-sm">เลือกประเภทขยะเพื่อเริ่มชั่งน้ำหนัก</p>
    </div>

    <!-- Weighing Modal (Bottom Sheet style) -->
    <div v-if="selectedType" class="fixed inset-0 z-50 flex items-end justify-center bg-black/60 backdrop-blur-sm"
      @click.self="closeWeighing">
      <div class="bg-white w-full max-w-md rounded-t-2xl p-6 shadow-2xl animate-slide-up">
        <div class="flex justify-between items-center mb-6">
          <div class="flex items-center gap-3">
            <div :class="`w-12 h-12 rounded-full flex items-center justify-center text-2xl ${selectedType.color}`">
              {{ selectedType.icon }}
            </div>
            <div>
              <h3 class="text-lg font-bold text-gray-800">{{ selectedType.name }}</h3>
              <p class="text-sm text-gray-500">ระบุน้ำหนัก (เช่น 34 + 3)</p>
            </div>
          </div>
          <button @click="closeWeighing" class="btn btn-circle btn-ghost btn-sm bg-gray-100">
            <IconX size="20" />
          </button>
        </div>

        <div class="form-control mb-6">
          <div class="relative">
            <input v-model="weightInput" type="text" placeholder="0"
              class="input input-bordered w-full text-3xl font-mono h-16 text-right pr-12 focus:outline-none focus:border-green-500 focus:ring-2 focus:ring-green-200"
              autofocus />
            <div class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 text-sm font-bold">KG</div>
          </div>
          <div class="flex justify-between items-center mt-2 px-1">
            <span class="text-xs text-gray-400">รองรับการบวกเลข (เช่น 5+2+1)</span>
            <div class="text-right text-green-600 font-bold font-mono text-lg">
              = {{ currentTotal }} kg
            </div>
          </div>
        </div>

        <button @click="addItem" class="btn btn-primary w-full btn-md text-white shadow-lg shadow-green-200"
          :disabled="currentTotal <= 0">
          <IconCheck />
          บันทึกรายการ
        </button>
        <div class="h-4"></div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
      @click.self="closeEditModal">
      <div class="bg-white w-full max-w-md rounded-2xl p-6 shadow-2xl mx-4 animate-slide-up">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-xl font-bold text-gray-800">แก้ไขรายการ</h3>
          <button @click="closeEditModal" class="btn btn-circle btn-ghost btn-sm bg-gray-100">
            <IconX :size="20" />
          </button>
        </div>

        <div v-if="editingItem" class="mb-6">
          <div class="flex items-center gap-3 mb-4 p-3 bg-gray-50 rounded-xl">
            <div
              :class="`w-12 h-12 rounded-xl flex items-center justify-center shadow-sm ${getWasteStyle(editingItem.name, editingItem.category || '').color}`">
              <component :is="getWasteStyle(editingItem.name, editingItem.category || '').iconComponent" :size="24"
                :stroke-width="2" />
            </div>
            <div>
              <div class="font-bold text-gray-800">{{ editingItem.name }}</div>
              <div class="text-sm text-gray-500">น้ำหนักปัจจุบัน: <span class="font-bold text-green-600">{{
                editingItem.totalWeight }} kg</span></div>
            </div>
          </div>

          <div class="form-control mb-4">
            <label class="label">
              <span class="label-text font-semibold">ระบุน้ำหนักที่ต้องการลบออก</span>
            </label>
            <div class="relative">
              <input v-model="reduceWeight" type="tel" inputmode="decimal" placeholder="0"
                class="input input-bordered w-full text-2xl font-mono h-14 text-right pr-12 focus:outline-none focus:border-red-500 focus:ring-2 focus:ring-red-200"
                autofocus />
              <div class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 text-sm font-bold">KG</div>
            </div>

            <!-- Suggestion Chips -->
            <div class="flex gap-2 mt-3 overflow-x-auto pb-2 scrollbar-hide">
              <button v-for="amount in [0.5, 1, 2, 5, 10]" :key="amount"
                @click="reduceWeight = (parseFloat(reduceWeight || '0') + amount).toString()"
                class="btn btn-xs btn-outline border-gray-300 text-gray-600 hover:bg-red-50 hover:text-red-600 hover:border-red-200">
                -{{ amount }}
              </button>
              <button @click="reduceWeight = ''" class="btn btn-xs btn-ghost text-gray-400">
                ล้าง
              </button>
            </div>
            <div class="flex justify-between items-center mt-2 px-1">
              <span class="text-xs text-gray-400">รองรับการบวกเลข (เช่น 1+0.5)</span>
              <div class="text-right">
                <span class="text-sm text-gray-500">คงเหลือ: </span>
                <span class="font-bold font-mono text-lg"
                  :class="remainingWeight <= 0 ? 'text-red-500' : 'text-green-600'">
                  {{ remainingWeight }} kg
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex gap-3">
          <button @click="wastePurchaseStore.removeItem(editingItem.id); closeEditModal()"
            class="btn btn-error btn-outline flex-none">
            <IconTrash :size="20" />
          </button>
          <button @click="confirmEdit" class="btn btn-primary flex-1 text-white"
            :class="remainingWeight <= 0 ? 'btn-error' : 'btn-primary'"
            :disabled="remainingWeight > 0 && (!reduceWeight || calculateTotal(reduceWeight) <= 0)">
            {{ remainingWeight <= 0 ? 'ลบรายการ' : 'บันทึกการแก้ไข' }} </button>
        </div>
      </div>
    </div>

    <!-- Payment Modal -->
    <div v-if="showPaymentModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
      @click.self="closePaymentModal">
      <div class="bg-white w-full max-w-md rounded-2xl p-6 shadow-2xl mx-4 animate-slide-up">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-xl font-bold text-gray-800">ยืนยันการชำระเงิน</h3>
          <button @click="closePaymentModal" class="btn btn-circle btn-ghost btn-sm bg-gray-100">
            <IconX :size="20" />
          </button>
        </div>

        <!-- Summary -->
        <div class="bg-gray-50 rounded-xl p-4 mb-6 space-y-3">
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">จำนวนรายการ:</span>
            <span class="font-bold">{{ items.length }} รายการ</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">น้ำหนักรวม:</span>
            <span class="font-bold text-green-600">{{ grandTotalWeight }} kg</span>
          </div>
          <div class="border-t border-gray-300 pt-2 space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ยอดรวม:</span>
              <span class="font-semibold">฿{{ formatCurrency(totalPrice) }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">VAT {{ (vatRate * 100).toFixed(0) }}%:</span>
              <span class="font-semibold">฿{{ formatCurrency(vatAmount) }}</span>
            </div>
            <div class="flex justify-between text-lg border-t border-gray-300 pt-2">
              <span class="font-bold text-gray-800">ยอดชำระทั้งหมด:</span>
              <span class="font-bold text-blue-600">฿{{ formatCurrency(grandTotal) }}</span>
            </div>
          </div>
        </div>

        <!-- Payment Method -->
        <div class="mb-6">
          <label class="label">
            <span class="label-text font-semibold">วิธีชำระเงิน</span>
          </label>
          <div class="p-4 border rounded-xl bg-green-50 border-green-200 flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center text-green-600">
              <IconCash :size="24" />
            </div>
            <div>
              <div class="font-bold text-gray-800">เงินสด (Cash)</div>
              <div class="text-xs text-gray-500">ชำระด้วยเงินสดเท่านั้น</div>
            </div>
            <div class="ml-auto text-green-600">
              <IconCheck :size="20" />
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-3">
          <button @click="closePaymentModal" class="btn btn-ghost flex-1">
            ยกเลิก
          </button>
          <button @click="confirmPayment" class="btn btn-primary flex-1 text-white" :disabled="isProcessingPayment">
            <span v-if="isProcessingPayment" class="loading loading-spinner"></span>
            {{ isProcessingPayment ? 'กำลังบันทึก...' : 'ยืนยันชำระเงิน' }}
          </button>
        </div>
      </div>
    </div>


    <!-- Success Modal -->
    <div v-if="showSuccessModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
      @click.self="closeSuccessModal">
      <div class="bg-white w-full max-w-sm rounded-2xl p-8 shadow-2xl mx-4 animate-slide-up text-center">
        <div class="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6 text-green-600">
          <IconCheck :size="48" stroke-width="3" />
        </div>
        <h3 class="text-2xl font-bold text-gray-800 mb-2">บันทึกสำเร็จ!</h3>
        <p class="text-gray-500 mb-8">รายการซื้อขยะถูกบันทึกเรียบร้อยแล้ว</p>
        <button @click="closeSuccessModal" class="btn btn-primary w-full text-white text-lg h-12 rounded-xl">
          ตกลง
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-slide-up {
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
  }

  to {
    transform: translateY(0);
  }
}
</style>
