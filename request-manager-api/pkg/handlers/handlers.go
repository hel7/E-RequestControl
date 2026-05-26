package handlers

import (
	"net/http"
	"request_manager_api/pkg/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *services.Service
}

func NewHandler(services *services.Service) *Handlers {
	return &Handlers{service: services}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowedOrigins := []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"127.0.0.1",
			"http://127.0.0.1:",
			"http://frontend.local",
			"http://localhost:5173/api",
			"http://req-front-service",
			"http://req-front-service.default.svc.cluster.local",
			"https://zyzel.de",
			"https://zyzel.de/api",
			"https://www.zyzel.de/api",
			"https://www.zyzel.de",
			"http://www.zyzel.de/api",
			"http://www.zyzel.de",
		}

		if contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (h *Handlers) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware())

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.register)
			auth.POST("/registerAdmin", h.registerAdmin)
			auth.POST("/login", h.login)
			auth.POST("/logout", h.logout)
		}

		api.Use(h.userIdentity)
		api.GET("/users", h.getAllUsers)
		tickets := api.Group("/tickets")
		{
			tickets.GET("/", h.getUserTickets)
			tickets.POST("/", h.createTicket)
			tickets.PUT("/:ticketID", h.updateTicket)
			tickets.DELETE("/:ticketID", h.deleteTicket)
		}
		notifications := api.Group("/notifications")
		{
			notifications.GET("/", h.getUserNotifications)
			notifications.DELETE("/:notificationID", h.markNotificationAsRead)
		}
	}
		admin := router.Group("api/admin", h.userIdentity, h.adminRequired)
		{
			adminTickets := admin.Group("/tickets")
			{
				adminTickets.GET("/", h.getTickets)
				adminTickets.GET("/:ticketID", h.getTicketByID)
				adminTickets.DELETE("/:ticketID", h.adminDeleteTicket)
			}
			adminNotifications := admin.Group("/notifications")
			{
				adminNotifications.GET("/", h.getAllNotifications)
				adminNotifications.POST("/", h.createNotification)
				adminNotifications.DELETE("/:notificationID", h.deleteNotification)
			}

			users := admin.Group("/users")
			{
				users.GET("/", h.getAllUsers)
				users.POST("/", h.createUser)
				users.GET("/:userID", h.getUserByID)
				users.PUT("/:userID", h.updateUser)
				users.DELETE("/:userID", h.deleteUser)
			}

			data := admin.Group("/data")
			{
				data.POST("/backup", h.backupData)
				data.POST("/restore", h.restoreData)
				data.GET("/export", h.exportData)
				data.POST("/import", h.importData)
			}
		}

	return router
}

