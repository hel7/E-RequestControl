import { defineStore } from 'pinia';
import { authApi } from '../api';
import { jwtDecode } from 'jwt-decode';

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null,
        token: localStorage.getItem('token') || '',
        role: localStorage.getItem('role') || null,
    }),
    actions: {
        async login(credentials) {
            try {
                const res = await authApi.login(credentials.username, credentials.password);

                this.user = res.data.user;
                this.token = res.data.token;

                const decodedToken = jwtDecode(this.token);
                this.role = decodedToken.RoleID;

                localStorage.setItem('token', this.token);
                localStorage.setItem('role', this.role);

                return true;
            } catch (error) {
                console.error("Помилка авторизації:", error);
                return false;
            }
        },
        async logout(silent = false) {
            try {
                if (this.token) {
                    await authApi.logout();
                }
            } catch (error) {
                if (!silent) {
                    console.error("Logout error:", error);
                }
            } finally {
                this.clearAuthData();
                window.location.href = '/login';
            }
        },

        clearAuthData() {
            this.user = null;
            this.token = '';
            this.role = null;
            localStorage.removeItem('token');
            localStorage.removeItem('role');
            sessionStorage.clear();
        },

    setRole(role) {
            this.role = role;
            localStorage.setItem('role', role);
        }
    },
    getters: {
        isAuthenticated: (state) => !!state.token,
        userRole: (state) => state.role,
    }
});
