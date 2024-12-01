<script setup lang="ts">
import { ref, computed } from 'vue';

const searchQuery = ref('');
const selectedCategory = ref('');

const items = [
    { category: 'พลาสติก', price: 15, location: 'สถานที่ A', hours: '09:00 - 18:00' },
    { category: 'โลหะ', price: 20, location: 'สถานที่ B', hours: '08:00 - 17:00' },
    { category: 'กระดาษ', price: 8, location: 'สถานที่ C', hours: '10:00 - 19:00' },
    { category: 'แก้ว', price: 12, location: 'สถานที่ D', hours: '08:30 - 18:30' },
    // เพิ่มรายการอื่นๆ ตามต้องการ
];

const filteredItems = computed(() => {
    return items.filter(item => {
        const matchesSearch = item.category.toLowerCase().includes(searchQuery.value.toLowerCase());
        const matchesCategory = selectedCategory.value ? item.category.toLowerCase().includes(selectedCategory.value.toLowerCase()) : true;
        return matchesSearch && matchesCategory;
    });
});
</script>

<template>
    <div class="container mx-auto px-4 py-6">
        <h1 class="text-3xl font-bold text-green-700 mb-4 text-center">ราคาขยะรีไซเคิล</h1>

        <div class="flex justify-center items-center gap-2"> <!-- ฟอร์มค้นหาขยะ -->
            <div class="mb-6 flex flex-col items-center justify-center">
                <label class="text-lg font-semibold text-green-700">ค้นหา</label>
                <input type="text" v-model="searchQuery" class="input input-bordered w-full max-w-xs mt-2"
                    placeholder="ค้นหาขยะ เช่น พลาสติก, โลหะ" />
            </div>

            <!-- ฟิลเตอร์เลือกประเภทขยะ -->
            <div class="mb-6 flex flex-col items-center justify-center">
                <label class="text-lg font-semibold text-green-700">เลือกประเภทขยะ</label>
                <select v-model="selectedCategory" class="select select-bordered w-full max-w-xs mt-2">
                    <option disabled selected>เลือกประเภทขยะ</option>
                    <option value="plastic">พลาสติก</option>
                    <option value="metal">โลหะ</option>
                    <option value="paper">กระดาษ</option>
                    <option value="glass">แก้ว</option>
                </select>
            </div>
        </div>

        <!-- รายการราคาขยะ -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            <!-- รายการตัวอย่างที่กรองด้วย searchQuery และ selectedCategory -->
            <div v-for="item in filteredItems" :key="item.category" class="card w-full bg-base-200 shadow-md">
                <div class="card-body">
                    <h2 class="card-title text-xl font-semibold text-green-700">{{ item.category }}</h2>
                    <p>ราคาต่อกิโลกรัม: ฿{{ item.price }}</p>
                    <p>สถานที่รับขยะ: {{ item.location }}</p>
                    <p>เวลาเปิดรับ: {{ item.hours }}</p>
                    <div class="card-actions justify-end">
                        <button class="btn btn-primary">ติดต่อ</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<style scoped></style>