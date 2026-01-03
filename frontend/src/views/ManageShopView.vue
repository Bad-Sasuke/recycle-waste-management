<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue';
import { useShopStore } from '@/stores/shop';
import { useRouter, useRoute } from 'vue-router';
import type { UpdateShopRequest } from '@/types/shop';
import { useI18n } from 'vue-i18n';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import ReceiptDetailModal from '@/components/ReceiptDetailModal.vue';
import EmployeeManagement from '@/components/shop/EmployeeManagement.vue';
import type { ReceiptDetailData } from '@/components/ReceiptDetailModal.vue';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js';
import { Line, Doughnut } from 'vue-chartjs';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const shopStore = useShopStore();

// States
const activeTab = ref('dashboard');
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

// Price Adjustment Modal States
const showPriceModal = ref(false);
const priceAdjustmentList = ref<PriceAdjustmentItem[]>([]);
const isSavingPrices = ref(false);

// Stock Out Modal States
const showStockOutModal = ref(false);
const stockOutItems = ref<StockOutItem[]>([]);
const isSubmittingStockOut = ref(false);

// Printer Settings Modal States
const showPrinterSettingsModal = ref(false);
const printerSettings = reactive({
  paperWidth: '80mm',
  autoPrint: true,
  showLogo: true,
  footerMessage: 'ขอบคุณที่ร่วมรักษ์โลกกับเรา\nRecycle Waste Management'
});

// Receipt Detail Modal States
const receiptModalRef = ref<InstanceType<typeof ReceiptDetailModal> | null>(null);
const selectedReceiptData = ref<ReceiptDetailData | null>(null);
const loadingReceiptDetail = ref(false);

// Expense Tracking Modal States
const showExpenseModal = ref(false);
const isSavingExpense = ref(false);
interface ExpenseItem {
  id: string;
  date: string;
  category: string;
  description: string;
  amount: number;
}
const expenseList = ref<ExpenseItem[]>([
  { id: '1', date: new Date().toISOString(), category: 'labor', description: 'ค่าแรงพนักงานรายวัน', amount: 500 },
  { id: '2', date: new Date().toISOString(), category: 'utility', description: 'ค่าน้ำมันรถรับของ', amount: 1200 }
]);
const newExpense = reactive({
  category: 'labor',
  description: '',
  amount: 0
});

// Loading states
const loadingStocks = ref(false);
const loadingReceipts = ref(false);

// Cache flags to prevent re-fetching
const stocksFetched = ref(false);
const receiptsFetched = ref(false);

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
  market_price?: number; // Optional market price
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

interface PriceAdjustmentItem extends StockItem {
  new_price: number;
}

interface StockOutItem extends StockItem {
  sell_quantity: number;
  sell_price: number;
  selected: boolean;
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

// Dashboard Data
const dashboardStats = reactive({
  todayPurchase: 15420,
  todayWeight: 450.5,
  estimatedProfit: 3200,
  customersCount: 12,
  totalExpenses: 8500,
  currentMonthWeight: 8500, // kg
  monthlyGoal: 10000 // kg
});

// Goal Setting States
const showGoalModal = ref(false);
const tempGoal = ref(10000);

const openGoalModal = () => {
  tempGoal.value = dashboardStats.monthlyGoal;
  showGoalModal.value = true;
};

const saveGoal = () => {
  if (tempGoal.value <= 0) {
    showModal('error', 'กรุณาระบุเป้าหมายที่ถูกต้อง');
    return;
  }
  dashboardStats.monthlyGoal = tempGoal.value;
  showGoalModal.value = false;
  showModal('success', 'บันทึกเป้าหมายเรียบร้อยแล้ว');
};

const revenueChartData = ref({
  labels: ['จ.', 'อ.', 'พ.', 'พฤ.', 'ศ.', 'ส.', 'อา.'],
  datasets: [{
    label: 'ยอดรับซื้อ (บาท)',
    backgroundColor: (ctx: { chart: { ctx: CanvasRenderingContext2D } }) => {
      const canvas = ctx.chart.ctx;
      const gradient = canvas.createLinearGradient(0, 0, 0, 400);
      gradient.addColorStop(0, 'rgba(34, 197, 94, 0.5)');
      gradient.addColorStop(1, 'rgba(34, 197, 94, 0.0)');
      return gradient;
    },
    borderColor: '#22c55e',
    data: [12000, 19000, 3000, 5000, 2000, 30000, 45000],
    fill: true,
    tension: 0.4
  }]
});

const categoryChartData = ref({
  labels: ['ขวดพลาสติก', 'กระดาษ', 'โลหะ/เหล็ก', 'ขวดแก้ว', 'อิเล็กทรอนิกส์'],
  datasets: [{
    backgroundColor: ['#22c55e', '#3b82f6', '#f59e0b', '#ef4444', '#8b5cf6'],
    data: [40, 25, 15, 15, 5]
  }]
});

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom' as const
    }
  }
};

const revenueChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      grid: {
        display: true,
        borderDash: [5, 5]
      }
    },
    x: {
      grid: {
        display: false
      }
    }
  }
};

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
  shop_code: shopStore.shop?.shop_code || '',
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

import { checkShopCode } from '@/services/shop';

// Shop Code Validation State
const shopCodeStatus = reactive({
  checking: false,
  available: true,
  message: ''
});

let debounceTimer: ReturnType<typeof setTimeout> | null = null;

// Watch for shop_code changes to validate
watch(() => editShopData.shop_code, (newCode) => {
  // Reset status
  shopCodeStatus.message = '';
  shopCodeStatus.available = true;

  // Clear previous timer
  if (debounceTimer) {
    clearTimeout(debounceTimer);
  }

  // If empty or same as current shop code, no need to check
  if (!newCode || newCode === shopStore.shop?.shop_code) {
    shopCodeStatus.checking = false;
    return;
  }

  // Validate format (only alphanumeric, -, _)
  const validFormat = /^[a-zA-Z0-9_-]+$/.test(newCode);
  if (!validFormat) {
    shopCodeStatus.available = false;
    shopCodeStatus.message = 'ใช้ได้เฉพาะภาษาอังกฤษ ตัวเลข ขีดกลาง (-) และขีดล่าง (_) เท่านั้น';
    shopCodeStatus.checking = false;
    return;
  }

  // Set checking state
  shopCodeStatus.checking = true;

  // Set new timer (2 seconds debounce)
  debounceTimer = setTimeout(async () => {
    try {
      const result = await checkShopCode(newCode);
      shopCodeStatus.available = result.available;
      if (!result.available) {
        shopCodeStatus.message = 'รหัสร้านนี้ถูกใช้งานแล้ว';
      } else {
        shopCodeStatus.message = 'สามารถใช้รหัสนี้ได้';
      }
    } catch (error) {
      console.error('Check shop code error:', error);
      shopCodeStatus.available = false; // Assume unavailable on error to be safe
      shopCodeStatus.message = 'เกิดข้อผิดพลาดในการตรวจสอบ';
    } finally {
      shopCodeStatus.checking = false;
    }
  }, 2000);
});

// Watch for changes in shop store to update edit form
watch(() => shopStore.shop, (newShop) => {
  if (newShop) {
    editShopData.shop_code = newShop.shop_code || '';
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
const fetchStocks = async (forceFetch = false) => {
  if (!shopStore.shop?.shop_id) return;

  // Skip if already fetched and not forcing refresh
  if (stocksFetched.value && !forceFetch) return;

  loadingStocks.value = true;
  try {
    const { fetchStocksByShopID } = await import('@/services/manageShop');
    const response = await fetchStocksByShopID(shopStore.shop.shop_id, stockPage.value, stockPageSize);
    console.log('fetchStocks response:', response);
    if (response.success && response.data) {
      stocks.value = response.data;
      stocksFetched.value = true;
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
const fetchReceipts = async (forceFetch = false) => {
  if (!shopStore.shop?.shop_id) return;

  // Skip if already fetched and not forcing refresh
  if (receiptsFetched.value && !forceFetch) return;

  loadingReceipts.value = true;
  try {
    const { fetchReceiptsByShopID } = await import('@/services/manageShop');
    const response = await fetchReceiptsByShopID(shopStore.shop.shop_id, receiptPage.value, receiptPageSize);
    console.log('fetchReceipts response:', response);
    if (response.success && response.data) {
      receipts.value = response.data;
      receiptsFetched.value = true;
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
  // Update URL hash
  router.replace({ hash: `#${newTab}` });

  if (newTab === 'stock') {
    fetchStocks();
  } else if (newTab === 'receipts') {
    fetchReceipts();
  } else if (newTab === 'dashboard') {
    // Optionally fetch dashboard data in future
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

const openPriceAdjustmentModal = () => {
  // Clone stocks to adjustment list
  // In a real app, maybe fetch fresh data or ensure we have all items
  priceAdjustmentList.value = stocks.value.map(item => ({
    ...item,
    // Start with current price
    new_price: item.current_price,
    // Mock market price if not exists (randomly slightly different from current price)
    market_price: item.market_price || item.current_price * (1 + (Math.random() * 0.1 - 0.05))
  }));
  showPriceModal.value = true;
};

const savePriceAdjustments = async () => {
  isSavingPrices.value = true;
  try {
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 1000));

    // In real implementation, call update API here
    // const { updateBatchPrices } = await import('@/services/manageShop');
    // await updateBatchPrices(shopStore.shop.shop_id, priceAdjustmentList.value.map(item => ({ id: item.id, price: item.new_price })));

    // Update local state to reflect changes
    stocks.value = stocks.value.map(item => {
      const adjustment = priceAdjustmentList.value.find(adj => adj.id === item.id);
      if (adjustment) {
        return {
          ...item,
          current_price: adjustment.new_price
        };
      }
      return item;
    });

    showPriceModal.value = false;
    showModal('success', 'ปรับราคาเรียบร้อยแล้ว');
  } catch (error) {
    console.error('Error saving prices:', error);
    showModal('error', 'เกิดข้อผิดพลาดในการบันทึกราคา');
  } finally {
    isSavingPrices.value = false;
  }
};

const openStockOutModal = () => {
  stockOutItems.value = stocks.value.map(item => ({
    ...item,
    sell_quantity: 0,
    sell_price: item.market_price ? (item.market_price * 1.1) : (item.current_price * 1.2), // Default selling price target (e.g. +20% margin)
    selected: false
  }));
  showStockOutModal.value = true;
};

const totalStockOutAmount = computed(() => {
  return stockOutItems.value
    .filter(item => item.selected)
    .reduce((sum, item) => sum + (item.sell_quantity * item.sell_price), 0);
});

const submitStockOut = async () => {
  const selectedItems = stockOutItems.value.filter(item => item.selected);

  if (selectedItems.length === 0) {
    showModal('error', 'กรุณาเลือกรายการที่ต้องการขาย');
    return;
  }

  // Validate quantities
  const invalidItems = selectedItems.filter(item => item.sell_quantity <= 0 || item.sell_quantity > item.quantity);
  if (invalidItems.length > 0) {
    showModal('error', `จำนวนสินค้าไม่ถูกต้องสำหรับ: ${invalidItems.map(i => i.name).join(', ')}`);
    return;
  }

  isSubmittingStockOut.value = true;
  try {
    // Simulate API
    await new Promise(resolve => setTimeout(resolve, 1500));

    let totalRealizedProfit = 0;

    // Update local stocks
    stocks.value = stocks.value.map(item => {
      const soldItem = selectedItems.find(s => s.id === item.id);
      if (soldItem) {
        // Calculate realized profit for this item
        const itemProfit = (soldItem.sell_price - item.purchase_price) * soldItem.sell_quantity;
        totalRealizedProfit += itemProfit;

        return {
          ...item,
          quantity: item.quantity - soldItem.sell_quantity,
          // Update profit tracking logic here if needed
        };
      }
      return item;
    });

    // Update Dashboard Stats (Real-time feedback)
    dashboardStats.estimatedProfit += totalRealizedProfit;

    showStockOutModal.value = false;
    showModal('success', `บันทึกรายการขายเรียบร้อย ยอดรวม ${formatCurrency(totalStockOutAmount.value)} (กำไร ${formatCurrency(totalRealizedProfit)})`);
    // Refresh dashboard stats if needed
  } catch (error) {
    console.error('Error stock out:', error);
    showModal('error', 'เกิดข้อผิดพลาดในการบันทึกรายการขาย');
  } finally {
    isSubmittingStockOut.value = false;
  }
};

const openPrinterSettings = () => {
  // Load settings from localStorage if available
  const saved = localStorage.getItem('printerSettings');
  if (saved) {
    Object.assign(printerSettings, JSON.parse(saved));
  }
  showPrinterSettingsModal.value = true;
};

const savePrinterSettings = () => {
  localStorage.setItem('printerSettings', JSON.stringify(printerSettings));
  showPrinterSettingsModal.value = false;
  showModal('success', 'บันทึกการตั้งค่าการพิมพ์เรียบร้อยแล้ว');
};

const openExpenseModal = () => {
  // Reset form
  newExpense.category = 'labor';
  newExpense.description = '';
  newExpense.amount = 0;
  showExpenseModal.value = true;
};

const saveExpense = async () => {
  if (newExpense.amount <= 0) {
    showModal('error', 'กรุณาระบุจำนวนเงินที่ถูกต้อง');
    return;
  }

  isSavingExpense.value = true;
  try {
    // Simulate API
    await new Promise(resolve => setTimeout(resolve, 800));

    const expense: ExpenseItem = {
      id: Date.now().toString(),
      date: new Date().toISOString(),
      category: newExpense.category,
      description: newExpense.description,
      amount: newExpense.amount
    };

    expenseList.value.unshift(expense);
    dashboardStats.totalExpenses += expense.amount;

    // Reset form for next entry
    newExpense.category = 'labor';
    newExpense.description = '';
    newExpense.amount = 0;

    showModal('success', `บันทึกรายจ่ายเรียบร้อย ${formatCurrency(expense.amount)}`);
  } catch (error) {
    console.error('Error saving expense:', error);
    showModal('error', 'บันทึกรายจ่ายไม่สำเร็จ');
  } finally {
    isSavingExpense.value = false;
  }
};

// Refresh data function
const refreshData = () => {
  if (activeTab.value === 'stock') {
    stocksFetched.value = false;
    fetchStocks(true);
  } else if (activeTab.value === 'receipts') {
    receiptsFetched.value = false;
    fetchReceipts(true);
  }
};

// Watch pagination changes
watch(stockPage, () => {
  if (activeTab.value === 'stock') {
    stocksFetched.value = false; // Reset cache for new page
    fetchStocks(true);
  }
});

watch(receiptPage, () => {
  if (activeTab.value === 'receipts') {
    receiptsFetched.value = false; // Reset cache for new page
    fetchReceipts(true);
  }
});

// Handle receipt row click
const handleReceiptClick = async (receiptId: string) => {
  loadingReceiptDetail.value = true;
  selectedReceiptData.value = null;

  // Open modal first
  receiptModalRef.value?.openModal();

  try {
    const { fetchReceiptByID } = await import('@/services/manageShop');
    const response = await fetchReceiptByID(receiptId);

    if (response.success && response.data) {
      selectedReceiptData.value = response.data;
    } else {
      console.error('Failed to fetch receipt details:', response.message);
      selectedReceiptData.value = null;
    }
  } catch (error) {
    console.error('Error fetching receipt details:', error);
    selectedReceiptData.value = null;
  } finally {
    loadingReceiptDetail.value = false;
  }
};

// Load shop data on component mount
onMounted(async () => {
  // Set active tab from hash if present
  const hash = route.hash.replace('#', '')
  if (hash && ['dashboard', 'info', 'stock', 'receipts', 'employees'].includes(hash)) {
    activeTab.value = hash;
  }

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

  // Load initial data based on active tab
  if (activeTab.value === 'stock') {
    fetchStocks();
  } else if (activeTab.value === 'receipts') {
    fetchReceipts();
  }
});

onUnmounted(() => {
  // Clean up map instance on component unmount
  cleanupMap();
});

// Expose refreshData for external use if needed
defineExpose({
  refreshData
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
          <button role="tab" :class="['tab', { 'tab-active': activeTab === 'dashboard' }]"
            @click="activeTab = 'dashboard'">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 mr-2" viewBox="0 0 24 24" stroke-width="2"
              stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
              <path stroke="none" d="M0 0h24v24H0z" fill="none" />
              <path d="M4 4h6v8h-6z" />
              <path d="M4 16h6v4h-6z" />
              <path d="M14 12h6v8h-6z" />
              <path d="M14 4h6v4h-6z" />
            </svg>
            ภาพรวม
          </button>
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
          <button role="tab" :class="['tab', { 'tab-active': activeTab === 'employees' }]"
            @click="activeTab = 'employees'">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 mr-2" viewBox="0 0 24 24" stroke-width="2"
              stroke="currentColor" fill="none">
              <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
              <path d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0 -8 0"></path>
              <path d="M6 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2"></path>
            </svg>
            พนักงาน
          </button>
        </div>

        <div class="p-6 md:p-8">
          <!-- Tab 0: Dashboard -->
          <div v-show="activeTab === 'dashboard'">

            <!-- Monthly Goal Section -->
            <div class="mb-8 bg-white p-6 rounded-2xl shadow-sm border border-gray-100 relative overflow-hidden">
              <div class="absolute top-0 right-0 p-4 opacity-10">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-32 h-32" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2l2.4 7.2h7.6l-6 4.8 2.4 7.2-6-4.8-6 4.8 2.4-7.2-6-4.8h7.6z" />
                </svg>
              </div>

              <div class="flex justify-between items-end mb-2 relative z-10">
                <div>
                  <h3 class="text-lg font-bold text-gray-800 flex items-center gap-2">
                    🎯 เป้าหมายลดโลกร้อนประจำเดือน
                    <span class="badge badge-primary badge-outline text-xs">Monthly Goal</span>
                  </h3>
                  <p class="text-gray-500 text-sm mt-1">ช่วยกันรับซื้อขยะรีไซเคิลให้ถึงเป้าหมาย!</p>
                </div>
                <div class="text-right">
                  <div class="text-3xl font-extrabold text-primary">
                    {{ Math.round((dashboardStats.currentMonthWeight / dashboardStats.monthlyGoal) * 100) }}%
                  </div>
                  <div class="text-xs text-gray-500">
                    {{ dashboardStats.currentMonthWeight.toLocaleString() }} / {{
                      dashboardStats.monthlyGoal.toLocaleString() }} kg
                  </div>
                </div>
              </div>

              <div class="w-full bg-gray-100 rounded-full h-6 relative z-10">
                <div
                  class="bg-gradient-to-r from-primary to-secondary h-6 rounded-full transition-all duration-1000 ease-out flex items-center justify-end pr-2 text-white text-xs font-bold shadow-lg"
                  :style="{ width: `${Math.min((dashboardStats.currentMonthWeight / dashboardStats.monthlyGoal) * 100, 100)}%` }">
                  <span v-if="(dashboardStats.currentMonthWeight / dashboardStats.monthlyGoal) * 100 >= 10">
                    {{ Math.round((dashboardStats.currentMonthWeight / dashboardStats.monthlyGoal) * 100) }}%
                  </span>
                </div>
              </div>

              <button class="btn btn-xs btn-ghost absolute top-2 right-2 z-20 text-gray-400 hover:text-gray-600"
                @click="openGoalModal">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                </svg>
                แก้ไขเป้าหมาย
              </button>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
              <!-- Total Purchase Today -->
              <div class="stats shadow bg-gradient-to-br from-blue-500 to-blue-600 text-white">
                <div class="stat">
                  <div class="stat-figure text-blue-200">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24"
                      stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  <div class="stat-title text-blue-100">ยอดรับซื้อวันนี้</div>
                  <div class="stat-value text-3xl">{{ formatCurrency(dashboardStats.todayPurchase) }}</div>
                  <div class="stat-desc text-blue-100 flex items-center gap-1">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M12 7a1 1 0 110-2 1 1 0 010 2zm-2 7a1 1 0 100-2 1 1 0 000 2z"
                        clip-rule="evenodd" />
                    </svg>
                    น้ำหนักรวม {{ dashboardStats.todayWeight }} กก.
                  </div>
                </div>
              </div>

              <!-- Estimated Profit -->
              <div class="stats shadow bg-gradient-to-br from-green-500 to-emerald-600 text-white">
                <div class="stat">
                  <div class="stat-figure text-green-200">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24"
                      stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                    </svg>
                  </div>
                  <div class="stat-title text-green-100">กำไรโดยประมาณ</div>
                  <div class="stat-value text-3xl">{{ formatCurrency(dashboardStats.estimatedProfit) }}</div>
                  <div class="stat-desc text-green-100">
                    ↗︎ 12% จากเมื่อวาน
                  </div>
                </div>
              </div>

              <!-- Customers Today -->
              <div class="stats shadow bg-gradient-to-br from-purple-500 to-indigo-600 text-white">
                <div class="stat">
                  <div class="stat-figure text-purple-200">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24"
                      stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                  </div>
                  <div class="stat-title text-purple-100">ลูกค้าวันนี้</div>
                  <div class="stat-value text-3xl">{{ dashboardStats.customersCount }} คน</div>
                </div>
              </div>

              <!-- Expense Card -->
              <div class="stats shadow bg-gradient-to-br from-red-500 to-pink-600 text-white">
                <div class="stat">
                  <div class="stat-figure text-red-200">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24"
                      stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M15 12H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  <div class="stat-title text-red-100">รายจ่ายดำเนินงาน</div>
                  <div class="stat-value text-3xl">{{ formatCurrency(dashboardStats.totalExpenses) }}</div>
                  <div class="stat-actions mt-2">
                    <button class="btn btn-xs btn-outline text-white border-white hover:bg-white hover:text-red-500"
                      @click="openExpenseModal">
                      + บันทึกรายจ่าย
                    </button>
                  </div>
                </div>
              </div>

              <!-- Customers Card -->
              <div class="stats shadow bg-gradient-to-br from-purple-500 to-indigo-600 text-white">
                <div class="stat">
                  <div class="stat-figure text-purple-200">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24"
                      stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                  </div>
                  <div class="stat-title text-purple-100">ลูกค้าวันนี้</div>
                  <div class="stat-value text-3xl">{{ dashboardStats.customersCount }}</div>
                  <div class="stat-desc text-purple-100">
                    เพิ่มขึ้น 2 คนจากเมื่อวาน
                  </div>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <!-- Revenue Chart -->
              <div class="lg:col-span-2 bg-white rounded-xl shadow-lg p-6 border border-gray-100">
                <div class="flex justify-between items-center mb-6">
                  <h3 class="font-bold text-lg text-gray-700 flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-green-600" fill="none"
                      viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
                    </svg>
                    แนวโน้มการรับซื้อ 7 วันย้อนหลัง
                  </h3>
                  <select class="select select-sm select-bordered">
                    <option>7 วันล่าสุด</option>
                    <option>30 วันล่าสุด</option>
                  </select>
                </div>
                <div class="h-80 w-full">
                  <Line :data="revenueChartData" :options="revenueChartOptions" />
                </div>
              </div>

              <!-- Top Categories -->
              <div class="bg-white rounded-xl shadow-lg p-6 border border-gray-100">
                <h3 class="font-bold text-lg text-gray-700 mb-6 flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-600" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z" />
                  </svg>
                  สัดส่วนขยะที่รับซื้อ
                </h3>
                <div class="h-64 flex justify-center items-center">
                  <Doughnut :data="categoryChartData" :options="chartOptions" />
                </div>
                <!-- Mini Stats List -->
                <div class="mt-4 space-y-2">
                  <div v-for="(label, idx) in categoryChartData.labels.slice(0, 3)" :key="idx"
                    class="flex justify-between items-center text-sm">
                    <div class="flex items-center gap-2">
                      <div class="w-3 h-3 rounded-full"
                        :style="{ backgroundColor: categoryChartData.datasets[0].backgroundColor[idx] }"></div>
                      <span>{{ label }}</span>
                    </div>
                    <span class="font-semibold">{{ categoryChartData.datasets[0].data[idx] }}%</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

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
                  <h3 class="font-semibold text-gray-600">รหัสร้าน (Shop Code)</h3>
                  <p class="text-lg font-mono bg-gray-100 px-3 py-2 rounded">{{ shopStore.shop?.shop_code || '-' }}</p>
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
                  <div class="form-control">
                    <label class="label">
                      <span class="label-text">รหัสร้าน (Shop Code) *</span>
                      <span class="label-text-alt text-gray-500">{{ editShopData.shop_code?.length || 0 }}/12</span>
                    </label>
                    <input v-model="editShopData.shop_code" type="text"
                      :class="['input input-bordered', { 'input-error': !shopCodeStatus.available, 'input-success': shopCodeStatus.available && editShopData.shop_code && !shopCodeStatus.checking && editShopData.shop_code !== shopStore.shop?.shop_code }]"
                      required maxlength="12" placeholder="เช่น SHOP001" />
                    <label class="label" v-if="shopCodeStatus.checking || shopCodeStatus.message">
                      <span class="label-text-alt flex items-center gap-1"
                        :class="{ 'text-error': !shopCodeStatus.available, 'text-success': shopCodeStatus.available }">
                        <span v-if="shopCodeStatus.checking" class="loading loading-spinner loading-xs"></span>
                        <svg v-else-if="shopCodeStatus.available" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4"
                          viewBox="0 0 20 20" fill="currentColor">
                          <path fill-rule="evenodd"
                            d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                            clip-rule="evenodd" />
                        </svg>
                        <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20"
                          fill="currentColor">
                          <path fill-rule="evenodd"
                            d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                            clip-rule="evenodd" />
                        </svg>
                        {{ shopCodeStatus.checking ? 'กำลังตรวจสอบ...' : shopCodeStatus.message }}
                      </span>
                    </label>
                  </div>
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
                <button class="btn btn-ghost btn-sm" @click="refreshData" :disabled="loadingStocks"
                  title="รีเฟรชข้อมูล">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="w-4 h-4" :class="{ 'animate-spin': loadingStocks }">
                    <path stroke-linecap="round" stroke-linejoin="round"
                      d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                  </svg>
                </button>
                <button class="btn btn-primary btn-sm text-white" @click="goToWastePurchase">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="w-4 h-4 mr-1">
                    <path stroke-linecap="round" stroke-linejoin="round"
                      d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  ชั่งชำระเงิน
                </button>
                <button class="btn btn-warning btn-sm text-white" @click="openPriceAdjustmentModal">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 mr-1" viewBox="0 0 24 24" stroke-width="2"
                    stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                    <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                    <path d="M6 9l6 6l6 -6"></path>
                  </svg>
                  ปรับราคา
                </button>
                <button class="btn btn-error btn-sm text-white" @click="openStockOutModal">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="w-4 h-4 mr-1">
                    <path stroke-linecap="round" stroke-linejoin="round"
                      d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
                  </svg>
                  ส่งโรงงาน (ขาย)
                </button>
              </div>
            </div>

            <!-- Mobile Action Buttons (Visible only on mobile) -->
            <div class="grid grid-cols-1 gap-3 mb-6 md:hidden">
              <button class="btn btn-ghost w-full" @click="refreshData" :disabled="loadingStocks">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                  stroke="currentColor" class="w-5 h-5 mr-2" :class="{ 'animate-spin': loadingStocks }">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                </svg>
                รีเฟรชข้อมูล
              </button>
              <button class="btn btn-primary w-full text-white" @click="goToWastePurchase">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                  stroke="currentColor" class="w-5 h-5 mr-2">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                ชั่งชำระเงิน
              </button>
              <button class="btn btn-warning w-full text-white" @click="openPriceAdjustmentModal">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 mr-2" viewBox="0 0 24 24" stroke-width="2"
                  stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                  <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                  <path d="M6 9l6 6l6 -6"></path>
                </svg>
                ปรับราคาประจำวัน
              </button>
              <button class="btn btn-error w-full text-white" @click="openStockOutModal">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                  stroke="currentColor" class="w-5 h-5 mr-2">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
                </svg>
                ส่งโรงงาน (ขาย)
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
            <div class="flex justify-between items-center mb-6">
              <h2 class="text-xl font-bold">ประวัติการทำรายการ</h2>
              <div class="flex gap-2">
                <button class="btn btn-outline btn-sm" @click="openPrinterSettings">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  ตั้งค่าใบเสร็จ
                </button>
                <button class="btn btn-ghost btn-sm" @click="refreshData" :disabled="loadingReceipts"
                  title="รีเฟรชข้อมูล">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="w-4 h-4" :class="{ 'animate-spin': loadingReceipts }">
                    <path stroke-linecap="round" stroke-linejoin="round"
                      d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                  </svg>
                </button>
              </div>
            </div>
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
                  <tr v-for="receipt in paginatedReceipts" :key="receipt.id"
                    class="cursor-pointer hover:bg-base-200 transition-colors" @click="handleReceiptClick(receipt.id)">
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

          <!-- Tab 4: Employees -->
          <div v-show="activeTab === 'employees'">
            <EmployeeManagement />
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

    <!-- Price Adjustment Modal -->
    <dialog id="price-adjustment-modal" class="modal modal-bottom sm:modal-middle"
      :class="{ 'modal-open': showPriceModal }">
      <div class="modal-box w-11/12 max-w-4xl">
        <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-warning" viewBox="0 0 24 24" stroke-width="2"
            stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M3 21l18 0"></path>
            <path d="M3 10l18 0"></path>
            <path d="M5 6l7 -3l7 3"></path>
            <path d="M4 10l0 11"></path>
            <path d="M20 10l0 11"></path>
            <path d="M8 14l0 3"></path>
            <path d="M12 14l0 3"></path>
            <path d="M16 14l0 3"></path>
          </svg>
          ปรับราคารับซื้อประจำวัน
        </h3>

        <div class="overflow-x-auto max-h-[60vh]">
          <table class="table table-pin-rows">
            <thead>
              <tr class="bg-base-200">
                <th>หมวดหมู่</th>
                <th>รายการ</th>
                <th class="text-right">ราคากลาง (บาท)</th>
                <th class="text-right">ราคาปัจจุบัน (บาท)</th>
                <th class="text-right w-32">ราคาใหม่ (บาท)</th>
                <th class="text-right">เปลี่ยนแปลง</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in priceAdjustmentList" :key="item.id">
                <td><span class="badge badge-ghost">{{ item.category }}</span></td>
                <td class="font-semibold">{{ item.name }}</td>
                <td class="text-right text-gray-500">{{ formatCurrency(item.market_price || 0) }}</td>
                <td class="text-right">{{ formatCurrency(item.current_price) }}</td>
                <td>
                  <input type="number" v-model.number="item.new_price" min="0" step="0.01"
                    class="input input-bordered input-sm w-full text-right" />
                </td>
                <td class="text-right">
                  <span v-if="item.new_price > item.current_price" class="text-success font-bold">
                    +{{ (item.new_price - item.current_price).toFixed(2) }}
                  </span>
                  <span v-else-if="item.new_price < item.current_price" class="text-error font-bold">
                    {{ (item.new_price - item.current_price).toFixed(2) }}
                  </span>
                  <span v-else class="text-gray-400">-</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="modal-action">
          <button type="button" class="btn" @click="showPriceModal = false" :disabled="isSavingPrices">ยกเลิก</button>
          <button type="button" class="btn btn-primary" @click="savePriceAdjustments" :disabled="isSavingPrices">
            <span v-if="isSavingPrices" class="loading loading-spinner"></span>
            บันทึกการเปลี่ยนแปลง
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="showPriceModal = false">close</button>
      </form>
    </dialog>

    <!-- Stock Out (Sales Order) Modal -->
    <dialog id="stock-out-modal" class="modal modal-bottom sm:modal-middle"
      :class="{ 'modal-open': showStockOutModal }">
      <div class="modal-box w-11/12 max-w-4xl">
        <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-error" viewBox="0 0 24 24" stroke-width="2"
            stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M14 3v4a1 1 0 0 0 1 1h4"></path>
            <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"></path>
            <path d="M12 11l0 6"></path>
            <path d="M9 14l3 3l3 -3"></path>
          </svg>
          ส่งโรงงาน / ขายออก (Sales Order)
        </h3>

        <div class="alert alert-info shadow-sm mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
            class="stroke-current shrink-0 w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <div>
            <h3 class="font-bold">สร้างรายการขาย</h3>
            <div class="text-xs">เลือกรายการที่ต้องการส่งขายโรงงานและระบุจำนวน/ราคาที่ตกลงขาย</div>
          </div>
        </div>

        <div class="overflow-x-auto max-h-[50vh]">
          <table class="table table-pin-rows">
            <thead>
              <tr class="bg-base-200">
                <th class="w-10">เลือก</th>
                <th>รายการ</th>
                <th class="text-right">คงเหลือ</th>
                <th class="text-right w-32">จำนวนที่ขาย</th>
                <th class="text-right w-32">ราคาขาย/กก.</th>
                <th class="text-right">รวมเงิน</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in stockOutItems" :key="item.id"
                :class="{ 'bg-base-100': !item.selected, 'bg-green-50': item.selected }">
                <td>
                  <input type="checkbox" v-model="item.selected" class="checkbox checkbox-primary" />
                </td>
                <td class="font-semibold" :class="{ 'opacity-50': !item.selected }">{{ item.name }}</td>
                <td class="text-right">{{ item.quantity }} กก.</td>
                <td>
                  <input type="number" v-model.number="item.sell_quantity" :disabled="!item.selected" min="0"
                    :max="item.quantity" class="input input-bordered input-sm w-full text-right" placeholder="0" />
                </td>
                <td>
                  <input type="number" v-model.number="item.sell_price" :disabled="!item.selected" min="0" step="0.01"
                    class="input input-bordered input-sm w-full text-right" placeholder="0.00" />
                </td>
                <td class="text-right font-bold text-green-700">
                  {{ item.selected ? formatCurrency(item.sell_quantity * item.sell_price) : '-' }}
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="bg-base-200 font-bold text-lg">
                <td colspan="5" class="text-right">ยอดรวมสุทธิ:</td>
                <td class="text-right text-success">{{ formatCurrency(totalStockOutAmount) }}</td>
              </tr>
            </tfoot>
          </table>
        </div>

        <div class="modal-action">
          <button type="button" class="btn" @click="showStockOutModal = false"
            :disabled="isSubmittingStockOut">ยกเลิก</button>
          <button type="button" class="btn btn-success text-white" @click="submitStockOut"
            :disabled="isSubmittingStockOut || totalStockOutAmount <= 0">
            <span v-if="isSubmittingStockOut" class="loading loading-spinner"></span>
            ยืนยันขายออก
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="showStockOutModal = false">close</button>
      </form>
    </dialog>

    <!-- Printer Settings Modal -->
    <dialog id="printer-settings-modal" class="modal modal-bottom sm:modal-middle"
      :class="{ 'modal-open': showPrinterSettingsModal }">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-gray-700" viewBox="0 0 24 24" stroke-width="2"
            stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M17 17h2a2 2 0 0 0 2 -2v-4a2 2 0 0 0 -2 -2h-14a2 2 0 0 0 -2 2v4a2 2 0 0 0 2 2h2"></path>
            <path d="M17 9v-4a2 2 0 0 0 -2 -2h-6a2 2 0 0 0 -2 2v4"></path>
            <path d="M7 13m0 2a2 2 0 0 1 2 -2h6a2 2 0 0 1 2 2v4a2 2 0 0 1 -2 2h-6a2 2 0 0 1 -2 -2z"></path>
          </svg>
          ตั้งค่าการพิมพ์ใบเสร็จ
        </h3>

        <div class="space-y-4">
          <div class="form-control">
            <label class="label cursor-pointer justify-start gap-4">
              <span class="label-text font-semibold">ขนาดกระดาษ</span>
              <div class="join">
                <input class="join-item btn btn-sm" type="radio" name="paperWidth" aria-label="58mm" value="58mm"
                  v-model="printerSettings.paperWidth" />
                <input class="join-item btn btn-sm" type="radio" name="paperWidth" aria-label="80mm" value="80mm"
                  v-model="printerSettings.paperWidth" />
              </div>
            </label>
          </div>

          <div class="form-control">
            <label class="label cursor-pointer justify-start gap-4">
              <span class="label-text font-semibold">พิมพ์อัตโนมัติเมื่อทำรายการสำเร็จ</span>
              <input type="checkbox" class="toggle toggle-success" v-model="printerSettings.autoPrint" />
            </label>
          </div>

          <div class="form-control">
            <label class="label cursor-pointer justify-start gap-4">
              <span class="label-text font-semibold">แสดงโลโก้ร้านบนใบเสร็จ</span>
              <input type="checkbox" class="toggle toggle-primary" v-model="printerSettings.showLogo" />
            </label>
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text font-semibold">ข้อความท้ายใบเสร็จ (Footer Message)</span>
            </label>
            <textarea class="textarea textarea-bordered h-24" placeholder="พิมพ์ข้อความที่ต้องการแสดงท้ายใบเสร็จ..."
              v-model="printerSettings.footerMessage"></textarea>
          </div>
        </div>

        <div class="modal-action">
          <button type="button" class="btn btn-outline"
            @click="showModal('success', 'ทดสอบพิมพ์ใบเสร็จเรียบร้อย (จำลอง)')">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            ทดสอบพิมพ์
          </button>
          <button class="btn btn-primary" @click="savePrinterSettings">บันทึกตั้งค่า</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="showPrinterSettingsModal = false">close</button>
      </form>
    </dialog>

    <!-- Expense Tracking Modal -->
    <dialog id="expense-modal" class="modal modal-bottom sm:modal-middle" :class="{ 'modal-open': showExpenseModal }">
      <div class="modal-box w-11/12 max-w-3xl">
        <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-red-500" viewBox="0 0 24 24" stroke-width="2"
            stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M14 3v4a1 1 0 0 0 1 1h4"></path>
            <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"></path>
            <path d="M10 14l4 0"></path>
          </svg>
          บันทึกรายจ่าย (Expense Tracking)
        </h3>

        <div class="grid md:grid-cols-2 gap-6">
          <!-- Form Section -->
          <div>
            <div class="card bg-base-100 border shadow-sm">
              <div class="card-body p-4">
                <h4 class="card-title text-sm mb-2">เพิ่มรายการใหม่</h4>
                <form @submit.prevent="saveExpense" class="space-y-3">
                  <div class="form-control">
                    <label class="label text-xs">หมวดหมู่</label>
                    <select class="select select-bordered select-sm w-full" v-model="newExpense.category">
                      <option value="labor">ค่าแรงพนักงาน</option>
                      <option value="utility">ค่าน้ำ/ค่าไฟ/ค่าน้ำมัน</option>
                      <option value="rent">ค่าเช่าสถานที่</option>
                      <option value="maintenance">ค่าซ่อมบำรุง</option>
                      <option value="other">อื่นๆ</option>
                    </select>
                  </div>

                  <div class="form-control">
                    <label class="label text-xs">รายละเอียด</label>
                    <input type="text" class="input input-bordered input-sm w-full" v-model="newExpense.description"
                      placeholder="เช่น ค่าแรงนาย A" required />
                  </div>

                  <div class="form-control">
                    <label class="label text-xs">จำนวนเงิน (บาท)</label>
                    <input type="number" class="input input-bordered input-sm w-full" v-model.number="newExpense.amount"
                      min="1" required />
                  </div>

                  <button type="submit" class="btn btn-error btn-sm text-white w-full mt-2" :disabled="isSavingExpense">
                    <span v-if="isSavingExpense" class="loading loading-spinner loading-xs"></span>
                    บันทึก
                  </button>
                </form>
              </div>
            </div>
          </div>

          <!-- History Section -->
          <div class="overflow-y-auto max-h-[400px]">
            <h4 class="font-bold mb-2">ประวัติรายจ่ายล่าสุด</h4>
            <div v-if="expenseList.length === 0" class="text-center text-gray-500 py-4 text-sm">
              ยังไม่มีรายการรายจ่าย
            </div>
            <ul v-else class="space-y-2">
              <li v-for="expense in expenseList" :key="expense.id"
                class="flex justify-between items-center p-3 bg-base-100 rounded-lg border hover:bg-gray-50 text-sm">
                <div>
                  <div class="font-semibold">{{ expense.description }}</div>
                  <div class="text-xs text-gray-500">
                    <span class="badge badge-xs badge-ghost mr-1">{{ expense.category }}</span>
                    {{ new Date(expense.date).toLocaleString('th-TH') }}
                  </div>
                </div>
                <div class="font-bold text-red-500">-{{ formatCurrency(expense.amount) }}</div>
              </li>
            </ul>
          </div>
        </div>

        <div class="modal-action">
          <form method="dialog">
            <button class="btn" @click="showExpenseModal = false">ปิดหน้าต่าง</button>
          </form>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="showExpenseModal = false">close</button>
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

    <!-- Receipt Detail Modal -->
    <ReceiptDetailModal ref="receiptModalRef" :receipt-data="selectedReceiptData" :is-loading="loadingReceiptDetail" />

    <!-- Goal Setting Modal -->
    <dialog id="goal-modal" class="modal modal-bottom sm:modal-middle" :class="{ 'modal-open': showGoalModal }">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
          🎯 ตั้งเป้าหมายประจำเดือน
        </h3>
        <div class="form-control">
          <label class="label">
            <span class="label-text">เป้าหมายน้ำหนักรวม (กิโลกรัม)</span>
          </label>
          <input type="number" v-model.number="tempGoal"
            class="input input-bordered w-full text-lg font-bold text-primary" placeholder="ระบุเป้าหมาย" />
          <label class="label">
            <span class="label-text-alt text-gray-500">เป้าหมายปัจจุบัน: {{ dashboardStats.monthlyGoal.toLocaleString()
            }} kg</span>
          </label>
        </div>
        <div class="modal-action">
          <button class="btn" @click="showGoalModal = false">ยกเลิก</button>
          <button class="btn btn-primary text-white" @click="saveGoal">บันทึกเป้าหมาย</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="showGoalModal = false">close</button>
      </form>
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
