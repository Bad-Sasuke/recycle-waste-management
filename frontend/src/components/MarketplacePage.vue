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
import { useShopStore } from '../stores/shop'
import type { GroupedRecyclableItem } from '../types/recycle_waste'
import { IconTagPlus, IconTagMinus, IconCalendar, IconPlus, IconSearch, IconCategory, IconSortAscending, IconSortDescending, IconArrowsSort, IconChevronLeft, IconChevronRight, IconChevronsLeft, IconChevronsRight, IconCamera, IconRecycle, IconBuildingStore, IconLeaf, IconChevronDown, IconFilter, IconX } from '@tabler/icons-vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const searchQuery = ref('')
const heroSearchQuery = ref('')
const selectedCategory = ref<string[]>([])
const showMobileFilter = ref(false)


const isLoading = ref(true)
const usersStore = useUsersStore()
const wastesStore = useWastesStore()
const categoryWasteStore = useCategoryWasteStore()
const shopStore = useShopStore()
const currentPage = ref(1)

// Animated stats
const animatedItemsCount = ref(0)
const animatedShopsCount = ref(0)
const animatedCO2Count = ref(0)

// Ref for product detail modal
const productDetailModal = ref<InstanceType<typeof ProductDetailModal> | null>(null)

// Stats computed
const totalItems = computed(() => wastesStore.groupedPagination.total_items || 0)
const totalShops = computed(() => shopStore.pagination.total_items || 0)
const estimatedCO2Saved = computed(() => {
  // Estimate: 2kg CO2 saved per recycled item on average
  return Math.round((totalItems.value * 2) / 1000 * 10) / 10
})

// Animate numbers
const animateValue = (start: number, end: number, duration: number, setter: (val: number) => void) => {
  const startTime = performance.now()
  const animate = (currentTime: number) => {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)
    const easeOutQuart = 1 - Math.pow(1 - progress, 4)
    const current = Math.round(start + (end - start) * easeOutQuart)
    setter(current)
    if (progress < 1) {
      requestAnimationFrame(animate)
    }
  }
  requestAnimationFrame(animate)
}

// Search from hero section
const handleHeroSearch = () => {
  searchQuery.value = heroSearchQuery.value
  // Scroll to products section
  const productsSection = document.getElementById('products-section')
  if (productsSection) {
    productsSection.scrollIntoView({ behavior: 'smooth' })
  }
}

onMounted(async () => {
  // Fetch shops for stats
  if (!shopStore.allShops.length) {
    await shopStore.fetchAllShops(1, 100)
  }

  if (!wastesStore.groupedWastes.length) {
    await wastesStore.fetchWastes(currentPage.value, 12)
  }
  if (!categoryWasteStore.category.length) {
    await categoryWasteStore.fetchCategoryWaste()
  }

  // Store the original data for filtering
  originalGroupedWastes.value = [...wastesStore.groupedWastes]

  isLoading.value = false

  // Animate stats after loading
  setTimeout(() => {
    animateValue(0, totalItems.value, 1500, (val) => animatedItemsCount.value = val)
    animateValue(0, totalShops.value, 1500, (val) => animatedShopsCount.value = val)
    animateValue(0, estimatedCO2Saved.value * 10, 1500, (val) => animatedCO2Count.value = val / 10)
  }, 300)
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

// Quick Categories with icons
const categoryIconMap: Record<string, string> = {
  'กระดาษ': '📦',
  'paper': '📦',
  'แก้ว': '🍾',
  'glass': '🍾',
  'โลหะ': '🔩',
  'metal': '🔩',
  'พลาสติก': '🥤',
  'plastic': '🥤',
  'อิเล็กทรอนิกส์': '⚡',
  'electronics': '⚡',
  'เหล็ก': '🔧',
  'steel': '🔧',
  'ทองแดง': '🟠',
  'copper': '🟠',
  'อลูมิเนียม': '🥫',
  'aluminum': '🥫',
}

const quickCategoryList = computed(() => {
  return categoryWasteStore.category.map((cat) => ({
    id: cat.id,
    name: cat.name ?? '',
    label: cat.name ?? '',
    icon: categoryIconMap[(cat.name ?? '').toLowerCase()] || '♻️'
  }))
})

const toggleQuickCategory = (categoryName: string) => {
  toggleCategory(categoryName)
  // Scroll to products section
  const productsSection = document.getElementById('products-section')
  if (productsSection) {
    productsSection.scrollIntoView({ behavior: 'smooth' })
  }
}

const clearCategoryFilter = () => {
  selectedCategory.value = []
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

  <div v-if="!isLoading">
    <!-- Hero Section -->
    <section class="hero-section relative overflow-hidden">
      <!-- Background with gradient -->
      <div class="absolute inset-0 bg-gradient-to-br from-emerald-600 via-green-500 to-teal-500"></div>

      <!-- Animated background circles -->
      <div class="absolute inset-0 overflow-hidden">
        <div class="floating-circle circle-1"></div>
        <div class="floating-circle circle-2"></div>
        <div class="floating-circle circle-3"></div>
      </div>

      <!-- Content -->
      <div class="relative z-10 container mx-auto px-4 py-12 md:py-20">
        <div class="text-center max-w-4xl mx-auto">
          <!-- Tagline -->
          <div
            class="inline-flex items-center gap-2 bg-white/20 backdrop-blur-sm px-4 py-2 rounded-full mb-6 animate-fade-in">
            <span class="text-white font-medium text-sm md:text-base">{{ $t('Marketplace.hero.tagline') }}</span>
          </div>

          <!-- Main Title -->
          <h1 class="text-3xl md:text-5xl lg:text-6xl font-bold text-white mb-4 leading-tight animate-slide-up">
            {{ $t('Marketplace.hero.title') }}
          </h1>

          <!-- Subtitle -->
          <p class="text-lg md:text-xl text-white/90 mb-8 max-w-2xl mx-auto animate-slide-up-delay">
            {{ $t('Marketplace.hero.subtitle') }}
          </p>

          <!-- Search Bar -->
          <div class="max-w-2xl mx-auto mb-10 animate-scale-in">
            <div class="relative flex items-center bg-white rounded-2xl shadow-2xl overflow-hidden p-1.5">
              <div class="flex-1 flex items-center px-4">
                <IconSearch stroke="2" size="24" class="text-gray-400 mr-3" />
                <input type="text" v-model="heroSearchQuery"
                  class="w-full py-4 text-lg outline-none placeholder-gray-400"
                  :placeholder="$t('Marketplace.hero.search_placeholder')" @keyup.enter="handleHeroSearch" />
              </div>
              <button @click="handleHeroSearch"
                class="hidden md:flex items-center gap-2 bg-gradient-to-r from-green-500 to-emerald-500 hover:from-green-600 hover:to-emerald-600 text-white px-6 py-4 rounded-xl font-semibold transition-all duration-300 transform hover:scale-105">
                <IconSearch stroke="2" size="20" />
                <span>{{ $t('Marketplace.search') }}</span>
              </button>
            </div>
            <!-- AI Scanner Button (Mobile Friendly) -->
            <div class="mt-4 flex justify-center gap-3">
              <button @click="router.push('/waste-scanner')"
                class="flex items-center gap-2 bg-white/20 hover:bg-white/30 backdrop-blur-sm text-white px-5 py-2.5 rounded-xl font-medium transition-all duration-300">
                <IconCamera stroke="2" size="20" />
                <span>{{ $t('Marketplace.hero.ai_scan') }}</span>
              </button>
              <button v-if="canAddProduct" @click="openModalWaste"
                class="flex items-center gap-2 bg-white/20 hover:bg-white/30 backdrop-blur-sm text-white px-5 py-2.5 rounded-xl font-medium transition-all duration-300">
                <IconPlus stroke="2" size="20" />
                <span>{{ $t('Marketplace.add_product') }}</span>
              </button>
            </div>
          </div>

          <!-- Stats Cards -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-3xl mx-auto animate-fade-in-delay">
            <!-- Items Count -->
            <div
              class="stat-card bg-white/15 backdrop-blur-md rounded-2xl p-5 border border-white/20 hover:bg-white/20 transition-all duration-300 transform hover:-translate-y-1">
              <div class="flex items-center justify-center gap-3">
                <div class="bg-white/20 p-3 rounded-xl">
                  <IconRecycle stroke="2" size="28" class="text-white" />
                </div>
                <div class="text-left">
                  <div class="text-3xl font-bold text-white">{{ animatedItemsCount.toLocaleString() }}</div>
                  <div class="text-white/80 text-sm">{{ $t('Marketplace.hero.stats.items') }}</div>
                </div>
              </div>
            </div>

            <!-- Shops Count -->
            <div
              class="stat-card bg-white/15 backdrop-blur-md rounded-2xl p-5 border border-white/20 hover:bg-white/20 transition-all duration-300 transform hover:-translate-y-1">
              <div class="flex items-center justify-center gap-3">
                <div class="bg-white/20 p-3 rounded-xl">
                  <IconBuildingStore stroke="2" size="28" class="text-white" />
                </div>
                <div class="text-left">
                  <div class="text-3xl font-bold text-white">{{ animatedShopsCount.toLocaleString() }}</div>
                  <div class="text-white/80 text-sm">{{ $t('Marketplace.hero.stats.shops') }}</div>
                </div>
              </div>
            </div>

            <!-- CO2 Saved -->
            <div
              class="stat-card bg-white/15 backdrop-blur-md rounded-2xl p-5 border border-white/20 hover:bg-white/20 transition-all duration-300 transform hover:-translate-y-1">
              <div class="flex items-center justify-center gap-3">
                <div class="bg-white/20 p-3 rounded-xl">
                  <IconLeaf stroke="2" size="28" class="text-white" />
                </div>
                <div class="text-left">
                  <div class="text-3xl font-bold text-white">{{ animatedCO2Count.toFixed(1) }}</div>
                  <div class="text-white/80 text-sm">{{ $t('Marketplace.hero.stats.co2_saved') }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Wave Divider -->
      <div class="absolute bottom-0 left-0 right-0">
        <svg viewBox="0 0 1440 100" fill="none" xmlns="http://www.w3.org/2000/svg" class="w-full">
          <path
            d="M0 50L48 45.7C96 41.3 192 32.7 288 35.8C384 39 480 54 576 58.3C672 62.7 768 56.3 864 47.5C960 38.7 1056 27.3 1152 27.5C1248 27.7 1344 39.3 1392 45.2L1440 51V101H1392C1344 101 1248 101 1152 101C1056 101 960 101 864 101C768 101 672 101 576 101C480 101 384 101 288 101C192 101 96 101 48 101H0V50Z"
            fill="white" />
        </svg>
      </div>
    </section>

    <!-- Quick Categories Section -->
    <section class="quick-categories-section bg-white py-6 shadow-sm">
      <div class="container mx-auto px-4">
        <div class="flex items-center gap-4 overflow-x-auto scrollbar-hide pb-2">
          <!-- Hot/All Button -->
          <button @click="clearCategoryFilter"
            class="quick-category-chip flex-shrink-0 flex items-center gap-2 px-5 py-2.5 rounded-full font-medium transition-all duration-300"
            :class="selectedCategory.length === 0 ? 'bg-gradient-to-r from-orange-500 to-red-500 text-white shadow-lg scale-105' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'">
            <span class="text-lg">🔥</span>
            <span>{{ $t('Marketplace.quickCategories.hot') }}</span>
          </button>

          <!-- Category Chips -->
          <button v-for="cat in quickCategoryList" :key="cat.id" @click="toggleQuickCategory(cat.name)"
            class="quick-category-chip flex-shrink-0 flex items-center gap-2 px-5 py-2.5 rounded-full font-medium transition-all duration-300"
            :class="selectedCategory.includes(cat.name.toLowerCase()) ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white shadow-lg scale-105' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'">
            <span class="text-lg">{{ cat.icon }}</span>
            <span>{{ cat.label }}</span>
          </button>
        </div>
      </div>
    </section>

    <!-- Products Section -->
    <div id="products-section" class="container px-4 py-8 grid grid-cols-1 lg:grid-cols-12 gap-8 max-w-full">
      <!-- ส่วนซ้าย : หมวดหมู่ (Sidebar) -->
      <!-- ส่วนซ้าย : หมวดหมู่ (Sidebar - Desktop Only) -->
      <div class="hidden lg:block lg:col-span-3">
        <div
          class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 sticky top-24 transition-all duration-300">
          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center gap-2">
              <IconCategory stroke="2" size="22" class="text-green-600" />
              <h2 class="text-lg font-bold text-gray-800">{{ $t('Marketplace.category') }}</h2>
            </div>
            <button v-if="selectedCategory.length > 0" @click="clearCategoryFilter"
              class="text-xs text-red-500 hover:text-red-600 hover:underline font-medium">
              {{ $t('Global.clear') || 'Clear' }}
            </button>
          </div>

          <div class="flex flex-col gap-2 max-h-[calc(100vh-200px)] overflow-y-auto pr-2 custom-scrollbar">
            <!-- All Categories Option -->
            <button @click="clearCategoryFilter"
              class="flex items-center justify-between w-full p-3 rounded-xl transition-all duration-200 group"
              :class="selectedCategory.length === 0 ? 'bg-green-50 text-green-700 border border-green-100' : 'hover:bg-gray-50 text-gray-600 border border-transparent'">
              <span class="font-medium">ทั้งหมด</span>
              <span class="text-xs px-2 py-0.5 rounded-full"
                :class="selectedCategory.length === 0 ? 'bg-green-200 text-green-800' : 'bg-gray-100 text-gray-400 group-hover:bg-gray-200'">
                {{ totalItems }}
              </span>
            </button>

            <div class="divider my-1"></div>

            <!-- Individual Categories -->
            <label v-for="category in categoryWasteStore.category" :key="category.id"
              class="flex items-center justify-between w-full p-3 rounded-xl cursor-pointer transition-all duration-200 border group select-none"
              :class="selectedCategory.includes((category.name ?? 'ไม่มี').toLowerCase()) ? 'bg-green-50 border-green-200 shadow-sm' : 'border-transparent hover:bg-gray-50'">
              <div class="flex items-center gap-3">
                <input type="checkbox" :value="category.name"
                  :checked="selectedCategory.includes((category.name ?? 'ไม่มี').toLowerCase())"
                  @change="toggleCategory(category?.name ?? 'ไม่มี')"
                  class="checkbox checkbox-xs checkbox-primary rounded-md" />
                <span class="text-sm font-medium"
                  :class="selectedCategory.includes((category.name ?? 'ไม่มี').toLowerCase()) ? 'text-green-700' : 'text-gray-600'">
                  {{ category?.name ?? 'ไม่มี' }}
                </span>
              </div>
            </label>
          </div>
        </div>
      </div>

      <!-- Mobile Filter Drawer (Overlay) -->
      <Transition name="fade">
        <div v-if="showMobileFilter" class="fixed inset-0 bg-black/50 z-50 lg:hidden" @click="showMobileFilter = false">
        </div>
      </Transition>

      <Transition name="slide-up">
        <div v-if="showMobileFilter"
          class="fixed bottom-0 left-0 right-0 bg-white z-50 rounded-t-3xl shadow-2xl lg:hidden max-h-[85vh] flex flex-col">
          <div
            class="p-4 border-b border-gray-100 flex items-center justify-between sticky top-0 bg-white rounded-t-3xl">
            <h2 class="text-xl font-bold text-gray-800 flex items-center gap-2">
              <IconCategory stroke="2" size="24" class="text-green-600" />
              {{ $t('Marketplace.category') }}
            </h2>
            <button @click="showMobileFilter = false" class="btn btn-circle btn-sm btn-ghost bg-gray-100">
              <IconX size="20" />
            </button>
          </div>

          <div class="p-4 overflow-y-auto flex-1 custom-scrollbar">
            <!-- Mobile category list content -->
            <div class="flex flex-col gap-2">
              <button @click="clearCategoryFilter(); showMobileFilter = false"
                class="flex items-center justify-between w-full p-4 rounded-xl transition-all duration-200 border"
                :class="selectedCategory.length === 0 ? 'bg-green-50 text-green-700 border-green-200' : 'bg-gray-50 text-gray-700 border-transparent'">
                <span class="font-bold">ทั้งหมด</span>
                <span class="badge badge-sm" :class="selectedCategory.length === 0 ? 'badge-primary' : 'badge-ghost'">{{
                  totalItems }}</span>
              </button>

              <label v-for="category in categoryWasteStore.category" :key="category.id"
                class="flex items-center justify-between w-full p-4 rounded-xl cursor-pointer transition-all duration-200 border"
                :class="selectedCategory.includes((category.name ?? 'ไม่มี').toLowerCase()) ? 'bg-green-50 border-green-200' : 'bg-white border-gray-100'">
                <span class="flex items-center gap-3">
                  <input type="checkbox" :value="category.name"
                    :checked="selectedCategory.includes((category.name ?? 'ไม่มี').toLowerCase())"
                    @change="toggleCategory(category?.name ?? 'ไม่มี')" class="checkbox checkbox-primary rounded-md" />
                  <span class="text-base font-medium">{{ category?.name ?? 'ไม่มี' }}</span>
                </span>
              </label>
            </div>
          </div>

          <div class="p-4 border-t border-gray-100 bg-white sticky bottom-0 safe-area-bottom">
            <button class="btn btn-primary w-full text-white text-lg rounded-xl shadow-lg shadow-green-200"
              @click="showMobileFilter = false">
              {{ $t('Global.close') }} ({{ filteredItems.length }})
            </button>
          </div>
        </div>
      </Transition>

      <!-- ส่วนขวา: รายการขยะ -->
      <div class="lg:col-span-9 col-span-1">
        <!-- Mobile Actions Bar (Filters & Sort) -->
        <div
          class="flex flex-col gap-3 mb-6 sticky top-[60px] lg:static z-20 bg-gray-50/95 lg:bg-transparent backdrop-blur-sm p-3 lg:p-0 rounded-xl lg:rounded-none border lg:border-none border-gray-200 shadow-sm lg:shadow-none transition-all duration-300">

          <!-- Mobile Sticky Search -->
          <div class="relative w-full lg:hidden">
            <IconSearch class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size="18" />
            <input type="text" v-model="searchQuery"
              :placeholder="$t('Marketplace.hero.search_placeholder') || 'Search...'"
              class="input input-bordered input-sm w-full pl-10 bg-white shadow-sm rounded-lg focus:outline-none focus:border-green-500 text-base" />
            <button v-if="searchQuery" @click="searchQuery = ''"
              class="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-400 p-1">
              <IconX size="14" />
            </button>
          </div>

          <div class="flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
            <div class="text-gray-500 text-sm font-medium hidden sm:block">
              {{ $t('Marketplace.pagination.page_info', {
                current: currentPage, total:
                  wastesStore.groupedPagination.total_pages, items: totalItems
              }) }}
            </div>

            <div class="flex items-center justify-between w-full sm:w-auto gap-2">
              <!-- Mobile Filter Button -->
              <button @click="showMobileFilter = true"
                class="btn btn-sm bg-white border-gray-200 text-gray-700 shadow-sm lg:hidden flex-1">
                <IconFilter size="18" />
                {{ $t('Marketplace.category') }}
                <div v-if="selectedCategory.length > 0" class="badge badge-xs badge-primary">{{ selectedCategory.length
                  }}
                </div>
              </button>

              <!-- Sort Dropdowns -->
              <div class="flex gap-2 flex-1 sm:flex-none justify-end">
                <!-- Date Filter -->
                <div class="dropdown dropdown-end">
                  <div tabindex="0" role="button"
                    class="btn btn-sm bg-white border border-gray-200 hover:border-green-500 hover:bg-white text-gray-700 font-normal gap-2 rounded-lg shadow-sm w-full sm:w-auto px-2 sm:px-3">
                    <IconCalendar stroke="1.5" size="18" />
                    <span class="hidden sm:inline">{{ $t('Marketplace.filter_date.title') }}</span>
                    <IconChevronDown class="w-3 h-3 opacity-50 ml-1" />
                  </div>
                  <ul tabindex="0"
                    class="dropdown-content z-[1] menu p-2 shadow-xl bg-white rounded-2xl w-48 border border-gray-100 mt-2">
                    <li class="menu-title text-xs text-gray-400 uppercase font-semibold px-3 py-2">Sort by Date</li>
                    <li>
                      <button @click="sortByLastUpdate" class="rounded-lg active:bg-green-50 active:text-green-700">
                        <IconSortAscending stroke="1.5" size="18" />
                        {{ $t('Marketplace.filter_date.lastest') }}
                      </button>
                    </li>
                    <li>
                      <button @click="sortByLastUpdateDesc" class="rounded-lg active:bg-green-50 active:text-green-700">
                        <IconSortDescending stroke="1.5" size="18" />
                        {{ $t('Marketplace.filter_date.oldest') }}
                      </button>
                    </li>
                  </ul>
                </div>

                <!-- Price Filter -->
                <div class="dropdown dropdown-end">
                  <div tabindex="0" role="button"
                    class="btn btn-sm bg-white border border-gray-200 hover:border-green-500 hover:bg-white text-gray-700 font-normal gap-2 rounded-lg shadow-sm w-full sm:w-auto px-2 sm:px-3">
                    <IconArrowsSort stroke="1.5" size="18" />
                    <span class="hidden sm:inline">{{ $t('Marketplace.filter_price.title') }}</span>
                    <IconChevronDown class="w-3 h-3 opacity-50 ml-1" />
                  </div>
                  <ul tabindex="0"
                    class="dropdown-content z-[1] menu p-2 shadow-xl bg-white rounded-2xl w-48 border border-gray-100 mt-2">
                    <li class="menu-title text-xs text-gray-400 uppercase font-semibold px-3 py-2">Sort by Price</li>
                    <li>
                      <button @click="sortByPrice" class="rounded-lg active:bg-green-50 active:text-green-700">
                        <IconTagMinus stroke="1.5" size="18" /> {{ $t('Marketplace.filter_price.asc') }}
                      </button>
                    </li>
                    <li>
                      <button @click="sortByPriceDesc" class="rounded-lg active:bg-green-50 active:text-green-700">
                        <IconTagPlus stroke="1.5" size="18" />{{ $t('Marketplace.filter_price.desc') }}
                      </button>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- รายการราคาขยะ -->
        <div class="min-h-[600px]">
          <!-- Empty State -->
          <div v-if="filteredItems.length === 0"
            class="flex flex-col items-center justify-center h-96 text-center p-8 bg-gray-50 rounded-3xl border border-dashed border-gray-200">
            <div class="bg-white p-4 rounded-full shadow-sm mb-4">
              <IconSearch class="w-12 h-12 text-gray-300" stroke="1.5" />
            </div>
            <h3 class="text-lg font-bold text-gray-700 mb-1">ไม่พบสินค้าที่คุณค้นหา</h3>
            <p class="text-gray-500 text-sm max-w-sm mx-auto">ลองเปลี่ยนคำค้นหา หรือ เลือกหมวดหมู่อื่นดูนะครับ</p>
            <button @click="clearCategoryFilter(); searchQuery = ''"
              class="btn btn-sm btn-outline btn-success mt-4">ล้างตัวกรองทั้งหมด</button>
          </div>

          <!-- Grid Items -->
          <TransitionGroup name="list" tag="div" class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 md:gap-6"
            v-else>
            <div v-for="item in filteredItems" :key="item.name" class="h-full">
              <Card :id="item.waste_ids && item.waste_ids[0] ? item.waste_ids[0] : ''" :name="item.name"
                :price="item.price" :category="item.category" :last_update="item.last_update" :url="item.url"
                @click="openProductDetailModal(item)" />
            </div>
          </TransitionGroup>
        </div>

        <!-- Pagination -->
        <div class="flex flex-col items-center mt-12 mb-4">
          <div class="flex items-center gap-2">
            <!-- First Page -->
            <button
              class="btn btn-sm btn-circle btn-ghost bg-white border border-gray-100 shadow-sm hover:bg-green-50 hover:text-green-600 disabled:bg-gray-50 disabled:text-gray-300 transition-all duration-200"
              :disabled="currentPage <= 1" @click="goToPage(1)" :title="$t('Marketplace.pagination.first')">
              <IconChevronsLeft stroke="2" size="18" />
            </button>

            <!-- Prev Page -->
            <button
              class="btn btn-sm btn-circle btn-ghost bg-white border border-gray-100 shadow-sm hover:bg-green-50 hover:text-green-600 disabled:bg-gray-50 disabled:text-gray-300 transition-all duration-200"
              :disabled="currentPage <= 1" @click="goToPrevPage" :title="$t('Marketplace.pagination.prev')">
              <IconChevronLeft stroke="2" size="18" />
            </button>

            <!-- Page Numbers -->
            <div class="flex items-center gap-1 mx-2">
              <button v-for="page in getVisiblePages()" :key="page"
                class="btn btn-sm w-9 h-9 border-none shadow-sm transition-all duration-200 rounded-xl text-sm font-bold"
                :class="page === currentPage
                  ? 'bg-green-500 text-white hover:bg-green-600 shadow-green-200 shadow-md transform scale-105'
                  : 'bg-white text-gray-600 hover:bg-gray-50 hover:text-green-600'" @click="goToPage(page)">
                {{ page }}
              </button>
            </div>

            <!-- Next Page -->
            <button
              class="btn btn-sm btn-circle btn-ghost bg-white border border-gray-100 shadow-sm hover:bg-green-50 hover:text-green-600 disabled:bg-gray-50 disabled:text-gray-300 transition-all duration-200"
              :disabled="currentPage >= wastesStore.groupedPagination.total_pages" @click="goToNextPage"
              :title="$t('Marketplace.pagination.next')">
              <IconChevronRight stroke="2" size="18" />
            </button>

            <!-- Last Page -->
            <button
              class="btn btn-sm btn-circle btn-ghost bg-white border border-gray-100 shadow-sm hover:bg-green-50 hover:text-green-600 disabled:bg-gray-50 disabled:text-gray-300 transition-all duration-200"
              :disabled="currentPage >= wastesStore.groupedPagination.total_pages"
              @click="goToPage(wastesStore.groupedPagination.total_pages)" :title="$t('Marketplace.pagination.last')">
              <IconChevronsRight stroke="2" size="18" />
            </button>
          </div>

          <div class="mt-4 text-xs font-medium text-gray-400 uppercase tracking-wider">
            {{ $t('Marketplace.pagination.page') }} <span class="text-gray-800 font-bold mx-1">{{ currentPage }}</span>
            {{ $t('Marketplace.pagination.of') }} <span class="text-gray-800 font-bold mx-1">{{
              wastesStore.groupedPagination.total_pages }}</span>
          </div>
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

/* Hero Section Styles */
.hero-section {
  min-height: 70vh;
  display: flex;
  align-items: center;
}

/* Floating Circles Animation */
.floating-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 15s infinite ease-in-out;
}

.circle-1 {
  width: 400px;
  height: 400px;
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.circle-2 {
  width: 300px;
  height: 300px;
  bottom: -50px;
  right: -50px;
  animation-delay: -5s;
}

.circle-3 {
  width: 200px;
  height: 200px;
  top: 50%;
  right: 20%;
  animation-delay: -10s;
}

@keyframes float {

  0%,
  100% {
    transform: translateY(0) rotate(0deg);
  }

  25% {
    transform: translateY(-20px) rotate(5deg);
  }

  50% {
    transform: translateY(0) rotate(0deg);
  }

  75% {
    transform: translateY(20px) rotate(-5deg);
  }
}

/* Animation Classes */
.animate-fade-in {
  animation: fadeIn 0.8s ease-out forwards;
}

.animate-fade-in-delay {
  animation: fadeIn 0.8s ease-out 0.5s forwards;
  opacity: 0;
}

.animate-slide-up {
  animation: slideUp 0.8s ease-out forwards;
}

.animate-slide-up-delay {
  animation: slideUp 0.8s ease-out 0.2s forwards;
  opacity: 0;
  transform: translateY(30px);
}

.animate-scale-in {
  animation: scaleIn 0.6s ease-out 0.4s forwards;
  opacity: 0;
  transform: scale(0.95);
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }

  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* Stat Card Hover Effect */
.stat-card {
  transition: all 0.3s ease;
}

.stat-card:hover {
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
}

/* Quick Categories Styles */
.quick-categories-section {
  position: relative;
  z-index: 10;
}

.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.quick-category-chip {
  white-space: nowrap;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.quick-category-chip:hover {
  transform: translateY(-2px);
}

.quick-category-chip:active {
  transform: scale(0.95);
}

/* Custom Scrollbar for Sidebar */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #cbd5e1;
  border-radius: 20px;
}

/* Safe Area for Mobile */
.safe-area-bottom {
  padding-bottom: env(safe-area-inset-bottom, 20px);
}
</style>
