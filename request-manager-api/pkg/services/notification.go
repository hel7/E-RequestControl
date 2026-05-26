package services

import (
	Request_Manager "request_manager_api"
	"request_manager_api/pkg/repository"
)

type NotificationService struct {
	repo repository.Notification
}

func NewNotificationService(repo repository.Notification) *NotificationService {
	return &NotificationService{repo: repo}
}
func (s *NotificationService) Create(notification Request_Manager.Notification) (int, error) {
	return s.repo.Create(notification)
}
func (s *NotificationService) GetAllUserNotification(userID int) ([]Request_Manager.Notification, error) {
	return s.repo.GetAllUserNotification(userID)
}
func (s *NotificationService) GetAll() ([]Request_Manager.Notification, error) {
	return s.repo.GetAll()
}
func (s *NotificationService) Delete(notificationID int) error {
	return s.repo.Delete(notificationID)
}
func (s *NotificationService) MarkNotificationAsRead(notificationID, userID int) error {
	return s.repo.MarkNotificationAsRead(notificationID, userID)
}
