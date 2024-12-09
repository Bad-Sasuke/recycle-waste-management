import { defineStore } from 'pinia';
import { useI18n } from 'vue-i18n';
import { setCookie, getCookie } from './cookie';

export const useI18nStore = defineStore('i18n', {
    state: () => ({
        locale: useI18n().locale
    }),
    actions: {
        setLocale(locale: string) {
            this.locale = locale;
            setCookie('locale', locale, 365);

        },
        getLocale() {
            const locale = getCookie('locale');
            if (locale) {
                this.locale = locale;
            }
        }
    }
});
