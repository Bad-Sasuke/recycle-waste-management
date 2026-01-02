<script setup lang="ts">
import { IconLanguageHiragana, IconChevronDown } from '@tabler/icons-vue'
import { useI18nStore } from '@/stores/i18n'
import { useRoute } from 'vue-router'
import { computed } from 'vue'

const i18nStore = useI18nStore()
const route = useRoute()
const isHomePage = computed(() => route.path === '/')

const changeLanguage = (lang: string) => {
  i18nStore.setLocale(lang)
}
</script>

<template>
  <div class="dropdown hidden md:block">
    <div tabindex="0" role="button" :class="[
      'btn btn-outline min-h-0 h-8',
      isHomePage ? 'text-white border-white/50 hover:bg-white/20 hover:border-white' : 'text-neutral'
    ]">
      <IconLanguageHiragana stroke="1.5" size="24" />
      <IconChevronDown stroke="1.5" />
    </div>
    <ul tabindex="0" class="dropdown-content menu bg-base-100 text-base-content rounded-box z-[1] w-32 p-2 shadow">
      <li>
        <button :class="{ active: $i18n.locale === 'en-US' }" @click="changeLanguage('en-US')">
          <span
            class="badge badge-sm badge-outline !pl-1.5 !pr-1 pt-px font-mono !text-[.6rem] font-bold tracking-widest opacity-50">EN</span>
          <span>English</span>
        </button>
      </li>
      <li>
        <button :class="{ active: $i18n.locale === 'th' }" @click="changeLanguage('th')">
          <span
            class="badge badge-sm badge-outline !pl-1.5 !pr-1 pt-px font-mono !text-[.6rem] font-bold tracking-widest opacity-50">TH</span>
          <span>ไทย</span>
        </button>
      </li>
    </ul>
  </div>
</template>
