<script setup lang="ts">
import { IconBrandGithub, IconBrandGoogleFilled } from '@tabler/icons-vue'
import { authLoginGithub, authLoginGoogle } from '@/services/auth'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const handleGithubLogin = async () => {
  // Close the modal before redirecting
  const modal = document.getElementById('popup-login') as HTMLDialogElement
  modal.close()
  await authLoginGithub()
}

const handleGoogleLogin = async () => {
  // Close the modal before redirecting
  const modal = document.getElementById('popup-login') as HTMLDialogElement
  modal.close()
  await authLoginGoogle()
}
</script>

<template>
  <dialog id="popup-login" class="modal">
    <div class="modal-box">
      <div class="flex flex-col gap-4">
        <h3 class="text-2xl font-bold text-center">{{ $t('PopupLogin.title') }}</h3>
        <div class="flex flex-col gap-2">
          <button class="btn btn-outline" @click="handleGithubLogin">
            <p class="mr-2"><IconBrandGithub stoke="1.5" /></p>
            {{ $t('PopupLogin.signInVia') }} Github
          </button>
          <button class="btn btn-outline" @click="handleGoogleLogin">
            <p class="mr-2"><IconBrandGoogleFilled stoke="1.5" /></p>
            {{ $t('PopupLogin.signInVia') }} Google
          </button>
        </div>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
</template>
