import { createRouter, createWebHistory } from 'vue-router';
import { jwtDecode } from 'jwt-decode';
import Login from '../views/Login.vue';
import Registration from '../views/Registration.vue';
import Dashboard from "../views/Dashboard.vue";
import AdminDashboard from "../views/AdminDashboard.vue";

const routes = [
    { path: '/', redirect: '/login' },
    {
        path: '/login',
        component: Login,
        meta: { requiresGuest: true }
    },
    {
        path: '/register',
        component: Registration,
        meta: { requiresGuest: true }
    },
    {
        path: '/dashboard',
        component: Dashboard,
        meta: { requiresUser: true, requiresGuest: false }
    },
    {
        path: '/admin-dashboard',
        component: AdminDashboard,
        meta: { requiresAdmin: true, requiresGuest: false  }
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

const isTokenValid = (token) => {
    if (!token) return false;

    try {
        const decoded = jwtDecode(token);
        const currentTime = Date.now() / 1000;

        if (decoded.exp < currentTime) {
            localStorage.removeItem('token');
            localStorage.removeItem('role');
            return false;
        }
        return true;
    } catch (error) {
        console.error('Invalid token:', error);
        localStorage.removeItem('token');
        localStorage.removeItem('role');
        return false;
    }
};

router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role');
    const isAuthenticated = token && isTokenValid(token);

    if (to.meta.requiresGuest && isAuthenticated) {
        return next(role === '1' ? '/admin-dashboard' : '/dashboard');
    }

    if ((to.meta.requiresUser || to.meta.requiresAdmin) && !isAuthenticated) {
        return next('/login');
    }

    if (to.meta.requiresAdmin && role !== '1') {
        return next('/dashboard');
    }

    if (to.meta.requiresUser && role === '1') {
        return next('/admin-dashboard');
    }

    next();
});

export default router;