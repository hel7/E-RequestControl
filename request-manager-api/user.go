package request_manager_api

import (
	"encoding/json"
	"errors"
	"regexp"
)

type User struct {
	UserID    int    `json:"UserID" db:"UserID"`
	Username  string `json:"Username" db:"Username"`
	Password  string `json:"Password,omitempty" db:"Password"`
	FirstName string `json:"FirstName" db:"FirstName"`
	LastName  string `json:"LastName" db:"LastName"`
	Email     string `json:"Email" db:"Email"`
	RoleID    int    `json:"RoleID" db:"RoleID"`
	CreatedAt string `json:"CreatedAt" db:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt" db:"UpdatedAt"`
}

// MarshalJSON ensures Password is never serialized in API responses
func (u User) MarshalJSON() ([]byte, error) {
	type SafeUser struct {
		UserID    int    `json:"UserID"`
		Username  string `json:"Username"`
		FirstName string `json:"FirstName"`
		LastName  string `json:"LastName"`
		Email     string `json:"Email"`
		RoleID    int    `json:"RoleID"`
		CreatedAt string `json:"CreatedAt"`
		UpdatedAt string `json:"UpdatedAt"`
	}
	return json.Marshal(SafeUser{
		UserID:    u.UserID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		RoleID:    u.RoleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}

type UpdateUserInput struct {
	UserID    *int    `json:"UserID" db:"UserID"`
	Username  *string `json:"Username" db:"Username"`
	Password  *string `json:"Password" db:"Password"`
	FirstName *string `json:"FirstName" db:"FirstName"`
	LastName  *string `json:"LastName" db:"LastName"`
	Email     *string `json:"Email" db:"Email"`
	RoleID    *int    `json:"RoleID" db:"RoleID"`
	CreatedAt *string `json:"CreatedAt" db:"CreatedAt"`
	UpdatedAt *string `json:"UpdatedAt" db:"UpdatedAt"`
}

func (i UpdateUserInput) Validate() error {
	if i.UserID == nil && i.Username == nil && i.Email == nil && i.Password == nil && i.RoleID == nil &&
		i.FirstName == nil && i.LastName == nil && i.CreatedAt == nil && i.UpdatedAt == nil {
		return errors.New("Update structure has no value")
	}
	return nil
}
func (u *User) ValidateEmail() error {
	if u.Email != "" {
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		match, _ := regexp.MatchString(emailRegex, u.Email)
		if !match {
			return errors.New("invalid email format")
		}
	}
	return nil
}
func (u *User) ValidatePassword() error {
	if u.Password != "" {
		if len(u.Password) < 8 {
			return errors.New("password should be at least 8 characters long")
		}

		letterRegex := regexp.MustCompile(`[a-zA-Z]`)
		if !letterRegex.MatchString(u.Password) {
			return errors.New("password should contain at least 1 letter")
		}

		digitRegex := regexp.MustCompile(`\d`)
		if !digitRegex.MatchString(u.Password) {
			return errors.New("password should contain at least 1 digit")
		}

		specialCharRegex := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
		if !specialCharRegex.MatchString(u.Password) {
			return errors.New("password should contain at least 1 special character")
		}
	}
	return nil
}
