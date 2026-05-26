<template>
  <div class="dashboard-container">
    <div class="dashboard-card">
      <header class="dashboard-header">
        <div class="language-switcher">
          <button @click="setLocale('ua')" :class="{ active: currentLocale === 'ua' }">{{ $t('ukrainian') }}</button>
          <button @click="setLocale('en')" :class="{ active: currentLocale === 'en' }">{{ $t('english') }}</button>
        </div>
        <h1>{{ $t('dashboardTitle') }}</h1>
        <button class="logout-btn" @click="authStore.logout()">
          <i class="fas fa-sign-out-alt"></i> {{ $t('logout') }}
        </button>
      </header>

      <!-- Tickets Section -->
      <section class="tickets-section" v-if="tickets.length > 0">
        <div class="section-header">
          <h2>{{ $t('tickets') }}</h2>
        </div>

        <div class="table-container">
          <table class="data-table">
            <thead>
            <tr>
              <th>{{ $t('title') }}</th>
              <th>{{ $t('description') }}</th>
              <th>{{ $t('status') }}</th>
              <th>{{ $t('assignee') }}</th>
              <th>{{ $t('createdAt') }}</th>
              <th>{{ $t('actions') }}</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="ticket in paginatedTickets" :key="ticket.TicketID">
              <td>{{ truncateTitle(ticket.Title) }}</td>
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
              <td>{{ ticket.AssigneeUsername || $t('notAssigned') }}</td>
              <td>{{ formatDate(ticket.CreatedAt) }}</td>
              <td class="actions">
                <button class="action-btn edit" @click="editTicket(ticket)" :title="$t('edit')">
                  <i class="fas fa-edit"></i>
                </button>
                <button class="action-btn delete" @click="confirmDelete('ticket', ticket.TicketID, ticket.Title)" :title="$t('delete')">
                  <i class="fas fa-trash-alt"></i>
                </button>
              </td>
            </tr>
            </tbody>
          </table>

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
      </section>

      <div class="empty-state" v-else>
        <i class="fas fa-ticket-alt empty-icon"></i>
        <p>{{ $t('noActiveTickets') }}</p>
      </div>

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
        <p>{{ $t('noNewNotifications') }}</p>
      </div>

      <!-- Create Ticket Button -->
      <button class="create-btn floating-btn" @click="showCreateTicketModal = true">
        <i class="fas fa-plus"></i> {{ $t('createTicket') }}
      </button>

      <!-- Modals -->
      <transition name="modal-fade">
        <div v-if="showCreateTicketModal" class="modal-overlay" @click.self="showCreateTicketModal = false">
          <div class="modal">
            <div class="modal-header">
              <h3>{{ $t('createTicketTitle') }}</h3>
              <button class="close-btn" @click="showCreateTicketModal = false">
                <i class="fas fa-times"></i>
              </button>
            </div>
            <div class="modal-body">
              <form @submit.prevent="addTicket">
                <div class="input-group">
                  <input
                      v-model="newTicket.Title"
                      type="text"
                      id="ticket-title"
                      required
                      class="form-input"
                      :placeholder="$t('ticketTitlePlaceholder')"
                      maxlength="60"
                  />
                  <label for="ticket-title">{{ $t('ticketTitle') }}</label>
                  <div class="input-border"></div>
                </div>

                <div class="input-group">
                  <textarea
                      v-model="newTicket.Description"
                      id="ticket-description"
                      required
                      class="form-input"
                      :placeholder="$t('problemDescription')"
                      rows="4"
                  ></textarea>
                  <label for="ticket-description">{{ $t('problemDescription') }}</label>
                  <div class="input-border"></div>
                </div>

                <div class="input-group">
                  <select
                      v-model="newTicket.AssignedTo"
                      id="ticket-assignee"
                      required
                      class="form-input"
                  >
                    <option v-for="user in filteredUsers" :key="user.UserID" :value="user.UserID">
                      {{ user.FirstName }} {{ user.LastName }} ({{ user.Username }})
                    </option>
                  </select>
                  <label for="ticket-assignee">{{ $t('assign') }}</label>
                  <div class="input-border"></div>
                </div>

                <div class="form-actions">
                  <button type="button" class="cancel-btn" @click="showCreateTicketModal = false">{{ $t('cancel') }}</button>
                  <button type="submit" class="submit-btn" :disabled="isLoading">
                    <span v-if="!isLoading">{{ $t('createButton') }}</span>
                    <span v-else class="loading-spinner"></span>
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </transition>

      <transition name="modal-fade">
        <div v-if="editTicketData" class="modal-overlay" @click.self="editTicketData = null">
          <div class="modal">
            <div class="modal-header">
              <h3>{{ $t('editTicketTitle') }}</h3>
              <button class="close-btn" @click="editTicketData = null">
                <i class="fas fa-times"></i>
              </button>
            </div>
            <div class="modal-body">
              <form @submit.prevent="updateTicket">
                <div class="input-group">
                  <input
                      v-model="editTicketData.Title"
                      type="text"
                      id="edit-ticket-title"
                      required
                      class="form-input"
                      :placeholder="$t('ticketTitlePlaceholder')"
                      maxlength="60"
                  />
                  <label for="edit-ticket-title">{{ $t('ticketTitle') }}</label>
                  <div class="input-border"></div>
                </div>

                <div class="input-group">
                  <textarea
                      v-model="editTicketData.Description"
                      id="edit-ticket-description"
                      required
                      class="form-input"
                      :placeholder="$t('problemDescription')"
                      rows="4"
                  ></textarea>
                  <label for="edit-ticket-description">{{ $t('problemDescription') }}</label>
                  <div class="input-border"></div>
                </div>

                <div class="input-group">
                  <select
                      v-model="editTicketData.AssignedTo"
                      id="edit-ticket-assignee"
                      required
                      class="form-input"
                  >
                    <option v-for="user in filteredUsers" :key="user.UserID" :value="user.UserID">
                      {{ user.Username }}
                    </option>
                  </select>
                  <label for="edit-ticket-assignee">{{ $t('assign') }}</label>
                  <div class="input-border"></div>
                </div>

                <div class="form-actions">
                  <button type="button" class="cancel-btn" @click="editTicketData = null">{{ $t('cancel') }}</button>
                  <button type="submit" class="submit-btn" :disabled="isLoading">
                    <span v-if="!isLoading">{{ $t('save') }}</span>
                    <span v-else class="loading-spinner"></span>
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </transition>

      <transition name="modal-fade">
        <div v-if="confirmModal.show" class="modal-overlay" @click.self="confirmModal.show = false">
          <div class="confirm-modal">
            <div class="modal-header">
              <h3>{{ $t('confirmationTitle') }}</h3>
            </div>
            <div class="modal-body">
              <p>{{ $t('confirmationMessage', {
                type: $t(confirmModal.type + 'Deleted'),
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
import { ref, onMounted, computed } from "vue";
import { requestApi, notificationApi } from "../api";
import { useAuthStore } from "../store/auth";
import { format } from 'date-fns';
import i18n from '../i18n';
import { useI18n } from 'vue-i18n';
import { useToast } from "vue-toastification";


export default {
  name: "TicketDashboard",
  setup() {
    const tickets = ref([]);
    const currentPage = ref(1);
    const notifications = ref([]);
    const users = ref([]);
    const authStore = useAuthStore();
    const expandedDescriptions = ref({});
    const itemsPerPage = ref(10);
    const currentNotificationsPage = ref(1);
    const isLoading = ref(false);
    const toast = useToast();
    const { t, locale } = useI18n({ useScope: 'global' });
    const currentLocale = computed(() => locale.value);

    const messageMap = {
      'Оновлено тікет': 'notificationMessages.ticketUpdated',
      'Створено новий тікет': 'notificationMessages.ticketCreated',
      'Новий':'statusMessages.statusCreated',
      'Оновлено':'statusMessages.statusUpdated'
    };

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

    const filteredUsers = computed(() => {
      return users.value.filter(user => user.RoleID !== 1);
    });

    const truncateTitle = (title) => {
      if (!title) return '';
      return title.length > 30 ? `${title.substring(0, 57)}...` : title;
    };

    const newTicket = ref({
      Title: '',
      Description: '',
      AssignedTo: null,
    });
    const editTicketData = ref(null);
    const showCreateTicketModal = ref(false);

    const confirmModal = ref({
      show: false,
      type: '',
      id: null,
      name: '',
      callback: null
    });

    const formatDate = (dateString) => {
      return format(new Date(dateString), 'dd.MM.yyyy HH:mm');
    };

    const availableUsers = computed(() => {
      return users.value.filter(user => user.RoleID !== 1);
    });

    const fetchUsers = async () => {
      try {
        const response = await requestApi.getAllUsers();
        users.value = response.data || [];
      } catch (error) {
        console.error("Помилка при завантаженні користувачів:", error);
        toast.error(t('userFetchError'), {
          icon: "fas fa-exclamation-triangle",
          timeout: 5000,
        });
      }
    };

    const fetchTickets = async () => {
      isLoading.value = true;
      try {
        const response = await requestApi.getUserTickets();
        tickets.value = (response.data || []).sort((a, b) =>
            new Date(b.CreatedAt) - new Date(a.CreatedAt)
        );
      } catch (error) {
        console.error("Помилка при завантаженні тікетів:", error);
        toast.error(t('ticketFetchError'), {
          icon: "fas fa-exclamation-triangle",
          timeout: 5000,
        });
        tickets.value = [];
      } finally {
        isLoading.value = false;
      }
    };

    const truncateDescription = (description, ticketId) => {
      if (!description) return '';
      const shouldTruncate = description.length > 100 && !expandedDescriptions.value[ticketId];
      return shouldTruncate ? `${description.substring(0, 100)}...` : description;
    };

    const toggleDescription = (ticketId) => {
      expandedDescriptions.value = {
        ...expandedDescriptions.value,
        [ticketId]: !expandedDescriptions.value[ticketId]
      };
    };

    const paginatedTickets = computed(() => {
      const start = (currentPage.value - 1) * itemsPerPage.value;
      const end = start + itemsPerPage.value;
      return tickets.value.slice(start, end);
    });

    const nextPage = () => {
      if (currentPage.value < totalPages.value) currentPage.value++;
    };

    const totalPages = computed(() => Math.ceil(tickets.value.length / itemsPerPage.value));

    const prevPage = () => {
      if (currentPage.value > 1) currentPage.value--;
    };

    const fetchNotifications = async () => {
      isLoading.value = true;
      try {
        const response = await notificationApi.getUserNotifications();
        notifications.value = (response.data || []).sort((a, b) =>
            new Date(b.CreatedAt) - new Date(a.CreatedAt)
        );
      } catch (error) {
        console.error("Помилка при завантаженні сповіщень:", error);
        toast.error(t('notificationFetchError'), {
          icon: "fas fa-exclamation-triangle",
          timeout: 5000,
        });
        notifications.value = [];
      } finally {
        isLoading.value = false;
      }
    };

    const paginatedNotifications = computed(() => {
      const start = (currentNotificationsPage.value - 1) * itemsPerPage.value;
      const end = start + itemsPerPage.value;
      return notifications.value.slice(start, end);
    });

    const totalNotificationsPages = computed(() =>
        Math.ceil(notifications.value.length / itemsPerPage.value)
    );

    const nextNotificationsPage = () => {
      if (currentNotificationsPage.value < totalNotificationsPages.value) currentNotificationsPage.value++;
    };

    const prevNotificationsPage = () => {
      if (currentNotificationsPage.value > 1) currentNotificationsPage.value--;
    };

    const confirmDelete = (type, id, name) => {
      confirmModal.value = {
        show: true,
        type,
        id,
        name,
        callback: type === 'ticket' ? deleteTicket : deleteNotification
      };
    };

    const executeDelete = async () => {
      try {
        await confirmModal.value.callback(confirmModal.value.id);
        confirmModal.value.show = false;
        await fetchTickets();
        await fetchNotifications();

        toast.success(t('deleteSuccess'), {
          icon: "fas fa-check-circle",
          timeout: 3000,
        });
      } catch (error) {
        console.error("Помилка при видаленні:", error);
        toast.error(t('deleteError'), {
          icon: "fas fa-exclamation-triangle",
          timeout: 5000,
        });
      }
    };

    const getStatusClass = (status) => {
      switch (status) {
        case 'Новий': return 'status-new';
        case 'Оновлено': return 'status-in-progress';
        default: return '';
      }
    };

    const deleteTicket = async (ticketID) => {
      try {
        await requestApi.deleteTicket(ticketID);
        tickets.value = tickets.value.filter(ticket => ticket.TicketID !== ticketID);
      } catch (error) {
        console.error("Помилка при видаленні тікета:", error);
        toast.error(t('ticketDeleteError'), {
          icon: "fas fa-exclamation-triangle",
          timeout: 5000,
        });
      }
    };

    const deleteNotification = async (notificationID) => {
      try {
        await notificationApi.markAsRead(notificationID);
        notifications.value = notifications.value.filter(notification => notification.NotificationID !== notificationID);
      } catch (error) {
        console.error("Помилка при видаленні сповіщення:", error);
        toast.error(t('notificationDeleteError'), {
          icon: "fas fa-exclamation-triangle",
          timeout: 5000,
        });
      }
    };

    const addTicket = async () => {
      isLoading.value = true;
      try {
        const assignedUser = users.value.find(u => u.UserID === newTicket.value.AssignedTo);
        if (assignedUser && assignedUser.RoleID === 1) {
          throw new Error(t('assignAdminError'));
        }

        await requestApi.createTicket(newTicket.value);
        newTicket.value = { Title: '', Description: '', AssignedTo: null };
        showCreateTicketModal.value = false;
        await fetchTickets();
        await fetchNotifications();

        toast.success(t('ticketCreateSuccess'), {
          icon: "fas fa-check-circle",
          timeout: 3000,
        });
      } catch (error) {
        console.error("Помилка при додаванні тікета:", error);
        toast.error(error.message || t('ticketCreateError'), {
          icon: "fas fa-exclamation-triangle",
          timeout: 5000,
        });
      } finally {
        isLoading.value = false;
      }
    };

    const editTicket = (ticket) => {
      editTicketData.value = { ...ticket };
    };

    const updateTicket = async () => {
      isLoading.value = true;
      try {
        const assignedUser = users.value.find(u => u.UserID === editTicketData.value.AssignedTo);
        if (assignedUser && assignedUser.RoleID === 1) {
          throw new Error(t('assignAdminError'));
        }

        await requestApi.updateTicket(editTicketData.value.TicketID, editTicketData.value);
        editTicketData.value = null;
        await Promise.all([
          fetchTickets(),
          fetchNotifications()
        ]);

        toast.success(t('ticketUpdateSuccess'), {
          icon: "fas fa-check-circle",
          timeout: 3000,
        });
      } catch (error) {
        console.error("Помилка при оновленні тікета:", error);
        toast.error(error.message || t('ticketUpdateError'), {
          icon: "fas fa-exclamation-triangle",
          timeout: 5000,
        });
      } finally {
        isLoading.value = false;
      }
    };

    const markAsRead = async (notificationID) => {
      try {
        await notificationApi.markAsRead(notificationID);
        await fetchNotifications();

        toast.success(t('notificationMarkedRead'), {
          icon: "fas fa-check-circle",
          timeout: 3000,
        });
      } catch (error) {
        console.error("Помилка при видаленні сповіщення:", error);
        toast.error(t('notificationMarkError'), {
          icon: "fas fa-exclamation-triangle",
          timeout: 5000,
        });
      }
    };

    onMounted(() => {
      fetchTickets();
      fetchNotifications();
      fetchUsers();
    });

    return {
      tickets,
      notifications,
      users,
      availableUsers,
      newTicket,
      editTicketData,
      showCreateTicketModal,
      confirmModal,
      authStore,
      isLoading,
      formatDate,
      confirmDelete,
      executeDelete,
      deleteTicket,
      addTicket,
      editTicket,
      updateTicket,
      markAsRead,
      paginatedNotifications,
      nextNotificationsPage,
      prevPage,
      nextPage,
      paginatedTickets,
      currentPage,
      totalPages,
      truncateTitle,
      truncateDescription,
      toggleDescription,
      getStatusClass,
      expandedDescriptions,
      currentNotificationsPage,
      totalNotificationsPages,
      prevNotificationsPage,
      itemsPerPage,

      filteredUsers,
      setLocale,
      currentLocale,
      messageMap
    };
  }
};
</script>

<style scoped>
.dashboard-container {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  min-height: 100vh;
  width: 100%;
  padding: 1rem;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.dashboard-card {
  width: 100%;
  max-width: 1200px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
  padding: 2rem;
  margin: auto;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #eaeaea;
}

.dashboard-header h1 {
  font-size: 1.8rem;
  color: #2c3e50;
  font-weight: 600;
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

.empty-state p {
  color: #718096;
  font-size: 1.05rem;
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

.ticket-description {
  max-width: 300px;
  word-break: break-word;
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

.status-badge {
  display: inline-block;
  border-radius: 50px;
  font-size: 0.8rem;
  font-weight: 500;
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
  gap: 0.4rem;
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

.pagination-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 1.2rem;
  gap: 1rem;
}

.page-info {
  font-size: 0.9rem;
  color: #718096;
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

/* Modal Styles */
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

/* Form Input Groups */
.input-group {
  position: relative;
  margin-bottom: 1.5rem;
}

.input-group label {
  position: absolute;
  top: -8px;
  left: 12px;
  color: #718096;
  font-size: 0.85rem;
  transition: all 0.3s ease;
  pointer-events: none;
  background: white;
  padding: 0 5px;
  z-index: 1;
}

.input-group .form-input {
  width: 100%;
  padding: 14px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.95rem;
  transition: all 0.3s ease;
  background: white;
}

.input-group textarea.form-input {
  min-height: 100px;
  resize: vertical;
}

.input-group select.form-input {
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%234a5568' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 12px center;
  background-size: 14px;
}

.input-group .form-input:focus {
  outline: none;
  border-color: #4361ee;
  box-shadow: 0 0 0 3px rgba(67, 97, 238, 0.1);
}

.input-group .form-input:focus + label {
  color: #4361ee;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.8rem;
  margin-top: 1.5rem;
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

.modal-footer {
  padding: 1.2rem 1.5rem;
  border-top: 1px solid #f0f4f8;
  display: flex;
  justify-content: flex-end;
  gap: 0.8rem;
  background: #f8fafc;
}

/* Transitions */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-active .modal,
.modal-fade-leave-active .modal,
.modal-fade-enter-active .confirm-modal,
.modal-fade-leave-active .confirm-modal {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.modal-fade-enter-from .modal,
.modal-fade-leave-to .modal,
.modal-fade-enter-from .confirm-modal,
.modal-fade-leave-to .confirm-modal {
  transform: translateY(-15px);
  opacity: 0;
}

/* Responsive */
@media (max-width: 768px) {
  .dashboard-container {
    padding: 0;
  }

  .dashboard-card {
    padding: 1rem;
    border-radius: 0;
  }

  .data-table {
    display: block;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    white-space: nowrap;
  }

  .ticket-description {
    max-width: 200px;
    word-break: break-all;
    white-space: normal;
  }

  .notification-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }

  .notification-content {
    width: 100%;
  }

  .action-btn {
    width: 28px;
    height: 28px;
    font-size: 0.8rem;
  }
}

@media (max-width: 480px) {
  .dashboard-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .logout-btn {
    margin-top: 1rem;
    width: 100%;
  }

  .data-table {
    font-size: 0.8rem;
  }

  .data-table th,
  .data-table td {
    padding: 6px 8px;
  }

  .ticket-description {
    max-width: 150px;
  }

  .status-badge {
    font-size: 0.7rem;
    padding: 0.2rem 0.5rem;
  }

  .create-btn {
    bottom: 0.5rem;
    right: 0.5rem;
    padding: 0.7rem 1rem;
    font-size: 0.8rem;
  }

  .modal, .confirm-modal {
    width: 100%;
    max-width: 100%;
    border-radius: 0;
    margin: 0;
  }

  .modal-body {
    padding: 1rem;
  }

  .form-input {
    padding: 12px;
    font-size: 0.9rem;
  }
}

@media (max-width: 360px) {
  .ticket-description {
    max-width: 120px;
  }

  .data-table th,
  .data-table td {
    padding: 4px 6px;
    font-size: 0.75rem;
  }

  .pagination-controls {
    gap: 0.5rem;
    padding: 0.8rem;
  }

  .pagination-btn {
    width: 30px;
    height: 30px;
  }
}

.description-text {
  word-break: break-word;
  overflow-wrap: break-word;
  hyphens: auto;
}

@media (max-width: 768px) {
  .description-text {
    word-break: break-all;
  }
}
</style>