<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  IconCamera,
  IconPhoto,
  IconX,
  IconRefresh,
  IconBolt,
  IconCurrencyBaht,
  IconCube,
} from '@tabler/icons-vue'
import { wasteClassifier, type WastePrediction } from '@/services/wasteClassifier'
import { useWastesStore } from '@/stores/wastes'

const router = useRouter()
const wastesStore = useWastesStore()

// State
const videoRef = ref<HTMLVideoElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)
const isCameraOpen = ref(false)
const stream = ref<MediaStream | null>(null)
const selectedImage = ref<string | null>(null)
const isAnalyzing = ref(false)
const predictions = ref<WastePrediction[]>([])
const error = ref<string | null>(null)
const showResult = ref(false)

// Wastes data
onMounted(async () => {
  // Pre-load model
  try {
    await wasteClassifier.loadModel()
    await wastesStore.fetchWastes() // Fetch wastes for pricing lookup
  } catch (e) {
    console.error('Failed to init:', e)
    error.value = 'Failed to initialize AI model'
  }
})

onUnmounted(() => {
  stopCamera()
})

// Camera Functions
const startCamera = async () => {
  try {
    error.value = null
    stream.value = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'environment' },
    })
    if (videoRef.value) {
      videoRef.value.srcObject = stream.value
      isCameraOpen.value = true
    }
  } catch (e) {
    console.error('Camera error:', e)
    error.value = 'Could not access camera. Please allow permissions.'
  }
}

const stopCamera = () => {
  if (stream.value) {
    stream.value.getTracks().forEach((track) => track.stop())
    stream.value = null
  }
  isCameraOpen.value = false
}

const captureImage = () => {
  if (!videoRef.value || !canvasRef.value) return

  const video = videoRef.value
  const canvas = canvasRef.value

  // Set canvas dimensions to match video
  canvas.width = video.videoWidth
  canvas.height = video.videoHeight

  const ctx = canvas.getContext('2d')
  if (ctx) {
    ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
    selectedImage.value = canvas.toDataURL('image/jpeg')
    stopCamera()
    analyzeImage(canvas)
  }
}

// File Upload
const handleFileUpload = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      selectedImage.value = e.target?.result as string
      // Create an image element to analyze
      const img = new Image()
      img.onload = () => {
        analyzeImage(img)
      }
      img.src = selectedImage.value
    }
    reader.readAsDataURL(file)
  }
}

// AI Analysis
const analyzeImage = async (element: HTMLImageElement | HTMLVideoElement | HTMLCanvasElement) => {
  isAnalyzing.value = true
  showResult.value = false
  predictions.value = []

  try {
    // Artificial delay for UX (to show analyzing animation)
    await new Promise(resolve => setTimeout(resolve, 1500))

    const results = await wasteClassifier.classify(element)
    predictions.value = results
    showResult.value = true
  } catch (e) {
    console.error('Analysis error:', e)
    error.value = 'Failed to analyze image'
  } finally {
    isAnalyzing.value = false
  }
}

const reset = () => {
  selectedImage.value = null
  showResult.value = false
  predictions.value = []
  error.value = null
  stopCamera()
}

// Pricing Logic
const getEstimatedPrice = (category?: string) => {
  if (!category) return 'N/A'

  // Find wastes in this category
  const categoryWastes = wastesStore.wastes.filter(
    w => (w.category && w.category.toLowerCase().includes(category.toLowerCase())) ||
      (w.category && category.toLowerCase().includes(w.category.toLowerCase()))
  )

  if (categoryWastes.length === 0) return 'Price depends on market'

  // Calculate average or range
  const prices = categoryWastes
    .map(w => w.price)
    .filter((p): p is number => p !== undefined)

  if (prices.length === 0) return 'Price not available'

  const min = Math.min(...prices)
  const max = Math.max(...prices)

  if (min === max) return `${min} THB/kg`
  return `${min} - ${max} THB/kg`
}

const getCategoryColor = (category?: string) => {
  switch (category?.toLowerCase()) {
    case 'plastic': return 'bg-blue-100 text-blue-800 border-blue-200'
    case 'glass': return 'bg-green-100 text-green-800 border-green-200'
    case 'paper': return 'bg-yellow-100 text-yellow-800 border-yellow-200'
    case 'metal': return 'bg-gray-100 text-gray-800 border-gray-200'
    default: return 'bg-purple-100 text-purple-800 border-purple-200'
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col items-center">
    <!-- Header -->
    <div class="w-full bg-white/80 backdrop-blur-md sticky top-0 z-10 border-b border-gray-100">
      <div class="max-w-md mx-auto px-4 py-4 flex items-center justify-between">
        <button @click="router.back()" class="p-2 rounded-full hover:bg-gray-100 transition-colors">
          <IconX class="w-6 h-6 text-gray-600" />
        </button>
        <h1 class="text-lg font-bold bg-gradient-to-r from-green-600 to-teal-500 bg-clip-text text-transparent">
          AI Waste Scanner
        </h1>
        <div class="w-10"></div> <!-- Spacer -->
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 w-full max-w-md p-4 flex flex-col gap-6 relative">

      <!-- Camera/Image Area -->
      <div
        class="relative aspect-[3/4] bg-gray-900 rounded-3xl overflow-hidden shadow-2xl border-4 border-white ring-1 ring-gray-200">
        <!-- hidden canvas for capture -->
        <canvas ref="canvasRef" class="hidden"></canvas>

        <!-- Default State -->
        <div v-if="!isCameraOpen && !selectedImage"
          class="absolute inset-0 flex flex-col items-center justify-center text-white space-y-4">
          <div
            class="w-20 h-20 rounded-full bg-white/10 backdrop-blur-sm flex items-center justify-center animate-pulse">
            <IconCamera class="w-10 h-10" />
          </div>
          <p class="text-gray-300 font-medium">Take a photo or upload</p>
        </div>

        <!-- Video Feed -->
        <video v-show="isCameraOpen && !selectedImage" ref="videoRef" autoplay playsinline
          class="absolute inset-0 w-full h-full object-cover"></video>

        <!-- Captured/Selected Image -->
        <img v-if="selectedImage" :src="selectedImage" class="absolute inset-0 w-full h-full object-cover" />

        <!-- Scanner Overlay Animation -->
        <div v-if="isAnalyzing" class="absolute inset-0 bg-black/30 z-20">
          <div
            class="absolute top-0 left-0 w-full h-1 bg-green-400 shadow-[0_0_15px_rgba(74,222,128,0.8)] animate-scan">
          </div>
          <div class="absolute inset-0 flex items-center justify-center">
            <div class="bg-black/60 backdrop-blur-md px-6 py-3 rounded-full flex items-center gap-3">
              <IconRefresh class="w-5 h-5 text-green-400 animate-spin" />
              <span class="text-white font-medium tracking-wide">Analyzing waste...</span>
            </div>
          </div>
        </div>

        <!-- Camera Controls -->
        <div v-if="!selectedImage && !isAnalyzing"
          class="absolute bottom-6 inset-x-0 flex items-center justify-center gap-8 z-20">
          <!-- Upload Button -->
          <label class="cursor-pointer group">
            <input type="file" accept="image/*" class="hidden" @change="handleFileUpload" />
            <div
              class="p-4 rounded-full bg-white/20 backdrop-blur-md hover:bg-white/30 transition-all border border-white/20 group-active:scale-95">
              <IconPhoto class="w-6 h-6 text-white" />
            </div>
          </label>

          <!-- Capture Button -->
          <button @click="isCameraOpen ? captureImage() : startCamera()"
            class="p-1 rounded-full border-4 border-white/50 hover:border-white transition-all active:scale-95">
            <div class="w-16 h-16 rounded-full bg-white flex items-center justify-center">
              <div class="w-14 h-14 rounded-full border-2 border-gray-100 bg-gradient-to-tr from-green-500 to-teal-500">
              </div>
            </div>
          </button>

          <!-- Spacer for balance -->
          <div class="w-14"></div>
        </div>

        <!-- Retake Button -->
        <button v-if="selectedImage && !isAnalyzing" @click="reset"
          class="absolute top-4 right-4 p-2 rounded-full bg-black/50 text-white backdrop-blur-md hover:bg-black/70 transition-all z-20">
          <IconRefresh class="w-5 h-5" />
        </button>
      </div>

      <!-- Results Section -->
      <transition enter-active-class="transition ease-out duration-300"
        enter-from-class="transform translate-y-10 opacity-0" enter-to-class="transform translate-y-0 opacity-100"
        leave-active-class="transition ease-in duration-200" leave-from-class="transform translate-y-0 opacity-100"
        leave-to-class="transform translate-y-10 opacity-0">
        <div v-if="showResult && predictions.length > 0" class="flex flex-col gap-4">
          <!-- Main Prediction Card -->
          <div class="bg-white rounded-3xl p-6 shadow-xl border border-gray-100 relative overflow-hidden">
            <div class="absolute top-0 right-0 p-4 opacity-10">
              <IconBolt class="w-24 h-24 text-green-600" />
            </div>

            <div class="relative z-10">
              <h2 class="text-sm uppercase tracking-wider text-gray-500 font-semibold mb-1">Detected Object</h2>
              <div class="flex items-start justify-between">
                <div>
                  <h3 class="text-3xl font-bold text-gray-800 capitalize mb-2">
                    {{ predictions[0].className }}
                  </h3>
                  <div class="flex items-center gap-2">
                    <span class="px-3 py-1 rounded-full text-xs font-bold border"
                      :class="getCategoryColor(predictions[0].predictedCategory || 'Unknown')">
                      {{ predictions[0].predictedCategory || 'Uncategorized' }}
                    </span>
                    <span class="text-xs text-gray-400">
                      {{ Math.round(predictions[0].probability * 100) }}% Confidence
                    </span>
                  </div>
                </div>
              </div>

              <!-- Estimated Price -->
              <div class="mt-6 p-4 bg-gray-50 rounded-2xl border border-gray-100">
                <div class="flex items-center gap-3 mb-1">
                  <IconCurrencyBaht class="w-5 h-5 text-green-600" />
                  <span class="text-sm font-medium text-gray-600">Estimated Price</span>
                </div>
                <p class="text-2xl font-bold text-gray-800 ml-8">
                  {{ getEstimatedPrice(predictions[0].predictedCategory) }}
                </p>
                <div v-if="predictions[0].predictedCategory" class="mt-2 ml-8 text-xs text-gray-400">
                  *Based on average {{ predictions[0].predictedCategory }} prices
                </div>
              </div>
            </div>
          </div>

          <!-- Other Possibilities -->
          <div v-if="predictions.length > 1"
            class="bg-white/60 backdrop-blur-sm rounded-2xl p-4 border border-gray-200">
            <h4 class="text-sm font-medium text-gray-500 mb-3 ml-1">Other possibilities</h4>
            <div class="space-y-2">
              <div v-for="(pred, idx) in predictions.slice(1, 4)" :key="idx"
                class="flex items-center justify-between p-2 hover:bg-white rounded-lg transition-colors">
                <div class="flex items-center gap-3">
                  <IconCube class="w-4 h-4 text-gray-400" />
                  <span class="text-gray-700 capitalize">{{ pred.className }}</span>
                </div>
                <span class="text-xs font-medium text-gray-400">
                  {{ Math.round(pred.probability * 100) }}%
                </span>
              </div>
            </div>
          </div>
        </div>
      </transition>

      <!-- Error Message -->
      <div v-if="error" class="p-4 bg-red-50 text-red-600 rounded-xl border border-red-100 text-sm text-center">
        {{ error }}
      </div>

    </div>
  </div>
</template>

<style scoped>
@keyframes scan {
  0% {
    top: 0%;
    opacity: 0;
  }

  10% {
    opacity: 1;
  }

  90% {
    opacity: 1;
  }

  100% {
    top: 100%;
    opacity: 0;
  }
}

.animate-scan {
  animation: scan 2s linear infinite;
}
</style>
