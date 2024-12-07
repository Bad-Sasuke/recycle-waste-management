<script setup lang="ts">
import { useWastesStore } from '../stores/wastes'
import { ref } from 'vue'

const wastesStore = useWastesStore()
const isDelete = ref(false)
const props = defineProps({
  id: { type: String },
  name: { type: String },
  price: { type: Number },
  category: { type: String },
  lastUpdate: { type: String },
  url: { type: String },
})

const deleteItem = async (id: string) => {
  isDelete.value = true
  await new Promise((resolve) => setTimeout(resolve, 1000))
  await wastesStore.deleteWaste(id)
  isDelete.value = false
}
</script>

<template>
  <div
    class="flex card w-full md:w-64 shadow-2xl border-neutral/30 rounded-lg h-96 flex-col gap-4"
    v-if="isDelete"
  >
    <div class="skeleton h-60 w-full object-cover bg-[#f2f2f2] rounded-lg"></div>
    <div class="mx-4 skeleton h-4 w-28"></div>
    <div class="mx-4 skeleton h-4"></div>
    <div class="mx-4 skeleton h-4 w-16"></div>
  </div>
  <div
    class="card w-full md:w-64 shadow-2xl border-neutral/30 rounded-lg border-[0.25px] cursor-pointer h-96"
    v-else
  >
    <figure>
      <img class="h-64 w-full object-cover bg-[#f2f2f2]" :src="props.url" />
    </figure>
    <div class="flex items-end justify-end absolute top-0 right-0">
      <button
        class="btn btn-sm btn-ghost mr-2 mt-2 hover:bg-base-600/50 hover:text-white/50"
        @click="deleteItem(props?.id ?? '')"
      >
        X
      </button>
    </div>
    <div class="card-body p-4" :data-id="props.id">
      <div class="flex flex-col items-between justify-between h-full gap-3">
        <h2 class="card-title text-base-content text-md">
          {{ props.name }}
          <!-- <div class="badge badge-secondary">NEW</div> -->
        </h2>
        <div class="flex flex-col items-between gap-1 w-full">
          <h2 class="text-md">{{ props.price }} บาท/กก.</h2>
          <div
            class="badge badge-outline border-green-700 text-green-700 hover:bg-green-700 hover:text-white badge-sm px-2"
          >
            {{ props.category }}
          </div>
        </div>
      </div>
    </div>
    <div class="bg-green-700 w-full text-center text-white rounded-b-lg text-xs p-1">
      อัปเดตล่าสุด: {{ props.lastUpdate }}
    </div>
  </div>
</template>
