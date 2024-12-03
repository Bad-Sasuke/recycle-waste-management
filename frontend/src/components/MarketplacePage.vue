<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import Card from './Card.vue';
import { fetchData } from '../services/api';
import config from '../config';

const searchQuery = ref('');
const selectedCategory = ref<string[]>([]);
const webAPI = config.webAPI;

// const items = [
// { category: 'พลาสติก', price: 15, lastUpdate: '11/05/2567', hours: '09:00 - 18:00' },
// { category: 'โลหะ', price: 20, lastUpdate: '12/01/2567', hours: '08:00 - 17:00' },
// { category: 'กระดาษ', price: 8, lastUpdate: '11/11/2567', hours: '10:00 - 19:00' },
// { category: 'แก้ว', price: 12, lastUpdate: '11/03/2567', hours: '08: 30 - 18: 30' },
// เพิ่มรายการอื่นๆ ตามต้องการ
// ];

const items = ref<any[]>([]);

onMounted(async () => {
    const dataJson = await fetchData(webAPI + "/api/recycle-waste/get-wastes");
    if (!dataJson) {
        return;
    }
    dataJson.data.forEach((item: any) => {
        item.lastUpdate = new Date(item.lastUpdate).toLocaleDateString('th-TH', {
            day: 'numeric',
            month: 'numeric',
            year: 'numeric'
        });
    })
    items.value = dataJson.data;

});

const filteredItems = computed(() => {
    if (items.value.length === 0) {
        return [];
    }
    return items.value.filter(item => {
        const matchesSearch = item.category.toLowerCase().includes(searchQuery.value.toLowerCase());
        const matchesCategory =
            selectedCategory.value.length === 0 ||
            selectedCategory.value.includes(item.category.toLowerCase());
        return matchesSearch && matchesCategory;
    });
});

const toggleCategory = (category: string) => {
    const index = selectedCategory.value.indexOf(category);
    if (index === -1) {
        selectedCategory.value.push(category);
    } else {
        selectedCategory.value.splice(index, 1);
    }
};
</script>

<template>
    <div class="container px-4 py-6 grid grid-cols-1 lg:grid-cols-6 gap-4 max-w-full">
        <!-- ส่วนซ้าย : หมวดหมู่ -->
        <div class="bg-white shadow-md p-4 rounded lg:h-screen">
            <h2 class="text-2xl font-bold text-green-700 mb-4">หมวดหมู่</h2>
            <div>
                <label class="flex items-center space-x-2 mb-2">
                    <input type="checkbox" value="พลาสติก" @change="toggleCategory('พลาสติก')" />
                    <span>พลาสติก</span>
                </label>
                <label class="flex items-center space-x-2 mb-2">
                    <input type="checkbox" value="โลหะ" @change="toggleCategory('โลหะ')" />
                    <span>โลหะ</span>
                </label>
                <label class="flex items-center space-x-2 mb-2">
                    <input type="checkbox" value="กระดาษ" @change="toggleCategory('กระดาษ')" />
                    <span>กระดาษ</span>
                </label>
                <label class="flex items-center space-x-2 mb-2">
                    <input type="checkbox" value="แก้ว" @change="toggleCategory('แก้ว')" />
                    <span>แก้ว</span>
                </label>
            </div>
        </div>
        <!-- ส่วนขวา: รายการขยะ -->
        <div class="lg:col-span-5">
            <h1 class="text-3xl font-bold text-green-700 mb-4 text-center">ราคาขยะรีไซเคิล</h1>

            <div class="flex justify-center items-center gap-2 mb-6"> <!-- ฟอร์มค้นหาขยะ -->
                <div class="w-full">
                    <input type="text" v-model="searchQuery" class="input input-bordered w-full"
                        placeholder="ค้นหาขยะ เช่น พลาสติก, โลหะ" />
                </div>
            </div>

            <!-- รายการราคาขยะ -->
            <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                <div v-for="item in filteredItems" :key="item.category" class="flex justify-center">
                    <Card :name=item.category :price=item.price :category=item.category :lastUpdate=item.lastUpdate
                        :url=item.url />
                </div>
                <!-- <div v-for="item in filteredItems" :key="item.category" class="card w-full bg-base-200 shadow-md">
                    <div class="card-body">
                        <h2 class="card-title text-xl font-semibold text-green-700">{{ item.category }}</h2>
                        <p>ราคาต่อกิโลกรัม: ฿{{ item.price }}</p>
                        <p>สถานที่รับขยะ: {{ item.location }}</p>
                        <p>เวลาเปิดรับ: {{ item.hours }}</p>
                        <div class="card-actions justify-end">
                            <button class="btn btn-primary">ติดต่อ</button>
                        </div>
                    </div>
                </div> -->



            </div>
        </div>


    </div>
</template>

<style scoped></style>
