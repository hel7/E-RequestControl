<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <div class="language-switcher">
          <button @click="setLocale('ua')" :class="{ active: currentLocale === 'ua' }">{{ $t('ukrainian') }}</button>
          <button @click="setLocale('en')" :class="{ active: currentLocale === 'en' }">{{ $t('english') }}</button>
        </div>
        <h1>{{ $t('welcome') }}</h1>
        <p>{{ $t('signInMessage') }}</p>
      </div>

      <form @submit.prevent="handleSubmit" class="auth-form">
        <div class="input-group">
          <input
              v-model="username"
              type="text"
              id="username"
              required
              class="form-input"
              placeholder=" "
          />
          <label for="username">{{ $t('username') }}</label>
          <div class="input-border"></div>
        </div>

        <div class="input-group password-group">
          <input
              v-model="password"
              :type="isPasswordVisible ? 'text' : 'password'"
              id="password"
              required
              class="form-input"
              placeholder=" "
          />
          <label for="password">{{ $t('password') }}</label>
          <div class="input-border"></div>
          <button
              type="button"
              class="password-toggle"
              @click="togglePasswordVisibility"
          >
            <i class="fas" :class="isPasswordVisible ? 'fa-eye-slash' : 'fa-eye'"></i>
          </button>
        </div>

        <div class="form-footer">
          <button type="submit" class="submit-btn" :disabled="isLoading">
            <span v-if="!isLoading">{{ $t('signIn') }}</span>
            <span v-else class="loading-spinner"></span>
          </button>

          <div class="auth-links">
            <router-link to="/register" class="auth-link">
              {{ $t('noAccount') }} <span>{{ $t('register') }}</span>
            </router-link>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import i18n from '../i18n';
import { useI18n } from 'vue-i18n';
import { useToast } from "vue-toastification";


const { t, locale } = useI18n({ useScope: 'global' });

const currentLocale = computed(() => locale.value);

const toast = useToast();
const setLocale = async (newLocale) => {
  if (locale.value === newLocale) return;

  try {
    const messages = import.meta.glob('../i18n/locales/*.json');
    const path = `../i18n/locales/${newLocale}.json`;
    const loader = messages[path];

    if (!loader) {
      toast.error(`Locale "${newLocale}" not found`);
      return;
    }

    const mod = await loader();
    i18n.global.setLocaleMessage(newLocale, mod.default);
    locale.value = newLocale;
    localStorage.setItem('locale', newLocale);
  } catch (error) {
    toast.error('Error loading locale file');
    console.error(error);
  }
};

const username = ref('');
const password = ref('');
const isPasswordVisible = ref(false);
const isLoading = ref(false);
const router = useRouter();
const authStore = useAuthStore();

const togglePasswordVisibility = () => {
  isPasswordVisible.value = !isPasswordVisible.value;
};

const handleSubmit = async () => {
  isLoading.value = true;

  try {
    const success = await authStore.login({
      username: (username.value).trim(),
      password: (password.value).trim(),
    });

    if (!success) {
      toast.error(t('invalidCredentials'), {
        icon: "fas fa-exclamation-triangle",
        timeout: 5000,
      });
      return;
    }

    toast.success(t('loginSuccess'), {
      timeout: 1500,
      icon: "fas fa-check-circle"
    });

    setTimeout(() => {
      router.push(authStore.role === 1 ? '/admin-dashboard' : '/dashboard');
    }, 1500);
  } catch (error) {
    let errorMessage = t('loginError');

    if (error.response?.status === 401) {
      errorMessage = t('invalidCredentials');
    } else if (error.response?.status === 400) {
      errorMessage = t('badRequest');
    } else if (error.response?.status === 500) {
      errorMessage = t('serverError');
    }

    toast.error(errorMessage, {
      icon: "fas fa-exclamation-triangle",
      timeout: 5000,
    });
  } finally {
    isLoading.value = false;
  }
};
</script>


<style scoped>
*,
*::before,
*::after {
  box-sizing: border-box;
}

.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  width: 100%;
  padding: 2rem;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.auth-card {
  width: 100%;
  max-width: 450px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  padding: 3rem;
  margin: auto;
}

.auth-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.auth-header h1 {
  font-size: 2.2rem;
  color: #2c3e50;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.auth-header p {
  color: #7f8c8d;
  font-size: 1.1rem;
}

.auth-form {
  margin-bottom: 1.5rem;
}

.input-group {
  position: relative;
  margin-bottom: 1.8rem;
}

.input-group label {
  position: absolute;
  top: 15px;
  left: 15px;
  color: #7f8c8d;
  font-size: 1rem;
  transition: all 0.3s ease;
  pointer-events: none;
  background: white;
  padding: 0 5px;
}

.input-group .form-input {
  width: 100%;
  padding: 16px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.input-group .form-input:focus {
  outline: none;
  border-color: #3498db;
}

.input-group .form-input:focus + label,
.input-group .form-input:not(:placeholder-shown) + label {
  top: -10px;
  left: 10px;
  font-size: 0.85rem;
  color: #3498db;
}

.input-border {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 0;
  height: 2px;
  background: #3498db;
  transition: width 0.3s ease;
}

.input-group .form-input:focus ~ .input-border {
  width: 100%;
}

.password-group {
  position: relative;
}

.password-toggle {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #7f8c8d;
  cursor: pointer;
  font-size: 1.1rem;
}

.submit-btn {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 6px rgba(52, 152, 219, 0.2);
}

.submit-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #2980b9 0%, #3498db 100%);
  transform: translateY(-2px);
  box-shadow: 0 6px 8px rgba(52, 152, 219, 0.3);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0);
}

.submit-btn:disabled {
  background: #bdc3c7;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.loading-spinner {
  display: inline-block;
  width: 22px;
  height: 22px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.form-footer {
  margin-top: 2.5rem;
}

.auth-links {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1.8rem;
  text-align: center;
}

.auth-link {
  color: #7f8c8d;
  text-decoration: none;
  font-size: 0.95rem;
  transition: color 0.3s ease;
}

.auth-link:hover {
  color: #3498db;
}

.auth-link span {
  color: #3498db;
  font-weight: 500;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 768px) {
  .auth-card {
    padding: 2.5rem;
  }
}

@media (max-width: 480px) {
  .auth-container {
    padding: 1.5rem;
  }

  .auth-card {
    padding: 2rem 1.5rem;
  }

  .auth-header h1 {
    font-size: 1.9rem;
  }

  .auth-header p {
    font-size: 1rem;
  }

  .input-group {
    margin-bottom: 1.5rem;
  }

  .submit-btn {
    padding: 15px;
    font-size: 1rem;
  }
}
</style>