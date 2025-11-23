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

      <div class="flex-1 overflow-y-auto p-4 space-y-4">
        <div v-if="isLoading" class="flex justify-center py-8">
          <span class="loading loading-spinner loading-md text-green-600"></span>
        </div>

        <div v-else-if="filteredShops.length === 0" class="text-center py-8 text-gray-500">
          {{ $t('ShopLocator.noShops') }}
        </div>

        <div v-else v-for="shop in filteredShops" :key="shop.shop_id"
          class="card bg-base-100 shadow-sm border border-base-200 hover:shadow-md transition-shadow cursor-pointer"
          @click="focusOnShop(shop)"
          :class="{ 'border-green-500 ring-1 ring-green-500': selectedShopId === shop.shop_id }">
          <div class="card-body p-4">
            <div class="flex gap-3">
              <!-- Shop Logo -->
              <div class="avatar flex-shrink-0">
                <div class="w-16 h-16 rounded-lg">
                  <img v-if="shop.image_url" :src="shop.image_url" :alt="shop.name" class="object-cover" />
                  <div v-else class="w-full h-full bg-green-100 flex items-center justify-center text-2xl">
                    🏪
                  </div>
                </div>
              </div>

              <!-- Shop Info -->
              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-lg truncate">{{ shop.name }}</h3>
                <p class="text-sm text-gray-600 line-clamp-2">{{ shop.address }}</p>
                <div class="flex flex-wrap items-center gap-2 mt-2 text-xs text-gray-500">
                  <span v-if="shop.opening_time" class="badge badge-ghost badge-sm">
                    {{ shop.opening_time }} - {{ shop.closing_time }}
                  </span>
                  <span v-if="shop.phone" class="badge badge-outline badge-sm">
                    {{ shop.phone }}
                  </span>
                </div>
              </div>
            </div>

            <div class="card-actions justify-end mt-2">
              <a :href="`https://www.google.com/maps/dir/?api=1&destination=${shop.latitude},${shop.longitude}`"
                target="_blank" class="btn btn-xs btn-primary text-white" @click.stop>
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch, onUnmounted } from 'vue';
import { useShopStore } from '@/stores/shop';
import { IconSearch, IconCurrentLocation } from '@tabler/icons-vue';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import type { Shop } from '@/types/shop';

// Fix Leaflet marker icons
// @ts-ignore
import iconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png';
// @ts-ignore
import iconUrl from 'leaflet/dist/images/marker-icon.png';
// @ts-ignore
import shadowUrl from 'leaflet/dist/images/marker-shadow.png';

// Fix for default marker icon
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl,
  iconUrl,
  shadowUrl,
});

const shopStore = useShopStore();
const searchQuery = ref('');
const map = ref<L.Map | null>(null);
const markers = ref<L.Marker[]>([]);
const selectedShopId = ref<string | null>(null);
const userLocation = ref<{ lat: number; lng: number } | null>(null);

const isLoading = computed(() => shopStore.isLoading);

const filteredShops = computed(() => {
  if (!searchQuery.value) return shopStore.allShops;
  const query = searchQuery.value.toLowerCase();
  return shopStore.allShops.filter(shop =>
    shop.name?.toLowerCase().includes(query) ||
    shop.address?.toLowerCase().includes(query)
  );
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
      // Create custom icon with store SVG
      const customIcon = L.divIcon({
        className: 'custom-shop-marker',
        html: `
          <div class="marker-container">
            <div class="marker-pin">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                <path d="M21 13.242V20H22V22H2V20H3V13.242L1.515 9.583C1.195 8.772 1.633 7.862 2.444 7.541C2.626 7.466 2.821 7.428 3.018 7.428H20.982C21.889 7.428 22.623 8.163 22.623 9.07C22.623 9.267 22.585 9.462 22.51 9.644L21 13.242ZM19 13.972L20.234 10.428H3.766L5 13.972V20H19V13.972ZM14 10H10V12C10 13.105 10.895 14 12 14C13.105 14 14 13.105 14 12V10Z"/>
              </svg>
            </div>
            <div class="marker-pulse"></div>
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
                <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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

// Focus map on specific shop
const focusOnShop = (shop: Shop) => {
  selectedShopId.value = shop.shop_id || null;
  if (map.value && shop.latitude && shop.longitude) {
    map.value.setView([shop.latitude, shop.longitude], 15);

    // Find and open popup
    const marker = markers.value.find(m => {
      const latLng = m.getLatLng();
      return Math.abs(latLng.lat - shop.latitude) < 0.0001 &&
        Math.abs(latLng.lng - shop.longitude) < 0.0001;
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

  map.value.on('locationfound', (e) => {
    userLocation.value = e.latlng;

    // Add user marker with glassmorphism popup
    L.circleMarker(e.latlng, {
      radius: 8,
      fillColor: '#3b82f6',
      color: '#fff',
      weight: 2,
      opacity: 1,
      fillOpacity: 0.8
    }).addTo(map.value!)
      .bindPopup(`
        <div class="glass-card">
          <div class="glass-card-content" style="padding: 16px; text-align: center;">
            <div style="width: 48px; height: 48px; margin: 0 auto 12px; background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%); border-radius: 50%; display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);">
              <svg style="width: 24px; height: 24px; color: white;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
              </svg>
            </div>
            <h3 class="shop-name" style="margin-bottom: 4px;">คุณอยู่ที่นี่</h3>
            <p class="shop-address" style="font-size: 0.75rem;">ตำแหน่งปัจจุบันของคุณ</p>
          </div>
        </div>
      `, {
        maxWidth: 220,
        className: 'glass-popup'
      })
      .openPopup();
  });

  map.value.on('locationerror', (e) => {
    console.error('Location error:', e.message);
    // Fallback to default view if location access denied
  });
};

watch(filteredShops, () => {
  updateMarkers();
});

onMounted(async () => {
  initMap();
  await shopStore.fetchAllShops();
  updateMarkers();

  // Try to locate user automatically
  locateUser();
});

onUnmounted(() => {
  if (map.value) {
    map.value.remove();
  }
});
</script>

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

:deep(.marker-container) {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

:deep(.marker-pin) {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  border-radius: 50% 50% 50% 0;
  transform: rotate(-45deg);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow:
    0 4px 12px rgba(16, 185, 129, 0.4),
    0 2px 6px rgba(0, 0, 0, 0.15),
    0 0 0 2px rgba(255, 255, 255, 0.9);
  position: relative;
  z-index: 2;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.marker-pin:hover) {
  transform: rotate(-45deg) scale(1.1);
  box-shadow:
    0 6px 20px rgba(16, 185, 129, 0.5),
    0 4px 10px rgba(0, 0, 0, 0.2),
    0 0 0 3px rgba(255, 255, 255, 1);
}

:deep(.marker-pin svg) {
  width: 18px;
  height: 18px;
  color: white;
  transform: rotate(45deg);
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.2));
}

/* Pulse Animation */
:deep(.marker-pulse) {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 36px;
  height: 36px;
  background: rgba(16, 185, 129, 0.4);
  border-radius: 50%;
  z-index: 1;
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {

  0%,
  100% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 0.6;
  }

  50% {
    transform: translate(-50%, -50%) scale(1.5);
    opacity: 0;
  }
}
</style>
