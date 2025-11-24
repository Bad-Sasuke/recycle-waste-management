<script setup lang="ts">
import { IconUserCircle, IconBell, IconMenu, IconTrash, IconBellOff, IconBuildingStore, IconHome, IconBuildingWarehouse, IconSettings, IconLogout, IconUser, IconMapPin, IconMapPinShare, IconUsers } from '@tabler/icons-vue'
import { RouterLink } from 'vue-router'
import { ref, onBeforeMount, onMounted, watch } from 'vue'
import { useI18nStore } from '@/stores/i18n'
import { useUsersStore } from '@/stores/users'
import { useShopStore } from '@/stores/shop'
import PopupLogin from './PopupLogin.vue'
import SwitchLangDesktop from './switchLang/SwitchLangDesktop.vue'
import SwitchLangMobile from './switchLang/SwitchLangMobile.vue'

const isDrawerOpen = ref(false)

const i18nStore = useI18nStore()
const usersStore = useUsersStore()
const shopStore = useShopStore()

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

onMounted(async () => {
  await usersStore.checkLogin()
  if (usersStore.isLogin) {
    await shopStore.checkUserShop()
  }
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
      <RouterLink to="/" class="btn btn-ghost normal-case text-xl md:text-2xl font-bold text-green-700">
        ♻️ Recycle Waste
      </RouterLink>
      <!-- เมนูหลัก -->
      <div class="hidden md:flex">
        <ul class="menu menu-horizontal px-1">
          <li class="mr-2">
            <RouterLink v-if="usersStore.isLogin" to="/marketplace"
              class="flex items-center gap-2 text-base font-medium hover:text-green-800 transition-colors">
              <IconBuildingWarehouse stroke="1.5" size="20" />
              {{ $t('Navbar.menu.marketplace') }}
            </RouterLink>
            <a v-else @click="handleLogin"
              class="flex items-center gap-2 text-base font-medium hover:text-green-800 transition-colors cursor-pointer">
              <IconBuildingWarehouse stroke="1.5" size="20" />
              {{ $t('Navbar.menu.marketplace') }}
            </a>
          </li>
          <li class="mr-2" v-if="usersStore.user?.role !== 'moderator'">
            <RouterLink v-if="usersStore.isLogin" to="/shop-locator"
              class="flex items-center gap-2 text-base font-medium hover:text-green-800 transition-colors">
              <IconMapPin stroke="1.5" size="20" />
              {{ $t('ShopLocator.title') }}
            </RouterLink>
            <a v-else @click="handleLogin"
              class="flex items-center gap-2 text-base font-medium hover:text-green-800 transition-colors cursor-pointer">
              <IconMapPin stroke="1.5" size="20" />
              {{ $t('ShopLocator.title') }}
            </a>
          </li>
          <li v-if="usersStore.isLogin && shopStore.hasShop" class="mr-2">
            <RouterLink to="/nearby-customers"
              class="flex items-center gap-2 text-base font-medium hover:text-green-800 transition-colors">
              <IconUsers stroke="1.5" size="20" />
              {{ $t('NearbyCustomers.title') }}
            </RouterLink>
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
        <ul tabindex="0"
          class="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-72 p-2 shadow h-64 flex flex-col">
          <span class="flex justify-between mx-2 py-4 border-b border-gray-200">
            <p class="text-lg font-semibold text-gray-800">{{ $t('Navbar.notif.title') }}</p>
            <div class="flex justify-between items-center gap-2">
              <li class="px-2">
                <a class="text-sm font-medium text-green-600 hover:text-green-700">{{ $t('Navbar.notif.read_all') }}</a>
              </li>
              <button class="hover:text-error p-1">
                <IconTrash stroke="1.5" size="20" />
              </button>
            </div>
          </span>
          <div class="flex flex-col items-center justify-center flex-grow gap-3">
            <IconBellOff stroke="1.5" size="28" class="text-gray-400" />
            <p class="text-base text-gray-500 font-medium">{{ $t('Navbar.notif.empty') }}</p>
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
              <a>ไปหน้ากล่องข้อหมู</a>
            </li>
            <li>
              <a>ไปหน้ากล่องข้อความ</a>
            </li> -->
          </div>
          <li class="flex items-center justify-center w-full mt-auto py-3">
            <a class="text-center text-base font-medium text-green-600 hover:text-green-700">{{
              $t('Navbar.notif.view_all') }}</a>
          </li>
        </ul>
      </div>

      <SwitchLangDesktop class="px-2" />

      <!-- โปรไฟล์ผู้ใช้งาน -->
      <div class="dropdown dropdown-end" style="z-index: 999;"
        @click="usersStore.isLogin !== true ? handleLogin() : ''">
        <div tabindex="0" role="button" class="flex items-center">
          <div v-if="usersStore.isLogin && usersStore.profileImage" class="avatar placeholder">
            <div class="bg-neutral-focus text-neutral-content rounded-full w-8">
              <img :src="usersStore.profileImage" alt="Profile" class="rounded-full" />
            </div>
          </div>
          <IconUserCircle v-else stroke="1.3" size="32" class="bg-neutral rounded-full text-base-100" />
        </div>
        <ul tabindex="0" class="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow"
          v-if="usersStore.isLogin">
          <span v-if="usersStore.isLogin">
            <li class="py-1">
              <RouterLink to="/profile" class="flex items-center gap-3 text-base font-medium hover:text-green-700">
                <div class="flex items-center gap-2">
                  <IconUser stroke="1.5" size="20" />
                  <span>{{ $t('Navbar.profile.textProfile') }}</span>
                </div>
              </RouterLink>
            </li>

            <li class="py-1" v-if="shopStore.hasShop">
              <RouterLink to="/manage-shop" class="flex items-center gap-3 text-base font-medium hover:text-green-700">
                <div class="flex items-center gap-2">
                  <IconBuildingStore stroke="1.5" size="20" />
                  <span>{{ $t('Navbar.profile.manageShop') }}</span>
                </div>
              </RouterLink>
            </li>

            <li class="py-1">
              <RouterLink to="/settings" class="flex items-center gap-3 text-base font-medium hover:text-green-700">
                <div class="flex items-center gap-2">
                  <IconSettings stroke="1.5" size="20" />
                  <span>{{ $t('Navbar.profile.textSettings') }}</span>
                </div>
              </RouterLink>
            </li>
            <li class="py-1">
              <a @click="handleLogout"
                class="flex items-center gap-3 text-base font-medium text-error hover:text-error-focus transition-colors">
                <div class="flex items-center gap-2">
                  <IconLogout stroke="1.5" size="20" />
                  <span>{{ $t('Navbar.profile.textLogout') }}</span>
                </div>
              </a>
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
        <p class="normal-case text-xl md:text-2xl font-bold text-green-700 m-2 mb-4">♻️ Recycle Waste</p>
        <li class="py-2">
          <RouterLink to="/" @click="closeDrawer" class="flex items-center gap-3 text-base font-medium">
            <IconHome stroke="1.5" size="20" />
            <span>{{ $t('Navbar.menu.home') }}</span>
          </RouterLink>
        </li>
        <li class="py-2">
          <RouterLink v-if="usersStore.isLogin" to="/marketplace" @click="closeDrawer"
            class="flex items-center gap-3 text-base font-medium">
            <IconBuildingWarehouse stroke="1.5" size="20" />
            <span>{{ $t('Navbar.menu.marketplace') }}</span>
          </RouterLink>
          <a v-else @click="() => { closeDrawer(); handleLogin(); }"
            class="flex items-center gap-3 text-base font-medium cursor-pointer">
            <IconBuildingWarehouse stroke="1.5" size="20" />
            <span>{{ $t('Navbar.menu.marketplace') }}</span>
          </a>
        </li>
        <li class="py-2" v-if="usersStore.user?.role !== 'moderator'">
          <RouterLink v-if="usersStore.isLogin" to="/shop-locator" @click="closeDrawer"
            class="flex items-center gap-3 text-base font-medium">
            <IconMapPin stroke="1.5" size="20" />
            <span>{{ $t('ShopLocator.title') }}</span>
          </RouterLink>
          <a v-else @click="() => { closeDrawer(); handleLogin(); }"
            class="flex items-center gap-3 text-base font-medium cursor-pointer">
            <IconMapPin stroke="1.5" size="20" />
            <span>{{ $t('ShopLocator.title') }}</span>
          </a>
        </li>
        <li class="py-2" v-if="usersStore.isLogin">
          <RouterLink to="/share-location" @click="closeDrawer" class="flex items-center gap-3 text-base font-medium">
            <IconMapPinShare stroke="1.5" size="20" />
            <span>{{ $t('ShareLocation.title') }}</span>
          </RouterLink>
        </li>
        <li class="py-2" v-if="usersStore.isLogin && shopStore.hasShop">
          <RouterLink to="/nearby-customers" @click="closeDrawer" class="flex items-center gap-3 text-base font-medium">
            <IconUsers stroke="1.5" size="20" />
            <span>{{ $t('NearbyCustomers.title') }}</span>
          </RouterLink>
        </li>
        <li class="py-2" v-if="usersStore.isLogin && shopStore.hasShop">
          <RouterLink to="/manage-shop" @click="closeDrawer" class="flex items-center gap-3 text-base font-medium">
            <IconBuildingStore stroke="1.5" size="20" />
            <span>{{ $t('Navbar.profile.manageShop') }}</span>
          </RouterLink>
        </li>
        <SwitchLangMobile class="mt-2 py-2" />
        <li v-if="usersStore.isLogin" class="mt-auto pt-4 border-t border-gray-400">
          <div class="flex items-center gap-3 p-2">
            <div v-if="usersStore.profileImage" class="avatar">
              <div class="bg-neutral-focus text-neutral-content rounded-full w-10">
                <img :src="usersStore.profileImage" alt="Profile" class="rounded-full" />
              </div>
            </div>
            <div v-else class="p-1">
              <IconUserCircle stroke="1.3" size="24" class="bg-neutral rounded-full text-base-100" />
            </div>
            <div class="flex-1">
              <p class="font-semibold text-base">{{ usersStore.username }}</p>
              <p class="text-xs opacity-60">Logged in</p>
            </div>
          </div>
        </li>
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
