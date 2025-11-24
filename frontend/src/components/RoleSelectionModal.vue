<script setup lang="ts">
import { ref, defineExpose } from 'vue';
import { IconUser, IconBuildingStore, IconCheck } from '@tabler/icons-vue';
import { useUsersStore } from '@/stores/users';
import { useRouter } from 'vue-router';

const usersStore = useUsersStore();
const router = useRouter();
const modal = ref<HTMLDialogElement | null>(null);
const selectedRole = ref<'user' | 'moderator' | null>(null);
const isLoading = ref(false);

const showModal = () => {
  if (modal.value) {
    modal.value.showModal();
  }
};

const closeModal = () => {
  if (modal.value) {
    modal.value.close();
  }
};

const selectRole = (role: 'user' | 'moderator') => {
  selectedRole.value = role;
};

const confirmRole = async () => {
  if (!selectedRole.value) return;

  isLoading.value = true;
  const success = await usersStore.updateUserRole(selectedRole.value);
  isLoading.value = false;

  if (success) {
    closeModal();
    if (selectedRole.value === 'moderator') {
      router.push({ name: 'create-shop' });
    }
  }
};

defineExpose({
  showModal,
  closeModal
});
</script>

<template>
  <dialog ref="modal" class="modal">
    <div class="modal-box max-w-3xl p-8">
      <h3 class="font-bold text-3xl text-center mb-2 text-gray-800">เลือกประเภทผู้ใช้งาน</h3>
      <p class="text-center text-gray-500 mb-8">กรุณาเลือกประเภทบัญชีของคุณเพื่อเริ่มต้นใช้งาน</p>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <!-- Customer Card -->
        <div @click="selectRole('user')"
          class="relative p-8 rounded-2xl border-2 cursor-pointer transition-all duration-300 hover:shadow-xl group"
          :class="selectedRole === 'user' ? 'border-green-500 bg-green-50' : 'border-gray-200 hover:border-green-300'">
          <div v-if="selectedRole === 'user'" class="absolute top-4 right-4 text-green-500">
            <IconCheck size="24" stroke="3" />
          </div>
          <div class="flex flex-col items-center gap-4">
            <div
              class="p-4 rounded-full bg-green-100 text-primary group-hover:scale-110 transition-transform duration-300">
              <IconUser size="48" stroke="1.5" />
            </div>
            <h4 class="text-xl font-bold text-gray-800">ลูกค้าทั่วไป</h4>
            <p class="text-center text-gray-500 text-sm">
              สำหรับผู้ที่ต้องการค้นหาร้านรับซื้อขยะรีไซเคิล และขายขยะ
            </p>
          </div>
        </div>

        <!-- Shop Owner Card -->
        <div @click="selectRole('moderator')"
          class="relative p-8 rounded-2xl border-2 cursor-pointer transition-all duration-300 hover:shadow-xl group"
          :class="selectedRole === 'moderator' ? 'border-green-500 bg-green-50' : 'border-gray-200 hover:border-green-300'">
          <div v-if="selectedRole === 'moderator'" class="absolute top-4 right-4 text-green-500">
            <IconCheck size="24" stroke="3" />
          </div>
          <div class="flex flex-col items-center gap-4">
            <div
              class="p-4 rounded-full bg-green-100 text-primary group-hover:scale-110 transition-transform duration-300">
              <IconBuildingStore size="48" stroke="1.5" />
            </div>
            <h4 class="text-xl font-bold text-gray-800">ร้านรับซื้อขยะ</h4>
            <p class="text-center text-gray-500 text-sm">
              สำหรับร้านค้าที่ต้องการรับซื้อขยะ และจัดการร้านค้าของคุณ
            </p>
          </div>
        </div>
      </div>

      <div class="flex justify-center gap-4">
        <button @click="confirmRole"
          class="btn btn-primary btn-lg text-white min-w-[200px] rounded-full shadow-lg hover:shadow-xl hover:scale-105 transition-all"
          :disabled="!selectedRole || isLoading">
          <span v-if="isLoading" class="loading loading-spinner"></span>
          {{ isLoading ? 'กำลังบันทึก...' : 'ยืนยันการเลือก' }}
        </button>
      </div>
    </div>
    <!-- Prevent closing by clicking outside (backdrop) to force selection -->
    <form method="dialog" class="modal-backdrop bg-black/50">
      <button class="cursor-default">close</button>
    </form>
  </dialog>
</template>
