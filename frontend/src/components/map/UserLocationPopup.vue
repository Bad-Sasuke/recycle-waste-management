<script setup lang="ts">
import GlassCard from '../common/GlassCard.vue'
import { IconHourglass } from '@tabler/icons-vue'

defineProps({
    status: {
        type: String,
        required: true, // 'idle', 'pending', 'accepted'
        default: 'idle'
    }
})
</script>

<template>
    <GlassCard>
        <div class="w-12 h-12 mx-auto mb-3 bg-gradient-info rounded-full flex items-center justify-center shadow-blue">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
        </div>
        <h3 class="shop-name mb-1 text-center">คุณอยู่ที่นี่</h3>
        <p class="shop-address text-xs mb-3 text-center">ตำแหน่งปัจจุบันของคุณ</p>

        <div v-if="status === 'idle'">
            <button id="share-location-btn"
                class="w-full px-4 py-2 bg-gradient-primary text-white font-semibold text-sm border-none rounded-lg cursor-pointer transition-all duration-200 hover:-translate-y-0.5 shadow-md">
                แชร์ตำแหน่งของคุณ
            </button>
        </div>

        <div v-else-if="status === 'pending'">
            <p class="text-xs text-warning mb-2 font-semibold flex items-center justify-center gap-1">
                <IconHourglass :size="16" class="animate-spin" />
                รอร้านค้าใกล้เคียง
            </p>
            <button disabled
                class="w-full px-4 py-2 bg-disabled text-white font-semibold text-sm border-none rounded-lg cursor-not-allowed opacity-70 transition-all duration-200">
                แชร์ตำแหน่งของคุณ
            </button>
        </div>

        <div v-else>
            <p class="text-xs text-success mb-2 font-semibold text-center">✓ ร้านรับคำขอแล้ว</p>
            <button id="chat-btn"
                class="w-full px-4 py-2 mb-2 bg-gradient-info text-white font-semibold text-sm border-none rounded-lg cursor-pointer transition-all duration-200 hover:-translate-y-0.5 shadow-md">
                💬 แชทกับร้าน
            </button>
            <button disabled
                class="w-full px-4 py-2 bg-gray-400 text-white font-semibold text-sm border-none rounded-lg cursor-not-allowed opacity-60">
                ยกเลิก
            </button>
        </div>
    </GlassCard>
</template>

<style scoped>
.shop-name {
    font-size: 1.125rem;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0;
    line-height: 1.4;
}

.shop-address {
    font-size: 0.875rem;
    color: var(--text-muted);
    margin: 0;
    line-height: 1.5;
}
</style>
