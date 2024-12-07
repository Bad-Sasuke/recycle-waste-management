<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import Card from './CardComponent.vue'
import ModelWaste from './ModelWaste.vue'
import { useWastesStore } from '../stores/wastes'
import { useCategoryWasteStore } from '../stores/category_waste'
import type RecycleWaste from '../types/recycle_waste'
import { IconFilter, IconTagPlus, IconTagMinus, IconCalendar, IconPlus } from '@tabler/icons-vue'

const searchQuery = ref('')
const selectedCategory = ref<string[]>([])

const items = ref<RecycleWaste[]>([])
const isLoading = ref(true)
const wastesStore = useWastesStore()
const categoryWasteStore = useCategoryWasteStore()

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

watch(
  () => wastesStore.wastes,
  (newValue) => {
    items.value = newValue
  },
)

const filteredItems = computed(() => {
  if (items.value.length === 0) {
    return []
  }

  return items.value.filter((item) => {
    const matchesSearch = item.name
      ? item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
      : false
    const matchesCategory = item.category
      ? selectedCategory.value.length === 0 ||
        selectedCategory.value.includes(item.category?.toLowerCase())
      : true
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
  items.value.sort((a, b) => (a?.price ?? 0) - (b?.price ?? 0))
}
const sortByPriceDesc = () => {
  items.value.sort((a, b) => (b?.price ?? 0) - (a?.price ?? 0))
}

const sortByLastUpdate = () => {
  items.value.sort((a, b) => {
    const dateA = new Date(a?.lastUpdate ?? 0)
    const dateB = new Date(b?.lastUpdate ?? 0)
    return dateB.getTime() - dateA.getTime()
  })
}

const sortByLastUpdateDesc = () => {
  items.value.sort((a, b) => {
    const dateA = new Date(a.lastUpdate ?? 0)
    const dateB = new Date(b.lastUpdate ?? 0)
    return dateA.getTime() - dateB.getTime()
  })
}

const openModalWaste = () => {
  const modal = document.getElementById('modal-waste') as HTMLDialogElement
  modal.showModal()
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

  <ModelWaste />

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
          <input
            type="checkbox"
            :value="category.name"
            @change="toggleCategory(category?.name ?? 'ไม่มี')"
          />
          <span>{{ category?.name ?? 'ไม่มี' }}</span>
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
      <div
        class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-6 xl:grid-cols-7 gap-6"
      >
        <div v-for="(item, index) in filteredItems" :key="index" class="flex justify-center">
          <Card
            :id="item.waste_id"
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
