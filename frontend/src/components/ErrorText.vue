<script setup lang="ts">
const props = defineProps({
  message: { type: String },
})
</script>

<template>
  <div class="label p-0 px-2" v-if="props.message">
    <span class="label-text-alt text-red-500">{{ props.message }}</span>
  </div>
</template>
