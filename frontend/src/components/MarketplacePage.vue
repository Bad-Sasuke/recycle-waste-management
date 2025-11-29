<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { debounce } from 'lodash-es'
import { useUsersStore } from '../stores/users'
import Card from './CardComponent.vue'
import PopupWaste from './PopupWaste.vue'
import EditWasteModal from './EditWasteModal.vue'
import ProductDetailModal from './ProductDetailModal.vue'
import { useWastesStore } from '../stores/wastes'
import { useCategoryWasteStore } from '../stores/category_waste'
import type { GroupedRecyclableItem } from '../types/recycle_waste'
import { IconTagPlus, IconTagMinus, IconCalendar, IconPlus, IconSearch, IconCategory, IconSortAscending, IconSortDescending, IconArrowsSort, IconChevronLeft, IconChevronRight, IconChevronsLeft, IconChevronsRight } from '@tabler/icons-vue'

const searchQuery = ref('')
const selectedCategory = ref<string[]>([])


const isLoading = ref(true)
const usersStore = useUsersStore()
const wastesStore = useWastesStore()
const categoryWasteStore = useCategoryWasteStore()
const currentPage = ref(1)

// Ref for product detail modal
const productDetailModal = ref<InstanceType<typeof ProductDetailModal> | null>(null)

onMounted(async () => {
  if (!wastesStore.groupedWastes.length) {
    await wastesStore.fetchWastes(currentPage.value, 12) // Fetch grouped data
  }
  if (!categoryWasteStore.category.length) {
    await categoryWasteStore.fetchCategoryWaste()
  }

  // Store the original data for filtering
  originalGroupedWastes.value = [...wastesStore.groupedWastes]

  isLoading.value = false
})



// Watch for pagination changes
watch(
  () => [wastesStore.groupedPagination.page, wastesStore.groupedPagination.total_pages],
  () => {
    currentPage.value = wastesStore.groupedPagination.page
  }
)

const goToPage = async (page: number) => {
  if (page >= 1 && page <= wastesStore.groupedPagination.total_pages) {
    currentPage.value = page
    isLoading.value = true
    await wastesStore.fetchWastes(page, 12) // Fetch grouped data
    originalGroupedWastes.value = [...wastesStore.groupedWastes] // Update original data
    isLoading.value = false
  }
}

const goToNextPage = async () => {
  if (currentPage.value < wastesStore.groupedPagination.total_pages) {
    await goToPage(currentPage.value + 1)
  }
}

const goToPrevPage = async () => {
  if (currentPage.value > 1) {
    await goToPage(currentPage.value - 1)
  }
}

const getVisiblePages = () => {
  const totalPages = wastesStore.groupedPagination.total_pages;
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

// Store original data to avoid refetching when filtering
const originalGroupedWastes = ref<GroupedRecyclableItem[]>([])

// Reactive variable to track the current sort method
const currentSortMethod = ref<string | null>(null)

// Update to use grouped items instead of individual items with local filtering
const filteredItems = computed(() => {
  if (originalGroupedWastes.value.length === 0) {
    return []
  }

  const result = originalGroupedWastes.value.filter((item) => {
    const matchesSearch = item.name
      ? item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
      : false
    const matchesCategory = item.category
      ? selectedCategory.value.length === 0 ||
      selectedCategory.value.includes(item.category?.toLowerCase())
      : true
    return matchesSearch && matchesCategory
  })

  // Apply sorting based on the current sort method
  if (currentSortMethod.value) {
    switch (currentSortMethod.value) {
      case 'price-asc':
        result.sort((a, b) => (a?.price ?? 0) - (b?.price ?? 0))
        break
      case 'price-desc':
        result.sort((a, b) => (b?.price ?? 0) - (a?.price ?? 0))
        break
      case 'date-asc':
        result.sort((a, b) => {
          const dateA = new Date(a?.last_update ?? 0)
          const dateB = new Date(b?.last_update ?? 0)
          return dateA.getTime() - dateB.getTime()
        })
        break
      case 'date-desc':
        result.sort((a, b) => {
          const dateA = new Date(a?.last_update ?? 0)
          const dateB = new Date(b?.last_update ?? 0)
          return dateB.getTime() - dateA.getTime()
        })
        break
    }
  }

  return result
})

// Watch for changes in the store data and update original data accordingly
watch(
  () => wastesStore.groupedWastes,
  (newData) => {
    originalGroupedWastes.value = [...newData]
  }
)

// Debounced search function to optimize performance
const debouncedSearch = debounce(() => {
  // No need to make API calls - the computed property will automatically update
  currentPage.value = 1 // Reset to first page when filtering
}, 300) // 300ms delay

// Update filter behavior to use local filtering instead of API calls with debounce
watch(
  [searchQuery, selectedCategory],
  () => {
    debouncedSearch()
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


const sortByPrice = () => {
  currentSortMethod.value = 'price-asc'
}

const sortByPriceDesc = () => {
  currentSortMethod.value = 'price-desc'
}

const sortByLastUpdate = () => {
  currentSortMethod.value = 'date-desc' // Latest first
}

const sortByLastUpdateDesc = () => {
  currentSortMethod.value = 'date-asc' // Oldest first
}

const openModalWaste = () => {
  const modal = document.getElementById('modal-waste') as HTMLDialogElement
  modal.showModal()
}

const canAddProduct = computed(() => {
  return usersStore.isLogin && (usersStore.user?.role === 'admin' || usersStore.user?.role === 'moderator')
})

// Function to open product detail modal
const openProductDetailModal = (item: GroupedRecyclableItem) => {
  if (productDetailModal.value) {
    productDetailModal.value.openModal(item)
  }
}

</script>

<template>
  <div class="hero min-h-screen" v-if="isLoading">
    <div class="hero-content text-center">
      <div class="max-w-md">
        <div class="flex flex-col justify-center w-full items-center h gap-6">
          <div class="flex items-center gap-2">
            <span class="loading loading-spinner loading-lg text-green-700"></span>
            <h1 class="text-xl md:text-2xl font-bold text-green-700">{{ $t('Global.loading') }}</h1>
          </div>
        </div>
      </div>
    </div>
  </div>

  <PopupWaste />
  <EditWasteModal />
  <ProductDetailModal ref="productDetailModal" />

  <div class="container px-4 py-6 grid grid-cols-1 lg:grid-cols-6 gap-4 max-w-full" v-if="!isLoading">
    <!-- ส่วนซ้าย : หมวดหมู่ -->
    <div class="bg-green-50 border border-green-200 p-4 rounded lg:min-h-screen">
      <div class="flex items-center gap-2 mb-4">
        <IconCategory stroke="2" size="24" class="text-green-700" />
        <h2 class="text-xl md:text-2xl font-bold text-green-700">{{ $t('Marketplace.category') }}</h2>
      </div>
      <div class="overflow-y-auto max-h-[250px] md:max-h-[600px]">
        <label class="flex items-center gap-2 mb-2 p-2 hover:bg-green-100 rounded cursor-pointer transition-all"
          v-for="category in categoryWasteStore.category" :key="category.id">
          <input type="checkbox" :value="category.name"
            :checked="selectedCategory.includes((category.name ?? 'ไม่มี').toLowerCase())"
            @change="toggleCategory(category?.name ?? 'ไม่มี')" class="checkbox checkbox-primary checkbox-sm" />
          <span class="text-sm md:text-base">{{ category?.name ?? 'ไม่มี' }}</span>
        </label>
      </div>
    </div>
    <!-- ส่วนขวา: รายการขยะ -->

    <div class="lg:col-span-5">
      <h1 class="text-2xl md:text-3xl font-bold text-green-700 mb-4 text-center">
        {{ $t('Marketplace.title') }}
      </h1>

      <div class="flex flex-col md:flex-row justify-center items-center gap-2 mb-4">
        <!-- ฟอร์มค้นหาขยะ -->
        <div class="w-full md:w-auto">
          <div class="flex items-center gap-2">
            <IconSearch stroke="2" size="20" class="text-green-600 ml-2" />
            <input type="text" v-model="searchQuery" class="input input-bordered w-full md:w-80"
              :placeholder="$t('Marketplace.search')" />
          </div>
        </div>
        <div v-if="canAddProduct">
          <button class="btn bg-green-700 hover:bg-green-600 text-white flex items-center gap-2"
            @click="openModalWaste">
            <IconPlus stroke="2" />
            <span>{{ $t('Marketplace.add_product') }}</span>
          </button>
        </div>
      </div>
      <div class="flex flex-wrap gap-2 mb-4">
        <div class="dropdown">
          <div tabindex="0" role="button" class="btn btn-outline btn-sm flex items-center gap-2">
            <IconCalendar stroke="2" /> {{ $t('Marketplace.filter_date.title') }}
          </div>
          <ul tabindex="0"
            class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow border border-neutral/50">
            <li>
              <button @click="sortByLastUpdate" class="flex items-center gap-2">
                <IconSortAscending stroke="2" />
                {{ $t('Marketplace.filter_date.lastest') }}
              </button>
            </li>
            <li>
              <button @click="sortByLastUpdateDesc" class="flex items-center gap-2">
                <IconSortDescending stroke="2" />
                {{ $t('Marketplace.filter_date.oldest') }}
              </button>
            </li>
          </ul>
        </div>

        <div class="dropdown">
          <div tabindex="0" role="button" class="btn btn-outline btn-sm flex items-center gap-2">
            <IconArrowsSort stroke="2" /> {{ $t('Marketplace.filter_price.title') }}
          </div>
          <ul tabindex="0"
            class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow border border-neutral/50">
            <li>
              <button @click="sortByPrice" class="flex items-center gap-2">
                <IconTagMinus stroke="2" /> {{ $t('Marketplace.filter_price.asc') }}
              </button>
            </li>
            <li>
              <button @click="sortByPriceDesc" class="flex items-center gap-2">
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
          <div v-for="item in filteredItems" :key="item.name" class="flex justify-center">
            <Card :id="item.waste_ids && item.waste_ids[0] ? item.waste_ids[0] : ''" :name="item.name"
              :price="item.price" :category="item.category" :last_update="item.last_update" :url="item.url"
              @click="openProductDetailModal(item)" />
          </div>
        </TransitionGroup>
      </div>

      <!-- Pagination -->
      <div class="flex flex-col items-center mt-6">
        <div class="join">
          <button class="join-item btn" :disabled="currentPage <= 1" @click="goToPage(1)"
            :title="$t('Marketplace.pagination.first')">
            <IconChevronsLeft stroke="2" />
          </button>
          <button class="join-item btn" :disabled="currentPage <= 1" @click="goToPrevPage"
            :title="$t('Marketplace.pagination.prev')">
            <IconChevronLeft stroke="2" />
          </button>
          <button v-for="page in getVisiblePages()" :key="page" class="join-item btn"
            :class="{ 'btn-active': page === currentPage }" @click="goToPage(page)">
            {{ page }}
          </button>
          <button class="join-item btn" :disabled="currentPage >= wastesStore.groupedPagination.total_pages"
            @click="goToNextPage" :title="$t('Marketplace.pagination.next')">
            <IconChevronRight stroke="2" />
          </button>
          <button class="join-item btn" :disabled="currentPage >= wastesStore.groupedPagination.total_pages"
            @click="goToPage(wastesStore.groupedPagination.total_pages)" :title="$t('Marketplace.pagination.last')">
            <IconChevronsRight stroke="2" />
          </button>
        </div>

        <div class="mt-3 text-sm text-gray-600">
          {{ $t('Marketplace.pagination.page') }} {{ currentPage }} {{ $t('Marketplace.pagination.of') }} {{
            wastesStore.groupedPagination.total_pages }}
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
