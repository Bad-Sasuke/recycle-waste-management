<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-80 p-4 backdrop-blur-sm"
    @click.self="$emit('close')">
    <div
      class="bg-white dark:bg-gray-900 rounded-xl max-w-md w-full overflow-hidden shadow-2xl flex flex-col max-h-[90vh]">

      <!-- Header / Controls -->
      <div class="p-4 flex justify-between items-center border-b border-gray-100 dark:border-gray-800">
        <h3 class="font-bold text-lg text-gray-800 dark:text-white flex items-center gap-2">
          Share Your Impact
          <IconWorld size="20" class="text-emerald-500" />
        </h3>
        <button @click="$emit('close')" class="btn btn-ghost btn-sm btn-circle">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Scrollable Preview Area -->
      <div class="overflow-y-auto flex-1 p-4 bg-gray-50 dark:bg-gray-800 flex justify-center">
        <!-- The Story Card (to be captured) -->
        <div ref="storyCard" id="impact-story-card"
          class="aspect-[9/16] w-full max-w-[320px] relative overflow-hidden rounded-2xl flex flex-col"
          style="background: linear-gradient(135deg, #10b981, #14b8a6, #0891b2); color: #ffffff; box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);">
          <!-- Decorative Background Elements -->
          <div class="absolute top-0 right-0 w-64 h-64 rounded-full -translate-y-1/2 translate-x-1/3"
            style="background-color: #ffffff; opacity: 0.05; filter: blur(64px);">
          </div>
          <div class="absolute bottom-0 left-0 w-64 h-64 rounded-full translate-y-1/3 -translate-x-1/3"
            style="background-color: #000000; opacity: 0.1; filter: blur(64px);">
          </div>

          <!-- Large Circle for Bottom Effect -->
          <div
            class="absolute bottom-[-10%] left-1/2 transform -translate-x-1/2 w-[120%] h-[50%] rounded-[100%] pointer-events-none"
            style="background-color: rgba(255, 255, 255, 0.05);">
          </div>

          <!-- Content -->
          <div class="relative z-10 flex flex-col h-full items-center text-center pb-32">

            <!-- Logo / Brand -->
            <div class="flex items-center gap-2 mt-8 px-6 py-2 rounded-full"
              style="background-color: rgba(255, 255, 255, 0.15); border: 1px solid rgba(255, 255, 255, 0.25);">
              <IconRecycle size="16" stroke="2" color="#ffffff" />
              <span class="font-bold tracking-[0.2em] text-xs uppercase" style="color: #ffffff;">Recycle Me</span>
            </div>

            <!-- Month -->
            <div class="mt-6 flex-shrink-0">
              <p class="uppercase tracking-[0.25em] text-[10px] font-bold mb-2" style="color: #d1fae5;">
                Monthly Report
              </p>
              <h2 class="text-4xl font-black font-serif italic" style="filter: drop-shadow(0 4px 4px rgba(0,0,0,0.1));">
                {{ currentMonth }}</h2>
              <div class="h-1.5 w-16 mx-auto mt-3 rounded-full" style="background-color: #6ee7b7;"></div>
            </div>

            <!-- Main Stat: Weight -->
            <div class="flex-1 flex flex-col justify-center items-center w-full relative z-20">
              <h1 class="text-[7rem] font-black leading-none mb-0"
                style="filter: drop-shadow(0 8px 16px rgba(0, 0, 0, 0.15)); line-height: 0.9;">
                {{ formatNumber(data?.total_weight_this_month || 0) }}
              </h1>
              <span class="text-2xl font-medium tracking-wide mt-4 block" style="color: #f0fdf4;">kg recycled</span>
            </div>

            <!-- Impact Footer -->
            <div class="absolute bottom-0 left-0 w-full pt-8 pb-8 px-6 rounded-t-[2.5rem] z-30"
              style="background: linear-gradient(to bottom, rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0.05)); border-top: 1px solid rgba(255, 255, 255, 0.2);">

              <div class="text-[10px] font-bold mb-2 uppercase tracking-[0.2em] opacity-90" style="color: #ecfdf5;">
                Equivalent to cleaning up
              </div>

              <div class="text-3xl font-black flex items-center justify-center gap-3" style="color: #ffffff;">
                <IconBottle size="32" stroke="2.5" color="#ffffff" />
                <span>{{ Math.round((data?.total_weight_this_month || 0) * 5) }} Waste Items</span>
              </div>
            </div>

          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="p-4 bg-white dark:bg-gray-900 border-t border-gray-100 dark:border-gray-800 flex flex-col gap-3">
        <button @click="generateImage" class="btn btn-primary w-full text-white shadow-lg shadow-primary/30 group"
          :disabled="generating">
          <span v-if="generating" class="loading loading-spinner"></span>
          <span v-else class="flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 group-hover:scale-110 transition-transform"
              fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            Download Story Image
          </span>
        </button>
        <p class="text-xs text-center text-gray-400">Save and share to your Instagram Story or Facebook!</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import html2canvas from 'html2canvas-pro';

import type { UserAnalyticsResponse } from '@/types/receipt';
import {
  IconWorld,
  IconRecycle,
  IconBottle
} from '@tabler/icons-vue';

// Props
defineProps<{
  data: UserAnalyticsResponse
}>();

defineEmits(['close']);

const storyCard = ref<HTMLElement | null>(null);
const generating = ref(false);

const currentMonth = computed(() => {
  return new Date().toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
});

const formatNumber = (num: number) => {
  if (num >= 1000) return (num / 1000).toFixed(1) + 'k';
  return Math.round(num); // Keep it clean for the story
};

const generateImage = async () => {
  if (!storyCard.value) return;

  generating.value = true;
  try {
    // Ensure fonts are loaded
    await document.fonts.ready;
    // Small delay to ensure rendering frames catch up
    await new Promise(resolve => setTimeout(resolve, 500));

    const canvas = await html2canvas(storyCard.value, {
      scale: 2,
      backgroundColor: null, // Transparent background
      logging: false,
      useCORS: true,
      allowTaint: true,
      imageTimeout: 15000,
    });

    const link = document.createElement('a');
    link.download = `my-impact-${new Date().toISOString().slice(0, 10)}.png`;
    link.href = canvas.toDataURL('image/png');
    link.click();
  } catch (e: unknown) {
    console.error('Failed to generate image:', e);
    const msg = e instanceof Error ? e.message : 'Unknown error';
    alert(`Failed to generate image: ${msg}`);
  } finally {
    generating.value = false;
  }
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Caveat:wght@600&family=Playfair+Display:ital,wght@1,700&display=swap');

.font-serif {
  font-family: 'Playfair Display', serif;
}

.font-handwriting {
  font-family: 'Caveat', cursive;
}
</style>
