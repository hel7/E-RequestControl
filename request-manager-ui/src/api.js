import axios from "axios";

const apiClient = axios.create({
  baseURL: "https://www.zyzel.de/api",
  allowedHosts: ["frontend.local","zyzel.de","https://zyzel.de/api","https://www.zyzel.de/api"],
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");

    if (token) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }

    return config;
  },
  (error) => Promise.reject(error)
);

apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;

    if (
      error.response?.status === 401 &&
      !originalRequest.url.includes("/auth/logout") &&
      !originalRequest.url.includes("/auth/login")
    ) {
      const authStore = useAuthStore();
      await authStore.logout();
    }

    return Promise.reject(error);
  }
);

export const authApi = {
  login: (username, password) =>
    apiClient.post("/auth/login", { username, password }),
  logout: () => apiClient.post("/auth/logout"),
  register: (firstname, lastname, email, username, password) =>
    apiClient.post("/auth/register", {
      firstname,
      lastname,
      email,
      username,
      password,
    }),
  registerAdmin: (firstname, lastname, email, username, password) =>
    apiClient.post("/auth/registerAdmin", {
      firstname,
      lastname,
      email,
      username,
      password,
    }),
};

export const requestApi = {
  getUserTickets: () => apiClient.get("/tickets/"),
  getAllUsers: () => apiClient.get("/users"),
  createTicket: async (ticketData) => {
    try {
      const response = await apiClient.post("/tickets/", ticketData);
      return response.data;
    } catch (error) {
      console.error("Помилка при створенні тікета:", error);
      throw error;
    }
  },
  updateTicket: (ticketID, ticketData) =>
    apiClient.put(`/tickets/${ticketID}`, ticketData),
  deleteTicket: async (ticketID) => {
    try {
      await apiClient.delete(`/tickets/${ticketID}`);
      return ticketID;
    } catch (error) {
      console.error("Помилка при видаленні тікета:", error);
      throw error;
    }
  },
};

export const notificationApi = {
  getUserNotifications: () => apiClient.get("/notifications/"),
  createNotification: async (notificationData) => {
    try {
      const response = await apiClient.post(
        "/notifications/",
        notificationData
      );
      return response.data;
    } catch (error) {
      console.error("Помилка при створенні сповіщення:", error);
      throw error;
    }
  },
  markAsRead: async (notificationID) => {
    try {
      await apiClient.delete(`/notifications/${notificationID}`);
      return notificationID;
    } catch (error) {
      console.error("Помилка при видаленні сповіщення:", error);
      throw error;
    }
  },
};

export const adminApi = {
  getFilteredTickets: (filters) => {
    const params = new URLSearchParams();

    if (filters.sender?.trim()) {
      params.append("sender", filters.sender.trim());
    }
    if (filters.assignee?.trim()) {
      params.append("assignee", filters.assignee.trim());
    }
    if (filters.status?.trim()) {
      params.append("status", filters.status.trim());
    }

    const queryString = params.toString();
    console.log("Request URL:", `/admin/tickets/?${queryString}`);

    if (queryString) {
      return apiClient.get(`/admin/tickets/?${queryString}`);
    }
  },
  getAllTickets: () => apiClient.get("/admin/tickets/"),

  getTicketByID: (ticketID) => apiClient.get(`/admin/tickets/${ticketID}`),
  adminDeleteTicket: (ticketID) =>
    apiClient.delete(`/admin/tickets/${ticketID}`),

  getAllNotifications: () => apiClient.get("/admin/notifications/"),
  createNotification: (notificationData) =>
    apiClient.post("/admin/notifications/", notificationData),
  deleteNotification: (notificationID) =>
    apiClient.delete(`/admin/notifications/${notificationID}`),

  getAllUsers: () => apiClient.get("/admin/users/"),

  createUser: async (userData) => {
    try {
      const response = await apiClient.post("/admin/users/", userData);
      return response.data;
    } catch (error) {
      console.error("Помилка при створенні користувача:", error);
      throw error;
    }
  },
  getUserByID: (userID) => apiClient.get(`/admin/users/${userID}`),
  updateUser: (userID, userData) =>
    apiClient.put(`/admin/users/${userID}`, userData),
  deleteUser: (userID) => apiClient.delete(`/admin/users/${userID}`),

  backupData: () =>
    apiClient.post("/admin/data/backup", null, {
      responseType: "blob",
      headers: {
        Accept: "application/octet-stream",
      },
    }),
  restoreData: (formData) =>
    apiClient.post("/admin/data/restore", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    }),
  exportData: async () => {
    const response = await apiClient.get("/admin/data/export", {
      responseType: "blob",
    });
    return response.data;
  },

  importData: async (importFile) => {
    const formData = new FormData();
    formData.append("file", importFile);
    await apiClient.post("/admin/data/import", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
  },
};
