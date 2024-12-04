<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import Card from './CardComponent.vue'
import { useWastesStore } from '../stores/wastes'
import * as entities from '../types/recycle_waste'

const searchQuery = ref('')
const selectedCategory = ref<string[]>([])

const items = ref<entities.RecycleWaste[]>([])
const wastesStore = useWastesStore()
const isLoading = ref(true)

onMounted(async () => {
  if (!wastesStore.wastes.length) {
    await wastesStore.fetchWastes()
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

  <div
    class="container px-4 py-6 grid grid-cols-1 lg:grid-cols-6 gap-4 max-w-full"
    v-if="!isLoading"
  >
    <!-- ส่วนซ้าย : หมวดหมู่ -->
    <div class="bg-white shadow-md p-4 rounded lg:h-screen">
      <h2 class="text-2xl font-bold text-green-700 mb-4">หมวดหมู่</h2>
      <div>
        <label class="flex items-center space-x-2 mb-2">
          <input type="checkbox" value="พลาสติก" @change="toggleCategory('พลาสติก')" />
          <span>พลาสติก</span>
        </label>
        <label class="flex items-center space-x-2 mb-2">
          <input type="checkbox" value="โลหะ" @change="toggleCategory('โลหะ')" />
          <span>โลหะ</span>
        </label>
        <label class="flex items-center space-x-2 mb-2">
          <input type="checkbox" value="กระดาษ" @change="toggleCategory('กระดาษ')" />
          <span>กระดาษ</span>
        </label>
        <label class="flex items-center space-x-2 mb-2">
          <input type="checkbox" value="แก้ว" @change="toggleCategory('แก้ว')" />
          <span>แก้ว</span>
        </label>
      </div>
    </div>
    <!-- ส่วนขวา: รายการขยะ -->

    <div class="lg:col-span-5">
      <h1 class="text-3xl font-bold text-green-700 mb-4 text-center">ราคาขยะรีไซเคิล</h1>

      <div class="flex justify-center items-center gap-2 mb-6">
        <!-- ฟอร์มค้นหาขยะ -->
        <div class="w-full">
          <input
            type="text"
            v-model="searchQuery"
            class="input input-bordered w-full"
            placeholder="ค้นหาขยะ เช่น พลาสติก, โลหะ"
          />
        </div>
      </div>

      <!-- รายการราคาขยะ -->
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
        <div v-for="item in filteredItems" :key="item.category" class="flex justify-center">
          <Card
            :name="item.category"
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
