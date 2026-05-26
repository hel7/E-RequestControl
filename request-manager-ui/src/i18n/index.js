import { createI18n } from 'vue-i18n';

const defaultLocale = localStorage.getItem('locale') || 'ua';

const i18n = createI18n({
    legacy: false,
    locale: defaultLocale,
    fallbackLocale: 'en',
    messages: {}
});

const loadLocaleMessages = async (locale) => {
    try {
        const messages = await import(`./locales/${locale}.json`);
        i18n.global.setLocaleMessage(locale, messages.default);
        i18n.global.locale.value = locale;
    } catch (e) {
        console.error(`Failed to load locale (${locale}):`, e);
    }
};

await loadLocaleMessages(defaultLocale);

export { i18n, loadLocaleMessages };
export default i18n;
