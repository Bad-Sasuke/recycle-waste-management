<script setup lang="ts">
import { ref, shallowRef, computed, onMounted, onUnmounted, nextTick, watch } from 'vue';
import { useShopStore } from '@/stores/shop';
import { IconCurrentLocation, IconSend } from '@tabler/icons-vue';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import type { CustomerRequest } from '@/types/customer_request';

// Fix Leaflet marker icons
// @ts-ignore
import iconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
// @ts-ignore
import iconUrl from 'leaflet/dist/images/marker-icon.png';
// @ts-ignore
import shadowUrl from 'leaflet/dist/images/marker-shadow.png';
import { useCustomerRequestStore } from '@/stores/customer_request';
import { useChatStore } from '@/stores/chat';
import { useUsersStore } from '../stores/users';
import { storeToRefs } from 'pinia';

delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl,
  iconUrl,
  shadowUrl,
});

// Mock Data
// const customerRequestsMyRequests = ref([
//   {
//     id: '1',
//     customerName: 'สมชาย ใจดี',
//     latitude: 13.7563,
//     longitude: 100.5018,
//     distance: 2.5,
//     wasteAmount: 'medium',
//     wasteTypes: ['plastic', 'paper'],
//     notes: 'มีขยะพลาสติกและกระดาษจำนวนมาก อยากให้มารับภายในวันนี้',
//     phone: '081-234-5678',
//     status: 'pending',
//     createdAt: new Date('2025-11-23T10:30:00'),
//   },
//   {
//     id: '2',
//     customerName: 'สมหญิง รักสะอาด',
//     latitude: 13.7663,
//     longitude: 100.5118,
//     distance: 3.2,
//     wasteAmount: 'large',
//     wasteTypes: ['plastic', 'metal', 'glass'],
//     notes: 'ขยะหลายประเภทค่อนข้างเยอะ',
//     phone: '082-345-6789',
//     status: 'accepted',
//     createdAt: new Date('2025-11-23T09:15:00'),
//   },
//   {
//     id: '3',
//     customerName: 'วิไล สุขใจ',
//     latitude: 13.7463,
//     longitude: 100.4918,
//     distance: 1.8,
//     wasteAmount: 'small',
//     wasteTypes: ['paper'],
//     notes: '',
//     phone: '083-456-7890',
//     status: 'rejected',
//     rejectionReason: 'ปริมาณขยะน้อยเกินไป ค่าขนส่งไม่คุ้ม',
//     createdAt: new Date('2025-11-22T15:45:00'),
//   },
// ]);

const customerRequestStore = useCustomerRequestStore();
const chatStore = useChatStore();
const usersStore = useUsersStore();
const customerRequestAll = computed(() => customerRequestStore.customerRequestsAll);
const { messages: chatMessages } = storeToRefs(chatStore);
const filteredChatMessages = computed(() => chatMessages.value.filter(msg => msg.type === 'message'));

// State
const map = shallowRef<L.Map | null>(null);
const markers = shallowRef<L.Marker[]>([]);
const radiusCircle = shallowRef<L.Circle | null>(null); // Circle overlay for distance radius
const isFirstLoad = ref(true);
const activeChat = ref<CustomerRequest | null>(null);
const newMessage = ref('');

// Distance filter
const maxDistance = ref(20); // Default 20 km
const isChangingDistance = ref(false);

// Modals
const rejectModal = ref<HTMLDialogElement | null>(null);
const cancelModal = ref<HTMLDialogElement | null>(null);
const completeModal = ref<HTMLDialogElement | null>(null);
const rejectionReason = ref('');
const cancellationReason = ref('');
const selectedRequest = ref<CustomerRequest | null>(null);
const chatMessagesContainer = ref<HTMLDivElement | null>(null);

// Computed
const acceptedRequestsList = computed(() =>
  customerRequestAll.value.filter(r => r.status === 'accepted')
);

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

// Watch for active chat changes to scroll to bottom
watch(activeChat, (newVal) => {
  if (newVal) {
    setTimeout(scrollToBottom, 100);
  }
});

// Methods
const initMap = () => {
  // Use shop location from store or fallback to default
  const shopStore = useShopStore();
  const defaultLat = shopStore.shop?.latitude || 13.7563;
  const defaultLng = shopStore.shop?.longitude || 100.5018;

  map.value = L.map('customers-map', {
    zoomControl: false // We can add custom zoom control if needed, or keep default
  }).setView([defaultLat, defaultLng], 14);

  L.control.zoom({
    position: 'topleft'
  }).addTo(map.value);

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  }).addTo(map.value);

  // Add distance radius circle
  radiusCircle.value = L.circle([defaultLat, defaultLng], {
    radius: maxDistance.value * 1000, // Convert km to meters
    color: '#10b981',
    fillColor: '#10b981',
    fillOpacity: 0.1,
    weight: 2,
    opacity: 0.5,
    className: 'radius-circle-animated'
  }).addTo(map.value);

  // Add Shop Marker
  const shopIcon = L.divIcon({
    className: 'custom-customer-marker',
    html: `
        <div class="marker-pin" style="background: #7c3aed; z-index: 1000;">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="white" width="20" height="20">
            <path d="M18.36 9l.6 3H5.04l.6-3h12.72M20 4H4v2h16V4zm0 3H4l-1 5v2h1v6h10v-6h4v6h2v-6h1l-1-5H20zM6 18v-4h6v4H6z"/>
          </svg>
        </div>
      `,
    iconSize: [32, 32],
    iconAnchor: [16, 32],
  });

  L.marker([defaultLat, defaultLng], { icon: shopIcon })
    .addTo(map.value)
    .bindPopup(`<b>${shopStore.shop?.name || 'ตำแหน่งร้านของคุณ'}</b><br>จุดรับซื้อขยะ`)
    .openPopup();

  // Add markers for all customer requests
  updateMapMarkers();
};

const updateMapMarkers = (shouldFitBounds = false) => {
  if (!map.value) return;

  // Clear existing markers
  markers.value.forEach(marker => marker.remove());
  markers.value = [];

  // Add new markers
  customerRequestAll.value.forEach(request => {

    // Skip rejected requests from map
    if (request.status === 'rejected') {
      return;
    }

    const icon = L.divIcon({
      className: 'custom-customer-marker',
      html: `
        <div class="liquid-marker">
          <div class="liquid-pin" style="background: ${getMarkerColor(request.status)}">
            <div class="liquid-glow"></div>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="white" width="20" height="20" style="position: relative; z-index: 2;">
              <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
            </svg>
          </div>
          <div class="liquid-pulse" style="background: ${getMarkerColor(request.status)}"></div>
        </div>
      `,
      iconSize: [40, 40],
      iconAnchor: [20, 40],
    });

    const marker = L.marker([request.latitude, request.longitude], { icon })
      .addTo(map.value!);

    // Create Popup Content
    const popupContent = document.createElement('div');
    popupContent.innerHTML = `
            <div class="glass-card">
                <div class="glass-card-content">
                    <div class="shop-image-container shop-placeholder">
                        <span class="shop-icon">👤</span>
                    </div>
                    <h3 class="shop-name">คำขอรับขยะ #${request.customer_request_id.substring(0, 8)}</h3>
                    <p class="shop-address">
                        <span class="px-2 py-0.5 rounded-full text-xs ${request.status === 'pending' ? 'bg-blue-100 text-blue-700' : 'bg-green-100 text-green-700'}">
                            ${getStatusText(request.status)}
                        </span>
                    </p>

                    <div class="shop-hours">
                        <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                        </svg>
                        <span>${request.description || 'ไม่มีรายละเอียด'}</span>
                    </div>

                    <div class="shop-phone">
                         <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                        </svg>
                        <span>${getDistanceFromShop(request).toFixed(1)} km</span>
                    </div>

                    <div class="shop-phone">
                         <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span>${formatDateTime(request.created_at)}</span>
                    </div>


                    ${request.status === 'pending' ? `
                        <button id="btn-navigate-${request.customer_request_id}" class="navigate-btn" style="background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%); box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);">
                            <svg class="w-4 h-4" style="display: inline-block; margin-right: 4px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
                            </svg>
                            ดูเส้นทาง
                        </button>
                        <button id="btn-accept-${request.customer_request_id}" class="navigate-btn">
                            รับงาน
                        </button>
                        <button id="btn-reject-${request.customer_request_id}" class="navigate-btn" style="background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3); margin-top: 6px; margin-bottom: 0;">
                            ปฏิเสธ
                        </button>
                    ` : ''}

                    ${request.status === 'accepted' ? `
                        <button id="btn-navigate-${request.customer_request_id}" class="navigate-btn" style="background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%); box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);">
                            <svg class="w-4 h-4" style="display: inline-block; margin-right: 4px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
                            </svg>
                            นำทาง
                        </button>
                        <button id="btn-chat-${request.customer_request_id}" class="navigate-btn" style="background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%); box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3); margin-bottom: 0;">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                            </svg>
                            แชท
                        </button>
                    ` : ''}
                </div>
            </div>
        `;

    marker.bindPopup(popupContent, {
      maxWidth: 280,
      className: 'glass-popup'
    });

    // Attach event listeners after popup opens
    marker.on('popupopen', () => {
      const acceptBtn = document.getElementById(`btn-accept-${request.customer_request_id}`);
      const rejectBtn = document.getElementById(`btn-reject-${request.customer_request_id}`);
      const navigateBtn = document.getElementById(`btn-navigate-${request.customer_request_id}`);
      const chatBtn = document.getElementById(`btn-chat-${request.customer_request_id}`);

      if (acceptBtn) {
        acceptBtn.onclick = () => {
          acceptRequest(request);
          marker.closePopup();
        };
      }
      if (rejectBtn) {
        rejectBtn.onclick = () => {
          openRejectModal(request);
          marker.closePopup();
        };
      }
      if (navigateBtn) {
        navigateBtn.onclick = () => {
          navigateToRequest(request);
        };
      }
      if (chatBtn) {
        chatBtn.onclick = () => {
          openChat(request);
          marker.closePopup();
        };
      }
    });

    markers.value.push(marker);
  });

  // Fit map bounds to show all markers (only on first load or when explicitly requested)
  if (shouldFitBounds && markers.value.length > 0) {
    const group = L.featureGroup(markers.value);
    map.value.fitBounds(group.getBounds().pad(0.1)); // 0.1 = 10% padding
  }
};

const getMarkerColor = (status: string) => {
  switch (status) {
    case 'pending': return '#3b82f6'; // blue
    case 'accepted': return '#10b981'; // green
    case 'rejected': return '#ef4444'; // red
    default: return '#6b7280'; // gray
  }
};

const getStatusText = (status: string) => {
  switch (status) {
    case 'pending': return 'รอดำเนินการ';
    case 'accepted': return 'รับงานแล้ว';
    case 'rejected': return 'ปฏิเสธ';
    default: return '-';
  }
};

// Calculate distance between two coordinates using Haversine formula (in km)
const calculateDistance = (lat1: number, lon1: number, lat2: number, lon2: number): number => {
  const R = 6371; // Radius of the Earth in km
  const dLat = (lat2 - lat1) * Math.PI / 180;
  const dLon = (lon2 - lon1) * Math.PI / 180;
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
    Math.sin(dLon / 2) * Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  const distance = R * c;
  return Math.round(distance * 10) / 10; // Round to 1 decimal place
};

const getDistanceFromShop = (request: CustomerRequest) => {
  const shopStore = useShopStore();
  const shopLat = shopStore.shop?.latitude || 13.7563;
  const shopLng = shopStore.shop?.longitude || 100.5018;
  return calculateDistance(shopLat, shopLng, request.latitude, request.longitude);
};

const formatDateTime = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

const formatTime = (date: Date | string) => {
  return new Date(date).toLocaleTimeString('th-TH', {
    hour: '2-digit',
    minute: '2-digit',
  });
};

const viewOnMap = (request: CustomerRequest) => {
  if (map.value) {
    map.value.setView([request.latitude, request.longitude], 16);

    // Find and open the marker popup
    const marker = markers.value.find(m => {
      const latlng = m.getLatLng();
      return Math.abs(latlng.lat - request.latitude) < 0.0001 &&
        Math.abs(latlng.lng - request.longitude) < 0.0001;
    });

    if (marker) {
      marker.openPopup();
    }
  }
};

const acceptRequest = async (request: CustomerRequest) => {
  try {
    await customerRequestStore.acceptRequest(request.customer_request_id);
    // Update map markers to reflect the status change
    updateMapMarkers();
    console.log('Customer request accepted successfully');
  } catch (error) {
    console.error('Failed to accept request:', error);
    alert('ไม่สามารถรับคำขอได้ กรุณาลองใหม่อีกครั้ง');
  }
};

const openRejectModal = (request: CustomerRequest) => {
  selectedRequest.value = request;
  rejectionReason.value = '';
  rejectModal.value?.showModal();
};

const closeRejectModal = () => {
  rejectModal.value?.close();
  selectedRequest.value = null;
  rejectionReason.value = '';
};

const confirmReject = async () => {
  if (!selectedRequest.value) return;
  selectedRequest.value.status = 'rejected';
  // Add custom property for rejection reason (not in API type)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  (selectedRequest.value as any).rejectionReason = rejectionReason.value;
  closeRejectModal();
  updateMapMarkers();
};

const openCancelModal = (request: CustomerRequest) => {
  selectedRequest.value = request;
  cancellationReason.value = '';
  cancelModal.value?.showModal();
};

const closeCancelModal = () => {
  cancelModal.value?.close();
  selectedRequest.value = null;
  cancellationReason.value = '';
};

const confirmCancel = async () => {
  if (!selectedRequest.value) return;

  // Validate reason length (max 100 characters)
  if (cancellationReason.value.length > 100) {
    alert('เหตุผลต้องไม่เกิน 100 ตัวอักษร');
    return;
  }

  if (!cancellationReason.value.trim()) {
    alert('กรุณาระบุเหตุผลในการยกเลิก');
    return;
  }

  try {
    await customerRequestStore.cancelRequest(
      selectedRequest.value.customer_request_id,
      cancellationReason.value
    );
    closeCancelModal();
    updateMapMarkers();
    activeChat.value = null;
    chatStore.disconnect();
    console.log('Customer request cancelled successfully');
  } catch (error) {
    console.error('Failed to cancel request:', error);
    alert(error instanceof Error ? error.message : 'ไม่สามารถยกเลิกคำขอได้ กรุณาลองใหม่อีกครั้ง');
  }
};

const openCompleteModal = (request: CustomerRequest) => {
  selectedRequest.value = request;
  completeModal.value?.showModal();
};

const closeCompleteModal = () => {
  completeModal.value?.close();
  selectedRequest.value = null;
};

const confirmComplete = async () => {
  if (!selectedRequest.value) return;

  try {
    await customerRequestStore.completeRequest(selectedRequest.value.customer_request_id);
    closeCompleteModal();
    updateMapMarkers();
    activeChat.value = null;
    chatStore.disconnect();
    console.log('Customer request completed successfully');
  } catch (error) {
    console.error('Failed to complete request:', error);
    alert(error instanceof Error ? error.message : 'ไม่สามารถปิดงานนี้ได้ กรุณาลองใหม่อีกครั้ง');
  }
};

const openChat = (request: CustomerRequest) => {
  activeChat.value = request;
  chatStore.connect(request.customer_request_id);
};

const sendMessage = () => {
  if (!newMessage.value.trim()) return;

  chatStore.sendMessage(newMessage.value);
  newMessage.value = '';
};

const navigateToRequest = (request: CustomerRequest) => {
  const shopStore = useShopStore();
  const shopLat = shopStore.shop?.latitude || 13.7563;
  const shopLng = shopStore.shop?.longitude || 100.5018;

  // Google Maps navigation URL
  const googleMapsUrl = `https://www.google.com/maps/dir/?api=1&origin=${shopLat},${shopLng}&destination=${request.latitude},${request.longitude}&travelmode=driving`;

  // Open in new tab
  window.open(googleMapsUrl, '_blank');
};

const panToMyLocation = () => {
  const shopStore = useShopStore();
  const lat = shopStore.shop?.latitude || 13.7563;
  const lng = shopStore.shop?.longitude || 100.5018;

  if (map.value) {
    map.value.setView([lat, lng], 16);
  }
};

const updateRadiusCircle = (newRadius: number) => {
  if (!radiusCircle.value) return;

  const currentRadius = radiusCircle.value.getRadius();
  const targetRadius = newRadius * 1000; // Convert km to meters
  const duration = 500; // Animation duration in ms
  const steps = 30;
  const stepDuration = duration / steps;
  const radiusDiff = targetRadius - currentRadius;
  const radiusStep = radiusDiff / steps;

  let currentStep = 0;

  const animate = () => {
    if (currentStep >= steps) {
      radiusCircle.value?.setRadius(targetRadius);
      return;
    }

    const newRadiusValue = currentRadius + (radiusStep * currentStep);
    radiusCircle.value?.setRadius(newRadiusValue);
    currentStep++;

    setTimeout(animate, stepDuration);
  };

  animate();
};

const updateDistanceFilter = async () => {
  isChangingDistance.value = true;
  try {
    // Animate circle radius change
    updateRadiusCircle(maxDistance.value);

    await customerRequestStore.updateMaxDistance(maxDistance.value);
    updateMapMarkers(true);
  } catch (error) {
    console.error('Failed to update distance filter:', error);
  } finally {
    isChangingDistance.value = false;
  }
};

let refreshInterval: ReturnType<typeof setInterval> | null = null;

onMounted(async () => {
  initMap();
  maxDistance.value = customerRequestStore.maxDistance;
  await customerRequestStore.fetchCustomerRequestsAll(1, 10, maxDistance.value)
  // Update markers after data is loaded and fit bounds on first load
  updateMapMarkers(true); // true = fit bounds to show all markers
  isFirstLoad.value = false;

  // Auto-refresh every 15 seconds
  refreshInterval = setInterval(async () => {
    try {
      await customerRequestStore.fetchCustomerRequestsAll(1, 10, maxDistance.value);
      updateMapMarkers(false); // false = don't fit bounds on auto-refresh
    } catch (error) {
      console.error('Failed to auto-refresh customer requests:', error);
    }
  }, 15000); // 15 seconds
});

onUnmounted(() => {
  // Clean up interval when component unmounts
  if (refreshInterval) {
    clearInterval(refreshInterval);
    refreshInterval = null;
  }

  // Clean up map
  if (map.value) {
    map.value.remove();
  }
});
</script>

<template>
  <div class="flex h-[calc(100vh-64px)] bg-white overflow-hidden">
    <!-- Left Sidebar: Accepted Requests -->
    <div class="w-96 flex-shrink-0 flex flex-col border-r border-gray-200 bg-white shadow-xl z-20">
      <!-- Header & Search -->
      <div class="p-4 border-b border-gray-100">
        <h2 class="text-xl font-bold text-green-800 mb-4">{{ $t('NearbyCustomers.title') }}</h2>

        <!-- Distance Filter -->
        <div class="mb-4 p-3 bg-gradient-to-r from-green-50 to-emerald-50 rounded-lg border border-green-100">
          <div class="flex items-center justify-between mb-2">
            <label class="text-sm font-semibold text-green-800 flex items-center gap-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              </svg>
              ระยะทางมองเห็น
            </label>
            <span class="badge badge-success text-white font-bold">{{ maxDistance }} km</span>
          </div>
          <div class="flex gap-2 items-center">
            <input v-model.number="maxDistance" type="range" min="1" max="500"
              class="range range-success range-xs flex-1" @change="updateDistanceFilter" />
          </div>
          <div class="flex justify-between text-xs text-gray-500 mt-1">
            <span>1 km</span>
            <span>500 km</span>
          </div>
          <button @click="updateDistanceFilter" :disabled="isChangingDistance"
            class="btn btn-xs btn-success text-white w-full mt-2">
            <svg v-if="!isChangingDistance" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span v-if="isChangingDistance" class="loading loading-spinner loading-xs"></span>
            {{ isChangingDistance ? 'กำลังโหลด...' : 'อัพเดท' }}
          </button>
        </div>

        <div class="relative">
          <input type="text" :placeholder="$t('NearbyCustomers.searchPlaceholder') || 'ค้นหารายการ...'"
            class="input input-bordered w-full pl-10 bg-gray-50 border-gray-200 focus:bg-white transition-colors" />
          <svg class="w-5 h-5 absolute left-3 top-3 text-gray-400" fill="none" stroke="currentColor"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
      </div>

      <!-- Accepted List -->
      <div class="flex-1 overflow-y-auto p-4 space-y-3 bg-gradient-to-br from-gray-50 to-gray-100">
        <div v-if="acceptedRequestsList.length === 0" class="text-center py-10 text-gray-400">
          <p>{{ $t('NearbyCustomers.noAcceptedRequests') || 'ไม่มีงานที่รับแล้ว' }}</p>
        </div>

        <div v-for="request in acceptedRequestsList" :key="request.customer_request_id"
          class="glass-card cursor-pointer group hover:scale-[1.02] transition-all duration-300"
          @click="viewOnMap(request)"
          :class="{ 'ring-2 ring-green-500': activeChat?.customer_request_id === request.customer_request_id }">
          <div class="p-4 relative z-10">
            <div class="flex items-center gap-3 mb-2">
              <div class="avatar placeholder">
                <div
                  class="bg-green-100 text-green-600 rounded-full w-10 group-hover:bg-green-600 group-hover:text-white transition-colors shadow-sm">
                  <span class="text-sm font-bold">#{{ request.customer_request_id.substring(0, 2).toUpperCase()
                    }}</span>
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-gray-800 truncate">คำขอ #{{ request.customer_request_id.substring(0, 8) }}
                </h3>
                <p class="text-xs text-gray-500">{{ formatDateTime(request.created_at) }}</p>
              </div>
              <div class="badge badge-success badge-sm text-white shadow-sm">
                {{ $t('NearbyCustomers.accepted') }}
              </div>
            </div>

            <div v-if="request.description" class="text-sm text-gray-600 mb-2 line-clamp-2">
              {{ request.description }}
            </div>

            <div class="flex items-center justify-end text-sm text-gray-600 mt-3 pt-3 border-t border-gray-100/50">
              <div class="flex gap-2">
                <button @click.stop="navigateToRequest(request)"
                  class="btn btn-xs btn-info btn-outline gap-1 hover:!text-white px-3 shadow-sm">
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
                  </svg>
                  นำทาง
                </button>
                <button @click.stop="openCompleteModal(request)"
                  class="btn btn-xs btn-success btn-outline gap-1 hover:!text-white px-3 shadow-sm">
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  เสร็จสิ้น
                </button>
                <button @click.stop="openChat(request)"
                  class="btn btn-xs btn-warning btn-outline gap-1 hover:!text-white px-3 shadow-sm">
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                  Chat
                </button>
                <button @click.stop="openCancelModal(request)"
                  class="btn btn-xs btn-error btn-outline gap-1 hover:!text-white px-3 shadow-sm">
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                  ยกเลิก
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Right Map Area -->
    <div class="flex-1 relative bg-gray-100">
      <div id="customers-map" class="w-full h-full z-0"></div>

      <!-- Chat Overlay -->
      <div v-if="activeChat" class="absolute bottom-6 right-6 w-96 z-[1000]">
        <div class="card bg-base-100 shadow-2xl border border-gray-200">
          <div class="card-body p-0 flex flex-col h-[500px]">
            <!-- Chat Header -->
            <div class="p-4 border-b border-gray-200 bg-gradient-to-r from-green-600 to-emerald-600 rounded-t-2xl">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="avatar placeholder">
                    <div class="bg-green-100 text-green-600 rounded-full w-10 shadow-sm">
                      <span class="text-sm font-bold">#{{ activeChat.customer_request_id?.substring(0, 2).toUpperCase()
                      }}</span>
                    </div>
                  </div>
                  <div class="text-white">
                    <h3 class="font-bold">คำขอ #{{ activeChat.customer_request_id?.substring(0, 8) }}</h3>
                    <p class="text-xs opacity-90 flex items-center gap-1">
                      <span class="w-2 h-2 rounded-full"
                        :class="chatStore.isPartnerOnline ? 'bg-green-300 animate-pulse' : 'bg-gray-400'"></span>
                      {{ chatStore.isPartnerOnline ? 'ออนไลน์' : 'ออฟไลน์' }}
                    </p>
                  </div>
                </div>
                <button @click="() => { activeChat = null; chatStore.disconnect(); }"
                  class="btn btn-ghost btn-sm btn-circle text-white hover:bg-white/20">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>

            <!-- Chat Messages -->
            <div ref="chatMessagesContainer" class="flex-1 overflow-y-auto p-4 space-y-3 bg-gray-50">
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

              <div v-for="(message, index) in filteredChatMessages" :key="index" class="flex"
                :class="message.sender_id === usersStore.user?.user_id ? 'justify-end' : 'justify-start'">
                <div class="max-w-[80%]">
                  <div class="px-4 py-2 shadow-sm" :class="message.sender_id === usersStore.user?.user_id
                    ? 'bg-blue-500 text-white rounded-2xl rounded-tr-sm'
                    : 'bg-gray-200 text-gray-800 border border-gray-100 rounded-2xl rounded-tl-sm'">
                    <p class="text-sm">{{ message.message }}</p>
                    <p class="text-[10px] mt-1 opacity-70"
                      :class="message.sender_id === usersStore.user?.user_id ? 'text-blue-100' : 'text-gray-500'">
                      {{ formatTime(message.timestamp) }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Chat Input -->
            <div class="p-4 border-t border-gray-200 bg-white rounded-b-2xl">
              <div v-if="!chatStore.isPartnerOnline" class="text-center text-xs text-gray-500 mb-2">
                คู่สนทนาออฟไลน์อยู่ ไม่สามารถส่งข้อความได้ในขณะนี้
              </div>
              <div class="flex gap-2">
                <input v-model="newMessage" type="text"
                  class="input input-bordered flex-1 input-sm focus:border-green-500 focus:ring-1 focus:ring-green-500"
                  :placeholder="chatStore.isPartnerOnline ? ($t('NearbyCustomers.typeMessage') || 'พิมพ์ข้อความ...') : 'รอคู่สนทนาออนไลน์...'"
                  @keyup.enter="sendMessage" :disabled="!chatStore.isPartnerOnline" />
                <button @click="sendMessage" class="btn btn-primary btn-sm text-white px-4"
                  :disabled="!newMessage.trim() || !chatStore.isPartnerOnline">
                  <IconSend :size="16" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- My Location Button -->
      <button @click="panToMyLocation"
        class="absolute bottom-6 right-6 z-[400] btn btn-circle btn-primary text-white shadow-lg hover:scale-105 transition-transform"
        :class="{ 'bottom-24': activeChat }" :title="$t('NearbyCustomers.myLocation') || 'ตำแหน่งของฉัน'">
        <IconCurrentLocation />
      </button>
    </div>

    <!-- Reject Modal -->
    <dialog ref="rejectModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg text-red-600 mb-4">{{ $t('NearbyCustomers.rejectRequest') }}</h3>
        <div class="form-control">
          <label class="label">
            <span class="label-text font-semibold">{{ $t('NearbyCustomers.rejectionReason') }}</span>
          </label>
          <textarea v-model="rejectionReason" class="textarea textarea-bordered h-24"
            :placeholder="$t('NearbyCustomers.reasonPlaceholder')"></textarea>
        </div>
        <div class="modal-action">
          <button @click="closeRejectModal" class="btn btn-ghost">{{ $t('NearbyCustomers.cancel') }}</button>
          <button @click="confirmReject" class="btn btn-error text-white" :disabled="!rejectionReason.trim()">
            {{ $t('NearbyCustomers.confirmReject') }}
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <!-- Cancel Pickup Modal -->
    <dialog ref="cancelModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg text-orange-600 mb-4">{{ $t('NearbyCustomers.cancelPickupTitle') }}</h3>
        <div class="form-control">
          <label class="label">
            <span class="label-text font-semibold">{{ $t('NearbyCustomers.cancellationReason') }}</span>
            <span class="label-text-alt" :class="cancellationReason.length > 100 ? 'text-error' : 'text-gray-500'">
              {{ cancellationReason.length }}/100
            </span>
          </label>
          <textarea v-model="cancellationReason" class="textarea textarea-bordered h-24"
            :class="{ 'textarea-error': cancellationReason.length > 100 }"
            :placeholder="$t('NearbyCustomers.reasonPlaceholder')" maxlength="100"></textarea>
        </div>
        <div class="modal-action">
          <button @click="closeCancelModal" class="btn btn-ghost">{{ $t('NearbyCustomers.cancel') }}</button>
          <button @click="confirmCancel" class="btn btn-warning text-white" :disabled="!cancellationReason.trim()">
            {{ $t('NearbyCustomers.confirmCancel') }}
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <!-- Complete Request Modal -->
    <dialog ref="completeModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg text-success mb-4">ยืนยันการเสร็จสิ้นงาน</h3>
        <p class="py-4 text-gray-600">
          คุณต้องการทำเครื่องหมายงานนี้เป็น <strong class="text-success">เสร็จสิ้น</strong> ใช่หรือไม่?
        </p>
        <p class="text-sm text-gray-500 mb-4">
          ✓ การดำเนินการนี้จะปิดงานและอัปเดตสถานะเป็นเสร็จสมบูรณ์
        </p>
        <div class="modal-action">
          <button @click="closeCompleteModal" class="btn btn-ghost">ยกเลิก</button>
          <button @click="confirmComplete" class="btn btn-success text-white">
            ยืนยัน เสร็จสิ้น
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>

<style scoped>
:deep(.custom-customer-marker) {
  background: transparent !important;
  border: none !important;
}

:deep(.marker-pin) {
  width: 32px;
  height: 32px;
  border-radius: 50% 50% 50% 0;
  transform: rotate(-45deg);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

:deep(.marker-pin svg) {
  transform: rotate(45deg);
}

/* Radius Circle Animation */
:deep(.radius-circle-animated) {
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.leaflet-interactive) {
  transition: r 0.5s ease-in-out;
}

/* Custom Scrollbar for Sidebar */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}
</style>

<style>
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

/* Liquid Marker Styles */
.liquid-marker {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.liquid-pin {
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

.liquid-pin svg {
  transform: rotate(45deg);
}

.liquid-glow {
  position: absolute;
  top: 5px;
  left: 5px;
  width: 15px;
  height: 15px;
  background: radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.8), transparent);
  border-radius: 50%;
  z-index: 3;
}

.liquid-pulse {
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


/* Navigate Button */
.navigate-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  padding: 10px 18px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  font-weight: 600;
  font-size: 0.875rem;
  border-radius: 10px;
  text-decoration: none;
  margin-top: 6px;
  box-shadow:
    0 4px 12px rgba(16, 185, 129, 0.3),
    0 0 0 1px rgba(16, 185, 129, 0.1) inset;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  border: none;
  cursor: pointer;
}

.navigate-btn:hover {
  transform: translateY(-2px);
  box-shadow:
    0 6px 16px rgba(16, 185, 129, 0.4),
    0 0 0 1px rgba(16, 185, 129, 0.2) inset;
}

.navigate-btn:active {
  transform: translateY(0);
  box-shadow:
    0 2px 8px rgba(16, 185, 129, 0.3),
    0 0 0 1px rgba(16, 185, 129, 0.1) inset;
}
</style>
