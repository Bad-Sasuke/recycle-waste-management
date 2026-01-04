<script setup lang="ts">
import { IconRecycle } from '@tabler/icons-vue'

defineProps({
    loading: {
        type: Boolean,
        default: false
    },
    text: {
        type: String,
        default: 'Loading...'
    },
    fullscreen: {
        type: Boolean,
        default: false
    }
})
</script>

<template>
    <transition enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0"
        enter-to-class="opacity-100" leave-active-class="transition duration-200 ease-in" leave-from-class="opacity-100"
        leave-to-class="opacity-0">
        <div v-if="loading" :class="[
            'flex flex-col items-center justify-center bg-gray-900/80 backdrop-blur-sm z-50',
            fullscreen ? 'fixed inset-0' : 'absolute inset-0'
        ]">
            <div class="relative w-24 h-24 mb-4">
                <div class="absolute inset-0 border-4 border-gray-700 rounded-full"></div>
                <div
                    class="absolute inset-0 border-4 border-t-green-500 border-r-transparent border-b-transparent border-l-transparent rounded-full animate-spin">
                </div>
                <div class="absolute inset-4 border-4 border-t-blue-500 border-r-transparent border-b-transparent border-l-transparent rounded-full animate-spin"
                    style="animation-direction: reverse; animation-duration: 1.5s;"></div>

                <div class="absolute inset-0 flex items-center justify-center">
                    <IconRecycle class="w-8 h-8 text-white animate-pulse" />
                </div>
            </div>
            <p
                class="text-xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-green-400 to-blue-500 animate-pulse">
                {{ text }}
            </p>
        </div>
    </transition>
</template>
