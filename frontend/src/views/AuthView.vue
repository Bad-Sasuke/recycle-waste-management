<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { setCookie } from '../stores/cookie'
import PopupLogin from '../components/PopupLogin.vue'

const route = useRoute()
const router = useRouter()
const token = ref('')

onMounted(() => {
  // Check if there's a token in the URL query params
  token.value = route.query.token as string

  if (token.value) {
    setCookie('token', token.value, 0, 6)
    
    // Check if there's a redirect parameter in the query
    const redirectPath = route.query.redirect as string
    if (redirectPath) {
      router.push(redirectPath)
    } else {
      router.push('/')
    }
  }
})
</script>

<template>
  <main>
    <PopupLogin />
  </main>
</template>
