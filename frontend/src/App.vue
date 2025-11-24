<script setup lang="ts">
import { RouterView } from 'vue-router'
import { ref, watch } from 'vue'
import { useUsersStore } from '@/stores/users'
import FooterComponent from './components/FooterComponent.vue'
import NavbarHeader from './components/NavbarHeader.vue'
import ShopNotification from './components/ShopNotification.vue'
import RoleSelectionModal from './components/RoleSelectionModal.vue'

const usersStore = useUsersStore()
const roleModal = ref<InstanceType<typeof RoleSelectionModal> | null>(null)

watch(
  () => usersStore.user,
  (newUser) => {
    if (usersStore.isLogin && newUser && !newUser.role) {
      // Small delay to ensure modal is mounted and ready
      setTimeout(() => {
        roleModal.value?.showModal()
      }, 500)
    }
  },
  { immediate: true }
)
</script>

<template>
  <div class="font-kanit flex flex-col min-h-screen">
    <header>
      <NavbarHeader />
    </header>
    <RouterView />
    <FooterComponent />
    <ShopNotification />
    <RoleSelectionModal ref="roleModal" />
  </div>
</template>
