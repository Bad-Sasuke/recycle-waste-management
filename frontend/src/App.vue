<script setup lang="ts">
import { RouterView } from 'vue-router'
import FooterComponent from './components/FooterComponent.vue'
import NavbarHeader from './components/NavbarHeader.vue'
</script>

<template>
  <div class="font-kanit">
    <header>
      <NavbarHeader />
    </header>
    <RouterView />
    <FooterComponent />
  </div>
</template>
