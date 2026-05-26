package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"path/filepath"
	Request_Manager "request_manager_api"
	"strconv"
	"time"
)

func (h *Handlers) getTickets(c *gin.Context) {
	SenderUsername := c.Query("sender")
	AssigneeUsername := c.Query("assignee")
	Status := c.Query("status")

	if SenderUsername != "" || AssigneeUsername != "" || Status != "" {
		filter := Request_Manager.TicketFilter{
			SenderUsername:   SenderUsername,
			AssigneeUsername: AssigneeUsername,
			Status:           Status,
		}

		tickets, err := h.service.Admin.GetFilteredTickets(filter)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, tickets)
		return
	}

	tickets, err := h.service.Ticket.GetAllTickets()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func (h *Handlers) getTicketByID(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("ticketID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ticket ID"})
		return
	}
	company, err := h.service.Ticket.GetTicketByID(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, company)
}
func (h *Handlers) adminDeleteTicket(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("ticketID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ticket ID"})
		return
	}
	if err := h.service.Ticket.DeleteTicketAdmin(ticketID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted"})
}

func (h *Handlers) getUserTickets(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}
	tickets, err := h.service.Ticket.GetUserTickets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}
func (h *Handlers) createTicket(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	var ticket Request_Manager.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket.UserID = userID

	id, err := h.service.Ticket.CreateTicket(userID, ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handlers) updateTicket(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ticketID := c.Param("ticketID")
	id, err := strconv.Atoi(ticketID)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid id param")
		return
	}

	var input Request_Manager.UpdateTicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Ticket.UpdateTicket(userID, id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

func (h *Handlers) deleteTicket(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	ticketID, err := strconv.Atoi(c.Param("ticketID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ticket ID"})
		return
	}

	if err := h.service.Ticket.DeleteUserTicket(ticketID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted"})
}
func (h *Handlers) getUserNotifications(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	notifications, err := h.service.Notification.GetAllUserNotification(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
func (h *Handlers) markNotificationAsRead(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	notificationID, err := strconv.Atoi(c.Param("notificationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.service.Notification.MarkNotificationAsRead(notificationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (h *Handlers) getAllUsers(c *gin.Context) {
	users, err := h.service.Admin.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
func (h *Handlers) createUser(c *gin.Context) {
	var user Request_Manager.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	id, err := h.service.Admin.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"UserID": id,
	})
}
func (h *Handlers) getUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handlers) updateUser(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid UserID")
		return
	}
	var input Request_Manager.UpdateUserInput
	var userInput Request_Manager.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.service.Admin.UpdateUser(userID, input, userInput)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

func (h *Handlers) deleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}
	if err := h.service.Admin.Delete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func (h *Handlers) createNotification(c *gin.Context) {
	var notification Request_Manager.Notification
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	id, err := h.service.Notification.Create(notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
func (h *Handlers) getAllNotifications(c *gin.Context) {
	notification, err := h.service.Notification.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notification)
}
func (h *Handlers) deleteNotification(c *gin.Context) {
	notificationID, err := strconv.Atoi(c.Param("notificationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.service.Notification.Delete(notificationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}
func (h *Handlers) backupData(c *gin.Context) {
	backupFile := "/root/backup_" + time.Now().Format("20060102_150405") + ".sql"

	if err := os.MkdirAll("/root", os.ModePerm); err != nil {
		log.Printf("Failed to create backup directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create backup directory",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.Admin.BackupData(backupFile); err != nil {
		log.Printf("Backup failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate backup",
			"details": err.Error(),
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(backupFile))
	c.Header("Content-Type", "application/octet-stream")
	c.File(backupFile)

	go func() {
		time.Sleep(2 * time.Second)
		if err := os.Remove(backupFile); err != nil {
			log.Printf("Failed to remove backup file: %v", err)
		}
	}()
}

func (h *Handlers) restoreData(c *gin.Context) {
	logrus.Info("Starting data restore operation")

	file, err := c.FormFile("file")
	if err != nil {
		logrus.Errorf("No file uploaded: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	logrus.Infof("Restoring from file: %s (size: %d bytes)", file.Filename, file.Size)

	f, err := file.Open()
	if err != nil {
		logrus.Errorf("Failed to open uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	if err := h.service.Admin.RestoreData(f); err != nil {
		logrus.Errorf("Data restore failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to restore data",
			"details": err.Error(),
		})
		return
	}

	logrus.Info("Data restored successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Data restored successfully"})
}

func (h *Handlers) exportData(c *gin.Context) {
	exportPath := "export.xlsx"

	if err := h.service.Admin.ExportData(exportPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export data"})
		return
	}
	c.File(exportPath)
	c.JSON(http.StatusOK, gin.H{"message": "Data export successful"})
}
func (h *Handlers) importData(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadPath := file.Filename
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	if err := h.service.Admin.ImportData(uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully"})
}
