<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { setCookie } from '../stores/cookie'
import LandingPage from '../components/LandingPage.vue'

const route = useRoute()
const token = ref('')

onMounted(() => {
  token.value = route.query.token as string

  if (token.value) {
    setCookie('token', token.value, 0, 6)
  }
})
</script>

<template>
  <main>
    <LandingPage />
  </main>
</template>
