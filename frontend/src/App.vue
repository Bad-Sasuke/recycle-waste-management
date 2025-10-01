<script setup lang="ts">
import { RouterView } from 'vue-router'
import FooterComponent from './components/FooterComponent.vue'
import NavbarHeader from './components/NavbarHeader.vue'
import ShopNotification from './components/ShopNotification.vue'
</script>

<template>
  <div class="font-kanit flex flex-col min-h-screen">
    <header>
      <NavbarHeader />
    </header>
    <RouterView />
    <FooterComponent />
    <ShopNotification />
  </div>
</template>
