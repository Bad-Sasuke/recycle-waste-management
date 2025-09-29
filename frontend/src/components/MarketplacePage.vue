<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useUsersStore } from '../stores/users'
import Card from './CardComponent.vue'
import PopupWaste from './PopupWaste.vue'
import EditWasteModal from './EditWasteModal.vue'
import { useWastesStore } from '../stores/wastes'
import { useCategoryWasteStore } from '../stores/category_waste'
import type RecycleWaste from '../types/recycle_waste'
import { IconFilter, IconTagPlus, IconTagMinus, IconCalendar, IconPlus } from '@tabler/icons-vue'


const searchQuery = ref('')
const selectedCategory = ref<string[]>([])

const items = ref<RecycleWaste[]>([])
const isLoading = ref(true)
const usersStore = useUsersStore()
const wastesStore = useWastesStore()
const categoryWasteStore = useCategoryWasteStore()
const currentPage = ref(1)

onMounted(async () => {
  if (!wastesStore.wastes.length) {
    await wastesStore.fetchWastes(currentPage.value)
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

// Watch for pagination changes
watch(
  () => [wastesStore.pagination.page, wastesStore.pagination.total_pages],
  () => {
    currentPage.value = wastesStore.pagination.page
  }
)

const goToPage = async (page: number) => {
  if (page >= 1 && page <= wastesStore.pagination.total_pages) {
    currentPage.value = page
    isLoading.value = true
    await wastesStore.fetchWastes(page)
    items.value = wastesStore.wastes
    isLoading.value = false
  }
}

const goToNextPage = async () => {
  if (currentPage.value < wastesStore.pagination.total_pages) {
    await goToPage(currentPage.value + 1)
  }
}

const goToPrevPage = async () => {
  if (currentPage.value > 1) {
    await goToPage(currentPage.value - 1)
  }
}

const getVisiblePages = () => {
  const totalPages = wastesStore.pagination.total_pages;
  const maxVisiblePages = 5;
  const halfMax = Math.floor(maxVisiblePages / 2);

  let startPage = Math.max(currentPage.value - halfMax, 1);
  const endPage = Math.min(startPage + maxVisiblePages - 1, totalPages);

  if (endPage - startPage + 1 < maxVisiblePages) {
    startPage = Math.max(endPage - maxVisiblePages + 1, 1);
  }

  const pages = [];
  for (let i = startPage; i <= endPage; i++) {
    pages.push(i);
  }

  return pages;
}

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

// Update filter behavior to reset pagination
watch(
  [searchQuery, selectedCategory],
  async () => {
    currentPage.value = 1
    isLoading.value = true
    await wastesStore.fetchWastes(currentPage.value)
    items.value = wastesStore.wastes
    isLoading.value = false
  },
  { deep: true }
)

const toggleCategory = (category: string) => {
  const categoryLower = category.toLowerCase();
  const index = selectedCategory.value.indexOf(categoryLower)
  if (index === -1) {
    selectedCategory.value.push(categoryLower)
  } else {
    selectedCategory.value.splice(index, 1)
  }
}

// Client-side sorting will not work with pagination, so we need to fetch sorted data from server
// For now, keeping the original sorting functions but they will only sort the current page
const sortByPrice = async () => {
  // For now, just sort the current page
  items.value.sort((a, b) => (a?.price ?? 0) - (b?.price ?? 0))
}

const sortByPriceDesc = async () => {
  // For now, just sort the current page
  items.value.sort((a, b) => (b?.price ?? 0) - (a?.price ?? 0))
}

const sortByLastUpdate = async () => {
  // For now, just sort the current page
  items.value.sort((a, b) => {
    const dateA = new Date(a?.lastUpdate ?? 0)
    const dateB = new Date(b?.lastUpdate ?? 0)
    return dateB.getTime() - dateA.getTime()
  })
}

const sortByLastUpdateDesc = async () => {
  // For now, just sort the current page
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

const canAddProduct = computed(() => {
  return usersStore.isLogin && (usersStore.user?.role === 'admin' || usersStore.user?.role === 'moderator')
})


</script>

<template>
  <div class="hero min-h-screen" v-if="isLoading">
    <div class="hero-content text-center">
      <div class="max-w-md">
        <div class="flex flex-col justify-center w-full items-center h gap-4">
          <span class="loading loading-spinner loading-lg text-green-700"></span>
          <h1 class="text-xl font-bold text-green-700">{{ $t('Global.loading') }}</h1>
        </div>
      </div>
    </div>
  </div>

  <PopupWaste />
  <EditWasteModal />

  <div class="container px-4 py-6 grid grid-cols-1 lg:grid-cols-6 gap-4 max-w-full" v-if="!isLoading">
    <!-- ส่วนซ้าย : หมวดหมู่ -->
    <div class="bg-white shadow-md p-4 rounded lg:h-screen">
      <h2 class="text-2xl font-bold text-green-700 mb-4">{{ $t('Marketplace.category') }}</h2>
      <div class="overflow-y-auto max-h-[250px] md:max-h-[600px]">
        <label class="flex items-center space-x-2 mb-2" v-for="category in categoryWasteStore.category"
          :key="category.id">
          <input type="checkbox" 
            :value="category.name" 
            :checked="selectedCategory.includes((category.name ?? 'ไม่มี').toLowerCase())"
            @change="toggleCategory(category?.name ?? 'ไม่มี')" />
          <span>{{ category?.name ?? 'ไม่มี' }}</span>
        </label>
      </div>
    </div>
    <!-- ส่วนขวา: รายการขยะ -->

    <div class="lg:col-span-5">
      <h1 class="text-3xl font-bold text-green-700 mb-4 text-center">
        {{ $t('Marketplace.title') }}
      </h1>

      <div class="flex justify-center items-center gap-2 mb-2">
        <!-- ฟอร์มค้นหาขยะ -->
        <div class="w-full">
          <input type="text" v-model="searchQuery" class="input input-bordered w-full"
            :placeholder="$t('Marketplace.search')" />
        </div>
        <div v-if="canAddProduct">
          <button class="btn bg-green-700 hover:bg-green-600 text-white" @click="openModalWaste">
            <IconPlus stroke="2" />
          </button>
        </div>
      </div>
      <div class="flex gap-2">
        <div class="dropdown mb-2">
          <div tabindex="0" role="button" class="btn btn-outline btn-sm">
            <IconCalendar stroke="2" /> {{ $t('Marketplace.filter_date.title') }}
          </div>
          <ul tabindex="0"
            class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow border border-neutral/50">
            <li>
              <button @click="sortByLastUpdate">{{ $t('Marketplace.filter_date.lastest') }}</button>
            </li>
            <li>
              <button @click="sortByLastUpdateDesc">
                {{ $t('Marketplace.filter_date.oldest') }}
              </button>
            </li>
          </ul>
        </div>

        <div class="dropdown mb-2">
          <div tabindex="0" role="button" class="btn btn-outline btn-sm">
            <IconFilter stroke="2" /> {{ $t('Marketplace.filter_price.title') }}
          </div>
          <ul tabindex="0"
            class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow border border-neutral/50">
            <li>
              <button @click="sortByPrice">
                <IconTagMinus stroke="2" /> {{ $t('Marketplace.filter_price.asc') }}
              </button>
            </li>
            <li>
              <button @click="sortByPriceDesc">
                <IconTagPlus stroke="2" />{{ $t('Marketplace.filter_price.desc') }}
              </button>
            </li>
          </ul>
        </div>
      </div>

      <!-- รายการราคาขยะ -->
      <div class="mb-4 min-h-[800px]">
        <TransitionGroup name="list" tag="div"
          class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-6 xl:grid-cols-7 gap-6">
          <div v-for="item in filteredItems" :key="item.waste_id" class="flex justify-center">
            <Card :id="item.waste_id" :name="item.name" :price="item.price" :category="item.category"
              :lastUpdate="item.lastUpdate" :url="item.url" />
          </div>
        </TransitionGroup>
      </div>

      <!-- Pagination -->
      <div class="flex flex-col items-center mt-6">
        <div class="join">
          <button class="join-item btn" :disabled="currentPage <= 1" @click="goToPage(1)"
            :title="$t('Marketplace.pagination.first')">
            &laquo;&laquo;
          </button>
          <button class="join-item btn" :disabled="currentPage <= 1" @click="goToPrevPage"
            :title="$t('Marketplace.pagination.prev')">
            «
          </button>
          <button v-for="page in getVisiblePages()" :key="page" class="join-item btn"
            :class="{ 'btn-active': page === currentPage }" @click="goToPage(page)">
            {{ page }}
          </button>
          <button class="join-item btn" :disabled="currentPage >= wastesStore.pagination.total_pages"
            @click="goToNextPage" :title="$t('Marketplace.pagination.next')">
            »
          </button>
          <button class="join-item btn" :disabled="currentPage >= wastesStore.pagination.total_pages"
            @click="goToPage(wastesStore.pagination.total_pages)" :title="$t('Marketplace.pagination.last')">
            &raquo;&raquo;
          </button>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

.list-leave-active {
  display: none;
}
</style>
