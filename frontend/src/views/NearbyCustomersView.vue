<template>
    <div class="min-h-screen bg-gradient-to-br from-emerald-50 via-white to-green-50">
        <div class="max-w-7xl mx-auto px-4 py-8">
            <!-- Header -->
            <div class="mb-8">
                <h1 class="text-3xl font-bold text-green-800 mb-2">
                    {{ $t('NearbyCustomers.title') }}
                </h1>
                <p class="text-gray-600">
                    {{ $t('NearbyCustomers.subtitle') }}
                </p>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <!-- Left: Customer List (2/3 width on large screens) -->
                <div class="lg:col-span-2 space-y-4">
                    <!-- Stats Cards -->
                    <div class="grid grid-cols-3 gap-4 mb-6">
                        <div class="card bg-gradient-to-br from-blue-500 to-blue-600 text-white shadow-lg">
                            <div class="card-body p-4">
                                <div class="flex items-center justify-between">
                                    <div>
                                        <p class="text-sm opacity-90">{{ $t('NearbyCustomers.pending') }}</p>
                                        <p class="text-3xl font-bold">{{ pendingCount }}</p>
                                    </div>
                                    <svg class="w-10 h-10 opacity-80" fill="none" stroke="currentColor"
                                        viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                </div>
                            </div>
                        </div>

                        <div class="card bg-gradient-to-br from-green-500 to-green-600 text-white shadow-lg">
                            <div class="card-body p-4">
                                <div class="flex items-center justify-between">
                                    <div>
                                        <p class="text-sm opacity-90">{{ $t('NearbyCustomers.accepted') }}</p>
                                        <p class="text-3xl font-bold">{{ acceptedCount }}</p>
                                    </div>
                                    <svg class="w-10 h-10 opacity-80" fill="none" stroke="currentColor"
                                        viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                </div>
                            </div>
                        </div>

                        <div class="card bg-gradient-to-br from-gray-500 to-gray-600 text-white shadow-lg">
                            <div class="card-body p-4">
                                <div class="flex items-center justify-between">
                                    <div>
                                        <p class="text-sm opacity-90">{{ $t('NearbyCustomers.total') }}</p>
                                        <p class="text-3xl font-bold">{{ customerRequests.length }}</p>
                                    </div>
                                    <svg class="w-10 h-10 opacity-80" fill="none" stroke="currentColor"
                                        viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                                    </svg>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Filter & Sort -->
                    <div class="flex gap-3 mb-4">
                        <select v-model="filterStatus" class="select select-bordered flex-1">
                            <option value="all">{{ $t('NearbyCustomers.allRequests') }}</option>
                            <option value="pending">{{ $t('NearbyCustomers.pendingOnly') }}</option>
                            <option value="accepted">{{ $t('NearbyCustomers.acceptedOnly') }}</option>
                            <option value="rejected">{{ $t('NearbyCustomers.rejectedOnly') }}</option>
                        </select>
                        <select v-model="sortBy" class="select select-bordered flex-1">
                            <option value="newest">{{ $t('NearbyCustomers.newest') }}</option>
                            <option value="closest">{{ $t('NearbyCustomers.closest') }}</option>
                            <option value="largest">{{ $t('NearbyCustomers.largestWaste') }}</option>
                        </select>
                    </div>

                    <!-- Customer Requests -->
                    <div v-if="filteredRequests.length === 0" class="card bg-base-100 shadow-md">
                        <div class="card-body text-center py-12">
                            <svg class="w-20 h-20 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                            </svg>
                            <p class="text-gray-500 text-lg">{{ $t('NearbyCustomers.noRequests') }}</p>
                        </div>
                    </div>

                    <div v-else class="space-y-4">
                        <div v-for="request in filteredRequests" :key="request.id"
                            class="card bg-base-100 shadow-md hover:shadow-lg transition-shadow border-l-4"
                            :class="getRequestBorderClass(request.status)">
                            <div class="card-body">
                                <!-- Header -->
                                <div class="flex items-start justify-between mb-3">
                                    <div class="flex items-center gap-3">
                                        <div class="avatar placeholder">
                                            <div class="bg-green-100 text-green-600 rounded-full w-12">
                                                <span class="text-xl">👤</span>
                                            </div>
                                        </div>
                                        <div>
                                            <h3 class="font-bold text-lg">{{ request.customerName || 'ลูกค้า' }}</h3>
                                            <p class="text-sm text-gray-500">{{ formatDateTime(request.createdAt) }}</p>
                                        </div>
                                    </div>
                                    <div class="badge badge-lg" :class="getStatusBadgeClass(request.status)">
                                        {{ getStatusText(request.status) }}
                                    </div>
                                </div>

                                <!-- Request Details -->
                                <div class="grid grid-cols-2 gap-4 mb-4">
                                    <!-- Location -->
                                    <div class="flex items-start gap-2">
                                        <svg class="w-5 h-5 text-green-600 mt-0.5" fill="none" stroke="currentColor"
                                            viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                                        </svg>
                                        <div class="text-sm">
                                            <p class="font-semibold text-gray-700">{{ $t('NearbyCustomers.location') }}
                                            </p>
                                            <p class="text-gray-600">{{ request.distance || '2.5' }} km</p>
                                        </div>
                                    </div>

                                    <!-- Waste Amount -->
                                    <div class="flex items-start gap-2">
                                        <svg class="w-5 h-5 text-green-600 mt-0.5" fill="none" stroke="currentColor"
                                            viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                                        </svg>
                                        <div class="text-sm">
                                            <p class="font-semibold text-gray-700">{{ $t('NearbyCustomers.wasteAmount')
                                                }}</p>
                                            <p class="text-gray-600">{{ getWasteAmountText(request.wasteAmount) }}</p>
                                        </div>
                                    </div>

                                    <!-- Waste Types -->
                                    <div class="flex items-start gap-2 col-span-2">
                                        <svg class="w-5 h-5 text-green-600 mt-0.5" fill="none" stroke="currentColor"
                                            viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                                        </svg>
                                        <div class="text-sm flex-1">
                                            <p class="font-semibold text-gray-700 mb-1">{{
                                                $t('NearbyCustomers.wasteTypes') }}</p>
                                            <div class="flex flex-wrap gap-2">
                                                <span v-for="type in request.wasteTypes" :key="type"
                                                    class="badge badge-success badge-sm">
                                                    {{ getWasteTypeText(type) }}
                                                </span>
                                            </div>
                                        </div>
                                    </div>

                                    <!-- Phone -->
                                    <div class="flex items-start gap-2 col-span-2">
                                        <svg class="w-5 h-5 text-green-600 mt-0.5" fill="none" stroke="currentColor"
                                            viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                                        </svg>
                                        <div class="text-sm">
                                            <p class="font-semibold text-gray-700">{{ $t('NearbyCustomers.phone') }}</p>
                                            <p class="text-gray-600">{{ request.phone }}</p>
                                        </div>
                                    </div>

                                    <!-- Notes -->
                                    <div v-if="request.notes" class="flex items-start gap-2 col-span-2">
                                        <svg class="w-5 h-5 text-green-600 mt-0.5" fill="none" stroke="currentColor"
                                            viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
                                        </svg>
                                        <div class="text-sm">
                                            <p class="font-semibold text-gray-700">{{ $t('NearbyCustomers.notes') }}</p>
                                            <p class="text-gray-600">{{ request.notes }}</p>
                                        </div>
                                    </div>
                                </div>

                                <!-- Actions -->
                                <div class="card-actions justify-end pt-3 border-t border-gray-100">
                                    <!-- Pending Actions -->
                                    <template v-if="request.status === 'pending'">
                                        <button @click="viewOnMap(request)" class="btn btn-sm btn-ghost">
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
                                            </svg>
                                            {{ $t('NearbyCustomers.viewOnMap') }}
                                        </button>
                                        <button @click="openRejectModal(request)"
                                            class="btn btn-sm btn-error btn-outline">
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M6 18L18 6M6 6l12 12" />
                                            </svg>
                                            {{ $t('NearbyCustomers.reject') }}
                                        </button>
                                        <button @click="acceptRequest(request)"
                                            class="btn btn-sm btn-success text-white">
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M5 13l4 4L19 7" />
                                            </svg>
                                            {{ $t('NearbyCustomers.accept') }}
                                        </button>
                                    </template>

                                    <!-- Accepted Actions -->
                                    <template v-else-if="request.status === 'accepted'">
                                        <button @click="openChat(request)" class="btn btn-sm btn-primary text-white">
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                                            </svg>
                                            {{ $t('NearbyCustomers.chat') }}
                                        </button>
                                        <button @click="openCancelModal(request)"
                                            class="btn btn-sm btn-error btn-outline">
                                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                    d="M6 18L18 6M6 6l12 12" />
                                            </svg>
                                            {{ $t('NearbyCustomers.cancelPickup') }}
                                        </button>
                                    </template>

                                    <!-- Rejected - Show Reason -->
                                    <template v-else-if="request.status === 'rejected'">
                                        <div class="text-sm text-gray-500">
                                            {{ $t('NearbyCustomers.rejectedReason') }}: {{ request.rejectionReason }}
                                        </div>
                                    </template>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Right: Map & Chat -->
                <div class="space-y-6">
                    <!-- Map -->
                    <div class="card bg-base-100 shadow-lg">
                        <div class="card-body p-0">
                            <div id="customers-map" class="w-full h-[400px] rounded-lg"></div>
                        </div>
                    </div>

                    <!-- Active Chat -->
                    <div v-if="activeChat" class="card bg-base-100 shadow-lg">
                        <div class="card-body p-0 flex flex-col h-[400px]">
                            <!-- Chat Header -->
                            <div class="p-4 border-b border-gray-200 bg-gradient-to-r from-green-500 to-emerald-600">
                                <div class="flex items-center justify-between">
                                    <div class="flex items-center gap-3">
                                        <div class="avatar placeholder">
                                            <div class="bg-white text-green-600 rounded-full w-10">
                                                <span class="text-lg">👤</span>
                                            </div>
                                        </div>
                                        <div class="text-white">
                                            <h3 class="font-bold">{{ activeChat.customerName || 'ลูกค้า' }}</h3>
                                            <p class="text-xs opacity-90">{{ $t('NearbyCustomers.chatting') }}</p>
                                        </div>
                                    </div>
                                    <button @click="activeChat = null"
                                        class="btn btn-ghost btn-sm btn-circle text-white">
                                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M6 18L18 6M6 6l12 12" />
                                        </svg>
                                    </button>
                                </div>
                            </div>

                            <!-- Chat Messages -->
                            <div class="flex-1 overflow-y-auto p-4 space-y-3 bg-gray-50">
                                <div v-for="(message, index) in chatMessages" :key="index" class="flex"
                                    :class="message.sender === 'shop' ? 'justify-end' : 'justify-start'">
                                    <div class="max-w-[70%]">
                                        <div class="px-4 py-2 rounded-2xl" :class="message.sender === 'shop'
                                            ? 'bg-green-500 text-white rounded-br-sm'
                                            : 'bg-white border border-gray-200 rounded-bl-sm'">
                                            <p class="text-sm">{{ message.text }}</p>
                                            <p class="text-xs mt-1 opacity-70">{{ formatTime(message.timestamp) }}</p>
                                        </div>
                                    </div>
                                </div>

                                <!-- Empty State -->
                                <div v-if="chatMessages.length === 0" class="text-center py-8">
                                    <p class="text-gray-400 text-sm">{{ $t('NearbyCustomers.startConversation') }}</p>
                                </div>
                            </div>

                            <!-- Chat Input -->
                            <div class="p-4 border-t border-gray-200 bg-white">
                                <div class="flex gap-2">
                                    <input v-model="newMessage" type="text" class="input input-bordered flex-1 input-sm"
                                        :placeholder="$t('NearbyCustomers.typeMessage')" @keyup.enter="sendMessage" />
                                    <button @click="sendMessage" class="btn btn-primary btn-sm text-white"
                                        :disabled="!newMessage.trim()">
                                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                                        </svg>
                                    </button>
                                </div>
                                <p class="text-xs text-gray-400 mt-2">{{ $t('NearbyCustomers.chatNotSaved') }}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
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
                    </label>
                    <textarea v-model="cancellationReason" class="textarea textarea-bordered h-24"
                        :placeholder="$t('NearbyCustomers.reasonPlaceholder')"></textarea>
                </div>
                <div class="modal-action">
                    <button @click="closeCancelModal" class="btn btn-ghost">{{ $t('NearbyCustomers.cancel') }}</button>
                    <button @click="confirmCancel" class="btn btn-warning text-white"
                        :disabled="!cancellationReason.trim()">
                        {{ $t('NearbyCustomers.confirmCancel') }}
                    </button>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop">
                <button>close</button>
            </form>
        </dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
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

// Mock Data
const customerRequests = ref([
    {
        id: '1',
        customerName: 'สมชาย ใจดี',
        latitude: 13.7563,
        longitude: 100.5018,
        distance: 2.5,
        wasteAmount: 'medium',
        wasteTypes: ['plastic', 'paper'],
        notes: 'มีขยะพลาสติกและกระดาษจำนวนมาก อยากให้มารับภายในวันนี้',
        phone: '081-234-5678',
        status: 'pending',
        createdAt: new Date('2025-11-23T10:30:00'),
    },
    {
        id: '2',
        customerName: 'สมหญิง รักสะอาด',
        latitude: 13.7663,
        longitude: 100.5118,
        distance: 3.2,
        wasteAmount: 'large',
        wasteTypes: ['plastic', 'metal', 'glass'],
        notes: 'ขยะหลายประเภทค่อนข้างเยอะ',
        phone: '082-345-6789',
        status: 'accepted',
        createdAt: new Date('2025-11-23T09:15:00'),
    },
    {
        id: '3',
        customerName: 'วิไล สุขใจ',
        latitude: 13.7463,
        longitude: 100.4918,
        distance: 1.8,
        wasteAmount: 'small',
        wasteTypes: ['paper'],
        notes: '',
        phone: '083-456-7890',
        status: 'rejected',
        rejectionReason: 'ปริมาณขยะน้อยเกินไป ค่าขนส่งไม่คุ้ม',
        createdAt: new Date('2025-11-22T15:45:00'),
    },
]);

// State
const map = ref<L.Map | null>(null);
const markers = ref<L.Marker[]>([]);
const filterStatus = ref('all');
const sortBy = ref('newest');
const activeChat = ref<any>(null);
const chatMessages = ref<any[]>([]);
const newMessage = ref('');

// Modals
const rejectModal = ref<HTMLDialogElement | null>(null);
const cancelModal = ref<HTMLDialogElement | null>(null);
const rejectionReason = ref('');
const cancellationReason = ref('');
const selectedRequest = ref<any>(null);

// Computed
const pendingCount = computed(() =>
    customerRequests.value.filter(r => r.status === 'pending').length
);

const acceptedCount = computed(() =>
    customerRequests.value.filter(r => r.status === 'accepted').length
);

const filteredRequests = computed(() => {
    let requests = customerRequests.value;

    // Filter by status
    if (filterStatus.value !== 'all') {
        requests = requests.filter(r => r.status === filterStatus.value);
    }

    // Sort
    if (sortBy.value === 'newest') {
        requests = [...requests].sort((a, b) =>
            new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
        );
    } else if (sortBy.value === 'closest') {
        requests = [...requests].sort((a, b) => a.distance - b.distance);
    } else if (sortBy.value === 'largest') {
        const wasteOrder = { very_large: 4, large: 3, medium: 2, small: 1 };
        requests = [...requests].sort((a, b) =>
            (wasteOrder[b.wasteAmount as keyof typeof wasteOrder] || 0) -
            (wasteOrder[a.wasteAmount as keyof typeof wasteOrder] || 0)
        );
    }

    return requests;
});

// Methods
const initMap = () => {
    const defaultLat = 13.7563;
    const defaultLng = 100.5018;

    map.value = L.map('customers-map').setView([defaultLat, defaultLng], 12);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map.value);

    // Add markers for all customer requests
    updateMapMarkers();
};

const updateMapMarkers = () => {
    if (!map.value) return;

    // Clear existing markers
    markers.value.forEach(marker => marker.remove());
    markers.value = [];

    // Add new markers
    customerRequests.value.forEach(request => {
        const icon = L.divIcon({
            className: 'custom-customer-marker',
            html: `
        <div class="marker-pin" style="background: ${getMarkerColor(request.status)}">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="white" width="20" height="20">
            <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
          </svg>
        </div>
      `,
            iconSize: [32, 32],
            iconAnchor: [16, 32],
        });

        const marker = L.marker([request.latitude, request.longitude], { icon })
            .addTo(map.value!)
            .bindPopup(`
        <div style="min-width: 200px;">
          <h4 style="font-weight: bold; margin-bottom: 8px;">${request.customerName}</h4>
          <p style="font-size: 12px; color: #666;">
            ${getWasteAmountText(request.wasteAmount)}<br/>
            ${request.distance} km
          </p>
        </div>
      `);

        markers.value.push(marker);
    });
};

const getMarkerColor = (status: string) => {
    switch (status) {
        case 'pending': return '#3b82f6'; // blue
        case 'accepted': return '#10b981'; // green
        case 'rejected': return '#ef4444'; // red
        default: return '#6b7280'; // gray
    }
};

const getRequestBorderClass = (status: string) => {
    switch (status) {
        case 'pending': return 'border-blue-500';
        case 'accepted': return 'border-green-500';
        case 'rejected': return 'border-red-500';
        default: return 'border-gray-300';
    }
};

const getStatusBadgeClass = (status: string) => {
    switch (status) {
        case 'pending': return 'badge-warning';
        case 'accepted': return 'badge-success';
        case 'rejected': return 'badge-error';
        default: return 'badge-ghost';
    }
};

const getStatusText = (status: string) => {
    switch (status) {
        case 'pending': return 'รอดำเนินการ';
        case 'accepted': return 'ยอมรับ';
        case 'rejected': return 'ปฏิเสธ';
        default: return '-';
    }
};

const getWasteAmountText = (amount: string) => {
    const amounts: Record<string, string> = {
        small: 'น้อย',
        medium: 'ปานกลาง',
        large: 'มาก',
        very_large: 'มากที่สุด',
    };
    return amounts[amount] || amount;
};

const getWasteTypeText = (type: string) => {
    const types: Record<string, string> = {
        plastic: 'พลาสติก',
        paper: 'กระดาษ',
        metal: 'โลหะ',
        glass: 'แก้ว',
    };
    return types[type] || type;
};

const formatDateTime = (date: Date) => {
    return new Date(date).toLocaleDateString('th-TH', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });
};

const formatTime = (date: Date) => {
    return new Date(date).toLocaleTimeString('th-TH', {
        hour: '2-digit',
        minute: '2-digit',
    });
};

const viewOnMap = (request: any) => {
    if (map.value) {
        map.value.setView([request.latitude, request.longitude], 15);

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

const acceptRequest = async (request: any) => {
    if (!confirm('คุณแน่ใจหรือไม่ว่าต้องการยอมรับคำขอนี้?')) {
        return;
    }

    // TODO: API call
    request.status = 'accepted';
    updateMapMarkers();
};

const openRejectModal = (request: any) => {
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

    // TODO: API call
    selectedRequest.value.status = 'rejected';
    selectedRequest.value.rejectionReason = rejectionReason.value;

    closeRejectModal();
    updateMapMarkers();
};

const openCancelModal = (request: any) => {
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

    // TODO: API call
    selectedRequest.value.status = 'rejected';
    selectedRequest.value.rejectionReason = cancellationReason.value;

    closeCancelModal();
    updateMapMarkers();
    activeChat.value = null;
};

const openChat = (request: any) => {
    activeChat.value = request;
    chatMessages.value = []; // Load messages from backend
};

const sendMessage = () => {
    if (!newMessage.value.trim()) return;

    chatMessages.value.push({
        sender: 'shop',
        text: newMessage.value,
        timestamp: new Date(),
    });

    newMessage.value = '';

    // Mock customer response
    setTimeout(() => {
        chatMessages.value.push({
            sender: 'customer',
            text: 'ขอบคุณครับ รอคุณอยู่นะครับ',
            timestamp: new Date(),
        });
    }, 1000);
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
</style>
