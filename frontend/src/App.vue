<script setup lang="ts">
import {RouterView } from 'vue-router'
import FooterComponent from './components/FooterComponent.vue'
import NavbarHeader from './components/NavbarHeader.vue'
</script>

<template>
  <header>
    <NavbarHeader />
  </header>
  <RouterView />
  <FooterComponent />
</template>
