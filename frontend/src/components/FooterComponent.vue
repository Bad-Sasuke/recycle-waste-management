<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router'
import { computed } from 'vue'

const route = useRoute()
const isHomePage = computed(() => route.path === '/')
</script>

<template>
  <footer :class="[
    'py-8 mt-auto transition-colors duration-300',
    isHomePage ? 'bg-slate-900 text-slate-400 border-t border-slate-800' : 'bg-footer text-footer'
  ]">
    <div class="container mx-auto px-6 text-center">
      <p class="mb-4">&copy; 2024 - {{ new Date().getFullYear() }} RecycleMarket | {{ $t('Footer.copyright') }}</p>
      <div class="flex justify-center gap-6">
        <RouterLink to="/about?tab=about" class="hover:underline transition-colors"
          :class="isHomePage ? 'hover:text-white' : 'hover:text-footer-hover'">{{
            $t('Footer.about_us')
          }}</RouterLink>
        <RouterLink to="/about?tab=privacy" class="hover:underline transition-colors"
          :class="isHomePage ? 'hover:text-white' : 'hover:text-footer-hover'">{{ $t('Footer.policy')
          }}</RouterLink>
        <RouterLink to="/about?tab=terms" class="hover:underline transition-colors"
          :class="isHomePage ? 'hover:text-white' : 'hover:text-footer-hover'">{{ $t('Footer.terms') }}
        </RouterLink>
      </div>
    </div>
  </footer>
</template>
