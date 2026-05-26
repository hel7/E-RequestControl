package repository

import (
	"github.com/jmoiron/sqlx"
	Request_Manager "request_manager_api"
)

type NotificationMysql struct {
	db *sqlx.DB
}

func NewNotificationMysql(db *sqlx.DB) *NotificationMysql {
	return &NotificationMysql{db: db}
}
func (r *NotificationMysql) Create(notification Request_Manager.Notification) (int, error) {
	createdAt := getCurrentTimeInUkraine()
	query := "INSERT INTO Notification (Message, UserID, CreatedAt) VALUES (?, ?, ?)"
	result, err := r.db.Exec(query, notification.Message, notification.UserID, createdAt)
	if err != nil {
		return 0, err
	}
	id64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	id := int(id64)
	return id, nil
}
func (r *NotificationMysql) GetAllUserNotification(userID int) ([]Request_Manager.Notification, error) {
	var userNotification []Request_Manager.Notification
	query := "SELECT * FROM Notification WHERE UserID=?"
	err := r.db.Select(&userNotification, query, userID)
	return userNotification, err
}
func (r *NotificationMysql) GetAll() ([]Request_Manager.Notification, error) {
	var notification []Request_Manager.Notification
	query := "SELECT * FROM Notification"
	err := r.db.Select(&notification, query)
	return notification, err
}

func (r *NotificationMysql) Delete(notificationID int) error {
	query := "DELETE FROM Notification WHERE NotificationID = ?"
	_, err := r.db.Exec(query, notificationID)
	return err
}
func (r *NotificationMysql) MarkNotificationAsRead(notificationID, userID int) error {
	query := "DELETE FROM Notification WHERE NotificationID = ? AND UserID=?"
	_, err := r.db.Exec(query, notificationID, userID)
	return err
}
