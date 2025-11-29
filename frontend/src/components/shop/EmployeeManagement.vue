<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useShopStore } from '@/stores/shop';
import type { Employee, CreateEmployeeRequest, UpdateEmployeeRequest } from '@/types/employee';
import * as employeeService from '@/services/employee';

const shopStore = useShopStore();

// States
const employees = ref<Employee[]>([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(10);
const totalPages = ref(1);
const totalItems = ref(0);

// Modal States
const showModal = ref(false);
const showDeleteModal = ref(false);
const isEditing = ref(false);
const isSubmitting = ref(false);
const selectedEmployee = ref<Employee | null>(null);
const showPassword = ref(false);

// Form Data
const formData = ref<CreateEmployeeRequest & UpdateEmployeeRequest>({
    shop_id: '',
    first_name: '',
    last_name: '',
    username: '',
    password: ''
});

// Fetch Employees
const fetchEmployees = async () => {
    if (!shopStore.shop?.shop_id) return;

    loading.value = true;
    try {
        const response = await employeeService.getEmployeesByShopID(
            shopStore.shop.shop_id,
            page.value,
            pageSize.value
        );

        if (response.data) {
            employees.value = response.data;
            totalPages.value = response.total_pages;
            totalItems.value = response.total_items;
        }
    } catch (error) {
        console.error('Error fetching employees:', error);
    } finally {
        loading.value = false;
    }
};

// Open Create Modal
const openCreateModal = () => {
    isEditing.value = false;
    selectedEmployee.value = null;
    showPassword.value = false;
    formData.value = {
        shop_id: shopStore.shop?.shop_id || '',
        first_name: '',
        last_name: '',
        username: '',
        password: ''
    };
    showModal.value = true;
};

// Open Edit Modal
const openEditModal = (employee: Employee) => {
    isEditing.value = true;
    selectedEmployee.value = employee;
    showPassword.value = false;
    formData.value = {
        shop_id: employee.shop_id,
        first_name: employee.first_name,
        last_name: employee.last_name,
        username: employee.username,
        password: '' // Don't show password
    };
    showModal.value = true;
};

// Open Delete Modal
const openDeleteModal = (employee: Employee) => {
    selectedEmployee.value = employee;
    showDeleteModal.value = true;
};

// Submit Form
const handleSubmit = async () => {
    if (!shopStore.shop?.shop_id) return;

    isSubmitting.value = true;
    try {
        if (isEditing.value && selectedEmployee.value) {
            // Update
            const updateData: UpdateEmployeeRequest = {
                first_name: formData.value.first_name,
                last_name: formData.value.last_name,
                username: formData.value.username
            };

            if (formData.value.password) {
                updateData.password = formData.value.password;
            }

            await employeeService.updateEmployee(selectedEmployee.value.employee_id, updateData);
        } else {
            // Create
            const createData: CreateEmployeeRequest = {
                shop_id: shopStore.shop.shop_id,
                first_name: formData.value.first_name,
                last_name: formData.value.last_name,
                username: formData.value.username,
                password: formData.value.password || ''
            };

            await employeeService.createEmployee(createData);
        }

        showModal.value = false;
        fetchEmployees();
    } catch (error) {
        console.error('Error saving employee:', error);
        alert('เกิดข้อผิดพลาดในการบันทึกข้อมูล');
    } finally {
        isSubmitting.value = false;
    }
};

// Delete Employee
const handleDelete = async () => {
    if (!selectedEmployee.value) return;

    try {
        await employeeService.deleteEmployee(selectedEmployee.value.employee_id);
        showDeleteModal.value = false;
        fetchEmployees();
    } catch (error) {
        console.error('Error deleting employee:', error);
        alert('เกิดข้อผิดพลาดในการลบข้อมูล');
    }
};

// Watch page changes
watch(page, () => {
    fetchEmployees();
});

onMounted(() => {
    fetchEmployees();
});
</script>

<template>
    <div>
        <div class="flex justify-between items-center mb-6">
            <h2 class="text-xl font-bold">จัดการพนักงาน</h2>
            <button @click="openCreateModal" class="btn btn-primary btn-sm text-white">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd"
                        d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z"
                        clip-rule="evenodd" />
                </svg>
                เพิ่มพนักงาน
            </button>
        </div>

        <!-- Table -->
        <div class="overflow-x-auto bg-white rounded-lg shadow">
            <table class="table table-zebra w-full">
                <thead>
                    <tr class="bg-base-200">
                        <th>ชื่อผู้ใช้</th>
                        <th>ชื่อ-นามสกุล</th>
                        <th>วันที่สร้าง</th>
                        <th class="text-center">จัดการ</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="loading">
                        <td colspan="4" class="text-center py-8">
                            <span class="loading loading-spinner loading-md text-primary"></span>
                        </td>
                    </tr>
                    <tr v-else-if="employees.length === 0">
                        <td colspan="4" class="text-center py-8 text-gray-500">
                            ไม่พบข้อมูลพนักงาน
                        </td>
                    </tr>
                    <tr v-else v-for="emp in employees" :key="emp.employee_id">
                        <td class="font-medium">{{ emp.username }}</td>
                        <td>{{ emp.first_name }} {{ emp.last_name }}</td>
                        <td class="text-sm text-gray-500">
                            {{ new Date(emp.created_at).toLocaleDateString('th-TH') }}
                        </td>
                        <td class="text-center">
                            <div class="flex justify-center gap-2">
                                <button @click="openEditModal(emp)"
                                    class="btn btn-square btn-ghost btn-sm text-warning">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20"
                                        fill="currentColor">
                                        <path
                                            d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                                    </svg>
                                </button>
                                <button @click="openDeleteModal(emp)"
                                    class="btn btn-square btn-ghost btn-sm text-error">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20"
                                        fill="currentColor">
                                        <path fill-rule="evenodd"
                                            d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                                            clip-rule="evenodd" />
                                    </svg>
                                </button>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex justify-center mt-6">
            <div class="join">
                <button class="join-item btn btn-sm" :disabled="page === 1" @click="page--">«</button>
                <button v-for="p in totalPages" :key="p" class="join-item btn btn-sm"
                    :class="{ 'btn-active': page === p }" @click="page = p">
                    {{ p }}
                </button>
                <button class="join-item btn btn-sm" :disabled="page === totalPages" @click="page++">»</button>
            </div>
        </div>

        <!-- Create/Edit Modal -->
        <dialog class="modal modal-bottom sm:modal-middle" :class="{ 'modal-open': showModal }">
            <div class="modal-box">
                <h3 class="font-bold text-lg mb-4">{{ isEditing ? 'แก้ไขข้อมูลพนักงาน' : 'เพิ่มพนักงานใหม่' }}</h3>
                <form @submit.prevent="handleSubmit">
                    <div class="grid grid-cols-2 gap-4 mb-4">
                        <div class="form-control">
                            <label class="label"><span class="label-text">ชื่อ *</span></label>
                            <input type="text" v-model="formData.first_name" class="input input-bordered w-full"
                                required />
                        </div>
                        <div class="form-control">
                            <label class="label"><span class="label-text">นามสกุล *</span></label>
                            <input type="text" v-model="formData.last_name" class="input input-bordered w-full"
                                required />
                        </div>
                    </div>

                    <div class="form-control mb-4">
                        <label class="label"><span class="label-text">ชื่อผู้ใช้ (Username) *</span></label>
                        <input type="text" v-model="formData.username" class="input input-bordered w-full" required />
                    </div>

                    <div class="form-control mb-6">
                        <label class="label">
                            <span class="label-text">รหัสผ่าน {{ isEditing ? '(เว้นว่างหากไม่ต้องการเปลี่ยน)' : '*'
                                }}</span>
                        </label>
                        <div class="relative">
                            <input :type="showPassword ? 'text' : 'password'" v-model="formData.password"
                                class="input input-bordered w-full pr-10" :required="!isEditing" minlength="6" />
                            <button type="button"
                                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500"
                                @click="showPassword = !showPassword">
                                <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none"
                                    viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                </svg>
                                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none"
                                    viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                                </svg>
                            </button>
                        </div>
                        <label class="label" v-if="!isEditing">
                            <span class="label-text-alt text-gray-500">รหัสผ่านต้องมีความยาวอย่างน้อย 6 ตัวอักษร</span>
                        </label>
                    </div>

                    <div class="modal-action">
                        <button type="button" class="btn" @click="showModal = false"
                            :disabled="isSubmitting">ยกเลิก</button>
                        <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
                            <span v-if="isSubmitting" class="loading loading-spinner"></span>
                            บันทึก
                        </button>
                    </div>
                </form>
            </div>
            <form method="dialog" class="modal-backdrop">
                <button @click="showModal = false">close</button>
            </form>
        </dialog>

        <!-- Delete Confirmation Modal -->
        <dialog class="modal modal-bottom sm:modal-middle" :class="{ 'modal-open': showDeleteModal }">
            <div class="modal-box">
                <h3 class="font-bold text-lg text-error flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                        stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                    ยืนยันการลบ
                </h3>
                <p class="py-4">คุณแน่ใจหรือไม่ที่จะลบพนักงาน <b>{{ selectedEmployee?.first_name }} {{
                    selectedEmployee?.last_name }}</b>?</p>
                <div class="modal-action">
                    <button class="btn" @click="showDeleteModal = false">ยกเลิก</button>
                    <button class="btn btn-error text-white" @click="handleDelete">ลบพนักงาน</button>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop">
                <button @click="showDeleteModal = false">close</button>
            </form>
        </dialog>
    </div>
</template>
