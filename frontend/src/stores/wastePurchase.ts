import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface WastePurchaseItem {
  id: string
  waste_id: string
  name: string
  icon: string
  iconComponent: any
  color: string
  totalWeight: number
  expression: string
  price?: number
  category?: string
}

export const useWastePurchaseStore = defineStore(
  'wastePurchase',
  () => {
    // State
    const items = ref<WastePurchaseItem[]>([])
    const currentRequestId = ref<string | null>(null)
    const isBottomSheetExpanded = ref(false)

    // Computed
    const grandTotalWeight = computed(() => {
      return items.value.reduce((sum, item) => sum + item.totalWeight, 0).toFixed(2)
    })

    const itemCount = computed(() => items.value.length)

    // Actions
    function addItem(item: WastePurchaseItem) {
      // Check if item with same waste_id already exists
      const existingIndex = items.value.findIndex((i) => i.waste_id === item.waste_id)

      if (existingIndex !== -1) {
        // Merge weights and expressions
        const existing = items.value[existingIndex]
        existing.totalWeight += item.totalWeight
        existing.expression = `${existing.expression} + ${item.expression}`
      } else {
        items.value.push(item)
      }
    }

    function updateItemWeight(itemId: string, newWeight: number, newExpression: string) {
      const item = items.value.find((i) => i.id === itemId)
      if (item) {
        item.totalWeight = newWeight
        item.expression = newExpression
      }
    }

    function removeItem(itemId: string) {
      items.value = items.value.filter((i) => i.id !== itemId)
    }

    function clearItems() {
      items.value = []
      isBottomSheetExpanded.value = false
    }

    function setRequestId(requestId: string | null) {
      // If request ID changes, clear items
      if (currentRequestId.value !== requestId) {
        clearItems()
      }
      currentRequestId.value = requestId
    }

    function toggleBottomSheet() {
      isBottomSheetExpanded.value = !isBottomSheetExpanded.value
    }

    function setBottomSheetExpanded(expanded: boolean) {
      isBottomSheetExpanded.value = expanded
    }

    return {
      // State
      items,
      currentRequestId,
      isBottomSheetExpanded,

      // Computed
      grandTotalWeight,
      itemCount,

      // Actions
      addItem,
      updateItemWeight,
      removeItem,
      clearItems,
      setRequestId,
      toggleBottomSheet,
      setBottomSheetExpanded,
    }
  },
  {
    persist: {
      key: 'wastePurchase',
      storage: localStorage,
    },
  },
)
