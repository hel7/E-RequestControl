package repository

import (
	"github.com/jmoiron/sqlx"
	"mime/multipart"
	Request_Manager "request_manager_api"
)

type Ticket interface {
	GetUserTickets(userID int) ([]Request_Manager.Ticket, error)
	CreateTicket(userID int, ticket Request_Manager.Ticket) (int, error)
	UpdateTicket(userID int, ticketID int, input Request_Manager.UpdateTicketInput) error
	DeleteUserTicket(ticketID, userID int) error
	DeleteTicketAdmin(ticketID int) error
	GetAllTickets() ([]Request_Manager.Ticket, error)
	GetTicketByID(ticketID int) (Request_Manager.Ticket, error)
}
type Notification interface {
	GetAllUserNotification(userID int) ([]Request_Manager.Notification, error)
	Create(notification Request_Manager.Notification) (int, error)
	Delete(notificationID int) error
	GetAll() ([]Request_Manager.Notification, error)
	MarkNotificationAsRead(notificationID, userID int) error
}
type Authorization interface {
	CreateAdmin(user Request_Manager.User) (int, error)
	CreateUser(user Request_Manager.User) (int, error)
	GetUser(username, password string) (Request_Manager.User, error)
}
type Admin interface {
	GetUserByID(userID int) (Request_Manager.User, error)
	CreateUser(user Request_Manager.User) (int, error)
	GetAllUsers() ([]Request_Manager.User, error)
	Delete(UserID int) error
	UpdateUser(UserID int, input Request_Manager.UpdateUserInput, user Request_Manager.User) error
	GetFilteredTickets(filter Request_Manager.TicketFilter) ([]Request_Manager.Ticket, error)
	BackupData(backupPath string) error
	RestoreData(backupPath multipart.File) error
	ImportData(backupPath string) error
	ExportData(backupPath string) error
}
type Repository struct {
	Ticket
	Notification
	Authorization
	Admin
}

func NewRepository(db *sqlx.DB, cfg Config) *Repository {
	notificationRepo := NewNotificationMysql(db)
	return &Repository{
		Authorization: NewAuthMysql(db),
		Admin:         NewAdminMysql(db, cfg),
		Ticket:        NewTicketMysql(db, notificationRepo),
		Notification:  NewNotificationMysql(db),
	}
}
