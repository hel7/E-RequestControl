<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <div class="language-switcher">
          <button @click="setLocale('ua')" :class="{ active: currentLocale === 'ua' }">{{ $t('ukrainian') }}</button>
          <button @click="setLocale('en')" :class="{ active: currentLocale === 'en' }">{{ $t('english') }}</button>
        </div>
        <h1>{{ $t('registerTitle') }}</h1>
        <p>{{ $t('registerSubtitle') }}</p>
      </div>

      <form @submit.prevent="registerUser" class="register-form">
        <div class="form-grid">
          <div class="input-group">
            <input
                v-model="firstname"
                type="text"
                id="firstname"
                required
                class="form-input"
                placeholder=" "
            />
            <label for="firstname">{{ $t('firstname') }}</label>
            <div class="input-border"></div>
          </div>

          <div class="input-group">
            <input
                v-model="lastname"
                type="text"
                id="lastname"
                required
                class="form-input"
                placeholder=" "
            />
            <label for="lastname">{{ $t('lastname') }}</label>
            <div class="input-border"></div>
          </div>

          <div class="input-group">
            <input
                v-model="email"
                type="email"
                id="email"
                required
                class="form-input"
                placeholder=" "
            />
            <label for="email">{{ $t('email') }}</label>
            <div class="input-border"></div>
          </div>

          <div class="input-group">
            <input
                v-model="username"
                type="text"
                id="username"
                required
                class="form-input"
                placeholder=" "
            />
            <label for="username">{{ $t('userNameLabel') }}</label>
            <div class="input-border"></div>
          </div>

          <div class="input-group password-group">
            <div class="password-input-wrapper">
              <input
                  v-model="password"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  id="password"
                  required
                  class="form-input"
                  placeholder=" "
                  @input="validatePassword"
              />
              <label for="password">{{ $t('passwordLabel') }}</label>
              <div class="input-border"></div>
              <button
                  type="button"
                  class="password-toggle"
                  @click="togglePasswordVisibility"
              >
                <i class="fas" :class="isPasswordVisible ? 'fa-eye-slash' : 'fa-eye'"></i>
              </button>
            </div>
            <transition name="fade">
              <div v-if="passwordErrors.length" class="password-hints">
                <div
                    v-for="error in passwordErrors"
                    :key="error"
                    class="hint-item"
                >
                  <i class="fas" :class="getHintIcon(error)"></i>
                  <span>{{ getHintText(error) }}</span>
                </div>
              </div>
            </transition>
          </div>

          <div class="input-group password-group">
            <div class="password-input-wrapper">
              <input
                  v-model="confirmPassword"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  id="confirmPassword"
                  required
                  class="form-input"
                  placeholder=" "
                  @input="validatePasswordMatch"
              />
              <label for="confirmPassword">{{ $t('confirmPasswordLabel') }}</label>
              <div class="input-border"></div>
              <button
                  type="button"
                  class="password-toggle"
                  @click="togglePasswordVisibility"
              >
                <i class="fas" :class="isPasswordVisible ? 'fa-eye-slash' : 'fa-eye'"></i>
              </button>
            </div>
            <transition name="fade">
              <div v-if="passwordMatchError" class="error-message">
                <i class="fas fa-exclamation-circle"></i>
                <span>{{ $t('passwordMismatch') }}</span>
              </div>
            </transition>
          </div>
        </div>

        <button type="submit" class="submit-btn" :disabled="isLoading">
          <span v-if="!isLoading">{{ $t('registerButton') }}</span>
          <span v-else class="loading-spinner"></span>
        </button>
      </form>

      <div class="login-redirect">
        <span>{{ $t('alreadyHaveAccount') }}</span>
        <router-link to="/login" class="login-link">{{ $t('loginLink') }}</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { computed, ref, watch } from 'vue';
import { authApi } from '../api';
import { useRouter } from 'vue-router';
import i18n from '../i18n';
import { useI18n } from 'vue-i18n';
import { useToast } from "vue-toastification";

export default {
  name: 'RegisterUser',
  setup() {
    const router = useRouter();
    const firstname = ref('');
    const lastname = ref('');
    const email = ref('');
    const username = ref('');
    const password = ref('');
    const confirmPassword = ref('');
    const isPasswordVisible = ref(false);
    const passwordErrors = ref([]);
    const passwordMatchError = ref('');
    const isLoading = ref(false);
    const toast = useToast();
    const { t, locale } = useI18n({ useScope: 'global' });



    const togglePasswordVisibility = () => {
      isPasswordVisible.value = !isPasswordVisible.value;
    };

    const currentLocale = computed(() => locale.value);

    const setLocale = async (newLocale) => {
      if (locale.value === newLocale) return;

      const messages = import.meta.glob('../i18n/locales/*.json');
      const path = `../i18n/locales/${newLocale}.json`;
      const loader = messages[path];

      if (!loader) {
        console.error(`Locale ${newLocale} not found`);
        return;
      }

      const mod = await loader();
      i18n.global.setLocaleMessage(newLocale, mod.default);
      locale.value = newLocale;
      localStorage.setItem('locale', newLocale);
    };

    const getHintIcon = (error) => {
      const isError = passwordErrors.value.includes(error);
      return isError ? 'fa-times-circle error' : 'fa-check-circle success';
    };


    const passwordRules = [
      { condition: () => password.value.length < 8, key: 'passwordMinLength' },
      { condition: () => !/[a-zA-Z]/.test(password.value), key: 'passwordLetter' },
      { condition: () => !/\d/.test(password.value), key: 'passwordDigit' },
      { condition: () => !/[^a-zA-Z0-9\s]/.test(password.value), key: 'passwordSpecialChar' }
    ];

    const validatePassword = () => {
      const errors = [];

      passwordRules.forEach(rule => {
        if (rule.condition() && password.value) {
          errors.push(rule.key);
        }
      });

      passwordErrors.value = errors;
      return errors.length === 0;
    };
    const getHintText = (key) => t(key);


    const validatePasswordMatch = () => {
      if (password.value && confirmPassword.value && password.value !== confirmPassword.value) {
        passwordMatchError.value = t('passwordMismatch');
        return false;
      } else {
        passwordMatchError.value = '';
        return true;
      }
    };
    const registerUser = async () => {
      isLoading.value = true;
      try {
        const response = await authApi.register(
            firstname.value,
            lastname.value,
            email.value,
            username.value,
            password.value
        );

        toast.success(t('registerSuccess'), {
          timeout: 2000,
          icon: "fas fa-check-circle",
        });

        setTimeout(() => {
          router.push('/login');
        }, 2000);

      } catch (error) {
        const errorMessage = error.response?.data?.message || error.message;
        toast.error(t('registerError') + ': ' + errorMessage, {
          icon: "fas fa-exclamation-triangle",
          timeout: 5000,
        });
      } finally {
        isLoading.value = false;
      }
    };

    watch(() => password.value, () => {
      validatePassword();
    });
    watch(() => password.value, () => {
      validatePasswordMatch();
    });
    return {
      firstname,
      lastname,
      email,
      username,
      password,
      confirmPassword,
      isPasswordVisible,
      passwordErrors,
      passwordMatchError,
      isLoading,
      togglePasswordVisibility,
      getHintIcon,
      getHintText,
      registerUser,
      currentLocale,
      setLocale
    };
  }
};
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 2rem;
}

.register-card {
  width: 100%;
  max-width: 900px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  padding: 2.5rem;
}

.register-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.register-header h1 {
  font-size: 2.2rem;
  color: #2c3e50;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.register-header p {
  color: #7f8c8d;
  font-size: 1.1rem;
}

.register-form {
  margin-bottom: 2rem;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
  margin-bottom: 2rem;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .register-card {
    padding: 1.5rem;
  }

  .register-header h1 {
    font-size: 1.8rem;
  }
}

.input-group {
  position: relative;
  margin-bottom: 1.5rem;
}

.input-group label {
  position: absolute;
  top: 12px;
  left: 15px;
  color: #7f8c8d;
  font-size: 0.9rem;
  transition: all 0.3s ease;
  pointer-events: none;
  background: white;
  padding: 0 5px;
}

.input-group .form-input {
  width: 100%;
  padding: 15px;
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
  font-size: 0.8rem;
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


.password-input-wrapper {
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
  font-size: 1rem;
  z-index: 2;
}

.password-hints,
.error-message {
  margin-top: 8px;
}
.password-hints {
  margin-top: 0.5rem;
  background: #f8f9fa;
  border-radius: 8px;
  padding: 0.8rem;
  border-left: 4px solid #e74c3c;
}

.hint-item {
  display: flex;
  align-items: center;
  margin-bottom: 0.3rem;
  font-size: 0.85rem;
}

.hint-item i {
  margin-right: 0.5rem;
}

.hint-item .error {
  color: #e74c3c;
}

.hint-item .success {
  color: #2ecc71;
}

.error-message {
  color: #e74c3c;
  font-size: 0.85rem;
  margin-top: 0.5rem;
  display: flex;
  align-items: center;
}

.error-message i {
  margin-right: 0.5rem;
}

.submit-btn {
  width: 100%;
  padding: 15px;
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 6px rgba(52, 152, 219, 0.2);
}

.submit-btn:hover {
  background: linear-gradient(135deg, #2980b9 0%, #3498db 100%);
  transform: translateY(-2px);
  box-shadow: 0 6px 8px rgba(52, 152, 219, 0.3);
}

.submit-btn:active {
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
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.login-redirect {
  text-align: center;
  color: #7f8c8d;
  font-size: 0.9rem;
}

.login-link {
  color: #3498db;
  text-decoration: none;
  margin-left: 0.5rem;
  font-weight: 500;
  transition: color 0.3s ease;
}

.login-link:hover {
  color: #2980b9;
  text-decoration: underline;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>