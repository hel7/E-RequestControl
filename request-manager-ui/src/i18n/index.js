import { createI18n } from 'vue-i18n';
import ua from './locales/ua.json';
import en from './locales/en.json';

const defaultLocale = localStorage.getItem('locale') || 'ua';

const i18n = createI18n({
    legacy: false,
    locale: defaultLocale,
    fallbackLocale: 'en',
    messages: {
        ua,
        en,
    },
});

const loadLocaleMessages = (locale) => {
    i18n.global.locale.value = locale;
    localStorage.setItem('locale', locale);
};

export { i18n, loadLocaleMessages };
export default i18n;
