<script setup lang="ts">
const props = defineProps({
    name: { type: String },
    price: { type: Number },
    category: { type: String },
    lastUpdate: { type: String },
    url: { type: String }
})
</script>

<template>
    <div class="card w-full md:w-64 shadow-2xl border-neutral/30 rounded-lg border-[0.25px] cursor-pointer h-96">
        <figure><img class="h-64 w-full object-cover bg-[#f2f2f2]" :src="props.url"></figure>
        <div class="card-body p-4">
            <div class="flex flex-col items-between justify-between h-full gap-3">

                <h2 class="card-title text-base-content text-md">
                    {{ props.name }}
                    <!-- <div class="badge badge-secondary">NEW</div> -->
                </h2>
                <div class="flex flex-col items-between gap-1 w-full">
                    <h2 class="text-md">
                        {{ props.price }} บาท/กก.
                    </h2>
                    <div class="badge badge-warning badge-sm px-2">{{ props.category }}</div>
                </div>
            </div>
        </div>
        <div class="bg-green-700 w-full text-center text-white rounded-b-lg text-xs p-1">อัปเดตล่าสุด:
            {{ props.lastUpdate }}</div>
    </div>
</template>