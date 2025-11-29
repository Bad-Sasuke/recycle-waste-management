<script setup lang="ts">
import { ref, shallowRef, onMounted, computed, watch, onUnmounted, h, render, nextTick } from 'vue';
import { IconSearch, IconCurrentLocation, IconList, IconSend, IconStar, IconStarFilled, IconFileText } from '@tabler/icons-vue';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

// Fix Leaflet marker icons
import iconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
import iconUrl from 'leaflet/dist/images/marker-icon.png';
import shadowUrl from 'leaflet/dist/images/marker-shadow.png';
import { useShopStore } from '../stores/shop';
import { useCustomerRequestStore } from '../stores/customer_request';
import { useChatStore } from '../stores/chat';
import { useUsersStore } from '../stores/users';
import { storeToRefs } from 'pinia';
import AlertModal from '../components/AlertModal.vue';
import LiquidMarker from '../components/map/LiquidMarker.vue';
import ShopPopup from '../components/map/ShopPopup.vue';
import UserLocationPopup from '../components/map/UserLocationPopup.vue';
import config from '../config';
import type { Shop } from '../types/shop';
import type { CustomerRequest } from '../types/customer_request';
import jsPDF from 'jspdf';
import autoTable from 'jspdf-autotable';

// Fix for default marker icon
// eslint-disable-next-line @typescript-eslint/no-explicit-any
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl,
  iconUrl,
  shadowUrl,
});


const shopStore = useShopStore();
const customerRequestsMyRequeststore = useCustomerRequestStore();
const chatStore = useChatStore();
const usersStore = useUsersStore();
const { customerRequestsMyRequests, errorCustomerRequest } = storeToRefs(customerRequestsMyRequeststore);
const { messages: chatMessages } = storeToRefs(chatStore);
const filteredChatMessages = computed(() => chatMessages.value.filter(msg => msg.type === 'message'));

const chatShop = computed(() => {
  // Try to find shop from chat messages (sender != current user)
  const otherUserMsg = filteredChatMessages.value.find(msg => msg.sender_id !== usersStore.user?.user_id);
  if (otherUserMsg) {
    return shopStore.allShops.find(s => s.user_id === otherUserMsg.sender_id);
  }

  // Fallback: if selectedShopId is set
  if (selectedShopId.value) {
    return shopStore.allShops.find(s => s.shop_id === selectedShopId.value);
  }

  return null;
});

const activeRequestId = ref<string | null>(null);
const searchQuery = ref('');
const map = shallowRef<L.Map | null>(null);
const markers = shallowRef<L.Marker[]>([]);
const selectedShopId = ref<string | null>(null);
const userLocation = ref<{ lat: number; lng: number } | null>(null);
const shareLocationModal = ref<HTMLDialogElement | null>(null);
const wasteDetails = ref('');
const userMarker = shallowRef<L.CircleMarker | null>(null);
const requestMarker = shallowRef<L.Marker | null>(null);
const shareStatus = ref<'idle' | 'pending' | 'accepted'>('idle');
const chatModal = ref<HTMLDialogElement | null>(null);
const chatMessagesContainer = ref<HTMLDivElement | null>(null);
const webAPI = config.webAPI;

const currentMessage = ref('');
const reviewModal = ref<HTMLDialogElement | null>(null);
const rating = ref(5);
const reviewComment = ref('');
const alertModal = ref<InstanceType<typeof AlertModal> | null>(null);
const requestsModal = ref<HTMLDialogElement | null>(null);
const receiptModal = ref<HTMLDialogElement | null>(null);
const refreshIntervalId = ref<number | null>(null);
interface ReceiptItem {
  name: string;
  category: string;
  weight: number;
  unit_price: number;
  price: number;
}

interface Receipt {
  id: string;
  created_at: string;
  payment_method: string;
  total_amount: number;
  vat_rate: number;
  vat: number;
  net_total: number;
}

interface ReceiptData {
  receipt: Receipt;
  items: ReceiptItem[];
  shop: Shop;
}

const currentReviewRequest = ref<CustomerRequest | null>(null);
const receiptData = ref<ReceiptData | null>(null);
const isLoadingReceipt = ref(false);
const reviewedRequests = ref<Map<string, boolean>>(new Map()); // Track which requests have been reviewed



const checkReviewStatus = async (requestId: string) => {
  try {
    const response = await fetch(`${webAPI}/api/reviews/check/${requestId}`);
    const data = await response.json();
    reviewedRequests.value.set(requestId, data.exists);
  } catch (error) {
    console.error('Error checking review status:', error);
  }
};

const openRequestsModal = async () => {
  if (requestsModal.value) {
    requestsModal.value.showModal();

    // Check review status for all done requests
    const doneRequests = customerRequestsMyRequests.value.filter(req => req.status === 'done');
    for (const req of doneRequests) {
      await checkReviewStatus(req.customer_request_id);
    }
  }
};

const closeRequestsModal = () => {
  if (requestsModal.value) {
    requestsModal.value.close();
  }
};

const openReceiptModal = async (requestId: string) => {
  isLoadingReceipt.value = true;
  try {
    const response = await fetch(`${webAPI}/api/receipts/by-request/${requestId}`, {
      headers: {
        'Authorization': `Bearer ${usersStore.jwt}`
      }
    });

    if (!response.ok) throw new Error('Receipt not found');

    const result = await response.json();
    receiptData.value = result.data;

    if (receiptModal.value) {
      receiptModal.value.showModal();
    }
  } catch (error) {
    console.error('Error fetching receipt:', error);
    if (alertModal.value) {
      alertModal.value.show('ไม่พบใบเสร็จสำหรับรายการนี้', 'error');
    }
  } finally {
    isLoadingReceipt.value = false;
  }
};

const closeReceiptModal = () => {
  if (receiptModal.value) {
    receiptModal.value.close();
  }
  receiptData.value = null;
};

const downloadPDF = () => {
  if (!receiptData.value) return;

  const doc = new jsPDF();
  const { receipt, items, shop } = receiptData.value;

  // Set font
  doc.setFont('helvetica');

  // Title
  doc.setFontSize(20);
  doc.text('Receipt / ใบเสร็จรับเงิน', 105, 20, { align: 'center' });

  // Receipt Info
  doc.setFontSize(10);
  let y = 35;
  doc.text(`Receipt No: ${receipt.id}`, 15, y);
  y += 7;
  doc.text(`Date: ${new Date(receipt.created_at).toLocaleString('th-TH')}`, 15, y);
  y += 7;
  doc.text(`Payment: ${receipt.payment_method}`, 15, y);
  y += 10;

  // Shop Info
  if (shop) {
    doc.setFontSize(12);
    doc.text('Shop Information', 15, y);
    doc.setFontSize(10);
    y += 7;
    doc.text(`Name: ${shop.name}`, 15, y);
    y += 7;
    doc.text(`Address: ${shop.address}`, 15, y);
    if (shop.phone) {
      y += 7;
      doc.text(`Phone: ${shop.phone}`, 15, y);
    }
    y += 10;
  }

  // Items Table
  const tableData = items.map((item) => [
    item.name,
    item.category,
    item.weight.toFixed(2),
    item.unit_price.toFixed(2),
    item.price.toFixed(2)
  ]);

  autoTable(doc, {
    startY: y,
    head: [['Item', 'Category', 'Weight (kg)', 'Price/Unit', 'Total']],
    body: tableData,
    theme: 'grid',
    headStyles: { fillColor: [66, 139, 202] },
  });

  // Summary
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const finalY = (doc as any).lastAutoTable.finalY + 10;
  doc.setFontSize(10);
  doc.text(`Subtotal: ฿${receipt.total_amount.toFixed(2)}`, 140, finalY, { align: 'right' });
  doc.text(`VAT (${(receipt.vat_rate * 100).toFixed(0)}%): ฿${receipt.vat.toFixed(2)}`, 140, finalY + 7, { align: 'right' });
  doc.setFontSize(12);
  doc.setFont('helvetica', 'bold');
  doc.text(`Net Total: ฿${receipt.net_total.toFixed(2)}`, 140, finalY + 14, { align: 'right' });

  // Save
  doc.save(`receipt-${receipt.id}.pdf`);
};


const isLoading = computed(() => shopStore.isLoading);

const filteredShops = computed(() => {
  if (!searchQuery.value) return shopStore.allShops;
  const query = searchQuery.value.toLowerCase();
  return shopStore.allShops.filter(shop =>
    shop.name?.toLowerCase().includes(query) ||
    shop.address?.toLowerCase().includes(query)
  );
});

// Generate popup content based on share status
// Generate popup content based on share status
const getUserPopupContent = computed(() => {
  const div = document.createElement('div');
  render(h(UserLocationPopup, { status: shareStatus.value }), div);
  return div.innerHTML;
});

// Initialize map
const initMap = () => {
  // Default center (Bangkok)
  const defaultLat = 13.7563;
  const defaultLng = 100.5018;

  map.value = L.map('map').setView([defaultLat, defaultLng], 10);

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  }).addTo(map.value);
};

// Add markers to map
const updateMarkers = () => {
  if (!map.value) return;

  // Clear existing markers
  markers.value.forEach(marker => marker.remove());
  markers.value = [];

  // Add new markers
  filteredShops.value.forEach(shop => {
    if (shop.latitude && shop.longitude) {
      // Create custom icon with liquid glass style
      // Create custom icon with liquid glass style
      const iconDiv = document.createElement('div');
      render(h(LiquidMarker, { colorClass: 'bg-primary' }), iconDiv);

      const customIcon = L.divIcon({
        className: 'custom-shop-marker',
        html: iconDiv.innerHTML,
        iconSize: [40, 40],
        iconAnchor: [20, 40],
        popupAnchor: [0, -40]
      });

      const popupDiv = document.createElement('div');
      render(h(ShopPopup, { shop: shop }), popupDiv);

      const marker = L.marker([shop.latitude, shop.longitude], {
        icon: customIcon
      })
        .addTo(map.value!)
        .bindPopup(popupDiv.innerHTML, {
          maxWidth: 280,
          className: 'glass-popup'
        });

      marker.on('click', () => {
        selectedShopId.value = shop.shop_id || null;
      });

      markers.value.push(marker);
    }
  });
};

// Focus map on specific request
const focusOnRequest = (req: CustomerRequest) => {
  if (!map.value) return;

  closeRequestsModal();
  map.value.setView([req.latitude, req.longitude], 16);

  if (requestMarker.value) {
    requestMarker.value.remove();
  }

  requestMarker.value = L.marker([req.latitude, req.longitude])
    .addTo(map.value)
    .bindPopup(`
      <div class="text-center p-2">
        <p class="font-bold text-sm mb-1">📍 ตำแหน่งคำขอ</p>
        <p class="text-xs text-gray-600">${req.description}</p>
        <div class="badge badge-sm mt-2 ${req.status === 'pending' ? 'badge-warning' : req.status === 'accepted' ? 'badge-success' : 'badge-error'} text-white">${req.status}</div>
      </div>
    `)
    .openPopup();
};

// Focus map on specific shop
const focusOnShop = (shop: Shop) => {
  selectedShopId.value = shop.shop_id || null;
  if (map.value && shop.latitude && shop.longitude) {
    map.value.setView([shop.latitude, shop.longitude], 15);

    // Find and open popup
    const marker = markers.value.find(m => {
      const latLng = m.getLatLng();
      return Math.abs(latLng.lat - shop.latitude!) < 0.0001 &&
        Math.abs(latLng.lng - shop.longitude!) < 0.0001;
    });

    if (marker) {
      marker.openPopup();
    }
  }
};

// Locate user
const locateUser = () => {
  if (!map.value) return;

  map.value.locate({ setView: true, maxZoom: 16 });

  map.value.on('locationfound', (e: L.LocationEvent) => {
    userLocation.value = e.latlng;

    // Remove old marker if exists
    if (userMarker.value) {
      userMarker.value.remove();
    }

    // Add user marker with glassmorphism popup
    userMarker.value = L.circleMarker(e.latlng, {
      radius: 8,
      fillColor: 'var(--color-info)',
      color: '#fff',
      weight: 2,
      opacity: 1,
      fillOpacity: 0.8
    }).addTo(map.value!);

    // Use computed content
    userMarker.value.bindPopup(getUserPopupContent.value, {
      maxWidth: 220,
      className: 'glass-popup'
    }).openPopup();
  });

  map.value.on('locationerror', (e: L.ErrorEvent) => {
    console.error('Location error:', e.message);
    // Fallback to default view if location access denied
  });
};

// Watch for status changes and update popup
watch(shareStatus, () => {
  if (userMarker.value) {
    userMarker.value.setPopupContent(getUserPopupContent.value);
  }
});

const scrollToBottom = async () => {
  await nextTick();
  if (chatMessagesContainer.value) {
    chatMessagesContainer.value.scrollTop = chatMessagesContainer.value.scrollHeight;
  }
};

// Watch for new messages and scroll to bottom
watch(filteredChatMessages, () => {
  scrollToBottom();
}, { deep: true });

// Scroll to bottom when chat modal opens
watch(chatModal, (isOpen) => {
  if (isOpen) {
    setTimeout(scrollToBottom, 100);
  }
});

// Use event delegation for share button (works even after popup recreates)
document.addEventListener('click', (e) => {
  const target = e.target as HTMLElement;

  if (target && target.id === 'share-location-btn') {
    openShareModal();
  } else if (target && target.id === 'cancel-share-btn') {
    cancelShareLocation();
  } else if (target && target.id === 'chat-btn') {
    openChatModal();
  }
});

// Modal handlers
const openShareModal = () => {
  if (shareLocationModal.value) {
    shareLocationModal.value.showModal();
  }
};

const closeShareModal = () => {
  if (shareLocationModal.value) {
    shareLocationModal.value.close();
    wasteDetails.value = '';
  }
};

const submitShareLocation = async () => {
  if (!wasteDetails.value.trim() || !userLocation.value) {
    return;
  }

  // TODO: Send location and details to backend API
  await customerRequestsMyRequeststore.addCustomerRequest({
    latitude: userLocation.value.lat,
    longitude: userLocation.value.lng,
    description: wasteDetails.value
  });

  if (errorCustomerRequest.value !== "customer request added successfully") {
    alertModal.value?.show(errorCustomerRequest.value, 'Error', 'error');
    return;
  }

  // Update status to pending
  shareStatus.value = 'pending';

  // Show success message
  alertModal.value?.show('แชร์ตำแหน่งสำเร็จ! รอร้านค้าใกล้เคียงตอบรับ', 'Success', 'success');

  closeShareModal();

  // Simulate shop accepting after 5 seconds (for demo)
  // setTimeout(() => {
  //   shareStatus.value = 'accepted';
  //   alertModal.value?.show('ร้านค้ารับคำขอแล้ว! คุณสามารถแชทกับร้านได้', 'Notification', 'success');
  // }, 5000);
};

const cancelShareLocation = () => {
  if (confirm('ต้องการยกเลิกการแชร์ตำแหน่งหรือไม่?')) {
    shareStatus.value = 'idle';
    wasteDetails.value = '';
    alertModal.value?.show('ยกเลิกการแชร์ตำแหน่งแล้ว', 'Info', 'info');
  }
};

const openChatModal = () => {
  if (chatModal.value) {
    chatModal.value.showModal();
  }
};

const closeChatModal = () => {
  if (chatModal.value) {
    chatModal.value.close();
  }
};

const sendMessage = () => {
  if (!currentMessage.value.trim()) return;

  chatStore.sendMessage(currentMessage.value);
  currentMessage.value = '';
};

const handleChatKeyPress = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault();
    sendMessage();
  }
};



const openReviewModalForRequest = (request: CustomerRequest) => {
  currentReviewRequest.value = request;
  rating.value = 5;
  reviewComment.value = '';
  if (reviewModal.value) {
    reviewModal.value.showModal();
  }
};

const closeReviewModal = () => {
  if (reviewModal.value) {
    reviewModal.value.close();
  }
  currentReviewRequest.value = null;
};

const submitReview = async () => {
  if (!reviewComment.value.trim()) {
    alertModal.value?.show('กรุณากรอกความคิดเห็น', 'Warning', 'warning');
    return;
  }

  if (!currentReviewRequest.value) {
    alertModal.value?.show('ไม่พบข้อมูลคำขอ', 'Error', 'error');
    return;
  }

  try {
    // Import postDataWithAuth dynamically
    const { postDataWithAuth } = await import('../services/api');

    const reviewData = {
      customer_request_id: currentReviewRequest.value.customer_request_id,
      rating: rating.value,
      comment: reviewComment.value.trim()
    };

    const response = await postDataWithAuth('/api/reviews', reviewData);

    if (typeof response === 'string') {
      // Error occurred
      alertModal.value?.show(response, 'Error', 'error');
      return;
    }

    // Success
    alertModal.value?.show(`ขอบคุณสำหรับรีวิว ${rating.value} ดาว!`, 'Thank You', 'success');

    // Mark this request as reviewed
    reviewedRequests.value.set(currentReviewRequest.value.customer_request_id, true);

    // Reset everything and close
    closeReviewModal();
    rating.value = 5;
    reviewComment.value = '';

    // Refresh customer requests to update status
    await customerRequestsMyRequeststore.fetchCustomerMyRequests();
  } catch (error) {
    console.error('Error submitting review:', error);
    alertModal.value?.show('เกิดข้อผิดพลาดในการส่งรีวิว', 'Error', 'error');
  }
};

const formatDate = (dateString: string) => {
  if (!dateString) return '-';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('th-TH', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
};

const formatDateTime = (dateString: string) => {
  if (!dateString) return '-';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('th-TH', {
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
};


// Function to refresh customer request status
const refreshCustomerRequestStatus = async () => {
  try {
    await customerRequestsMyRequeststore.fetchCustomerMyRequests();

    // Update shareStatus based on latest request status
    const pendingOrAcceptedReq = customerRequestsMyRequests.value.find(req => req.status === 'pending' || req.status === 'accepted');

    if (pendingOrAcceptedReq) {
      activeRequestId.value = pendingOrAcceptedReq.customer_request_id;

      if (pendingOrAcceptedReq.status === 'pending') {
        shareStatus.value = 'pending';
      } else {
        shareStatus.value = 'accepted';
      }
    } else {
      shareStatus.value = 'idle';
      activeRequestId.value = null;
    }

    console.log('[Auto-refresh] Customer requests updated at', new Date().toLocaleTimeString('th-TH'));
  } catch (error) {
    console.error('[Auto-refresh] Failed to refresh customer requests:', error);
  }
};

watch(filteredShops, () => {
  updateMarkers();
});

// Watch shareStatus to connect/disconnect chat
watch(shareStatus, (newStatus) => {
  if (newStatus === 'accepted' && activeRequestId.value) {
    chatStore.connect(activeRequestId.value);
  } else if (newStatus === 'idle') {
    chatStore.disconnect();
  }
});

// Watch activeRequestId to reconnect if changed (e.g. new request accepted)
watch(activeRequestId, (newId) => {
  if (newId && shareStatus.value === 'accepted') {
    chatStore.connect(newId);
  }
});

onMounted(async () => {
  initMap();
  await shopStore.fetchAllShops();
  await customerRequestsMyRequeststore.fetchCustomerMyRequests();

  // Check for pending requests
  const pendingOrAcceptedReq = customerRequestsMyRequests.value.find(req => req.status === 'pending' || req.status === 'accepted');
  if (pendingOrAcceptedReq) {
    activeRequestId.value = pendingOrAcceptedReq.customer_request_id;

    if (pendingOrAcceptedReq.status === 'pending') {
      shareStatus.value = 'pending';
    } else {
      shareStatus.value = 'accepted';
    }
  }


  console.log("customerRequestsMyRequests", customerRequestsMyRequests.value);


  updateMarkers();

  // Try to locate user automatically
  locateUser();

  // Set up auto-refresh every 15 seconds
  refreshIntervalId.value = window.setInterval(() => {
    refreshCustomerRequestStatus();
  }, 15000); // 15 seconds

  console.log('[Auto-refresh] Started refreshing customer requests every 15 seconds');
});

onUnmounted(() => {
  // Clear auto-refresh interval
  if (refreshIntervalId.value !== null) {
    clearInterval(refreshIntervalId.value);
    console.log('[Auto-refresh] Stopped refreshing customer requests');
  }

  if (map.value) {
    map.value.remove();
  }
});
</script>

<template>
  <div class="h-[calc(100vh-64px)] flex flex-col md:flex-row">
    <!-- Sidebar / List View -->
    <div class="w-full md:w-1/3 lg:w-1/4 bg-base-100 border-r border-base-200 flex flex-col h-1/3 md:h-full">
      <div class="p-4 border-b border-base-200">
        <h1 class="text-xl font-bold text-green-700 mb-2">{{ $t('ShopLocator.title') }}</h1>
        <div class="relative">
          <input type="text" v-model="searchQuery" :placeholder="$t('ShopLocator.search')"
            class="input input-bordered w-full pl-10" />
          <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
            <IconSearch size="18" />
          </span>
        </div>
      </div>

      <div class="flex-1 overflow-y-auto p-4 space-y-4 bg-gradient-to-br from-gray-50 to-gray-100">
        <div v-if="isLoading" class="flex justify-center py-8">
          <span class="loading loading-spinner loading-md text-green-600"></span>
        </div>

        <div v-else-if="filteredShops.length === 0" class="text-center py-8 text-gray-500">
          {{ $t('ShopLocator.noShops') }}
        </div>

        <div v-else v-for="shop in filteredShops" :key="shop.shop_id"
          class="glass-card cursor-pointer group hover:scale-[1.02] transition-all duration-300"
          @click="focusOnShop(shop)" :class="{ 'ring-2 ring-green-500': selectedShopId === shop.shop_id }">
          <div class="p-4 relative z-10">
            <div class="flex gap-3">
              <!-- Shop Logo -->
              <div class="avatar flex-shrink-0">
                <div class="w-16 h-16 rounded-lg shadow-sm overflow-hidden">
                  <img v-if="shop.image_url" :src="shop.image_url" :alt="shop.name"
                    class="object-cover w-full h-full" />
                  <div v-else
                    class="w-full h-full bg-gradient-to-br from-green-100 to-green-200 flex items-center justify-center text-2xl">
                    🏪
                  </div>
                </div>
              </div>

              <!-- Shop Info -->
              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-lg truncate text-gray-800">{{ shop.name }}</h3>
                <p class="text-sm text-gray-600 line-clamp-2">{{ shop.address }}</p>
                <div class="flex flex-wrap items-center gap-2 mt-2 text-xs text-gray-500">
                  <span v-if="shop.opening_time" class="badge badge-ghost badge-sm bg-white/50 gap-1">
                    <IconClock size="12" />
                    {{ shop.opening_time }} - {{ shop.closing_time }}
                  </span>
                  <span v-if="shop.phone" class="badge badge-outline badge-sm border-gray-300 gap-1">
                    <IconPhone size="12" />
                    {{ shop.phone }}
                  </span>
                </div>
              </div>
            </div>

            <div class="flex justify-end mt-3 pt-3 border-t border-gray-100/50 gap-2">
              <router-link :to="{ name: 'shop-profile', params: { shop_id: shop.shop_id } }"
                class="btn btn-xs btn-outline btn-primary" @click.stop>
                Profile
              </router-link>
              <a :href="`https://www.google.com/maps/dir/?api=1&destination=${shop.latitude},${shop.longitude}`"
                target="_blank" class="btn btn-xs btn-primary text-white shadow-sm" @click.stop>
                {{ $t('ShopLocator.navigate') }}
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Map View -->
    <div class="w-full md:w-2/3 lg:w-3/4 h-2/3 md:h-full relative">
      <div id="map" class="w-full h-full z-0"></div>

      <!-- Loading Overlay -->
      <div v-if="isLoading" class="absolute inset-0 bg-white/80 z-10 flex items-center justify-center">
        <div class="flex flex-col items-center gap-2">
          <span class="loading loading-spinner loading-lg text-green-600"></span>
          <span class="text-green-800 font-medium">{{ $t('ShopLocator.loading') }}</span>
        </div>
      </div>

      <!-- Locate Me Button -->
      <button @click="locateUser"
        class="absolute bottom-6 right-6 z-[400] btn btn-circle btn-primary text-white shadow-lg hover:scale-105 transition-transform"
        :title="$t('ShopLocator.myLocation')">
        <IconCurrentLocation />
      </button>

      <!-- View Requests Button -->
      <button @click="openRequestsModal"
        class="absolute top-6 right-6 z-[400] btn btn-primary text-white shadow-lg hover:scale-105 transition-transform gap-2">
        <IconList size="20" />
        ดูคำขอทั้งหมด
      </button>
    </div>

    <!-- Share Location Modal -->
    <dialog ref="shareLocationModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">{{ $t('ShopLocator.shareLocationTitle') }}</h3>

        <div class="form-control">
          <label class="label">
            <span class="label-text">{{ $t('ShopLocator.wasteDetails') }}</span>
            <span class="label-text-alt text-gray-500">{{ wasteDetails.length }}/100</span>
          </label>
          <textarea v-model="wasteDetails" :placeholder="$t('ShopLocator.wasteDetailsPlaceholder')"
            class="textarea textarea-bordered h-24 resize-none" maxlength="100"></textarea>
        </div>

        <div class="modal-action">
          <button @click="closeShareModal" class="btn btn-ghost">
            {{ $t('ShopLocator.cancel') }}
          </button>
          <button @click="submitShareLocation" class="btn btn-primary text-white" :disabled="!wasteDetails.trim()">
            {{ $t('ShopLocator.submit') }}
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="closeShareModal">close</button>
      </form>
    </dialog>

    <!-- Chat Modal -->
    <dialog ref="chatModal" class="modal">
      <div class="modal-box max-w-2xl h-[600px] flex flex-col">
        <div class="flex items-center justify-between mb-4 pb-3 border-b">
          <div class="flex items-center gap-3">
            <div class="avatar online">
              <div class="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center overflow-hidden">
                <img v-if="chatShop?.image_url" :src="chatShop.image_url" :alt="chatShop.name"
                  class="w-full h-full object-cover" />
                <span v-else class="text-xl">🏪</span>
              </div>
            </div>
            <div>
              <h3 class="font-bold text-lg">{{ chatShop?.name || 'แชทกับร้าน' }}</h3>
              <p class="text-xs" :class="chatStore.isPartnerOnline ? 'text-green-500' : 'text-gray-500'">
                {{ chatStore.isPartnerOnline ? 'ออนไลน์' : 'ออฟไลน์' }}
              </p>
            </div>
          </div>
          <button @click="closeChatModal" class="btn btn-sm btn-circle btn-ghost">✕</button>
        </div>

        <div ref="chatMessagesContainer" class="flex-1 overflow-y-auto mb-4 space-y-3 px-2">
          <!-- Privacy Notice -->
          <div class="alert alert-info shadow-sm text-xs py-2 mb-2 bg-blue-50 border-blue-100 text-blue-800">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
              class="stroke-current shrink-0 w-4 h-4">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <span>การสนทนานี้เป็นแบบเรียลไทม์และเน้นความเป็นส่วนตัว ข้อความจะไม่ถูกบันทึก
              ทั้งสองฝ่ายจำเป็นต้องออนไลน์พร้อมกันเพื่อเริ่มการสนทนา</span>
          </div>

          <div v-for="(msg, index) in filteredChatMessages" :key="index"
            :class="msg.sender_id === usersStore.user?.user_id ? 'chat chat-end' : 'chat chat-start'">
            <div class="chat-bubble" :class="msg.sender_id === usersStore.user?.user_id
              ? 'bg-blue-500 text-white'
              : 'bg-gray-200 text-gray-800'">
              {{ msg.message }}
            </div>
            <div class="chat-footer opacity-50 text-xs">
              {{ formatDateTime(msg.timestamp) }}
            </div>
          </div>
        </div>

        <div v-if="!chatStore.isPartnerOnline" class="text-center text-xs text-gray-500 mb-2">
          คู่สนทนาออฟไลน์อยู่ ไม่สามารถส่งข้อความได้ในขณะนี้
        </div>
        <div class="flex gap-2">
          <input v-model="currentMessage" @keypress="handleChatKeyPress" type="text"
            :placeholder="chatStore.isPartnerOnline ? 'พิมพ์ข้อความ...' : 'รอคู่สนทนาออนไลน์...'"
            class="input input-bordered flex-1" :disabled="!chatStore.isPartnerOnline" />
          <button @click="sendMessage" class="btn btn-primary text-white px-4"
            :disabled="!currentMessage.trim() || !chatStore.isPartnerOnline">
            <IconSend :size="20" />
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="closeChatModal">close</button>
      </form>
    </dialog>

    <!-- Review Modal -->
    <dialog ref="reviewModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">⭐ ให้คะแนนบริการ</h3>
        <p class="text-sm text-gray-600 mb-4">บริการเสร็จสิ้น! กรุณาให้คะแนนและรีวิวบริการของเรา</p>

        <!-- Star Rating -->
        <div class="form-control mb-4">
          <label class="label">
            <span class="label-text">ให้คะแนน</span>
          </label>
          <div class="flex gap-2 justify-center">
            <button v-for="star in 5" :key="star" type="button" @click="rating = star"
              class="transition-all hover:scale-125 active:scale-95">
              <IconStarFilled v-if="star <= rating" :size="40" class="text-yellow-400 drop-shadow-md" />
              <IconStar v-else :size="40" class="text-gray-300 hover:text-yellow-200" />
            </button>
          </div>
          <p class="text-center mt-2 text-sm text-gray-500">{{ rating }} ดาว</p>
        </div>

        <!-- Review Comment -->
        <div class="form-control mb-4">
          <label class="label">
            <span class="label-text">ความคิดเห็น</span>
            <span class="label-text-alt text-gray-500">{{ reviewComment.length }}/200</span>
          </label>
          <textarea v-model="reviewComment" placeholder="แบ่งปันประสบการณ์การใช้บริการของคุณ..."
            class="textarea textarea-bordered h-24 resize-none" maxlength="200"></textarea>
        </div>

        <div class="modal-action">
          <button @click="closeReviewModal" class="btn btn-ghost">ข้าม</button>
          <button @click="submitReview" class="btn btn-primary text-white">ส่งรีวิว</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="closeReviewModal">close</button>
      </form>
    </dialog>

    <!-- Requests List Modal -->
    <dialog ref="requestsModal" class="modal">
      <div class="modal-box w-11/12 max-w-5xl">
        <div class="flex justify-between items-center mb-6">
          <h3 class="font-bold text-2xl text-gray-800">รายการคำขอทั้งหมด</h3>
          <button @click="closeRequestsModal" class="btn btn-sm btn-circle btn-ghost">✕</button>
        </div>

        <div class="overflow-x-auto">
          <table class="table w-full">
            <!-- head -->
            <thead>
              <tr class="bg-base-200 text-gray-600 uppercase text-sm">
                <th class="rounded-tl-lg">วันที่</th>
                <th>รายละเอียด</th>
                <th>ตำแหน่ง (Lat, Lng)</th>
                <th>สถานะ</th>
                <th>รีวิว</th>
                <th>ใบเสร็จ</th>
                <th class="rounded-tr-lg text-center">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="customerRequestsMyRequests.length === 0">
                <td colspan="7" class="text-center py-8 text-gray-500">
                  <div class="flex flex-col items-center gap-2">
                    <IconList size="48" class="text-gray-300" />
                    <p>ไม่พบรายการคำขอ</p>
                  </div>
                </td>
              </tr>
              <tr v-else v-for="req in customerRequestsMyRequests" :key="req.customer_request_id"
                class="hover:bg-base-50 transition-colors">
                <td class="text-sm text-gray-500 whitespace-nowrap">
                  {{ formatDate(req.created_at) }}
                </td>
                <td>
                  <div class="font-medium text-gray-800">{{ req.description }}</div>
                </td>
                <td class="text-sm font-mono text-gray-600">
                  {{ req.latitude.toFixed(6) }}, {{ req.longitude.toFixed(6) }}
                </td>
                <td>
                  <div class="badge gap-2" :class="{
                    'badge-warning': req.status === 'pending',
                    'badge-success text-white': req.status === 'accepted',
                    'badge-error text-white': req.status === 'rejected'
                  }">
                    <span v-if="req.status === 'pending'" class="loading loading-spinner loading-xs"></span>
                    {{ req.status }}
                  </div>
                </td>
                <td class="text-center">
                  <!-- Show "รีวิวแล้ว" if reviewed -->
                  <span v-if="req.status === 'done' && reviewedRequests.get(req.customer_request_id)"
                    class="text-xs text-green-600 font-semibold flex items-center justify-center gap-1">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" />
                    </svg>
                    รีวิวแล้ว
                  </span>
                  <!-- Show review button if not reviewed -->
                  <button v-else-if="req.status === 'done'" class="btn btn-xs btn-warning gap-1"
                    @click="openReviewModalForRequest(req)">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
                      <path
                        d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                    </svg>
                    รีวิว
                  </button>
                  <span v-else class="text-xs text-gray-400">-</span>
                </td>
                <td class="text-center">
                  <button v-if="req.status === 'done'" class="btn btn-xs btn-info text-white gap-1"
                    @click="openReceiptModal(req.customer_request_id)">
                    <IconFileText size="14" />
                    ดูใบเสร็จ
                  </button>
                  <span v-else class="text-xs text-gray-400">-</span>
                </td>
                <td class="text-center">
                  <button class="btn btn-ghost btn-xs text-primary" @click="focusOnRequest(req)">
                    ดูตำแหน่ง
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="closeRequestsModal">close</button>
      </form>
    </dialog>

    <!-- Receipt Modal -->
    <dialog ref="receiptModal" class="modal">
      <div class="modal-box w-11/12 max-w-2xl">
        <div class="flex justify-between items-center mb-6">
          <h3 class="font-bold text-2xl text-gray-800">ใบเสร็จรับเงิน</h3>
          <button @click="closeReceiptModal" class="btn btn-sm btn-circle btn-ghost">✕</button>
        </div>

        <div v-if="isLoadingReceipt" class="flex justify-center items-center py-12">
          <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>

        <div v-else-if="receiptData" class="space-y-6">
          <!-- Receipt Info -->
          <div class="bg-base-200 rounded-lg p-4 space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">เลขที่ใบเสร็จ:</span>
              <span class="font-mono font-semibold">{{ receiptData.receipt.id }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">วันที่:</span>
              <span>{{ new Date(receiptData.receipt.created_at).toLocaleString('th-TH') }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">วิธีชำระเงิน:</span>
              <span class="badge badge-success text-white">{{ receiptData.receipt.payment_method }}</span>
            </div>
          </div>

          <!-- Shop Info -->
          <div v-if="receiptData.shop" class="bg-base-200 rounded-lg p-4 space-y-2">
            <h4 class="font-semibold text-base mb-2 text-gray-700">ร้านที่รับซื้อ</h4>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ชื่อร้าน:</span>
              <span class="font-medium">{{ receiptData.shop.name }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ที่อยู่:</span>
              <span class="text-right max-w-xs">{{ receiptData.shop.address }}</span>
            </div>
            <div v-if="receiptData.shop.phone" class="flex justify-between text-sm">
              <span class="text-gray-600">เบอร์โทร:</span>
              <span>{{ receiptData.shop.phone }}</span>
            </div>
          </div>

          <!-- Items -->
          <div>
            <h4 class="font-semibold text-lg mb-3">รายการสินค้า</h4>
            <div class="overflow-x-auto">
              <table class="table w-full">
                <thead>
                  <tr class="bg-base-200">
                    <th>รายการ</th>
                    <th>หมวดหมู่</th>
                    <th class="text-right">น้ำหนัก (kg)</th>
                    <th class="text-right">ราคา/หน่วย (฿)</th>
                    <th class="text-right">รวม (฿)</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(item, idx) in receiptData.items" :key="idx">
                    <td class="font-medium">{{ item.name }}</td>
                    <td><span class="badge badge-sm">{{ item.category }}</span></td>
                    <td class="text-right">{{ item.weight.toFixed(2) }}</td>
                    <td class="text-right">{{ item.unit_price.toFixed(2) }}</td>
                    <td class="text-right font-semibold">{{ item.price.toFixed(2) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Summary -->
          <div class="border-t-2 border-gray-300 pt-4 space-y-2">
            <div class="flex justify-between text-lg">
              <span class="text-gray-600">ยอดรวม (ก่อน VAT):</span>
              <span class="font-semibold">฿{{ receiptData.receipt.total_amount.toFixed(2) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">VAT ({{ (receiptData.receipt.vat_rate * 100).toFixed(0) }}%):</span>
              <span class="font-semibold">฿{{ receiptData.receipt.vat.toFixed(2) }}</span>
            </div>
            <div class="flex justify-between text-xl border-t-2 border-primary pt-3 mt-2">
              <span class="text-gray-800 font-bold">ยอดรวมสุทธิ:</span>
              <span class="text-primary font-bold">฿{{ receiptData.receipt.net_total.toFixed(2) }}</span>
            </div>
          </div>
        </div>

        <div class="modal-action">
          <button @click="downloadPDF" class="btn btn-success text-white gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            ดาวน์โหลด PDF
          </button>
          <button @click="closeReceiptModal" class="btn btn-ghost">ปิด</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="closeReceiptModal">close</button>
      </form>
    </dialog>

    <!-- Alert Modal -->
    <AlertModal ref="alertModal" />
  </div>
</template>

<style scoped>
/* Glassmorphism Popup Styles */
:deep(.glass-popup) {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
}

:deep(.glass-popup .leaflet-popup-content-wrapper) {
  background: transparent !important;
  padding: 0 !important;
  border-radius: 20px !important;
  box-shadow: none !important;
}

:deep(.glass-popup .leaflet-popup-content) {
  margin: 0 !important;
  min-width: 250px !important;
}

:deep(.glass-popup .leaflet-popup-tip-container) {
  display: none;
}

/* Custom Shop Marker Styles */
:deep(.custom-shop-marker) {
  background: transparent !important;
  border: none !important;
}

/* Global Styles for Leaflet Popups (Non-scoped) */
.glass-popup .leaflet-popup-content-wrapper {
  background: transparent !important;
  padding: 0 !important;
  border-radius: 20px !important;
  box-shadow: none !important;
}

.glass-popup .leaflet-popup-content {
  margin: 0 !important;
  min-width: 280px !important;
}

.glass-popup .leaflet-popup-tip-container {
  display: none;
}
</style>
