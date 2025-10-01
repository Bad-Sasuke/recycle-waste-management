<template>
  <div v-if="showNotification" class="fixed bottom-4 right-4 z-50">
    <div class="alert alert-warning shadow-lg max-w-sm">
      <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current flex-shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4v.01m4.5-7.5a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
      <div>
        <span>You haven't created your recycling shop yet!</span>
        <div class="mt-2">
          <button @click="goToCreateShop" class="btn btn-sm btn-primary">
            Create Shop Now
          </button>
          <button @click="closeNotification" class="btn btn-sm btn-ghost ml-2">
            Dismiss
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useShopStore } from '@/stores/shop';
import { useRouter } from 'vue-router';

const shopStore = useShopStore();
const router = useRouter();

const showNotification = computed(() => {
  // Show notification if user is logged in, has checked for shop, but doesn't have one
  return shopStore.checked && !shopStore.hasShop;
});

const goToCreateShop = () => {
  router.push({ name: 'create-shop' });
};

const closeNotification = () => {
  // For now, we just hide the notification
  // In a more complex app, you might want to store this preference
};
</script>