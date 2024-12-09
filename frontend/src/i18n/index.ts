import { createI18n } from 'vue-i18n';
import en from './en.json';
import th from './th.json';

// กำหนดประเภทของข้อความสำหรับ i18n
type MessageSchema = typeof en;

const messages: Record<string, MessageSchema> = {
    en,
    th,
};

const i18n = createI18n<MessageSchema>({
    locale: 'en', // กำหนดภาษาเริ่มต้น
    fallbackLocale: 'en', // ภาษา fallback
    messages,
});

export default i18n;
