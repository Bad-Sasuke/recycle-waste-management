<template>
  <dialog id="modal-edit-waste" class="modal">
    <div class="modal-box">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
      </form>
      <h3 class="text-lg font-bold">แก้ไขขยะรีไซเคิล</h3>
      <form class="py-4 flex flex-col gap-2">
        <label class="input input-bordered flex items-center gap-2">
          ชื่อขยะ
          <input type="text" class="grow" placeholder="เช่น ขวดพลาสติก" v-model="formDataWaste.name"
            :disabled="isUpdatingWaste" />
        </label>
        <ErrorText :message="errName" />
        <label class="input input-bordered flex items-center gap-2">
          ราคา
          <input type="number" class="grow" placeholder="เช่น 100" v-model="formDataWaste.price"
            :disabled="isUpdatingWaste" />
        </label>
        <ErrorText :message="errPrice" />

        <label class="flex items-center justify-between w-full gap-2">
          <p class="px-2">หมวดหมู่</p>
        </label>
        <select class="select select-bordered w-full" v-model="formDataWaste.category" :disabled="isUpdatingWaste">
          <option disabled selected>เลือกหมวดหมู่</option>
          <option v-for="category in categoryWasteStore.category" :key="category.id" :value="category.name">
            {{ category.name }}
          </option>
        </select>
        <ErrorText :message="errCategory" />
        <label class="flex items-center justify-between w-full gap-2">
          <p class="px-2">รูปภาพ</p>
        </label>
        <div class="flex justify-center items-center border-2 border-dashed border-gray-300/50 p-2">
          <label class="block">
            <input type="file" @change="loadFile"
              class="block w-full text-sm text-slate-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-violet-50 file:text-violet-700 hover:file:bg-violet-100"
              accept="image/*" :class="{
                'cursor-not-allowed': isUpdatingWaste,
                'file:text-gray-700 hover:file:bg-gray-100': isUpdatingWaste,
              }" :disabled="isUpdatingWaste" />
          </label>
          <img :src="previewImage" class="h-32 rounded" alt="Current profile photo" />
        </div>
      </form>
      <button class="mt-4 btn bg-green-700 hover:bg-green-600 text-white w-full" @click="updateWaste"
        :disabled="isUpdatingWaste">
        <span class="loading loading-spinner loading-sm" v-if="isUpdatingWaste"></span>
        <p v-if="!isUpdatingWaste">อัปเดตขยะ</p>
      </button>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useCategoryWasteStore } from '../stores/category_waste'
import { useWastesStore } from '../stores/wastes'
import ErrorText from './ErrorText.vue'

interface FormDataWaste {
  name: string
  price: string
  category: string
  imageFile: File | null
}

const previewImage = ref('/favicon/favicon-96x96.png')
const errName = ref('')
const errPrice = ref('')
const errCategory = ref('')
const isUpdatingWaste = ref(false)
const categoryWasteStore = useCategoryWasteStore()
const wastesStore = useWastesStore()

const formDataWaste = reactive<FormDataWaste>({
  name: '',
  price: '',
  category: 'เลือกหมวดหมู่',
  imageFile: null,
})

// Watch for changes in wasteToEdit from the store
watch(() => wastesStore.wasteToEdit, (newWasteToEdit) => {
  if (newWasteToEdit) {
    // Initialize the form with the waste data
    formDataWaste.name = newWasteToEdit.name
    formDataWaste.price = newWasteToEdit.price.toString()
    formDataWaste.category = newWasteToEdit.category
    previewImage.value = newWasteToEdit.url
  }
}, { immediate: true })

// Function to load file for preview
const loadFile = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    // สร้าง URL สำหรับไฟล์ที่อัพโหลด
    const objectURL = URL.createObjectURL(file)
    previewImage.value = objectURL

    // ใช้ ref เพื่ออ้างอิงไปที่ <img> และตั้ง onload ที่นั่น
    const imgElement = document.getElementById('img') as HTMLImageElement | null
    if (imgElement) {
      imgElement.onload = () => {
        URL.revokeObjectURL(objectURL) // free memory
      }
    }

    formDataWaste.imageFile = file as File | null
  }
}

// Function to validate form data
const validateForm = (): boolean => {
  errName.value = ''
  errPrice.value = ''
  errCategory.value = ''

  if (!formDataWaste.name) {
    errName.value = 'กรุณากรอกชื่อขยะรีไซเคิล'
    return false
  }

  if (!formDataWaste.price) {
    errPrice.value = 'กรุณากรอกราคาขยะรีไซเคิล'
    return false
  } else if (Number(formDataWaste.price) < 0) {
    errPrice.value = 'กรุณากรอกราคาขยะรีไซเคิลให้ถูกต้อง'
    return false
  }

  if (!formDataWaste.category || formDataWaste.category === 'เลือกหมวดหมู่') {
    errCategory.value = 'กรุณาเลือกหมวดหมู่ขยะรีไซเคิล'
    return false
  }

  return true
}

// Function to update the waste item
const updateWaste = async () => {
  if (!validateForm()) {
    return
  }

  if (!wastesStore.wasteToEdit) {
    console.error('No waste ID provided for update')
    return
  }

  isUpdatingWaste.value = true

  // Prepare form data for API request
  const updateData = {
    name: formDataWaste.name,
    price: Number(formDataWaste.price),
    category: formDataWaste.category,
    imageFile: formDataWaste.imageFile
  }

  try {
    // Call the update function in the wastes store
    await wastesStore.updateWaste(wastesStore.wasteToEdit.id, updateData)

    // Reset form and clear waste to edit
    resetForm()

    // Close modal
    const modal = document.getElementById('modal-edit-waste') as HTMLDialogElement
    if (modal) {
      modal.close()
    }
  } catch (error) {
    console.error('Error updating waste:', error)
  } finally {
    isUpdatingWaste.value = false
  }
}

// Reset form data
const resetForm = () => {
  formDataWaste.name = ''
  formDataWaste.price = ''
  formDataWaste.category = 'เลือกหมวดหมู่'
  formDataWaste.imageFile = null
  errName.value = ''
  errPrice.value = ''
  errCategory.value = ''
  previewImage.value = '/favicon/favicon-96x96.png'
  wastesStore.clearWasteToEdit()
}
</script>
