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
  <div class="flex card w-full md:w-64 border-green-700  rounded-xl h-96 flex-col gap-4" v-if="isDelete">
    <div class="skeleton h-60 w-full object-cover bg-[#f2f2f2] rounded-xl"></div>
    <div class="mx-4 skeleton h-4 w-28"></div>
    <div class="mx-4 skeleton h-4"></div>
    <div class="mx-4 skeleton h-4 w-16"></div>
  </div>
  <div class="card w-full md:w-64 border-green-700 rounded-xl border-[2px] cursor-pointer h-96" v-else
    @click="$emit('click', { id: props.id, name: props.name, price: props.price, category: props.category, last_update: props.last_update, url: props.url })">
    <figure>
      <img class="h-64 w-full object-cover bg-[#f2f2f2] cursor-pointer" :src="props.url" />
    </figure>
    <div class="flex items-end justify-end absolute top-0 right-0">
      <!-- Edit button for authorized users -->
      <button v-if="canEdit" class="btn btn-sm btn-ghost mr-1 mt-2 hover:bg-yellow-500 hover:text-white"
        @click.stop="openEditModal">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
        </svg>
      </button>
      <button class="btn btn-sm btn-ghost mr-2 mt-2 hover:bg-red-500 hover:text-white"
        @click.stop="confirmDelete(props?.id ?? '')">
        X
      </button>
    </div>
    <div class="card-body p-4" :data-id="props.id">
      <div class="flex flex-col items-between justify-between h-full gap-3">
        <h2 class="card-title text-base-content text-md cursor-pointer">
          {{ props.name }}
          <!-- <div class="badge badge-secondary">NEW</div> -->
        </h2>
        <div class="flex flex-col items-between gap-1 w-full">
          <h2 class="text-md cursor-pointer">{{ props.price }} บาท/กก.</h2>
          <div
            class="badge badge-outline border-green-700 text-green-700 hover:bg-green-700 hover:text-white badge-sm px-2 cursor-pointer">
            {{ props.category }}
          </div>
        </div>
      </div>
    </div>
    <div class="bg-green-700 w-full text-center text-white rounded-b-lg text-xs p-1 cursor-pointer">
      อัปเดตล่าสุด: {{ formatDateThai(props.last_update) }}
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
