<template>
  <div v-if="showNotification" class="fixed bottom-4 right-4 z-50">
    <div class="alert alert-warning shadow-lg max-w-sm">
      <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current flex-shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4v.01m4.5-7.5a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
      <div>
        <span class="font-bold">{{ t('Shop.notification.noShopTitle') }}</span>
        <div class="text-sm mt-1">{{ t('Shop.notification.noShopMessage') }}</div>
        <div class="mt-2 flex gap-2">
          <button @click="goToCreateShop" class="btn btn-sm btn-primary">
            {{ t('Shop.notification.createNow') }}
          </button>
          <button @click="closeNotification" class="btn btn-sm btn-ghost">
            {{ t('Shop.notification.dismiss') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useShopStore } from '@/stores/shop';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const shopStore = useShopStore();
const router = useRouter();

// Check if user has dismissed the notification before
const DISMISS_KEY = 'shop_notification_dismissed';
const dismissed = ref(localStorage.getItem(DISMISS_KEY) === 'true');

const showNotification = computed(() => {
  // Show notification if user is logged in, has checked for shop, but doesn't have one, and hasn't dismissed it
  return shopStore.checked && !shopStore.hasShop && !dismissed.value;
});

const goToCreateShop = () => {
  // Clear the dismiss flag when user goes to create shop
  localStorage.removeItem(DISMISS_KEY);
  // Include the current route as redirect parameter
  const currentRoute = router.currentRoute.value.fullPath;
  router.push({ name: 'create-shop', query: { redirect: currentRoute } });
};

const closeNotification = () => {
  // Store the dismiss preference in localStorage
  localStorage.setItem(DISMISS_KEY, 'true');
  dismissed.value = true;
};
</script>