<script setup lang="ts">
import { IconUserCircle, IconBell, IconMenu, IconTrash, IconBellOff } from '@tabler/icons-vue'
import { RouterLink } from 'vue-router'
import { ref, onBeforeMount, onMounted, watch } from 'vue'
import { useI18nStore } from '@/stores/i18n'
import { useUsersStore } from '@/stores/users'
import PopupLogin from './PopupLogin.vue'
import SwitchLangDesktop from './switchLang/SwitchLangDesktop.vue'
import SwitchLangMobile from './switchLang/SwitchLangMobile.vue'

const isDrawerOpen = ref(false)

const i18nStore = useI18nStore()
const usersStore = useUsersStore()

const closeDrawer = () => {
  isDrawerOpen.value = !isDrawerOpen.value
}

onBeforeMount(() => {
  i18nStore.getLocale()
})

const cookies = ref(document.cookie)

watch(cookies, (newCookie, oldCookie) => {
  if (newCookie !== oldCookie) {
    usersStore.checkLogin()
  }
})

onMounted(() => {
  usersStore.checkLogin()
  setInterval(() => {
    const currentCookies = document.cookie
    if (cookies.value !== currentCookies) {
      cookies.value = currentCookies
    }
  }, 500)
})

const handleLogin = () => {
  const modal = document.getElementById('popup-login') as HTMLDialogElement
  modal?.showModal()
  console.log(usersStore.isLogin)
}
const handleLogout = () => {
  usersStore.logout()
}
</script>

<template>
  <PopupLogin />
  <div class="navbar bg-green-100 shadow-md">
    <!-- เมนูสำหรับ Mobile -->
    <div class="flex-none md:hidden">
      <label for="menu-mobile" class="btn btn-ghost btn-circle">
        <IconMenu stroke="1.5" size="24" />
      </label>
    </div>

    <!-- โลโก้และชื่อเว็บ -->
    <div class="flex-1">
      <RouterLink to="/" class="btn btn-ghost normal-case text-2xl font-bold text-green-700">
        ♻️ Recycle Waste
      </RouterLink>
      <!-- เมนูหลัก -->
      <div class="hidden md:flex">
        <ul class="menu menu-horizontal px-1">
          <li v-if="usersStore.isLogin">
            <RouterLink to="/marketplace" class="font-semibold hover:text-green-600">{{
              $t('Navbar.menu.marketplace')
            }}</RouterLink>
          </li>
        </ul>
      </div>
    </div>

    <div class="flex-none gap-4">
      <!-- ปุ่มแจ้งเตือน -->
      <div class="dropdown dropdown-end">
        <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
          <div class="rounded-full">
            <IconBell stroke="1.5" size="24" />
          </div>
        </div>
        <ul
          tabindex="0"
          class="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-72 p-2 shadow h-64 flex flex-col"
        >
          <span class="flex justify-between mx-2 py-4">
            <p class="text-[1rem]">{{ $t('Navbar.notif.title') }}</p>
            <div class="flex justify-between items-center gap-2">
              <li>
                <a class="text-xs">{{ $t('Navbar.notif.read_all') }}</a>
              </li>
              <button class="hover:text-error">
                <IconTrash stroke="1.5" size="24" />
              </button>
            </div>
          </span>
          <div class="flex flex-col items-center justify-center flex-grow gap-2">
            <IconBellOff stroke="1.5" size="24" class="text-gray-500" />
            <p class="text-sm text-gray-500">{{ $t('Navbar.notif.empty') }}</p>
          </div>
          <div class="flex flex-col overflow-auto h-38">
            <!-- <li>
              <a>ไปหน้ากล่องข้อความ</a>
            </li>
            <li>
              <a>ไปหน้ากล่องข้อความ</a>
            </li>
            <li>
              <a>ไปหน้ากล่องข้อความ</a>
            </li>
            <li>
              <a>ไปหน้ากล่องข้อความ</a>
            </li>
            <li>
              <a>ไปหน้ากล่องข้อความ</a>
            </li> -->
          </div>
          <li class="flex items-center justify-center w-full mt-auto">
            <a>{{ $t('Navbar.notif.view_all') }}</a>
          </li>
        </ul>
      </div>

      <SwitchLangDesktop />

      <!-- โปรไฟล์ผู้ใช้งาน -->
      <div class="dropdown dropdown-end" @click="usersStore.isLogin !== true ? handleLogin() : ''">
        <div tabindex="0" role="button" class="">
          <IconUserCircle stroke="1.3" size="32" class="bg-neutral rounded-full text-base-100" />
        </div>
        <ul
          tabindex="0"
          class="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow"
          v-if="usersStore.isLogin"
        >
          <span v-if="usersStore.isLogin">
            <li>
              <RouterLink to="/profile" class="justify-between">
                {{ $t('Navbar.profile.textProfile') }}
              </RouterLink>
            </li>

            <li class="disabled">
              <a>
                {{ $t('Navbar.profile.textSettings') }}
                <span class="badge">{{ $t('Global.comingSoon') }} </span>
              </a>
            </li>
            <li>
              <a @click="handleLogout">{{ $t('Navbar.profile.textLogout') }}</a>
            </li>
          </span>
        </ul>
      </div>
    </div>
  </div>

  <!-- เมนูสำหรับ Mobile -->
  <div class="drawer z-50">
    <input id="menu-mobile" type="checkbox" class="drawer-toggle" v-model="isDrawerOpen" />
    <div class="drawer-side">
      <label for="menu-mobile" aria-label="close sidebar" class="drawer-overlay"></label>
      <ul class="menu bg-base-200 text-base-content min-h-full w-80 p-4">
        <p class="normal-case text-2xl font-bold text-green-700 m-2">♻️ Recycle Waste</p>
        <li>
          <RouterLink to="/" @click="closeDrawer">{{ $t('Navbar.menu.home') }}</RouterLink>
        </li>
        <li v-if="usersStore.isLogin">
          <RouterLink to="/marketplace" @click="closeDrawer">{{
            $t('Navbar.menu.marketplace')
          }}</RouterLink>
        </li>
        <SwitchLangMobile />
      </ul>
    </div>
  </div>
</template>

<style scoped>
.navbar {
  background-color: #f0fdf4;
  /* สีเขียวอ่อน */
  border-bottom: 2px solid #bbf7d0;
  /* ขอบด้านล่าง */
}
</style>
