<script setup lang="ts">
import { IconUserCircle, IconBell, IconMenu } from "@tabler/icons-vue";


</script>

<template>
  <div class="navbar bg-green-100 shadow-md">

    <!-- เมนูสำหรับ Mobile -->
    <div class="flex-none md:hidde ">
      <label for="menu-mobile" class="btn btn-ghost btn-circle">
        <IconMenu stroke="1.5" size="24" />
      </label>
    </div>




    <!-- โลโก้และชื่อเว็บ -->
    <div class="flex-1">
      <a class="btn btn-ghost normal-case text-2xl font-bold text-green-700">
        ♻️ Recycle Waste
      </a>
    </div>



    <!-- เมนูหลัก -->
    <div class="hidden md:flex">
      <ul class="menu menu-horizontal px-1">
        <li><a href="#home" class="font-semibold hover:text-green-600">หน้าแรก</a></li>
        <li><a href="#marketplace" class="font-semibold hover:text-green-600">ราคาขยะรีไซเคิล</a></li>
      </ul>
    </div>



    <!-- ไอคอน -->
    <div class="flex-none gap-4">

      <!-- ปุ่มแจ้งเตือน -->
      <button class="btn btn-ghost btn-circle">
        <IconBell stroke="1.5" size="24" />
      </button>

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
  <div class="drawer">
    <input id="menu-mobile" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content">
      <!-- Page content here -->
    </div>
    <div class="drawer-side">
      <label for="menu-mobile" aria-label="close sidebar" class="drawer-overlay"></label>
      <ul class="menu bg-base-200 text-base-content min-h-full w-80 p-4">
        <p >
          <a class="btn btn-ghost normal-case text-2xl font-bold text-green-700">
            ♻️ Recycle Waste
          </a>
        </p>
        <!-- Sidebar content here -->
        <li><a href="#home">หน้าแรก</a></li>
        <li><a href="#marketplace">ราคาขยะรีไซเคิล</a></li>
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
