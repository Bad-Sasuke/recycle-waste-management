<script setup lang="ts">
import { IconLanguageHiragana, IconChevronDown } from '@tabler/icons-vue'
import { useI18nStore } from '@/stores/i18n'

const i18nStore = useI18nStore()
const changeLanguage = (lang: string) => {
  i18nStore.setLocale(lang)
}
</script>

<template>
  <div class="dropdown mx-4 md:hidden">
    <div tabindex="0" role="button" class="btn btn-outline p-1 text-neutral w-full">
      <IconLanguageHiragana stroke="1.5" />
      <IconChevronDown stroke="1.5" />
    </div>
    <ul
      tabindex="0"
      class="dropdown-content menu bg-base-100 rounded-box z-[1] w-32 p-2 shadow w-full"
    >
      <li>
        <button :class="{ active: $i18n.locale === 'en-US' }" @click="changeLanguage('en-US')">
          <span
            class="badge badge-sm badge-outline !pl-1.5 !pr-1 pt-px font-mono !text-[.6rem] font-bold tracking-widest opacity-50"
            >EN</span
          >
          <span>English</span>
        </button>
      </li>
      <li>
        <button :class="{ active: $i18n.locale === 'th' }" @click="changeLanguage('th')">
          <span
            class="badge badge-sm badge-outline !pl-1.5 !pr-1 pt-px font-mono !text-[.6rem] font-bold tracking-widest opacity-50"
            >TH</span
          >
          <span>ไทย</span>
        </button>
      </li>
    </ul>
  </div>
</template>
