package repository

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	Request_Manager "request_manager_api"
	"strings"
)

type AdminMysql struct {
	db *sqlx.DB
	config Config
}

func NewAdminMysql(db *sqlx.DB, cfg Config) *AdminMysql {
	return &AdminMysql{
	    db: db,
	    config: cfg,
	}
}

func (r *AdminMysql) GetUserByID(userID int) (Request_Manager.User, error) {
	var user Request_Manager.User
	query := `SELECT * FROM User WHERE UserID = ?`
	err := r.db.Get(&user, query, userID)
	if err != nil {
		log.Printf("Error fetching user by ID %d: %s", userID, err)
	}
	return user, err
}
func (r *AdminMysql) GetAllUsers() ([]Request_Manager.User, error) {
	var users []Request_Manager.User
	query := fmt.Sprintf("SELECT * FROM User")
	err := r.db.Select(&users, query)
	if err != nil {
		log.Printf("Error fetching all users: %s", err)
	}
	return users, err
}
func (r *AdminMysql) Delete(UserID int) error {
	query := fmt.Sprintf("DELETE FROM User WHERE UserID = ? ")
	_, err := r.db.Exec(query, UserID)
	return err
}

func (r *AdminMysql) CreateUser(user Request_Manager.User) (int, error) {
	checkQuery := fmt.Sprintf("SELECT COUNT(*) FROM User WHERE email=? OR username=?")
	var count int
	err := r.db.Get(&count, checkQuery, user.Email, user.Username)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, fmt.Errorf("user with this email or username already exists")
	}

	createdAt := getCurrentTimeInUkraine()
	updatedAt := createdAt

	query := fmt.Sprintf(`INSERT INTO User (Username, Email, Password, RoleID, FirstName, LastName, CreatedAt, UpdatedAt)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)

	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.RoleID, user.FirstName, user.LastName, createdAt, updatedAt)

	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok && mysqlErr.Number == 1062 {
			return 0, fmt.Errorf("user with this email or username already exists")
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *AdminMysql) UpdateUser(UserID int, input Request_Manager.UpdateUserInput, user Request_Manager.User) error {
	existingUser, err := r.GetUserByID(UserID)
	if err != nil {
		return err
	}

	updatedAt := getCurrentTimeInUkraine()

	if input.Email != nil {
		checkEmailQuery := "SELECT UserID FROM User WHERE Email=? AND UserID != ?"
		var otherUserID int
		err := r.db.Get(&otherUserID, checkEmailQuery, *input.Email, UserID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if otherUserID != 0 {
			return errors.New("user with this email already exists")
		}
		existingUser.Email = *input.Email
	}

	if input.Password != nil {
		existingUser.Password = *input.Password
	}

	if input.Username != nil {
		checkUsernameQuery := "SELECT UserID FROM User WHERE Username=? AND UserID != ?"
		var otherUserID int
		err := r.db.Get(&otherUserID, checkUsernameQuery, *input.Username, UserID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if otherUserID != 0 {
			return errors.New("user with this username already exists")
		}
		existingUser.Username = *input.Username
	}

	if input.RoleID != nil {
		existingUser.RoleID = *input.RoleID
	}
	if input.FirstName != nil {
		existingUser.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		existingUser.LastName = *input.LastName
	}

	query := `UPDATE User
	          SET Username=?, Email=?, Password=?, RoleID=?, FirstName=?, LastName=?, UpdatedAt=?
	          WHERE UserID=?`
	_, err = r.db.Exec(query, existingUser.Username, existingUser.Email, existingUser.Password, existingUser.RoleID, existingUser.FirstName, existingUser.LastName, updatedAt, UserID)
	if err != nil {
		return err
	}

	return nil
}
func (r *AdminMysql) GetFilteredTickets(filter Request_Manager.TicketFilter) ([]Request_Manager.Ticket, error) {
	var tickets []Request_Manager.Ticket

	query := `
        SELECT 
            t.TicketID, 
            t.Title, 
            t.Description, 
            ts.Status,
            t.CreatedAt,
            t.UpdatedAt,
            t.StatusID,
            t.AssignedTo,
            t.UserID,
            u.Username AS SenderUsername,
            a.Username AS AssigneeUsername
        FROM Ticket t
        JOIN TicketStatus ts ON t.StatusID = ts.StatusID
        JOIN User u ON t.UserID = u.UserID
        LEFT JOIN User a ON t.AssignedTo = a.UserID
        WHERE 1=1`

	var args []interface{}

	if filter.SenderUsername != "" {
		query += " AND u.Username LIKE ?"
		args = append(args, "%"+filter.SenderUsername+"%")
	}

	if filter.AssigneeUsername == "unassigned" {
		query += " AND t.AssignedTo IS NULL"
	} else if filter.AssigneeUsername != "" {
		query += " AND a.Username LIKE ?"
		args = append(args, "%"+filter.AssigneeUsername+"%")
	}

	if filter.Status != "" {
		query += " AND ts.Status = ? COLLATE utf8mb4_unicode_ci"
		args = append(args, filter.Status)
	}

	query += " ORDER BY t.CreatedAt DESC"

	err := r.db.Select(&tickets, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("database error: %v", err)
	}

	return tickets, nil
}

func (r *AdminMysql) BackupData(backupPath string) error {
	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	outputFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer outputFile.Close()

	cmd := exec.Command(
    	"mysqldump",
    	"--host="+r.config.Host,
    	"--port="+r.config.Port,
    	"--user="+r.config.Username,
    	"--password="+r.config.Password,
    	"--skip-ssl",
    	r.config.Dbname,
    )

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = outputFile

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("backup failed: %v, details: %s", err, stderr.String())
	}

	if info, err := os.Stat(backupPath); err != nil || info.Size() == 0 {
		return fmt.Errorf("backup file is empty or not created")
	}

	return nil
}

func (r *AdminMysql) RestoreData(backupFile multipart.File) error {
	tempFile, err := os.CreateTemp("", "restore-*.sql")
	if err != nil {
		logrus.Errorf("Failed to create temp file: %v", err)
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, backupFile); err != nil {
		logrus.Errorf("Failed to save backup file: %v", err)
		return fmt.Errorf("failed to save backup file: %v", err)
	}

	logrus.Infof("Starting database restore from temporary file: %s", tempFile.Name())

	dumpFile, err := os.Open(tempFile.Name())
	if err != nil {
		logrus.Errorf("Failed to open dump file: %v", err)
		return fmt.Errorf("failed to open dump file: %v", err)
	}
	defer dumpFile.Close()

	cmd := exec.Command(
    	"mysql",
    	"-h", r.config.Host,
    	"-P", r.config.Port,
    	"-u", r.config.Username,
    	"-p"+r.config.Password,
    	"--skip-ssl",
    	r.config.Dbname,
    )

	cmd.Stdin = dumpFile
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	logrus.Info("Executing database restore command")
	if err := cmd.Run(); err != nil {
		errorDetails := stderr.String()
		logrus.Errorf("Restore command failed: %v, details: %s", err, errorDetails)
		return fmt.Errorf("restore failed: %v, details: %s", err, errorDetails)
	}

	logrus.Info("Database restore completed successfully")
	return nil
}

func (r *AdminMysql) ExportData(exportPath string) error {
	logger := logrus.New()

	file := xlsx.NewFile()

	tables := map[string]string{
		"User":         "SELECT * FROM User;",
		"Role":         "SELECT * FROM Role;",
		"TicketStatus": "SELECT * FROM TicketStatus;",
		"Ticket":       "SELECT * FROM Ticket;",
		"Notification": "SELECT * FROM Notification;",
	}

	for tableName, query := range tables {
		rows, err := r.db.Query(query)
		if err != nil {
			logger.Errorf("Error querying table %s: %s", tableName, err)
			return fmt.Errorf("error querying table %s: %w", tableName, err)
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			logger.Errorf("Error fetching columns for table %s: %s", tableName, err)
			return fmt.Errorf("error fetching columns for table %s: %w", tableName, err)
		}

		sheet, err := file.AddSheet(tableName)
		if err != nil {
			logger.Errorf("Error adding sheet %s: %s", tableName, err)
			return fmt.Errorf("error adding sheet %s: %w", tableName, err)
		}

		headerRow := sheet.AddRow()
		for _, column := range columns {
			cell := headerRow.AddCell()
			cell.Value = column
		}

		for rows.Next() {
			rowData := make([]sql.NullString, len(columns))
			valuePointers := make([]interface{}, len(columns))
			for i := range rowData {
				valuePointers[i] = &rowData[i]
			}

			err := rows.Scan(valuePointers...)
			if err != nil {
				logger.Errorf("Error scanning rows for table %s: %s", tableName, err)
				return fmt.Errorf("error scanning rows for table %s: %w", tableName, err)
			}

			row := sheet.AddRow()
			for _, value := range rowData {
				cell := row.AddCell()
				if value.Valid {
					cell.Value = value.String
				} else {
					cell.Value = "NULL"
				}
			}
		}
	}

	err := file.Save(exportPath)
	if err != nil {
		logger.Errorf("Error saving Excel file: %s", err)
		return fmt.Errorf("error saving Excel file: %w", err)
	}

	return nil
}

func (r *AdminMysql) ImportData(importPath string) error {
	file, err := xlsx.OpenFile(importPath)
	if err != nil {
		logrus.Errorf("Error opening Excel file: %s", err)
		return fmt.Errorf("error opening Excel file: %w", err)
	}

	for _, sheet := range file.Sheets {
		tableName := sheet.Name

		rows := sheet.Rows
		if len(rows) < 2 {
			continue
		}

		columns := make([]string, len(rows[0].Cells))
		for i, cell := range rows[0].Cells {
			columns[i] = cell.String()
		}

		query := fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES ", tableName, strings.Join(columns, ","))

		var valueStrings []string
		var valueArgs []interface{}

		for _, row := range rows[1:] {
			var values []interface{}

			for _, cell := range row.Cells {
				values = append(values, cell.Value)
			}

			placeholders := make([]string, len(columns))
			for i := range placeholders {
				placeholders[i] = "?"
			}

			valueStrings = append(valueStrings, "("+strings.Join(placeholders, ",")+")")
			valueArgs = append(valueArgs, values...)
		}

		query += strings.Join(valueStrings, ",")
		_, err := r.db.Exec(query, valueArgs...)
		if err != nil {
			logrus.Errorf("Error inserting data into table %s: %s", tableName, err)
			return fmt.Errorf("error inserting data into table %s: %w", tableName, err)
		}
	}

	return nil
}
