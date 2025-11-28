<script setup lang="ts">
import { ref, shallowRef, onMounted, computed, watch, onUnmounted, h, render, nextTick } from 'vue';
import { IconSearch, IconCurrentLocation, IconHourglass, IconList, IconSend, IconStar, IconStarFilled } from '@tabler/icons-vue';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

// Fix Leaflet marker icons
// @ts-ignore
import iconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
// @ts-ignore
import iconUrl from 'leaflet/dist/images/marker-icon.png';
// @ts-ignore
import shadowUrl from 'leaflet/dist/images/marker-shadow.png';
import { useShopStore } from '../stores/shop';
import { useCustomerRequestStore } from '../stores/customer_request';
import { useChatStore } from '../stores/chat';
import { useUsersStore } from '../stores/users';
import { storeToRefs } from 'pinia';
import AlertModal from '../components/AlertModal.vue';
import type { Shop } from '../types/shop';

// Fix for default marker icon
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

const currentMessage = ref('');
const reviewModal = ref<HTMLDialogElement | null>(null);
const rating = ref(5);
const reviewComment = ref('');
const messageCount = ref(0);
const alertModal = ref<InstanceType<typeof AlertModal> | null>(null);
const requestsModal = ref<HTMLDialogElement | null>(null);
const refreshIntervalId = ref<number | null>(null);
const currentReviewRequest = ref<any | null>(null);

const openRequestsModal = () => {
  if (requestsModal.value) {
    requestsModal.value.showModal();
  }
};

const closeRequestsModal = () => {
  if (requestsModal.value) {
    requestsModal.value.close();
  }
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
const getUserPopupContent = computed(() => {
  const baseContent = `
    <div class="glass-card">
      <div class="glass-card-content" style="padding: 16px; text-align: center;">
        <div style="width: 48px; height: 48px; margin: 0 auto 12px; background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%); border-radius: 50%; display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);">
          <svg style="width: 24px; height: 24px; color: white;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
        </div>
        <h3 class="shop-name" style="margin-bottom: 4px;">คุณอยู่ที่นี่</h3>
        <p class="shop-address" style="font-size: 0.75rem; margin-bottom: 12px;">ตำแหน่งปัจจุบันของคุณ</p>`;

  if (shareStatus.value === 'idle') {
    // ปุ่มแชร์ตำแหน่งปกติ (สีเขียว)
    return baseContent + `
        <button id="share-location-btn" style="width: 100%; padding: 8px 16px; background: linear-gradient(135deg, #10b981 0%, #059669 100%); color: white; font-weight: 600; font-size: 0.875rem; border: none; border-radius: 8px; cursor: pointer; transition: all 0.2s;">
          แชร์ตำแหน่งของคุณ
        </button>
      </div>
    </div>`;
  } else if (shareStatus.value === 'pending') {
    // Render icon to HTML string
    const div = document.createElement('div');
    render(h(IconHourglass, { size: 16, class: 'animate-spin' }), div);
    const iconHtml = div.innerHTML;

    // สถานะรอร้าน: ปุ่มรอ (สีเขียวแต่ disabled)
    return baseContent + `
        <p style="font-size: 0.75rem; color: #f59e0b; margin-bottom: 8px; font-weight: 600; display: flex; align-items: center; justify-content: center; gap: 4px;">
          ${iconHtml}
          รอร้านค้าใกล้เคียง
        </p>
        <button disabled style="width: 100%; padding: 8px 16px; background: linear-gradient(135deg, #cccccc 0%, #cccccc 100%); color: white; font-weight: 600; font-size: 0.875rem; border: none; border-radius: 8px; cursor: not-allowed; opacity: 0.7; transition: all 0.2s;">
          แชร์ตำแหน่งของคุณ
        </button>
      </div>
    </div>`;
  } else {
    // สถานะร้านรับแล้ว: ปุ่มแชท + ปุ่มยกเลิก (disabled)
    return baseContent + `
        <p style="font-size: 0.75rem; color: #10b981; margin-bottom: 8px; font-weight: 600;">✓ ร้านรับคำขอแล้ว</p>
        <button id="chat-btn" style="width: 100%; padding: 8px 16px; margin-bottom: 8px; background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%); color: white; font-weight: 600; font-size: 0.875rem; border: none; border-radius: 8px; cursor: pointer; transition: all 0.2s;">
          💬 แชทกับร้าน
        </button>
        <button disabled style="width: 100%; padding: 8px 16px; background: #9ca3af; color: white; font-weight: 600; font-size: 0.875rem; border: none; border-radius: 8px; cursor: not-allowed; opacity: 0.6;">
          ยกเลิก
        </button>
      </div>
    </div>`;
  }
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
      const customIcon = L.divIcon({
        className: 'custom-shop-marker',
        html: `
          <div class="liquid-marker">
            <div class="liquid-pin" style="background: #10b981;">
              <div class="liquid-glow"></div>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="white" width="20" height="20" style="position: relative; z-index: 2;">
                <path d="M21 13.242V20H22V22H2V20H3V13.242L1.515 9.583C1.195 8.772 1.633 7.862 2.444 7.541C2.626 7.466 2.821 7.428 3.018 7.428H20.982C21.889 7.428 22.623 8.163 22.623 9.07C22.623 9.267 22.585 9.462 22.51 9.644L21 13.242ZM19 13.972L20.234 10.428H3.766L5 13.972V20H19V13.972ZM14 10H10V12C10 13.105 10.895 14 12 14C13.105 14 14 13.105 14 12V10Z"/>
              </svg>
            </div>
            <div class="liquid-pulse" style="background: #10b981;"></div>
          </div>
        `,
        iconSize: [40, 40],
        iconAnchor: [20, 40],
        popupAnchor: [0, -40]
      });

      const marker = L.marker([shop.latitude, shop.longitude], {
        icon: customIcon
      })
        .addTo(map.value!)
        .bindPopup(`
          <div class="glass-card">
            <div class="glass-card-content">
              ${shop.image_url
            ? `<div class="shop-image-container">
                     <img src="${shop.image_url}" alt="${shop.name}" class="shop-image" />
                   </div>`
            : `<div class="shop-image-container shop-placeholder">
                     <span class="shop-icon">🏪</span>
                   </div>`
          }
              <h3 class="shop-name">${shop.name}</h3>
              <p class="shop-address">${shop.address}</p>
              ${shop.opening_time ? `
                <div class="shop-hours">
                  <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  <span>${shop.opening_time} - ${shop.closing_time}</span>
                </div>
              ` : ''}
              ${shop.phone ? `
                <div class="shop-phone">
                  <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/>
                  </svg>
                  <span>${shop.phone}</span>
                </div>
              ` : ''}
              <a href="https://www.google.com/maps/dir/?api=1&destination=${shop.latitude},${shop.longitude}"
                 target="_blank"
                 class="navigate-btn">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"/>
                </svg>
                Navigate
              </a>
            </div>
          </div>
        `, {
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
const focusOnRequest = (req: any) => {
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
      fillColor: '#3b82f6',
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

const openReviewModal = () => {
  if (reviewModal.value) {
    reviewModal.value.showModal();
  }
};

const openReviewModalForRequest = (request: any) => {
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
                  <span v-if="shop.opening_time" class="badge badge-ghost badge-sm bg-white/50">
                    {{ shop.opening_time }} - {{ shop.closing_time }}
                  </span>
                  <span v-if="shop.phone" class="badge badge-outline badge-sm border-gray-300">
                    {{ shop.phone }}
                  </span>
                </div>
              </div>
            </div>

            <div class="flex justify-end mt-3 pt-3 border-t border-gray-100/50">
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
                <th class="rounded-tr-lg text-center">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="customerRequestsMyRequests.length === 0">
                <td colspan="5" class="text-center py-8 text-gray-500">
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
                  <button v-if="req.status === 'done'" class="btn btn-xs btn-warning gap-1"
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

/* Glass Card */
:deep(.glass-card) {
  background: rgba(255, 255, 255, 0.75);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow:
    0 8px 32px 0 rgba(31, 38, 135, 0.15),
    0 0 0 1px rgba(255, 255, 255, 0.1) inset,
    0 1px 2px 0 rgba(0, 0, 0, 0.05);
  overflow: hidden;
  position: relative;
}

:deep(.glass-card::before) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg,
      transparent,
      rgba(255, 255, 255, 0.8),
      transparent);
}

:deep(.glass-card-content) {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Shop Image */
:deep(.shop-image-container) {
  width: 100%;
  height: 140px;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

:deep(.shop-image) {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

:deep(.shop-placeholder) {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

:deep(.shop-icon) {
  font-size: 3.5rem;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

/* Shop Info */
:deep(.shop-name) {
  font-size: 1.125rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
  line-height: 1.4;
}

:deep(.shop-address) {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

/* Info Row with Icons */
:deep(.shop-hours),
:deep(.shop-phone) {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.75rem;
  color: #4b5563;
  padding: 6px 12px;
  background: rgba(16, 185, 129, 0.08);
  border-radius: 8px;
  border: 1px solid rgba(16, 185, 129, 0.15);
}

:deep(.shop-hours .icon),
:deep(.shop-phone .icon) {
  width: 14px;
  height: 14px;
  color: #10b981;
  flex-shrink: 0;
}

/* Navigate Button */
:deep(.navigate-btn) {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px 20px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  font-weight: 600;
  font-size: 0.875rem;
  border-radius: 12px;
  text-decoration: none;
  margin-top: 8px;
  box-shadow:
    0 4px 12px rgba(16, 185, 129, 0.3),
    0 1px 3px rgba(0, 0, 0, 0.1) inset;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

:deep(.navigate-btn::before) {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.5s;
}

:deep(.navigate-btn:hover) {
  transform: translateY(-2px);
  box-shadow:
    0 6px 20px rgba(16, 185, 129, 0.4),
    0 2px 6px rgba(0, 0, 0, 0.15) inset;
}

:deep(.navigate-btn:hover::before) {
  left: 100%;
}

:deep(.navigate-btn:active) {
  transform: translateY(0);
  box-shadow:
    0 2px 8px rgba(16, 185, 129, 0.4),
    0 1px 2px rgba(0, 0, 0, 0.1) inset;
}

:deep(.navigate-btn .icon) {
  width: 16px;
  height: 16px;
}

/* Custom Shop Marker Styles */
:deep(.custom-shop-marker) {
  background: transparent !important;
  border: none !important;
}

/* Liquid Marker Styles */
:deep(.liquid-marker) {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

:deep(.liquid-pin) {
  width: 40px;
  height: 40px;
  border-radius: 50% 50% 50% 0;
  transform: rotate(-45deg);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow:
    0 4px 12px rgba(0, 0, 0, 0.2),
    inset 0 2px 4px rgba(255, 255, 255, 0.3),
    inset 0 -2px 4px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 2;
  animation: float 3s ease-in-out infinite;
}

:deep(.liquid-pin svg) {
  transform: rotate(45deg);
}

:deep(.liquid-glow) {
  position: absolute;
  top: 5px;
  left: 5px;
  width: 15px;
  height: 15px;
  background: radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.8), transparent);
  border-radius: 50%;
  z-index: 3;
}

:deep(.liquid-pulse) {
  position: absolute;
  bottom: -10px;
  left: 50%;
  transform: translateX(-50%) rotateX(60deg);
  width: 30px;
  height: 30px;
  border-radius: 50%;
  opacity: 0.5;
  filter: blur(4px);
  animation: pulse-shadow 3s ease-in-out infinite;
  z-index: 1;
}

@keyframes float {

  0%,
  100% {
    transform: translateY(0) rotate(-45deg);
  }

  50% {
    transform: translateY(-6px) rotate(-45deg);
  }
}

@keyframes pulse-shadow {

  0%,
  100% {
    transform: translateX(-50%) rotateX(60deg) scale(1);
    opacity: 0.5;
  }

  50% {
    transform: translateX(-50%) rotateX(60deg) scale(0.8);
    opacity: 0.2;
  }
}

/* Share Location Button Styles */
:deep(#share-location-btn:hover) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(16, 185, 129, 0.4);
}

:deep(#share-location-btn:active) {
  transform: translateY(0);
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

/* Liquid Glass Card */
.glass-card {
  background: rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow:
    0 8px 32px 0 rgba(31, 38, 135, 0.15),
    0 0 0 1px rgba(255, 255, 255, 0.1) inset,
    0 1px 2px 0 rgba(0, 0, 0, 0.05);
  overflow: hidden;
  position: relative;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.glass-card:hover {
  transform: translateY(-2px);
  box-shadow:
    0 12px 40px 0 rgba(31, 38, 135, 0.2),
    0 0 0 1px rgba(255, 255, 255, 0.2) inset;
}

.glass-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -50%;
  width: 200%;
  height: 100%;
  background: linear-gradient(to right,
      transparent,
      rgba(255, 255, 255, 0.3),
      transparent);
  transform: skewX(-25deg);
  animation: shine 6s infinite;
  pointer-events: none;
}

@keyframes shine {
  0% {
    transform: translateX(-100%) skewX(-25deg);
  }

  20% {
    transform: translateX(100%) skewX(-25deg);
  }

  100% {
    transform: translateX(100%) skewX(-25deg);
  }
}

.glass-card-content {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Shop Image */
.shop-image-container {
  width: 100%;
  height: 140px;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.shop-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.shop-placeholder {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.shop-icon {
  font-size: 3.5rem;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

/* Shop Info */
.shop-name {
  font-size: 1.125rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
  line-height: 1.4;
}

.shop-address {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

/* Info Row with Icons */
.shop-hours,
.shop-phone {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.75rem;
  color: #4b5563;
  padding: 6px 12px;
  background: rgba(16, 185, 129, 0.08);
  border-radius: 8px;
  border: 1px solid rgba(16, 185, 129, 0.15);
}

.shop-hours .icon,
.shop-phone .icon {
  width: 14px;
  height: 14px;
  color: #10b981;
  flex-shrink: 0;
}
</style>
