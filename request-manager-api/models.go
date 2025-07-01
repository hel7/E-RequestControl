package request_manager_api

import "errors"

type Ticket struct {
	TicketID         int    `json:"TicketID" db:"TicketID"`
	Title            string `json:"Title" binding:"required" db:"Title"`
	Description      string `json:"Description" db:"Description"`
	Status           string `json:"Status" db:"Status"`
	CreatedAt        string `json:"CreatedAt" db:"CreatedAt"`
	UpdatedAt        string `json:"UpdatedAt" db:"UpdatedAt"`
	StatusID         string `json:"-" db:"StatusID"`
	AssignedTo       int    `json:"AssignedTo" db:"AssignedTo"`
	UserID           int    `json:"-" db:"UserID"`
	SenderUsername   string `json:"SenderUsername" db:"SenderUsername"`
	AssigneeUsername string `json:"AssigneeUsername" db:"AssigneeUsername"`
}
type TicketFilter struct {
	Status           string `json:"Status" db:"Status"`
	SenderUsername   string `json:"SenderUsername" db:"SenderUsername"`
	AssigneeUsername string `json:"AssigneeUsername" db:"AssigneeUsername"`
}

type UpdateTicketInput struct {
	TicketID    *int    `json:"TicketID" db:"TicketID"`
	Title       *string `json:"Title" binding:"required" db:"Title"`
	Description *string `json:"Description" db:"Description"`
	StatusID    *int    `json:"StatusID" db:"StatusID"`
	CreatedAt   *string `json:"CreatedAt" db:"CreatedAt"`
	UpdatedAt   *string `json:"UpdatedAt" db:"UpdatedAt"`
	AssignedTo  *int    `json:"AssignedTo" db:"AssignedTo"`
	UserID      *int    `json:"UserID" db:"UserID"`
}

type TicketStatus struct {
	StatusID  int    `json:"StatusID" db:"StatusID"`
	Status    string `json:"Status" db:"Status"`
	CreatedAt string `json:"CreatedAt" db:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt" db:"UpdatedAt"`
}

type Role struct {
	RoleID    int    `json:"RoleID" db:"RoleID"`
	RoleName  string `json:"RoleName" db:"RoleName"`
	CreatedAt string `json:"CreatedAt" db:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt" db:"UpdatedAt"`
}

type Notification struct {
	NotificationID int    `json:"NotificationID" db:"NotificationID"`
	Message        string `json:"Message" db:"Message"`
	UserID         int    `json:"UserID" db:"UserID"`
	CreatedAt      string `json:"CreatedAt" db:"CreatedAt"`
}

func (i UpdateTicketInput) Validate() error {
	if i.TicketID == nil && i.Title == nil && i.Description == nil && i.StatusID == nil && i.CreatedAt == nil &&
		i.UpdatedAt == nil && i.AssignedTo == nil && i.UserID == nil {
		return errors.New("Update structure has no value")
	}
	return nil
}
