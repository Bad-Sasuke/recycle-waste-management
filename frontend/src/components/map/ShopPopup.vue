<script setup lang="ts">
import GlassCard from '../common/GlassCard.vue'
import NavigateBtn from '../common/NavigateBtn.vue'
import { IconClock, IconPhone, IconBuildingStore, IconNavigation } from '@tabler/icons-vue'
import type { PropType } from 'vue'

// Define a minimal Shop interface for the prop since we might not have the full type available here easily
// or we can import it if we know the path. Let's try to be generic or import if possible.
// The user has 'frontend/src/types/shop.ts'.
import type { Shop } from '@/types/shop'

defineProps({
    shop: {
        type: Object as PropType<Shop>,
        required: true
    }
})
</script>

<template>
    <GlassCard>
        <div class="shop-image-container" :class="{ 'shop-placeholder': !shop.image_url }">
            <img v-if="shop.image_url" :src="shop.image_url" :alt="shop.name" class="shop-image" />
            <span v-else class="shop-icon">🏪</span>
        </div>

        <h3 class="shop-name">{{ shop.name }}</h3>
        <p class="shop-address">{{ shop.address }}</p>

        <div v-if="shop.opening_time" class="shop-hours">
            <IconClock class="icon" />
            <span>{{ shop.opening_time }} - {{ shop.closing_time }}</span>
        </div>

        <div v-if="shop.phone" class="shop-phone">
            <IconPhone class="icon" />
            <span>{{ shop.phone }}</span>
        </div>

        <div class="flex flex-col gap-2 mt-4">
            <NavigateBtn :href="`/shop/${shop.shop_id}`" variant="primary">
                <IconBuildingStore class="w-4 h-4" />
                <span>ดูโปรไฟล์ร้านค้า</span>
            </NavigateBtn>

            <a :href="`https://www.google.com/maps/dir/?api=1&destination=${shop.latitude},${shop.longitude}`"
                target="_blank"
                class="btn btn-sm btn-outline btn-primary w-full gap-2 hover:bg-primary hover:text-white shadow-sm rounded-lg no-underline flex items-center justify-center">
                <IconNavigation class="w-4 h-4" />
                นำทางไปยังร้าน
            </a>
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
