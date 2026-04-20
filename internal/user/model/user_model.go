package model

import (
	"errors"
	"net/mail"
	"strings"
)

// User represents a user record stored in the database.
type User struct {
	ID       int64  `db:"id" json:"id"`
	UserName string `db:"user_name" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

func (u *User) Validate() error {
	u.UserName = strings.TrimSpace(u.UserName)
	u.Email = strings.TrimSpace(u.Email)

	if u.UserName == "" {
		return errors.New("username is required")
	}

	if len(u.UserName) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if strings.Contains(u.UserName, " ") {
		return errors.New("username cannot contain spaces")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	// simple email validation
	err := validateEmail(u.Email)
	if err != nil {
		return err
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	return nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}
	return nil
}
