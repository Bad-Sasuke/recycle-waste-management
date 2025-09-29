<script setup lang="ts"></script>

<template>
  <footer class="bg-[#c8edc5] text-[#0c270c] py-4 mt-auto">
    <div class="text-center">
      <p>&copy; 2024 - {{ new Date().getFullYear() }} RecycleMarket | {{ $t('Footer.copyright') }}</p>
      <div class="mt-2">
        <a href="#" class="hover:underline px-2 hover:text-[#328c2d]">{{
          $t('Footer.about_us')
          }}</a>
        <a href="#" class="hover:underline px-2 hover:text-[#328c2d]">{{ $t('Footer.policy') }}</a>
        <a href="#" class="hover:underline px-2 hover:text-[#328c2d]">{{ $t('Footer.terms') }}</a>
      </div>
    </div>
  </footer>
</template>
