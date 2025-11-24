<script setup lang="ts">
import { ref } from 'vue';

const dialog = ref<HTMLDialogElement | null>(null);
const title = ref('');
const message = ref('');
const type = ref<'success' | 'error' | 'info' | 'warning'>('info');

const show = (msg: string, t: string = 'Alert', msgType: 'success' | 'error' | 'info' | 'warning' = 'info') => {
  message.value = msg;
  title.value = t;
  type.value = msgType;
  dialog.value?.showModal();
};

const close = () => {
  dialog.value?.close();
};

defineExpose({ show, close });
</script>

<template>
  <dialog ref="dialog" class="modal">
    <div class="modal-box">
      <h3 class="font-bold text-lg" :class="{
        'text-error': type === 'error',
        'text-success': type === 'success',
        'text-info': type === 'info',
        'text-warning': type === 'warning'
      }">
        {{ title }}
      </h3>
      <p class="py-4 whitespace-pre-wrap">{{ message }}</p>
      <div class="modal-action">
        <form method="dialog">
          <button class="btn" @click="close">{{ $t('close') }}</button>
        </form>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button @click="close">close</button>
    </form>
  </dialog>
</template>
