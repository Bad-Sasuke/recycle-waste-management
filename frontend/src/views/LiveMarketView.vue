<template>
  <div class="min-h-screen bg-[#0b0e11] text-white font-sans overflow-x-hidden relative">
    <!-- Background Effects -->
    <div class="absolute top-0 left-0 w-full h-96 bg-gradient-to-b from-[#1e222d] to-transparent opacity-50 z-0"></div>

    <!-- Loading Overlay -->
    <LoadingOverlay :loading="wastesStore.loading && loadingInitial" text="Initializing Trading Terminals..." />

    <!-- Ticker Tape -->
    <div
      class="bg-[#171b26] border-b border-gray-800 h-10 flex items-center overflow-hidden relative z-10 text-xs font-mono tracking-wide">
      <div class="animate-marquee whitespace-nowrap flex gap-12 items-center">
        <template v-for="(item, index) in tickerItems" :key="index">
          <div class="flex items-center gap-3">
            <span class="font-bold text-gray-400">{{ item.name }}</span>
            <span :class="item.change >= 0 ? 'text-[#00ff7f]' : 'text-[#ff3b3b]'" class="flex items-center">
              {{ item.change >= 0 ? '▲' : '▼' }} {{ Math.abs(item.change).toFixed(2) }}%
            </span>
            <span class="text-gray-500">฿{{ item.price.toFixed(2) }}</span>
          </div>
        </template>
        <!-- Duplicate for loop -->
        <template v-for="(item, index) in tickerItems" :key="`dup-${index}`">
          <div class="flex items-center gap-3">
            <span class="font-bold text-gray-400">{{ item.name }}</span>
            <span :class="item.change >= 0 ? 'text-[#00ff7f]' : 'text-[#ff3b3b]'" class="flex items-center">
              {{ item.change >= 0 ? '▲' : '▼' }} {{ Math.abs(item.change).toFixed(2) }}%
            </span>
            <span class="text-gray-500">฿{{ item.price.toFixed(2) }}</span>
          </div>
        </template>
      </div>
    </div>

    <!-- Main Dashboard -->
    <div class="container mx-auto px-4 py-6 relative z-10 max-w-[1600px]">
      <header class="mb-6 flex justify-between items-end">
        <div>
          <h1 class="text-3xl font-bold text-white mb-1 flex items-center gap-2">
            <span class="w-2 h-8 bg-[#00ff7f] rounded-sm"></span>
            Global Waste Exchange
          </h1>
          <p class="text-gray-500 text-sm font-mono">LIVE MARKET DATA • BANGKOK SESSION</p>
        </div>
        <div class="flex gap-4">
          <div class="text-right">
            <div class="text-xs text-gray-500">MARKET STATUS</div>
            <div class="text-[#00ff7f] font-mono text-sm tracking-wider flex items-center gap-2 justify-end">
              <span class="w-2 h-2 rounded-full bg-[#00ff7f] animate-pulse"></span>
              OPEN
            </div>
          </div>
        </div>
      </header>

      <!-- Indices Section -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <div v-for="(index, i) in marketIndices" :key="i"
          class="bg-[#1e222d] border border-gray-800 rounded-lg p-4 hover:border-gray-600 transition-colors group cursor-pointer"
          @click="selectIndex(index)">
          <div class="flex justify-between items-start mb-2">
            <div class="text-sm text-gray-400 font-bold tracking-wider">{{ index.name }}</div>
            <component :is="getCategoryIcon(index.name)"
              class="w-5 h-5 text-gray-600 group-hover:text-white transition-colors" />
          </div>
          <div class="flex items-baseline gap-2 mb-2">
            <span class="text-2xl font-mono font-bold">{{ index.value.toFixed(2) }}</span>
            <span :class="index.change >= 0 ? 'text-[#00ff7f] bg-[#00ff7f]/10' : 'text-[#ff3b3b] bg-[#ff3b3b]/10'"
              class="text-xs px-1.5 py-0.5 rounded font-mono">
              {{ index.change >= 0 ? '+' : '' }}{{ index.change.toFixed(2) }}%
            </span>
          </div>
          <!-- Mini Sparkline using simple SVG -->
          <svg class="w-full h-8 stroke-current" :class="index.change >= 0 ? 'text-[#00ff7f]' : 'text-[#ff3b3b]'"
            fill="none" stroke-width="2">
            <polyline :points="generateSparkline(index.history)" />
          </svg>
        </div>
      </div>

      <!-- Main Trading Interface -->
      <div class="grid grid-cols-12 gap-4 mb-6 h-[600px]">
        <!-- Chart Area -->
        <div class="col-span-12 lg:col-span-9 bg-[#1e222d] border border-gray-800 rounded-lg flex flex-col">
          <!-- Toolbar -->
          <div class="h-12 border-b border-gray-800 flex items-center px-4 justify-between bg-[#171b26] rounded-t-lg">
            <div class="flex items-center gap-4">
              <span class="font-bold text-gray-200">{{ selectedAsset }} / THB</span>
              <div class="flex gap-1 text-xs text-gray-500">
                <button class="px-2 py-1 rounded hover:bg-gray-700 hover:text-white transition">15m</button>
                <button class="px-2 py-1 rounded bg-gray-700 text-white transition">1H</button>
                <button class="px-2 py-1 rounded hover:bg-gray-700 hover:text-white transition">4H</button>
                <button class="px-2 py-1 rounded hover:bg-gray-700 hover:text-white transition">1D</button>
              </div>
            </div>
            <div class="flex gap-4 text-xs font-mono">
              <div class="flex gap-2">
                <span class="text-gray-500">O:</span> <span class="text-[#00ff7f]">{{ currentOHLC.open }}</span>
              </div>
              <div class="flex gap-2">
                <span class="text-gray-500">H:</span> <span class="text-[#00ff7f]">{{ currentOHLC.high }}</span>
              </div>
              <div class="flex gap-2">
                <span class="text-gray-500">L:</span> <span class="text-[#00ff7f]">{{ currentOHLC.low }}</span>
              </div>
              <div class="flex gap-2">
                <span class="text-gray-500">C:</span> <span class="text-[#00ff7f]">{{ currentOHLC.close }}</span>
              </div>
            </div>
          </div>

          <!-- Charts Container -->
          <div class="flex-1 relative p-2 flex flex-col">
            <!-- Candlestick Chart -->
            <div class="h-2/3 w-full">
              <VueApexCharts type="candlestick" height="100%" :options="candleOptions" :series="candleSeries" />
            </div>
            <!-- Volume Oscillator -->
            <div class="h-1/3 w-full border-t border-gray-800 pt-2">
              <VueApexCharts type="bar" height="100%" :options="volumeOptions" :series="volumeSeries" />
            </div>
          </div>
        </div>

        <!-- Order Book (Market Depth) -->
        <div
          class="col-span-12 lg:col-span-3 bg-[#1e222d] border border-gray-800 rounded-lg flex flex-col h-full overflow-hidden">
          <div class="h-12 border-b border-gray-800 flex items-center px-4 bg-[#171b26] font-bold text-sm rounded-t-lg">
            Order Book
          </div>
          <div class="flex-1 overflow-y-auto custom-scrollbar font-mono text-xs">
            <!-- Asks (Sellers) -->
            <div class="flex flex-col-reverse"> <!-- Flex reverse to put lowest ask at bottom -->
              <div v-for="(ask, i) in orderBook.asks" :key="`ask-${i}`"
                class="relative py-1 px-2 hover:bg-gray-800 cursor-pointer flex justify-between">
                <div class="absolute inset-y-0 right-0 bg-[#ff3b3b]/10 transition-all duration-300"
                  :style="`width: ${ask.percent}%`"></div>
                <span class="text-[#ff3b3b] relative z-10">{{ ask.price.toFixed(2) }}</span>
                <span class="text-gray-400 relative z-10">{{ ask.volume }}</span>
              </div>
            </div>

            <!-- Current Price Divider -->
            <div
              class="py-2 border-y border-gray-700 bg-[#171b26] text-center text-lg font-bold text-white my-1 sticky top-0 bottom-0 z-20">
              {{ currentPrice.toFixed(2) }} <span class="text-xs text-gray-500">THB</span>
            </div>

            <!-- Bids (Buyers) -->
            <div>
              <div v-for="(bid, i) in orderBook.bids" :key="`bid-${i}`"
                class="relative py-1 px-2 hover:bg-gray-800 cursor-pointer flex justify-between">
                <div class="absolute inset-y-0 right-0 bg-[#00ff7f]/10 transition-all duration-300"
                  :style="`width: ${bid.percent}%`"></div>
                <span class="text-[#00ff7f] relative z-10">{{ bid.price.toFixed(2) }}</span>
                <span class="text-gray-400 relative z-10">{{ bid.volume }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Analysis Section -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Sector Performance (Treemap) -->
        <div class="bg-[#1e222d] border border-gray-800 rounded-lg p-6">
          <h3 class="text-lg font-bold text-white mb-4">Sector Performance Heatmap</h3>
          <div class="h-80 w-full">
            <VueApexCharts type="treemap" height="100%" :options="treemapOptions" :series="treemapSeries" />
          </div>
        </div>

        <!-- Correlation Matrix -->
        <div class="bg-[#1e222d] border border-gray-800 rounded-lg p-6">
          <h3 class="text-lg font-bold text-white mb-4">Market Correlation Matrix</h3>
          <div class="h-80 w-full">
            <VueApexCharts type="heatmap" height="100%" :options="heatmapOptions" :series="heatmapSeries" />
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, markRaw, computed } from 'vue'
import { IconBottle, IconGlassFull, IconBox, IconCpu, IconRecycle } from '@tabler/icons-vue'
import VueApexCharts from 'vue3-apexcharts'
import { generateCandles, generateOrderBook, generateCorrelationMatrix, generateVolumeData, type CandleData } from '@/services/marketSimulator'
import LoadingOverlay from '@/components/common/LoadingOverlay.vue'
import config from '@/config'
import { useWastesStore } from '@/stores/wastes'

const wastesStore = useWastesStore()

// State
const loadingInitial = ref(true)
const selectedAsset = ref('Plastic Market')
const currentPrice = ref(12.50)
const currentOHLC = ref({ open: 0, high: 0, low: 0, close: 0 })

// Order Book State
const orderBook = ref(generateOrderBook(12.50))

// Icons
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const getCategoryIcon = (catName: string): any => {
  if (catName.includes('Plastic')) return markRaw(IconBottle)
  if (catName.includes('Glass')) return markRaw(IconGlassFull)
  if (catName.includes('Paper')) return markRaw(IconBox)
  if (catName.includes('Metal') || catName.includes('Alu')) return markRaw(IconCpu)
  return markRaw(IconRecycle)
}

// Sparkline Helper
const generateSparkline = (history: number[]) => {
  if (!history.length) return ''
  const min = Math.min(...history)
  const max = Math.max(...history)
  const range = (max - min) === 0 ? 1 : (max - min)
  const width = 250 // Arbitrary SVG width
  const height = 30 // Arbitrary SVG height

  return history.map((val, i) => {
    const x = (i / (history.length - 1)) * width
    const y = height - ((val - min) / range) * height
    return `${x},${y}`
  }).join(' ')
}

// ----------------------------------------------------------------------
// DATA: Ticker & Indices
// ----------------------------------------------------------------------
// Use real data from store if available
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const tickerItems = computed<any[]>(() => {
  if (wastesStore.groupedWastes.length > 0) {
    return wastesStore.groupedWastes.map(item => ({
      name: item.name,
      price: item.price,
      change: (Math.random() * 5) - 2.5 // Simulated change
    }))
  }
  return [
    { name: 'Loading Market...', price: 0, change: 0 }
  ]
})

const marketIndices = ref([
  { name: 'Plastic Index', query: 'Plastic', value: 0, change: 0, history: [0] },
  { name: 'Paper Index', query: 'Paper', value: 0, change: 0, history: [0] },
  { name: 'Metal Index', query: 'Metal', value: 0, change: 0, history: [0] },
  { name: 'Glass Index', query: 'Glass', value: 0, change: 0, history: [0] },
])

// ----------------------------------------------------------------------
// CHART: Candlestick
// ----------------------------------------------------------------------
const candleData = ref<CandleData[]>(generateCandles(12.50, 60))

const candleSeries = computed(() => [{
  data: candleData.value
}])

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const candleOptions = {
  chart: {
    type: 'candlestick' as const,
    height: 350,
    background: 'transparent',
    toolbar: { show: false },
    animations: { enabled: false }
  },
  theme: { mode: 'dark' as const },
  grid: {
    borderColor: '#2B2D3E',
    strokeDashArray: 4,
  },
  plotOptions: {
    candlestick: {
      colors: {
        upward: '#00ff7f',
        downward: '#ff3b3b'
      },
      wick: { useFillColor: true }
    }
  },
  xaxis: {
    type: 'datetime' as const,
    tooltip: { enabled: true },
    axisBorder: { show: false },
    axisTicks: { show: false }
  },
  yaxis: {
    tooltip: { enabled: true },
    labels: {
      formatter: (val: number) => val.toFixed(2),
      style: { colors: '#9ca3af' }
    }
  }
}

// ----------------------------------------------------------------------
// CHART: Volume Oscillator
// ----------------------------------------------------------------------
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const volumeDataState = ref<any[]>(generateVolumeData(60))

const volumeSeries = computed(() => [{
  name: 'Volume',
  data: volumeDataState.value
}])

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const volumeOptions = {
  chart: {
    type: 'bar' as const,
    height: 150,
    background: 'transparent',
    toolbar: { show: false },
    animations: { enabled: false }
  },
  theme: { mode: 'dark' as const },
  grid: {
    borderColor: '#2B2D3E',
    strokeDashArray: 4,
  },
  colors: ['#3b82f6'],
  dataLabels: { enabled: false },
  xaxis: {
    type: 'datetime' as const,
    axisBorder: { show: false },
    axisTicks: { show: false },
    labels: { show: false } // Hide time on volume (sync with candle visually)
  },
  yaxis: {
    show: false,
  }
}

// ----------------------------------------------------------------------
// CHART: Treemap (Sector)
// ----------------------------------------------------------------------
const treemapSeries = [
  {
    data: [
      { x: 'PET Clear', y: 400 },
      { x: 'HDPE', y: 230 },
      { x: 'PP Injection', y: 150 },
      { x: 'LDPE Film', y: 110 },
      { x: 'Cardboard', y: 350 },
      { x: 'Office Paper', y: 100 },
      { x: 'Aluminium', y: 200 },
      { x: 'Steel', y: 280 },
      { x: 'Copper', y: 50 },
    ]
  }
]

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const treemapOptions = {
  legend: { show: false },
  chart: {
    type: 'treemap' as const,
    background: 'transparent',
    toolbar: { show: false }
  },
  theme: { mode: 'dark' as const },
  colors: ['#00ff7f', '#00cc66', '#00994d', '#ff3b3b', '#cc2f2f', '#3b82f6', '#2563eb'],
  plotOptions: {
    treemap: {
      distributed: true,
      enableShades: true
    }
  }
}

// ----------------------------------------------------------------------
// CHART: Heatmap (Correlation)
// ----------------------------------------------------------------------
const correlationValues = generateCorrelationMatrix(['Plastic', 'Paper', 'Metal', 'Oil', 'Gold'])
const heatmapSeries = computed(() => {
  // Group by Y to form the series format Apex expect
  const keys = ['Plastic', 'Paper', 'Metal', 'Oil', 'Gold']
  return keys.map(key => ({
    name: key,
    data: correlationValues.filter(v => v.y === key).map(v => ({ x: v.x, y: v.value }))
  }))
})

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const heatmapOptions = {
  chart: {
    type: 'heatmap' as const,
    background: 'transparent',
    toolbar: { show: false }
  },
  theme: { mode: 'dark' as const },
  dataLabels: { enabled: false },
  plotOptions: {
    heatmap: {
      colorScale: {
        ranges: [
          { from: -1, to: -0.5, color: '#ff3b3b', name: 'Negative' }, // Strong Negative
          { from: -0.5, to: 0.5, color: '#3b82f6', name: 'Neutral' },
          { from: 0.5, to: 1, color: '#00ff7f', name: 'Positive' }
        ]
      }
    }
  }
}

// ----------------------------------------------------------------------
// Logic
// ----------------------------------------------------------------------

// Fetch raw data helper
const getMarketData = async (name: string) => {
  try {
    const response = await fetch(`${config.webAPI}/api/receipts/market/history?name=${name}&interval=1d`)
    const res = await response.json()
    if (res.success && res.data && res.data.length > 0) {
      return res.data
    }
  } catch (e) {
    console.error("Error fetching market data for", name, e)
  }
  return null
}

const updateMainChart = async (name: string) => {
  const data = await getMarketData(name)

  if (data) {
    candleData.value = data
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    volumeDataState.value = data.map((d: any) => ({ x: d.x, y: d.volume }))

    const last = data[data.length - 1]
    currentPrice.value = last.y[3]

    currentOHLC.value = {
      open: last.y[0],
      high: last.y[1],
      low: last.y[2],
      close: last.y[3]
    }
  } else {
    // Fallback to simulation
    console.log("No real data for main chart", name)
    candleData.value = generateCandles(currentPrice.value || 10, 60)
    volumeDataState.value = generateVolumeData(60)
  }
}

const initIndices = async () => {
  for (const index of marketIndices.value) {
    const data = await getMarketData(index.query)
    if (data && data.length > 0) {
      // Get last candle
      const last = data[data.length - 1]
      const close = last.y[3]

      // Calculate change
      let change = 0
      if (data.length > 1) {
        const prev = data[data.length - 2]
        const prevClose = prev.y[3]
        change = ((close - prevClose) / prevClose) * 100
      }

      // Generate History (Sparkline)
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const history = data.slice(-10).map((d: any) => d.y[3])

      index.value = close
      index.change = change
      index.history = history
    } else {
      // Fallback to simulation so UI doesn't look broken (0.00)
      console.log("No data for index:", index.name, "using simulation")

      // Pick a base price roughly reasonable for the material
      let basePrice = 10
      if (index.name.includes('Metal')) basePrice = 30
      if (index.name.includes('Paper')) basePrice = 8
      if (index.name.includes('Plastic')) basePrice = 14
      if (index.name.includes('Glass')) basePrice = 4

      const simData = generateCandles(basePrice, 30) // Generate 30 candles
      const last = simData[simData.length - 1]
      const close = last.y[3]

      let change = 0
      if (simData.length > 1) {
        const prev = simData[simData.length - 2]
        const prevClose = prev.y[3]
        change = ((close - prevClose) / prevClose) * 100
      }

      const history = simData.slice(-10).map(d => d.y[3])

      index.value = close
      index.change = change
      index.history = history
    }
  }
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const selectIndex = (index: any) => {
  selectedAsset.value = index.name
  // Optimistically set price from index card, updated by real data in updateMainChart if available
  if (index.value > 0) {
    currentPrice.value = index.value
  }

  updateMainChart(index.query)
  orderBook.value = generateOrderBook(index.value || 12.50)
}

// Live Updates Simulation
let interval: ReturnType<typeof setInterval>

onMounted(async () => {
  // Fetch real waste items
  await wastesStore.fetchWastes(1, 100)

  // Init Indices
  await initIndices()

  // Fetch default market history
  updateMainChart('Plastic')

  setTimeout(() => {
    loadingInitial.value = false
  }, 1500)

  // Initial OHLC display
  if (candleData.value.length > 0) {
    const last = candleData.value[candleData.value.length - 1]
    currentOHLC.value = {
      open: last.y[0],
      high: last.y[1],
      low: last.y[2],
      close: last.y[3]
    }
  }

  interval = setInterval(() => {
    // 1. Update Order Book
    orderBook.value = generateOrderBook(currentPrice.value)

    // 2. Simulate Ticker movement (if needed)

    // 3. Update Last Candle (simulated real-time)
    if (candleData.value.length > 0) {
      const lastCandle = candleData.value[candleData.value.length - 1]
      const volatility = currentPrice.value * 0.005
      const movement = (Math.random() - 0.5) * volatility

      const newClose = lastCandle.y[3] + movement
      const newHigh = Math.max(lastCandle.y[1], newClose)
      const newLow = Math.min(lastCandle.y[2], newClose)

      const newCandle = {
        x: lastCandle.x,
        y: [lastCandle.y[0], newHigh, newLow, newClose] as [number, number, number, number]
      }

      // In Vue3 Apex, mutating deeply might not always trigger.
      // We replace the last element.
      const newData = [...candleData.value]
      newData[newData.length - 1] = newCandle
      candleData.value = newData

      currentPrice.value = newClose
      currentOHLC.value = {
        open: newCandle.y[0],
        high: newCandle.y[1],
        low: newCandle.y[2],
        close: newCandle.y[3]
      }
    }

  }, 1000)
})

onUnmounted(() => {
  clearInterval(interval)
})
</script>

<style scoped>
@keyframes marquee {
  0% {
    transform: translateX(0%);
  }

  100% {
    transform: translateX(-50%);
  }
}

.animate-marquee {
  animation: marquee 120s linear infinite;
}

/* Hide scrollbar for Chrome, Safari and Opera */
.custom-scrollbar::-webkit-scrollbar {
  display: none;
}

/* Hide scrollbar for IE, Edge and Firefox */
.custom-scrollbar {
  -ms-overflow-style: none;
  /* IE and Edge */
  scrollbar-width: none;
  /* Firefox */
}
</style>
