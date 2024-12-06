<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import Card from './CardComponent.vue'
import { useWastesStore } from '../stores/wastes'
import { useCategoryWasteStore } from '../stores/category_waste'
import * as entities from '../types/recycle_waste'
import { IconFilter, IconTagPlus, IconTagMinus, IconCalendar, IconPlus } from '@tabler/icons-vue'

const searchQuery = ref('')
const selectedCategory = ref<string[]>([])
const previewImage = ref('/favicon/favicon-96x96.png')

const items = ref<entities.RecycleWaste[]>([])
const isLoading = ref(true)
const wastesStore = useWastesStore()
const categoryWasteStore = useCategoryWasteStore()
const formDataWaste = reactive({
  name: '',
  price: '',
  category: 'เลือกหมวดหมู่',
  imageFile: null,
})

onMounted(async () => {
  if (!wastesStore.wastes.length) {
    await wastesStore.fetchWastes()
  }
  if (!categoryWasteStore.category.length) {
    await categoryWasteStore.fetchCategoryWaste()
  }

  items.value = wastesStore.wastes
  isLoading.value = false
})

const filteredItems = computed(() => {
  if (items.value.length === 0) {
    return []
  }

  return items.value.filter((item) => {
    const matchesSearch = item.category.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesCategory =
      selectedCategory.value.length === 0 ||
      selectedCategory.value.includes(item.category.toLowerCase())
    return matchesSearch && matchesCategory
  })
})

const toggleCategory = (category: string) => {
  const index = selectedCategory.value.indexOf(category)
  if (index === -1) {
    selectedCategory.value.push(category)
  } else {
    selectedCategory.value.splice(index, 1)
  }
}

const sortByPrice = () => {
  items.value.sort((a, b) => a.price - b.price)
}
const sortByPriceDesc = () => {
  items.value.sort((a, b) => b.price - a.price)
}

const sortByLastUpdate = () => {
  items.value.sort((a, b) => {
    const dateA = new Date(a.lastUpdate)
    const dateB = new Date(b.lastUpdate)
    return dateB.getTime() - dateA.getTime()
  })
}

const sortByLastUpdateDesc = () => {
  items.value.sort((a, b) => {
    const dateA = new Date(a.lastUpdate)
    const dateB = new Date(b.lastUpdate)
    return dateA.getTime() - dateB.getTime()
  })
}

const openModalWaste = () => {
  const modal = document.getElementById('modal-waste') as HTMLDialogElement
  modal.showModal()
}

const loadFile = async (event) => {
  const file = event.target.files[0]
  if (file) {
    // สร้าง URL สำหรับไฟล์ที่อัพโหลด
    const objectURL = URL.createObjectURL(file)
    previewImage.value = objectURL

    // ใช้ ref เพื่ออ้างอิงไปที่ <img> และตั้ง onload ที่นั่น
    const imgElement = document.querySelector('img')
    imgElement.onload = () => {
      // รีลีส URL หลังจากโหลดเสร็จ
      URL.revokeObjectURL(objectURL) // free memory
    }

    formDataWaste.imageFile = file
  }
}

const createWaste = async () => {
  await wastesStore.addWaste(formDataWaste)
  items.value = wastesStore.wastes
  const modal = document.getElementById('modal-waste') as HTMLDialogElement
  formDataWaste.name = ''
  formDataWaste.price = ''
  formDataWaste.category = 'เลือกหมวดหมู่'
  formDataWaste.imageFile = null
  modal.close()
}
</script>

<template>
  <div class="hero min-h-screen" v-if="isLoading">
    <div class="hero-content text-center">
      <div class="max-w-md">
        <div class="flex flex-col justify-center w-full items-center h gap-4">
          <span class="loading loading-spinner loading-lg text-green-700"></span>
          <h1 class="text-xl font-bold text-green-700">กำลังโหลดข้อมูล</h1>
        </div>
      </div>
    </div>
  </div>

  <dialog id="modal-waste" class="modal">
    <div class="modal-box">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
      </form>
      <h3 class="text-lg font-bold">เพิ่มขยะรีไซเคิล</h3>
      <form class="py-4 flex flex-col gap-2">
        <label class="input input-bordered flex items-center gap-2">
          ชื่อขยะ
          <input
            type="text"
            class="grow"
            placeholder="เช่น ขวดพลาสติก"
            v-model="formDataWaste.name"
          />
        </label>
        <label class="input input-bordered flex items-center gap-2">
          ราคา
          <input type="number" class="grow" placeholder="เช่น 100" v-model="formDataWaste.price" />
        </label>

        <label class="flex items-center justify-between w-full gap-2">
          <p class="px-2">หมวดหมู่</p>
        </label>
        <select class="select select-bordered w-full" v-model="formDataWaste.category">
          <option disabled selected>เลือกหมวดหมู่</option>
          <option
            v-for="category in categoryWasteStore.category"
            :key="category.id"
            :value="category.name"
          >
            {{ category.name }}
          </option>
        </select>
        <label class="flex items-center justify-between w-full gap-2">
          <p class="px-2">รูปภาพ</p>
        </label>
        <div class="flex justify-center items-center border-2 border-dashed border-gray-300/50 p-2">
          <label class="block">
            <input
              type="file"
              @change="loadFile"
              class="block w-full text-sm text-slate-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-violet-50 file:text-violet-700 hover:file:bg-violet-100"
              accept="image/*"
              required
            />
          </label>
          <img :src="previewImage" class="h-32 rounded" alt="Current profile photo" />
        </div>
      </form>
      <button
        class="mt-4 btn bg-green-700 hover:bg-green-600 text-white w-full"
        @click="createWaste"
      >
        เพิ่มขยะ
      </button>
    </div>
  </dialog>

  <div
    class="container px-4 py-6 grid grid-cols-1 lg:grid-cols-6 gap-4 max-w-full"
    v-if="!isLoading"
  >
    <!-- ส่วนซ้าย : หมวดหมู่ -->
    <div class="bg-white shadow-md p-4 rounded lg:h-screen">
      <h2 class="text-2xl font-bold text-green-700 mb-4">หมวดหมู่</h2>
      <div class="overflow-y-auto max-h-[250px] md:max-h-[600px]">
        <label
          class="flex items-center space-x-2 mb-2"
          v-for="category in categoryWasteStore.category"
          :key="category.id"
        >
          <input type="checkbox" :value="category.name" @change="toggleCategory(category.name)" />
          <span>{{ category.name }}</span>
        </label>
      </div>
    </div>
    <!-- ส่วนขวา: รายการขยะ -->

    <div class="lg:col-span-5">
      <h1 class="text-3xl font-bold text-green-700 mb-4 text-center">ราคาขยะรีไซเคิล</h1>

      <div class="flex justify-center items-center gap-2 mb-2">
        <!-- ฟอร์มค้นหาขยะ -->
        <div class="w-full">
          <input
            type="text"
            v-model="searchQuery"
            class="input input-bordered w-full"
            placeholder="ค้นหาขยะ เช่น พลาสติก, โลหะ"
          />
        </div>
        <div>
          <button class="btn bg-green-700 hover:bg-green-600 text-white" @click="openModalWaste">
            <IconPlus stroke="2" />
          </button>
        </div>
      </div>
      <div class="flex gap-2">
        <div class="dropdown mb-2">
          <div tabindex="0" role="button" class="btn btn-outline btn-sm">
            <IconCalendar stroke="2" /> วันที่อัพเดท
          </div>
          <ul
            tabindex="0"
            class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow border border-neutral/50"
          >
            <li>
              <button @click="sortByLastUpdate">ล่าสุดไปเก่า</button>
            </li>
            <li>
              <button @click="sortByLastUpdateDesc">เก่าไปล่าสุด</button>
            </li>
          </ul>
        </div>

        <div class="dropdown mb-2">
          <div tabindex="0" role="button" class="btn btn-outline btn-sm">
            <IconFilter stroke="2" /> ราคา
          </div>
          <ul
            tabindex="0"
            class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow border border-neutral/50"
          >
            <li>
              <button @click="sortByPrice"><IconTagMinus stroke="2" /> น้อยไปมาก</button>
            </li>
            <li>
              <button @click="sortByPriceDesc"><IconTagPlus stroke="2" />มากไปน้อย</button>
            </li>
          </ul>
        </div>
      </div>

      <!-- รายการราคาขยะ -->
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
        <div v-for="item in filteredItems" :key="item.category" class="flex justify-center">
          <Card
            :name="item.name"
            :price="item.price"
            :category="item.category"
            :lastUpdate="item.lastUpdate"
            :url="item.url"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
