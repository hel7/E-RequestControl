<template>
  <div class="admin-dashboard">
    <div class="admin-dashboard-card">
      <header class="admin-header">
        <div class="language-switcher">
          <button @click="setLocale('ua')" :class="{ active: currentLocale === 'ua' }">{{ $t('ukrainian') }}</button>
          <button @click="setLocale('en')" :class="{ active: currentLocale === 'en' }">{{ $t('english') }}</button>
        </div>
        <h1>{{ $t('adminPanelTitle') }}</h1>
        <button class="logout-btn" @click="authStore.logout()">
          <i class="fas fa-sign-out-alt"></i> {{ $t('logout') }}
        </button>
      </header>

      <!-- Users Section -->
      <section class="users-section" v-if="users.length > 0">
        <div class="section-header">
          <h2>{{ $t('users') }}</h2>
          <button class="create-btn" @click="showCreateUserModal = true">
            <i class="fas fa-plus"></i> {{ $t('addUser') }}
          </button>
        </div>
        <div class="table-container">
          <table class="data-table">
            <thead>
            <tr>
              <th>{{ $t('firstName') }}</th>
              <th>{{ $t('lastName') }}</th>
              <th>{{ $t('username') }}</th>
              <th>{{ $t('email') }}</th>
              <th>{{ $t('role') }}</th>
              <th>{{ $t('createdAt') }}</th>
              <th>{{ $t('updatedAt') }}</th>
              <th>{{ $t('actions') }}</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="user in paginatedUsers" :key="user.UserID">
              <td>{{ user.FirstName }}</td>
              <td>{{ user.LastName }}</td>
              <td>{{ user.Username }}</td>
              <td>{{ user.Email }}</td>
              <td>
                  <span :class="['role-badge', {'admin': user.RoleID === 1, 'user': user.RoleID === 2}]">
                    {{ user.RoleID === 1 ? $t('admin') : user.RoleID === 2 ? $t('user') : $t('unknownRole') }}
                  </span>
              </td>
              <td>{{ formatDate(user.CreatedAt) }}</td>
              <td>{{ formatDate(user.UpdatedAt) }}</td>
              <td class="actions">
                <button class="edit-btn" @click="editUser(user)">
                  <i class="fas fa-edit"></i>
                </button>
                <button class="delete-btn" @click="confirmDelete('user', user.UserID, user.Username)">
                  <i class="fas fa-trash-alt"></i>
                </button>
              </td>
            </tr>
            </tbody>
          </table>
          <div class="pagination-controls" v-if="totalUsersPages > 1">
            <button class="pagination-btn" @click="prevUsersPage" :disabled="currentUsersPage === 1">
              <i class="fas fa-chevron-left"></i>
            </button>
            <span class="page-info">
              {{ $t('pageInfo', { current: currentUsersPage, total: totalUsersPages }) }}
            </span>
            <button class="pagination-btn" @click="nextUsersPage" :disabled="currentUsersPage === totalUsersPages">
              <i class="fas fa-chevron-right"></i>
            </button>
          </div>
        </div>
      </section>
      <div class="empty-state" v-else>
        <i class="fas fa-users empty-icon"></i>
        <p>{{ $t('noUsers') }}</p>
      </div>

      <!-- Tickets Section -->
      <section class="tickets-section">
        <div class="section-header">
          <h2>{{ $t('tickets') }}</h2>
        </div>

        <!-- Filter Controls -->
        <div class="ticket-filters">
          <div class="filter-group">
            <label for="assignee-filter">{{ $t('assignee') }}:</label>
            <input
                id="assignee-filter"
                v-model="filters.assignee"
                type="text"
                :placeholder="$t('assigneeFilterPlaceholder')"
                @input="applyFilters"
            >
          </div>
          <div class="filter-group">
            <label for="sender-filter">{{ $t('sender') }}:</label>
            <input
                id="sender-filter"
                v-model="filters.sender"
                type="text"
                :placeholder="$t('senderFilterPlaceholder')"
                @input="applyFilters"
            >
          </div>
          <div class="filter-group">
            <label for="status-filter">{{ $t('status') }}:</label>
            <select
                id="status-filter"
                v-model="filters.status"
                @change="applyFilters"
            >
              <option value="">{{ $t('allStatuses') }}</option>
              <option v-for="status in statusOptions" :key="status" :value="status">
                {{ $t(messageMap[status] || status) }}
              </option>
            </select>
          </div>

          <button class="reset-btn" @click="resetFilters">
            <i class="fas fa-times"></i> {{ $t('resetFilters') }}
          </button>
        </div>

        <!-- Tickets Table -->
        <div class="table-container" v-if="tickets.length > 0">
          <table class="data-table">
            <thead>
            <tr>
              <th>{{ $t('title') }}</th>
              <th>{{ $t('description') }}</th>
              <th>{{ $t('status') }}</th>
              <th>{{ $t('assignee') }}</th>
              <th>{{ $t('sender') }}</th>
              <th>{{ $t('createdAt') }}</th>
              <th>{{ $t('updatedAt') }}</th>
              <th>{{ $t('actions') }}</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="ticket in paginatedTickets" :key="ticket.TicketID">
              <td>{{ ticket.Title }}</td>
              <td class="ticket-description">
                  <span class="description-text">
                    {{ truncateDescription(ticket.Description, ticket.TicketID) }}
                  </span>
                <button
                    v-if="ticket.Description && ticket.Description.length > 100"
                    class="show-more-btn"
                    @click="toggleDescription(ticket.TicketID)"
                >
                  {{ expandedDescriptions[ticket.TicketID] ? '▲' : '▼' }}
                </button>
              </td>
              <td>
                  <span :class="['status-badge', getStatusClass(ticket.Status)]">
                    {{ $t(messageMap[ticket.Status] || ticket.Status) }}
                  </span>
              </td>
              <td>{{ ticket.AssigneeUsername }}</td>
              <td>{{ ticket.SenderUsername }}</td>
              <td>{{ formatDate(ticket.CreatedAt) }}</td>
              <td>{{ formatDate(ticket.UpdatedAt) }}</td>
              <td class="actions">
                <button class="action-btn delete-btn" @click="confirmDelete('ticket', ticket.TicketID, ticket.Title)" :title="$t('delete')">
                  <i class="fas fa-trash-alt"></i>
                </button>
              </td>
            </tr>
            </tbody>
          </table>

          <!-- Pagination Controls -->
          <div class="pagination-controls" v-if="totalPages > 1">
            <button class="pagination-btn" @click="prevPage" :disabled="currentPage === 1">
              <i class="fas fa-chevron-left"></i>
            </button>
            <span class="page-info">
              {{ $t('pageInfo', { current: currentPage, total: totalPages }) }}
            </span>
            <button class="pagination-btn" @click="nextPage" :disabled="currentPage === totalPages">
              <i class="fas fa-chevron-right"></i>
            </button>
          </div>
        </div>

        <div class="empty-state" v-else>
          <i class="fas fa-ticket-alt empty-icon"></i>
          <p>{{ $t('noTickets') }}</p>
          <button class="reset-btn" @click="resetFilters">
            <i class="fas fa-times"></i> {{ $t('resetFiltersButton') }}
          </button>
        </div>
      </section>

      <!-- Notifications Section -->
      <section class="notifications-section" v-if="notifications.length > 0">
        <h2>{{ $t('notifications') }}</h2>
        <div class="notifications-list">
          <div v-for="notification in paginatedNotifications" :key="notification.NotificationID" class="notification-item">
            <div class="notification-content">
              <p class="notification-message">{{ $t(messageMap[notification.Message] || notification.Message) }}</p>
              <span class="notification-time">{{ formatDate(notification.CreatedAt) }}</span>
            </div>
            <button class="action-btn small delete" @click="confirmDelete('notification', notification.NotificationID, $t('notification'))">
              <i class="fas fa-trash-alt"></i>
            </button>
          </div>
          <div class="pagination-controls" v-if="notifications.length > itemsPerPage">
            <button class="pagination-btn" @click="prevNotificationsPage" :disabled="currentNotificationsPage === 1">
              <i class="fas fa-chevron-left"></i>
            </button>
            <span class="page-info">
              {{ $t('pageInfo', { current: currentNotificationsPage, total: totalNotificationsPages }) }}
            </span>
            <button class="pagination-btn" @click="nextNotificationsPage" :disabled="currentNotificationsPage === totalNotificationsPages">
              <i class="fas fa-chevron-right"></i>
            </button>
          </div>
        </div>
      </section>
      <div class="empty-state" v-else>
        <i class="fas fa-bell empty-icon"></i>
        <p>{{ $t('noNotifications') }}</p>
      </div>

      <!-- Data Management Section -->
      <section class="data-management-section">
        <div class="section-header">
          <h2>{{ $t('dataManagement') }}</h2>
        </div>

        <div class="data-grid">
          <div class="data-card" @click="backupData">
            <div class="card-icon backup">
              <i class="fas fa-database"></i>
            </div>
            <h3>{{ $t('backup') }}</h3>
            <p>{{ $t('backupDescription') }}</p>
          </div>

          <div class="data-card" @click="restoreData">
            <div class="card-icon restore">
              <i class="fas fa-redo"></i>
            </div>
            <h3>{{ $t('restore') }}</h3>
            <p>{{ $t('restoreDescription') }}</p>
          </div>

          <div class="data-card" @click="exportData">
            <div class="card-icon export">
              <i class="fas fa-file-export"></i>
            </div>
            <h3>{{ $t('export') }}</h3>
            <p>{{ $t('exportDescription') }}</p>
          </div>

          <div class="data-card" @click="importData">
            <div class="card-icon import">
              <i class="fas fa-file-import"></i>
            </div>
            <h3>{{ $t('import') }}</h3>
            <p>{{ $t('importDescription') }}</p>
          </div>
        </div>
      </section>

      <!-- Create User Modal -->
      <transition name="modal-fade">
        <div
            v-if="showCreateUserModal"
            class="modal-overlay"
            @mousedown="closeModalOnOutsideClick($event, 'create')"
        >
          <div class="modal">
            <div class="modal-header">
              <h3>{{ $t('createUserTitle') }}</h3>
              <button class="close-btn" @click="showCreateUserModal = false">
                <i class="fas fa-times"></i>
              </button>
            </div>
            <div class="modal-body">
              <form @submit.prevent="submitCreateUser">
                <div class="form-group">
                  <label>{{ $t('firstName') }}</label>
                  <input v-model="newUser.FirstName" type="text" :placeholder="$t('enterFirstName')" required />
                </div>
                <div class="form-group">
                  <label>{{ $t('lastName') }}</label>
                  <input v-model="newUser.LastName" type="text" :placeholder="$t('enterLastName')" required />
                </div>
                <div class="form-group">
                  <label>{{ $t('username') }}</label>
                  <input v-model="newUser.Username" type="text" :placeholder="$t('enterUsername')" required />
                </div>
                <div class="form-group">
                  <label>{{ $t('email') }}</label>
                  <input v-model="newUser.Email" type="email" :placeholder="$t('enterEmail')" required />
                </div>
                <div class="form-group">
                  <label>{{ $t('passwordLabel') }}</label>
                  <div class="input-group">
                    <div class="password-wrapper">
                      <input
                          :type="isPasswordVisible.create ? 'text' : 'password'"
                          v-model="newUser.Password"
                          @input="validatePassword(newUser.Password, 'create')"
                          :placeholder="$t('passwordMinLength')"
                      />
                      <button
                          type="button"
                          @click="togglePasswordVisibility('create')"
                          class="eye-button"
                          :class="{ 'active': isPasswordVisible.create  }"
                      >
                        <i class="fas" :class="isPasswordVisible.create  ? 'fa-eye-slash' : 'fa-eye'"></i>
                      </button>
                    </div>

                    <div class="password-strength">
                      <div class="strength-bar" :class="{
                          'weak': passwordStrength === 1,
                          'medium': passwordStrength === 2 || passwordStrength === 3,
                          'strong': passwordStrength === 4
                        }"></div>
                      <div class="strength-labels">
                        <span :class="{ 'active': passwordStrength > 0 }">{{ $t('passwordMinLength') }}</span>
                        <span :class="{ 'active': passwordStrength > 1 }">{{ $t('passwordLetter') }}</span>
                        <span :class="{ 'active': passwordStrength > 2 }">{{ $t('passwordDigit') }}</span>
                        <span :class="{ 'active': passwordStrength > 3 }">{{ $t('passwordSpecialChar') }}</span>
                      </div>
                    </div>

                    <div v-if="passwordErrors.create?.length" class="error-messages">
                      <p v-for="error in passwordErrors.create" :key="error" class="error-text">{{ $t(error) }}</p>
                    </div>
                  </div>
                </div>
                <div class="form-group">
                  <label>{{ $t('role') }}</label>
                  <select v-model="newUser.RoleID" required>
                    <option value="2">{{ $t('user') }}</option>
                    <option value="1">{{ $t('admin') }}</option>
                  </select>
                </div>
                <div class="form-actions">
                  <button type="button" class="cancel-btn" @click="resetNewUserForm; showCreateUserModal = false">
                    {{ $t('cancel') }}
                  </button>
                  <button type="submit" class="submit-btn">{{ $t('create') }}</button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </transition>

      <!-- Edit User Modal -->
      <transition name="modal-fade">
        <div v-if="editUserData" class="modal-overlay" @click.self="editUserData = null">
          <div class="modal">
            <div class="modal-header">
              <h3>{{ $t('editUserTitle') }}</h3>
              <button class="close-btn" @click="editUserData = null">
                <i class="fas fa-times"></i>
              </button>
            </div>
            <div class="modal-body">
              <form @submit.prevent="submitEditUser">
                <div class="form-group">
                  <label>{{ $t('firstName') }}</label>
                  <input v-model="editUserData.FirstName" type="text" :placeholder="$t('enterFirstName')" required />
                </div>
                <div class="form-group">
                  <label>{{ $t('lastName') }}</label>
                  <input v-model="editUserData.LastName" type="text" :placeholder="$t('enterLastName')" required />
                </div>
                <div class="form-group">
                  <label>{{ $t('username') }}</label>
                  <input v-model="editUserData.Username" type="text" :placeholder="$t('enterUsername')" required />
                </div>
                <div class="form-group">
                  <label>{{ $t('email') }}</label>
                  <input v-model="editUserData.Email" type="email" :placeholder="$t('enterEmail')" required />
                </div>

                <div class="form-group">
                  <label>{{ $t('newPasswordPlaceholder') }}</label>
                  <div class="input-group">
                    <div class="password-wrapper">
                      <input
                          :type="isPasswordVisible.edit ? 'text' : 'password'"
                          v-model="editUserData.Password"
                          @input="validatePassword(editUserData.Password, 'edit')"
                          :placeholder="$t('passwordMinLength')"
                      />
                      <button
                          type="button"
                          @click="togglePasswordVisibility('edit')"
                          class="eye-button"
                          :class="{ 'active': isPasswordVisible.edit }"
                      >
                        <i class="fas" :class="isPasswordVisible.edit ? 'fa-eye-slash' : 'fa-eye'"></i>
                      </button>
                    </div>

                    <div class="password-strength" v-if="editUserData.Password">
                      <div class="strength-bar" :class="{
                          'weak': passwordStrength === 1,
                          'medium': passwordStrength === 2 || passwordStrength === 3,
                          'strong': passwordStrength === 4
                        }"></div>
                      <div class="strength-labels">
                        <span :class="{ 'active': passwordStrength > 0 }">{{ $t('passwordMinLength') }}</span>
                        <span :class="{ 'active': passwordStrength > 1 }">{{ $t('passwordLetter') }}</span>
                        <span :class="{ 'active': passwordStrength > 2 }">{{ $t('passwordDigit') }}</span>
                        <span :class="{ 'active': passwordStrength > 3 }">{{ $t('passwordSpecialChar') }}</span>
                      </div>
                    </div>

                    <div v-if="passwordErrors.edit?.length" class="error-messages">
                      <p v-for="error in passwordErrors.edit" :key="error" class="error-text">{{ $t(error) }}</p>
                    </div>
                  </div>
                </div>

                <div class="form-group">
                  <label>{{ $t('role') }}</label>
                  <select v-model="editUserData.RoleID" required>
                    <option value="2">{{ $t('user') }}</option>
                    <option value="1">{{ $t('admin') }}</option>
                  </select>
                </div>
                <div class="form-actions">
                  <button type="button" class="cancel-btn" @click="editUserData = null">{{ $t('cancel') }}</button>
                  <button type="submit" class="submit-btn">{{ $t('save') }}</button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </transition>

      <!-- Confirmation Modal -->
      <transition name="modal-fade">
        <div v-if="confirmModal.show" class="modal-overlay" @click.self="confirmModal.show = false">
          <div class="confirm-modal">
            <div class="modal-header">
              <h3>{{ $t('confirmationTitle') }}</h3>
            </div>
            <div class="modal-body">
              <p>{{ $t('confirmationMessage', {
                type: $t(confirmModal.type),
                name: confirmModal.name
              }) }}</p>
            </div>
            <div class="modal-footer">
              <button class="cancel-btn" @click="confirmModal.show = false">{{ $t('cancel') }}</button>
              <button class="confirm-btn" @click="executeDelete">{{ $t('delete') }}</button>
            </div>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script>
import { reactive, ref, onMounted,watch, computed } from "vue";
import { adminApi } from "../api";
import { format } from 'date-fns';
import { useAuthStore } from '../store/auth';
import i18n from '../i18n';
import { useI18n } from 'vue-i18n';
import { useToast } from "vue-toastification";
export default {
  name: "AdminDashboard",
  setup() {
    const tickets = ref([]);
    const filteredTickets = ref([]);
    const notifications = ref([]);
    const users = ref([]);
    const editUserData = ref(null);
    const isPasswordVisible = reactive({
      create: false,
      edit: false
    });
    const showCreateUserModal = ref(false);
    const authStore = useAuthStore();
    const { t, locale } = useI18n({ useScope: 'global' });
    const statusOptions = ref(['Новий', 'Оновлено']);
    const expandedDescriptions = ref({});
    const currentPage = ref(1);
    const itemsPerPage = ref(10);
    const currentNotificationsPage = ref(1);
    const currentUsersPage = ref(1);
    const toast = useToast();
    const messageMap = {
      'Оновлено тікет': 'notificationMessages.ticketUpdated',
      'Створено новий тікет': 'notificationMessages.ticketCreated',
      'Новий':'statusMessages.statusCreated',
      'Оновлено':'statusMessages.statusUpdated'
    };
    const passwordErrors = reactive({
      create: [],
      edit: []
    });
    const passwordStrength = ref(0);
    const hasMinLength = ref(false);
    const hasLetter = ref(false);
    const hasNumber = ref(false);
    const hasSpecialChar = ref(false);

    const closeModalOnOutsideClick = (event, modalType) => {
      const isClickInside = event.target.closest('.modal');
      if (!isClickInside) {
        if (modalType === 'create') {
          resetNewUserForm();
          showCreateUserModal.value = false;
        } else {
          editUserData.value = null;
        }
      }
    };

    const handleApiError = (error, defaultMessage) => {
      const serverMessage = error.response?.data?.message ||
          error.response?.data?.error ||
          error.message;
      const message = `${t(defaultMessage)}: ${serverMessage}`;
      toast.error(message);
      console.error(message, error);
      return message;
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

    const resetNewUserForm = () => {
      newUser.value = {
        Username: '',
        Password: '',
        FirstName: '',
        LastName: '',
        Email: '',
        RoleID: '2',
      };
      passwordStrength.value = 0;
      passwordErrors.create = [];
      hasMinLength.value = false;
      hasLetter.value = false;
      hasNumber.value = false;
      hasSpecialChar.value = false;
    };

    watch(showCreateUserModal, (newVal) => {
      if (!newVal) {
        resetNewUserForm();
      }
    });
    const validatePassword = (password, type) => {
      passwordErrors[type] = [];
      if (!password) {
        passwordStrength.value = 0;
        return;
      }

      let strength = 0;

      if (password.length >= 8) {
        strength++;
      } else {
        passwordErrors[type].push(t('passwordValidation.minLength'));
      }

      if (/[a-zA-Z]/.test(password)) {
        strength++;
      } else {
        passwordErrors[type].push(t('passwordValidation.letter'));
      }

      if (/\d/.test(password)) {
        strength++;
      } else {
        passwordErrors[type].push(t('passwordValidation.digit'));
      }

      if (/[^a-zA-Z0-9\s]/.test(password)) {
        strength++;
      } else {
        passwordErrors[type].push(t('passwordValidation.specialChar'));
      }

      passwordStrength.value = strength;
    };


    const resetEditPasswordFields = () => {
      if (editUserData.value) {
        editUserData.value.Password = '';
      }
      isPasswordVisible.value = false;
      passwordErrors.edit = [];
      passwordStrength.value = 0;
    };
    const togglePasswordVisibility = (type) => {
      isPasswordVisible[type] = !isPasswordVisible[type];
    };

    const filters = ref({
      sender: '',
      assignee: '',
      status: ''
    });

    const newUser = ref({
      Username: '',
      Password: '',
      FirstName: '',
      LastName: '',
      Email: '',
      RoleID: '2',
    });

    const confirmModal = ref({
      show: false,
      type: '',
      id: null,
      name: '',
      callback: null
    });

    const formatDate = (dateString) => {
      if (!dateString) return 'Немає дати';
      return format(new Date(dateString), 'dd.MM.yyyy HH:mm');
    };

    const fetchTickets = async (forceUseFilters = false) => {
      try {
        const shouldUseFilters = forceUseFilters && (
            filters.value.sender?.trim() ||
            filters.value.assignee?.trim() ||
            filters.value.status?.trim()
        );

        const response = shouldUseFilters
            ? await adminApi.getFilteredTickets(filters.value)
            : await adminApi.getAllTickets();

        tickets.value = (response.data || []).sort((a, b) =>
            new Date(b.CreatedAt) - new Date(a.CreatedAt)
        );
        filteredTickets.value = [...tickets.value];
        currentPage.value = 1;
      } catch (error) {
        toast.error(t('ticketFetchError'));
        console.error("Помилка при завантаженні тікетів:", error);
        tickets.value = [];
        filteredTickets.value = [];
      }
    };
    const applyFilters = async () => {
      await fetchTickets(true);
    };

    const toggleDescription = (ticketId) => {
      expandedDescriptions.value = {
        ...expandedDescriptions.value,
        [ticketId]: !expandedDescriptions.value[ticketId]
      };
    };

    const truncateDescription = (description, ticketId) => {
      if (!description) return '';
      const shouldTruncate = description.length > 100 && !expandedDescriptions.value[ticketId];
      return shouldTruncate ? `${description.substring(0, 100)}...` : description;
    };

    const getStatusClass = (status) => {
      switch (status) {
        case 'Новий': return 'status-new';
        case 'Оновлено': return 'status-in-progress';
        default: return '';
      }
    };
    const resetFilters = async () => {
      filters.value = {
        sender: null,
        assignee: null,
        status: ''
      };
      await fetchTickets(false);
    };

    const paginatedUsers = computed(() => {
      const start = (currentUsersPage.value - 1) * itemsPerPage.value;
      const end = start + itemsPerPage.value;
      return users.value.slice(start, end);
    });

    const totalUsersPages = computed(() =>
        Math.ceil(users.value.length / itemsPerPage.value)
    );

    const paginatedNotifications = computed(() => {
      const start = (currentNotificationsPage.value - 1) * itemsPerPage.value;
      const end = start + itemsPerPage.value;
      return notifications.value.slice(start, end);
    });

    const totalNotificationsPages = computed(() =>
        Math.ceil(notifications.value.length / itemsPerPage.value)
    );

    const nextUsersPage = () => {
      if (currentUsersPage.value < totalUsersPages.value) currentUsersPage.value++;
    };

    const prevUsersPage = () => {
      if (currentUsersPage.value > 1) currentUsersPage.value--;
    };

    const nextNotificationsPage = () => {
      if (currentNotificationsPage.value < totalNotificationsPages.value) currentNotificationsPage.value++;
    };

    const prevNotificationsPage = () => {
      if (currentNotificationsPage.value > 1) currentNotificationsPage.value--;
    };
    const totalPages = computed(() =>
        Math.ceil(filteredTickets.value.length / itemsPerPage.value)
    );

    const paginatedTickets = computed(() => {
      const start = (currentPage.value - 1) * itemsPerPage.value;
      const end = start + itemsPerPage.value;
      return filteredTickets.value.slice(start, end);
    });

    const nextPage = () => {
      if (currentPage.value < totalPages.value) currentPage.value++;
    };

    const prevPage = () => {
      if (currentPage.value > 1) currentPage.value--;
    };

    const fetchNotifications = async () => {
      try {
        const response = await adminApi.getAllNotifications();
        notifications.value = (response.data || []).sort((a, b) =>
            new Date(b.CreatedAt) - new Date(a.CreatedAt)
        );
      } catch (error) {
        console.error("Помилка при завантаженні сповіщень:", error);
        handleApiError(error, 'notificationFetchError');
      }
    };

    const fetchUsers = async () => {
      try {
        const response = await adminApi.getAllUsers();
        users.value = (response.data || []).sort((a, b) =>
            new Date(b.CreatedAt) - new Date(a.CreatedAt)
        );
      } catch (error) {
        console.error("Помилка при завантаженні користувачів:", error);
        handleApiError(error, 'userFetchError');
      }
    };

    const confirmDelete = (type, id, name) => {
      confirmModal.value = {
        show: true,
        type,
        id,
        name,
        callback: type === 'user' ? deleteUser :
            type === 'ticket' ? deleteTicket :
                deleteNotification
      };
    };

    const executeDelete = async () => {
      try {
        await confirmModal.value.callback(confirmModal.value.id);
        confirmModal.value.show = false;
        if (confirmModal.value.type === 'user') await fetchUsers();
        else if (confirmModal.value.type === 'ticket') await fetchTickets();
        else await fetchNotifications();
      } catch (error) {
        console.error("Помилка при видаленні:", error);
        handleApiError(error, 'delete Error');
      }
    };

    const deleteTicket = async (ticketID) => {
      try {
        await adminApi.adminDeleteTicket(ticketID);
        tickets.value = tickets.value.filter(ticket => ticket.TicketID !== ticketID);
        toast.success(t('ticketDeleteSuccess'));
      } catch (error) {
        toast.error(t('ticketDeleteError'));
        console.error("Помилка при видаленні тікета:", error);
        handleApiError(error, 'delete Error');
      }
    };

    const deleteUser = async (userID) => {
      try {
      await adminApi.deleteUser(userID);
      users.value = users.value.filter(user => user.UserID !== userID);
        toast.success(t('userDeleteSuccess'), {
          icon: "fas fa-check-circle",
          timeout: 3000,
        });
      }
      catch (error){
        toast.error(t('userDeleteError'), {
          icon: "fas fa-check-circle",
          timeout: 3000,
        });
      }
    };

    const deleteNotification = async (notificationID) => {
      try {
        await adminApi.deleteNotification(notificationID);
        notifications.value = notifications.value.filter(notification => notification.NotificationID !== notificationID);
        toast.success(t('notificationMarkedRead'), {
          icon: "fas fa-check-circle",
          timeout: 3000,
        });
      }
      catch (error){
        toast.error(t('notificationMarkError'), {
          icon: "fas fa-check-circle",
          timeout: 3000,
        });
      }
    };

    const editUser = (user) => {
      editUserData.value = {
        ...user,
        Password: ''
      };
      resetEditPasswordFields();
    };

    const submitEditUser = async () => {
      if (editUserData.value.Password) {
        validatePassword(editUserData.value.Password, 'edit');
        if (passwordErrors.edit.length > 0) {
          toast.error(t('passwordChangeError:\n') + passwordErrors.edit.join("\n"));
          return;
        }
      }

      try {
        const userDataToSend = { ...editUserData.value };

        if (!userDataToSend.Password) {
          delete userDataToSend.Password;
        }

        userDataToSend.RoleID = parseInt(userDataToSend.RoleID);

        await adminApi.updateUser(editUserData.value.UserID, userDataToSend);
        editUserData.value = null;
        toast.success(t('userUpdateSuccess'));
        await fetchUsers();
      } catch (error) {
        console.error("Помилка при оновленні користувача:", error);
        handleApiError(error, 'userUpdateError');
      }
    };

    watch(editUserData, (newVal) => {
      if (!newVal) {
        resetEditPasswordFields();
      }
    }, { deep: true });
    watch(() => editUserData.value?.Password, (newPassword) => {
      validatePassword(newPassword, 'edit');
    });
    watch(() => newUser.value?.Password, (newPassword) => {
      validatePassword(newPassword, 'create');
    })
    const submitCreateUser = async () => {
      validatePassword(newUser.value.Password, 'create');
      if (passwordErrors.create.length > 0) {
        toast.error(passwordErrors.create.join("\n"));
        return;
      }

      try {
        await adminApi.createUser({
          ...newUser.value,
          RoleID: parseInt(newUser.value.RoleID)
        });

        showCreateUserModal.value = false;
        newUser.value = {
          Username: '',
          Password: '',
          FirstName: '',
          LastName: '',
          Email: '',
          RoleID: '2'
        };

        passwordStrength.value = 0;
        hasMinLength.value = false;
        hasLetter.value = false;
        hasNumber.value = false;
        hasSpecialChar.value = false;
        passwordErrors.value = [];

        toast.success(t('userCreateSuccess'));
        await fetchUsers();
      } catch (error) {
        console.error('Помилка при створенні користувача:', error);
        handleApiError(error,'userCreateError');
      }
    };
    const backupData = async () => {
      try {
        const response = await adminApi.backupData();

        const blob = new Blob([response.data], { type: 'application/sql' });
        const downloadUrl = URL.createObjectURL(blob);

        const link = document.createElement('a');
        link.href = downloadUrl;

        const contentDisposition = response.headers['content-disposition'];
        let filename = 'backup_' + new Date().toISOString().slice(0, 10) + '.sql';

        if (contentDisposition) {
          const filenameMatch = contentDisposition.match(/filename="?(.+)"?/);
          if (filenameMatch && filenameMatch[1]) {
            filename = filenameMatch[1];
          }
        }

        link.download = filename;
        document.body.appendChild(link);
        link.click();

        setTimeout(() => {
          document.body.removeChild(link);
          URL.revokeObjectURL(downloadUrl);
        }, 100);
        toast.success(t('backupSuccess'));
      } catch (error) {
        toast.error(t('backupError') + ": " + error.message);
      }
    };


    const restoreData = async () => {
      const fileInput = document.createElement('input');
      fileInput.type = 'file';
      fileInput.accept = '.sql';

      fileInput.onchange = async (e) => {
        const file = e.target.files[0];
        if (!file) return;
        if (!file.name.endsWith('.sql')) {
          toast.error(t('invalidFileFormat'));
          return;
        }
        try {
          const formData = new FormData();
          formData.append('file', file);
          await adminApi.restoreData(formData);
          toast.success(t('restoreSuccess'));
          await Promise.all([fetchUsers(), fetchTickets(), fetchNotifications()]);
        } catch (error) {
          handleApiError(error, 'restoreError');
        }
      };

      fileInput.click();
    };

    const exportData = async () => {
      try {
        const data = await adminApi.exportData();
        const blob = new Blob([data], {type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'});
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = `export_${new Date().toISOString().slice(0, 10)}.xlsx`;
        link.click();
        toast.success(t('exportSuccess'));
      } catch (error) {
        console.error('Помилка при експорті даних:', error);
        toast.error(t('exportError'));
      }
    };

    const importData = async () => {
      const fileInput = document.createElement('input');
      fileInput.type = 'file';
      fileInput.accept = '.xlsx,.xls';

      fileInput.onchange = async (e) => {
        const file = e.target.files[0];
        if (file) {
          try {
            await adminApi.importData(file);
            toast.success(t('importSuccess'));
            await Promise.all([fetchUsers(), fetchTickets(), fetchNotifications()]);
          } catch (error) {
            console.error('Помилка при імпорті даних:', error);
            toast.error(t('importError'));
          }
        }
      };

      fileInput.click();
    };

    onMounted(() => {
      Promise.all([fetchUsers(), fetchTickets(), fetchNotifications()]);
    });

    return {
      tickets,
      notifications,
      users,
      editUserData,
      newUser,
      showCreateUserModal,
      confirmModal,
      authStore,
      formatDate,
      confirmDelete,
      executeDelete,
      editUser,
      submitEditUser,
      submitCreateUser,

      closeModalOnOutsideClick,
      filters,
      statusOptions,
      filteredTickets,
      applyFilters,
      resetFilters,
      getStatusClass,
      toggleDescription,
      truncateDescription,
      paginatedTickets,
      totalPages,
      currentPage,
      nextPage,
      prevPage,
      paginatedUsers,
      currentUsersPage,
      totalUsersPages,
      nextUsersPage,
      prevUsersPage,
      paginatedNotifications,
      currentNotificationsPage,
      totalNotificationsPages,
      nextNotificationsPage,
      prevNotificationsPage,

      backupData,
      restoreData,
      exportData,
      importData,

      expandedDescriptions,
      itemsPerPage,
      isPasswordVisible,
      togglePasswordVisibility,

      passwordErrors,
      passwordStrength,
      hasMinLength,
      hasLetter,
      hasNumber,
      hasSpecialChar,

      setLocale,
      currentLocale,
      messageMap
    };
  }
};
</script>

<style scoped>

.admin-dashboard {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  min-height: 100vh;
  width: 100%;
  padding: 1rem;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}
.admin-dashboard-card {
  width: 100%;
  max-width: 1200px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
  padding: 2rem;
  margin: auto;
}
.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #eaeaea;
}
h1 {
  font-size: 1.8rem;
  color: #2c3e50;
  font-weight: 600;
}

.section-header {
  margin-bottom: 1.2rem;
}

.section-header h2 {
  font-size: 1.4rem;
  color: #2c3e50;
  font-weight: 600;
}
.empty-state {
  text-align: center;
  padding: 2.5rem;
  background: #f9fbfd;
  border-radius: 10px;
  margin: 1.5rem 0;
  border: 1px dashed #e1e8f0;
}

.empty-icon {
  font-size: 2.5rem;
  color: #a0aec0;
  margin-bottom: 1rem;
}
.password-strength {
  margin-top: 0.8rem;
}

.strength-bar {
  height: 6px;
  border-radius: 3px;
  background: #e0e0e0;
  overflow: hidden;
  position: relative;
}

.strength-bar::after {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  transition: width 0.3s ease;
}

.strength-bar.weak::after {
  width: 25%;
  background: #ff3860;
}

.strength-bar.medium::after {
  width: 50%;
  background: #ffdd57;
}

.strength-bar.strong::after {
  width: 100%;
  background: #23d160;
}

.strength-labels {
  display: flex;
  justify-content: space-between;
  margin-top: 0.4rem;
  font-size: 0.75rem;
}

.strength-labels span {
  color: #b5b5b5;
  flex: 1;
  text-align: center;
}

.strength-labels span.active {
  color: #2d3748;
  font-weight: 500;
}

.error-messages {
  margin-top: 0.5rem;
}

.error-text {
  color: #ff3860;
  font-size: 0.85rem;
  margin: 0.2rem 0;
}
.empty-state p {
  color: #718096;
  font-size: 1.05rem;
}

.password-strength {
  margin-top: 0.8rem;
}

.strength-bar {
  height: 6px;
  border-radius: 3px;
  background: #e0e0e0;
  overflow: hidden;
  position: relative;
}

.strength-bar::after {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  transition: width 0.3s ease;
}

.strength-bar.weak::after {
  width: 33%;
  background: #ff3860;
}

.strength-bar.medium::after {
  width: 66%;
  background: #ffdd57;
}

.strength-bar.strong::after {
  width: 100%;
  background: #23d160;
}

.strength-labels {
  display: flex;
  justify-content: space-between;
  margin-top: 0.4rem;
  font-size: 0.75rem;
}

.strength-labels span {
  color: #b5b5b5;
}

.strength-labels span.active {
  color: var(--dark-color);
  font-weight: 500;
}

.password-rules {
  margin-top: 0.8rem;
  padding-left: 1.2rem;
}

.password-rules li {
  color: #b5b5b5;
  font-size: 0.85rem;
  margin-bottom: 0.3rem;
  list-style-type: none;
  display: flex;
  align-items: center;
}

.password-rules li.valid {
  color: #23d160;
}

.password-rules li i {
  margin-right: 0.5rem;
  font-size: 0.9rem;
}
.table-container {
  overflow-x: auto;
  background: white;
  border-radius: 10px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
  margin-bottom: 1.8rem;
  border: 1px solid #edf2f7;
}

.data-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
}

.data-table th {
  background: #f8f9fa;
  color: #4a5568;
  padding: 10px 12px;
  text-align: left;
  font-weight: 600;
  font-size: 0.9rem;
  border-bottom: 2px solid #e2e8f0;
}

.data-table th,
.data-table td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid #edf2f7;
}

.data-table tr:hover {
  background-color: #f8fafc;
}

.role-badge {
  display: inline-block;
  padding: 0.35rem 0.65rem;
  border-radius: 50px;
  font-size: 0.8rem;
  font-weight: 500;
}

.role-badge.admin {
  background-color: rgba(244, 63, 94, 0.1);
  color: #f43f5e;
}

.role-badge.user {
  background-color: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.status-badge {
  display: inline-block;
  border-radius: 4px;
  font-weight: 500;
  font-size: 0.8rem;
}

.status-new {
  background-color: #ebf8ff;
  color: #3182ce;
}

.status-in-progress {
  background-color: #fff5eb;
  color: #dd6b20;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.6rem 1.2rem;
  border-radius: 6px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

button i {
  font-size: 0.9rem;
}

.create-btn {
  position: fixed;
  bottom: 1.5rem;
  right: 1.5rem;
  background: #4361ee;
  color: white;
  border: none;
  border-radius: 50px;
  padding: 0.9rem 1.4rem;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.6rem;
  box-shadow: 0 4px 10px rgba(67, 97, 238, 0.25);
  z-index: 100;
}

.create-btn:hover {
  background: #3a56e0;
  transform: translateY(-3px);
  box-shadow: 0 6px 15px rgba(67, 97, 238, 0.3);
}
.action-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.action-btn.edit {
  background: #e6f7ff;
  color: #1890ff;
}

.action-btn.delete {
  background: #fff1f0;
  color: #f5222d;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}
.data-management-section h2{
  padding-top: 30px;
}
.data-management-section {
  margin-top: 3rem;
  padding-top: 2rem;
  border-top: 1px solid var(--border-color);
}

.data-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-top: 1.5rem;
}

.data-card {
  background: var(--white);
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: center;
  border: 2px solid transparent;
}

.data-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
}

.data-card:hover .card-icon {
  transform: scale(1.1);
}

.card-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.2rem;
  font-size: 2rem;
  transition: all 0.3s ease;
}

.card-icon.backup {
  background-color: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.card-icon.restore {
  background-color: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
}

.card-icon.export {
  background-color: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.card-icon.import {
  background-color: rgba(139, 92, 246, 0.15);
  color: #8b5cf6;
}

.data-card h3 {
  font-size: 1.3rem;
  margin-bottom: 0.5rem;
  color: var(--dark-color);
}

.data-card p {
  color: var(--gray-color);
  font-size: 0.95rem;
  line-height: 1.5;
}
.backup {
  background-color: #10b981;
  color: white;
}

.backup:hover {
  background-color: #0d9e6e;
}

.restore {
  background-color: #3b82f6;
  color: white;
}

.restore:hover {
  background-color: #2563eb;
}

.export {
  background-color: #f59e0b;
  color: white;
}

.export:hover {
  background-color: #d97706;
}

.import {
  background-color: #8b5cf6;
  color: white;
}

.import:hover {
  background-color: #7c3aed;
}

.data-actions {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  margin-bottom: 2rem;
}

.notifications-list {
  background: white;
  border-radius: 10px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
  padding: 1rem;
  border: 1px solid #edf2f7;
}

.notification-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.9rem;
  border-bottom: 1px solid #f0f4f8;
}

.notification-item:last-child {
  border-bottom: none;
}


.notification-content {
  flex: 1;
}
.notification-message {
  font-weight: 500;
  margin-bottom: 0.2rem;
  color: #2d3748;
}

.notification-time {
  font-size: 0.8rem;
  color: #a0aec0;
}
.password-wrapper {
  position: relative;
}

.eye-button {
  position: absolute;
  top: 50%;
  right: 12px;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: var(--gray-color);
  font-size: 1rem;
  cursor: pointer;
  padding: 0.5rem;
  transition: color 0.2s ease;
}

.eye-button:hover,
.eye-button.active {
  color: var(--primary-color);
}
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-active .modal,
.modal-fade-leave-active .modal {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.modal-fade-enter-from .modal,
.modal-fade-leave-to .modal {
  transform: translateY(-20px);
  opacity: 0;
}

.form-group input,
.form-group select,
.form-group textarea {
  border: 1px solid #ced4da !important;
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  border-color: #4361ee !important;
  box-shadow: 0 0 0 0.2rem rgba(67, 97, 238, 0.25) !important;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.4);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-fade-enter-active .modal,
.modal-fade-leave-active .modal,
.modal-fade-enter-active .confirm-modal,
.modal-fade-leave-active .confirm-modal {
  transition: all 0.3s ease;
}

.modal-fade-enter-from .modal,
.modal-fade-leave-to .modal,
.modal-fade-enter-from .confirm-modal,
.modal-fade-leave-to .confirm-modal {
  transform: translateY(-20px);
  opacity: 0;
}
.modal, .confirm-modal {
  background-color: white;
  border-radius: 14px;
  width: 100%;
  max-width: 480px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.confirm-modal {
  max-width: 380px;
}

.modal-header {
  padding: 1.2rem 1.5rem;
  border-bottom: 1px solid #f0f4f8;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f8fafc;
}
.modal-header h3 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #2c3e50;
}

.close-btn {
  background: none;
  border: none;
  color: #a0aec0;
  font-size: 1.1rem;
  cursor: pointer;
  padding: 0.2rem;
  transition: color 0.3s ease;
}

.close-btn:hover {
  color: #f5222d;
}

.modal-body {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.25rem;
  border-color: #1ebeff;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: var(--dark-color);
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(67, 97, 238, 0.1);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}


.cancel-btn:hover {
  background-color: #420001;
}

.submit-btn:hover {
  background-color: #074200;
}

.confirm-btn:hover {
  background-color: #e5177e;
}

.modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}
.logout-btn {
  background: #4361ee;
  color: white;
  padding: 0.7rem 1.2rem;
  border: none;
  border-radius: 6px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  box-shadow: 0 2px 5px rgba(67, 97, 238, 0.2);
}

.logout-btn:hover {
  background: #3a56e0;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(67, 97, 238, 0.25);
}

.ticket-filters {
  display: flex;
  gap: 1rem;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 1.5rem;
  padding: 1rem;
  background-color: var(--white);
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-group label {
  font-weight: 500;
  color: var(--dark-color);
  white-space: nowrap;
}

.filter-group input,
.filter-group select {
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  min-width: 150px;
}

.reset-btn {
  background-color: var(--gray-color);
  color: var(--white);
  padding: 0.5rem 1rem;
}

.reset-btn:hover {
  background-color: #5a6268;
}
.ticket-description {
  max-width: 300px;
  position: relative;
}

.description-text {
  display: inline;
  word-wrap: break-word;
}

.show-more-btn {
  background: none;
  border: none;
  color: #4299e1;
  cursor: pointer;
  font-size: 0.8rem;
  padding: 0.2rem 0.5rem;
  margin-left: 0.5rem;
  white-space: nowrap;
  font-weight: 600;
}

.pagination-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 1.2rem;
  gap: 1rem;
}

.pagination-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: #edf2f7;
  color: #4a5568;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-btn:hover:not(:disabled) {
  background: #e2e8f0;
  transform: translateY(-1px);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.page-info {
  font-size: 0.9rem;
  color: #718096;
}

.confirm-btn, .submit-btn {
  background-color: #5a6268;
  color: gray;
  padding: 0.75rem 1.5rem;
  border-radius: 6px;
  border: none;
  font-weight: 600;
}
button {
  font-family: 'Roboto', sans-serif;
  font-weight: 500;
  transition: all 0.2s ease;
}

button:disabled {
  background-color: var(--gray-color) !important;
  color: var(--light-color) !important;
  cursor: not-allowed;
  opacity: 0.7;
}

.cancel-btn {
  background: #edf2f7;
  color: #4a5568;
  padding: 0.7rem 1.2rem;
  border: none;
  border-radius: 6px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.submit-btn, .confirm-btn {
  padding: 0.7rem 1.4rem;
  border: none;
  border-radius: 6px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.submit-btn {
  background: #4361ee;
  color: white;
}

.confirm-btn {
  background: #f5222d;
  color: white;
}

.cancel-btn:hover {
  background: #e2e8f0;
}

.submit-btn:hover {
  background: #3a56e0;
}

.confirm-btn:hover {
  background: #cf1322;
}

.submit-btn:disabled, .confirm-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .admin-dashboard {
    padding: 1rem;
  }

  .data-actions {
    flex-direction: column;
  }

  .action-btn {
    width: 100%;
  }

  .modal {
    width: 95%;
    user-select: none;
  }

  .form-input, .form-group label, .password-wrapper {
    user-select: text;
  }
  .admin-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  .ticket-filters {
    flex-direction: column;
    align-items: flex-start;
  }
  .ticket-description {
    max-width: 200px;
  }
  .filter-group {
    width: 100%;
  }

  .filter-group input,
  .filter-group select {
    width: 100%;
  }
  .reset-btn {
    margin-left: 0;
    width: 100%;
  }

  .ticket-description {
    max-width: 200px;
  }
  .logout-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>