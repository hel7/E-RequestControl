package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	Request_Manager "request_manager_api"
)

type TicketMysql struct {
	db           *sqlx.DB
	notification Notification
}

func NewTicketMysql(db *sqlx.DB, notification Notification) *TicketMysql {
	return &TicketMysql{db: db, notification: notification}
}

func (r *TicketMysql) GetUserTickets(userID int) ([]Request_Manager.Ticket, error) {
	var tickets []Request_Manager.Ticket
	query := `SELECT 
        t.TicketID, 
        t.Title, 
        t.Description, 
        ts.Status, 
        t.CreatedAt, 
        t.UpdatedAt, 
        t.AssignedTo, 
        sender.Username AS SenderUsername, 
        assignee.Username AS AssigneeUsername
    FROM Ticket t
    JOIN TicketStatus ts ON t.StatusID = ts.StatusID
    JOIN User sender ON t.UserID = sender.UserID
    LEFT JOIN User assignee ON t.AssignedTo = assignee.UserID
    WHERE t.UserID = ?`

	err := r.db.Select(&tickets, query, userID)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketMysql) CreateTicket(userID int, ticket Request_Manager.Ticket) (int, error) {
	createdAt := getCurrentTimeInUkraine()

	if ticket.AssignedTo == 0 {
		return 0, errors.New("AssignedTo не може бути пустим")
	}

	var assignedTo int
	err := r.db.QueryRow("SELECT UserID FROM User WHERE UserID = ?", ticket.AssignedTo).Scan(&assignedTo)
	if err != nil {
		return 0, errors.New("Обраного користувача не існує")
	}

	var statusID int
	err = r.db.QueryRow("SELECT StatusID FROM TicketStatus WHERE Status = ?", "Новий").Scan(&statusID)

	if err == sql.ErrNoRows {
		_, err := r.db.Exec("INSERT INTO TicketStatus (Status) VALUES (?)", "Новий")
		if err != nil {
			return 0, err
		}

		err = r.db.QueryRow("SELECT StatusID FROM TicketStatus WHERE Status = ?", "Новий").Scan(&statusID)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO Ticket (Title, Description, StatusID, CreatedAt, UpdatedAt, AssignedTo, UserID)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.Exec(query, ticket.Title, ticket.Description, statusID, createdAt, createdAt, ticket.AssignedTo, ticket.UserID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	notificationMessage := "Створено новий тікет"

	notification := Request_Manager.Notification{
		Message:   notificationMessage,
		UserID:    userID,
		CreatedAt: getCurrentTimeInUkraine(),
	}

	_, err = r.notification.Create(notification)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *TicketMysql) UpdateTicket(userID int, ticketID int, input Request_Manager.UpdateTicketInput) error {
	updatedAt := getCurrentTimeInUkraine()

	existingTicket, err := r.GetByID(ticketID)
	if err != nil {
		return err
	}

	if input.Title != nil {
		existingTicket.Title = *input.Title
	}
	if input.Description != nil {
		existingTicket.Description = *input.Description
	}
	if input.AssignedTo != nil {
		existingTicket.AssignedTo = *input.AssignedTo
	}

	var statusID int
	err = r.db.Get(&statusID, `SELECT StatusID FROM TicketStatus WHERE Status = ?`, "Оновлено")
	if err != nil {
		result, err := r.db.Exec(`INSERT INTO TicketStatus (Status) VALUES (?)`, "Оновлено")
		if err != nil {
			return err
		}
		statusID64, err := result.LastInsertId()
		if err != nil {
			return err
		}
		statusID = int(statusID64)
	}

	query := `UPDATE Ticket SET Title=?, Description=?, AssignedTo=?, StatusID=?, UpdatedAt=? WHERE TicketID=?`
	_, err = r.db.Exec(query, existingTicket.Title, existingTicket.Description, existingTicket.AssignedTo, statusID, updatedAt, ticketID)

	notificationMessage := "Оновлено тікет"

	notification := Request_Manager.Notification{
		Message:   notificationMessage,
		UserID:    userID,
		CreatedAt: getCurrentTimeInUkraine(),
	}

	_, err = r.notification.Create(notification)
	if err != nil {
		return err
	}

	return err
}

func (r *TicketMysql) GetByID(ticketID int) (Request_Manager.Ticket, error) {
	var ticket Request_Manager.Ticket
	query := `SELECT t.TicketID, t.Title, t.Description, ts.Status, t.AssignedTo, t.UserID, t.CreatedAt, t.UpdatedAt
	          FROM Ticket t 
	          JOIN TicketStatus ts ON t.StatusID = ts.StatusID
	          WHERE t.TicketID = ?`
	err := r.db.Get(&ticket, query, ticketID)
	if err != nil {
		return ticket, err
	}
	return ticket, nil
}

func (r *TicketMysql) DeleteUserTicket(ticketID, userID int) error {
	query := `DELETE FROM Ticket WHERE TicketID=? AND UserID=?`
	result, err := r.db.Exec(query, ticketID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("ticket not found or does not belong to user")
	}

	return nil
}

func (r *TicketMysql) DeleteTicketAdmin(ticketID int) error {
	query := `DELETE FROM Ticket WHERE TicketID=?`
	_, err := r.db.Exec(query, ticketID)
	return err
}
func (r *TicketMysql) GetAllTickets() ([]Request_Manager.Ticket, error) {
	var tickets []Request_Manager.Ticket
	query := `
		SELECT t.TicketID, t.Title, t.Description, ts.Status, t.CreatedAt, t.UpdatedAt,
		       t.AssignedTo, t.UserID, u.Username AS SenderUsername, a.Username AS AssigneeUsername
		FROM Ticket t
		JOIN TicketStatus ts ON t.StatusID = ts.StatusID
		JOIN User u ON t.UserID = u.UserID
		LEFT JOIN User a ON t.AssignedTo = a.UserID`

	err := r.db.Select(&tickets, query)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketMysql) GetTicketByID(ticketID int) (Request_Manager.Ticket, error) {
	var ticket Request_Manager.Ticket
	query := `
		SELECT t.TicketID, t.Title, t.Description, ts.Status, t.CreatedAt, t.UpdatedAt,
		       t.AssignedTo, t.UserID, u.Username AS SenderUsername, a.Username AS AssigneeUsername
		FROM Ticket t
		JOIN TicketStatus ts ON t.StatusID = ts.StatusID
		JOIN User u ON t.UserID = u.UserID
		LEFT JOIN User a ON t.AssignedTo = a.UserID
		WHERE TicketID=?`

	err := r.db.Get(&ticket, query, ticketID)
	return ticket, err
}
