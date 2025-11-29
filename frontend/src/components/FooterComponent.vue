<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <footer class="bg-footer text-footer py-4 mt-auto">
    <div class="text-center">
      <p>&copy; 2024 - {{ new Date().getFullYear() }} RecycleMarket | {{ $t('Footer.copyright') }}</p>
      <div class="mt-2">
        <RouterLink to="/about?tab=about" class="hover:underline px-2 hover:text-footer-hover">{{
          $t('Footer.about_us')
        }}</RouterLink>
        <RouterLink to="/about?tab=privacy" class="hover:underline px-2 hover:text-footer-hover">{{ $t('Footer.policy')
          }}</RouterLink>
        <RouterLink to="/about?tab=terms" class="hover:underline px-2 hover:text-footer-hover">{{ $t('Footer.terms') }}
        </RouterLink>
      </div>
    </div>
  </footer>
</template>
