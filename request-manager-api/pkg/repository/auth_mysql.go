package repository

import (
	"github.com/jmoiron/sqlx"
	"log"
	Request_Manager "request_manager_api"
	"time"
)

type AuthMysql struct {
	db *sqlx.DB
}

func NewAuthMysql(db *sqlx.DB) *AuthMysql {
	return &AuthMysql{db: db}
}

func getCurrentTimeInUkraine() string {
	ukraineLocation, err := time.LoadLocation("Europe/Kiev")
	if err != nil {
		log.Fatal("Error loading Ukraine timezone: ", err)
	}
	return time.Now().In(ukraineLocation).Format("2006-01-02 15:04:05")
}

func (r *AuthMysql) CreateUser(user Request_Manager.User) (int, error) {

	currentTime := getCurrentTimeInUkraine()

	query := "INSERT INTO User (username, email, password, roleID, firstName, lastName, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, 2, user.FirstName, user.LastName, currentTime, currentTime)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *AuthMysql) CreateAdmin(user Request_Manager.User) (int, error) {
	currentTime := getCurrentTimeInUkraine()

	query := "INSERT INTO User (username, email, password, roleID, firstName, lastName, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, 1, user.FirstName, user.LastName, currentTime, currentTime)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *AuthMysql) GetUser(username, password string) (Request_Manager.User, error) {
	var user Request_Manager.User
	query := "SELECT UserID, RoleID FROM User WHERE Username=? AND PASSWORD=?"
	err := r.db.Get(&user, query, username, password)
	return user, err
}
