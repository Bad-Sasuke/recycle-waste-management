<script setup lang="ts">
import GlassCard from '../common/GlassCard.vue'
import { IconClock, IconNavigation, IconCheck, IconX, IconMessageCircle } from '@tabler/icons-vue'
import type { PropType } from 'vue'

// We need to define or import CustomerRequest type.
// Assuming it's available or we define a minimal one.
// The user has 'frontend/src/types/customer_request.ts'.
import type { CustomerRequest } from '@/types/customer_request'

defineProps({
    request: {
        type: Object as PropType<CustomerRequest>,
        required: true
    },
    distance: {
        type: Number,
        required: true
    }
})

const getStatusText = (status: string) => {
    switch (status) {
        case 'pending': return 'รอดำเนินการ';
        case 'accepted': return 'รับงานแล้ว';
        case 'rejected': return 'ปฏิเสธ';
        default: return '-';
    }
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
</script>

<template>
    <GlassCard>
        <div class="shop-image-container shop-placeholder">
            <span class="shop-icon">👤</span>
        </div>
        <h3 class="shop-name">คำขอรับขยะ #{{ request.customer_request_id.substring(0, 8) }}</h3>
        <p class="shop-address">
            <span class="px-2 py-0.5 rounded-full text-xs"
                :class="request.status === 'pending' ? 'bg-blue-100 text-blue-700' : 'bg-green-100 text-green-700'">
                {{ getStatusText(request.status) }}
            </span>
        </p>

        <div class="shop-hours">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
            </svg>
            <span>{{ request.description || 'ไม่มีรายละเอียด' }}</span>
        </div>

        <div class="shop-phone">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
            </svg>
            <span>{{ distance.toFixed(1) }} km</span>
        </div>

        <div class="shop-phone">
            <IconClock class="icon" />
            <span>{{ formatDateTime(request.created_at) }}</span>
        </div>

        <div v-if="request.status === 'pending'" class="flex flex-col gap-2 mt-2">
            <button :id="`btn-navigate-${request.customer_request_id}`"
                class="navigate-btn bg-gradient-warning shadow-orange">
                <IconNavigation class="w-4 h-4" />
                ดูเส้นทาง
            </button>
            <button :id="`btn-accept-${request.customer_request_id}`"
                class="navigate-btn bg-gradient-success shadow-green">
                <IconCheck class="w-4 h-4" />
                รับงาน
            </button>
            <button :id="`btn-reject-${request.customer_request_id}`"
                class="navigate-btn bg-gradient-danger shadow-red">
                <IconX class="w-4 h-4" />
                ปฏิเสธ
            </button>
        </div>

        <div v-if="request.status === 'accepted'" class="flex flex-col gap-2 mt-2">
            <button :id="`btn-navigate-${request.customer_request_id}`"
                class="navigate-btn bg-gradient-warning shadow-orange">
                <IconNavigation class="w-4 h-4" />
                นำทาง
            </button>
            <button :id="`btn-chat-${request.customer_request_id}`" class="navigate-btn bg-gradient-info shadow-blue">
                <IconMessageCircle class="w-4 h-4" />
                แชท
            </button>
        </div>
    </GlassCard>
</template>

<style scoped>
/* Shop Image */
.shop-image-container {
    width: 100%;
    height: 140px;
    border-radius: 12px;
    overflow: hidden;
    margin-bottom: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
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
    color: white;
    font-weight: 600;
    font-size: 0.875rem;
    border-radius: 10px;
    text-decoration: none;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    border: none;
    cursor: pointer;
}

.navigate-btn:hover {
    transform: translateY(-2px);
}

.navigate-btn:active {
    transform: translateY(0);
}
</style>
