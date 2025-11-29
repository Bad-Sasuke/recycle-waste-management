<script setup lang="ts">
import { ref } from 'vue'
import { IconBrandGithub, IconBrandGoogleFilled } from '@tabler/icons-vue'
import { authLoginGithub, authLoginGoogle } from '@/services/auth'
import { employeeLogin } from '@/services/employee'
import { setCookie } from '@/stores/cookie'
import { useRouter } from 'vue-router'
import { useUsersStore } from '@/stores/users'

const router = useRouter()
const usersStore = useUsersStore()
const showEmployeeModal = ref(false)
const isLoggingIn = ref(false)
const employeeForm = ref({
  shop_code: '',
  username: '',
  password: ''
})
const showPassword = ref(false)
const errorMessage = ref('')

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

const openEmployeeLogin = () => {
  // Close the main login modal first
  const modal = document.getElementById('popup-login') as HTMLDialogElement
  modal.close()

  // Then open employee login modal
  showEmployeeModal.value = true
  errorMessage.value = ''
  employeeForm.value = {
    shop_code: '',
    username: '',
    password: ''
  }
}

const closeEmployeeModal = () => {
  showEmployeeModal.value = false
  errorMessage.value = ''
}

const handleEmployeeLogin = async () => {
  try {
    isLoggingIn.value = true
    errorMessage.value = ''

    const result = await employeeLogin(
      employeeForm.value.shop_code,
      employeeForm.value.username,
      employeeForm.value.password
    )

    // Save token
    setCookie('token', result.token, 7)

    // Save employee data to store
    usersStore.setEmployeeData(result.employee)

    // Close modals
    showEmployeeModal.value = false
    const modal = document.getElementById('popup-login') as HTMLDialogElement
    modal.close()

    // Re-check login status instead of reload
    await usersStore.checkLogin()

    // Redirect to home
    router.push({ name: 'home' })
  } catch (error: unknown) {
    console.error('Employee login error:', error)
    errorMessage.value = 'รหัสร้าน ชื่อผู้ใช้ หรือรหัสผ่านไม่ถูกต้อง'
  } finally {
    isLoggingIn.value = false
  }
}
</script>

<template>
  <dialog id="popup-login" class="modal">
    <div class="modal-box">
      <div class="flex flex-col gap-4">
        <h3 class="text-2xl font-bold text-center">{{ $t('PopupLogin.title') }}</h3>
        <div class="flex flex-col gap-2">
          <button class="btn btn-outline" @click="handleGithubLogin">
            <p class="mr-2">
              <IconBrandGithub stoke="1.5" />
            </p>
            {{ $t('PopupLogin.signInVia') }} Github
          </button>
          <button class="btn btn-outline" @click="handleGoogleLogin">
            <p class="mr-2">
              <IconBrandGoogleFilled stoke="1.5" />
            </p>
            {{ $t('PopupLogin.signInVia') }} Google
          </button>

          <!-- Employee Login Button -->
          <div class="divider">หรือ</div>
          <button class="btn btn-primary" @click="openEmployeeLogin">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
              <path
                d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z" />
            </svg>
            เข้าสู่ระบบพนักงาน
          </button>
        </div>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>

  <!-- Employee Login Modal -->
  <dialog class="modal modal-bottom sm:modal-middle" :class="{ 'modal-open': showEmployeeModal }">
    <div class="modal-box">
      <h3 class="font-bold text-lg mb-4">เข้าสู่ระบบพนักงาน</h3>
      <form @submit.prevent="handleEmployeeLogin">
        <div class="form-control mb-4">
          <label class="label">
            <span class="label-text">รหัสร้าน *</span>
          </label>
          <input type="text" v-model="employeeForm.shop_code" class="input input-bordered" required maxlength="12"
            placeholder="เช่น SHOP001" />
        </div>

        <div class="form-control mb-4">
          <label class="label">
            <span class="label-text">ชื่อผู้ใช้ *</span>
          </label>
          <input type="text" v-model="employeeForm.username" class="input input-bordered" required
            placeholder="Username" />
        </div>

        <div class="form-control mb-6">
          <label class="label">
            <span class="label-text">รหัสผ่าน *</span>
          </label>
          <div class="relative">
            <input :type="showPassword ? 'text' : 'password'" v-model="employeeForm.password"
              class="input input-bordered w-full pr-10" required minlength="6" placeholder="Password" />
            <button type="button" class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500"
              @click="showPassword = !showPassword">
              <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none"
                viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
                stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
        </div>

        <div v-if="errorMessage" class="alert alert-error mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{{ errorMessage }}</span>
        </div>

        <div class="modal-action">
          <button type="button" class="btn" @click="closeEmployeeModal" :disabled="isLoggingIn">ยกเลิก</button>
          <button type="submit" class="btn btn-primary" :disabled="isLoggingIn">
            <span v-if="isLoggingIn" class="loading loading-spinner"></span>
            {{ isLoggingIn ? 'กำลังเข้าสู่ระบบ...' : 'เข้าสู่ระบบ' }}
          </button>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button @click="closeEmployeeModal">close</button>
    </form>
  </dialog>
</template>
