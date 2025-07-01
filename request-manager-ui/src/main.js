import { createApp } from 'vue';
import { createPinia } from 'pinia';
import router from './router';
import App from './App.vue';
import './assets/styles.css';
import i18n from './i18n';
import Toast from "vue-toastification";
import "/node_modules/vue-toastification/dist/index.css";

const app = createApp(App);
app.use(createPinia());
app.use(router);
app.use(i18n);
app.use(Toast, {
    transition: "Vue-Toastification__fade",
    maxToasts: 3,
    newestOnTop: true
});
app.mount('#app');
