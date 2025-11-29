<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

export interface ReceiptDetailData {
    receipt: {
        id: string
        created_at: string
        payment_method: string
        total_amount: number
        vat_rate: number
        vat: number
        net_total: number
    }
    items: Array<{
        id: string
        name: string
        category: string
        weight: number
        unit_price: number
        price: number
    }>
    shop: {
        name: string
    }
}

export interface ReceiptDetailProps {
    isLoading?: boolean
    receiptData?: ReceiptDetailData | null
}

const _props = withDefaults(defineProps<ReceiptDetailProps>(), {
    isLoading: false,
    receiptData: null
})

const showModal = ref(false)

// Format currency
const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('th-TH', {
        style: 'currency',
        currency: 'THB',
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    }).format(value)
}

// Function to open the modal
const openModal = () => {
    showModal.value = true
}

// Function to close the modal
const closeModal = () => {
    showModal.value = false
}

// Handle Escape key to close modal
const handleEscKey = (e: KeyboardEvent) => {
    if (e.key === 'Escape' && showModal.value) {
        closeModal()
    }
}

// Close modal when clicking outside
const closeOnOutsideClick = (e: Event) => {
    if (e.target === e.currentTarget) {
        closeModal()
    }
}

onMounted(() => {
    document.addEventListener('keydown', handleEscKey)
})

onUnmounted(() => {
    document.removeEventListener('keydown', handleEscKey)
})

// Expose methods to parent components
defineExpose({
    openModal,
    closeModal
})
</script>

<template>
    <!-- Receipt Detail Modal -->
    <dialog id="receipt-detail-modal" class="modal modal-bottom sm:modal-middle" :class="{ 'modal-open': showModal }"
        @click="closeOnOutsideClick">
        <div class="modal-box max-w-4xl">
            <h3 class="font-bold text-2xl mb-4 flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                รายละเอียดใบเสร็จ
            </h3>

            <!-- Loading State -->
            <div v-if="isLoading" class="flex justify-center items-center py-12">
                <span class="loading loading-spinner loading-lg text-primary"></span>
            </div>

            <!-- Receipt Content -->
            <div v-else-if="receiptData" class="space-y-6">
                <!-- Receipt Header Info -->
                <div class="bg-base-200 rounded-lg p-4 space-y-2">
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <p class="text-sm text-gray-500">เลขที่ใบเสร็จ</p>
                            <p class="font-mono font-semibold">#{{ receiptData.receipt?.id || '-' }}</p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-500">วันที่ออกใบเสร็จ</p>
                            <p class="font-semibold">
                                {{
                                    receiptData.receipt?.created_at
                                        ? new Date(receiptData.receipt.created_at).toLocaleString('th-TH')
                                        : '-'
                                }}
                            </p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-500">ร้านค้า</p>
                            <p class="font-semibold">{{ receiptData.shop?.name || '-' }}</p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-500">วิธีชำระเงิน</p>
                            <p class="font-semibold">{{ receiptData.receipt?.payment_method || '-' }}</p>
                        </div>
                    </div>
                </div>

                <!-- Receipt Items -->
                <div>
                    <h4 class="font-bold text-lg mb-3">รายการขยะรีไซเคิล</h4>
                    <div class="overflow-x-auto">
                        <table class="table table-zebra w-full">
                            <thead>
                                <tr class="bg-base-200">
                                    <th>ชื่อวัสดุ</th>
                                    <th>หมวดหมู่</th>
                                    <th class="text-right">น้ำหนัก (กก.)</th>
                                    <th class="text-right">ราคา/กก.</th>
                                    <th class="text-right">รวม</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="item in receiptData.items" :key="item.id">
                                    <td class="font-semibold">{{ item.name }}</td>
                                    <td><span class="badge badge-primary badge-sm">{{ item.category }}</span></td>
                                    <td class="text-right">{{ item.weight }}</td>
                                    <td class="text-right">{{ formatCurrency(item.unit_price) }}</td>
                                    <td class="text-right font-bold">{{ formatCurrency(item.price) }}</td>
                                </tr>
                                <tr v-if="!receiptData.items || receiptData.items.length === 0">
                                    <td colspan="5" class="text-center text-gray-500 py-4">ไม่มีรายการ</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>

                <!-- Receipt Summary -->
                <div class="bg-gradient-to-br from-green-50 to-emerald-50 rounded-lg p-4 space-y-2">
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600">ยอดรวม</span>
                        <span class="font-semibold text-lg">{{
                            formatCurrency(receiptData.receipt?.total_amount || 0)
                            }}</span>
                    </div>
                    <div class="flex justify-between items-center">
                        <span class="text-gray-600">VAT ({{ ((receiptData.receipt?.vat_rate || 0) * 100).toFixed(0)
                            }}%)</span>
                        <span class="font-semibold text-lg">{{
                            formatCurrency(receiptData.receipt?.vat || 0)
                            }}</span>
                    </div>
                    <div class="divider my-2"></div>
                    <div class="flex justify-between items-center">
                        <span class="text-xl font-bold">ยอดสุทธิ</span>
                        <span class="text-2xl font-bold text-green-600">{{
                            formatCurrency(receiptData.receipt?.net_total || 0)
                            }}</span>
                    </div>
                </div>
            </div>

            <!-- Error State -->
            <div v-else class="text-center py-8">
                <svg xmlns="http://www.w3.org/2000/svg" class="mx-auto h-16 w-16 text-error" fill="none"
                    viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p class="mt-4 text-gray-600">ไม่สามารถโหลดข้อมูลใบเสร็จได้</p>
            </div>

            <div class="modal-action">
                <button class="btn" @click="closeModal">ปิด</button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button @click="closeModal">close</button>
        </form>
    </dialog>
</template>
