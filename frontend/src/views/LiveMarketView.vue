<template>
  <div class="min-h-screen bg-gray-900 text-white font-sans overflow-hidden relative">
    <!-- Background Effects -->
    <div class="absolute top-0 left-0 w-96 h-96 bg-primary/20 rounded-full blur-3xl -translate-x-1/2 -translate-y-1/2">
    </div>
    <div
      class="absolute bottom-0 right-0 w-96 h-96 bg-secondary/20 rounded-full blur-3xl translate-x-1/2 translate-y-1/2">
    </div>

    <!-- Ticker Tape -->
    <div
      class="bg-gray-800/80 backdrop-blur-md border-b border-gray-700 h-10 flex items-center overflow-hidden relative z-10">
      <div class="animate-marquee whitespace-nowrap flex gap-8 items-center">
        <div v-for="(item, index) in tickerItems" :key="index" class="flex items-center gap-2">
          <span class="font-bold text-gray-300">{{ item.name }}</span>
          <span :class="item.change >= 0 ? 'text-green-400' : 'text-red-400'"
            class="flex items-center text-sm font-mono">
            <span v-if="item.change >= 0">▲</span>
            <span v-else>▼</span>
            {{ Math.abs(item.change.toFixed(2)) }}%
          </span>
          <span class="text-xs text-gray-500">฿{{ item.price.toFixed(2) }}</span>
        </div>
        <!-- Repeat for seamless loop -->
        <div v-for="(item, index) in tickerItems" :key="`dup-${index}`" class="flex items-center gap-2">
          <span class="font-bold text-gray-300">{{ item.name }}</span>
          <span :class="item.change >= 0 ? 'text-green-400' : 'text-red-400'"
            class="flex items-center text-sm font-mono">
            <span v-if="item.change >= 0">▲</span>
            <span v-else>▼</span>
            {{ Math.abs(item.change.toFixed(2)) }}%
          </span>
          <span class="text-xs text-gray-500">฿{{ item.price.toFixed(2) }}</span>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="container mx-auto px-4 py-8 relative z-10">
      <header class="mb-8 flex justify-between items-end">
        <div>
          <h1
            class="text-4xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-green-400 to-blue-500 mb-2">
            Live Waste Market
          </h1>
          <p class="text-gray-400">Real-time global recycling prices & trends</p>
        </div>
        <div class="flex gap-2 text-sm text-gray-500">
          <span class="w-2 h-2 rounded-full bg-green-500 animate-pulse mt-1.5"></span>
          Live Connection
        </div>
      </header>

      <div class="space-y-6">
        <!-- Market Overview Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div v-for="category in categories" :key="category.id"
            class="bg-gray-800/50 backdrop-blur-xl border border-gray-700/50 rounded-2xl p-6 hover:border-gray-600 transition-all duration-300 group">
            <div class="flex justify-between items-start mb-4">
              <div class="flex items-center gap-3">
                <div :class="`p-3 rounded-xl bg-gradient-to-br ${category.color}`">
                  <component :is="category.icon" class="w-6 h-6 text-white" />
                </div>
                <div>
                  <h3 class="font-bold text-lg text-gray-200">{{ category.name }}</h3>
                  <p class="text-xs text-gray-400">Global Avg.</p>
                </div>
              </div>
              <div class="text-right">
                <div class="text-2xl font-mono font-bold text-white transition-all duration-300" :key="category.price">
                  ฿{{ category.price.toFixed(2) }}
                </div>
                <div :class="category.trend >= 0 ? 'text-green-400' : 'text-red-400'"
                  class="text-xs font-bold flex justify-end items-center gap-1">
                  {{ category.trend >= 0 ? '+' : '' }}{{ category.trend.toFixed(2) }}%
                </div>
              </div>
            </div>

            <!-- Mini Chart Area (Simulated) -->
            <div class="h-16 flex items-end gap-1 opacity-50 group-hover:opacity-100 transition-opacity">
              <div v-for="(bar, i) in category.history" :key="i" class="flex-1 rounded-t-sm transition-all duration-500"
                :class="category.trend >= 0 ? 'bg-green-500/50' : 'bg-red-500/50'" :style="{ height: `${bar}%` }">
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Live Chart Section -->
      <div class="mt-6 bg-gray-800/50 backdrop-blur-xl border border-gray-700/50 rounded-2xl p-6">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-bold text-lg text-gray-200">Price Trend Analysis (7 Days)</h3>
          <div class="flex gap-2">
            <button
              class="px-3 py-1 rounded-lg bg-gray-700 text-xs text-gray-300 hover:bg-gray-600 transition">1H</button>
            <button
              class="px-3 py-1 rounded-lg bg-green-500/20 text-green-400 border border-green-500/30 text-xs">1D</button>
            <button
              class="px-3 py-1 rounded-lg bg-gray-700 text-xs text-gray-300 hover:bg-gray-600 transition">1W</button>
          </div>
        </div>
        <div class="h-64 w-full">
          <Line v-if="chartData.datasets.length > 0" :data="chartData" :options="chartOptions" />
          <div v-else class="flex items-center justify-center h-full text-gray-500">
            Initializing Real-time Data...
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, markRaw, computed } from 'vue'
import { IconBottle, IconGlassFull, IconBox, IconCpu, IconRecycle } from '@tabler/icons-vue'
import { useWastesStore } from '@/stores/wastes'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import { Line } from 'vue-chartjs'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const wastesStore = useWastesStore()

// Mapped icons for known categories
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const categoryIcons: Record<string, any> = {
  'Plastic': IconBottle,
  'Glass': IconGlassFull,
  'Paper': IconBox,
  'Metal': IconCpu,
  'Aluminium': IconCpu,
}

const getCategoryIcon = (catName: string) => {
  for (const key in categoryIcons) {
    if (catName.toLowerCase().includes(key.toLowerCase())) {
      return markRaw(categoryIcons[key])
    }
  }
  return markRaw(IconRecycle)
}

const getCategoryColor = (index: number) => {
  const colors = [
    'from-blue-500 to-cyan-400',
    'from-orange-500 to-yellow-400',
    'from-green-500 to-emerald-400',
    'from-purple-500 to-pink-400',
    'from-red-500 to-pink-500',
    'from-indigo-500 to-purple-500'
  ]
  return colors[index % colors.length]
}

const getChartColor = (index: number) => {
  const colors = [
    'rgba(59, 130, 246, 1)',   // Blue
    'rgba(249, 115, 22, 1)',   // Orange
    'rgba(16, 185, 129, 1)',   // Emerald
    'rgba(168, 85, 247, 1)',   // Purple
    'rgba(236, 72, 153, 1)',   // Pink
    'rgba(99, 102, 241, 1)'    // Indigo
  ]
  return colors[index % colors.length]
}

// Reactive data derived from store
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const tickerItems = computed<any[]>(() => {
  if (wastesStore.groupedWastes.length > 0) {
    return wastesStore.groupedWastes.slice(0, 10).map(item => ({
      name: item.name,
      price: item.price,
      change: (Math.random() * 5) - 2.5 // Mock change
    }))
  }
  return [
    { name: 'Loading Market Data...', price: 0, change: 0 }
  ]
})

// Categories aggregation
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const categories = ref<any[]>([])

// Chart Data Logic
const chartLabels = ref<string[]>(Array.from({ length: 20 }, (_, i) => `${20 - i}m ago`))

const chartData = computed(() => {
  return {
    labels: chartLabels.value,
    datasets: categories.value.map((cat, index) => ({
      label: cat.name,
      borderColor: getChartColor(index),
      backgroundColor: getChartColor(index).replace('1)', '0.1)'),
      data: cat.history,
      tension: 0.4,
      fill: true,
      pointRadius: 0
    }))
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: {
      grid: {
        color: 'rgba(255, 255, 255, 0.1)',
        drawBorder: false
      },
      ticks: {
        color: '#9ca3af'
      }
    },
    y: {
      grid: {
        color: 'rgba(255, 255, 255, 0.1)',
        drawBorder: false
      },
      ticks: {
        color: '#9ca3af'
      }
    }
  },
  plugins: {
    legend: {
      labels: {
        color: '#e5e7eb',
        usePointStyle: true
      }
    },
    tooltip: {
      mode: 'index' as const,
      intersect: false,
      backgroundColor: 'rgba(17, 24, 39, 0.9)',
      titleColor: '#fff',
      bodyColor: '#e5e7eb',
      borderColor: 'rgba(75, 85, 99, 0.5)',
      borderWidth: 1
    }
  },
  interaction: {
    mode: 'nearest' as const,
    axis: 'x' as const,
    intersect: false
  }
}

const processCategories = () => {
  const groups: Record<string, { total: number, count: number, name: string }> = {}

  // Group items by their category field
  wastesStore.groupedWastes.forEach(item => {
    const cat = item.category || 'Other'
    if (!groups[cat]) {
      groups[cat] = { total: 0, count: 0, name: cat }
    }
    groups[cat].total += item.price
    groups[cat].count++
  })

  // Helper to keep existing history if available to avoid graph reset
  const createHistory = () => {
    let prevPrice = Math.random() * 50
    return Array.from({ length: 20 }, () => {
      prevPrice += (Math.random() - 0.5) * 5
      return Math.abs(prevPrice)
    })
  }

  // Convert to array for UI
  categories.value = Object.values(groups).map((group, index) => {
    // Try to preserve existing history if category already existed
    const existing = categories.value.find(c => c.name === group.name)
    const avgPrice = group.total / group.count

    return {
      id: index,
      name: group.name,
      price: avgPrice, // Average price
      trend: (Math.random() * 5) - 2.5, // Mock trend
      color: getCategoryColor(index),
      icon: getCategoryIcon(group.name),
      history: existing ? existing.history : createHistory()
    }
  })

  // If empty (no data), show placeholder
  if (categories.value.length === 0) {
    categories.value = [
      { id: 1, name: 'Market Offline', price: 0, trend: 0, color: 'from-gray-700 to-gray-600', icon: markRaw(IconRecycle), history: [] }
    ]
  }
}

let interval: ReturnType<typeof setInterval>

onMounted(async () => {
  await wastesStore.fetchWastes(1, 100)
  processCategories()

  // Simulate Live Ticker
  interval = setInterval(() => {
    if (categories.value.length === 0) return

    categories.value.forEach(cat => {
      // Mock price movement based on last history value
      const lastValue = cat.history[cat.history.length - 1]
      const move = (Math.random() - 0.5) * 2
      let newValue = lastValue + move
      if (newValue < 0) newValue = 0

      cat.history.shift()
      cat.history.push(newValue)

      // Update current price visual
      cat.price = newValue
      cat.trend = updateTrend(cat.trend, move)
    })

  }, 1500)
})

onUnmounted(() => {
  clearInterval(interval)
})

function updateTrend(currentTrend: number, move: number) {
  // Simple smoothing
  return (currentTrend * 0.9) + (move * 10)
}
</script>

<style scoped>
@keyframes marquee {
  0% {
    transform: translateX(0);
  }

  100% {
    transform: translateX(-50%);
  }
}

.animate-marquee {
  animation: marquee 30s linear infinite;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}
</style>
