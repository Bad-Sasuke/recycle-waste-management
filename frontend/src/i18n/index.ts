import { createI18n } from 'vue-i18n';
import en from './en.json';
import th from './th.json';

type MessageSchema = typeof en;

const messages: Record<'en-US' | 'th', MessageSchema> = {
    'en-US': en,
    'th': th,
};

// สร้าง instance ของ i18n
const i18n = createI18n({
    legacy: false,
    locale: 'en-US',
    fallbackLocale: 'th',
    messages,
});

export default i18n;
