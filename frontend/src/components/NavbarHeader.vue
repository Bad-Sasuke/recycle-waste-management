<script setup lang="ts">
import { IconUserCircle, IconBell, IconMenu, IconTrash, IconBellOff } from "@tabler/icons-vue";
import { RouterLink } from 'vue-router';
import { ref } from "vue";

const isDrawerOpen = ref(false);

const closeDrawer = () => {

  isDrawerOpen.value = !isDrawerOpen.value;
};

</script>

<template>
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
          <li>
            <RouterLink to="/marketplace" class="font-semibold hover:text-green-600">ราคาขยะรีไซเคิล</RouterLink>
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
          <span class="flex justify-between mx-2 py-4">
            <p class="text-[1rem]">กล่องข้อความ</p>
            <div class="flex justify-between items-center gap-2">
              <li><a class="text-xs">อ่านทั้งหมด</a></li>
              <button class="hover:text-error">
                <IconTrash stroke="1.5" size="24" />
              </button>
            </div>
          </span>
          <div class="flex flex-col items-center justify-center flex-grow gap-2">
            <IconBellOff stroke="1.5" size="24" class="text-gray-500" />
            <p class="text-sm text-gray-500">คุณยังไม่มีแจ้งเตือน</p>
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
            <a>ไปหน้ากล่องข้อความ</a>
          </li>
        </ul>
      </div>

      <!-- โปรไฟล์ผู้ใช้งาน -->
      <div class="dropdown dropdown-end">
        <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
          <div class="rounded-full">
            <IconUserCircle stroke="1.5" size="32" />
          </div>
        </div>
        <ul tabindex="0" class="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow">
          <li>
            <a class="justify-between">
              โปรไฟล์
              <span class="badge">New</span>
            </a>
          </li>

          <li><a>การตั้งค่า</a></li>
          <li><a>ออกจากระบบ</a></li>
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
        <p class="normal-case text-2xl font-bold text-green-700 m-2">
          ♻️ Recycle Waste
        </p>
        <li>
          <RouterLink to="/" @click="closeDrawer">หน้าแรก</RouterLink>
        </li>
        <li>
          <RouterLink to="/marketplace" @click="closeDrawer">ราคาขยะรีไซเคิล</RouterLink>
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
