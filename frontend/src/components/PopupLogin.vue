<script setup lang="ts">
import { IconBrandGithub, IconBrandGoogleFilled } from '@tabler/icons-vue'
import { authLoginGithub,authLoginGoogle } from '@/services/auth'
</script>

<template>
  <dialog id="popup-login" class="modal">
    <div class="modal-box">
      <div class="flex flex-col gap-4">
        <h3 class="text-2xl font-bold text-center">{{ $t('PopupLogin.title') }}</h3>
        <div class="flex flex-col gap-2">
          <button class="btn btn-outline" @click="authLoginGithub">
            <p class="mr-2"><IconBrandGithub stoke="1.5" /></p>
            {{ $t('PopupLogin.signInVia') }} Github
          </button>
          <button class="btn btn-outline" @click="authLoginGoogle">
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
