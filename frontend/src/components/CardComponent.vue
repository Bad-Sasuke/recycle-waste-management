<script setup lang="ts">
import { useWastesStore } from '../stores/wastes'
import { useUsersStore } from '../stores/users'
import { ref, computed } from 'vue'

const wastesStore = useWastesStore()
const usersStore = useUsersStore()
const isDelete = ref(false)
const showDeleteConfirm = ref(false)
const itemToDelete = ref('')
const props = defineProps({
  id: { type: String },
  name: { type: String },
  price: { type: Number },
  category: { type: String },
  last_update: { type: String },
  url: { type: String },
})

// Define emits
defineEmits(['click'])

const confirmDelete = (id: string) => {
  itemToDelete.value = id
  showDeleteConfirm.value = true
}

const deleteItem = async (id: string) => {
  showDeleteConfirm.value = false // Close the modal
  isDelete.value = true
  await new Promise((resolve) => setTimeout(resolve, 1000))
  await wastesStore.deleteWaste(id)
  isDelete.value = false
  itemToDelete.value = '' // Reset the item to delete
}

const cancelDelete = () => {
  showDeleteConfirm.value = false
  itemToDelete.value = ''
}

// Check if current user can edit (admin or moderator)
const canEdit = computed(() => {
  return usersStore.isLogin && (usersStore.user?.role === 'admin' || usersStore.user?.role === 'moderator')
})

// Function to open the edit modal
const openEditModal = () => {
  // Set the waste to edit in the store
  wastesStore.setWasteToEdit(props.id ?? '', {
    name: props.name ?? '',
    price: props.price ?? 0,
    category: props.category ?? '',
    url: props.url ?? '',
    last_update: props.last_update ?? ''
  });

  // Show the modal
  const modal = document.getElementById('modal-edit-waste') as HTMLDialogElement
  if (modal) {
    modal.showModal()
  }
}

// Format date to Thai format (e.g., 1 มกราคม 2567 เวลา 10:30 น.)
const formatDateThai = (dateString: string | undefined): string => {
  if (!dateString) return '';

  try {
    const date = new Date(dateString);

    // Check if the date is valid
    if (isNaN(date.getTime())) {
      return dateString; // Return original if parsing fails
    }

    const day = date.getDate();
    const months = [
      'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน',
      'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม'
    ];
    const month = months[date.getMonth()];
    const year = date.getFullYear() + 543; // Convert to Buddhist Era

    // Format time
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');

    return `${day} ${month} ${year} เวลา ${hours}:${minutes} น.`;
  } catch (error) {
    console.error('Error formatting date:', error);
    return dateString;
  }
}
</script>

<template>
  <div class="flex card w-full bg-base-100 shadow-md rounded-2xl h-full flex-col gap-4 p-4 animate-pulse"
    v-if="isDelete">
    <div class="skeleton h-48 w-full rounded-xl"></div>
    <div class="skeleton h-4 w-28"></div>
    <div class="skeleton h-4 w-full"></div>
    <div class="skeleton h-4 w-16"></div>
  </div>

  <div
    class="group card w-full bg-white hover:shadow-2xl transition-all duration-300 rounded-2xl overflow-hidden border border-gray-100 h-full flex flex-col hover:-translate-y-1"
    v-else
    @click="$emit('click', { id: props.id, name: props.name, price: props.price, category: props.category, last_update: props.last_update, url: props.url })">

    <figure class="relative h-56 overflow-hidden">
      <img class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
        :src="props.url || 'https://placehold.co/400x300?text=No+Image'" :alt="props.name" />

      <!-- Badges overlay -->
      <div class="absolute top-3 left-3 flex flex-col gap-2">
        <div v-if="props.category"
          class="badge border-none text-white font-medium shadow-sm backdrop-blur-md bg-green-600/90 text-xs py-3 px-3">
          {{ props.category }}
        </div>
      </div>

      <!-- Action buttons -->
      <div class="absolute top-3 right-3 flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
        <button v-if="canEdit"
          class="btn btn-circle btn-sm bg-white/90 border-none text-amber-500 hover:bg-amber-500 hover:text-white shadow-md backdrop-blur-sm"
          @click.stop="openEditModal" title="แก้ไข">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
        <button
          class="btn btn-circle btn-sm bg-white/90 border-none text-red-500 hover:bg-red-500 hover:text-white shadow-md backdrop-blur-sm"
          @click.stop="confirmDelete(props?.id ?? '')" title="ลบ">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Soft gradient overlay at bottom -->
      <div class="absolute inset-x-0 bottom-0 h-16 bg-gradient-to-t from-black/50 to-transparent opacity-60"></div>
    </figure>

    <div class="card-body p-5 flex-1 flex flex-col justify-between relative bg-white">
      <div>
        <h2 class="font-bold text-lg text-gray-800 line-clamp-2 mb-2 group-hover:text-green-700 transition-colors">
          {{ props.name }}
        </h2>

        <div class="flex items-baseline gap-1 mt-1">
          <span class="text-2xl font-bold text-green-600">{{ props.price }}</span>
          <span class="text-gray-500 text-sm">บาท/กก.</span>
        </div>
      </div>

      <div class="border-t border-gray-100 mt-4 pt-3 flex items-center justify-between text-xs text-gray-400">
        <div class="flex items-center gap-1">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <polyline points="12 6 12 12 16 14"></polyline>
          </svg>
          <span>{{ formatDateThai(props.last_update).split('เวลา')[0] }}</span>
        </div>
      </div>
    </div>
  </div>

  <!-- Delete Confirmation Modal -->
  <dialog id="delete_confirm_modal" class="modal" :class="{ 'modal-open': showDeleteConfirm }">
    <div class="modal-box">
      <h3 class="font-bold text-lg ">ยืนยันการลบ</h3>
      <p class="py-4">คุณแน่ใจใช่ไหมว่าต้องการลบ "{{ props.name }}"?</p>
      <div class="modal-action">
        <button class="btn btn-outline" @click="cancelDelete">
          ยกเลิก
        </button>
        <button class="btn btn-error text-white" @click="deleteItem(itemToDelete)">
          ยืนยันการลบ
        </button>
      </div>
    </div>
    <form class="modal-backdrop" @click="cancelDelete">
      <button type="button">close</button>
    </form>
  </dialog>
</template>
